<script lang="ts">
    import { onMount } from "svelte";
    import WorkflowEditor from "../lib/components/WorkflowEditor.svelte";
    import Modal from "../lib/components/Modal.svelte";
    import { toast } from "../lib/stores/toast";

    interface Workflow {
        id: number;
        site_id: number;
        site_name: string;
        name: string;
        trigger_event?: string;
        status: string;
        created_at: string;
        steps: string;
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

    interface EmailTemplate {
        id: number;
        name: string;
        subject: string;
        body: string;
    }

    let workflows = $state<Workflow[]>([]);
    let sites = $state<Site[]>([]);
    let eventDefinitions = $state<EventDefinition[]>([]);
    let audiences = $state<Audience[]>([]);
    let emailTemplates = $state<EmailTemplate[]>([]);
    let loading = $state(true);

    // Create Modal State
    let showCreateModal = $state(false);
    let newWorkflowName = $state("");
    let selectedSiteId = $state("");
    let selectedTrigger = $state("");
    let triggerType = $state("EVENT"); // 'EVENT' or 'SCHEDULE'
    let schedule = $state("");
    let selectedAudienceId = $state("");

    // Create Template State
    let showCreateTemplateModal = $state(false);
    let newTemplateName = $state("");
    let newTemplateSubject = $state("");
    let newTemplateBody = $state("");

    // Template Preview State
    let testFirstName = $state("John");
    let testLastName = $state("Doe");

    let previewSubject = $derived(
        newTemplateSubject
            .replace(/{{first_name}}/g, testFirstName)
            .replace(/{{last_name}}/g, testLastName),
    );
    let previewBody = $derived(
        newTemplateBody
            .replace(/{{first_name}}/g, testFirstName)
            .replace(/{{last_name}}/g, testLastName),
    );

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

    async function loadEmailTemplates() {
        try {
            const res = await fetch("/api/email-templates?organization_id=1");
            if (res.ok) {
                emailTemplates = await res.json();
            }
        } catch (e) {
            console.error("Failed to load email templates", e);
        }
    }

    async function loadAudiences() {
        try {
            const res = await fetch("/api/audiences?organization_id=1");
            if (res.ok) {
                audiences = await res.json();
            }
        } catch (e) {
            console.error("Failed to load audiences", e);
        }
    }

    async function createEmailTemplate() {
        if (!newTemplateName || !newTemplateSubject || !newTemplateBody) {
            toast.error("Please fill all template fields");
            return;
        }
        try {
            const res = await fetch("/api/email-templates", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    organization_id: 1,
                    name: newTemplateName,
                    subject: newTemplateSubject,
                    body: newTemplateBody,
                }),
            });
            if (res.ok) {
                const newTemplate = await res.json();
                emailTemplates = [...emailTemplates, newTemplate];
                showCreateTemplateModal = false;
                // Auto-select if a node is selected
                if (selectedNode && selectedNode.data) {
                    selectedNode.data.template_id = newTemplate.id;
                }
                newTemplateName = "";
                newTemplateSubject = "";
                newTemplateBody = "";
            } else {
                toast.error("Failed to create template");
            }
        } catch (e) {
            console.error("Error creating template", e);
            toast.error("Error creating template");
        }
    }

    // ... existing fetchEvents ...
    async function loadAllEvents() {
        try {
            const res = await fetch("/api/events/definitions");
            if (res.ok) {
                eventDefinitions = await res.json();
            }
        } catch (e) {
            console.error("Failed to load events", e);
        }
    }

    // Computed unique event names for dropdown
    let uniqueEventNames = $derived(
        Array.from(new Set(eventDefinitions.map((e) => e.name))).sort(),
    );

    function handleSiteChange() {
        selectedTrigger = ""; // Reset trigger
        // reset graph?
        initializeGraph();
    }

    // ... existing functions ...

    // ... addNode ...

    // ... createWorkflow ...

    function initializeGraph() {
        graphNodes = [];
        graphEdges = [];
        addNode("TRIGGER", "Start Trigger");
        if (graphNodes[0]) {
            graphNodes[0].position = { x: 50, y: 300 };
        }
    }

    function addNode(type: string, label: string = "") {
        const id = `n-${Date.now()}-${Math.floor(Math.random() * 1000)}`;
        let pos = { x: 100, y: 100 };
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
                trigger_type: "EVENT",
                site_ids: [],
                audience_ids: [],
                cron: "",
                trigger_event: "",
            },
        };

        graphNodes = [...graphNodes, newNode];

        if (graphNodes.length > 1) {
            const sourceNode = graphNodes[graphNodes.length - 2];
            graphEdges = [
                ...graphEdges,
                {
                    id: `e-${Date.now()}`,
                    source: sourceNode.id,
                    target: id,
                },
            ];
        }
    }

    let editingWorkflowId = $state(0);

    function editWorkflow(wf: Workflow) {
        editingWorkflowId = wf.id;
        newWorkflowName = wf.name;

        // Parse steps
        try {
            if (wf.steps) {
                const steps = JSON.parse(wf.steps);
                graphNodes = steps.nodes || [];
                graphEdges = steps.edges || [];
            } else {
                initializeGraph();
            }
        } catch (e) {
            console.error("Failed to parse workflow steps", e);
            initializeGraph();
        }

        showCreateModal = true;
    }

    function resetEditor() {
        newWorkflowName = "";
        selectedSiteId = "";
        selectedTrigger = "";
        schedule = "";
        selectedAudienceId = "";
        triggerType = "EVENT";
        editingWorkflowId = 0;
        initializeGraph();
    }

    async function createWorkflow() {
        if (!newWorkflowName) {
            toast.error("Please enter a workflow name");
            return;
        }

        const workflowData: any = {
            name: newWorkflowName,
            trigger_type: "MULTIPLE",
            steps: JSON.stringify({ nodes: graphNodes, edges: graphEdges }),
            status: "ACTIVE",
            organization_id: 1,
        };

        // Trigger validation is now implicit in the graph configuration
        // We could validate that at least one trigger exists in graphNodes
        const triggerNode = graphNodes.find((n) => n.type === "TRIGGER");
        if (!triggerNode) {
            toast.error("Workflow must have a Trigger node");
            return;
        }

        try {
            let url = "/api/workflows";
            let method = "POST";

            if (editingWorkflowId) {
                url = `/api/workflows/${editingWorkflowId}`;
                method = "PUT";
            }

            const res = await fetch(url, {
                method: method,
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(workflowData),
            });

            if (res.ok) {
                showCreateModal = false;
                resetEditor();
                loadWorkflows();
                toast.success("Workflow saved successfully");
            } else {
                toast.error("Failed to save workflow");
            }
        } catch (error) {
            console.error("Error saving workflow", error);
            toast.error("Error saving workflow");
        }
    }

    onMount(() => {
        loadWorkflows();
        loadSites();
        loadEmailTemplates();
        loadAllEvents();
        loadAudiences();
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
                resetEditor();
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
                            ><span class="sr-only">Actions</span></th
                        >
                    </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                    {#each workflows as wf}
                        <tr
                            class="hover:bg-gray-50 cursor-pointer"
                            onclick={() => editWorkflow(wf)}
                        >
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
                                        e.stopPropagation();
                                        // router.navigate(`/sites/${wf.site_id}`);
                                    }}
                                    class="text-gray-600 hover:text-gray-900"
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
                    {editingWorkflowId ? "Edit" : "Create"} Graph Workflow
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

                    <div class="space-y-4">
                        <div>
                            <label
                                class="block text-sm font-medium text-gray-700"
                                >Workflow Name</label
                            >
                            <input
                                type="text"
                                bind:value={newWorkflowName}
                                placeholder="My Workflow"
                                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm sm:text-sm"
                            />
                        </div>
                        <div
                            class="p-4 bg-blue-50 text-blue-800 text-xs rounded"
                        >
                            Configure triggers by clicking on the <strong
                                >Start Trigger</strong
                            > node in the graph.
                        </div>
                    </div>
                    <hr class="my-4" />

                    <!-- Node Properties -->
                    <div>
                        <h4 class="font-medium text-gray-900 mb-4">
                            Node Properties
                        </h4>
                        {#if selectedNode}
                            <div class="mb-4">
                                <label
                                    class="block text-xs font-medium text-gray-500 uppercase"
                                    >Label</label
                                >
                                <input
                                    type="text"
                                    bind:value={selectedNode.data.label}
                                    maxlength="50"
                                    class="mt-1 block w-full border-gray-300 rounded-md shadow-sm sm:text-sm"
                                />
                                {#if selectedNode.data.label && selectedNode.data.label.length >= 50}
                                    <p class="mt-1 text-xs text-red-600">
                                        Maximum length of 50 characters reached.
                                    </p>
                                {/if}
                            </div>

                            {#if selectedNode.type === "TRIGGER"}
                                <div class="space-y-4">
                                    <div>
                                        <label
                                            class="block text-sm font-medium text-gray-700"
                                            >Trigger Type</label
                                        >
                                        <select
                                            bind:value={
                                                selectedNode.data.trigger_type
                                            }
                                            class="block w-full border-gray-300 rounded-md shadow-sm sm:text-sm"
                                        >
                                            <option value="EVENT"
                                                >Event Based</option
                                            >
                                            <option value="SCHEDULE"
                                                >Scheduled</option
                                            >
                                        </select>
                                    </div>

                                    {#if selectedNode.data.trigger_type === "EVENT"}
                                        <div>
                                            <label
                                                class="block text-sm font-medium text-gray-700"
                                                >Event Name</label
                                            >
                                            <select
                                                bind:value={
                                                    selectedNode.data
                                                        .trigger_event
                                                }
                                                class="block w-full border-gray-300 rounded-md shadow-sm sm:text-sm"
                                            >
                                                <option value=""
                                                    >Select Event...</option
                                                >
                                                {#each uniqueEventNames as name}
                                                    <option value={name}
                                                        >{name}</option
                                                    >
                                                {/each}
                                            </select>
                                        </div>
                                        <div>
                                            <label
                                                class="block text-sm font-medium text-gray-700 mb-2"
                                                >Filter Sites</label
                                            >
                                            <div
                                                class="space-y-2 max-h-40 overflow-y-auto border p-2 rounded bg-gray-50"
                                            >
                                                {#each sites as site}
                                                    <label
                                                        class="flex items-center text-sm"
                                                    >
                                                        <input
                                                            type="checkbox"
                                                            bind:group={
                                                                selectedNode
                                                                    .data
                                                                    .site_ids
                                                            }
                                                            value={site.id}
                                                            class="rounded border-gray-300 text-indigo-600 focus:ring-indigo-500 mr-2"
                                                        />
                                                        {site.name}
                                                    </label>
                                                {/each}
                                            </div>
                                            <p
                                                class="text-xs text-gray-500 mt-1"
                                            >
                                                Select logic: OR (Run if event
                                                happens on any selected site).
                                            </p>
                                        </div>
                                    {:else}
                                        <div>
                                            <label
                                                class="block text-sm font-medium text-gray-700"
                                                >Cron Expression</label
                                            >
                                            <input
                                                type="text"
                                                bind:value={
                                                    selectedNode.data.cron
                                                }
                                                placeholder="* * * * *"
                                                class="block w-full border-gray-300 rounded-md shadow-sm sm:text-sm"
                                            />
                                            <p
                                                class="text-xs text-gray-500 mt-1"
                                            >
                                                Example: `0 9 * * *` (Every day
                                                at 9am)
                                            </p>
                                        </div>
                                    {/if}
                                </div>

                                <!-- Audience Filter (Common to Event & Schedule) -->
                                <div class="mt-4 border-t pt-4">
                                    <label
                                        class="block text-sm font-medium text-gray-700 mb-2"
                                        >Filter / Apply to Audiences</label
                                    >
                                    <div
                                        class="space-y-2 max-h-40 overflow-y-auto border p-2 rounded bg-gray-50"
                                    >
                                        {#each audiences as audience}
                                            <label
                                                class="flex items-center text-sm"
                                            >
                                                <input
                                                    type="checkbox"
                                                    bind:group={
                                                        selectedNode.data
                                                            .audience_ids
                                                    }
                                                    value={audience.id}
                                                    class="rounded border-gray-300 text-indigo-600 focus:ring-indigo-500 mr-2"
                                                />
                                                {audience.name}
                                            </label>
                                        {/each}
                                    </div>
                                    <p class="text-xs text-gray-500 mt-1">
                                        {#if selectedNode.data.trigger_type === "EVENT"}
                                            Trigger only if user is in selected
                                            audiences.
                                        {:else}
                                            Apply schedule context to selected
                                            audiences.
                                        {/if}
                                    </p>
                                </div>
                            {:else if selectedNode.type === "ACTION"}
                                <div>
                                    <label
                                        class="block text-xs font-medium text-gray-500 uppercase"
                                        >Action</label
                                    >
                                    {#if selectedNode.data}
                                        <select
                                            bind:value={
                                                selectedNode.data.action
                                            }
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

                                {#if selectedNode.data && selectedNode.data.action === "Send Email"}
                                    <div
                                        class="mt-4 p-3 bg-gray-50 rounded border border-gray-200"
                                    >
                                        <div
                                            class="flex justify-between items-center mb-2"
                                        >
                                            <label
                                                class="block text-xs font-medium text-gray-700 uppercase"
                                                >Email Template</label
                                            >
                                            <button
                                                onclick={() => {
                                                    newTemplateName = "";
                                                    newTemplateSubject = "";
                                                    newTemplateBody = "";
                                                    showCreateTemplateModal = true;
                                                }}
                                                class="text-indigo-600 hover:text-indigo-800 text-xs font-medium"
                                            >
                                                + New
                                            </button>
                                        </div>

                                        <select
                                            bind:value={
                                                selectedNode.data.template_id
                                            }
                                            class="block w-full border-gray-300 rounded-md shadow-sm sm:text-sm"
                                        >
                                            <option value=""
                                                >Select Template...</option
                                            >
                                            {#each emailTemplates as t}
                                                <option value={t.id}
                                                    >{t.name}</option
                                                >
                                            {/each}
                                        </select>
                                    </div>
                                {/if}

                                <div class="mt-4">
                                    <label
                                        class="block text-sm font-medium text-gray-700"
                                        >Delay (Days)</label
                                    >
                                    <input
                                        type="number"
                                        bind:value={
                                            selectedNode.data.delay_days
                                        }
                                        min="0"
                                        class="mt-1 block w-full border-gray-300 rounded-md shadow-sm sm:text-sm"
                                    />
                                </div>
                                <div class="mt-2">
                                    <label
                                        class="block text-sm font-medium text-gray-700"
                                        >Delay (Hours)</label
                                    >
                                    <input
                                        type="number"
                                        bind:value={
                                            selectedNode.data.delay_hours
                                        }
                                        min="0"
                                        class="mt-1 block w-full border-gray-300 rounded-md shadow-sm sm:text-sm"
                                    />
                                </div>
                                <div class="mt-2">
                                    <label
                                        class="block text-sm font-medium text-gray-700"
                                        >Delay (Minutes)</label
                                    >
                                    <input
                                        type="number"
                                        bind:value={
                                            selectedNode.data.delay_minutes
                                        }
                                        min="0"
                                        class="mt-1 block w-full border-gray-300 rounded-md shadow-sm sm:text-sm"
                                    />
                                </div>
                            {/if}
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
                                    addNode("TRIGGER", "New Trigger")}
                                class="w-full text-left px-3 py-2 bg-green-50 text-green-700 rounded hover:bg-green-100 text-sm font-medium"
                                >+ Trigger Node</button
                            >
                        </div>
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
                    {editingWorkflowId ? "Save Changes" : "Create Workflow"}
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

{#if showCreateTemplateModal}
    <Modal onclose={() => (showCreateTemplateModal = false)} size="xl">
        <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <h3
                class="text-lg leading-6 font-medium text-gray-900"
                id="modal-title"
            >
                Create Email Template
            </h3>

            <div class="mt-4 grid grid-cols-1 md:grid-cols-2 gap-6">
                <!-- Editor Column -->
                <div class="space-y-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700"
                            >Template Name</label
                        >
                        <input
                            type="text"
                            bind:value={newTemplateName}
                            class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                            placeholder="e.g., Welcome Email"
                        />
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700"
                            >Subject</label
                        >
                        <input
                            type="text"
                            bind:value={newTemplateSubject}
                            class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                            placeholder={"Subject line with {{first_name}}"}
                        />
                        <p class="mt-1 text-xs text-gray-500">
                            Supports variables: <code>{`{{first_name}}`}</code>,
                            <code>{`{{last_name}}`}</code>
                        </p>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700"
                            >Body</label
                        >
                        <textarea
                            bind:value={newTemplateBody}
                            rows="10"
                            class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm font-mono"
                        ></textarea>
                        <p class="mt-1 text-xs text-gray-500">
                            Supports variables: <code>{`{{first_name}}`}</code>,
                            <code>{`{{last_name}}`}</code>
                        </p>
                    </div>
                </div>

                <!-- Preview Column -->
                <div class="bg-gray-50 p-4 rounded-lg flex flex-col h-full">
                    <h4 class="text-sm font-medium text-gray-700 mb-2">
                        Preview
                    </h4>

                    <div class="grid grid-cols-2 gap-2 mb-4">
                        <div>
                            <label
                                class="block text-xs font-medium text-gray-500"
                                >Test First Name</label
                            >
                            <input
                                type="text"
                                bind:value={testFirstName}
                                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-xs"
                            />
                        </div>
                        <div>
                            <label
                                class="block text-xs font-medium text-gray-500"
                                >Test Last Name</label
                            >
                            <input
                                type="text"
                                bind:value={testLastName}
                                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-xs"
                            />
                        </div>
                    </div>

                    <div
                        class="bg-white border border-gray-200 rounded-md p-4 shadow-sm flex-1 overflow-y-auto hidden-scrollbar"
                    >
                        <div class="border-b border-gray-100 pb-2 mb-2">
                            <span class="text-xs text-gray-500 block"
                                >Subject:</span
                            >
                            <p
                                class="text-sm font-medium text-gray-900 break-words"
                            >
                                {previewSubject || "(No Subject)"}
                            </p>
                        </div>
                        <div>
                            <span class="text-xs text-gray-500 block mb-1"
                                >Body:</span
                            >
                            <div
                                class="text-sm text-gray-800 whitespace-pre-wrap font-sans"
                            >
                                {@html previewBody ||
                                    "<span class='text-gray-400 italic'>(Empty body)</span>"}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
            <button
                onclick={createEmailTemplate}
                type="button"
                class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 sm:ml-3 sm:w-auto sm:text-sm"
            >
                Save
            </button>
            <button
                onclick={() => (showCreateTemplateModal = false)}
                type="button"
                class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
            >
                Cancel
            </button>
        </div>
    </Modal>
{/if}
