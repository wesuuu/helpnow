<script lang="ts">
    import { onMount } from "svelte";
    import Modal from "../lib/components/Modal.svelte";

    interface EmailDomain {
        id: number;
        domain: string;
        is_verified: boolean;
        dkim_record: string;
        spf_record: string;
        created_at: string;
    }

    interface EmailTemplate {
        id: number;
        name: string;
        subject: string;
        body: string;
        updated_at: string;
    }

    let activeTab = $state("domains"); // 'domains' | 'templates'

    // -- Domains State --
    let emailDomains = $state<EmailDomain[]>([]);
    let showCreateDomainModal = $state(false);
    let newDomainName = $state("");
    let viewingDomain = $state<EmailDomain | null>(null);
    let showDomainDetailsModal = $state(false);

    // -- Templates State --
    let emailTemplates = $state<EmailTemplate[]>([]);
    let showTemplateModal = $state(false);
    let editingTemplateId = $state<number | null>(null);
    let templateName = $state("");
    let templateSubject = $state("");
    let templateBody = $state("");

    // Template Preview State
    let testFirstName = $state("John");
    let testLastName = $state("Doe");

    // Derived preview (simple reactive statements)
    let previewSubject = $derived(
        templateSubject
            .replace(/{{first_name}}/g, testFirstName)
            .replace(/{{last_name}}/g, testLastName),
    );
    let previewBody = $derived(
        templateBody
            .replace(/{{first_name}}/g, testFirstName)
            .replace(/{{last_name}}/g, testLastName),
    );

    async function fetchEmailDomains() {
        try {
            const res = await fetch("/api/email-domains?organization_id=1");
            if (res.ok) {
                emailDomains = await res.json();
            }
        } catch (error) {
            console.error("Failed to fetch email domains", error);
        }
    }

    async function createDomain() {
        if (!newDomainName) return;
        try {
            const res = await fetch("/api/email-domains", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    organization_id: 1, // Hardcoded
                    domain: newDomainName,
                }),
            });
            if (res.ok) {
                showCreateDomainModal = false;
                newDomainName = "";
                fetchEmailDomains();
            } else {
                alert("Failed to create domain");
            }
        } catch (error) {
            console.error("Failed to create domain", error);
        }
    }

    async function verifyDomain(id: number) {
        try {
            const res = await fetch(`/api/email-domains/${id}/verify`, {
                method: "POST",
            });
            if (res.ok) {
                const data = await res.json();
                if (data.is_verified) {
                    fetchEmailDomains(); // Refresh list
                    if (viewingDomain && viewingDomain.id === id) {
                        viewingDomain.is_verified = true;
                    }
                }
            }
        } catch (error) {
            console.error("Failed to verify domain", error);
        }
    }

    async function deleteDomain(id: number) {
        if (!confirm("Are you sure?")) return;
        try {
            const res = await fetch(`/api/email-domains/${id}`, {
                method: "DELETE",
            });
            if (res.ok) {
                fetchEmailDomains();
                if (viewingDomain?.id === id) {
                    showDomainDetailsModal = false;
                    viewingDomain = null;
                }
            }
        } catch (error) {
            console.error("Failed to delete domain", error);
        }
    }

    function viewDomainDetails(domain: EmailDomain) {
        viewingDomain = domain;
        showDomainDetailsModal = true;
    }

    // -- Template Functions --

    async function fetchEmailTemplates() {
        try {
            const res = await fetch("/api/email-templates?organization_id=1");
            if (res.ok) {
                emailTemplates = await res.json();
            }
        } catch (error) {
            console.error("Failed to fetch templates", error);
        }
    }

    function openCreateTemplateModal() {
        editingTemplateId = null;
        templateName = "";
        templateSubject = "";
        templateBody = "";
        showTemplateModal = true;
    }

    function openEditTemplateModal(tmpl: EmailTemplate) {
        editingTemplateId = tmpl.id;
        templateName = tmpl.name;
        templateSubject = tmpl.subject;
        templateBody = tmpl.body;
        showTemplateModal = true;
    }

    async function saveTemplate() {
        if (!templateName || !templateSubject || !templateBody) {
            alert("Please fill in all fields.");
            return;
        }

        const payload = {
            organization_id: 1,
            name: templateName,
            subject: templateSubject,
            body: templateBody,
        };

        try {
            let res;
            if (editingTemplateId) {
                res = await fetch(`/api/email-templates/${editingTemplateId}`, {
                    method: "PUT",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(payload),
                });
            } else {
                res = await fetch("/api/email-templates", {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(payload),
                });
            }

            if (res.ok) {
                showTemplateModal = false;
                fetchEmailTemplates();
            } else {
                alert("Failed to save template");
            }
        } catch (error) {
            console.error("Failed to save template", error);
        }
    }

    async function deleteTemplate(id: number) {
        if (!confirm("Are you sure you want to delete this template?")) return;
        try {
            const res = await fetch(`/api/email-templates/${id}`, {
                method: "DELETE",
            });
            if (res.ok) {
                fetchEmailTemplates();
            }
        } catch (error) {
            console.error("Failed to delete template", error);
        }
    }

    onMount(() => {
        fetchEmailDomains();
        fetchEmailTemplates();
    });
</script>

<div class="space-y-6">
    <div
        class="flex justify-between items-center border-b border-gray-200 dark:border-gray-700 pb-4"
    >
        <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">
            Email Settings
        </h1>
    </div>

    <!-- Tabs -->
    <div class="border-b border-gray-200 dark:border-gray-700">
        <nav class="-mb-px flex space-x-8" aria-label="Tabs">
            <button
                onclick={() => (activeTab = "domains")}
                class="{activeTab === 'domains'
                    ? 'border-indigo-500 text-indigo-600 dark:text-indigo-400'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'} whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm"
            >
                Sender Domains
            </button>
            <button
                onclick={() => (activeTab = "templates")}
                class="{activeTab === 'templates'
                    ? 'border-indigo-500 text-indigo-600 dark:text-indigo-400'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'} whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm"
            >
                Email Templates
            </button>
        </nav>
    </div>

    {#if activeTab === "domains"}
        <div class="mt-4">
            <div class="flex justify-end mb-4">
                <button
                    onclick={() => (showCreateDomainModal = true)}
                    class="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700"
                >
                    Add Domain
                </button>
            </div>
            <div class="overflow-x-auto">
                <table class="min-w-full divide-y divide-gray-200">
                    <thead class="bg-gray-50">
                        <tr>
                            <th
                                scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                                >Domain</th
                            >
                            <th
                                scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                                >Status</th
                            >
                            <th
                                scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                                >Created</th
                            >
                            <th scope="col" class="relative px-6 py-3"
                                ><span class="sr-only">Actions</span></th
                            >
                        </tr>
                    </thead>
                    <tbody class="bg-white divide-y divide-gray-200">
                        {#each emailDomains as domain}
                            <tr>
                                <td
                                    class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900"
                                    >{domain.domain}</td
                                >
                                <td
                                    class="px-6 py-4 whitespace-nowrap text-sm text-gray-500"
                                >
                                    {#if domain.is_verified}
                                        <span
                                            class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800"
                                            >Verified</span
                                        >
                                    {:else}
                                        <span
                                            class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-yellow-100 text-yellow-800"
                                            >Unverified</span
                                        >
                                    {/if}
                                </td>
                                <td
                                    class="px-6 py-4 whitespace-nowrap text-sm text-gray-500"
                                    >{new Date(
                                        domain.created_at,
                                    ).toLocaleDateString()}</td
                                >
                                <td
                                    class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium"
                                >
                                    <button
                                        onclick={() =>
                                            viewDomainDetails(domain)}
                                        class="text-indigo-600 hover:text-indigo-900 mr-4"
                                        >Verify/View</button
                                    >
                                    <button
                                        onclick={() => deleteDomain(domain.id)}
                                        class="text-red-600 hover:text-red-900"
                                        >Delete</button
                                    >
                                </td>
                            </tr>
                        {/each}
                        {#if emailDomains.length === 0}
                            <tr>
                                <td
                                    colspan="4"
                                    class="px-6 py-4 text-center text-sm text-gray-500"
                                >
                                    No email domains configured. Add one to
                                    start sending emails.
                                </td>
                            </tr>
                        {/if}
                    </tbody>
                </table>
            </div>
        </div>
    {:else if activeTab === "templates"}
        <div class="mt-4">
            <div class="flex justify-end mb-4">
                <button
                    onclick={openCreateTemplateModal}
                    class="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700"
                >
                    Create Template
                </button>
            </div>
            <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
                {#each emailTemplates as tmpl}
                    <div
                        class="bg-white overflow-hidden shadow rounded-lg divide-y divide-gray-200 flex flex-col"
                    >
                        <div class="px-4 py-5 sm:px-6 flex-1">
                            <h3
                                class="text-lg leading-6 font-medium text-gray-900"
                            >
                                {tmpl.name}
                            </h3>
                            <p
                                class="mt-1 max-w-2xl text-sm text-gray-500 truncate"
                            >
                                Subject: {tmpl.subject}
                            </p>
                        </div>
                        <div
                            class="px-4 py-4 sm:px-6 flex justify-end space-x-3 bg-gray-50"
                        >
                            <button
                                onclick={() => openEditTemplateModal(tmpl)}
                                class="text-indigo-600 hover:text-indigo-900 text-sm font-medium"
                                >Edit</button
                            >
                            <button
                                onclick={() => deleteTemplate(tmpl.id)}
                                class="text-red-600 hover:text-red-900 text-sm font-medium"
                                >Delete</button
                            >
                        </div>
                    </div>
                {/each}
                {#if emailTemplates.length === 0}
                    <div class="col-span-full text-center py-12 text-gray-500">
                        No templates found. Create one to use in your campaigns.
                    </div>
                {/if}
            </div>
        </div>
    {/if}
</div>

<!-- Create/Edit Domain Modal -->
{#if showCreateDomainModal}
    <Modal onclose={() => (showCreateDomainModal = false)}>
        <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <h3 class="text-lg leading-6 font-medium text-gray-900">
                Add Sender Domain
            </h3>
            <div class="mt-2">
                <p class="text-sm text-gray-500">
                    Enter the domain you want to send emails from.
                </p>
                <input
                    type="text"
                    bind:value={newDomainName}
                    placeholder="example.com"
                    class="mt-2 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                />
            </div>
        </div>
        <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
            <button
                onclick={createDomain}
                type="button"
                class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 sm:ml-3 sm:w-auto sm:text-sm"
                >Add</button
            >
            <button
                onclick={() => (showCreateDomainModal = false)}
                type="button"
                class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
                >Cancel</button
            >
        </div>
    </Modal>
{/if}

<!-- Domain Details Modal -->
{#if showDomainDetailsModal && viewingDomain}
    <Modal onclose={() => (showDomainDetailsModal = false)}>
        <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <h3 class="text-lg leading-6 font-medium text-gray-900">
                Domain Verification
            </h3>
            <div class="mt-4 space-y-4">
                <p class="text-sm text-gray-500">
                    Add these records to your DNS settings to verify ownership.
                </p>
                <div>
                    <label
                        for="dkim-record"
                        class="block text-sm font-medium text-gray-700"
                        >DKIM Record (TXT)</label
                    >
                    <div class="mt-1 flex rounded-md shadow-sm">
                        <input
                            id="dkim-record"
                            type="text"
                            readonly
                            value={viewingDomain?.dkim_record}
                            class="flex-1 block w-full rounded-none rounded-l-md border-gray-300 sm:text-sm bg-gray-50"
                        />
                        <button
                            onclick={() =>
                                viewingDomain &&
                                navigator.clipboard.writeText(
                                    viewingDomain.dkim_record,
                                )}
                            class="-ml-px relative inline-flex items-center space-x-2 px-4 py-2 border border-gray-300 text-sm font-medium rounded-r-md text-gray-700 bg-gray-50 hover:bg-gray-100 focus:outline-none focus:ring-1 focus:ring-indigo-500 focus:border-indigo-500"
                            >Copy</button
                        >
                    </div>
                </div>
                <div>
                    <label
                        for="spf-record"
                        class="block text-sm font-medium text-gray-700"
                        >SPF Record (TXT)</label
                    >
                    <div class="mt-1 flex rounded-md shadow-sm">
                        <input
                            id="spf-record"
                            type="text"
                            readonly
                            value={viewingDomain?.spf_record}
                            class="flex-1 block w-full rounded-none rounded-l-md border-gray-300 sm:text-sm bg-gray-50"
                        />
                        <button
                            onclick={() =>
                                viewingDomain &&
                                navigator.clipboard.writeText(
                                    viewingDomain.spf_record,
                                )}
                            class="-ml-px relative inline-flex items-center space-x-2 px-4 py-2 border border-gray-300 text-sm font-medium rounded-r-md text-gray-700 bg-gray-50 hover:bg-gray-100 focus:outline-none focus:ring-1 focus:ring-indigo-500 focus:border-indigo-500"
                            >Copy</button
                        >
                    </div>
                </div>
                {#if viewingDomain?.is_verified}
                    <div class="rounded-md bg-green-50 p-4">
                        <div class="flex">
                            <div class="flex-shrink-0">
                                <svg
                                    class="h-5 w-5 text-green-400"
                                    xmlns="http://www.w3.org/2000/svg"
                                    viewBox="0 0 20 20"
                                    fill="currentColor"
                                    aria-hidden="true"
                                >
                                    <path
                                        fill-rule="evenodd"
                                        d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                                        clip-rule="evenodd"
                                    />
                                </svg>
                            </div>
                            <div class="ml-3">
                                <p class="text-sm font-medium text-green-800">
                                    Domain Verified
                                </p>
                            </div>
                        </div>
                    </div>
                {:else}
                    <button
                        onclick={() =>
                            viewingDomain && verifyDomain(viewingDomain.id)}
                        type="button"
                        class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:text-sm"
                        >Verify DNS Records</button
                    >
                {/if}
            </div>
        </div>
        <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
            <button
                onclick={() => (showDomainDetailsModal = false)}
                type="button"
                class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
                >Close</button
            >
        </div>
    </Modal>
{/if}

<!-- Template Modal -->
{#if showTemplateModal}
    <Modal onclose={() => (showTemplateModal = false)} size="xl">
        <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <h3
                class="text-lg leading-6 font-medium text-gray-900"
                id="modal-title"
            >
                {editingTemplateId ? "Edit" : "Create"} Email Template
            </h3>

            <div class="mt-4 grid grid-cols-1 md:grid-cols-2 gap-6">
                <!-- Editor Column -->
                <div class="space-y-4">
                    <div>
                        <label
                            for="template-name"
                            class="block text-sm font-medium text-gray-700"
                            >Template Name</label
                        >
                        <input
                            type="text"
                            id="template-name"
                            bind:value={templateName}
                            class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                            placeholder="e.g., Welcome Email"
                        />
                    </div>
                    <div>
                        <label
                            for="template-subject"
                            class="block text-sm font-medium text-gray-700"
                            >Subject</label
                        >
                        <input
                            type="text"
                            id="template-subject"
                            bind:value={templateSubject}
                            class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                            placeholder="Subject line with {`{{first_name}}`}"
                        />
                        <p class="mt-1 text-xs text-gray-500">
                            Supports variables: <code>{`{{first_name}}`}</code>,
                            <code>{`{{last_name}}`}</code>
                        </p>
                    </div>
                    <div>
                        <label
                            for="template-body"
                            class="block text-sm font-medium text-gray-700"
                            >Body</label
                        >
                        <textarea
                            id="template-body"
                            bind:value={templateBody}
                            rows="10"
                            class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm font-mono"
                        ></textarea>
                        <p class="mt-1 text-xs text-gray-500">
                            Supports variables: <code>{`{{first_name}}`}</code>,
                            <code>{`{{last_name}}`}</code>
                        </p>
                    </div>
                </div>

                <!-- Preview Column -->
                <div class="bg-gray-50 p-4 rounded-lg flex flex-col h-full">
                    <h4 class="text-sm font-medium text-gray-700 mb-2">
                        Preview
                    </h4>

                    <div class="grid grid-cols-2 gap-2 mb-4">
                        <div>
                            <label
                                for="test-first-name"
                                class="block text-xs font-medium text-gray-500"
                                >Test First Name</label
                            >
                            <input
                                type="text"
                                id="test-first-name"
                                bind:value={testFirstName}
                                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-xs"
                            />
                        </div>
                        <div>
                            <label
                                for="test-last-name"
                                class="block text-xs font-medium text-gray-500"
                                >Test Last Name</label
                            >
                            <input
                                type="text"
                                id="test-last-name"
                                bind:value={testLastName}
                                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-xs"
                            />
                        </div>
                    </div>

                    <div
                        class="bg-white border border-gray-200 rounded-md p-4 shadow-sm flex-1 overflow-y-auto hidden-scrollbar"
                    >
                        <div class="border-b border-gray-100 pb-2 mb-2">
                            <span class="text-xs text-gray-500 block"
                                >Subject:</span
                            >
                            <p
                                class="text-sm font-medium text-gray-900 break-words"
                            >
                                {previewSubject || "(No Subject)"}
                            </p>
                        </div>
                        <div>
                            <span class="text-xs text-gray-500 block mb-1"
                                >Body:</span
                            >
                            <div
                                class="text-sm text-gray-800 whitespace-pre-wrap font-sans"
                            >
                                {#if previewBody}
                                    {@html previewBody}
                                {:else}
                                    <span class="text-gray-400 italic"
                                        >(Empty body)</span
                                    >
                                {/if}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
            <button
                onclick={saveTemplate}
                type="button"
                class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 sm:ml-3 sm:w-auto sm:text-sm"
            >
                Save
            </button>
            <button
                onclick={() => (showTemplateModal = false)}
                type="button"
                class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
            >
                Cancel
            </button>
        </div>
    </Modal>
{/if}
