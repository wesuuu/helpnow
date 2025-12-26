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

    // Layout Calculation
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
        const _ = handlePos;
        const __ = data.swapYesNo; // Track this too
        setTimeout(() => {
            // @ts-ignore
            updateNodeInternals(id);
        }, 0);
    });

    // CSS Classes for Container positioning
    const containerClass = $derived.by(() => {
        if (handlePos === "right")
            return "absolute top-1/2 transform -translate-y-1/2 right-0 -mr-3 flex flex-col gap-6";
        if (handlePos === "left")
            return "absolute top-1/2 transform -translate-y-1/2 left-0 -ml-3 flex flex-col gap-6";
        if (handlePos === "bottom")
            return "absolute left-1/2 transform -translate-x-1/2 bottom-0 -mb-3 flex flex-row gap-6";
        if (handlePos === "top")
            return "absolute left-1/2 transform -translate-x-1/2 top-0 -mt-3 flex flex-row gap-6";
        return "";
    });

    // Label positioning (Yes/No tags)
    const labelClass = $derived.by(() => {
        if (handlePos === "right")
            return "left-full ml-2 top-1/2 -translate-y-1/2";
        if (handlePos === "left")
            return "right-full mr-2 top-1/2 -translate-y-1/2";
        if (handlePos === "bottom")
            return "top-full mt-2 left-1/2 -translate-x-1/2";
        if (handlePos === "top")
            return "bottom-full mb-2 left-1/2 -translate-x-1/2";
        return "";
    });
</script>

<div
    class="px-3 py-2 rounded-md bg-white border-2 border-orange-500 shadow-md min-w-[150px]"
    style="width: fit-content;"
>
    <!-- Input Handle -->
    <Handle
        type="target"
        position={targetPosition}
        class="w-3 h-3 bg-gray-400"
    />

    <div class="font-bold text-sm text-orange-700 mb-1">
        {data.label || "Condition"}
    </div>

    <div class="text-xs text-gray-800 mb-2">
        {data.description || "True/False Check"}
    </div>

    <!-- Output Handles Container -->
    <div class={containerClass}>
        {#if data.swapYesNo}
            <!-- False (Red) First -->
            <div class="relative group/handle">
                <Handle
                    type="source"
                    position={sourcePosition}
                    id="false"
                    class="!w-3 !h-3 !bg-red-500 !border-2 !border-white"
                    style="top: 0; position: relative;"
                />
                <span
                    class={`absolute text-[10px] text-red-600 font-bold bg-white px-1 rounded shadow-sm whitespace-nowrap ${labelClass}`}
                    >NO</span
                >
            </div>
            <!-- True (Green) Second -->
            <div class="relative group/handle">
                <Handle
                    type="source"
                    position={sourcePosition}
                    id="true"
                    class="!w-3 !h-3 !bg-green-500 !border-2 !border-white"
                    style="top: 0; position: relative;"
                />
                <span
                    class={`absolute text-[10px] text-green-600 font-bold bg-white px-1 rounded shadow-sm whitespace-nowrap ${labelClass}`}
                    >YES</span
                >
            </div>
        {:else}
            <!-- True (Green) First -->
            <div class="relative group/handle">
                <Handle
                    type="source"
                    position={sourcePosition}
                    id="true"
                    class="!w-3 !h-3 !bg-green-500 !border-2 !border-white"
                    style="top: 0; position: relative;"
                />
                <span
                    class={`absolute text-[10px] text-green-600 font-bold bg-white px-1 rounded shadow-sm whitespace-nowrap ${labelClass}`}
                    >YES</span
                >
            </div>
            <!-- False (Red) Second -->
            <div class="relative group/handle">
                <Handle
                    type="source"
                    position={sourcePosition}
                    id="false"
                    class="!w-3 !h-3 !bg-red-500 !border-2 !border-white"
                    style="top: 0; position: relative;"
                />
                <span
                    class={`absolute text-[10px] text-red-600 font-bold bg-white px-1 rounded shadow-sm whitespace-nowrap ${labelClass}`}
                    >NO</span
                >
            </div>
        {/if}
    </div>
</div>
