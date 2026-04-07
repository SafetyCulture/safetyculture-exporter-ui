<script lang="ts">
  import dayjs from 'dayjs';
  import { shadowConfig, templateCache } from '../lib/store';
  import { GetTemplates } from '../../wailsjs/go/main/App.js';
  import { push } from '@keenmate/svelte-spa-router';
  import { trim } from '../lib/utils';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';

  import * as Dialog from '$lib/components/ui/dialog';

  import SearchIcon from '@lucide/svelte/icons/search';
  import XCircleIcon from '@lucide/svelte/icons/x-circle';

  let searchFilter = $state('');
  let isChecked = $state(false);
  let templatesLoaded = $state(false);

  if (Array.isArray($templateCache)) {
    if ($templateCache.length === 0) {
      GetTemplates().then((result: any[]) => {
        let niceFormat = result
          .map((elem) => ({
            id: elem.id,
            name: elem.name.length > 90 ? `${elem.name.substring(0, 90)}…` : elem.name,
            modified_at: dayjs(elem.modified_at).format('DD-MMM-YYYY'),
          }))
          .slice(0, 3000);
        templatesLoaded = true;
        templateCache.set(niceFormat);
        checkAllSelected();
      });
    } else {
      templatesLoaded = true;
    }
  }

  let filteredTemplates = $derived(
    searchFilter.length >= 2
      ? ($templateCache as any[]).filter((v: any) =>
          v.name.toLowerCase().includes(searchFilter.toLowerCase()),
        )
      : ($templateCache as any[]),
  );

  let showEmptyFilter = $derived(searchFilter.length >= 2 && filteredTemplates.length === 0);

  function checkAllSelected() {
    if (($shadowConfig as any)['Export']['TemplateIds'].length === 0) {
      ($shadowConfig as any)['Export']['TemplateIds'] = ($templateCache as any[]).map(
        (e: any) => e.id,
      );
      isChecked = true;
    }
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
    if (showEmptyFilter) {
      push('/config');
      return;
    }

    const checkboxes = document.querySelectorAll(
      '.table-body input[type="checkbox"]',
    ) as NodeListOf<HTMLInputElement>;
    let selectedTemplates: string[] = [];

    for (const checkbox of checkboxes) {
      if (checkbox.checked) selectedTemplates.push((checkbox as any).__value);
    }

    ($shadowConfig as any)['Export']['TemplateIds'] =
      ($templateCache as any[]).length === selectedTemplates.length ? [] : selectedTemplates;

    push('/config');
  }

  checkAllSelected();
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
      <div class="flex h-10 items-center justify-between bg-muted px-3">
        <div class="flex items-center">
          <input
            type="checkbox"
            class="size-4 accent-primary"
            onclick={toggleBodyCheckboxes}
            bind:checked={isChecked}
          />
          <span class="ml-4 text-sm font-medium text-muted-foreground">Template</span>
        </div>
        <span class="text-sm font-medium text-muted-foreground">Last modified</span>
      </div>
      <div class="table-body max-h-[calc(100vh-280px)] overflow-y-auto">
        {#each filteredTemplates as { id, name, modified_at } (id)}
          <div class="flex h-11 items-center justify-between border-t border-border px-3">
            <div class="flex items-center">
              <input
                type="checkbox"
                class="size-4 accent-primary"
                onclick={toggleHeaderCheckbox}
                bind:group={($shadowConfig as any)['Export']['TemplateIds']}
                value={id}
              />
              <img class="ml-4 size-4" src="../images/template-icon.svg" alt="template" />
              <span class="ml-2 text-sm">{trim(name as string, 90)}</span>
            </div>
            <span class="text-sm text-muted-foreground">{modified_at}</span>
          </div>
        {/each}
      </div>
    </div>
  {/if}
</div>
