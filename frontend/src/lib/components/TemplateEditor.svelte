<script lang="ts">
    import grapesjs from "grapesjs";
    import "grapesjs/dist/css/grapes.min.css";
    import webpagePreset from "grapesjs-preset-webpage";
    import basicBlocks from "grapesjs-blocks-basic";
    import formsPlugin from "grapesjs-plugin-forms";
    import { onMount, onDestroy } from "svelte";

    export let content = "";
    export let onChange = (html: string, css: string) => {};

    let editorContainer: HTMLElement;
    let editor: any;

    onMount(() => {
        editor = grapesjs.init({
            container: editorContainer,
            height: "100%",
            width: "100%",
            storageManager: false, // We'll manage state manually via props/events
            plugins: [webpagePreset, basicBlocks, formsPlugin],
            pluginsOpts: {
                "grapesjs-preset-webpage": {},
                "grapesjs-blocks-basic": {},
                "grapesjs-plugin-forms": {},
            },
        });

        // Set initial content if provided
        if (content) {
            editor.setComponents(content);
        }

        editor.on("update", () => {
            const html = editor.getHtml();
            const css = editor.getCss();
            onChange(html, css);
        });
    });

    onDestroy(() => {
        if (editor) {
            editor.destroy();
        }
    });

    // Expose a method to update content programmatically if needed (e.g. initial load async)
    export function setContent(newContent: string) {
        if (editor) {
            editor.setComponents(newContent);
        }
    }
</script>

<div
    class="h-full w-full border border-gray-300 rounded"
    bind:this={editorContainer}
></div>

<style>
    /* Ensure editor container has height */
    :global(.gjs-cv-canvas) {
        top: 0;
        width: 100%;
        height: 100%;
    }
</style>
