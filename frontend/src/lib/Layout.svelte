<script lang="ts">
    import Navbar from "./Navbar.svelte";
    import Sidebar from "./Sidebar.svelte";
    let { children } = $props();
</script>

<div
    class="min-h-screen bg-gray-50 dark:bg-gray-900 flex flex-col transition-colors duration-200"
>
    <Navbar />

    <div class="flex flex-1 overflow-hidden">
        <Sidebar />

        <main class="flex-1 overflow-y-auto p-4 md:p-8">
            {@render children?.()}
        </main>
    </div>
</div>
