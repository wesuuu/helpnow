<script lang="ts">
    import type { Node } from "@xyflow/svelte";
    import { onMount } from "svelte";

    // --- Types ---
    interface FieldSchema {
        name: string;
        type: string;
        required: boolean;
        description?: string;
        validations?: string[];
    }

    interface ComponentSchema {
        name: string;
        type: string;
        description?: string;
        fields: FieldSchema[];
    }

    // --- Props ---
    let {
        selectedNode = $bindable(null),
        sites = [],
        audiences = [],
        emailTemplates = [],
        uniqueEventNames = [],
        onCreateTemplate,
        onClose, // New prop
    } = $props<{
        selectedNode: Node | null;
        sites: any[];
        audiences: any[];
        emailTemplates: any[];
        uniqueEventNames: string[];
        onCreateTemplate: () => void;
        onClose?: () => void;
    }>();

    // --- State ---
    let actionSchemas = $state<ComponentSchema[]>([]);
    let logicSchemas = $state<ComponentSchema[]>([]);
    let triggerSchemas = $state<ComponentSchema[]>([]);
    let loadingSchemas = $state(false);

    // --- Derived ---
    let nodeTypeLabel = $derived(
        selectedNode?.data?.label || selectedNode?.type || "Settings",
    );

    // Determine the current schema based on node selection
    let currentSchema = $derived.by(() => {
        if (!selectedNode) return null;

        if (selectedNode.type === "ACTION") {
            const actionName = selectedNode.data.action;
            return actionSchemas.find((s) => s.name === actionName);
        }
        if (selectedNode.type === "TRIGGER") {
            const triggerType = selectedNode.data.trigger_type;
            return triggerSchemas.find((s) => s.name === triggerType); // e.g. "EVENT"
        }
        if (selectedNode.type === "CONDITION") {
            // Default or selected logic type
            const logicName = selectedNode.data.logic_type || "Condition";
            return logicSchemas.find((s) => s.name === logicName);
        }
        return null;
    });

    // --- Fetch Schemas ---
    onMount(async () => {
        loadingSchemas = true;
        try {
            const [actionsRes, logicRes, triggersRes] = await Promise.all([
                fetch("/api/workflow-components/actions"),
                fetch("/api/workflow-components/logic"),
                fetch("/api/workflow-components/triggers"),
            ]);

            if (actionsRes.ok) actionSchemas = await actionsRes.json();
            if (logicRes.ok) logicSchemas = await logicRes.json();
            if (triggersRes.ok) triggerSchemas = await triggersRes.json();
        } catch (e) {
            console.error("Failed to load component schemas", e);
        } finally {
            loadingSchemas = false;
        }
    });

    // --- Helpers ---
    function isOneOf(field: FieldSchema): string[] | null {
        if (!field.validations) return null;
        for (const v of field.validations) {
            if (v.startsWith("oneof=")) {
                return v.replace("oneof=", "").split(" ");
            }
        }
        return null;
    }
</script>

<div class="h-full bg-white text-gray-900 flex flex-col">
    {#if !selectedNode}
        <div
            class="flex-1 flex items-center justify-center p-8 text-center text-gray-500 text-sm"
        >
            Select a node to configure settings.
        </div>
    {:else}
        <!-- Header -->
        <div
            class="p-4 border-b border-gray-200 flex justify-between items-center"
        >
            <div>
                <div
                    class="text-xs font-bold text-indigo-600 uppercase tracking-wider mb-0.5"
                >
                    {selectedNode.type}
                </div>
                <h2 class="font-semibold text-lg text-gray-900">
                    {nodeTypeLabel}
                </h2>
            </div>
            {#if onClose}
                <button
                    onclick={onClose}
                    class="text-gray-400 hover:text-gray-600 p-1"
                    aria-label="Close Properties"
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
                        ><line x1="18" y1="6" x2="6" y2="18"></line><line
                            x1="6"
                            y1="6"
                            x2="18"
                            y2="18"
                        ></line></svg
                    >
                </button>
            {/if}
        </div>

        <!-- content -->
        <div class="flex-1 overflow-y-auto p-4 space-y-6">
            {#if loadingSchemas}
                <div class="text-sm text-gray-500 text-center py-4">
                    Loading configuration options...
                </div>
            {:else}
                <!-- Global Node Settings -->
                <div class="space-y-3">
                    <label
                        for="node-name"
                        class="block text-xs font-medium text-gray-500 uppercase"
                        >Node Name</label
                    >
                    <input
                        id="node-name"
                        type="text"
                        bind:value={selectedNode.data.label}
                        class="block w-full bg-white border-gray-300 rounded-md text-sm text-gray-900 focus:ring-indigo-500 focus:border-indigo-500 placeholder-gray-400"
                        placeholder="Enter node name"
                    />
                </div>

                <hr class="border-gray-200" />

                <!-- Component Type Selection -->
                <div class="space-y-4">
                    {#if selectedNode.type === "ACTION"}
                        <div>
                            <label
                                for="action-type"
                                class="block text-xs font-medium text-gray-500 uppercase mb-2"
                                >Action Type</label
                            >
                            <select
                                id="action-type"
                                bind:value={selectedNode.data.action}
                                class="block w-full bg-white border-gray-300 rounded-md text-sm text-gray-900 focus:ring-indigo-500 focus:border-indigo-500"
                            >
                                <option value="">Select Action...</option>
                                {#each actionSchemas as schema}
                                    <option value={schema.name}
                                        >{schema.name}</option
                                    >
                                {/each}
                            </select>
                        </div>
                    {:else if selectedNode.type === "TRIGGER"}
                        <div>
                            <label
                                for="trigger-type"
                                class="block text-xs font-medium text-gray-500 uppercase mb-2"
                                >Trigger Type</label
                            >
                            <select
                                id="trigger-type"
                                bind:value={selectedNode.data.trigger_type}
                                class="block w-full bg-white border-gray-300 rounded-md text-sm text-gray-900 focus:ring-indigo-500 focus:border-indigo-500"
                            >
                                <option value="">Select Trigger...</option>
                                {#each triggerSchemas as schema}
                                    <option value={schema.name}
                                        >{schema.name}</option
                                    >
                                {/each}
                            </select>
                        </div>
                    {:else if selectedNode.type === "CONDITION"}
                        <!-- Logic types are less varied for now, but keeping dynamic -->
                    {/if}

                    <!-- Dynamic Fields -->
                    {#if currentSchema}
                        <div class="space-y-5 pt-2">
                            <!-- Helper: Schema Description -->
                            {#if currentSchema.description}
                                <p class="text-xs text-gray-500 italic">
                                    {currentSchema.description}
                                </p>
                            {/if}

                            {#each currentSchema.fields as field}
                                <div class="space-y-2">
                                    <label
                                        for={field.name}
                                        class="block text-xs font-medium text-gray-500 uppercase"
                                    >
                                        {field.name.replace(/_/g, " ")}
                                        {#if field.required}
                                            <span class="text-red-500">*</span>
                                        {/if}
                                    </label>

                                    <!-- Render Input Based on Field Name (Rich Widgets) or Type -->
                                    {#if field.name === "site_ids"}
                                        <!-- SITE SELECTOR -->
                                        <div
                                            class="bg-white border border-gray-300 rounded p-2 max-h-32 overflow-y-auto"
                                        >
                                            {#each sites as site}
                                                <label
                                                    class="flex items-center text-sm py-1"
                                                >
                                                    <input
                                                        type="checkbox"
                                                        bind:group={
                                                            selectedNode.data[
                                                                field.name
                                                            ]
                                                        }
                                                        value={site.id}
                                                        class="rounded border-gray-300 text-indigo-600 focus:ring-indigo-500 mr-2"
                                                    />
                                                    <span class="text-gray-700"
                                                        >{site.name}</span
                                                    >
                                                </label>
                                            {/each}
                                        </div>
                                    {:else if field.name === "audience_ids"}
                                        <!-- AUDIENCE SELECTOR -->
                                        <div
                                            class="bg-white border border-gray-300 rounded p-2 max-h-32 overflow-y-auto"
                                        >
                                            {#each audiences as audience}
                                                <label
                                                    class="flex items-center text-sm py-1"
                                                >
                                                    <input
                                                        type="checkbox"
                                                        bind:group={
                                                            selectedNode.data[
                                                                field.name
                                                            ]
                                                        }
                                                        value={audience.id}
                                                        class="rounded border-gray-300 text-indigo-600 focus:ring-indigo-500 mr-2"
                                                    />
                                                    <span class="text-gray-700"
                                                        >{audience.name}</span
                                                    >
                                                </label>
                                            {/each}
                                        </div>
                                    {:else if field.name === "template_id"}
                                        <!-- TEMPLATE SELECTOR -->
                                        <div class="flex flex-col gap-2">
                                            <div
                                                class="flex justify-between items-center"
                                            >
                                                <span
                                                    class="text-[10px] text-gray-400"
                                                    >Select ID</span
                                                >
                                                {#if onCreateTemplate}
                                                    <button
                                                        type="button"
                                                        onclick={onCreateTemplate}
                                                        class="text-xs text-indigo-600 hover:text-indigo-500 font-medium"
                                                        >+ New Template</button
                                                    >
                                                {/if}
                                            </div>
                                            <select
                                                id={field.name}
                                                bind:value={
                                                    selectedNode.data[
                                                        field.name
                                                    ]
                                                }
                                                class="block w-full bg-white border-gray-300 rounded-md text-sm text-gray-900 focus:ring-indigo-500 focus:border-indigo-500"
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
                                    {:else if field.name === "trigger_event"}
                                        <!-- EVENT NAME SELECTOR -->
                                        <select
                                            id={field.name}
                                            bind:value={
                                                selectedNode.data[field.name]
                                            }
                                            class="block w-full bg-white border-gray-300 rounded-md text-sm text-gray-900 focus:ring-indigo-500 focus:border-indigo-500"
                                        >
                                            <option value=""
                                                >Select Event...</option
                                            >
                                            {#each uniqueEventNames as name}
                                                <option value={name}
                                                    >{name}</option
                                                >
                                            {/each}
                                            <!-- Fallbacks if not in definitions yet -->
                                            {#if !uniqueEventNames.includes("email_opened")}
                                                <option value="email_opened"
                                                    >email_opened</option
                                                >
                                            {/if}
                                            {#if !uniqueEventNames.includes("webhook_received")}
                                                <option value="webhook_received"
                                                    >webhook_received</option
                                                >
                                            {/if}
                                        </select>
                                    {:else}
                                        <!-- GENERIC FIELDS -->
                                        {#if isOneOf(field)}
                                            <!-- ONEOF DROPDOWN -->
                                            <select
                                                id={field.name}
                                                bind:value={
                                                    selectedNode.data[
                                                        field.name
                                                    ]
                                                }
                                                class="block w-full bg-white border-gray-300 rounded-md text-sm text-gray-900 focus:ring-indigo-500 focus:border-indigo-500"
                                            >
                                                {#each isOneOf(field) || [] as opt}
                                                    <option value={opt}
                                                        >{opt}</option
                                                    >
                                                {/each}
                                            </select>
                                        {:else if field.type.includes("int") || field.type.includes("float")}
                                            <!-- NUMBER INPUT -->
                                            <input
                                                id={field.name}
                                                type="number"
                                                bind:value={
                                                    selectedNode.data[
                                                        field.name
                                                    ]
                                                }
                                                class="block w-full bg-white border-gray-300 rounded-md text-sm text-gray-900 focus:ring-indigo-500 focus:border-indigo-500"
                                            />
                                        {:else}
                                            <!-- DEFAULT STRING INPUT -->
                                            <input
                                                id={field.name}
                                                type="text"
                                                bind:value={
                                                    selectedNode.data[
                                                        field.name
                                                    ]
                                                }
                                                class="block w-full bg-white border-gray-300 rounded-md text-sm text-gray-900 focus:ring-indigo-500 focus:border-indigo-500"
                                            />
                                        {/if}
                                    {/if}

                                    <!-- Helper: Field Description -->
                                    {#if field.description}
                                        <p class="text-xs text-gray-500">
                                            {field.description}
                                        </p>
                                    {/if}
                                </div>
                            {/each}
                        </div>
                    {/if}
                </div>
            {/if}
        </div>
    {/if}
</div>
