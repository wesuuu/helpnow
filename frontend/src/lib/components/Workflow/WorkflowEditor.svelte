<script lang="ts">
    import { SvelteFlowProvider } from "@xyflow/svelte";
    import WorkflowCanvas from "./WorkflowCanvas.svelte";
    import type { Node, Edge } from "@xyflow/svelte";

    // Props
    let {
        nodes = $bindable([]),
        edges = $bindable([]),
        selectedNodeId = $bindable(null),
    } = $props<{
        nodes: Node[];
        edges: Edge[];
        selectedNodeId?: string | null;
    }>();
</script>

<div class="h-full w-full relative group">
    <SvelteFlowProvider>
        <WorkflowCanvas bind:nodes bind:edges bind:selectedNodeId />
    </SvelteFlowProvider>
</div>
