<script lang="ts">
    import {
        Handle,
        Position,
        useUpdateNodeInternals,
        type NodeProps,
    } from "@xyflow/svelte";

    const updateNodeInternals = useUpdateNodeInternals();

    type $$Props = NodeProps;

    let { id, data }: NodeProps = $props();

    // Parse position
    const handlePos = $derived(data.handlePosition || "right");
    const position = $derived(
        handlePos === "left"
            ? Position.Left
            : handlePos === "top"
              ? Position.Top
              : handlePos === "bottom"
                ? Position.Bottom
                : Position.Right,
    );

    $effect(() => {
        // Access position to track changes
        const _ = handlePos;
        // Run updateNodeInternals
        setTimeout(() => {
            // @ts-ignore
            updateNodeInternals(id);
        }, 0);
    });
</script>

<div
    class="px-3 py-2 rounded-md bg-white border-2 border-purple-500 shadow-md min-w-[150px]"
    style="width: fit-content;"
>
    <!-- Main Label (Custom Name) -->
    <div class="font-bold text-sm text-purple-700 mb-1">
        {data.label || "Trigger"}
    </div>

    <!-- Subtitle (Type) -->
    <div class="text-xs text-gray-600">
        {data.trigger_type === "EVENT" && data.trigger_event
            ? data.trigger_event
            : data.trigger_type || "Select Type"}
    </div>

    <!-- Output handle -->
    <Handle type="source" {position} id="default" class="w-3 h-3 bg-gray-400" />
</div>
