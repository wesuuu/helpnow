<script lang="ts">
    import { onMount } from "svelte";
    import Modal from "../lib/components/Modal.svelte";
    import Table from "../lib/components/Table.svelte";
    import Tooltip from "../lib/components/Tooltip.svelte";
    import { toast } from "../lib/stores/toast";

    type Agent = {
        id: number;
        organization_id: number;
        name: string;
        description: string;
        model_config: string;
        created_at: string;
    };

    let agents: Agent[] = [];
    let loading = true;
    let editingAgent: Agent | null = null;
    let isCreating = false;

    // Form inputs
    let formName = "";
    let formDescription = "";
    let formSystemPrompt = "";
    let formTools: string[] = [];

    // Constants
    const MAX_NAME_LEN = 20;
    const MAX_DESC_LEN = 100;
    const AVAILABLE_TOOLS = [
        { id: "http_request", label: "HTTP Request" },
        { id: "send_email", label: "Send Email" }, // Added hypothetical tools for multi-select demo
        { id: "db_query", label: "DB Query" },
    ];

    async function fetchAgents() {
        loading = true;
        try {
            const res = await fetch("/api/agents?org_id=1");
            if (res.ok) {
                agents = await res.json();
            } else {
                toast.error("Failed to load agents");
            }
        } catch (e) {
            toast.error("Error loading agents");
        } finally {
            loading = false;
        }
    }

    function openModal(agent?: Agent) {
        if (agent) {
            isCreating = false;
            editingAgent = agent;
            formName = agent.name;
            formDescription = agent.description;

            // Parse config
            try {
                const config = JSON.parse(agent.model_config || "{}");
                formSystemPrompt = config.system_prompt || "";
                formTools = Array.isArray(config.tools)
                    ? config.tools
                    : ["http_request"];
            } catch (e) {
                formSystemPrompt = "";
                formTools = ["http_request"];
            }
        } else {
            isCreating = true;
            editingAgent = { id: 0 } as Agent;
            formName = "";
            formDescription = "";
            formSystemPrompt = "";
            formTools = ["http_request"];
        }
    }

    function closeModal() {
        editingAgent = null;
        isCreating = false;
    }

    async function saveAgent() {
        if (!formName) return toast.error("Name is required");
        if (formName.length > MAX_NAME_LEN) return toast.error("Name too long");
        if (formDescription.length > MAX_DESC_LEN)
            return toast.error("Description too long");

        // Build config
        const config = {
            system_prompt: formSystemPrompt,
            tools: formTools,
        };

        const payload = {
            organization_id: 1,
            name: formName,
            description: formDescription,
            model_config: JSON.stringify(config),
        };

        try {
            const url = isCreating
                ? "/api/agents"
                : `/api/agents/${editingAgent!.id}`;
            const method = isCreating ? "POST" : "PUT";

            const res = await fetch(url, {
                method,
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(payload),
            });

            if (res.ok) {
                toast.success(isCreating ? "Agent created" : "Agent updated");
                closeModal();
                fetchAgents();
            } else {
                toast.error("Failed to save agent");
            }
        } catch (e) {
            toast.error("Error saving agent");
        }
    }

    async function deleteAgent(id: number) {
        if (!confirm("Are you sure you want to delete this agent?")) return;

        try {
            const res = await fetch(`/api/agents/${id}`, { method: "DELETE" });
            if (res.ok) {
                toast.success("Agent deleted");
                // Correct logic: if we are in modal for this agent, close it.
                if (editingAgent && editingAgent.id === id) {
                    closeModal();
                }
                fetchAgents();
            } else {
                toast.error("Failed to delete agent");
            }
        } catch (e) {
            toast.error("Error deleting agent");
        }
    }

    function toggleTool(toolId: string) {
        if (formTools.includes(toolId)) {
            formTools = formTools.filter((t) => t !== toolId);
        } else {
            formTools = [...formTools, toolId];
        }
    }

    onMount(() => {
        fetchAgents();
    });

    const columns = [
        {
            key: "name",
            label: "Name",
            class: "font-medium text-gray-900 dark:text-white",
        },
        {
            key: "description",
            label: "Description",
            class: "text-gray-500 dark:text-gray-400",
        },
        {
            key: "created_at",
            label: "Created",
            format: (val: string) => new Date(val).toLocaleDateString(),
            class: "text-gray-500 dark:text-gray-400",
        },
    ];
</script>

<div class="space-y-6">
    <div class="flex justify-between items-center">
        <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">
            Agents
        </h1>
        <button
            class="px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700"
            onclick={() => openModal()}
        >
            Create Agent
        </button>
    </div>

    {#if loading}
        <div class="text-center py-4">Loading...</div>
    {:else if agents.length === 0}
        <div class="text-center py-10 bg-gray-50 dark:bg-gray-800 rounded-lg">
            <p class="text-gray-500">
                No agents found. Create one to get started.
            </p>
        </div>
    {:else}
        <Table {columns} data={agents} onRowClick={openModal} />
    {/if}

    {#if editingAgent}
        <Modal size="md" onclose={closeModal}>
            <div class="-m-6">
                <!-- Header -->
                <div
                    class="flex justify-between items-center p-4 border-b dark:border-gray-700"
                >
                    <h3
                        class="text-lg font-medium text-gray-900 dark:text-white"
                    >
                        {isCreating ? "Create Agent" : "Edit Agent"}
                    </h3>
                    <button
                        class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
                        onclick={closeModal}
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="h-6 w-6"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M6 18L18 6M6 6l12 12"
                            />
                        </svg>
                    </button>
                </div>

                <!-- Body -->
                <div class="p-6 space-y-4">
                    <!-- Name -->
                    <div>
                        <label
                            class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                        >
                            Name
                        </label>
                        <input
                            type="text"
                            bind:value={formName}
                            maxlength={MAX_NAME_LEN}
                            class="w-full border rounded px-3 py-2 dark:bg-gray-700 dark:border-gray-600 dark:text-white {formName.length ===
                            MAX_NAME_LEN
                                ? 'border-red-500 focus:ring-red-500'
                                : ''}"
                            placeholder="e.g. Sales Assistant"
                        />
                        <div class="flex justify-between text-xs mt-1">
                            <span
                                class="text-red-500 {formName.length <
                                MAX_NAME_LEN
                                    ? 'invisible'
                                    : ''}">Max length check</span
                            >
                            <span
                                class={formName.length === MAX_NAME_LEN
                                    ? "text-red-500"
                                    : "text-gray-500"}
                                >{formName.length}/{MAX_NAME_LEN}</span
                            >
                        </div>
                    </div>

                    <!-- Description -->
                    <div>
                        <label
                            class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                        >
                            Description
                        </label>
                        <textarea
                            bind:value={formDescription}
                            maxlength={MAX_DESC_LEN}
                            rows="2"
                            class="w-full border rounded px-3 py-2 dark:bg-gray-700 dark:border-gray-600 dark:text-white {formDescription.length ===
                            MAX_DESC_LEN
                                ? 'border-red-500 focus:ring-red-500'
                                : ''}"
                            placeholder="Brief description of the agent"
                        ></textarea>
                        <div class="flex justify-between text-xs mt-1">
                            <span
                                class="text-red-500 {formDescription.length <
                                MAX_DESC_LEN
                                    ? 'invisible'
                                    : ''}">Max length check</span
                            >
                            <span
                                class={formDescription.length === MAX_DESC_LEN
                                    ? "text-red-500"
                                    : "text-gray-500"}
                                >{formDescription.length}/{MAX_DESC_LEN}</span
                            >
                        </div>
                    </div>

                    <!-- System Prompt -->
                    <div>
                        <label
                            class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1 flex items-center"
                        >
                            System Prompt
                            <Tooltip
                                text="The core instructions that define the agent's behavior and personality."
                            >
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    class="h-4 w-4 ml-1 text-gray-400 cursor-help"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                    stroke="currentColor"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                                    />
                                </svg>
                            </Tooltip>
                        </label>
                        <textarea
                            bind:value={formSystemPrompt}
                            rows="4"
                            class="w-full border rounded px-3 py-2 text-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                            placeholder="You are a helpful assistant..."
                        ></textarea>
                    </div>

                    <!-- Tools -->
                    <div>
                        <label
                            class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1 flex items-center"
                        >
                            Tools
                            <Tooltip
                                text="Select the capabilities available to this agent."
                            >
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    class="h-4 w-4 ml-1 text-gray-400 cursor-help"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                    stroke="currentColor"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                                    />
                                </svg>
                            </Tooltip>
                        </label>
                        <div
                            class="border rounded p-3 dark:border-gray-600 dark:bg-gray-700 max-h-40 overflow-y-auto"
                        >
                            {#each AVAILABLE_TOOLS as tool}
                                <label
                                    class="flex items-center space-x-2 py-1 cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-600 rounded px-1"
                                >
                                    <input
                                        type="checkbox"
                                        checked={formTools.includes(tool.id)}
                                        onchange={() => toggleTool(tool.id)}
                                        class="rounded text-indigo-600 focus:ring-indigo-500"
                                    />
                                    <span
                                        class="text-sm text-gray-700 dark:text-gray-300"
                                        >{tool.label}</span
                                    >
                                </label>
                            {/each}
                        </div>
                    </div>
                </div>

                <!-- Footer -->
                <div
                    class="p-4 border-t bg-gray-50 dark:bg-gray-700 dark:border-gray-600 flex justify-between rounded-b-lg"
                >
                    <div>
                        {#if !isCreating}
                            <button
                                class="px-4 py-2 text-red-600 hover:text-red-900 hover:bg-red-50 rounded"
                                onclick={() => deleteAgent(editingAgent!.id)}
                            >
                                Delete
                            </button>
                        {/if}
                    </div>
                    <div class="flex space-x-3">
                        <button
                            class="px-4 py-2 bg-white border border-gray-300 rounded text-gray-700 hover:bg-gray-50"
                            onclick={closeModal}
                        >
                            Cancel
                        </button>
                        <button
                            class="px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700"
                            onclick={saveAgent}
                        >
                            Save
                        </button>
                    </div>
                </div>
            </div>
        </Modal>
    {/if}
</div>
