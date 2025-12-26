<script lang="ts">
    let { text, children } = $props<{
        text: string;
        children: import("svelte").Snippet;
    }>();
</script>

<div class="group relative inline-block">
    {@render children()}
    <div
        class="invisible group-hover:visible absolute z-10 w-48 p-2 mt-2 text-sm text-white bg-gray-900 rounded-lg shadow-lg -left-1/2 transform translate-x-1/2 dark:bg-gray-700"
    >
        {text}
        <div
            class="absolute top-0 left-1/2 w-2 h-2 -mt-1 bg-gray-900 transform -translate-x-1/2 rotate-45 dark:bg-gray-700"
        ></div>
    </div>
</div>
