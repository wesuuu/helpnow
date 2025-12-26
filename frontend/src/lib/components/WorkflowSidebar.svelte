<script lang="ts">
    // WorkflowSidebar.svelte
    // Displays a list of available nodes that can be dragged onto the canvas.
    // Also updated to support click-to-add via callback.

    let { onAddNode } = $props<{
        onAddNode?: (type: string, data?: any) => void;
    }>();

    function onDragStart(event: DragEvent, nodeType: string, data?: any) {
        if (!event.dataTransfer) return;
        event.dataTransfer.setData("application/svelteflow", nodeType);
        if (data) {
            event.dataTransfer.setData(
                "application/svelteflow-data",
                JSON.stringify(data),
            );
        }
        event.dataTransfer.effectAllowed = "move";
    }

    function handleNodeClick(nodeType: string, data?: any) {
        if (onAddNode) {
            onAddNode(nodeType, data);
        }
    }
</script>

<div
    class="h-full flex flex-col bg-white border-r border-gray-200 text-gray-900 w-64"
>
    <div class="p-4 border-b border-gray-200">
        <h2
            class="font-semibold text-sm uppercase tracking-wider text-gray-500"
        >
            Nodes
        </h2>
    </div>

    <div class="flex-1 overflow-y-auto p-4 space-y-6">
        <!-- TRIGGERS -->
        <div>
            <h3 class="text-xs font-bold text-gray-500 uppercase mb-3 px-1">
                Triggers
            </h3>
            <div class="space-y-2">
                <div
                    class="bg-white p-3 rounded cursor-move hover:bg-gray-50 hover:ring-1 hover:ring-indigo-500 transition-all border border-gray-200 flex items-center gap-3 shadow-sm"
                    draggable={true}
                    role="button"
                    tabindex="0"
                    aria-label="Webhook Trigger"
                    ondragstart={(e) =>
                        onDragStart(e, "TRIGGER", {
                            label: "Webhook",
                            trigger_type: "EVENT",
                            trigger_event: "webhook_received",
                        })}
                    onclick={() =>
                        handleNodeClick("TRIGGER", {
                            label: "Webhook",
                            trigger_type: "EVENT",
                            trigger_event: "webhook_received",
                        })}
                    onkeydown={(e) =>
                        e.key === "Enter" &&
                        handleNodeClick("TRIGGER", {
                            label: "Webhook",
                            trigger_type: "EVENT",
                            trigger_event: "webhook_received",
                        })}
                >
                    <div class="p-1.5 bg-purple-100 rounded text-purple-600">
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="w-4 h-4"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            ><path
                                d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"
                            ></path><polyline points="15 3 21 3 21 9"
                            ></polyline><line x1="10" y1="14" x2="21" y2="3"
                            ></line></svg
                        >
                    </div>
                    <span class="text-sm font-medium">Webhook</span>
                </div>

                <div
                    class="bg-white p-3 rounded cursor-move hover:bg-gray-50 hover:ring-1 hover:ring-indigo-500 transition-all border border-gray-200 flex items-center gap-3 shadow-sm"
                    draggable={true}
                    role="button"
                    tabindex="0"
                    aria-label="Schedule Trigger"
                    ondragstart={(e) =>
                        onDragStart(e, "TRIGGER", {
                            label: "Schedule",
                            trigger_type: "SCHEDULE",
                            cron: "0 9 * * *",
                        })}
                    onclick={() =>
                        handleNodeClick("TRIGGER", {
                            label: "Schedule",
                            trigger_type: "SCHEDULE",
                            cron: "0 9 * * *",
                        })}
                    onkeydown={(e) =>
                        e.key === "Enter" &&
                        handleNodeClick("TRIGGER", {
                            label: "Schedule",
                            trigger_type: "SCHEDULE",
                            cron: "0 9 * * *",
                        })}
                >
                    <div class="p-1.5 bg-blue-100 rounded text-blue-600">
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="w-4 h-4"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            ><circle cx="12" cy="12" r="10"></circle><polyline
                                points="12 6 12 12 16 14"
                            ></polyline></svg
                        >
                    </div>
                    <span class="text-sm font-medium">Schedule</span>
                </div>

                <!-- STUB: Email Opened -->
                <div
                    class="bg-white p-3 rounded cursor-move hover:bg-gray-50 hover:ring-1 hover:ring-indigo-500 transition-all border border-gray-200 flex items-center gap-3 shadow-sm"
                    draggable={true}
                    role="button"
                    tabindex="0"
                    aria-label="Email Opened Trigger"
                    ondragstart={(e) =>
                        onDragStart(e, "TRIGGER", {
                            label: "Email Opened",
                            trigger_type: "EVENT",
                            trigger_event: "email_opened",
                        })}
                    onclick={() =>
                        handleNodeClick("TRIGGER", {
                            label: "Email Opened",
                            trigger_type: "EVENT",
                            trigger_event: "email_opened",
                        })}
                    onkeydown={(e) =>
                        e.key === "Enter" &&
                        handleNodeClick("TRIGGER", {
                            label: "Email Opened",
                            trigger_type: "EVENT",
                            trigger_event: "email_opened",
                        })}
                >
                    <div class="p-1.5 bg-yellow-100 rounded text-yellow-600">
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="w-4 h-4"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            ><path
                                d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"
                            ></path><polyline points="22,6 12,13 2,6"
                            ></polyline></svg
                        >
                    </div>
                    <span class="text-sm font-medium">Email Opened</span>
                </div>
            </div>
        </div>

        <!-- ACTIONS -->
        <div>
            <h3 class="text-xs font-bold text-gray-500 uppercase mb-3 px-1">
                Actions
            </h3>
            <div class="space-y-2">
                <div
                    class="bg-white p-3 rounded cursor-move hover:bg-gray-50 hover:ring-1 hover:ring-indigo-500 transition-all border border-gray-200 flex items-center gap-3 shadow-sm"
                    draggable={true}
                    role="button"
                    tabindex="0"
                    aria-label="Send Email Action"
                    ondragstart={(e) =>
                        onDragStart(e, "ACTION", {
                            label: "Send Email",
                            action: "Send Email",
                        })}
                    onclick={() =>
                        handleNodeClick("ACTION", {
                            label: "Send Email",
                            action: "Send Email",
                        })}
                    onkeydown={(e) =>
                        e.key === "Enter" &&
                        handleNodeClick("ACTION", {
                            label: "Send Email",
                            action: "Send Email",
                        })}
                >
                    <div class="p-1.5 bg-indigo-100 rounded text-indigo-600">
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="w-4 h-4"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            ><path
                                d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"
                            ></path><polyline points="22,6 12,13 2,6"
                            ></polyline></svg
                        >
                    </div>
                    <span class="text-sm font-medium">Send Email</span>
                </div>

                <div
                    class="bg-white p-3 rounded cursor-move hover:bg-gray-50 hover:ring-1 hover:ring-indigo-500 transition-all border border-gray-200 flex items-center gap-3 shadow-sm"
                    draggable={true}
                    role="button"
                    tabindex="0"
                    aria-label="Update DB Action"
                    ondragstart={(e) =>
                        onDragStart(e, "ACTION", {
                            label: "Update DB",
                            action: "Update DB",
                        })}
                    onclick={() =>
                        handleNodeClick("ACTION", {
                            label: "Update DB",
                            action: "Update DB",
                        })}
                    onkeydown={(e) =>
                        e.key === "Enter" &&
                        handleNodeClick("ACTION", {
                            label: "Update DB",
                            action: "Update DB",
                        })}
                >
                    <div class="p-1.5 bg-teal-100 rounded text-teal-600">
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="w-4 h-4"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            ><ellipse cx="12" cy="5" rx="9" ry="3"
                            ></ellipse><path
                                d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"
                            ></path><path
                                d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"
                            ></path></svg
                        >
                    </div>
                    <span class="text-sm font-medium">Update DB</span>
                </div>

                <div
                    class="bg-white p-3 rounded cursor-move hover:bg-gray-50 hover:ring-1 hover:ring-indigo-500 transition-all border border-gray-200 flex items-center gap-3 shadow-sm"
                    draggable={true}
                    role="button"
                    tabindex="0"
                    aria-label="HTTP Request Action"
                    ondragstart={(e) =>
                        onDragStart(e, "ACTION", {
                            label: "HTTP Request",
                            action: "HTTP Request",
                        })}
                    onclick={() =>
                        handleNodeClick("ACTION", {
                            label: "HTTP Request",
                            action: "HTTP Request",
                        })}
                    onkeydown={(e) =>
                        e.key === "Enter" &&
                        handleNodeClick("ACTION", {
                            label: "HTTP Request",
                            action: "HTTP Request",
                        })}
                >
                    <div class="p-1.5 bg-sky-100 rounded text-sky-600">
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="w-4 h-4"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            ><circle cx="12" cy="12" r="10"></circle><line
                                x1="2"
                                y1="12"
                                x2="22"
                                y2="12"
                            ></line><path
                                d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"
                            ></path></svg
                        >
                    </div>
                    <span class="text-sm font-medium">HTTP Request</span>
                </div>
            </div>
        </div>

        <!-- LOGIC -->
        <div>
            <h3 class="text-xs font-bold text-gray-500 uppercase mb-3 px-1">
                Logic
            </h3>
            <div class="space-y-2">
                <div
                    class="bg-white p-3 rounded cursor-move hover:bg-gray-50 hover:ring-1 hover:ring-indigo-500 transition-all border border-gray-200 flex items-center gap-3 shadow-sm"
                    draggable={true}
                    role="button"
                    tabindex="0"
                    aria-label="If / Else Logic"
                    ondragstart={(e) =>
                        onDragStart(e, "CONDITION", {
                            label: "If / Else",
                            description: "Conditional Logic",
                        })}
                    onclick={() =>
                        handleNodeClick("CONDITION", {
                            label: "If / Else",
                            description: "Conditional Logic",
                        })}
                    onkeydown={(e) =>
                        e.key === "Enter" &&
                        handleNodeClick("CONDITION", {
                            label: "If / Else",
                            description: "Conditional Logic",
                        })}
                >
                    <div class="p-1.5 bg-orange-100 rounded text-orange-600">
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="w-4 h-4"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            ><polyline points="16 3 21 3 21 8"></polyline><line
                                x1="4"
                                y1="20"
                                x2="21"
                                y2="3"
                            ></line><polyline points="21 16 21 21 16 21"
                            ></polyline><line x1="15" y1="15" x2="21" y2="21"
                            ></line><line x1="4" y1="4" x2="9" y2="9"
                            ></line></svg
                        >
                    </div>
                    <span class="text-sm font-medium">If / Else</span>
                </div>

                <div
                    class="bg-white p-3 rounded cursor-move hover:bg-gray-50 hover:ring-1 hover:ring-indigo-500 transition-all border border-gray-200 flex items-center gap-3 shadow-sm"
                    draggable={true}
                    role="button"
                    tabindex="0"
                    aria-label="Delay Logic"
                    ondragstart={(e) =>
                        onDragStart(e, "ACTION", {
                            label: "Delay",
                            action: "Delay",
                            delay_minutes: 60,
                        })}
                    onclick={() =>
                        handleNodeClick("ACTION", {
                            label: "Delay",
                            action: "Delay",
                            delay_minutes: 60,
                        })}
                    onkeydown={(e) =>
                        e.key === "Enter" &&
                        handleNodeClick("ACTION", {
                            label: "Delay",
                            action: "Delay",
                            delay_minutes: 60,
                        })}
                >
                    <div class="p-1.5 bg-amber-100 rounded text-amber-600">
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="w-4 h-4"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            ><circle cx="12" cy="12" r="10"></circle><polyline
                                points="12 6 12 12 16 14"
                            ></polyline></svg
                        >
                    </div>
                    <span class="text-sm font-medium">Delay</span>
                </div>
            </div>
        </div>
    </div>
</div>
