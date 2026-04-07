<script>
  import { Router } from '@keenmate/svelte-spa-router';
  import Init from './routes/Init.svelte';
  import Welcome from './routes/Welcome.svelte';
  import Config from './routes/Config.svelte';
  import Page404 from './routes/Page404.svelte';
  import TemplateFilter from './routes/TemplateFilter.svelte';
  import DataSetFilter from './routes/DataSetFilter.svelte';
  import ExportStatus from './routes/ExportStatus.svelte';
  import Update from './routes/Update.svelte';
  import StatusBar from './components/StatusBar.svelte';
  import { Toaster } from '$lib/components/ui/sonner';
</script>

<div class="flex h-screen flex-col">
  <div class="flex-1 overflow-y-auto pb-10">
    <Router
      routes={{
        '/': Init,
        '/update': Update,
        '/welcome': Welcome,
        '/config': Config,
        '/config/templates': TemplateFilter,
        '/config/datasets': DataSetFilter,
        '/export/status': ExportStatus,
        '*': Page404,
      }}
    />
  </div>
  <StatusBar />
</div>
<Toaster position="top-center" richColors />
