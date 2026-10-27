<script lang="ts">
  import type { Column } from '@tanstack/svelte-table';
  import { Button } from '$lib/components/ui/button';
  import ArrowUpDownIcon from '@lucide/svelte/icons/arrow-up-down';
  import ArrowUpIcon from '@lucide/svelte/icons/arrow-up';
  import ArrowDownIcon from '@lucide/svelte/icons/arrow-down';

  let {
    column,
    label,
  }: {
    column: Column<any, unknown>;
    label: string;
  } = $props();

  let sorted = $derived(column.getIsSorted());
</script>

<Button variant="ghost" class="-ml-3 h-8" onclick={() => column.toggleSorting(sorted === 'asc')}>
  {label}
  {#if sorted === 'asc'}
    <ArrowUpIcon class="ml-2 size-4" />
  {:else if sorted === 'desc'}
    <ArrowDownIcon class="ml-2 size-4" />
  {:else}
    <ArrowUpDownIcon class="ml-2 size-4" />
  {/if}
</Button>
