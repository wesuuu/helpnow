<script lang="ts">
    import {
        SvelteFlow,
        Background,
        Controls,
        SvelteFlowProvider,
        type Node,
        type Edge,
        type NodeEventWithPointer,
        type Connection,
    } from "@xyflow/svelte";
    import "@xyflow/svelte/dist/style.css";

    import TriggerNode from "./nodes/TriggerNode.svelte";
    import ActionNode from "./nodes/ActionNode.svelte";
    import DeletableEdge from "./edges/DeletableEdge.svelte";
    import ContextMenu from "./ContextMenu.svelte";

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

    // Node Types
    const nodeTypes = {
        TRIGGER: TriggerNode,
        ACTION: ActionNode,
    };

    // Edge Types
    const edgeTypes = {
        deletable: DeletableEdge,
    };

    // Connection Validation
    const isValidConnection = (connection: Connection | Edge) => {
        // Enforce strict single output from a node.
        // User request: "You should only be able to only create an edge to one of those [handles], not all or some, just one"
        // This means if ANY edge exists from this source node, we block new ones.
        const existingEdgesFromSource = edges.filter(
            (e) => e.source === connection.source,
        );

        if (existingEdgesFromSource.length > 0) return false;

        return true;
    };

    // Context Menu State
    let contextMenu = $state<{
        x: number;
        y: number;
        cx: number;
        cy: number;
        id: string; // Optional for pane
        type: "node" | "edge" | "pane";
    } | null>(null);
    let containerRef: HTMLDivElement;

    function onNodeClick({ node }: { node: Node }) {
        selectedNodeId = node.id;
    }

    const onNodeContextMenu: NodeEventWithPointer<MouseEvent> = ({
        event,
        node,
    }) => {
        event.preventDefault();
        const mouseEvent = event as MouseEvent;
        if (mouseEvent) {
            mouseEvent.preventDefault();
        }

        if (node && (node.type === "ACTION" || node.type === "TRIGGER")) {
            const rect = containerRef?.getBoundingClientRect();
            const x = rect
                ? mouseEvent.clientX - rect.left
                : mouseEvent.clientX;
            const y = rect ? mouseEvent.clientY - rect.top : mouseEvent.clientY;

            contextMenu = {
                x,
                y,
                cx: mouseEvent.clientX,
                cy: mouseEvent.clientY,
                id: node.id,
                type: "node",
            };
        } else {
            contextMenu = null;
        }
    };

    const onEdgeContextMenu = ({
        event,
        edge,
    }: {
        event: MouseEvent;
        edge: Edge;
    }) => {
        event.preventDefault();
        const rect = containerRef?.getBoundingClientRect();
        const x = rect ? event.clientX - rect.left : event.clientX;
        const y = rect ? event.clientY - rect.top : event.clientY;

        contextMenu = {
            x,
            y,
            cx: event.clientX,
            cy: event.clientY,
            id: edge.id,
            type: "edge",
        };
    };

    const onPaneContextMenu = ({ event }: { event: MouseEvent }) => {
        event.preventDefault();
        const rect = containerRef?.getBoundingClientRect();
        const x = rect ? event.clientX - rect.left : event.clientX;
        const y = rect ? event.clientY - rect.top : event.clientY;

        contextMenu = {
            x,
            y,
            cx: event.clientX,
            cy: event.clientY,
            id: "", // No specific ID for pane
            type: "pane",
        };
    };

    function onPaneClick() {
        selectedNodeId = null;
        contextMenu = null;
    }
</script>

<div
    bind:this={containerRef}
    class="h-[600px] w-full border border-gray-200 rounded-lg bg-slate-50 relative group"
>
    <SvelteFlowProvider>
        <SvelteFlow
            bind:nodes
            bind:edges
            {nodeTypes}
            {edgeTypes}
            defaultEdgeOptions={{ type: "deletable", animated: true }}
            {isValidConnection}
            fitView
            onnodeclick={onNodeClick}
            onpointerdown={onPaneClick}
            onnodecontextmenu={onNodeContextMenu}
            onedgecontextmenu={onEdgeContextMenu}
            onpanecontextmenu={onPaneContextMenu}
            onpaneclick={onPaneClick}
        >
            <Background />
            <Controls />
        </SvelteFlow>
        {#if contextMenu}
            <ContextMenu
                onclick={onPaneClick}
                id={contextMenu.id}
                type={contextMenu.type}
                top={contextMenu.y}
                left={contextMenu.x}
                cx={contextMenu.cx}
                cy={contextMenu.cy}
                right={undefined}
                bottom={undefined}
            />
        {/if}
    </SvelteFlowProvider>
</div>
