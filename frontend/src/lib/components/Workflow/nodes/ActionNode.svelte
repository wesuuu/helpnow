<script lang="ts">
    import {
        Handle,
        Position,
        useUpdateNodeInternals,
        type NodeProps,
    } from "@xyflow/svelte";

    type $$Props = NodeProps;

    let { id, data }: NodeProps = $props();

    const updateNodeInternals = useUpdateNodeInternals();

    // Default to 'left' if not specified
    // Default to 'right' if not specified
    const handlePos = $derived(data.handlePosition || "bottom");

    const sourcePosition = $derived(
        handlePos === "left"
            ? Position.Left
            : handlePos === "top"
              ? Position.Top
              : handlePos === "bottom"
                ? Position.Bottom
                : Position.Right,
    );

    const targetPosition = $derived(
        handlePos === "left"
            ? Position.Right
            : handlePos === "top"
              ? Position.Bottom
              : handlePos === "bottom"
                ? Position.Top
                : Position.Left,
    );

    $effect(() => {
        // Access position to track changes
        const _ = handlePos;
        // Run updateNodeInternals in the next tick to ensure DOM handles are rendered
        setTimeout(() => {
            updateNodeInternals(id);
        }, 0);
    });

    // Label positioning based on handle position
    const inputLabelClass = $derived.by(() => {
        if (handlePos === "right")
            return "absolute right-full mr-2 top-1/2 -translate-y-1/2 text-[10px] bg-slate-100 text-slate-500 px-1 rounded opacity-0 peer-hover:opacity-100 transition-opacity pointer-events-none whitespace-nowrap z-50 shadow-sm border border-slate-200";
        if (handlePos === "left")
            return "absolute left-full ml-2 top-1/2 -translate-y-1/2 text-[10px] bg-slate-100 text-slate-500 px-1 rounded opacity-0 peer-hover:opacity-100 transition-opacity pointer-events-none whitespace-nowrap z-50 shadow-sm border border-slate-200";
        if (handlePos === "bottom")
            return "absolute bottom-full mb-2 left-1/2 -translate-x-1/2 text-[10px] bg-slate-100 text-slate-500 px-1 rounded opacity-0 peer-hover:opacity-100 transition-opacity pointer-events-none whitespace-nowrap z-50 shadow-sm border border-slate-200";
        if (handlePos === "top")
            return "absolute top-full mt-2 left-1/2 -translate-x-1/2 text-[10px] bg-slate-100 text-slate-500 px-1 rounded opacity-0 peer-hover:opacity-100 transition-opacity pointer-events-none whitespace-nowrap z-50 shadow-sm border border-slate-200";
        return "";
    });

    const outputLabelClass = $derived.by(() => {
        if (handlePos === "right")
            return "absolute left-full ml-2 top-1/2 -translate-y-1/2 text-[10px] bg-slate-100 text-slate-500 px-1 rounded opacity-0 peer-hover:opacity-100 transition-opacity pointer-events-none whitespace-nowrap z-50 shadow-sm border border-slate-200";
        if (handlePos === "left")
            return "absolute right-full mr-2 top-1/2 -translate-y-1/2 text-[10px] bg-slate-100 text-slate-500 px-1 rounded opacity-0 peer-hover:opacity-100 transition-opacity pointer-events-none whitespace-nowrap z-50 shadow-sm border border-slate-200";
        if (handlePos === "bottom")
            return "absolute top-full mt-2 left-1/2 -translate-x-1/2 text-[10px] bg-slate-100 text-slate-500 px-1 rounded opacity-0 peer-hover:opacity-100 transition-opacity pointer-events-none whitespace-nowrap z-50 shadow-sm border border-slate-200";
        if (handlePos === "top")
            return "absolute bottom-full mb-2 left-1/2 -translate-x-1/2 text-[10px] bg-slate-100 text-slate-500 px-1 rounded opacity-0 peer-hover:opacity-100 transition-opacity pointer-events-none whitespace-nowrap z-50 shadow-sm border border-slate-200";
        return "";
    });
</script>

<div
    class="px-3 py-2 rounded-md bg-white border-2 border-blue-500 shadow-md min-w-[180px] relative"
    style="width: fit-content;"
>
    <!-- Input Handle: Opposite to Output -->
    <Handle
        type="target"
        position={targetPosition}
        class="w-3 h-3 bg-gray-400 peer"
    />
    <span class={inputLabelClass}>Input</span>

    <div class="font-bold text-sm text-blue-700 mb-1">
        {data.label || "Action"}
    </div>

    <div class="text-xs text-gray-800 mb-2">
        {data.action || "Select Action"}
    </div>

    <!-- Output Handle -->
    <Handle
        type="source"
        position={sourcePosition}
        id="default"
        class="w-3 h-3 bg-gray-400 peer"
    />
    <span class={outputLabelClass}>Output</span>
</div>
