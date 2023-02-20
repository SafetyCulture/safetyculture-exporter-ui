<script>
    import Router from 'svelte-spa-router';
    import Init from "./routes/Init.svelte"
    import Welcome from "./routes/Welcome.svelte"
    import Config from "./routes/Config.svelte";
    import Page404 from "./routes/Page404.svelte";
    import TemplateFilter from "./routes/TemplateFilter.svelte";
    import DataSetFilter from "./routes/DataSetFilter.svelte";
    import ExportStatus from "./routes/ExportStatus.svelte";
</script>

<div class="container">
<Router routes={{
	'/': Init,
	'/welcome': Welcome,
	'/config': Config,
	'/config/templates': TemplateFilter,
	'/config/datasets': DataSetFilter,
	'/export/status': ExportStatus,
	'*': Page404
}} />
</div>

<style>
    .container {
        height: 100vh;
        display: flex;
        flex-direction: column;
        justify-content: space-between;
    }
</style>
