<script lang="ts">
  import dayjs from 'dayjs';
  import { shadowConfig, templateCache } from '../lib/store';
  import { GetTemplates } from '../../wailsjs/go/main/App.js';
  import { push } from '@keenmate/svelte-spa-router';
  import { trim } from '../lib/utils';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import { Checkbox } from '$lib/components/ui/checkbox';
  import * as Dialog from '$lib/components/ui/dialog';
  import * as Table from '$lib/components/ui/table';
  import { SortButton } from '$lib/components/ui/data-table';
  import {
    type ColumnDef,
    type SortingState,
    getCoreRowModel,
    getSortedRowModel,
    getFilteredRowModel,
    createTable,
    FlexRender,
  } from '@tanstack/svelte-table';
  import { renderComponent } from '@tanstack/svelte-table';

  import SearchIcon from '@lucide/svelte/icons/search';
  import XCircleIcon from '@lucide/svelte/icons/x-circle';

  interface Template {
    id: string;
    name: string;
    modified_at: string;
    modified_at_raw: string;
  }

  let searchFilter = $state('');
  let templatesLoaded = $state(false);
  let tableData = $state<Template[]>([]);
  let selectedIds = $state<Set<string>>(new Set());

  if (Array.isArray($templateCache)) {
    if ($templateCache.length === 0) {
      GetTemplates().then((result: any[]) => {
        let niceFormat = result
          .map((elem) => ({
            id: elem.id,
            name: elem.name.length > 90 ? `${elem.name.substring(0, 90)}…` : elem.name,
            modified_at: dayjs(elem.modified_at).format('DD-MMM-YYYY'),
            modified_at_raw: elem.modified_at,
          }))
          .slice(0, 3000);
        templatesLoaded = true;
        templateCache.set(niceFormat);
        tableData = niceFormat;
        initSelection();
      });
    } else {
      templatesLoaded = true;
      tableData = $templateCache as Template[];
      initSelection();
    }
  }

  function initSelection() {
    const existing = ($shadowConfig as any)['Export']['TemplateIds'] as string[];
    if (existing.length === 0) {
      selectedIds = new Set(tableData.map((t) => t.id));
      ($shadowConfig as any)['Export']['TemplateIds'] = tableData.map((t) => t.id);
    } else {
      selectedIds = new Set(existing);
    }
  }

  let allSelected = $derived(tableData.length > 0 && selectedIds.size === tableData.length);

  function toggleAll() {
    if (allSelected) {
      selectedIds = new Set();
    } else {
      selectedIds = new Set(tableData.map((t) => t.id));
    }
  }

  function toggleRow(id: string) {
    const next = new Set(selectedIds);
    if (next.has(id)) next.delete(id);
    else next.add(id);
    selectedIds = next;
  }

  function handleDone() {
    const selected = [...selectedIds];
    ($shadowConfig as any)['Export']['TemplateIds'] =
      selected.length === tableData.length ? [] : selected;
    push('/config');
  }

  // TanStack Table setup
  let sorting = $state<SortingState>([]);
  let globalFilter = $state('');

  // Keep globalFilter in sync with searchFilter
  $effect(() => {
    globalFilter = searchFilter.length >= 2 ? searchFilter : '';
  });

  const columns: ColumnDef<Template>[] = [
    {
      id: 'select',
      header: () => '',
      cell: ({ row }) => row.original.id,
      enableSorting: false,
      enableGlobalFilter: false,
    },
    {
      accessorKey: 'name',
      header: 'Template',
      cell: ({ row }) => trim(row.original.name, 90),
      enableSorting: false,
    },
    {
      accessorKey: 'modified_at_raw',
      header: ({ column }) =>
        renderComponent(SortButton, {
          column: column,
          label: 'Last modified',
        }),
      cell: ({ row }) => row.original.modified_at,
    },
  ];

  const table = createTable<Template>({
    get data() {
      return tableData;
    },
    columns,
    state: {
      get sorting() {
        return sorting;
      },
      get globalFilter() {
        return globalFilter;
      },
    },
    onSortingChange: (updater) => {
      sorting = typeof updater === 'function' ? updater(sorting) : updater;
    },
    onGlobalFilterChange: (updater) => {
      globalFilter = typeof updater === 'function' ? updater(globalFilter) : updater;
    },
    globalFilterFn: (row, _columnId, filterValue) => {
      return row.original.name.toLowerCase().includes(filterValue.toLowerCase());
    },
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
  });

  let visibleRows = $derived(table.getRowModel().rows);
  let showEmptyFilter = $derived(searchFilter.length >= 2 && visibleRows.length === 0);
</script>

{#if !templatesLoaded}
  <Dialog.Root open={true}>
    <Dialog.Content class="sm:max-w-md">
      <Dialog.Header>
        <Dialog.Title>Loading templates</Dialog.Title>
        <Dialog.Description>Please wait while we process your request...</Dialog.Description>
      </Dialog.Header>
      <div class="flex justify-center py-4">
        <img src="../images/loading.gif" alt="loading" class="size-12" />
      </div>
    </Dialog.Content>
  </Dialog.Root>
{/if}

<div class="px-8 pt-6">
  <div class="flex items-center justify-between">
    <h1 class="text-xl font-semibold">Template selection</h1>
    <Button variant="outline" onclick={handleDone}>Done</Button>
  </div>

  <div class="mt-5 flex items-center justify-between">
    <h2 class="text-base font-medium">Select templates</h2>
    <div class="relative w-[260px]">
      <SearchIcon
        class="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground"
      />
      <Input class="pl-9 pr-9" placeholder="Search" bind:value={searchFilter} />
      {#if searchFilter.length > 0}
        <button
          class="absolute top-1/2 right-3 -translate-y-1/2 cursor-pointer text-muted-foreground hover:text-foreground"
          onclick={() => (searchFilter = '')}
        >
          <XCircleIcon class="size-4" />
        </button>
      {/if}
    </div>
  </div>

  {#if showEmptyFilter}
    <div class="mt-8 flex h-[50vh] flex-col items-center justify-center text-center">
      <img src="../images/empty_page.svg" alt="empty page" class="mb-8" />
      <p class="text-sm">
        Your search - <span class="font-semibold"
          >{searchFilter.length > 15 ? searchFilter.substring(0, 15) + '...' : searchFilter}</span
        > - did not match any template names.
      </p>
      <p class="mt-2 text-sm text-muted-foreground">
        Make sure all the words are spelled correctly, or try different keywords.
      </p>
    </div>
  {:else}
    <div class="mt-4 overflow-hidden rounded-lg border border-border">
      <Table.Root>
        <Table.Header>
          {#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
            <Table.Row class="bg-muted hover:bg-muted">
              {#each headerGroup.headers as header (header.id)}
                <Table.Head>
                  {#if header.id === 'select'}
                    <Checkbox checked={allSelected} onCheckedChange={toggleAll} />
                  {:else if !header.isPlaceholder}
                    <FlexRender
                      content={header.column.columnDef.header}
                      context={header.getContext()}
                    />
                  {/if}
                </Table.Head>
              {/each}
            </Table.Row>
          {/each}
        </Table.Header>
        <Table.Body class="max-h-[calc(100vh-280px)] overflow-y-auto">
          {#each visibleRows as row (row.id)}
            <Table.Row>
              {#each row.getVisibleCells() as cell (cell.id)}
                <Table.Cell>
                  {#if cell.column.id === 'select'}
                    <Checkbox
                      checked={selectedIds.has(row.original.id)}
                      onCheckedChange={() => toggleRow(row.original.id)}
                    />
                  {:else if cell.column.id === 'name'}
                    <div class="flex items-center">
                      <img class="mr-2 size-4" src="../images/template-icon.svg" alt="template" />
                      <span class="text-sm">{trim(row.original.name, 90)}</span>
                    </div>
                  {:else}
                    <span class="text-sm text-muted-foreground">{row.original.modified_at}</span>
                  {/if}
                </Table.Cell>
              {/each}
            </Table.Row>
          {:else}
            <Table.Row>
              <Table.Cell colspan={3} class="h-24 text-center">No results.</Table.Cell>
            </Table.Row>
          {/each}
        </Table.Body>
      </Table.Root>
    </div>
  {/if}
</div>
