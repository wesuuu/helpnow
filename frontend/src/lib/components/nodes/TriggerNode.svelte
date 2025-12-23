<script lang="ts">
    import { Handle, Position, type NodeProps } from "@xyflow/svelte";

    type $$Props = NodeProps;
    export let data: any;
</script>

<div
    class="px-3 py-2 rounded-md bg-white border-2 border-purple-500 shadow-md min-w-[150px]"
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

    <Handle
        type="source"
        position={Position.Bottom}
        class="w-3 h-3 bg-gray-400"
    />
</div>
