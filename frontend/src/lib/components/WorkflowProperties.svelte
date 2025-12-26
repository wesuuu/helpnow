<script lang="ts">
    import type { Node } from "@xyflow/svelte";
    import { onMount } from "svelte";

    // Props
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

    // Derived Label
    let nodeTypeLabel = $derived(
        selectedNode?.data?.label || selectedNode?.type || "Settings",
    );
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

            <!-- Trigger Specific -->
            {#if selectedNode.type === "TRIGGER"}
                <div class="space-y-4">
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
                            <option value="EVENT">Event Based</option>
                            <option value="SCHEDULE">Scheduled</option>
                        </select>
                    </div>

                    {#if selectedNode.data.trigger_type === "EVENT"}
                        <div>
                            <label
                                for="event-name"
                                class="block text-xs font-medium text-gray-500 uppercase mb-2"
                                >Event Name</label
                            >
                            <select
                                id="event-name"
                                bind:value={selectedNode.data.trigger_event}
                                class="block w-full bg-white border-gray-300 rounded-md text-sm text-gray-900 focus:ring-indigo-500 focus:border-indigo-500"
                            >
                                <option value="">Select Event...</option>
                                {#each uniqueEventNames as name}
                                    <option value={name}>{name}</option>
                                {/each}
                                <!-- Manually add Email Opened if not in list yet? -->
                                <option value="email_opened"
                                    >email_opened</option
                                >
                                <option value="webhook_received"
                                    >webhook_received</option
                                >
                            </select>
                        </div>

                        <div>
                            <span
                                class="block text-xs font-medium text-gray-500 uppercase mb-2"
                                >Filter Sites</span
                            >
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
                                                selectedNode.data.site_ids
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
                        </div>
                    {:else}
                        <div>
                            <label
                                for="cron-expression"
                                class="block text-xs font-medium text-gray-500 uppercase mb-2"
                                >Cron Expression</label
                            >
                            <input
                                id="cron-expression"
                                type="text"
                                bind:value={selectedNode.data.cron}
                                placeholder="0 9 * * *"
                                class="block w-full bg-white border-gray-300 rounded-md text-sm text-gray-900 focus:ring-indigo-500 focus:border-indigo-500 font-mono"
                            />
                            <p class="text-xs text-gray-500 mt-1">
                                e.g. 0 9 * * * (Daily at 9am)
                            </p>
                        </div>
                    {/if}

                    <div class="pt-4 border-t border-gray-200">
                        <span
                            class="block text-xs font-medium text-gray-500 uppercase mb-2"
                            >Audiences</span
                        >
                        <div
                            class="bg-white border border-gray-300 rounded p-2 max-h-32 overflow-y-auto"
                        >
                            {#each audiences as audience}
                                <label class="flex items-center text-sm py-1">
                                    <input
                                        type="checkbox"
                                        bind:group={
                                            selectedNode.data.audience_ids
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
                    </div>
                </div>

                <!-- Action Specific -->
            {:else if selectedNode.type === "ACTION"}
                <div class="space-y-4">
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
                            <option value="Send Email">Send Email</option>
                            <option value="Log Event">Log Event</option>
                            <option value="Update DB">Update Database</option>
                            <option value="HTTP Request">HTTP Request</option>
                            <option value="Delay">Delay</option>
                            <option value="FAIL">Simulate Failure</option>
                        </select>
                    </div>

                    {#if selectedNode.data.action === "Send Email"}
                        <div
                            class="p-3 bg-gray-50 rounded border border-gray-200 space-y-3"
                        >
                            <div class="flex justify-between items-center">
                                <label
                                    for="email-template"
                                    class="text-xs font-medium text-gray-500 uppercase"
                                    >Email Template</label
                                >
                                <button
                                    type="button"
                                    onclick={onCreateTemplate}
                                    class="text-xs text-indigo-600 hover:text-indigo-500 font-medium"
                                    >+ New</button
                                >
                            </div>
                            <select
                                id="email-template"
                                bind:value={selectedNode.data.template_id}
                                class="block w-full bg-white border-gray-300 rounded-md text-sm text-gray-900 focus:ring-indigo-500 focus:border-indigo-500"
                            >
                                <option value="">Select Template...</option>
                                {#each emailTemplates as t}
                                    <option value={t.id}>{t.name}</option>
                                {/each}
                            </select>
                        </div>
                    {/if}

                    {#if selectedNode.data.action === "Delay"}
                        <div
                            class="p-3 bg-gray-50 rounded border border-gray-200 space-y-3"
                        >
                            <span
                                class="text-xs font-medium text-gray-500 uppercase"
                                >Delay Duration</span
                            >
                            <div class="grid grid-cols-2 gap-2">
                                <div>
                                    <span
                                        class="text-xs text-gray-500 block mb-1"
                                        >Days</span
                                    >
                                    <input
                                        type="number"
                                        min="0"
                                        bind:value={
                                            selectedNode.data.delay_days
                                        }
                                        class="block w-full bg-white border-gray-300 rounded-md text-sm text-gray-900"
                                    />
                                </div>
                                <div>
                                    <span
                                        class="text-xs text-gray-500 block mb-1"
                                        >Hours</span
                                    >
                                    <input
                                        type="number"
                                        min="0"
                                        bind:value={
                                            selectedNode.data.delay_hours
                                        }
                                        class="block w-full bg-white border-gray-300 rounded-md text-sm text-gray-900"
                                    />
                                </div>
                                <div class="col-span-2">
                                    <span
                                        class="text-xs text-gray-500 block mb-1"
                                        >Minutes</span
                                    >
                                    <input
                                        type="number"
                                        min="0"
                                        bind:value={
                                            selectedNode.data.delay_minutes
                                        }
                                        class="block w-full bg-white border-gray-300 rounded-md text-sm text-gray-900"
                                    />
                                </div>
                            </div>
                        </div>
                    {/if}
                </div>

                <!-- Condition Specific -->
            {:else if selectedNode.type === "CONDITION"}
                <div class="space-y-4">
                    <div>
                        <label
                            for="debug-force"
                            class="block text-xs font-medium text-gray-500 uppercase mb-2"
                            >Debug Force Outcome</label
                        >
                        <select
                            id="debug-force"
                            bind:value={selectedNode.data.force}
                            class="block w-full bg-white border-gray-300 rounded-md text-sm text-gray-900 focus:ring-indigo-500 focus:border-indigo-500"
                        >
                            <option value="">Random (Default)</option>
                            <option value="true">Force TRUE</option>
                            <option value="false">Force FALSE</option>
                        </select>
                        <p class="text-xs text-gray-500 mt-2">
                            Advanced logic coming soon via expression builder.
                        </p>
                    </div>
                </div>
            {/if}
        </div>
    {/if}
</div>
