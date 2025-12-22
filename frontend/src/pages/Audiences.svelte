<script lang="ts">
    import { onMount } from "svelte";
    import Modal from "../lib/components/Modal.svelte";

    interface Audience {
        id: number;
        name: string;
        description?: string;
    }

    interface Person {
        id: number;
        full_name: string;
        email?: string;
        age?: number;
        gender?: string;
        ethnicity?: string;
        location?: string;
        created_at: string;
    }

    let activeTab = $state("audiences"); // 'audiences' or 'people'

    // Audiences State
    let audiences = $state<Audience[]>([]);
    let selectedAudience = $state<Audience | null>(null);
    let audienceMembers = $state<Person[]>([]);
    let showAddMemberModal = $state(false);

    // People State
    let people = $state<Person[]>([]);
    let showAddPersonModal = $state(false);
    let newPerson = $state({
        full_name: "",
        email: "",
        age: "",
        gender: "",
        ethnicity: "",
        location: "",
    });

    // Add Member State
    let selectedPersonIdToAdd = $state("");

    async function fetchAudiences() {
        // TODO: Use real org ID
        const res = await fetch(`/api/audiences?organization_id=1`);
        if (res.ok) audiences = await res.json();
    }

    async function fetchPeople() {
        // TODO: Use real org ID
        const res = await fetch(`/api/people?organization_id=1`);
        if (res.ok) people = await res.json();
    }

    async function fetchAudienceMembers(audienceId: number) {
        // TODO: Use real org ID
        const res = await fetch(
            `/api/audiences/${audienceId}/members?organization_id=1`,
        );
        if (res.ok) audienceMembers = await res.json();
    }

    async function createPerson() {
        const res = await fetch("/api/people", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                organization_id: 1, // Hardcoded for MVP
                ...newPerson,
                age: newPerson.age ? parseInt(newPerson.age) : null,
            }),
        });

        if (res.ok) {
            showAddPersonModal = false;
            fetchPeople();
            newPerson = {
                full_name: "",
                email: "",
                age: "",
                gender: "",
                ethnicity: "",
                location: "",
            };
        }
    }

    async function addPersonToAudience() {
        if (!selectedAudience || !selectedPersonIdToAdd) return;

        const res = await fetch(
            `/api/audiences/${selectedAudience.id}/members`,
            {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    person_id: parseInt(selectedPersonIdToAdd),
                }),
            },
        );

        if (res.ok) {
            showAddMemberModal = false;
            fetchAudienceMembers(selectedAudience.id);
            selectedPersonIdToAdd = "";
        }
    }

    function selectAudience(audience: Audience) {
        selectedAudience = audience;
        fetchAudienceMembers(audience.id);
    }

    onMount(() => {
        fetchAudiences();
        fetchPeople();
    });
</script>

<div class="space-y-6">
    <div
        class="flex justify-between items-center border-b border-gray-200 dark:border-gray-700 pb-4"
    >
        <div class="flex space-x-4">
            <button
                class={`px-4 py-2 font-medium rounded-md ${activeTab === "audiences" ? "bg-indigo-100 text-indigo-700 dark:bg-indigo-900 dark:text-indigo-200" : "text-gray-500 hover:text-gray-700 dark:text-gray-400"}`}
                onclick={() => (activeTab = "audiences")}
            >
                My Audiences
            </button>
            <button
                class={`px-4 py-2 font-medium rounded-md ${activeTab === "people" ? "bg-indigo-100 text-indigo-700 dark:bg-indigo-900 dark:text-indigo-200" : "text-gray-500 hover:text-gray-700 dark:text-gray-400"}`}
                onclick={() => (activeTab = "people")}
            >
                People Database
            </button>
        </div>
    </div>

    <!-- AUDIENCES TAB -->
    {#if activeTab === "audiences"}
        <div class="flex gap-6 h-[calc(100vh-200px)]">
            <!-- List -->
            <div
                class="w-1/3 border-r border-gray-200 dark:border-gray-700 pr-6 overflow-y-auto"
            >
                <h3
                    class="text-lg font-bold mb-4 text-gray-900 dark:text-white"
                >
                    Audiences
                </h3>
                <div class="space-y-2">
                    {#each audiences as audience}
                        <button
                            onclick={() => selectAudience(audience)}
                            class={`w-full text-left p-4 rounded-lg border transition ${selectedAudience?.id === audience.id ? "border-indigo-500 bg-indigo-50 dark:bg-indigo-900/20" : "border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800"}`}
                        >
                            <div
                                class="font-medium text-gray-900 dark:text-white"
                            >
                                {audience.name}
                            </div>
                            <div class="text-xs text-gray-500 mt-1">
                                {audience.description || "No description"}
                            </div>
                        </button>
                    {/each}
                    {#if audiences.length === 0}
                        <div class="text-gray-500 text-sm">
                            No audiences created yet.
                        </div>
                    {/if}
                </div>
            </div>

            <!-- Details -->
            <div class="w-2/3 pl-2 overflow-y-auto">
                {#if selectedAudience}
                    <div class="flex justify-between items-center mb-6">
                        <h3
                            class="text-xl font-bold text-gray-900 dark:text-white"
                        >
                            {selectedAudience.name}
                            <span class="text-gray-400 font-normal text-base"
                                >Members</span
                            >
                        </h3>
                        <button
                            onclick={() => (showAddMemberModal = true)}
                            class="text-sm bg-indigo-600 text-white px-3 py-1.5 rounded hover:bg-indigo-700"
                        >
                            + Add Person
                        </button>
                    </div>

                    <div
                        class="bg-white dark:bg-gray-800 shadow rounded-lg overflow-hidden"
                    >
                        <table
                            class="min-w-full divide-y divide-gray-200 dark:divide-gray-700"
                        >
                            <thead class="bg-gray-50 dark:bg-gray-900">
                                <tr>
                                    <th
                                        class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                                        >Name</th
                                    >
                                    <th
                                        class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                                        >Email</th
                                    >
                                    <th
                                        class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                                        >Location</th
                                    >
                                </tr>
                            </thead>
                            <tbody
                                class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700"
                            >
                                {#each audienceMembers as member}
                                    <tr>
                                        <td
                                            class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white"
                                            >{member.full_name}</td
                                        >
                                        <td
                                            class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300"
                                            >{member.email || "-"}</td
                                        >
                                        <td
                                            class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300"
                                            >{member.location || "-"}</td
                                        >
                                    </tr>
                                {/each}
                                {#if audienceMembers.length === 0}
                                    <tr>
                                        <td
                                            colspan="3"
                                            class="px-6 py-4 text-center text-sm text-gray-500"
                                            >No members in this audience yet.</td
                                        >
                                    </tr>
                                {/if}
                            </tbody>
                        </table>
                    </div>
                {:else}
                    <div
                        class="flex h-full items-center justify-center text-gray-400"
                    >
                        Select an audience to view members.
                    </div>
                {/if}
            </div>
        </div>
    {/if}

    <!-- PEOPLE TAB -->
    {#if activeTab === "people"}
        <div>
            <div class="flex justify-between items-center mb-4">
                <h3 class="text-lg font-bold text-gray-900 dark:text-white">
                    All People
                </h3>
                <button
                    onclick={() => (showAddPersonModal = true)}
                    class="bg-green-600 text-white px-4 py-2 rounded-md hover:bg-green-700"
                >
                    Add New Person
                </button>
            </div>

            <div
                class="bg-white dark:bg-gray-800 shadow rounded-lg overflow-hidden"
            >
                <table
                    class="min-w-full divide-y divide-gray-200 dark:divide-gray-700"
                >
                    <thead class="bg-gray-50 dark:bg-gray-900">
                        <tr>
                            <th
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                                >Name</th
                            >
                            <th
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                                >Demographics</th
                            >
                            <th
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                                >Location</th
                            >
                            <th
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                                >Created</th
                            >
                        </tr>
                    </thead>
                    <tbody
                        class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700"
                    >
                        {#each people as person}
                            <tr>
                                <td class="px-6 py-4 whitespace-nowrap">
                                    <div
                                        class="text-sm font-medium text-gray-900 dark:text-white"
                                    >
                                        {person.full_name}
                                    </div>
                                    <div class="text-sm text-gray-500">
                                        {person.email}
                                    </div>
                                </td>
                                <td
                                    class="px-6 py-4 text-sm text-gray-500 dark:text-gray-300"
                                >
                                    {[
                                        person.age ? `${person.age}yo` : "",
                                        person.gender,
                                        person.ethnicity,
                                    ]
                                        .filter(Boolean)
                                        .join(", ") || "-"}
                                </td>
                                <td
                                    class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300"
                                    >{person.location || "-"}</td
                                >
                                <td
                                    class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300"
                                    >{new Date(
                                        person.created_at,
                                    ).toLocaleDateString()}</td
                                >
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
        </div>
    {/if}

    <!-- MODALS -->
    {#if showAddPersonModal}
        <Modal onclose={() => (showAddPersonModal = false)}>
            <div
                class="bg-white dark:bg-gray-800 p-6 rounded-lg max-w-md w-full"
            >
                <h3
                    class="text-lg font-bold mb-4 text-gray-900 dark:text-white"
                >
                    Add Person
                </h3>
                <div class="space-y-4">
                    <input
                        type="text"
                        placeholder="Full Name"
                        bind:value={newPerson.full_name}
                        class="block w-full rounded-md border-gray-300 dark:bg-gray-700 px-3 py-2 border mb-2"
                    />
                    <input
                        type="email"
                        placeholder="Email"
                        bind:value={newPerson.email}
                        class="block w-full rounded-md border-gray-300 dark:bg-gray-700 px-3 py-2 border mb-2"
                    />
                    <div class="grid grid-cols-2 gap-4">
                        <input
                            type="number"
                            placeholder="Age"
                            bind:value={newPerson.age}
                            class="block w-full rounded-md border-gray-300 dark:bg-gray-700 px-3 py-2 border"
                        />
                        <input
                            type="text"
                            placeholder="Gender"
                            bind:value={newPerson.gender}
                            class="block w-full rounded-md border-gray-300 dark:bg-gray-700 px-3 py-2 border"
                        />
                    </div>
                    <input
                        type="text"
                        placeholder="Ethnicity"
                        bind:value={newPerson.ethnicity}
                        class="block w-full rounded-md border-gray-300 dark:bg-gray-700 px-3 py-2 border mb-2"
                    />
                    <input
                        type="text"
                        placeholder="Location"
                        bind:value={newPerson.location}
                        class="block w-full rounded-md border-gray-300 dark:bg-gray-700 px-3 py-2 border mb-2"
                    />
                </div>
                <div class="mt-6 flex justify-end space-x-3">
                    <button
                        onclick={() => (showAddPersonModal = false)}
                        class="px-3 py-2 text-gray-700">Cancel</button
                    >
                    <button
                        onclick={createPerson}
                        class="px-4 py-2 bg-green-600 text-white rounded-md"
                        >Save Person</button
                    >
                </div>
            </div>
        </Modal>
    {/if}

    {#if showAddMemberModal}
        <Modal onclose={() => (showAddMemberModal = false)}>
            <div
                class="bg-white dark:bg-gray-800 p-6 rounded-lg max-w-sm w-full"
            >
                <h3
                    class="text-lg font-bold mb-4 text-gray-900 dark:text-white"
                >
                    Add Person to {selectedAudience?.name}
                </h3>
                <select
                    bind:value={selectedPersonIdToAdd}
                    class="block w-full rounded-md border-gray-300 dark:bg-gray-700 px-3 py-2 border mb-4"
                >
                    <option value="">Select a person...</option>
                    {#each people as person}
                        <option value={person.id}
                            >{person.full_name} ({person.email})</option
                        >
                    {/each}
                </select>
                <div class="flex justify-end space-x-3">
                    <button
                        onclick={() => (showAddMemberModal = false)}
                        class="px-3 py-2 text-gray-700">Cancel</button
                    >
                    <button
                        onclick={addPersonToAudience}
                        class="px-4 py-2 bg-indigo-600 text-white rounded-md"
                        >Add</button
                    >
                </div>
            </div>
        </Modal>
    {/if}
</div>
