<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import TemplateEditor from "../lib/components/TemplateEditor.svelte";
    // import StudioEditor from "../lib/components/StudioEditor.svelte";
    import { toast } from "../lib/stores/toast";
    import { router } from "../lib/router.svelte.js";

    type Template = {
        id: number;
        organization_id: number;
        name: string;
        type: string;
        content: string; // HTML content
        schema: string; // JSON schema string
        created_at: string;
    };

    const DEFAULTS: Record<string, string> = {
        AD: `<div style="text-align: center; padding: 20px; border: 1px solid #ccc; max-width: 300px; margin: auto;">
    <h2 style="color: #333;">Special Offer!</h2>
    <p style="color: #666;">Get 50% off on your first purchase.</p>
    <a href="#" style="background-color: #007bff; color: white; padding: 10px 20px; text-decoration: none; display: inline-block; border-radius: 5px;">Shop Now</a>
</div>`,
        FORM: `<form style="max-width: 400px; margin: auto; padding: 20px; border: 1px solid #eee; border-radius: 8px;">
    <div style="margin-bottom: 15px;">
        <label style="display: block; margin-bottom: 5px; font-weight: bold;">Name</label>
        <input type="text" placeholder="John Doe" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px;" />
    </div>
    <div style="margin-bottom: 15px;">
        <label style="display: block; margin-bottom: 5px; font-weight: bold;">Email</label>
        <input type="email" placeholder="john@example.com" style="width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px;" />
    </div>
    <button type="submit" style="width: 100%; background-color: #28a745; color: white; padding: 10px; border: none; border-radius: 4px; cursor: pointer;">Subscribe</button>
</form>`,
        LANDING_PAGE: `<div style="font-family: Arial, sans-serif;">
    <header style="background-color: #f8f9fa; padding: 40px 20px; text-align: center;">
        <h1 style="color: #333; font-size: 2.5em;">Welcome to Our Service</h1>
        <p style="color: #666; font-size: 1.2em;">We provide the best solutions for your business.</p>
        <a href="#contact" style="background-color: #007bff; color: white; padding: 12px 24px; text-decoration: none; border-radius: 5px; display: inline-block; margin-top: 20px;">Get Started</a>
    </header>
    <section style="padding: 40px 20px; max-width: 800px; margin: auto; display: flex; justify-content: space-around;">
        <div style="text-align: center; max-width: 30%;">
            <h3>Feature 1</h3>
            <p>Description of feature 1.</p>
        </div>
        <div style="text-align: center; max-width: 30%;">
            <h3>Feature 2</h3>
            <p>Description of feature 2.</p>
        </div>
        <div style="text-align: center; max-width: 30%;">
            <h3>Feature 3</h3>
            <p>Description of feature 3.</p>
        </div>
    </section>
</div>`,
        EMAIL: `<table width="100%" border="0" cellspacing="0" cellpadding="0" style="font-family: sans-serif;">
    <tr>
        <td align="center" style="padding: 20px; background-color: #f6f6f6;">
            <table width="600" border="0" cellspacing="0" cellpadding="0" style="background-color: #ffffff; border-radius: 4px; overflow: hidden; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
                <tr>
                    <td style="padding: 30px; text-align: center;">
                        <h2 style="margin: 0; color: #333333;">Welcome Aboard!</h2>
                        <p style="color: #666666; line-height: 1.5;">Thank you for joining us. We are excited to have you on board.</p>
                        <a href="#" style="display: inline-block; margin-top: 20px; background-color: #007bff; color: #ffffff; padding: 10px 20px; text-decoration: none; border-radius: 4px;">Verify Email</a>
                    </td>
                </tr>
                <tr>
                    <td style="background-color: #eeeeee; padding: 10px; text-align: center; font-size: 12px; color: #999999;">
                        &copy; 2024 HelpNow. All rights reserved.
                    </td>
                </tr>
            </table>
        </td>
    </tr>
</table>`,
    };

    let templates: Template[] = [];
    let loading = true;
    let editingTemplate: Template | null = null;
    let isCreating = false;
    let editorMode: "visual" | "studio" | "code" = "visual";
    let showPreview = false;
    let isDirty = false;
    let visualEditor: any;

    // Form bind variables
    let formName = "";
    let formType = "AD";

    async function fetchTemplates() {
        loading = true;
        try {
            // Hardcoding org_id=1 for now as per previous context/pattern or we should get from auth
            const res = await fetch("/api/templates?organization_id=1");
            if (res.ok) {
                templates = await res.json();
            } else {
                toast.error("Failed to load templates");
            }
        } catch (e) {
            toast.error("Error loading templates");
        } finally {
            loading = false;
        }
    }

    async function saveTemplate() {
        if (!editingTemplate) return;

        try {
            const url = isCreating
                ? "/api/templates"
                : `/api/templates/${editingTemplate.id}`;
            const method = isCreating ? "POST" : "PUT";

            // For creating, we need org_id. For updating, we generally act on the resource
            const payload = {
                organization_id: 1, // hardcoded for now
                name: editingTemplate.name,
                type: editingTemplate.type,
                content: editingTemplate.content,
                schema: editingTemplate.schema || "{}",
            };

            const res = await fetch(url, {
                method,
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(payload),
            });

            if (res.ok) {
                toast.success("Template saved successfully");
                editingTemplate = null;
                isCreating = false;
                setDirty(false);
                fetchTemplates();
            } else {
                toast.error("Failed to save template");
            }
        } catch (e) {
            toast.error("Error saving template");
        }
    }

    function startCreate() {
        editingTemplate = {
            id: 0,
            organization_id: 1,
            name: "New Template",
            type: "AD",
            content: DEFAULTS["AD"],
            schema: "{}",
            created_at: "",
        };
        isCreating = true;
        editorMode = "visual";
        setDirty(false);
    }

    function startEdit(t: Template) {
        if (isDirty) {
            if (
                !confirm(
                    "You have unsaved changes. Are you sure you want to discard them?",
                )
            ) {
                return;
            }
        }
        // Clone to avoid mutating list directly
        editingTemplate = { ...t };
        isCreating = false;
        editorMode = "visual";
        setDirty(false);
    }

    function cancelEdit() {
        if (isDirty) {
            if (
                !confirm(
                    "You have unsaved changes. Are you sure you want to discard them?",
                )
            ) {
                return;
            }
        }
        editingTemplate = null;
        isCreating = false;
        setDirty(false);
    }

    function setDirty(dirty: boolean) {
        isDirty = dirty;
        if (isDirty) {
            router.registerGuard(() => {
                return confirm(
                    "You have unsaved changes. Are you sure you want to leave?",
                );
            });
        } else {
            router.unregisterGuard();
        }
    }

    function handleTypeChange() {
        if (!editingTemplate) return;
        if (DEFAULTS[editingTemplate.type]) {
            const newContent = DEFAULTS[editingTemplate.type];
            editingTemplate.content = newContent;

            // Update editor explicitly if active
            if (editorMode === "visual" && visualEditor) {
                // @ts-ignore
                visualEditor.setContent(newContent);
            }
            // Mark as dirty because content changed
            setDirty(true);
        }
    }

    function handleEditorChange(html: string, styles: string) {
        if (editingTemplate && editorMode === "visual") {
            try {
                // Combine HTML and CSS for storage and preview
                // Obfuscate style tag to prevent Svelte preprocessor confusion
                const styleOpen = "<" + "style" + ">";
                const styleClose = "</" + "style" + ">";
                const cssBlock = styles
                    ? `${styleOpen}${styles}${styleClose}`
                    : "";
                editingTemplate.content = html + cssBlock;
                setDirty(true);
            } catch (e) {
                console.error("Error updating editor content:", e);
            }
        }
    }

    onDestroy(() => {
        router.unregisterGuard();
    });

    onMount(() => {
        fetchTemplates();
    });
</script>

<svelte:window
    on:beforeunload={(e) => {
        if (isDirty && editingTemplate) {
            e.preventDefault();
            e.returnValue = "";
        }
    }}
/>

<div class="space-y-6 h-full flex flex-col">
    {#if !editingTemplate}
        <div class="flex justify-between items-center">
            <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">
                Content Templates
            </h1>
            <button
                class="px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700"
                onclick={startCreate}>Create Template</button
            >
        </div>

        {#if loading}
            <div class="text-center py-4">Loading...</div>
        {:else if templates.length === 0}
            <div
                class="text-center py-10 bg-gray-50 dark:bg-gray-800 rounded-lg"
            >
                <p class="text-gray-500">
                    No templates found. Create one to get started.
                </p>
            </div>
        {:else}
            <div
                class="bg-white dark:bg-gray-800 shadow rounded-lg overflow-hidden"
            >
                <table
                    class="min-w-full divide-y divide-gray-200 dark:divide-gray-700"
                >
                    <thead class="bg-gray-50 dark:bg-gray-700">
                        <tr>
                            <th
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
                                >Name</th
                            >
                            <th
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
                                >Type</th
                            >
                            <th
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
                                >Created</th
                            >
                            <th
                                class="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
                                >Actions</th
                            >
                        </tr>
                    </thead>
                    <tbody
                        class="divide-y divide-gray-200 dark:divide-gray-700"
                    >
                        {#each templates as t}
                            <tr
                                class="hover:bg-gray-50 dark:hover:bg-gray-700 cursor-pointer"
                                onclick={() => startEdit(t)}
                            >
                                <td
                                    class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white"
                                    >{t.name}</td
                                >
                                <td
                                    class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400"
                                    >{t.type}</td
                                >
                                <td
                                    class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400"
                                    >{new Date(
                                        t.created_at,
                                    ).toLocaleDateString()}</td
                                >
                                <td
                                    class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium"
                                >
                                    <button
                                        class="text-indigo-600 hover:text-indigo-900 dark:text-indigo-400 dark:hover:text-indigo-300"
                                        >Edit</button
                                    >
                                </td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
        {/if}
    {:else}
        <!-- Editor Mode -->
        <div class="flex flex-col h-full relative">
            <div
                class="flex justify-between items-center mb-4 bg-white dark:bg-gray-800 p-2 rounded shadow-sm"
            >
                <div class="flex items-center space-x-4">
                    <input
                        type="text"
                        bind:value={editingTemplate.name}
                        class="px-3 py-2 border rounded text-lg font-medium text-gray-900 dark:bg-gray-700 dark:text-white dark:border-gray-600"
                    />
                    <select
                        bind:value={editingTemplate.type}
                        onchange={handleTypeChange}
                        class="px-3 py-2 border rounded text-sm text-gray-900 dark:bg-gray-700 dark:text-white dark:border-gray-600"
                    >
                        <option value="AD">Ad</option>
                        <option value="FORM">Form</option>
                        <option value="LANDING_PAGE">Landing Page</option>
                        <option value="EMAIL">Email</option>
                    </select>
                </div>

                <!-- Toolbar -->
                <div class="flex items-center space-x-2">
                    <div
                        class="bg-gray-100 dark:bg-gray-700 rounded p-1 flex mr-4"
                    >
                        <button
                            class="px-3 py-1 rounded text-sm {editorMode ===
                            'visual'
                                ? 'bg-white shadow dark:bg-gray-600 dark:text-white'
                                : 'text-gray-500 hover:text-gray-700 dark:text-gray-400'}"
                            onclick={() => (editorMode = "visual")}
                            >Visual</button
                        >
                        <button
                            class="px-3 py-1 rounded text-sm {editorMode ===
                            'code'
                                ? 'bg-white shadow dark:bg-gray-600 dark:text-white'
                                : 'text-gray-500 hover:text-gray-700 dark:text-gray-400'}"
                            onclick={() => (editorMode = "code")}>Code</button
                        >
                    </div>

                    <button
                        class="px-4 py-2 border border-gray-300 text-gray-700 rounded hover:bg-gray-50 dark:text-gray-300 dark:border-gray-600 dark:hover:bg-gray-700"
                        onclick={() => (showPreview = true)}>Preview</button
                    >

                    <button
                        class="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300"
                        onclick={cancelEdit}>Cancel</button
                    >
                    <button
                        class="px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700"
                        onclick={saveTemplate}>Save</button
                    >
                </div>
            </div>

            <!-- Content Area -->
            <div
                class="flex-1 border rounded relative overflow-hidden"
                style="min-height: 600px;"
            >
                {#if editorMode === "visual"}
                    <TemplateEditor
                        bind:this={visualEditor}
                        content={editingTemplate.content}
                        onChange={handleEditorChange}
                    />
                    <div
                        class="md:hidden p-4 bg-yellow-50 text-yellow-800 rounded absolute bottom-0 w-full text-center opacity-75"
                    >
                        Editor is optimized for desktop view.
                    </div>
                {:else}
                    <textarea
                        class="w-full h-full p-4 font-mono text-sm bg-gray-50 dark:bg-gray-900 dark:text-gray-100 resize-none focus:outline-none"
                        bind:value={editingTemplate.content}
                    ></textarea>
                {/if}
            </div>

            <!-- Preview Modal -->
            {#if showPreview}
                <div
                    class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 p-4"
                >
                    <div
                        class="bg-white dark:bg-gray-800 rounded-lg shadow-xl w-full max-w-4xl h-[80vh] flex flex-col"
                    >
                        <div
                            class="flex justify-between items-center p-4 border-b dark:border-gray-700"
                        >
                            <h3
                                class="text-lg font-medium text-gray-900 dark:text-white"
                            >
                                Preview
                            </h3>
                            <button
                                class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
                                onclick={() => (showPreview = false)}
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
                                        d="M6 18L18 6M6 6l12 12"
                                    />
                                </svg>
                            </button>
                        </div>
                        <div class="flex-1 bg-gray-100 p-4 overflow-auto">
                            <div
                                class="bg-white shadow mx-auto"
                                style="min-height: 100%;"
                            >
                                <iframe
                                    srcdoc={editingTemplate.content}
                                    title="Preview"
                                    class="w-full h-full border-none"
                                    style="min-height: 500px;"
                                ></iframe>
                            </div>
                        </div>
                    </div>
                </div>
            {/if}
        </div>
    {/if}
</div>
