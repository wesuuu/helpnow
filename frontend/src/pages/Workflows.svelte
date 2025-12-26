<script lang="ts">
    import { onMount } from "svelte";
    import { router } from "../lib/router.svelte.js";
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

    let workflows = $state<Workflow[]>([]);
    let loading = $state(true);

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

    function createWorkflow() {
        router.navigate("/workflows/new");
    }

    function editWorkflow(wf: Workflow) {
        router.navigate(`/workflows/${wf.id}`);
    }

    onMount(() => {
        loadWorkflows();
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
            onclick={createWorkflow}
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
