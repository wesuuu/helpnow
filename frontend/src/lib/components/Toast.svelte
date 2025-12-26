<script lang="ts">
    import { toast, type Toast } from "../stores/toast";
    import { fade, fly } from "svelte/transition";

    let { item } = $props<{ item: Toast }>();

    const bgColors = {
        success: "bg-green-500",
        error: "bg-red-500",
        info: "bg-blue-500",
    };

    const iconColors = {
        success: "text-green-100",
        error: "text-red-100",
        info: "text-blue-100",
    };
</script>

<div
    class="{bgColors[
        item.type
    ]} text-white px-6 py-4 rounded shadow-lg flex items-center justify-between mb-4 min-w-[300px]"
    in:fly={{ y: 20, duration: 300 }}
    out:fade={{ duration: 200 }}
    role="alert"
>
    <div class="flex items-center">
        {#if item.type === "success"}
            <!-- Heroicon: check-circle -->
            <svg
                class="w-6 h-6 mr-3 {iconColors.success}"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
            >
                <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                ></path>
            </svg>
        {:else if item.type === "error"}
            <!-- Heroicon: exclamation-circle -->
            <svg
                class="w-6 h-6 mr-3 {iconColors.error}"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
            >
                <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                ></path>
            </svg>
        {:else}
            <!-- Heroicon: information-circle -->
            <svg
                class="w-6 h-6 mr-3 {iconColors.info}"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
            >
                <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                ></path>
            </svg>
        {/if}
        <span>{item.message}</span>
    </div>
    <button
        onclick={() => toast.remove(item.id)}
        class="ml-4 text-white hover:text-gray-200 focus:outline-none"
        aria-label="Close"
    >
        <svg
            class="w-4 h-4"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
        >
            <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M6 18L18 6M6 6l12 12"
            ></path>
        </svg>
    </button>
</div>
