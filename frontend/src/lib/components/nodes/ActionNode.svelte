<script lang="ts">
    import { Handle, Position, type NodeProps } from "@xyflow/svelte";

    type $$Props = NodeProps;

    export let data: {
        label: string;
        action?: string;
        handlePosition?: "left" | "right"; // New property
    };

    // Default to 'right' if not specified
    $: position = data.handlePosition || "right";
    $: isRight = position === "right";
</script>

<div
    class="px-3 py-2 rounded-md bg-white border-2 border-blue-500 shadow-md min-w-[180px]"
>
    <!-- Input Handle: Opposite to Output -->
    <Handle
        type="target"
        position={isRight ? Position.Left : Position.Right}
        class="w-3 h-3 bg-gray-400"
    />

    <div class="font-bold text-sm text-blue-700 mb-1">ACTION</div>

    <div class="text-xs text-gray-800 mb-4">
        {data.label || "Action Node"}
    </div>

    <!-- Output Handles Container -->
    <div
        class={`absolute top-1/2 transform -translate-y-1/2 flex flex-col gap-3 ${isRight ? "right-0 -mr-3" : "left-0 -ml-3"}`}
    >
        <!-- Success -->
        <div class="relative group">
            <Handle
                type="source"
                position={isRight ? Position.Right : Position.Left}
                id="success"
                class="!w-3 !h-3 !bg-green-500 !border-2 !border-white"
                style="top: 0; position: relative;"
            />
            <span
                class={`absolute top-1/2 -translate-y-1/2 text-[10px] text-green-600 font-bold opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap bg-white px-1 rounded shadow-sm ${isRight ? "left-full ml-2" : "right-full mr-2"}`}
                >OK</span
            >
        </div>

        <!-- Failure -->
        <div class="relative group">
            <Handle
                type="source"
                position={isRight ? Position.Right : Position.Left}
                id="failure"
                class="!w-3 !h-3 !bg-red-500 !border-2 !border-white"
                style="top: 0; position: relative;"
            />
            <span
                class={`absolute top-1/2 -translate-y-1/2 text-[10px] text-red-600 font-bold opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap bg-white px-1 rounded shadow-sm ${isRight ? "left-full ml-2" : "right-full mr-2"}`}
                >FAIL</span
            >
        </div>

        <!-- Always -->
        <div class="relative group">
            <Handle
                type="source"
                position={isRight ? Position.Right : Position.Left}
                id="always"
                class="!w-3 !h-3 !bg-gray-500 !border-2 !border-white"
                style="top: 0; position: relative;"
            />
            <span
                class={`absolute top-1/2 -translate-y-1/2 text-[10px] text-gray-600 font-bold opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap bg-white px-1 rounded shadow-sm ${isRight ? "left-full ml-2" : "right-full mr-2"}`}
                >ANY</span
            >
        </div>
    </div>
</div>
