<script>
    import Router from 'svelte-spa-router';
    import Init from "./routes/Init.svelte"
    import Welcome from "./routes/Welcome.svelte"
    import Config from "./routes/Config.svelte";
    import Page404 from "./routes/Page404.svelte";
    import TemplateFilter from "./routes/TemplateFilter.svelte";
    import FeedFilter from "./routes/FeedFilter.svelte";
    import StatusBar from "./components/StatusBar.svelte";
    import ExportStatus from "./routes/ExportStatus.svelte";
    import DebugMem from "./routes/DebugMem.svelte";
</script>

<div class="container">
<Router routes={{
	'/': Init,
	'/welcome': Welcome,
	'/config': Config,
	'/config/templates': TemplateFilter,
	'/config/tables': FeedFilter,
	'/exportStatus': ExportStatus,
	'/debug': DebugMem,
	'*': Page404
}} />

<StatusBar/>
</div>

<style>
    .container {
        height: 100vh;
        display: flex;
        flex-direction: column;
        justify-content: space-between;
    }
</style>
