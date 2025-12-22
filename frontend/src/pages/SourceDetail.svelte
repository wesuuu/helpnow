<script lang="ts">
    import { onMount } from "svelte";
    import { router } from "../lib/router.svelte.js";
    import Modal from "../lib/components/Modal.svelte";

    // Params handling (mock since router doesn't pass params easily yet, needs parsing)
    // Actually router stores path. typical pattern `router.path`.
    // We need to parse ID from path.
    const pathParts = router.path.split("/");
    const sourceId = pathParts[pathParts.indexOf("sources") + 1];

    interface DataSync {
        id: number;
        audience_id: number;
        sync_type: string;
        schedule: string;
        status: string;
        last_run_at: string;
        created_at: string;
    }

    interface Audience {
        id: number;
        name: string;
    }

    let syncs: DataSync[] = [];
    let audiences: Audience[] = [];
    let loading = true;
    let showModal = false;

    // Form
    let selectedAudienceId = "";
    let syncType = "ONE_OFF";
    let schedule = "0 0 * * *"; // Daily
    let query = "SELECT email, full_name, location FROM users";

    async function fetchData() {
        loading = true;
        try {
            const [syncsRes, audRes] = await Promise.all([
                fetch(`/api/syncs?source_id=${sourceId}`),
                fetch(`/api/audiences`),
            ]);

            if (syncsRes.ok) {
                syncs = await syncsRes.json();
            }
            if (audRes.ok) {
                audiences = await audRes.json();
            }
        } catch (e) {
            console.error(e);
        } finally {
            loading = false;
        }
    }

    async function createSync() {
        const payload = {
            source_id: parseInt(sourceId),
            audience_id: parseInt(selectedAudienceId),
            sync_type: syncType,
            schedule: syncType === "CRON" ? schedule : "",
            query: query,
        };

        try {
            const res = await fetch("/api/syncs", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(payload),
            });
            if (res.ok) {
                showModal = false;
                fetchData();
            }
        } catch (e) {
            console.error(e);
        }
    }

    onMount(() => {
        fetchData();
    });
</script>

<div class="max-w-6xl mx-auto">
    <div class="mb-6">
        <nav class="sm:hidden" aria-label="Back">
            <a
                href="/sources"
                on:click|preventDefault={() => router.navigate("/sources")}
                class="flex items-center text-sm font-medium text-gray-500 hover:text-gray-700"
            >
                <svg
                    class="-ml-1 mr-1 h-5 w-5 flex-shrink-0 text-gray-400"
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 20 20"
                    fill="currentColor"
                    aria-hidden="true"
                >
                    <path
                        fill-rule="evenodd"
                        d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z"
                        clip-rule="evenodd"
                    />
                </svg>
                Back
            </a>
        </nav>
        <nav class="hidden sm:flex" aria-label="Breadcrumb">
            <ol role="list" class="flex items-center space-x-4">
                <li>
                    <div class="flex">
                        <a
                            href="/sources"
                            on:click|preventDefault={() =>
                                router.navigate("/sources")}
                            class="text-sm font-medium text-gray-500 hover:text-gray-700"
                            >Data Sources</a
                        >
                    </div>
                </li>
                <li>
                    <div class="flex items-center">
                        <svg
                            class="flex-shrink-0 h-5 w-5 text-gray-300"
                            xmlns="http://www.w3.org/2000/svg"
                            fill="currentColor"
                            viewBox="0 0 20 20"
                            aria-hidden="true"
                        >
                            <path
                                d="M5.555 17.776l8-16 .894.448-8 16-.894-.448z"
                            />
                        </svg>
                        <span class="ml-4 text-sm font-medium text-gray-500"
                            >Source #{sourceId}</span
                        >
                    </div>
                </li>
            </ol>
        </nav>
    </div>

    <div class="flex justify-between items-center mb-8">
        <div>
            <h1 class="text-2xl font-bold text-gray-900">Sync Jobs</h1>
            <p class="mt-1 text-sm text-gray-500">
                Manage data synchronization tasks for this source.
            </p>
        </div>
        <button
            on:click={() => (showModal = true)}
            class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700"
        >
            Create Sync Job
        </button>
    </div>

    <div class="bg-white shadow overflow-hidden sm:rounded-lg">
        <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
                <tr>
                    <th
                        scope="col"
                        class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                        >Type</th
                    >
                    <th
                        scope="col"
                        class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                        >Target Audience</th
                    >
                    <th
                        scope="col"
                        class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                        >Schedule</th
                    >
                    <th
                        scope="col"
                        class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                        >Status</th
                    >
                    <th
                        scope="col"
                        class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                        >Last Run</th
                    >
                </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
                {#each syncs as sync}
                    <tr>
                        <td
                            class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900"
                            >{sync.sync_type}</td
                        >
                        <td
                            class="px-6 py-4 whitespace-nowrap text-sm text-gray-500"
                            >Audience #{sync.audience_id}</td
                        >
                        <td
                            class="px-6 py-4 whitespace-nowrap text-sm text-gray-500"
                            >{sync.schedule || "-"}</td
                        >
                        <td class="px-6 py-4 whitespace-nowrap">
                            <span
                                class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-blue-100 text-blue-800"
                            >
                                {sync.status}
                            </span>
                        </td>
                        <td
                            class="px-6 py-4 whitespace-nowrap text-sm text-gray-500"
                            >{sync.last_run_at
                                ? new Date(sync.last_run_at).toLocaleString()
                                : "Never"}</td
                        >
                    </tr>
                {:else}
                    <tr>
                        <td
                            colspan="5"
                            class="px-6 py-8 text-center text-gray-500"
                            >No sync jobs configured.</td
                        >
                    </tr>
                {/each}
            </tbody>
        </table>
    </div>
</div>

{#if showModal}
    <Modal onclose={() => (showModal = false)}>
        <div>
            <h3
                class="text-lg leading-6 font-medium text-gray-900"
                id="modal-title"
            >
                Create Data Sync
            </h3>
            <div class="mt-4 space-y-4">
                <div>
                    <label
                        for="audience"
                        class="block text-sm font-medium text-gray-700"
                        >Target Audience</label
                    >
                    <select
                        id="audience"
                        bind:value={selectedAudienceId}
                        class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
                    >
                        <option value="">Select Audience</option>
                        {#each audiences as aud}
                            <option value={aud.id}>{aud.name}</option>
                        {/each}
                    </select>
                </div>
                <div>
                    <label
                        for="synctype"
                        class="block text-sm font-medium text-gray-700"
                        >Sync Type</label
                    >
                    <select
                        id="synctype"
                        bind:value={syncType}
                        class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                    >
                        <option value="ONE_OFF">One-off Run</option>
                        <option value="CRON">Scheduled (Cron)</option>
                    </select>
                </div>
                {#if syncType === "CRON"}
                    <div>
                        <label
                            for="schedule"
                            class="block text-sm font-medium text-gray-700"
                            >Cron Schedule</label
                        >
                        <input
                            type="text"
                            id="schedule"
                            bind:value={schedule}
                            placeholder="0 0 * * *"
                            class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                        />
                        <p class="text-xs text-gray-500 mt-1">
                            Format: Minute Hour Day Month DayOfWeek
                        </p>
                    </div>
                {/if}
                <div>
                    <label
                        for="query"
                        class="block text-sm font-medium text-gray-700"
                        >SQL Query</label
                    >
                    <textarea
                        id="query"
                        rows="4"
                        bind:value={query}
                        class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm font-mono"
                    ></textarea>
                    <p class="text-xs text-gray-500 mt-1">
                        Must select columns mapping to Audience fields (email,
                        full_name, location).
                    </p>
                </div>
            </div>
        </div>
        <div
            class="mt-5 sm:mt-6 sm:grid sm:grid-cols-2 sm:gap-3 sm:grid-flow-row-dense"
        >
            <button
                type="button"
                on:click={createSync}
                class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:col-start-2 sm:text-sm"
            >
                Create Job
            </button>
            <button
                type="button"
                on:click={() => (showModal = false)}
                class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:col-start-1 sm:text-sm"
            >
                Cancel
            </button>
        </div>
    </Modal>
{/if}
