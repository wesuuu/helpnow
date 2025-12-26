<script lang="ts">
    import { onMount } from "svelte";
    import { router } from "../lib/router.svelte";
    import { toast } from "../lib/stores/toast";

    let personId = $derived(router.path.split("/")[2]);
    let person = $state<any>(null);
    let loading = $state(true);

    async function fetchPerson(id: string) {
        if (!id) return;
        loading = true;
        try {
            const res = await fetch(`/api/people/${id}`);
            if (res.ok) {
                person = await res.json();
                if (!person) {
                    toast.error("Person not found");
                    router.navigate("/audiences");
                }
            } else {
                toast.error("Failed to load person");
            }
        } catch (e) {
            toast.error("Error loading person");
        } finally {
            loading = false;
        }
    }

    $effect(() => {
        if (personId) {
            fetchPerson(personId);
        }
    });

    // onMount not strictly needed if effect covers init, but effect runs on mount too in Svelte 5.
</script>

<div class="space-y-6">
    <div class="flex items-center space-x-4 mb-6">
        <button
            onclick={() => router.navigate("/audiences")}
            class="p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-500 dark:text-gray-400"
            aria-label="Back"
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
                    d="M10 19l-7-7m0 0l7-7m-7 7h18"
                />
            </svg>
        </button>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
            Person Details
        </h1>
    </div>

    {#if loading}
        <div class="text-center py-10">Loading...</div>
    {:else if person}
        <div
            class="bg-white dark:bg-gray-800 shadow rounded-lg overflow-hidden"
        >
            <div
                class="px-6 py-5 border-b border-gray-200 dark:border-gray-700"
            >
                <h3 class="text-lg font-medium text-gray-900 dark:text-white">
                    {person.first_name}
                    {person.last_name}
                </h3>
                <div class="text-sm text-gray-500 dark:text-gray-400">
                    {person.email}
                </div>
            </div>

            <!-- Events History -->
            <div
                class="px-6 py-5 border-t border-gray-200 dark:border-gray-700"
            >
                <h3
                    class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-4"
                >
                    Event History
                </h3>
                {#if person.events && person.events.length > 0}
                    <div class="space-y-4">
                        {#each person.events as event}
                            <div
                                class="bg-gray-50 dark:bg-gray-700 p-4 rounded-lg"
                            >
                                <div class="flex justify-between items-start">
                                    <div class="space-y-1">
                                        <pre
                                            class="text-xs text-gray-700 dark:text-gray-300 overflow-x-auto">{JSON.stringify(
                                                event.event,
                                                null,
                                                2,
                                            )}</pre>
                                    </div>
                                    <span
                                        class="text-xs text-gray-400 whitespace-nowrap ml-4"
                                    >
                                        {new Date(
                                            event.created_at,
                                        ).toLocaleString()}
                                    </span>
                                </div>
                            </div>
                        {/each}
                    </div>
                {:else}
                    <div
                        class="bg-gray-50 dark:bg-gray-700 p-4 rounded-lg text-sm text-gray-500 dark:text-gray-400"
                    >
                        No events found.
                    </div>
                {/if}
            </div>
            <div class="px-6 py-5 space-y-6">
                <!-- Demographics -->
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div>
                        <h4
                            class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-2"
                        >
                            Demographics
                        </h4>
                        <dl
                            class="grid grid-cols-1 gap-x-4 gap-y-4 sm:grid-cols-2"
                        >
                            <div class="sm:col-span-1">
                                <dt
                                    class="text-sm font-medium text-gray-500 dark:text-gray-400"
                                >
                                    Age
                                </dt>
                                <dd
                                    class="mt-1 text-sm text-gray-900 dark:text-white"
                                >
                                    {person.age || "-"}
                                </dd>
                            </div>
                            <div class="sm:col-span-1">
                                <dt
                                    class="text-sm font-medium text-gray-500 dark:text-gray-400"
                                >
                                    Gender
                                </dt>
                                <dd
                                    class="mt-1 text-sm text-gray-900 dark:text-white"
                                >
                                    {person.gender || "-"}
                                </dd>
                            </div>
                            <div class="sm:col-span-1">
                                <dt
                                    class="text-sm font-medium text-gray-500 dark:text-gray-400"
                                >
                                    Ethnicity
                                </dt>
                                <dd
                                    class="mt-1 text-sm text-gray-900 dark:text-white"
                                >
                                    {person.ethnicity || "-"}
                                </dd>
                            </div>
                            <div class="sm:col-span-1">
                                <dt
                                    class="text-sm font-medium text-gray-500 dark:text-gray-400"
                                >
                                    Location
                                </dt>
                                <dd
                                    class="mt-1 text-sm text-gray-900 dark:text-white"
                                >
                                    {person.location || "-"}
                                </dd>
                            </div>
                        </dl>
                    </div>

                    <!-- Interests -->
                    <div>
                        <h4
                            class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-2"
                        >
                            Interests
                        </h4>
                        {#if person.interests && person.interests.length > 0}
                            <div class="flex flex-wrap gap-2">
                                {#each person.interests as interest}
                                    <span
                                        class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-indigo-100 text-indigo-800 dark:bg-indigo-900 dark:text-indigo-200"
                                    >
                                        {interest}
                                    </span>
                                {/each}
                            </div>
                        {:else}
                            <p class="text-sm text-gray-500 dark:text-gray-400">
                                No interests listed.
                            </p>
                        {/if}
                    </div>
                </div>

                <!-- Meta Information -->
                {#if person.meta}
                    <div
                        class="border-t border-gray-200 dark:border-gray-700 pt-6"
                    >
                        <h4
                            class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-4"
                        >
                            System Metadata
                        </h4>
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                            <!-- IP Addresses -->
                            <div>
                                <h5
                                    class="text-sm font-medium text-gray-900 dark:text-white mb-2"
                                >
                                    IP Addresses
                                </h5>
                                {#if person.meta.ip_addresses && person.meta.ip_addresses.length > 0}
                                    <ul
                                        class="space-y-2 text-sm text-gray-600 dark:text-gray-300"
                                    >
                                        {#each person.meta.ip_addresses as ip}
                                            <li
                                                class="bg-gray-50 dark:bg-gray-700/50 p-2 rounded"
                                            >
                                                <div class="font-mono">
                                                    {ip.address}
                                                </div>
                                                <div
                                                    class="text-xs text-gray-500 mt-1"
                                                >
                                                    Used: {new Date(
                                                        ip.used,
                                                    ).toLocaleString()}
                                                    <br />
                                                    UA: {ip.user_agent}
                                                </div>
                                            </li>
                                        {/each}
                                    </ul>
                                {:else}
                                    <p
                                        class="text-sm text-gray-500 dark:text-gray-400"
                                    >
                                        No IP history.
                                    </p>
                                {/if}
                            </div>

                            <!-- User Agents -->
                            <div>
                                <h5
                                    class="text-sm font-medium text-gray-900 dark:text-white mb-2"
                                >
                                    User Agents (All)
                                </h5>
                                {#if person.meta.user_agents && person.meta.user_agents.length > 0}
                                    <ul
                                        class="space-y-2 text-sm text-gray-600 dark:text-gray-300"
                                    >
                                        {#each person.meta.user_agents as ua}
                                            <li
                                                class="bg-gray-50 dark:bg-gray-700/50 p-2 rounded break-all"
                                            >
                                                {ua}
                                            </li>
                                        {/each}
                                    </ul>
                                {:else}
                                    <p
                                        class="text-sm text-gray-500 dark:text-gray-400"
                                    >
                                        No User Agent history.
                                    </p>
                                {/if}
                            </div>
                        </div>
                    </div>
                {/if}

                <!-- Event History -->
                <div class="border-t border-gray-200 dark:border-gray-700 pt-6">
                    <h4
                        class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-4"
                    >
                        Event History
                    </h4>
                    {#if person.event_history && person.event_history.length > 0}
                        <!-- Assuming event_history is an array of objects/strings, simplistic render for now -->
                        <pre
                            class="bg-gray-50 dark:bg-gray-900 p-4 rounded text-xs overflow-auto max-h-60 text-gray-700 dark:text-gray-300">{JSON.stringify(
                                person.event_history,
                                null,
                                2,
                            )}</pre>
                    {:else}
                        <p class="text-sm text-gray-500 dark:text-gray-400">
                            No events recorded.
                        </p>
                    {/if}
                </div>
            </div>
        </div>
    {/if}
</div>
