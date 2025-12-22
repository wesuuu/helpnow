<script lang="ts">
    import { onMount } from "svelte";
    import { auth } from "../stores/auth.js";
    import Modal from "../lib/components/Modal.svelte";

    interface Campaign {
        id: number;
        name: string;
        status: string;
        prompt: string;
        content?: string;
    }

    interface SignupCampaign {
        id: number;
        name: string;
        token: string;
        target_audience_id: number;
        created_at: string;
    }

    interface Audience {
        id: number;
        name: string;
        segments?: { id: number; name: string }[];
    }

    let activeTab = $state("email"); // 'email' | 'signup'

    // Email Campaigns State
    let campaigns = $state<Campaign[]>([]);
    let audiences = $state<Audience[]>([]);
    let showCreateModal = $state(false);
    let newCampaignName = $state("");
    let newCampaignPrompt = $state("");
    let selectedAudienceId = $state("");
    let expandedCampaignId = $state<number | null>(null);
    let isGenerating = $state(false);

    // Signup Campaigns State
    let signupCampaigns = $state<SignupCampaign[]>([]);
    let showCreateSignupModal = $state(false);
    let newSignupName = $state("");
    let selectedSignupAudienceId = $state("");

    async function fetchCampaigns() {
        // TODO: Use real org ID
        const res = await fetch(`/api/campaigns?organization_id=1`);
        if (res.ok) {
            campaigns = await res.json();
        }
    }

    async function fetchSignupCampaigns() {
        const res = await fetch(`/api/signup-campaigns?organization_id=1`);
        if (res.ok) {
            signupCampaigns = await res.json();
        }
    }

    async function fetchAudiences() {
        const res = await fetch(`/api/audiences?organization_id=1`);
        if (res.ok) {
            audiences = await res.json();
        }
    }

    async function createCampaign() {
        const audience = audiences.find(
            (a) => a.id == parseInt(selectedAudienceId),
        );
        let segmentId = 1; // Default fallback
        if (audience && audience.segments && audience.segments.length > 0) {
            segmentId = audience.segments[0].id; // Just pick first segment for MVP
        }

        const res = await fetch("/api/campaigns", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                organization_id: 1, // Hardcoded for MVP
                audience_segment_id: segmentId,
                name: newCampaignName,
                prompt: newCampaignPrompt,
            }),
        });

        if (res.ok) {
            showCreateModal = false;
            fetchCampaigns();
            newCampaignName = "";
            newCampaignPrompt = "";
        }
    }

    async function createSignupCampaign() {
        const res = await fetch("/api/signup-campaigns", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                organization_id: 1,
                target_audience_id: parseInt(selectedSignupAudienceId),
                name: newSignupName,
            }),
        });

        if (res.ok) {
            showCreateSignupModal = false;
            fetchSignupCampaigns();
            newSignupName = "";
            selectedSignupAudienceId = "";
        }
    }

    async function generateContent(campaignId: number) {
        isGenerating = true;
        const res = await fetch(`/api/campaigns/${campaignId}/content`, {
            method: "POST",
        });
        if (res.ok) {
            await fetchCampaigns();
        }
        isGenerating = false;
    }

    async function activateCampaign(campaign: Campaign) {
        const res = await fetch(`/api/campaigns/${campaign.id}`, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                content: campaign.content,
                status: "ACTIVE",
                schedule_interval: "DAILY",
            }),
        });
        if (res.ok) {
            fetchCampaigns();
        }
    }

    function copyIntegrationUrl(token: string) {
        const url = `${window.location.origin}/api/public/capture/${token}`; // In real life, might be different domain
        navigator.clipboard.writeText(url);
        alert("Copied URL: " + url);
    }

    onMount(() => {
        fetchCampaigns();
        fetchSignupCampaigns();
        fetchAudiences();
    });
</script>

<div class="space-y-6">
    <div
        class="flex justify-between items-center border-b border-gray-200 dark:border-gray-700 pb-4"
    >
        <div class="flex space-x-4">
            <button
                class={`px-4 py-2 font-medium rounded-md ${activeTab === "email" ? "bg-indigo-100 text-indigo-700 dark:bg-indigo-900 dark:text-indigo-200" : "text-gray-500 hover:text-gray-700 dark:text-gray-400"}`}
                onclick={() => (activeTab = "email")}
            >
                Email Campaigns
            </button>
            <button
                class={`px-4 py-2 font-medium rounded-md ${activeTab === "signup" ? "bg-indigo-100 text-indigo-700 dark:bg-indigo-900 dark:text-indigo-200" : "text-gray-500 hover:text-gray-700 dark:text-gray-400"}`}
                onclick={() => (activeTab = "signup")}
            >
                Signup Forms
            </button>
        </div>
        <div>
            {#if activeTab === "email"}
                <button
                    onclick={() => (showCreateModal = true)}
                    class="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700"
                >
                    New Email Campaign
                </button>
            {:else}
                <button
                    onclick={() => (showCreateSignupModal = true)}
                    class="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700"
                >
                    New Signup Form
                </button>
            {/if}
        </div>
    </div>

    <!-- EMAIL CAMPAIGNS TAB -->
    {#if activeTab === "email"}
        {#if showCreateModal}
            <Modal onclose={() => (showCreateModal = false)}>
                <div
                    class="bg-white dark:bg-gray-800 p-6 rounded-lg max-w-md w-full"
                >
                    <h3
                        class="text-lg font-bold mb-4 text-gray-900 dark:text-white"
                    >
                        Create New Email Campaign
                    </h3>
                    <div class="space-y-4">
                        <div>
                            <label
                                for="campaign-name"
                                class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                                >Name</label
                            >
                            <input
                                id="campaign-name"
                                type="text"
                                bind:value={newCampaignName}
                                class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white px-3 py-2 border"
                            />
                        </div>
                        <div>
                            <label
                                for="campaign-audience"
                                class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                                >Audience</label
                            >
                            <select
                                id="campaign-audience"
                                bind:value={selectedAudienceId}
                                class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white px-3 py-2 border"
                            >
                                <option value="">Select Audience...</option>
                                {#each audiences as audience}
                                    <option value={audience.id.toString()}
                                        >{audience.name}</option
                                    >
                                {/each}
                            </select>
                        </div>
                        <div>
                            <label
                                for="campaign-prompt"
                                class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                                >AI Prompt</label
                            >
                            <textarea
                                id="campaign-prompt"
                                bind:value={newCampaignPrompt}
                                rows="3"
                                class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white px-3 py-2 border"
                                placeholder="e.g., Write a newsletter about our new feature..."
                            ></textarea>
                        </div>
                    </div>
                    <div class="mt-6 flex justify-end space-x-3">
                        <button
                            onclick={() => (showCreateModal = false)}
                            class="px-3 py-2 text-gray-700 hover:text-gray-900 dark:text-gray-300"
                            >Cancel</button
                        >
                        <button
                            onclick={createCampaign}
                            class="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700"
                            >Create</button
                        >
                    </div>
                </div>
            </Modal>
        {/if}

        <div class="grid gap-6">
            {#each campaigns as campaign}
                <div
                    class="bg-white dark:bg-gray-800 shadow rounded-lg p-6 border border-gray-200 dark:border-gray-700"
                >
                    <div class="flex justify-between items-start">
                        <div>
                            <h3
                                class="text-xl font-semibold text-gray-900 dark:text-white"
                            >
                                {campaign.name}
                            </h3>
                            <p
                                class="text-sm text-gray-500 dark:text-gray-400 mt-1"
                            >
                                Status: <span
                                    class={`font-medium ${campaign.status === "ACTIVE" ? "text-green-600" : "text-yellow-600"}`}
                                    >{campaign.status}</span
                                >
                            </p>
                        </div>
                        <button
                            class="text-indigo-600 hover:text-indigo-800 text-sm font-medium"
                            onclick={() =>
                                (expandedCampaignId =
                                    expandedCampaignId === campaign.id
                                        ? null
                                        : campaign.id)}
                        >
                            {expandedCampaignId === campaign.id
                                ? "Hide Details"
                                : "Manage"}
                        </button>
                    </div>
                    {#if expandedCampaignId === campaign.id}
                        <div
                            class="mt-6 border-t border-gray-100 dark:border-gray-700 pt-4 space-y-4"
                        >
                            <div>
                                <h4
                                    class="text-sm font-medium text-gray-900 dark:text-white"
                                >
                                    Prompt
                                </h4>
                                <p
                                    class="text-gray-600 dark:text-gray-300 text-sm bg-gray-50 dark:bg-gray-900 p-3 rounded mt-1"
                                >
                                    {campaign.prompt}
                                </p>
                            </div>
                            <div>
                                <div
                                    class="flex justify-between items-center mb-2"
                                >
                                    <h4
                                        class="text-sm font-medium text-gray-900 dark:text-white"
                                    >
                                        Email Content
                                    </h4>
                                    <button
                                        onclick={() =>
                                            generateContent(campaign.id)}
                                        disabled={isGenerating}
                                        class="text-xs bg-indigo-50 text-indigo-700 px-2 py-1 rounded hover:bg-indigo-100 disabled:opacity-50"
                                    >
                                        {isGenerating
                                            ? "Generating..."
                                            : "Generate with AI"}
                                    </button>
                                </div>
                                <div
                                    class="bg-gray-50 dark:bg-gray-900 p-4 rounded text-sm text-gray-800 dark:text-gray-200 whitespace-pre-wrap min-h-[100px] border border-gray-200 dark:border-gray-700"
                                >
                                    {campaign.content ||
                                        "No content generated yet."}
                                </div>
                            </div>
                            {#if campaign.status === "DRAFT"}
                                <div class="flex justify-end pt-2">
                                    <button
                                        onclick={() =>
                                            activateCampaign(campaign)}
                                        disabled={!campaign.content}
                                        class="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed"
                                    >
                                        Activate & Schedule Daily
                                    </button>
                                </div>
                            {/if}
                        </div>
                    {/if}
                </div>
            {/each}
            {#if campaigns.length === 0}
                <div class="text-center py-12 text-gray-500 dark:text-gray-400">
                    No email campaigns yet.
                </div>
            {/if}
        </div>
    {/if}

    <!-- SIGNUP CAMPAIGNS TAB -->
    {#if activeTab === "signup"}
        {#if showCreateSignupModal}
            <Modal onclose={() => (showCreateSignupModal = false)}>
                <div
                    class="bg-white dark:bg-gray-800 p-6 rounded-lg max-w-md w-full"
                >
                    <h3
                        class="text-lg font-bold mb-4 text-gray-900 dark:text-white"
                    >
                        Create Lead Capture Form
                    </h3>
                    <div class="space-y-4">
                        <div>
                            <label
                                for="signup-name"
                                class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                                >Name</label
                            >
                            <input
                                id="signup-name"
                                type="text"
                                bind:value={newSignupName}
                                class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white px-3 py-2 border"
                                placeholder="e.g. Winter Sale Landing Page"
                            />
                        </div>
                        <div>
                            <label
                                for="signup-audience"
                                class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                                >Target Audience</label
                            >
                            <select
                                id="signup-audience"
                                bind:value={selectedSignupAudienceId}
                                class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white px-3 py-2 border"
                            >
                                <option value="">Select Audience...</option>
                                {#each audiences as audience}
                                    <option value={audience.id.toString()}
                                        >{audience.name}</option
                                    >
                                {/each}
                            </select>
                            <p class="text-xs text-gray-500 mt-1">
                                Leads from this form will be added to this
                                audience automatically.
                            </p>
                        </div>
                    </div>
                    <div class="mt-6 flex justify-end space-x-3">
                        <button
                            onclick={() => (showCreateSignupModal = false)}
                            class="px-3 py-2 text-gray-700 hover:text-gray-900 dark:text-gray-300"
                            >Cancel</button
                        >
                        <button
                            onclick={createSignupCampaign}
                            class="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700"
                            >Create</button
                        >
                    </div>
                </div>
            </Modal>
        {/if}

        <div class="grid gap-6">
            {#each signupCampaigns as campaign}
                <div
                    class="bg-white dark:bg-gray-800 shadow rounded-lg p-6 border border-gray-200 dark:border-gray-700"
                >
                    <div
                        class="flex flex-col sm:flex-row justify-between items-start sm:items-center"
                    >
                        <div>
                            <h3
                                class="text-xl font-semibold text-gray-900 dark:text-white"
                            >
                                {campaign.name}
                            </h3>
                            <p
                                class="text-sm text-gray-500 dark:text-gray-400 mt-1"
                            >
                                Target Audience ID: {campaign.target_audience_id}
                            </p>
                        </div>
                        <button
                            onclick={() => copyIntegrationUrl(campaign.token)}
                            class="mt-4 sm:mt-0 text-sm bg-gray-100 text-gray-700 px-3 py-1.5 rounded border border-gray-300 hover:bg-gray-200 flex items-center"
                        >
                            Copy Integration URL
                        </button>
                    </div>
                    <div
                        class="mt-4 p-3 bg-gray-50 dark:bg-gray-900 rounded border border-gray-200 dark:border-gray-700 text-xs font-mono text-gray-600 dark:text-gray-300 break-all"
                    >
                        POST {window.location
                            .origin}/api/public/capture/{campaign.token}
                    </div>
                </div>
            {/each}
            {#if signupCampaigns.length === 0}
                <div class="text-center py-12 text-gray-500 dark:text-gray-400">
                    No signup forms yet. Create one to start capturing leads!
                </div>
            {/if}
        </div>
    {/if}
</div>
