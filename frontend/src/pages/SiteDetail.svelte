<script lang="ts">
    import { onMount } from "svelte";
    import { router } from "../lib/router.svelte.js";

    interface Site {
        id: number;
        organization_id: number;
        name: string;
        url: string;
        tracking_id: string;
        created_at: string;
    }

    interface SiteStats {
        active: boolean;
        last_event_at: string;
        stats: {
            impressions: number;
            clicks: number;
            conversions: number;
            cost: number;
        };
    }

    interface Workflow {
        id: number;
        name: string;
        trigger_event: string;
        steps: string;
        status: string;
        created_at: string;
    }

    interface EventDefinition {
        id: number;
        name: string;
        description: string;
    }

    let site: Site | null = null;
    let siteStats: SiteStats | null = null;
    let workflows: Workflow[] = [];
    let eventDefinitions: EventDefinition[] = [];
    let loading = true;
    let error = "";
    let activeTab = "overview";

    // Workflow Builder State
    let showCreateWorkflow = false;
    let newWorkflowName = "";
    let newWorkflowTrigger = "";
    let newWorkflowSteps = [{ action: "", delay: { days: 0, hours: 0 } }];

    // Event Definition State
    let showCreateEvent = false;
    let newEventName = "";
    let newEventDesc = "";

    // Extract ID from path manually since router params aren't passed automatically yet
    const siteId = router.path.split("/").pop();

    async function fetchSite() {
        try {
            const [siteRes, statsRes] = await Promise.all([
                fetch(`/api/sites/${siteId}`),
                fetch(`/api/sites/${siteId}/stats`),
            ]);

            if (siteRes.ok) {
                site = await siteRes.json();
            } else {
                error = "Site not found";
            }

            if (statsRes.ok) {
                siteStats = await statsRes.json();
            }

            // Fetch Workflows and Events
            fetchWorkflows();
            fetchEvents();
        } catch (e) {
            error = "Failed to load site data";
            console.error(e);
        } finally {
            loading = false;
        }
    }

    function getSnippet(trackingId: string) {
        // Break script tags to avoid parser issues
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

    async function fetchWorkflows() {
        const res = await fetch(`/api/workflows?site_id=${siteId}`);
        if (res.ok) {
            workflows = await res.json();
        }
    }

    async function fetchEvents() {
        const res = await fetch(`/api/events/definitions?site_id=${siteId}`);
        if (res.ok) {
            eventDefinitions = await res.json();
        }
    }

    async function createWorkflow() {
        if (
            !newWorkflowName ||
            !newWorkflowTrigger ||
            newWorkflowSteps.length === 0
        )
            return;

        const payload = {
            site_id: parseInt(siteId as string),
            name: newWorkflowName,
            trigger_type: "EVENT",
            trigger_event: newWorkflowTrigger,
            steps: JSON.stringify(newWorkflowSteps),
            status: "ACTIVE",
        };

        const res = await fetch("/api/workflows", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload),
        });

        if (res.ok) {
            showCreateWorkflow = false;
            newWorkflowName = "";
            newWorkflowTrigger = "";
            newWorkflowSteps = [{ action: "", delay: { days: 0, hours: 0 } }];
            fetchWorkflows();
        }
    }

    function addStep() {
        newWorkflowSteps = [
            ...newWorkflowSteps,
            { action: "", delay: { days: 0, hours: 0 } },
        ];
    }

    function removeStep(index: number) {
        newWorkflowSteps = newWorkflowSteps.filter((_, i) => i !== index);
    }

    async function createEvent() {
        if (!newEventName) return;

        const payload = {
            site_id: parseInt(siteId as string),
            name: newEventName,
            description: newEventDesc,
        };

        const res = await fetch("/api/events/definitions", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload),
        });

        if (res.ok) {
            showCreateEvent = false;
            newEventName = "";
            newEventDesc = "";
            fetchEvents();
        }
    }

    function copyToClipboard(text: string) {
        navigator.clipboard.writeText(text);
        alert("Copied to clipboard!");
    }

    onMount(() => {
        if (siteId) {
            fetchSite();
        } else {
            error = "Invalid Site ID";
            loading = false;
        }
    });
</script>

<div class="max-w-4xl mx-auto">
    {#if loading}
        <div class="text-center py-12">Loading...</div>
    {:else if error}
        <div class="text-center py-12 text-red-600">{error}</div>
    {:else if site}
        <div class="mb-4">
            <a
                href="/sites"
                on:click|preventDefault={() => router.navigate("/sites")}
                class="text-indigo-600 hover:text-indigo-900 mb-4 inline-block"
                >&larr; Back to Sites</a
            >
            <div class="flex items-start justify-between">
                <div>
                    <h1
                        class="text-3xl font-bold text-gray-900 flex items-center gap-3"
                    >
                        {site.name}
                        {#if siteStats}
                            <span
                                class={`px-2.5 py-0.5 rounded-full text-xs font-medium ${siteStats.active ? "bg-green-100 text-green-800" : "bg-gray-100 text-gray-800"}`}
                            >
                                {siteStats.active
                                    ? "Integration Active"
                                    : "Pending Data"}
                            </span>
                        {/if}
                    </h1>
                    <p class="mt-2 text-sm text-gray-500">
                        <a
                            href={site.url}
                            target="_blank"
                            class="hover:underline">{site.url}</a
                        >
                    </p>
                </div>
            </div>
        </div>

        <div class="border-b border-gray-200 mb-6">
            <nav class="-mb-px flex space-x-8" aria-label="Tabs">
                <button
                    on:click={() => (activeTab = "overview")}
                    class="{activeTab === 'overview'
                        ? 'border-indigo-500 text-indigo-600'
                        : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}
                                whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm focus:outline-none"
                >
                    Overview
                </button>
                <button
                    on:click={() => (activeTab = "stats")}
                    class="{activeTab === 'stats'
                        ? 'border-indigo-500 text-indigo-600'
                        : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}
                                whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm focus:outline-none"
                >
                    Stats
                </button>
                <button
                    on:click={() => (activeTab = "workflows")}
                    class="{activeTab === 'workflows'
                        ? 'border-indigo-500 text-indigo-600'
                        : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}
                                whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm focus:outline-none"
                >
                    Workflows
                </button>
            </nav>
        </div>

        {#if activeTab === "overview"}
            <div class="bg-white shadow overflow-hidden sm:rounded-lg mb-8">
                <div class="px-4 py-5 sm:px-6">
                    <h3 class="text-lg leading-6 font-medium text-gray-900">
                        Integration
                    </h3>
                    <p class="mt-1 max-w-2xl text-sm text-gray-500">
                        Add this code to your website's <code>&lt;head&gt;</code
                        >.
                    </p>
                </div>
                <div class="border-t border-gray-200 px-4 py-5 sm:px-6">
                    <div class="mb-4">
                        <h4 class="block text-sm font-medium text-gray-700">
                            Tracking ID
                        </h4>
                        <div class="mt-1 flex items-center">
                            <span
                                class="inline-flex items-center px-3 py-1 rounded-md text-sm font-medium bg-gray-100 text-gray-800"
                            >
                                {site!.tracking_id}
                            </span>
                        </div>
                    </div>

                    <div>
                        <div class="flex justify-between items-center mb-2">
                            <h4 class="block text-sm font-medium text-gray-700">
                                Code Snippet
                            </h4>
                            <button
                                class="text-sm text-indigo-600 hover:text-indigo-900"
                                on:click={() =>
                                    copyToClipboard(
                                        getSnippet(site!.tracking_id),
                                    )}
                            >
                                Copy Code
                            </button>
                        </div>
                        <pre
                            class="bg-gray-800 text-gray-100 p-4 rounded-md overflow-x-auto text-sm">{getSnippet(
                                site!.tracking_id,
                            )}</pre>
                    </div>
                </div>
            </div>

            <div class="bg-white shadow overflow-hidden sm:rounded-lg">
                <div class="px-4 py-5 sm:px-6">
                    <h3 class="text-lg leading-6 font-medium text-gray-900">
                        Recent Activity
                    </h3>
                    <p class="mt-1 max-w-2xl text-sm text-gray-500">
                        Live capability coming soon.
                    </p>
                </div>
                <div
                    class="border-t border-gray-200 px-4 py-12 sm:px-6 text-center text-gray-500"
                >
                    No recent data available.
                </div>
            </div>
        {:else if activeTab === "stats" && siteStats}
            <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
                <!-- Impressions -->
                <div class="bg-white overflow-hidden shadow rounded-lg">
                    <div class="px-4 py-5 sm:p-6">
                        <dt class="text-sm font-medium text-gray-500 truncate">
                            Total Impressions
                        </dt>
                        <dd class="mt-1 text-3xl font-semibold text-gray-900">
                            {siteStats.stats.impressions.toLocaleString()}
                        </dd>
                    </div>
                </div>
                <!-- Clicks -->
                <div class="bg-white overflow-hidden shadow rounded-lg">
                    <div class="px-4 py-5 sm:p-6">
                        <dt class="text-sm font-medium text-gray-500 truncate">
                            Total Clicks
                        </dt>
                        <dd class="mt-1 text-3xl font-semibold text-gray-900">
                            {siteStats.stats.clicks.toLocaleString()}
                        </dd>
                    </div>
                </div>
                <!-- Conversions -->
                <div class="bg-white overflow-hidden shadow rounded-lg">
                    <div class="px-4 py-5 sm:p-6">
                        <dt class="text-sm font-medium text-gray-500 truncate">
                            Conversions
                        </dt>
                        <dd class="mt-1 text-3xl font-semibold text-gray-900">
                            {siteStats.stats.conversions.toLocaleString()}
                        </dd>
                    </div>
                </div>
                <!-- Cost -->
                <div class="bg-white overflow-hidden shadow rounded-lg">
                    <div class="px-4 py-5 sm:p-6">
                        <dt class="text-sm font-medium text-gray-500 truncate">
                            Total Cost
                        </dt>
                        <dd class="mt-1 text-3xl font-semibold text-gray-900">
                            ${siteStats.stats.cost.toLocaleString()}
                        </dd>
                    </div>
                </div>
            </div>
            {#if !siteStats.active}
                <div class="rounded-md bg-yellow-50 p-4 mt-6">
                    <div class="flex">
                        <div class="flex-shrink-0">
                            <!-- Heroicon name: solid/exclamation -->
                            <svg
                                class="h-5 w-5 text-yellow-400"
                                xmlns="http://www.w3.org/2000/svg"
                                viewBox="0 0 20 20"
                                fill="currentColor"
                                aria-hidden="true"
                            >
                                <path
                                    fill-rule="evenodd"
                                    d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
                                    clip-rule="evenodd"
                                />
                            </svg>
                        </div>
                        <div class="ml-3">
                            <h3 class="text-sm font-medium text-yellow-800">
                                No data received yet
                            </h3>
                            <div class="mt-2 text-sm text-yellow-700">
                                <p>
                                    Once you install the tracking snippet on
                                    your site, data will appear here.
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
            {/if}
        {:else if activeTab === "workflows"}
            <div class="bg-white shadow overflow-hidden sm:rounded-lg mb-8">
                <div
                    class="px-4 py-5 sm:px-6 flex justify-between items-center"
                >
                    <div>
                        <h3 class="text-lg leading-6 font-medium text-gray-900">
                            Workflows
                        </h3>
                        <p class="mt-1 max-w-2xl text-sm text-gray-500">
                            Automate actions based on events.
                        </p>
                    </div>
                    <div class="flex space-x-2">
                        <button
                            on:click={() =>
                                (showCreateEvent = !showCreateEvent)}
                            class="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none"
                        >
                            Register Event
                        </button>
                        <button
                            on:click={() =>
                                (showCreateWorkflow = !showCreateWorkflow)}
                            class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none"
                        >
                            Create Workflow
                        </button>
                    </div>
                </div>

                {#if showCreateEvent}
                    <div class="bg-gray-50 p-4 border-t border-gray-200">
                        <h4 class="text-sm font-medium text-gray-900 mb-3">
                            Define New Event
                        </h4>
                        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
                            <div>
                                <label
                                    for="event-name"
                                    class="block text-sm font-medium text-gray-700"
                                    >Event Name</label
                                >
                                <input
                                    id="event-name"
                                    type="text"
                                    bind:value={newEventName}
                                    placeholder="e.g. purchased_item"
                                    class="mt-1 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
                                />
                            </div>
                            <div>
                                <label
                                    for="event-desc"
                                    class="block text-sm font-medium text-gray-700"
                                    >Description</label
                                >
                                <input
                                    id="event-desc"
                                    type="text"
                                    bind:value={newEventDesc}
                                    placeholder="Optional description"
                                    class="mt-1 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md"
                                />
                            </div>
                        </div>
                        <div class="mt-4 flex justify-end">
                            <button
                                on:click={createEvent}
                                class="px-4 py-2 bg-indigo-600 text-white rounded-md text-sm"
                                >Save Event</button
                            >
                        </div>
                    </div>
                {/if}

                {#if showCreateWorkflow}
                    <div class="bg-indigo-50 p-4 border-t border-indigo-100">
                        <h4 class="text-sm font-medium text-indigo-900 mb-3">
                            New Workflow
                        </h4>
                        <div class="space-y-4">
                            <div>
                                <label
                                    for="workflow-name"
                                    class="block text-sm font-medium text-gray-700"
                                    >Name</label
                                >
                                <input
                                    id="workflow-name"
                                    type="text"
                                    bind:value={newWorkflowName}
                                    class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                                />
                            </div>
                            <div>
                                <label
                                    for="workflow-trigger"
                                    class="block text-sm font-medium text-gray-700"
                                    >Trigger Event</label
                                >
                                <select
                                    id="workflow-trigger"
                                    bind:value={newWorkflowTrigger}
                                    class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
                                >
                                    <option value="">Select an event...</option>
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
                            <div>
                                <h5
                                    class="block text-sm font-medium text-gray-700"
                                >
                                    Steps
                                </h5>
                                <div class="space-y-3 mt-2">
                                    {#each newWorkflowSteps as step, i}
                                        <div
                                            class="flex gap-2 items-start border p-3 rounded-md bg-white"
                                        >
                                            <div class="flex-1 space-y-2">
                                                <div class="flex-1">
                                                    <label
                                                        for={`step-action-${i}`}
                                                        class="block text-xs font-medium text-gray-500"
                                                        >Action</label
                                                    >
                                                    <input
                                                        id={`step-action-${i}`}
                                                        type="text"
                                                        bind:value={step.action}
                                                        placeholder="e.g. Send Email"
                                                        class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                                                    />
                                                </div>
                                                <div class="flex gap-2">
                                                    <div class="w-24">
                                                        <label
                                                            for={`step-days-${i}`}
                                                            class="block text-xs font-medium text-gray-500"
                                                            >Delay Days</label
                                                        >
                                                        <input
                                                            id={`step-days-${i}`}
                                                            type="number"
                                                            min="0"
                                                            bind:value={
                                                                step.delay.days
                                                            }
                                                            class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                                                        />
                                                    </div>
                                                    <div class="w-24">
                                                        <label
                                                            for={`step-hours-${i}`}
                                                            class="block text-xs font-medium text-gray-500"
                                                            >Delay Hours</label
                                                        >
                                                        <input
                                                            id={`step-hours-${i}`}
                                                            type="number"
                                                            min="0"
                                                            bind:value={
                                                                step.delay.hours
                                                            }
                                                            class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                                                        />
                                                    </div>
                                                </div>
                                            </div>
                                            {#if newWorkflowSteps.length > 1}
                                                <button
                                                    on:click={() =>
                                                        removeStep(i)}
                                                    class="text-red-600 hover:text-red-900 mt-6"
                                                    aria-label="Remove Step"
                                                >
                                                    <svg
                                                        xmlns="http://www.w3.org/2000/svg"
                                                        class="h-5 w-5"
                                                        viewBox="0 0 20 20"
                                                        fill="currentColor"
                                                    >
                                                        <path
                                                            fill-rule="evenodd"
                                                            d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 000-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z"
                                                            clip-rule="evenodd"
                                                        />
                                                    </svg>
                                                </button>
                                            {/if}
                                        </div>
                                    {/each}
                                </div>
                                <button
                                    on:click={addStep}
                                    class="mt-2 text-sm text-indigo-600 hover:text-indigo-900 font-medium"
                                    >+ Add Step</button
                                >
                            </div>
                            <div class="flex justify-end">
                                <button
                                    on:click={createWorkflow}
                                    class="px-4 py-2 bg-indigo-600 text-white rounded-md text-sm"
                                    >Create Workflow</button
                                >
                            </div>
                        </div>
                    </div>
                {/if}

                <ul class="divide-y divide-gray-200 border-t border-gray-200">
                    {#each workflows as wf}
                        <li class="px-4 py-4 sm:px-6 hover:bg-gray-50">
                            <div class="flex items-center justify-between">
                                <div
                                    class="text-sm font-medium text-indigo-600 truncate"
                                >
                                    {wf.name}
                                </div>
                                <div class="ml-2 flex-shrink-0 flex">
                                    <span
                                        class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800"
                                    >
                                        {wf.status}
                                    </span>
                                </div>
                            </div>
                            <div class="mt-2 text-sm text-gray-500">
                                <p>
                                    Triggers on: <span
                                        class="font-medium text-gray-900"
                                        >{wf.trigger_event}</span
                                    >
                                </p>
                            </div>
                        </li>
                    {:else}
                        <li class="px-4 py-8 text-center text-gray-500 text-sm">
                            No workflows defined.
                        </li>
                    {/each}
                </ul>
            </div>
        {/if}
    {/if}
</div>
