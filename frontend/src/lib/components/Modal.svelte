<script lang="ts">
    let {
        onclose,
        children,
        size = "md",
    } = $props<{
        onclose?: () => void;
        children?: import("svelte").Snippet;
        size?: "sm" | "md" | "lg" | "xl" | "2xl";
    }>();

    function close() {
        if (onclose) onclose();
    }

    const maxWidthClasses: Record<string, string> = {
        sm: "sm:max-w-sm",
        md: "sm:max-w-lg",
        lg: "sm:max-w-xl",
        xl: "sm:max-w-4xl",
        "2xl": "sm:max-w-6xl",
    };
</script>

<div
    class="fixed z-50 inset-0 overflow-y-auto"
    aria-labelledby="modal-title"
    role="dialog"
    aria-modal="true"
>
    <div
        class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0"
    >
        <div
            class="fixed inset-0 bg-gray-500/75 transition-opacity"
            aria-hidden="true"
            onclick={close}
            role="presentation"
        ></div>

        <span
            class="hidden sm:inline-block sm:align-middle sm:h-screen"
            aria-hidden="true">&#8203;</span
        >

        <div
            class="relative inline-block align-bottom bg-white rounded-lg px-4 pt-5 pb-4 text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:w-full sm:p-6 z-50 {maxWidthClasses[
                size
            ] || maxWidthClasses.md}"
        >
            {@render children?.()}
        </div>
    </div>
</div>
