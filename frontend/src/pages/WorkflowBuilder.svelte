<script lang="ts">
    import { onMount } from "svelte";
    import WorkflowBuilder from "../lib/components/Workflow/WorkflowBuilder.svelte";
    import { router } from "../lib/router.svelte.js";

    let workflowId = $state(0);

    onMount(() => {
        const match = window.location.pathname.match(/\/workflows\/(\d+|new)/);
        if (match && match[1] !== "new") {
            workflowId = parseInt(match[1]);
        }
    });

    function handleBack() {
        router.navigate("/workflows");
    }

    function handleWorkflowCreated(id: number) {
        window.history.replaceState({}, "", `/workflows/${id}`);
    }
</script>

<WorkflowBuilder
    {workflowId}
    onBack={handleBack}
    onWorkflowCreated={handleWorkflowCreated}
/>
