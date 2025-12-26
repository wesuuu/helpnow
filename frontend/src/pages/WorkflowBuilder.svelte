<script lang="ts">
    import { onMount } from "svelte";
    import { router } from "../lib/router.svelte.js";
    import WorkflowSidebar from "../lib/components/WorkflowSidebar.svelte";
    import WorkflowEditor from "../lib/components/WorkflowEditor.svelte";
    import WorkflowProperties from "../lib/components/WorkflowProperties.svelte";
    import Modal from "../lib/components/Modal.svelte";
    import { toast } from "../lib/stores/toast";
    import { useSvelteFlow } from "@xyflow/svelte";

    // Data State
    let workflowId = $state(0);
    let workflowName = $state("New Workflow");
    let status = $state("DRAFT");
    let loading = $state(true);

    // Graph Data
    let graphNodes = $state<any[]>([]);
    let graphEdges = $state<any[]>([]);
    let selectedNodeId = $state<string | null>(null);

    // Reference Data
    let sites = $state<any[]>([]);
    let audiences = $state<any[]>([]);
    let emailTemplates = $state<any[]>([]);
    let eventDefinitions = $state<any[]>([]);

    // Derived
    let selectedNode = $derived(
        selectedNodeId ? graphNodes.find((n) => n.id === selectedNodeId) : null,
    );

    let uniqueEventNames = $derived(
        Array.from(new Set(eventDefinitions.map((e) => e.name))).sort(),
    );

    // Modal State
    let showCreateTemplateModal = $state(false);
    let newTemplateName = $state("");
    let newTemplateSubject = $state("");
    let newTemplateBody = $state("");

    let showNameModal = $state(false);
    let tempWorkflowName = $state("");
    let isDuplicateError = $state(false);

    // --- Loading ---
    onMount(async () => {
        // Parse ID from URL
        // Simple regex to get ID from /workflows/:id
        const match = window.location.pathname.match(/\/workflows\/(\d+|new)/);
        if (match && match[1] !== "new") {
            workflowId = parseInt(match[1]);
            await loadWorkflow(workflowId);
        } else {
            // New Workflow
            await calculateDefaultName();
            initializeGraph();
            loading = false;
        }

        await Promise.all([
            loadSites(),
            loadAudiences(),
            loadEmailTemplates(),
            loadEventDefinitions(),
        ]);
    });

    let allWorkflows = $state<any[]>([]);

    function getUniqueName(baseName: string): string {
        let candidate = baseName;
        let counter = 1;

        // Helper to check if name exists (excluding current workflow)
        const exists = (name: string) => {
            return allWorkflows.some(
                (w) => w.name === name && w.id !== workflowId,
            );
        };

        if (!exists(candidate)) return candidate;

        // Try increments
        while (exists(`${baseName} ${counter}`)) {
            counter++;
        }
        return `${baseName} ${counter}`;
    }

    async function calculateDefaultName() {
        try {
            const res = await fetch("/api/workflows");
            if (res.ok) {
                allWorkflows = await res.json();
                workflowName = getUniqueName("New Workflow");
            }
        } catch (e) {
            console.error("Error calculating default name", e);
            workflowName = "New Workflow";
        }
    }

    async function loadWorkflow(id: number) {
        try {
            loading = true;
            // Fetch all workflows first for validation context
            fetch("/api/workflows")
                .then((r) => r.json())
                .then((data) => (allWorkflows = data));

            const res = await fetch(`/api/workflows/${id}`);
            if (res.ok) {
                const wf = await res.json();
                workflowName = wf.name;
                status = wf.status;
                if (wf.steps) {
                    const steps = JSON.parse(wf.steps);
                    graphNodes = steps.nodes || [];
                    graphEdges = steps.edges || [];
                } else {
                    initializeGraph();
                }
            } else {
                toast.error("Failed to load workflow");
                router.navigate("/workflows");
            }
        } catch (e) {
            console.error("Load flow error", e);
            toast.error("Error loading workflow");
        } finally {
            loading = false;
        }
    }

    async function loadSites() {
        const res = await fetch("/api/sites?organization_id=1");
        if (res.ok) sites = await res.json();
    }
    async function loadAudiences() {
        const res = await fetch("/api/audiences?organization_id=1");
        if (res.ok) audiences = await res.json();
    }
    async function loadEmailTemplates() {
        const res = await fetch("/api/email-templates?organization_id=1");
        if (res.ok) emailTemplates = await res.json();
    }
    async function loadEventDefinitions() {
        const res = await fetch("/api/events/definitions");
        if (res.ok) eventDefinitions = await res.json();
    }

    function initializeGraph() {
        graphNodes = [];
        graphEdges = [];
    }

    function addNode(
        type: string,
        label: string,
        position: { x: number; y: number },
        dataStr?: string,
    ) {
        const id = `n-${Date.now()}`;
        const defaultData = {
            label: label,
            handlePosition: "right",
            trigger_type: "EVENT",
            site_ids: [],
            audience_ids: [],
        };

        let extraData = {};
        if (dataStr) {
            try {
                extraData = JSON.parse(dataStr);
            } catch (e) {}
        }

        const newNode = {
            id,
            type,
            position,
            data: { ...defaultData, ...extraData },
        };

        graphNodes = [...graphNodes, newNode];
    }

    function onDrop(event: DragEvent) {
        event.preventDefault();
        if (!event.dataTransfer) return;

        const type = event.dataTransfer.getData("application/svelteflow");
        const dataStr = event.dataTransfer.getData(
            "application/svelteflow-data",
        );

        if (typeof type === "undefined" || !type) {
            return;
        }

        // Drop logic handled internally by WorkflowCanvas or via append here if canvas fails?
        // Actually, the refactoring in WorkflowCanvas handles drop.
        // We do NOT need duplicate logic here unless we want to control it from parent.
        // But for Click-To-Add, we definitely need logic here.
    }

    // --- Click to Add ---
    function handleAddNode(type: string, data?: any) {
        // Add to center of graph (approximation without getting current viewport from flow instance which is hard here)
        // Ideally we get center from flow instance. For MVP, we add at offset.
        // Better: We add it to 'graphNodes' array directly.
        const id = `n-${Date.now()}`;

        // Random slight offset to prevent perfect stacking if clicked multiple times fast
        const offset = Math.random() * 20;

        const defaultData = {
            label: type,
            handlePosition: "bottom",
            trigger_type: "EVENT",
            site_ids: [],
            audience_ids: [],
        };

        const newNode = {
            id,
            type,
            position: { x: 250 + offset, y: 250 + offset }, // Center-ish
            data: { ...defaultData, ...data },
        };

        graphNodes = [...graphNodes, newNode];
        toast.success(`Added ${data?.label || type} node`);
    }

    // --- Actions ---
    function performSaveValidation() {
        console.log("Checking validation...");
        console.log("Nodes:", graphNodes);
        console.log("Edges:", graphEdges);

        // Validation: At least one connected TRIGGER
        const triggers = graphNodes.filter((n) => n.type === "TRIGGER");
        if (triggers.length === 0) {
            toast.error("Workflow must have at least one trigger.");
            return false;
        }

        const connectedTriggers = triggers.filter((t) =>
            graphEdges.some((e) => e.source === t.id),
        );
        console.log("Connected Triggers:", connectedTriggers);

        if (connectedTriggers.length === 0) {
            toast.error(
                "Your trigger needs to be connected to an action to save.",
            );
            return false;
        }
        return true;
    }

    async function saveWorkflow() {
        console.log("Save clicked. Name:", workflowName);
        if (!performSaveValidation()) {
            console.log("Validation failed");
            return;
        }

        // Check generic name or duplicate
        const isGeneric = /^New Workflow(?: \d+)?$/.test(workflowName);
        const isDuplicate = allWorkflows.some(
            (w) => w.name === workflowName && w.id !== workflowId,
        );

        if (isGeneric || isDuplicate) {
            console.log(
                "Name issue detected. Generic:",
                isGeneric,
                "Duplicate:",
                isDuplicate,
            );
            isDuplicateError = isDuplicate;

            if (isDuplicate) {
                tempWorkflowName = getUniqueName(workflowName);
            } else {
                tempWorkflowName = workflowName;
            }
            showNameModal = true;
            return;
        }

        console.log("Executing save...");
        await executeSave();
    }

    async function confirmSaveName() {
        console.log("confirmSaveName called. tempName:", tempWorkflowName);
        if (!tempWorkflowName.trim()) {
            toast.error("Please enter a name.");
            return;
        }

        // Check for duplicate name
        const isDuplicate = allWorkflows.some(
            (w) => w.name === tempWorkflowName && w.id !== workflowId,
        );
        if (isDuplicate) {
            toast.error("Name already exists, cannot save");
            isDuplicateError = true;
            return;
        }

        workflowName = tempWorkflowName;
        showNameModal = false;
        console.log("Calling executeSave from modal...");
        await executeSave();
    }

    async function executeSave() {
        console.log("executeSave called. Name:", workflowName);
        const payload = {
            name: workflowName,
            status: status,
            organization_id: 1,
            trigger_type: "MULTIPLE",
            steps: JSON.stringify({ nodes: graphNodes, edges: graphEdges }),
        };

        try {
            let url = "/api/workflows";
            let method = "POST";
            if (workflowId) {
                url = `/api/workflows/${workflowId}`;
                method = "PUT";
            }

            console.log("Sending save request to:", url, method);

            const res = await fetch(url, {
                method,
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(payload),
            });

            console.log("Save response status:", res.status);

            if (res.ok) {
                const data = await res.json();
                if (!workflowId) {
                    workflowId = data.id;
                    // Update URL without reload
                    window.history.replaceState(
                        {},
                        "",
                        `/workflows/${workflowId}`,
                    );
                }

                // Refresh list on success to keep valid client-side cache
                // Don't await strictly necessary but good practice
                fetch("/api/workflows")
                    .then((r) => r.json())
                    .then((data) => (allWorkflows = data));

                toast.success("Workflow saved!");
            } else if (res.status === 409) {
                const text = await res.json(); // or text
                console.error("Conflict:", text);
                // toast.error("Name already exists. Suggesting unique name."); // Suppressed as per request

                // Suggest unique name and reopen modal
                isDuplicateError = true;
                tempWorkflowName = getUniqueName(workflowName);
                showNameModal = true;
            } else {
                const text = await res.text();
                console.error("Save failed response:", text);
                toast.error(`Failed to save (${res.status}): ${text}`);
            }
        } catch (e) {
            console.error("Save error:", e);
            toast.error("Error saving.");
        }
    }

    function backToList() {
        router.navigate("/workflows");
    }

    // --- Template Creation ---
    async function createTemplate() {
        // ...Logic from previous Workflows.svelte...
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
                if (selectedNode && selectedNode.data) {
                    selectedNode.data.template_id = newTemplate.id; // Auto select
                }
                newTemplateName = "";
                newTemplateSubject = "";
                newTemplateBody = "";
            }
        } catch (e) {
            console.error(e);
        }
    }
</script>

<div class="flex flex-col h-screen text-gray-900 bg-gray-50">
    <!-- HEADER -->
    <header
        class="h-14 bg-white border-b border-gray-200 flex items-center justify-between px-4 shrink-0"
    >
        <div class="flex items-center gap-4">
            <button
                onclick={backToList}
                class="text-gray-400 hover:text-gray-900"
                aria-label="Back to Workflows"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="20"
                    height="20"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    ><path d="M19 12H5" /><path d="M12 19l-7-7 7-7" /></svg
                >
            </button>
            <div class="flex flex-col">
                <input
                    type="text"
                    bind:value={workflowName}
                    class="bg-transparent border-none text-gray-900 font-semibold text-sm focus:ring-0 p-0 placeholder-gray-500 w-64 focus:outline-none"
                    placeholder="Workflow Name"
                />
                <span class="text-[10px] text-gray-500 uppercase tracking-wider"
                    >{status} • {workflowId ? "Saved" : "New"}</span
                >
            </div>
        </div>

        <div class="flex items-center gap-3">
            <button
                class="px-3 py-1.5 text-xs font-medium text-gray-700 hover:text-gray-900 border border-gray-300 rounded bg-white hover:bg-gray-50"
            >
                Test Run
            </button>
            <button
                onclick={saveWorkflow}
                class="px-3 py-1.5 text-xs font-medium text-white border border-transparent rounded bg-indigo-600 hover:bg-indigo-700 shadow-sm flex items-center gap-2"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="14"
                    height="14"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    ><path
                        d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"
                    ></path><polyline points="17 21 17 13 7 13 7 21"
                    ></polyline><polyline points="7 3 7 8 15 8"></polyline></svg
                >
                Save
            </button>
            <button
                class="px-3 py-1.5 text-xs font-medium text-white border border-transparent rounded bg-blue-600 hover:bg-blue-700 shadow-sm"
            >
                Publish
            </button>
        </div>
    </header>

    <!-- CONTENT -->
    <div class="flex-1 flex overflow-hidden">
        <!-- SIDEBAR -->
        <WorkflowSidebar onAddNode={handleAddNode} />

        <!-- CANVAS WRAPPER -->
        <div class="flex-1 bg-gray-50 relative" role="application">
            {#if loading}
                <div
                    class="absolute inset-0 flex items-center justify-center text-gray-500"
                >
                    Loading...
                </div>
            {:else}
                <WorkflowEditor
                    bind:nodes={graphNodes}
                    bind:edges={graphEdges}
                    bind:selectedNodeId
                />
            {/if}

            <!-- PROPERTIES OVERLAY -->
            {#if selectedNode}
                <div
                    class="absolute top-4 right-4 bottom-4 w-80 bg-white rounded-lg shadow-xl border border-gray-200 z-10 overflow-hidden"
                >
                    <WorkflowProperties
                        bind:selectedNode
                        {sites}
                        {audiences}
                        {emailTemplates}
                        {uniqueEventNames}
                        onCreateTemplate={() =>
                            (showCreateTemplateModal = true)}
                        onClose={() => (selectedNodeId = null)}
                    />
                </div>
            {/if}
        </div>
    </div>
</div>

<!-- Template Creation Modal -->
{#if showCreateTemplateModal}
    <Modal size="lg" onclose={() => (showCreateTemplateModal = false)}>
        <div class="bg-white p-6 rounded-lg text-slate-900">
            <h3 class="text-lg font-medium mb-4">Create Email Template</h3>
            <div class="space-y-4">
                <input
                    bind:value={newTemplateName}
                    placeholder="Template Name"
                    class="w-full border p-2 rounded"
                />
                <input
                    bind:value={newTemplateSubject}
                    placeholder="Subject"
                    class="w-full border p-2 rounded"
                />
                <textarea
                    bind:value={newTemplateBody}
                    placeholder="Body HTML"
                    class="w-full border p-2 rounded h-32"
                ></textarea>
                <div class="flex justify-end gap-2">
                    <button
                        onclick={() => (showCreateTemplateModal = false)}
                        class="px-4 py-2 text-sm border rounded">Cancel</button
                    >
                    <button
                        onclick={createTemplate}
                        class="px-4 py-2 text-sm bg-indigo-600 text-white rounded"
                        >create</button
                    >
                </div>
            </div>
        </div>
    </Modal>
{/if}

<!-- Naming Confirmation Modal -->
{#if showNameModal}
    <Modal size="sm" onclose={() => (showNameModal = false)}>
        <div class="bg-white p-6 rounded-lg text-gray-900">
            <h3 class="text-lg font-medium mb-4 text-gray-900">
                Name your Workflow
            </h3>
            <p class="text-sm text-gray-500 mb-4">
                {#if isDuplicateError}
                    The name is already taken. We've suggested a unique name
                    below.
                {:else}
                    This workflow still has a default name. Give it a unique
                    name to save.
                {/if}
            </p>
            <div class="space-y-4">
                <div>
                    <label
                        for="confirm-name"
                        class="block text-sm font-medium text-gray-700 mb-1"
                        >Workflow Name</label
                    >
                    <input
                        id="confirm-name"
                        type="text"
                        bind:value={tempWorkflowName}
                        placeholder="My Awesome Workflow"
                        class="block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm text-gray-900"
                    />
                </div>
                <div class="flex justify-end gap-3 pt-2">
                    <button
                        type="button"
                        onclick={() => (showNameModal = false)}
                        class="px-4 py-2 text-sm font-medium text-gray-700 hover:text-gray-900 bg-white border border-gray-300 rounded-md shadow-sm hover:bg-gray-50"
                    >
                        Cancel
                    </button>
                    <button
                        type="button"
                        onclick={() => {
                            console.log("Modal Save Clicked");
                            confirmSaveName();
                        }}
                        class="px-4 py-2 text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 border border-transparent rounded-md shadow-sm"
                    >
                        Save Workflow
                    </button>
                </div>
            </div>
        </div>
    </Modal>
{/if}
