<script lang="ts">
    import {
        BaseEdge,
        EdgeLabel,
        getBezierPath,
        useSvelteFlow,
        type EdgeProps,
    } from "@xyflow/svelte";

    let {
        id,
        sourceX,
        sourceY,
        targetX,
        targetY,
        sourcePosition,
        targetPosition,
        style = undefined,
        markerEnd = undefined,
    }: EdgeProps = $props();

    const { deleteElements } = useSvelteFlow();

    let [edgePath, labelX, labelY] = $derived(
        getBezierPath({
            sourceX,
            sourceY,
            sourcePosition,
            targetX,
            targetY,
            targetPosition,
        }),
    );
</script>

<BaseEdge path={edgePath} {markerEnd} {style} />

<EdgeLabel x={labelX} y={labelY}>
    <button
        class="w-5 h-5 bg-gray-200 hover:bg-red-500 hover:text-white rounded-full flex items-center justify-center text-xs leading-none border border-gray-400 transition-colors shadow-sm nodrag nopan pointer-events-auto"
        onclick={(evt) => {
            evt.stopPropagation();
            deleteElements({ edges: [{ id }] });
        }}
        aria-label="delete edge"
    >
        ×
    </button>
</EdgeLabel>
