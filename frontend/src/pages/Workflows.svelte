<script lang="ts">
    import { onMount } from "svelte";
    import WorkflowEditor from "../lib/components/WorkflowEditor.svelte";
    import Modal from "../lib/components/Modal.svelte";

    interface Workflow {
        id: number;
        site_id: number;
        site_name: string;
        name: string;
        trigger_event?: string;
        status: string;
        created_at: string;
    }

    interface Site {
        id: number;
        name: string;
    }

    interface EventDefinition {
        id: number;
        name: string;
    }

    interface Audience {
        id: number;
        name: string;
    }

    let workflows = $state<Workflow[]>([]);
    let sites = $state<Site[]>([]);
    let eventDefinitions = $state<EventDefinition[]>([]);
    let audiences = $state<Audience[]>([]);
    let loading = $state(true);

    // Create Modal State
    let showCreateModal = $state(false);
    let newWorkflowName = $state("");
    let selectedSiteId = $state("");
    let selectedTrigger = $state("");
    let triggerType = $state("EVENT"); // 'EVENT' or 'SCHEDULE'
    let schedule = $state("");
    let selectedAudienceId = $state("");

    // Graph State
    let graphNodes = $state<any[]>([]);
    let graphEdges = $state<any[]>([]);
    let selectedNodeId = $state<string | null>(null);

    // Helper to get selected node
    let selectedNode = $derived(
        selectedNodeId ? graphNodes.find((n) => n.id === selectedNodeId) : null,
    );

    async function loadWorkflows() {
        try {
            const res = await fetch("/api/workflows");
            if (res.ok) {
                workflows = await res.json();
            }
        } catch (e) {
            console.error("Failed to load workflows", e);
        } finally {
            loading = false;
        }
    }

    async function loadSites() {
        try {
            const res = await fetch("/api/sites?organization_id=1"); // Hardcoded org for now
            if (res.ok) {
                sites = await res.json();
            }
        } catch (e) {
            console.error("Failed to load sites", e);
        }
    }

    async function fetchEvents(siteId: string) {
        if (!siteId) {
            eventDefinitions = [];
            return;
        }
        try {
            const res = await fetch(
                `/api/events/definitions?site_id=${siteId}`,
            );
            if (res.ok) {
                eventDefinitions = await res.json();
            }
        } catch (e) {
            console.error("Failed to load events", e);
        }
    }

    function handleSiteChange() {
        selectedTrigger = ""; // Reset trigger
        fetchEvents(selectedSiteId);
        // Reset/Recreate Start Node based on Trigger?
        // For simplicity, let's keep graph manual but maybe auto-add Start Node
        initializeGraph();
    }

    function initializeGraph() {
        graphNodes = [];
        graphEdges = [];
        addNode("TRIGGER", "Start Trigger");
        // Fix position of first node
        if (graphNodes[0]) {
            graphNodes[0].position = { x: 50, y: 300 };
        }
    }

    function addNode(type: string, label: string = "") {
        const id = `n-${Date.now()}-${Math.floor(Math.random() * 1000)}`;

        let pos = { x: 100, y: 100 };
        // Simple auto-placement
        if (graphNodes.length > 0) {
            const lastNode = graphNodes[graphNodes.length - 1];
            pos = { x: lastNode.position.x + 250, y: lastNode.position.y };
        }

        const newNode = {
            id,
            type,
            position: pos,
            data: {
                label: label || type,
                action: "",
                force: "", // For condition nodes
                handlePosition: "right",
            },
        };

        graphNodes = [...graphNodes, newNode];

        // Auto-connect not strictly required for this step but good UX
        if (graphNodes.length > 1) {
            const sourceNode = graphNodes[graphNodes.length - 2];
            graphEdges = [
                ...graphEdges,
                {
                    id: `e-${Date.now()}`,
                    source: sourceNode.id,
                    target: id,
                    // Default handle logic could normally go here
                },
            ];
        }
    }

    async function createWorkflow() {
        if (!newWorkflowName) {
            alert("Please enter a workflow name");
            return;
        }

        const workflowData: any = {
            name: newWorkflowName,
            trigger_type: triggerType,
            steps: JSON.stringify({ nodes: graphNodes, edges: graphEdges }), // Graph JSON
            status: "ACTIVE",
            organization_id: 1,
        };

        if (triggerType === "EVENT") {
            if (!selectedSiteId || !selectedTrigger) {
                alert("Please select a site and trigger event");
                return;
            }
            workflowData.site_id = parseInt(selectedSiteId);
            workflowData.trigger_event = selectedTrigger;
        } else if (triggerType === "SCHEDULE") {
            if (!schedule) {
                alert("Please enter a schedule (Cron expression)");
                return;
            }
            workflowData.schedule = schedule;
            if (selectedAudienceId) {
                workflowData.audience_id = parseInt(selectedAudienceId);
            }
        }

        try {
            const res = await fetch("/api/workflows", {
                // Relative path OK
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(workflowData),
            });

            if (res.ok) {
                showCreateModal = false;
                newWorkflowName = "";
                selectedSiteId = "";
                selectedTrigger = "";
                schedule = "";
                selectedAudienceId = "";
                triggerType = "EVENT";
                initializeGraph();
                loadWorkflows();
            } else {
                alert("Failed to create workflow");
            }
        } catch (error) {
            console.error("Error creating workflow", error);
            alert("Error creating workflow");
        }
    }

    onMount(() => {
        loadWorkflows();
        loadSites();
        initializeGraph();
    });
</script>

<div class="max-w-6xl mx-auto">
    <div class="flex justify-between items-center mb-8">
        <div>
            <h1 class="text-2xl font-bold text-gray-900">All Workflows</h1>
            <p class="mt-1 text-sm text-gray-500">
                Manage automation workflows across all your sites.
            </p>
        </div>
        <button
            onclick={() => {
                showCreateModal = true;
                initializeGraph();
            }}
            class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        >
            Create Workflow
        </button>
    </div>

    {#if loading}
        <div class="text-center py-12">Loading...</div>
    {:else}
        <!-- Workflow List (Unchanged) -->
        <div class="bg-white shadow overflow-hidden sm:rounded-lg">
            <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                    <tr>
                        <th
                            class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                            >Workflow Name</th
                        >
                        <th
                            class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                            >Site</th
                        >
                        <th
                            class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                            >Trigger</th
                        >
                        <th
                            class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                            >Status</th
                        >
                        <th class="relative px-6 py-3"
                            ><span class="sr-only">Edit</span></th
                        >
                    </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                    {#each workflows as wf}
                        <tr class="hover:bg-gray-50">
                            <td
                                class="px-6 py-4 whitespace-nowrap text-sm font-medium text-indigo-600"
                                >{wf.name}</td
                            >
                            <td
                                class="px-6 py-4 whitespace-nowrap text-sm text-gray-500"
                                >{wf.site_name || `Site #${wf.site_id}`}</td
                            >
                            <td
                                class="px-6 py-4 whitespace-nowrap text-sm text-gray-500"
                                >{wf.trigger_event}</td
                            >
                            <td class="px-6 py-4 whitespace-nowrap"
                                ><span
                                    class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800"
                                    >{wf.status}</span
                                ></td
                            >
                            <td
                                class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium"
                            >
                                <a
                                    href={`/sites/${wf.site_id}`}
                                    onclick={(e) => {
                                        e.preventDefault();
                                        // router.navigate(`/sites/${wf.site_id}`); // Removed router import
                                    }}
                                    class="text-indigo-600 hover:text-indigo-900"
                                    >View Site</a
                                >
                            </td>
                        </tr>
                    {:else}
                        <tr
                            ><td
                                colspan="5"
                                class="px-6 py-8 text-center text-gray-500 text-sm"
                                >No workflows found.</td
                            ></tr
                        >
                    {/each}
                </tbody>
            </table>
        </div>
    {/if}
</div>

{#if showCreateModal}
    <Modal size="2xl" onclose={() => (showCreateModal = false)}>
        <div
            class="bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:w-full h-[85vh] flex flex-col"
        >
            <!-- Header -->
            <div
                class="bg-gray-50 px-4 py-3 border-b flex justify-between items-center"
            >
                <h3 class="text-lg leading-6 font-medium text-gray-900">
                    Create Graph Workflow
                </h3>
                <input
                    type="text"
                    bind:value={newWorkflowName}
                    placeholder="Workflow Name"
                    class="ml-4 flex-1 border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                />
            </div>

            <!-- Main Content Area -->
            <div class="flex-1 flex overflow-hidden">
                <!-- Sidebar: Configuration -->
                <div
                    class="w-1/4 min-w-[300px] border-r border-gray-200 bg-white p-4 overflow-y-auto z-10"
                >
                    <h4 class="font-medium text-gray-900 mb-4">Settings</h4>

                    <!-- Trigger Config -->
                    <div class="space-y-4 mb-6">
                        <div>
                            <label
                                class="block text-xs font-medium text-gray-500 uppercase"
                                >Trigger Type</label
                            >
                            <div class="mt-1 flex gap-2">
                                <label class="inline-flex items-center text-sm">
                                    <input
                                        type="radio"
                                        bind:group={triggerType}
                                        value="EVENT"
                                        class="form-radio text-indigo-600"
                                    />
                                    <span class="ml-1">Event</span>
                                </label>
                                <label class="inline-flex items-center text-sm">
                                    <input
                                        type="radio"
                                        bind:group={triggerType}
                                        value="SCHEDULE"
                                        class="form-radio text-indigo-600"
                                    />
                                    <span class="ml-1">Schedule</span>
                                </label>
                            </div>
                        </div>

                        {#if triggerType === "EVENT"}
                            <div>
                                <label
                                    class="block text-sm font-medium text-gray-700"
                                    >Site</label
                                >
                                <select
                                    bind:value={selectedSiteId}
                                    onchange={handleSiteChange}
                                    class="mt-1 block w-full border-gray-300 rounded-md shadow-sm sm:text-sm"
                                >
                                    <option value="">Select...</option>
                                    {#each sites as site}
                                        <option value={site.id}
                                            >{site.name}</option
                                        >
                                    {/each}
                                </select>
                            </div>
                            <div>
                                <label
                                    class="block text-sm font-medium text-gray-700"
                                    >Event</label
                                >
                                <select
                                    bind:value={selectedTrigger}
                                    class="mt-1 block w-full border-gray-300 rounded-md shadow-sm sm:text-sm"
                                >
                                    <option value="">Select...</option>
                                    <option value="impression"
                                        >impression</option
                                    >
                                    <option value="click">click</option>
                                    <option value="conversion"
                                        >conversion</option
                                    >
                                    {#each eventDefinitions as event}
                                        <option value={event.name}
                                            >{event.name}</option
                                        >
                                    {/each}
                                </select>
                            </div>
                        {:else}
                            <div>
                                <label
                                    class="block text-sm font-medium text-gray-700"
                                    >Schedule (Cron)</label
                                >
                                <input
                                    type="text"
                                    bind:value={schedule}
                                    placeholder="0 12 * * *"
                                    class="mt-1 block w-full border-gray-300 rounded-md shadow-sm sm:text-sm"
                                />
                            </div>
                        {/if}
                    </div>

                    <hr class="my-4" />

                    <!-- Node Properties -->
                    {#if selectedNode}
                        <h4 class="font-medium text-gray-900 mb-2">
                            Node Properties
                        </h4>
                        <div class="text-xs text-gray-500 mb-2">
                            ID: {selectedNode.id}
                        </div>
                        <div class="space-y-3">
                            <div class="space-y-4">
                                <label
                                    class="block text-sm font-medium text-gray-700"
                                    >Action Type</label
                                >
                                {#if selectedNode.data}
                                    <select
                                        bind:value={selectedNode.data.action}
                                        class="mt-1 block w-full border-gray-300 rounded-md shadow-sm sm:text-sm"
                                    >
                                        <option value=""
                                            >Select Action...</option
                                        >
                                        <option value="Send Email"
                                            >Send Email</option
                                        >
                                        <option value="Log Event"
                                            >Log Event</option
                                        >
                                        <option value="FAIL"
                                            >Simulate Failure</option
                                        >
                                    </select>
                                {/if}
                            </div>

                            <div class="mt-4">
                                <label
                                    class="block text-sm font-medium text-gray-700"
                                    >Label</label
                                >
                                {#if selectedNode.data}
                                    <input
                                        type="text"
                                        bind:value={selectedNode.data.label}
                                        class="mt-1 block w-full border-gray-300 rounded-md shadow-sm sm:text-sm"
                                    />
                                {/if}
                            </div>

                            <div class="mt-4">
                                <label
                                    class="block text-sm font-medium text-gray-700"
                                    >Handle Position</label
                                >
                                {#if selectedNode.data}
                                    <select
                                        bind:value={
                                            selectedNode.data.handlePosition
                                        }
                                        class="mt-1 block w-full border-gray-300 rounded-md shadow-sm sm:text-sm"
                                    >
                                        <option value="right"
                                            >Right (Outputs)</option
                                        >
                                        <option value="left"
                                            >Left (Outputs)</option
                                        >
                                    </select>
                                {/if}
                            </div>
                        </div>
                    {:else}
                        <div class="text-sm text-gray-500 italic">
                            Select a node to edit properties.
                        </div>
                    {/if}

                    <div class="mt-8 space-y-2">
                        <label
                            class="block text-xs font-medium text-gray-500 uppercase"
                            >Add Node</label
                        >
                        <button
                            onclick={() => addNode("ACTION", "New Action")}
                            class="w-full text-left px-3 py-2 bg-blue-50 text-blue-700 rounded hover:bg-blue-100 text-sm font-medium"
                            >+ Action Node</button
                        >
                        <button
                            onclick={() =>
                                addNode("CONDITION", "New Condition")}
                            class="w-full text-left px-3 py-2 bg-yellow-50 text-yellow-700 rounded hover:bg-yellow-100 text-sm font-medium"
                            >+ Condition Node</button
                        >
                    </div>
                </div>

                <!-- Editor Canvas -->
                <div class="flex-1 bg-gray-100 relative">
                    <WorkflowEditor
                        bind:nodes={graphNodes}
                        bind:edges={graphEdges}
                        bind:selectedNodeId
                    />
                </div>
            </div>

            <!-- Footer -->
            <div
                class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse border-t"
            >
                <button
                    type="button"
                    onclick={createWorkflow}
                    class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:ml-3 sm:w-auto sm:text-sm"
                >
                    Create Workflow
                </button>
                <button
                    type="button"
                    onclick={() => (showCreateModal = false)}
                    class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
                >
                    Cancel
                </button>
            </div>
        </div>
    </Modal>
{/if}
