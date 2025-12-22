<script lang="ts">
    import { router } from "./router.svelte.js";
    import type { Component as ComponentType } from "svelte";
    import type { Snippet } from "svelte";

    interface Props {
        path: string;
        component?: ComponentType;
        children?: Snippet;
    }

    let { path, component: Component, children }: Props = $props();

    function matchPath(routePath: string, currentPath: string): boolean {
        if (!routePath) return false;

        // Exact match
        if (routePath === currentPath) return true;

        // Parameter match (simple)
        // e.g. /sites/:id matches /sites/123
        const paramRegex = new RegExp(
            "^" + routePath.replace(/:[^\s/]+/g, "([\\w-]+)") + "$",
        );
        return paramRegex.test(currentPath);
    }

    let isMatch = $derived(matchPath(path, router.path));
</script>

{#if isMatch}
    {#if Component}
        <Component />
    {/if}
    {@render children?.()}
{/if}
