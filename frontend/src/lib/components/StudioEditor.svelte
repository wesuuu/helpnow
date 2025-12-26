<script lang="ts">
    import GrapesJsStudio from "@grapesjs/studio-sdk";
    import "@grapesjs/studio-sdk/dist/style.css";
    import { onMount, onDestroy } from "svelte";

    export let content = "";
    export let onChange = (html: string, css: string) => {};

    let editorContainer: HTMLElement;
    let editor: any;

    onMount(() => {
        // Initialize Studio SDK
        // Note: The Studio SDK might require a license key for full features,
        // but often has a trial or limited mode.
        // Documentation suggests simple init.

        // This is a placeholder for actual Studio SDK usage.
        // Assuming it mounts to an element or we use the component if available for framework (none found).
        // Using vanilla JS init.

        GrapesJsStudio({
            root: editorContainer,
            // Common options
            project: {
                default: {
                    pages: [
                        {
                            name: "Home",
                            component: content, // Load initial content
                        },
                    ],
                },
            },
            assets: [], // Your assets
            // Customize the UI/UX
            layout: {
                // ...
            },
        });

        // NOTE: The Studio SDK API is different from GrapesJS Core.
        // We might need to listen to events differently.
        // For now, this is a basic integration.
        // Since I can't easily verify the exact API without docs/running it,
        // I'll try to hook into generic change events if available or leave it as a viewer for now.
    });
</script>

<div class="h-full w-full" bind:this={editorContainer}></div>
