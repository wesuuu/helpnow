<script lang="ts">
    import Route from "./lib/Route.svelte";
    import Login from "./pages/Login.svelte";
    import Register from "./pages/Register.svelte";
    import Dashboard from "./pages/Dashboard.svelte";
    import Campaigns from "./pages/Campaigns.svelte";
    import Sites from "./pages/Sites.svelte";
    import SiteDetail from "./pages/SiteDetail.svelte";
    import Workflows from "./pages/Workflows.svelte";
    import WorkflowBuilder from "./pages/WorkflowBuilder.svelte";
    import Sources from "./pages/Sources.svelte";
    import SourceDetail from "./pages/SourceDetail.svelte";
    import Audiences from "./pages/Audiences.svelte";
    import OrgSettings from "./pages/OrgSettings.svelte";
    import TeamSettings from "./pages/TeamSettings.svelte";
    import UserSettings from "./pages/UserSettings.svelte";
    import EmailSettings from "./pages/EmailSettings.svelte";
    import ContentTemplates from "./pages/ContentTemplates.svelte";
    import Agents from "./pages/Agents.svelte";
    import PersonDetail from "./pages/PersonDetail.svelte";
    import Layout from "./lib/Layout.svelte";
    import { auth } from "./stores/auth.js";
    import { router } from "./lib/router.svelte.js";
    import ToastContainer from "./lib/components/ToastContainer.svelte";

    // Redirect to login if not authenticated and trying to access protected route
    $effect(() => {
        if (
            !$auth.isAuthenticated &&
            router.path !== "/login" &&
            router.path !== "/register"
        ) {
            router.navigate("/login");
        }
    });
</script>

<div>
    <ToastContainer />
    <Route path="/login" component={Login} />
    <Route path="/register" component={Register} />

    {#if $auth.isAuthenticated}
        <Layout>
            <Route path="/dashboard" component={Dashboard} />
            <Route path="/campaigns" component={Campaigns} />
            <Route path="/sites" component={Sites} />
            <Route path="/sites/:id" component={SiteDetail} />
            <Route path="/workflows" component={Workflows} />
            <Route path="/workflows/:id" component={WorkflowBuilder} />
            <Route path="/sources" component={Sources} />
            <Route path="/sources/:id" component={SourceDetail} />
            <Route path="/audiences" component={Audiences} />
            <Route path="/people/:id" component={PersonDetail} />
            <Route path="/organization-settings" component={OrgSettings} />
            <Route path="/team-settings" component={TeamSettings} />
            <Route path="/user-settings" component={UserSettings} />
            <Route path="/email-settings" component={EmailSettings} />
            <Route path="/templates" component={ContentTemplates} />
            <Route path="/agents" component={Agents} />
            <Route path="/" component={Dashboard} />
        </Layout>
    {/if}
</div>
