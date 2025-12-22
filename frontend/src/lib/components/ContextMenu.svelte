<script lang="ts">
    import { useNodes } from "@xyflow/svelte";

    let {
        id,
        top,
        left,
        right,
        bottom,
        onclick,
    }: {
        id: string;
        top: number | undefined;
        left: number | undefined;
        right: number | undefined;
        bottom: number | undefined;
        onclick: () => void;
    } = $props();

    const nodes = useNodes();

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
        }
        if (onclick) onclick();
    }
</script>

<div
    style="top: {top}px; left: {left}px; right: {right}px; bottom: {bottom}px;"
    class="context-menu"
    {onclick}
    onpointerdown={(e) => e.stopPropagation()}
>
    <p style="margin: 0.5em;">
        <small>node: {id}</small>
    </p>
    <button onclick={switchHandlePosition}>Switch Handles Side</button>
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
