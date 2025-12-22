<script lang="ts">
    import { onMount } from "svelte";
    import { router } from "../lib/router.svelte.js";
    import Modal from "../lib/components/Modal.svelte";

    interface Site {
        id: number;
        organization_id: number;
        name: string;
        url: string;
        tracking_id: string;
        created_at: string;
    }

    let sites: Site[] = [];
    let showModal = false;
    let newSiteName = "";
    let newSiteUrl = "";
    let orgId = 1; // Hardcoded for now

    async function fetchSites() {
        try {
            const res = await fetch(`/api/sites?organization_id=${orgId}`);
            if (res.ok) {
                sites = await res.json();
            }
        } catch (e) {
            console.error(e);
        }
    }

    async function createSite() {
        try {
            const res = await fetch("/api/sites", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    organization_id: orgId,
                    name: newSiteName,
                    url: newSiteUrl,
                }),
            });
            if (res.ok) {
                showModal = false;
                newSiteName = "";
                newSiteUrl = "";
                fetchSites();
            }
        } catch (e) {
            console.error(e);
        }
    }

    function getSnippet(trackingId: string) {
        return (
            `<` +
            `script src="https://cdn.helpnow.com/analytics.js"><` +
            `/script>
<` +
            `script>
  const analytics = new HelpNowAnalytics('https://api.helpnow.com');
  analytics.identify('${trackingId}');
  analytics.trackImpression();
<` +
            `/script>`
        );
    }

    function copyToClipboard(text: string) {
        navigator.clipboard.writeText(text);
        alert("Copied to clipboard!");
    }

    onMount(() => {
        fetchSites();
    });
</script>

<div class="max-w-6xl mx-auto">
    <div class="flex justify-between items-center mb-8">
        <div>
            <h1 class="text-2xl font-bold text-gray-900">My Sites</h1>
            <p class="mt-1 text-sm text-gray-500">
                Manage your connected websites and landing pages.
            </p>
        </div>
        <button
            on:click={() => (showModal = true)}
            class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        >
            Register Site
        </button>
    </div>

    <div class="bg-white shadow overflow-hidden sm:rounded-md">
        <ul role="list" class="divide-y divide-gray-200">
            {#each sites as site}
                <li
                    class="relative hover:bg-gray-50 transition-colors duration-150"
                >
                    <div class="px-4 py-4 sm:px-6">
                        <div class="flex items-center justify-between">
                            <div>
                                <p
                                    class="text-lg font-medium text-indigo-600 truncate"
                                >
                                    <a
                                        href={`/sites/${site.id}`}
                                        on:click|preventDefault={() =>
                                            router.navigate(
                                                `/sites/${site.id}`,
                                            )}
                                        class="focus:outline-none"
                                    >
                                        <span
                                            class="absolute inset-0"
                                            aria-hidden="true"
                                        ></span>
                                        {site.name}
                                    </a>
                                </p>
                                <p class="text-sm text-gray-500">
                                    <a
                                        href={site.url}
                                        target="_blank"
                                        class="hover:underline relative z-10"
                                        on:click|stopPropagation>{site.url}</a
                                    >
                                </p>
                            </div>
                            <div
                                class="ml-2 flex-shrink-0 flex flex-col items-end"
                            >
                                <span
                                    class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800 mb-2"
                                >
                                    Active
                                </span>
                                {#if site.tracking_id}
                                    <div class="flex items-center space-x-2">
                                        <span
                                            class="inline-flex items-center px-2.5 py-0.5 rounded-md text-sm font-medium bg-gray-100 text-gray-800"
                                        >
                                            {site.tracking_id}
                                        </span>
                                        <button
                                            class="text-xs text-indigo-600 hover:text-indigo-900 relative z-10"
                                            on:click={() =>
                                                copyToClipboard(
                                                    getSnippet(
                                                        site.tracking_id,
                                                    ),
                                                )}
                                        >
                                            Copy Snippet
                                        </button>
                                    </div>
                                {/if}
                            </div>
                        </div>
                    </div>
                </li>
            {/each}
            {#if sites.length === 0}
                <li class="px-4 py-8 text-center text-gray-500">
                    No sites registered yet. Click "Register Site" to add one.
                </li>
            {/if}
        </ul>
    </div>
</div>

{#if showModal}
    <Modal onclose={() => (showModal = false)}>
        <div>
            <div
                class="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-indigo-100"
            >
                <svg
                    class="h-6 w-6 text-indigo-600"
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                </svg>
            </div>
            <div class="mt-3 text-center sm:mt-5">
                <h3
                    class="text-lg leading-6 font-medium text-gray-900"
                    id="modal-title"
                >
                    Register New Site
                </h3>
                <div class="mt-2">
                    <p class="text-sm text-gray-500">
                        Register a new website or landing page to track.
                    </p>
                </div>
                <div class="mt-4 space-y-4">
                    <div>
                        <label
                            for="site-name"
                            class="block text-sm font-medium text-gray-700 text-left"
                            >Site Name</label
                        >
                        <div class="mt-1">
                            <input
                                type="text"
                                id="site-name"
                                bind:value={newSiteName}
                                class="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md"
                                placeholder="e.g. Main Corporate Site"
                            />
                        </div>
                    </div>
                    <div>
                        <label
                            for="site-url"
                            class="block text-sm font-medium text-gray-700 text-left"
                            >Site URL</label
                        >
                        <div class="mt-1">
                            <input
                                type="text"
                                id="site-url"
                                bind:value={newSiteUrl}
                                class="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md"
                                placeholder="e.g. https://example.com"
                            />
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div
            class="mt-5 sm:mt-6 sm:grid sm:grid-cols-2 sm:gap-3 sm:grid-flow-row-dense"
        >
            <button
                type="button"
                class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:col-start-2 sm:text-sm"
                on:click={createSite}
            >
                Register
            </button>
            <button
                type="button"
                class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:col-start-1 sm:text-sm"
                on:click={() => (showModal = false)}
            >
                Cancel
            </button>
        </div>
    </Modal>
{/if}
