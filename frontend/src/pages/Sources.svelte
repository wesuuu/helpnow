<script lang="ts">
    import { onMount } from "svelte";
    import { router } from "../lib/router.svelte.js";
    import Modal from "../lib/components/Modal.svelte";

    interface DataSource {
        id: number;
        name: string;
        type: string;
        created_at: string;
    }

    let sources = $state<DataSource[]>([]);
    let loading = $state(true);
    let showModal = $state(false);
    let modalStep = $state("SELECT_TYPE"); // 'SELECT_TYPE' | 'CONFIGURE'

    // Form
    let name = $state("");
    let type = $state("");
    let config = $state<any>({});

    async function fetchSources() {
        loading = true;
        try {
            const res = await fetch("/api/sources");
            if (res.ok) {
                sources = await res.json();
            }
        } catch (e) {
            console.error(e);
        } finally {
            loading = false;
        }
    }

    function selectType(selectedType: string) {
        type = selectedType;
        modalStep = "CONFIGURE";
        // Reset defaults
        name = "";
        config = {};
    }

    async function createSource() {
        const payload = {
            name,
            type,
            config: JSON.stringify(config),
        };
        try {
            const res = await fetch("/api/sources", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(payload),
            });
            if (res.ok) {
                closeModal();
                fetchSources();
            }
        } catch (e) {
            console.error(e);
        }
    }

    function closeModal() {
        showModal = false;
        modalStep = "SELECT_TYPE";
        type = "";
    }

    onMount(() => {
        fetchSources();
    });
</script>

<div class="max-w-6xl mx-auto">
    <div class="flex justify-between items-center mb-8">
        <div>
            <h1 class="text-2xl font-bold text-gray-900">Data Sources</h1>
            <p class="mt-1 text-sm text-gray-500">
                Connect your data sources to sync audiences.
            </p>
        </div>
        <button
            onclick={() => (showModal = true)}
            class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700"
        >
            Add Source
        </button>
    </div>

    {#if loading}
        <div>Loading...</div>
    {:else}
        <div class="bg-white shadow overflow-hidden sm:rounded-md">
            <ul role="list" class="divide-y divide-gray-200">
                {#each sources as source}
                    <li class="px-4 py-4 sm:px-6 hover:bg-gray-50 transition">
                        <div class="flex items-center justify-between">
                            <div class="flex items-center gap-4">
                                <div class="bg-indigo-100 rounded-md p-2">
                                    <svg
                                        xmlns="http://www.w3.org/2000/svg"
                                        class="h-6 w-6 text-indigo-600"
                                        fill="none"
                                        viewBox="0 0 24 24"
                                        stroke="currentColor"
                                    >
                                        <path
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                            stroke-width="2"
                                            d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4"
                                        />
                                    </svg>
                                </div>
                                <div>
                                    <p
                                        class="text-lg font-medium text-indigo-600 truncate"
                                    >
                                        {source.name}
                                    </p>
                                    <p class="text-sm text-gray-500">
                                        {source.type}
                                    </p>
                                </div>
                            </div>
                            <div class="flex gap-2">
                                <button
                                    onclick={() =>
                                        router.navigate(
                                            `/sources/${source.id}`,
                                        )}
                                    class="text-indigo-600 hover:text-indigo-900 text-sm font-medium"
                                >
                                    Syncs
                                </button>
                            </div>
                        </div>
                    </li>
                {:else}
                    <li class="px-4 py-8 text-center text-gray-500">
                        No data sources connected.
                    </li>
                {/each}
            </ul>
        </div>
    {/if}
</div>

{#if showModal}
    <Modal onclose={closeModal}>
        {#if modalStep === "SELECT_TYPE"}
            <div>
                <h3
                    class="text-lg leading-6 font-medium text-gray-900"
                    id="modal-title"
                >
                    Select Integration
                </h3>
                <div class="mt-4 grid grid-cols-1 gap-4">
                    <!-- Postgres Card -->
                    <button
                        onclick={() => selectType("POSTGRES")}
                        class="flex items-center p-4 border border-gray-300 rounded-lg hover:border-indigo-500 hover:ring-1 hover:ring-indigo-500 transition-all text-left group"
                    >
                        <div
                            class="bg-blue-100 p-3 rounded-md group-hover:bg-blue-200"
                        >
                            <!-- Database Icon -->
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="h-6 w-6 text-blue-600"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4"
                                />
                            </svg>
                        </div>
                        <div class="ml-4">
                            <p class="text-sm font-medium text-gray-900">
                                PostgreSQL
                            </p>
                            <p class="text-sm text-gray-500">
                                Connect to a Postgres database.
                            </p>
                        </div>
                    </button>

                    <!-- Databricks Card -->
                    <button
                        onclick={() => selectType("DATABRICKS")}
                        class="flex items-center p-4 border border-gray-300 rounded-lg hover:border-indigo-500 hover:ring-1 hover:ring-indigo-500 transition-all text-left group"
                    >
                        <div
                            class="bg-orange-100 p-3 rounded-md group-hover:bg-orange-200"
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="h-6 w-6 text-orange-600"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
                                />
                            </svg>
                        </div>
                        <div class="ml-4">
                            <p class="text-sm font-medium text-gray-900">
                                Databricks
                            </p>
                            <p class="text-sm text-gray-500">
                                Connect to a Databricks cluster.
                            </p>
                        </div>
                    </button>

                    <!-- Webhook Card -->
                    <button
                        onclick={() => selectType("WEBHOOK")}
                        class="flex items-center p-4 border border-gray-300 rounded-lg hover:border-indigo-500 hover:ring-1 hover:ring-indigo-500 transition-all text-left group"
                    >
                        <div
                            class="bg-green-100 p-3 rounded-md group-hover:bg-green-200"
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="h-6 w-6 text-green-600"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M13 10V3L4 14h7v7l9-11h-7z"
                                />
                            </svg>
                        </div>
                        <div class="ml-4">
                            <p class="text-sm font-medium text-gray-900">
                                Webhook
                            </p>
                            <p class="text-sm text-gray-500">
                                Push data via HTTP webhooks.
                            </p>
                        </div>
                    </button>
                </div>
            </div>
        {:else if modalStep === "CONFIGURE"}
            <div>
                <h3
                    class="text-lg leading-6 font-medium text-gray-900"
                    id="modal-title"
                >
                    Configure {type === "POSTGRES"
                        ? "PostgreSQL"
                        : type === "DATABRICKS"
                          ? "Databricks"
                          : "Webhook"}
                </h3>
                <div class="mt-4 space-y-4">
                    <div>
                        <label
                            for="name"
                            class="block text-sm font-medium text-gray-700"
                            >Name</label
                        >
                        <input
                            type="text"
                            id="name"
                            bind:value={name}
                            class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                            placeholder="My Data Source"
                        />
                    </div>

                    {#if type === "POSTGRES" || type === "DATABRICKS"}
                        <div>
                            <label
                                for="host"
                                class="block text-sm font-medium text-gray-700"
                                >Host</label
                            >
                            <input
                                type="text"
                                id="host"
                                bind:value={config.host}
                                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                            />
                        </div>
                        <div>
                            <label
                                for="port"
                                class="block text-sm font-medium text-gray-700"
                                >Port</label
                            >
                            <input
                                type="number"
                                id="port"
                                bind:value={config.port}
                                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                            />
                        </div>
                        <div>
                            <label
                                for="user"
                                class="block text-sm font-medium text-gray-700"
                                >User</label
                            >
                            <input
                                type="text"
                                id="user"
                                bind:value={config.user}
                                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                            />
                        </div>
                        <div>
                            <label
                                for="password"
                                class="block text-sm font-medium text-gray-700"
                                >Password</label
                            >
                            <input
                                type="password"
                                id="password"
                                bind:value={config.password}
                                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                            />
                        </div>
                        <div>
                            <label
                                for="dbname"
                                class="block text-sm font-medium text-gray-700"
                                >Database Name</label
                            >
                            <input
                                type="text"
                                id="dbname"
                                bind:value={config.dbname}
                                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                            />
                        </div>
                    {:else if type === "WEBHOOK"}
                        <div>
                            <p class="text-sm text-gray-500">
                                A webhook URL will be generated after creation.
                            </p>
                        </div>
                    {/if}
                </div>
            </div>
            <div
                class="mt-5 sm:mt-6 sm:grid sm:grid-cols-2 sm:gap-3 sm:grid-flow-row-dense"
            >
                <button
                    type="button"
                    onclick={createSource}
                    class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:col-start-2 sm:text-sm"
                >
                    Create
                </button>
                <button
                    type="button"
                    onclick={() => (modalStep = "SELECT_TYPE")}
                    class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:col-start-1 sm:text-sm"
                >
                    Back
                </button>
            </div>
        {/if}
    </Modal>
{/if}
