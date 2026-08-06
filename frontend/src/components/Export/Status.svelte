<script lang="ts">
  import { onMount } from 'svelte';
  import { EventsOn } from '../../../wailsjs/runtime/runtime.js';
  import { Badge } from '$lib/components/ui/badge';
  import * as Table from '$lib/components/ui/table';

  const statusQueued = 'Queued';
  const statusFailed = 'Failed';
  const statusComplete = 'Complete';
  const statusDownloading = 'Downloading';
  const statusExporting = 'Exporting';

  let {
    name = '',
    cancelled = false,
  }: {
    name?: string;
    cancelled?: boolean;
  } = $props();

  let status = $state(statusQueued);
  let counter = $state(-1);
  let counterDecremental = $state(true);
  let pillType = $state<'success' | 'info' | 'neutral' | 'error' | 'cancelled'>('neutral');

  function formatExportItemName(str: string): string {
    return str
      .toLowerCase()
      .replace(/_/g, ' ')
      .replace(/(^|\s)\S/g, (L) => L.toUpperCase());
  }

  const badgeStyles = {
    success: 'bg-emerald-50 text-emerald-800 hover:bg-emerald-50',
    info: 'bg-sky-50 text-sky-700 hover:bg-sky-50',
    neutral: 'bg-muted text-muted-foreground hover:bg-muted',
    error: 'bg-destructive text-white hover:bg-destructive',
    cancelled: 'bg-red-50 text-red-700 hover:bg-red-50',
  } as const;

  onMount(() => {
    EventsOn('update-' + name, (newValue: Record<string, unknown>) => {
      counterDecremental = newValue['counter_decremental'] as boolean;
      if (newValue['started'] === true && newValue['finished'] === false) {
        switch (newValue['stage']) {
          case 'API_DOWNLOAD':
            status = statusDownloading;
            break;
          case 'CSV_EXPORT':
            status = statusExporting;
            break;
        }
        counter = newValue['counter'] as number;
      }

      if (newValue['started'] === true && newValue['finished'] === true) {
        status = newValue['has_error'] === false ? statusComplete : statusFailed;
        counter = 0;
      }

      switch (status) {
        case 'Complete':
          pillType = 'success';
          break;
        case 'Downloading':
        case 'Exporting':
          pillType = 'info';
          break;
        case 'Queued':
          pillType = 'neutral';
          break;
        case 'Failed':
          pillType = 'error';
          break;
      }
    });
  });
</script>

<Table.Cell class="w-[40%]">{formatExportItemName(name)}</Table.Cell>
<Table.Cell class="w-[20%]">
  {#if cancelled}
    <Badge variant="secondary" class={badgeStyles.cancelled}>Cancelled</Badge>
  {:else}
    <Badge variant="secondary" class={badgeStyles[pillType]}>{status}</Badge>
  {/if}
</Table.Cell>
<Table.Cell>
  {#if !cancelled}
    {counter === -1 || counter === 0
      ? ''
      : counter + ' ' + (counterDecremental ? 'remaining' : 'exported')}
  {/if}
</Table.Cell>
