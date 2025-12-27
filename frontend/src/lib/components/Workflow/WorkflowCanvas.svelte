<script lang="ts">
    import {
        SvelteFlow,
        Background,
        Controls,
        useSvelteFlow,
        type Node,
        type Edge,
        type NodeEventWithPointer,
        type Connection,
        addEdge,
    } from "@xyflow/svelte";
    import DeletableEdge from "./edges/DeletableEdge.svelte";
    import ContextMenu from "../ContextMenu.svelte";
    import "@xyflow/svelte/dist/style.css";
    import TriggerNode from "./nodes/TriggerNode.svelte";
    import ActionNode from "./nodes/ActionNode.svelte";
    import ConditionNode from "./nodes/ConditionNode.svelte";

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

    const { screenToFlowPosition } = useSvelteFlow();

    // Node Types
    const nodeTypes = {
        TRIGGER: TriggerNode,
        ACTION: ActionNode,
        CONDITION: ConditionNode,
    };

    // Edge Types
    const edgeTypes = {
        deletable: DeletableEdge,
    };

    // Connection Validation
    const isValidConnection = (connection: Connection | Edge) => {
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
        id: string;
        type: "node" | "edge" | "pane";
        nodeType?: string;
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
        if (mouseEvent) mouseEvent.preventDefault();

        if (node) {
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
                nodeType: node.type,
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
            id: "",
            type: "pane",
        };
    };

    function onPaneClick() {
        selectedNodeId = null;
        contextMenu = null;
    }

    // --- Drag and Drop ---
    function onDragOver(event: DragEvent) {
        event.preventDefault();
        if (event.dataTransfer) {
            event.dataTransfer.dropEffect = "move";
        }
    }

    function onDrop(event: DragEvent) {
        event.preventDefault();
        if (!event.dataTransfer) return;

        const type = event.dataTransfer.getData("application/svelteflow");
        const dataStr = event.dataTransfer.getData(
            "application/svelteflow-data",
        );

        if (typeof type === "undefined" || !type) {
            return;
        }

        const position = screenToFlowPosition({
            x: event.clientX,
            y: event.clientY,
        });

        // Add Node
        const id = `n-${Date.now()}`;
        const defaultData = {
            label: type, // temporary fallback
            handlePosition: "bottom",
            trigger_type: "EVENT",
            site_ids: [],
            audience_ids: [],
        };

        let extraData = {};
        if (dataStr) {
            try {
                extraData = JSON.parse(dataStr);
            } catch (e) {}
        }

        const newNode = {
            id,
            type,
            position,
            data: { ...defaultData, ...extraData },
        };

        nodes = [...nodes, newNode];
    }
</script>

<div
    bind:this={containerRef}
    class="h-full w-full relative group"
    ondragover={onDragOver}
    ondrop={onDrop}
    role="application"
>
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
        onconnect={(params) => {
            edges = addEdge(params, edges);
        }}
    >
        <Background />
        <Controls />
    </SvelteFlow>
    {#if contextMenu}
        <ContextMenu
            {...contextMenu}
            onclick={onPaneClick}
            top={contextMenu.y}
            left={contextMenu.x}
            right={undefined}
            bottom={undefined}
            nodeType={contextMenu.nodeType}
        />
    {/if}
</div>
