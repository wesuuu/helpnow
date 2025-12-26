<script lang="ts">
    import {
        useNodes,
        useUpdateNodeInternals,
        useSvelteFlow,
    } from "@xyflow/svelte";

    let {
        id,
        top,
        left,
        right,
        bottom,
        onclick,
        type = "node",
        nodeType, // Add nodeType prop
        cx,
        cy,
    }: {
        id: string;
        top: number | undefined;
        left: number | undefined;
        right: number | undefined;
        bottom: number | undefined;
        onclick: () => void;
        type?: "node" | "edge" | "pane";
        nodeType?: string; // TRIGGER, ACTION, CONDITION
        cx?: number;
        cy?: number;
    } = $props();

    const nodes = useNodes();
    const updateNodeInternals = useUpdateNodeInternals();
    const { deleteElements, screenToFlowPosition } = useSvelteFlow();

    function setHandlePosition(pos: "top" | "right" | "bottom" | "left") {
        const nodeIndex = nodes.current.findIndex((n) => n.id === id);
        if (nodeIndex !== -1) {
            const node = nodes.current[nodeIndex];

            // Update the node data
            const updatedNode = {
                ...node,
                data: {
                    ...node.data,
                    handlePosition: pos,
                },
            };

            // Update the nodes store
            const newNodes = [...nodes.current];
            newNodes[nodeIndex] = updatedNode;
            nodes.set(newNodes);

            // Force update node internals to fix edge positions
            setTimeout(() => {
                updateNodeInternals(id);
            }, 0);
        }
        if (onclick) onclick();
    }

    // For Action/Condition nodes: set axis (Vertical vs Horizontal)
    // Horizontal -> Right (default)
    // Vertical -> Bottom (default for vertical)
    function setHandleAxis(axis: "horizontal" | "vertical") {
        if (axis === "horizontal") {
            setHandlePosition("right");
        } else {
            setHandlePosition("bottom");
        }
    }

    function swapYesNo() {
        const nodeIndex = nodes.current.findIndex((n) => n.id === id);
        if (nodeIndex !== -1) {
            const node = nodes.current[nodeIndex];
            const currentSwap = node.data.swapYesNo || false;

            const updatedNode = {
                ...node,
                data: {
                    ...node.data,
                    swapYesNo: !currentSwap,
                },
            };

            const newNodes = [...nodes.current];
            newNodes[nodeIndex] = updatedNode;
            nodes.set(newNodes);

            setTimeout(() => {
                updateNodeInternals(id);
            }, 0);
        }
        if (onclick) onclick();
    }

    function swapHandlePosition() {
        const nodeIndex = nodes.current.findIndex((n) => n.id === id);
        if (nodeIndex !== -1) {
            const node = nodes.current[nodeIndex];
            const currentPos = node.data.handlePosition || "right";
            let newPos;

            switch (currentPos) {
                case "right":
                    newPos = "left";
                    break;
                case "left":
                    newPos = "right";
                    break;
                case "top":
                    newPos = "bottom";
                    break;
                case "bottom":
                    newPos = "top";
                    break;
                default:
                    newPos = "right";
            }

            // Reuse setHandlePosition logic essentially
            setHandlePosition(newPos as "top" | "right" | "bottom" | "left");
        }
    }

    async function deleteItem() {
        if (type === "node") {
            // We need to pass the full node object or just ID? deleteElements takes { nodes: [{ id: '...' }], edges: [] }
            // Actually it accepts partials with id.
            await deleteElements({ nodes: [{ id }] });
        } else if (type === "edge") {
            await deleteElements({ edges: [{ id }] });
        }
        if (onclick) onclick();
    }

    function addActionNode() {
        if (cx === undefined || cy === undefined) return;
        const pos = screenToFlowPosition({ x: cx, y: cy });
        const newNode = {
            id: `n-${Date.now()}`,
            type: "ACTION",
            position: pos,
            data: {
                label: "New Action",
                handlePosition: "right", // Default
            },
        };

        const newNodes = [...nodes.current, newNode];
        nodes.set(newNodes);
        if (onclick) onclick();
    }
</script>

<div
    style="top: {top}px; left: {left}px; right: {right}px; bottom: {bottom}px;"
    class="context-menu"
    {onclick}
    onpointerdown={(e) => e.stopPropagation()}
>
    {#if type === "pane"}
        <p style="margin: 0.5em;"><small>Menu</small></p>
        <button onclick={addActionNode}>Add Action Node</button>
    {:else}
        <p style="margin: 0.5em;">
            <small>{type}: {id}</small>
        </p>
    {/if}

    {#if type === "node"}
        <div class="mb-2 pb-2 border-b border-gray-100">
            <p class="px-4 py-1 text-xs text-gray-500 font-semibold">
                Handle Position
            </p>

            {#if nodeType === "TRIGGER"}
                <!-- Trigger: Top, Bottom, Left, Right -->
                <div class="grid grid-cols-2 gap-1 px-2 mb-2">
                    <button
                        class="text-xs border rounded hover:bg-gray-50 text-center justify-center pointer-events-auto"
                        onclick={() => setHandlePosition("top")}>Top</button
                    >
                    <button
                        class="text-xs border rounded hover:bg-gray-50 text-center justify-center pointer-events-auto"
                        onclick={() => setHandlePosition("bottom")}
                        >Bottom</button
                    >
                    <button
                        class="text-xs border rounded hover:bg-gray-50 text-center justify-center pointer-events-auto"
                        onclick={() => setHandlePosition("left")}>Left</button
                    >
                    <button
                        class="text-xs border rounded hover:bg-gray-50 text-center justify-center pointer-events-auto"
                        onclick={() => setHandlePosition("right")}>Right</button
                    >
                </div>
            {:else if nodeType === "ACTION"}
                <!-- Action: Vertical (Top/Bottom), Horizontal (Left/Right) -->
                <div class="grid grid-cols-1 gap-1 px-2 mb-2">
                    <button
                        class="text-xs border rounded hover:bg-gray-50 text-center justify-center pointer-events-auto"
                        onclick={() => setHandleAxis("vertical")}
                        >Top - Bottom (Vertical)</button
                    >
                    <button
                        class="text-xs border rounded hover:bg-gray-50 text-center justify-center pointer-events-auto"
                        onclick={() => setHandleAxis("horizontal")}
                        >Left - Right (Horizontal)</button
                    >
                </div>
            {:else if nodeType === "CONDITION"}
                <!-- Condition: Vertical, Horizontal, Swap Yes/No -->
                <div class="grid grid-cols-1 gap-1 px-2 mb-2">
                    <button
                        class="text-xs border rounded hover:bg-gray-50 text-center justify-center pointer-events-auto"
                        onclick={() => setHandleAxis("vertical")}
                        >Top - Bottom (Vertical)</button
                    >
                    <button
                        class="text-xs border rounded hover:bg-gray-50 text-center justify-center pointer-events-auto"
                        onclick={() => setHandleAxis("horizontal")}
                        >Left - Right (Horizontal)</button
                    >
                </div>
                <!-- Separate section for "Swap Yes/No" -->
                <button
                    onclick={swapYesNo}
                    class="text-blue-600 hover:bg-blue-50 font-medium w-full text-center pointer-events-auto"
                >
                    Swap Yes/No ⇄
                </button>
            {:else}
                <!-- Default Fallback (if any other node type) -->
                <div class="grid grid-cols-2 gap-1 px-2 mb-2">
                    <button
                        class="text-xs border rounded hover:bg-gray-50 text-center justify-center pointer-events-auto"
                        onclick={() => setHandlePosition("top")}>Top</button
                    >
                    <button
                        class="text-xs border rounded hover:bg-gray-50 text-center justify-center pointer-events-auto"
                        onclick={() => setHandlePosition("bottom")}
                        >Bottom</button
                    >
                    <button
                        class="text-xs border rounded hover:bg-gray-50 text-center justify-center pointer-events-auto"
                        onclick={() => setHandlePosition("left")}>Left</button
                    >
                    <button
                        class="text-xs border rounded hover:bg-gray-50 text-center justify-center pointer-events-auto"
                        onclick={() => setHandlePosition("right")}>Right</button
                    >
                </div>
            {/if}
        </div>
        <button
            onclick={deleteItem}
            class="text-red-500 hover:!bg-red-50 pointer-events-auto"
            >Delete Node</button
        >
    {/if}
    {#if type === "edge"}
        <button onclick={deleteItem} class="text-red-500 hover:!bg-red-50"
            >Delete Edge</button
        >
    {/if}
</div>

<style>
    .context-menu {
        background: white;
        border-style: solid;
        box-shadow: 10px 19px 20px rgba(0, 0, 0, 10%);
        position: absolute;
        z-index: 10;
        min-width: 150px;
        padding: 0.5rem 0;
        border-radius: 0.375rem; /* rounded-md */
        border-color: #e5e7eb; /* gray-200 */
        border-width: 1px;
    }

    .context-menu button {
        border: none;
        display: block;
        padding: 0.5em 1em;
        text-align: left;
        width: 100%;
        background: transparent;
        cursor: pointer;
    }

    .context-menu button:hover {
        background: #f3f4f6; /* gray-100 */
    }
</style>
