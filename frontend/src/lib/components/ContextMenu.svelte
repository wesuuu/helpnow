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
        cx?: number;
        cy?: number;
    } = $props();

    const nodes = useNodes();
    const updateNodeInternals = useUpdateNodeInternals();
    const { deleteElements, screenToFlowPosition } = useSvelteFlow();

    function switchHandlePosition() {
        const nodeIndex = nodes.current.findIndex((n) => n.id === id);
        if (nodeIndex !== -1) {
            const node = nodes.current[nodeIndex];
            const currentPos = node.data.handlePosition || "right";
            const newPos = currentPos === "right" ? "left" : "right";

            // Update the node data
            const updatedNode = {
                ...node,
                data: {
                    ...node.data,
                    handlePosition: newPos,
                },
            };

            // Update the nodes store
            const newNodes = [...nodes.current];
            newNodes[nodeIndex] = updatedNode;
            nodes.set(newNodes);

            // Force update node internals to fix edge positions
            // We need to wait for the DOM to update, so requestAnimationFrame or setTimeout might be needed
            // but usually calling it after set is enough if reactive.
            // Let's wrap in setTimeout to be safe as the DOM needs to reflow first.
            setTimeout(() => {
                updateNodeInternals(id);
            }, 0);
        }
        if (onclick) onclick();
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
        <button onclick={switchHandlePosition}>Switch Handles Side</button>
        <button onclick={deleteItem} class="text-red-500 hover:!bg-red-50"
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
