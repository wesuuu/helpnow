<script lang="ts">
    import { toast } from "../stores/toast";
    import ToastItem from "./Toast.svelte";
</script>

<div
    class="fixed top-4 right-4 z-[100] flex flex-col items-end pointer-events-none"
>
    <div class="pointer-events-auto">
        {#each $toast as item (item.id)}
            <ToastItem {item} />
        {/each}
    </div>
</div>
