<script lang="ts">
    import { Handle, Position, type NodeProps } from "@xyflow/svelte";

    type $$Props = NodeProps;
</script>

<div
    class="px-3 py-2 rounded-md bg-white border-2 border-purple-500 shadow-md min-w-[150px]"
>
    <div class="font-bold text-sm text-purple-700 mb-1">TRIGGER</div>

    <div class="text-xs text-gray-600">
        <!-- Label would typically be passed via data -->
        Event Trigger
    </div>

    <Handle
        type="source"
        position={Position.Bottom}
        class="w-3 h-3 bg-gray-400"
    />
</div>
