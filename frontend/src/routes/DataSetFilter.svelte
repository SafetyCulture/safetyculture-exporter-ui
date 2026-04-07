<script lang="ts">
  import { shadowConfig } from '../lib/store';
  import { push } from '@keenmate/svelte-spa-router';
  import { Button } from '$lib/components/ui/button';

  interface DataItem {
    id: string;
    name: string;
  }
  interface DataRow {
    left: DataItem | null;
    right: DataItem | null;
  }

  const data: DataRow[] = [
    {
      left: { id: 'inspections', name: 'Inspections' },
      right: { id: 'inspection_items', name: 'Inspection Items' },
    },
    {
      left: { id: 'templates', name: 'Templates' },
      right: { id: 'template_permissions', name: 'Template Permissions' },
    },
    { left: { id: 'sites', name: 'Sites' }, right: { id: 'site_members', name: 'Site Members' } },
    { left: { id: 'groups', name: 'Groups' }, right: { id: 'group_users', name: 'Group Users' } },
    { left: { id: 'users', name: 'Users' }, right: { id: 'schedules', name: 'Schedules' } },
    {
      left: { id: 'schedule_assignees', name: 'Schedule Assignees' },
      right: { id: 'schedule_occurrences', name: 'Schedule Occurrences' },
    },
    { left: { id: 'actions', name: 'Actions' }, right: { id: 'issues', name: 'Issues' } },
    {
      left: { id: 'action_timeline_items', name: 'Action Timeline Items' },
      right: { id: 'issue_timeline_items', name: 'Issue Timeline Items' },
    },
    {
      left: { id: 'action_assignees', name: 'Action Assignees' },
      right: { id: 'issue_assignees', name: 'Issue Assignees' },
    },
    {
      left: { id: 'assets', name: 'Assets' },
      right: { id: 'training_course_progresses', name: 'Training course completions' },
    },
  ];

  const allItems: DataItem[] = data.flatMap((row) =>
    [row.left, row.right].filter((x): x is DataItem => x !== null),
  );

  let isChecked = $state(false);
  if (($shadowConfig as any)['Export']['Tables'].length === 0) {
    ($shadowConfig as any)['Export']['Tables'] = allItems.map((e) => e.id);
    isChecked = true;
  }

  function toggleHeaderCheckbox() {
    if (isChecked) isChecked = false;
  }

  function toggleBodyCheckboxes() {
    const checkboxes = document.querySelectorAll(
      '.table-body input[type="checkbox"]',
    ) as NodeListOf<HTMLInputElement>;
    for (const checkbox of checkboxes) {
      checkbox.checked = !isChecked;
    }
  }

  function handleDone() {
    const checkboxes = document.querySelectorAll(
      '.table-body input[type="checkbox"]',
    ) as NodeListOf<HTMLInputElement>;
    let selectedTables: string[] = [];
    for (const checkbox of checkboxes) {
      if (checkbox.checked) selectedTables.push((checkbox as any).__value);
    }

    ($shadowConfig as any)['Export']['Tables'] =
      selectedTables.length === allItems.length ? [] : selectedTables;

    push('/config');
  }
</script>

<div class="px-8 pt-6">
  <section class="flex items-center justify-between">
    <h1 class="text-xl font-semibold">Data set selection</h1>
    <Button variant="outline" onclick={handleDone}>Done</Button>
  </section>

  <div class="mt-5 overflow-hidden rounded-lg border border-border">
    <div class="flex h-10 items-center bg-muted px-3">
      <input
        type="checkbox"
        class="size-4 accent-primary"
        onclick={toggleBodyCheckboxes}
        bind:checked={isChecked}
      />
      <span class="ml-4 text-sm font-medium text-muted-foreground">Data set</span>
    </div>
    <div class="table-body">
      {#each data as { left, right }, i (i)}
        <div class="flex h-11 border-t border-border">
          {#if left}
            <div class="flex w-1/2 items-center px-3">
              <input
                type="checkbox"
                class="size-4 accent-primary"
                onclick={toggleHeaderCheckbox}
                bind:group={($shadowConfig as any)['Export']['Tables']}
                value={left.id}
              />
              <img class="ml-4 size-4" src="../images/template-icon.svg" alt="dataset" />
              <span class="ml-2 text-sm">{left.name}</span>
            </div>
          {/if}
          {#if right}
            <div class="flex w-1/2 items-center px-3">
              <input
                type="checkbox"
                class="size-4 accent-primary"
                onclick={toggleHeaderCheckbox}
                bind:group={($shadowConfig as any)['Export']['Tables']}
                value={right.id}
              />
              <img class="ml-4 size-4" src="../images/template-icon.svg" alt="dataset" />
              <span class="ml-2 text-sm">{right.name}</span>
            </div>
          {/if}
        </div>
      {/each}
    </div>
  </div>
</div>
