<script lang="ts">
    import {
        SvelteFlow,
        Background,
        Controls,
        SvelteFlowProvider,
        type Node,
        type Edge,
        type NodeEventWithPointer,
    } from "@xyflow/svelte";
    import "@xyflow/svelte/dist/style.css";

    import TriggerNode from "./nodes/TriggerNode.svelte";
    import ActionNode from "./nodes/ActionNode.svelte";
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

    // Context Menu State
    let contextMenu = $state<{ x: number; y: number; nodeId: string } | null>(
        null,
    );
    let containerRef: HTMLDivElement;

    function onNodeClick(e: CustomEvent) {
        selectedNodeId = e.detail.node.id;
    }

    const onNodeContextMenu: NodeEventWithPointer<MouseEvent> = ({
        event,
        node,
    }) => {
        // e.detail.event is the original MouseEvent
        event.preventDefault();
        const mouseEvent = event as MouseEvent;
        if (mouseEvent) {
            mouseEvent.preventDefault();
            console.log("Context Menu Triggered", mouseEvent);
        }

        // Helper to find the node type
        if (node && (node.type === "ACTION" || node.type === "TRIGGER")) {
            // Get the container rect to position relative to it
            const rect = containerRef?.getBoundingClientRect();

            // Calculate x/y relative to the flow container
            const x = rect
                ? mouseEvent.clientX - rect.left
                : mouseEvent.clientX;
            const y = rect ? mouseEvent.clientY - rect.top : mouseEvent.clientY;

            console.log("Setting Context Menu at", { x, y });

            contextMenu = {
                x,
                y,
                nodeId: node.id,
            };
        } else {
            console.log("Node not actionable for context menu");
            contextMenu = null;
        }
    };

    function onPaneClick() {
        selectedNodeId = null;
        contextMenu = null;
    }

    function toggleHandlePosition() {
        if (!contextMenu) return;
        const index = nodes.findIndex((n) => n.id === contextMenu!.nodeId);
        if (index !== -1) {
            const node = nodes[index];
            const current = node.data.handlePosition || "right";
            const newVal = current === "right" ? "left" : "right";

            // Re-assign to trigger reactivity
            nodes[index] = {
                ...node,
                data: {
                    ...node.data,
                    handlePosition: newVal,
                },
            };
        }
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
            fitView
            onpointerdown={onPaneClick}
            onnodecontextmenu={onNodeContextMenu}
            onpaneclick={onPaneClick}
        >
            <Background />
            <Controls />
        </SvelteFlow>
        {#if contextMenu}
            <ContextMenu
                onclick={onPaneClick}
                id={contextMenu.nodeId}
                top={contextMenu.y}
                left={contextMenu.x}
                right={undefined}
                bottom={undefined}
            />
        {/if}
    </SvelteFlowProvider>
</div>
