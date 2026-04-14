<script lang="ts">
  import {
    type ColumnDef,
    type SortingState,
    getCoreRowModel,
    getSortedRowModel,
    createTable,
    FlexRender,
  } from '@tanstack/svelte-table';
  import { renderComponent } from '@tanstack/svelte-table';
  import * as Table from '$lib/components/ui/table';
  import SortButton from './sort-button.svelte';

  interface Item {
    name: string;
    date: string;
  }

  let { empty = false }: { empty?: boolean } = $props();

  const allData: Item[] = [
    { name: 'Alpha', date: '2025-06-15' },
    { name: 'Bravo', date: '2025-12-01' },
    { name: 'Charlie', date: '2025-01-10' },
  ];

  let sampleData = $derived(empty ? [] : allData);

  let sorting = $state<SortingState>([]);

  const columns: ColumnDef<Item>[] = [
    {
      accessorKey: 'name',
      header: 'Name',
      enableSorting: false,
    },
    {
      accessorKey: 'date',
      header: ({ column }) =>
        renderComponent(SortButton, {
          column: column,
          label: 'Date',
        }),
    },
  ];

  const table = createTable<Item>({
    get data() {
      return sampleData;
    },
    columns,
    state: {
      get sorting() {
        return sorting;
      },
    },
    onSortingChange: (updater) => {
      sorting = typeof updater === 'function' ? updater(sorting) : updater;
    },
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
  });
</script>

<Table.Root>
  <Table.Header>
    {#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
      <Table.Row>
        {#each headerGroup.headers as header (header.id)}
          <Table.Head>
            {#if !header.isPlaceholder}
              <FlexRender content={header.column.columnDef.header} context={header.getContext()} />
            {/if}
          </Table.Head>
        {/each}
      </Table.Row>
    {/each}
  </Table.Header>
  <Table.Body>
    {#each table.getRowModel().rows as row (row.id)}
      <Table.Row>
        {#each row.getVisibleCells() as cell (cell.id)}
          <Table.Cell>
            <FlexRender content={cell.column.columnDef.cell} context={cell.getContext()} />
          </Table.Cell>
        {/each}
      </Table.Row>
    {:else}
      <Table.Row>
        <Table.Cell colspan={columns.length} class="h-24 text-center">No results.</Table.Cell>
      </Table.Row>
    {/each}
  </Table.Body>
</Table.Root>
