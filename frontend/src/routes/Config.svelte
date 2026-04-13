<script lang="ts">
  import dayjs from 'dayjs';
  import utc from 'dayjs/plugin/utc';
  import timezone from 'dayjs/plugin/timezone';
  import {
    CalendarDate,
    today as calToday,
    getLocalTimeZone,
    parseDate,
  } from '@internationalized/date';
  import { toast } from 'svelte-sonner';
  import {
    CheckDBConnection,
    ExportCSV,
    ExportReports,
    ExportSQL,
    ExportSQLite,
    ReadBuild,
    SaveSettings,
    SelectDirectory,
    ValidateExportDirectory,
  } from '../../wailsjs/go/main/App.js';
  import { exportConfig, shadowConfig } from '../lib/store';
  import { allTables } from '../lib/utils';
  import { Quit } from '../../wailsjs/runtime/runtime.js';
  import { push } from '@keenmate/svelte-spa-router';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Checkbox } from '$lib/components/ui/checkbox';
  import { Calendar } from '$lib/components/ui/calendar';
  import * as Select from '$lib/components/ui/select';
  import * as Dialog from '$lib/components/ui/dialog';
  import * as Popover from '$lib/components/ui/popover';
  import { Separator } from '$lib/components/ui/separator';

  import ChevronRightIcon from '@lucide/svelte/icons/chevron-right';
  import CalendarIcon from '@lucide/svelte/icons/calendar';
  import FolderIcon from '@lucide/svelte/icons/folder';
  import ArrowLeftIcon from '@lucide/svelte/icons/arrow-left';

  let build = $state('');
  ReadBuild().then((it: string) => {
    build = it;
  });

  const statusItems = [
    { value: 'true', label: 'Completed only' },
    { value: 'false', label: 'Incomplete only' },
    { value: 'both', label: 'Both - completed and incomplete' },
  ];
  let selectedStatus = $state(
    ($shadowConfig as any)['Export']['Inspection']['Completed'] as string,
  );

  const archivedItems = [
    { value: 'true', label: 'Archived only' },
    { value: 'false', label: 'Active only' },
    { value: 'both', label: 'Both - active and archived' },
  ];
  let selectedArchived = $state(
    ($shadowConfig as any)['Export']['Inspection']['Archived'] as string,
  );

  const dataExportFormatItems = [
    { value: 'csv', label: 'CSV' },
    { value: 'mysql', label: 'MySQL' },
    { value: 'postgres', label: 'Postgres' },
    { value: 'sqlserver', label: 'SQL Server' },
    { value: 'sqlite', label: 'SQLite' },
    { value: 'reports', label: 'Reports' },
  ];

  const connectionStrings: Record<string, string> = {
    postgres: 'postgresql://${dbUser}:${dbPassword}@${dbHost}:${dbPort}/${dbName}',
    sqlserver: 'sqlserver://${dbUser}:${dbPassword}@${dbHost}:${dbPort}?database=${dbName}',
    mysql:
      '${dbUser}:${dbPassword}@tcp(${dbHost}:${dbPort})/${dbName}?charset=utf8mb4&parseTime=True&loc=Local',
  };
  const dialects: Record<string, string> = {
    postgres: 'postgres',
    sqlserver: 'sqlserver',
    mysql: 'mysql',
  };

  let dbHost = $state(''),
    dbHostShowError = $state(false);
  let dbPort = $state(''),
    dbPortPlaceholder = $state(
      'e.g. ' + getDefaultSQLPort(($shadowConfig as any)['Db']['Dialect']),
    ),
    dbPortShowError = $state(false);
  let dbUser = $state(''),
    dbUserShowError = $state(false);
  let dbPassword = $state(''),
    dbPasswordShowError = $state(false);
  let dbName = $state(''),
    dbNameShowError = $state(false);
  let formError = $state(false);
  let dbError = $state(false);
  let showBanner = $state(false);
  let exportLocationError = $state(false);

  let selectedExportFormat = $state(($shadowConfig as any)['Session']['ExportType'] as string);

  const reportFormatItems = [
    { value: 'PDF', label: 'PDF' },
    { value: 'WORD', label: 'Word' },
    { value: 'BOTH', label: 'Both PDF and Word' },
  ];
  let selectedReportFormat = $state(readReportFormat());

  let selectedTimeZone = $state(($shadowConfig as any)['Export']['TimeZone'] as string);

  function readReportFormat(): string {
    const fmt = ($shadowConfig as any)['Report']['Format'];
    if (Array.isArray(fmt) && fmt.includes('PDF') && fmt.includes('WORD')) return 'BOTH';
    if (Array.isArray(fmt) && fmt.includes('WORD')) return 'WORD';
    return 'PDF';
  }

  function prepareReportFormatForSave(): string[] {
    switch (selectedReportFormat) {
      case 'PDF':
        return ['PDF'];
      case 'WORD':
        return ['WORD'];
      case 'BOTH':
        return ['PDF', 'WORD'];
      default:
        return ['PDF'];
    }
  }

  // Date picker
  dayjs.extend(utc);
  dayjs.extend(timezone);

  const minCalDate = calToday(getLocalTimeZone()).subtract({ years: 10 });
  const maxCalDate = calToday(getLocalTimeZone());

  function initCalendarDate(): CalendarDate {
    const raw = ($shadowConfig as any)['Export']['ModifiedAfter'] as string;
    if (raw === '' || raw === '0001-01-01T00:00:00Z') return minCalDate;
    const iso = dayjs(raw).tz(selectedTimeZone).format('YYYY-MM-DD');
    try {
      return parseDate(iso);
    } catch {
      return minCalDate;
    }
  }

  let calendarValue = $state<CalendarDate>(initCalendarDate());
  let datePickerOpen = $state(false);
  let dateInputValue = $state(formatCalDate(initCalendarDate()));
  let dateInputError = $state(false);

  function formatCalDate(d: CalendarDate): string {
    return `${String(d.day).padStart(2, '0')}/${String(d.month).padStart(2, '0')}/${d.year}`;
  }

  function tryParseDate(input: string): CalendarDate | null {
    const match = input.match(/^(\d{1,2})\/(\d{1,2})\/(\d{4})$/);
    if (!match) return null;
    const [, dd, mm, yyyy] = match;
    try {
      return new CalendarDate(+yyyy, +mm, +dd);
    } catch {
      return null;
    }
  }

  function handleDateInput() {
    const parsed = tryParseDate(dateInputValue);
    if (!parsed) {
      dateInputError = true;
      return;
    }
    dateInputError = false;

    if (parsed.compare(minCalDate) < 0) {
      toast.error(`Date cannot be earlier than ${formatCalDate(minCalDate)}`);
      calendarValue = minCalDate;
      dateInputValue = formatCalDate(minCalDate);
    } else if (parsed.compare(maxCalDate) > 0) {
      toast.error(`Date cannot be later than ${formatCalDate(maxCalDate)}`);
      calendarValue = maxCalDate;
      dateInputValue = formatCalDate(maxCalDate);
    } else {
      calendarValue = parsed;
    }
  }

  function handleDatePickerClose(open: boolean) {
    if (!open) {
      const parsed = tryParseDate(dateInputValue);
      if (!parsed || parsed.compare(minCalDate) < 0 || parsed.compare(maxCalDate) > 0) {
        toast.error('Invalid date (DD/MM/YYYY), reverting to previous value');
        dateInputValue = formatCalDate(calendarValue);
      }
      dateInputError = false;
    }
  }

  function handleCalendarSelect(value: CalendarDate | undefined) {
    if (value) {
      calendarValue = value;
      dateInputValue = formatCalDate(value);
      datePickerOpen = false;
    }
  }

  // For saving — convert CalendarDate back to YYYY-MM-DD string
  let dateValue = $derived(
    `${calendarValue.year}-${String(calendarValue.month).padStart(2, '0')}-${String(calendarValue.day).padStart(2, '0')}`,
  );

  function generateTemplateName(): string {
    const num = ($shadowConfig as any)['Export']['TemplateIds'].length;
    if (num === 0) return 'All templates selected';
    if (num === 1) return '1 template selected';
    return num + ' templates selected';
  }

  function generateDataSetName(): string {
    const num = ($shadowConfig as any)['Export']['Tables'].length;
    if (num === 0) return 'All data sets selected';
    if (num === 1) return '1 data set selected';
    return num + ' data sets selected';
  }

  function handleExportFormatUpdate(value: string) {
    selectedExportFormat = value;
    if (['mysql', 'postgres', 'sqlserver'].includes(selectedExportFormat)) {
      dbPortPlaceholder = 'e.g. ' + getDefaultSQLPort(selectedExportFormat);
      dbPort = getDefaultSQLPort(selectedExportFormat);
    }
  }

  function setConnString() {
    if (selectedExportFormat !== '') {
      const connectionString = connectionStrings[selectedExportFormat];
      ($shadowConfig as any)['Db']['ConnectionString'] = connectionString
        .replace(/\${dbUser}/g, dbUser)
        .replace(/\${dbPassword}/g, dbPassword)
        .replace(/\${dbHost}/g, dbHost)
        .replace(/\${dbPort}/g, dbPort)
        .replace(/\${dbName}/g, dbName);
      ($shadowConfig as any)['Db']['Dialect'] = dialects[selectedExportFormat];
    }
  }

  function parseDbConnectionString() {
    const url = ($shadowConfig as any)['Db']['ConnectionString'];
    function mapFields(m: RegExpMatchArray) {
      dbUser = m[1];
      dbPassword = m[2];
      dbHost = m[3];
      dbPort = m[4];
      dbName = m[5];
    }
    let m;
    m = url.match(/^sqlserver:\/\/(.+):(.+)@(.+):(\d+)\?database=(.+)$/);
    if (m) {
      mapFields(m);
      return;
    }
    m = url.match(/^(.+):(.+)@tcp\((.+):(\d+)\)\/(.+)\?/);
    if (m) {
      mapFields(m);
      return;
    }
    m = url.match(/^postgresql:\/\/(.+):(.+)@(.+):(\d+)\/(.+)$/);
    if (m) {
      mapFields(m);
      return;
    }
    if (dbPort === '') dbPort = getDefaultSQLPort(($shadowConfig as any)['Db']['Dialect']);
  }

  function saveConfiguration(): Promise<void> {
    if (['mysql', 'postgres', 'sqlserver'].includes(selectedExportFormat)) setConnString();

    if (selectedTimeZone !== '') {
      ($shadowConfig as any)['Export']['TimeZone'] = selectedTimeZone;
      ($shadowConfig as any)['Export']['ModifiedAfter'] = dayjs(dateValue)
        .tz(selectedTimeZone)
        .format();
    }

    ($shadowConfig as any)['Export']['Inspection']['Completed'] = selectedStatus;
    ($shadowConfig as any)['Export']['Inspection']['Archived'] = selectedArchived;
    ($shadowConfig as any)['Report']['Format'] = prepareReportFormatForSave();
    ($shadowConfig as any)['Session']['ExportType'] = selectedExportFormat;

    if (($shadowConfig as any)['Export']['Tables'].length > 0) {
      if (
        ($shadowConfig as any)['Export']['Media'] === true &&
        !($shadowConfig as any)['Export']['Tables'].includes('inspection_items')
      ) {
        ($shadowConfig as any)['Export']['Tables'].push('inspection_items');
      }
      if (
        ($shadowConfig as any)['Export']['Tables'].includes('inspection_items') &&
        !($shadowConfig as any)['Export']['Tables'].includes('inspections')
      ) {
        ($shadowConfig as any)['Export']['Tables'].push('inspections');
      }
    }

    if ($shadowConfig !== null && Object.keys($shadowConfig).length > 0) {
      return SaveSettings($shadowConfig);
    }
    return Promise.reject('empty configuration');
  }

  function getDefaultSQLPort(flavour: string): string {
    switch (flavour) {
      case 'mysql':
        return '3306';
      case 'postgres':
        return '5432';
      case 'sqlserver':
        return '1433';
      default:
        return '';
    }
  }

  function isValidPortNumber(input: string): boolean {
    return Number.isFinite(+input) && +input >= 1 && +input <= 65535;
  }

  function validateExport(): boolean {
    let hasError = false;
    if (['mysql', 'postgres', 'sqlserver'].includes(selectedExportFormat)) {
      if (dbHost.trim() === '') {
        dbHostShowError = true;
        hasError = true;
      } else {
        dbHostShowError = false;
      }
      if (dbPort.trim() === '' || !isValidPortNumber(dbPort)) {
        dbPortShowError = true;
        hasError = true;
      } else {
        dbPortShowError = false;
      }
      if (dbUser.trim() === '') {
        dbUserShowError = true;
        hasError = true;
      } else {
        dbUserShowError = false;
      }
      if (dbPassword.trim() === '') {
        dbPasswordShowError = true;
        hasError = true;
      } else {
        dbPasswordShowError = false;
      }
      if (dbName.trim() === '') {
        dbNameShowError = true;
        hasError = true;
      } else {
        dbNameShowError = false;
      }
    } else if (!['csv', 'sqlite', 'reports'].includes(selectedExportFormat)) {
      hasError = true;
    }
    return hasError;
  }

  function handleSaveAndExport() {
    showBanner = true;
    formError = validateExport();
    if (formError) return;

    saveConfiguration()
      .then(() => {
        switch (selectedExportFormat) {
          case 'csv':
          case 'sqlite':
            ValidateExportDirectory()
              .then((result: boolean) => {
                if (result) {
                  showBanner = false;
                  exportLocationError = false;
                  ($exportConfig as any)['items'] = getFeedsForExport();
                  if (selectedExportFormat === 'csv') {
                    ExportCSV();
                  } else {
                    ExportSQLite();
                  }
                  push('/export/status');
                } else {
                  showBanner = false;
                  exportLocationError = true;
                }
              })
              .catch(() => {
                exportLocationError = true;
              });
            break;
          case 'mysql':
          case 'postgres':
          case 'sqlserver':
            CheckDBConnection()
              .then(() => {
                ($exportConfig as any)['items'] = getFeedsForExport();
                showBanner = false;
                ExportSQL();
                push('/export/status');
              })
              .catch(() => {
                dbError = true;
              });
            break;
          case 'reports':
            ValidateExportDirectory()
              .then((result: boolean) => {
                if (result) {
                  showBanner = false;
                  exportLocationError = false;
                  ($exportConfig as any)['items'] = ['inspections', 'reports'];
                  ExportReports();
                  push('/export/status');
                } else {
                  showBanner = false;
                  exportLocationError = true;
                }
              })
              .catch(() => {
                exportLocationError = true;
              });
            break;
        }
      })
      .catch((e: unknown) => {
        console.debug('saveConfiguration err', e);
      });
  }

  function getFeedsForExport(): string[] {
    let feedsToExport: string[] = [];
    const tables = ($shadowConfig as any)['Export']['Tables'];
    if (tables !== null && tables.length > 0) feedsToExport = Array.from(tables);
    if (feedsToExport.length === 0) feedsToExport = Array.from(allTables);
    if (($shadowConfig as any)['Export']['Media'] === true && !feedsToExport.includes('media')) {
      feedsToExport.push('media');
    }
    return feedsToExport;
  }

  function handleSaveAndClose() {
    saveConfiguration().then(() => {
      Quit();
    });
  }

  function openFolderDialog() {
    if (build === 'windows' || build === '') return;
    SelectDirectory(($shadowConfig as any)['Export']['Path']).then((result: string) => {
      if (result !== '') {
        ($shadowConfig as any)['Export']['Path'] = result;
        ($shadowConfig as any)['Export']['MediaPath'] = result + '/media/';
      }
    });
  }

  function handlePortInput(e: Event) {
    const target = e.target as HTMLInputElement;
    dbPort = target.value.replace(/[^0-9]/g, '');
  }

  function trimPath(value: string, size: number = 35): string {
    if (value.length <= size) return value;
    return '...' + value.substring(value.length - size);
  }

  parseDbConnectionString();
</script>

{#if !formError && showBanner}
  <Dialog.Root open={true}>
    <Dialog.Content class="sm:max-w-md">
      <Dialog.Header>
        <Dialog.Title>Exporting...</Dialog.Title>
        <Dialog.Description>This might take a while...</Dialog.Description>
      </Dialog.Header>
    </Dialog.Content>
  </Dialog.Root>
{/if}

{#if dbError}
  <Dialog.Root open={true}>
    <Dialog.Content class="sm:max-w-md">
      <Dialog.Header>
        <Dialog.Title>Database connection error</Dialog.Title>
        <Dialog.Description
          >Please ensure the database details are correct and try again.</Dialog.Description
        >
      </Dialog.Header>
      <Dialog.Footer>
        <Button
          onclick={() => {
            dbError = false;
            showBanner = false;
          }}>Go back</Button
        >
      </Dialog.Footer>
    </Dialog.Content>
  </Dialog.Root>
{/if}

<div class="px-8 pt-6">
  <section class="flex items-center gap-3">
    <Button variant="ghost" size="icon" onclick={() => push('/welcome')}>
      <ArrowLeftIcon class="size-5" />
    </Button>
    <h1 class="text-xl font-semibold">Export configuration</h1>
  </section>

  <div class="mt-6 flex gap-6">
    <!-- Left: Filters -->
    <section class="flex-1 space-y-5">
      <div>
        <h2 class="text-sm font-semibold">Filters</h2>
        <p class="mt-1 text-sm text-muted-foreground">
          Select which sets of data you want to export from your organization.
        </p>
      </div>

      <!-- Template selector -->
      <div class="space-y-1.5">
        <Label class="text-sm font-medium">Select templates</Label>
        <Button
          variant="outline"
          class="h-10 w-full justify-between text-sm font-normal"
          onclick={() => push('/config/templates')}
        >
          {generateTemplateName()}
          <ChevronRightIcon class="size-4 text-muted-foreground" />
        </Button>
      </div>

      <!-- Data set selector -->
      <div class="space-y-1.5">
        <Label class="text-sm font-medium">Select data sets</Label>
        <Button
          variant="outline"
          class="h-10 w-full justify-between text-sm font-normal"
          onclick={() => push('/config/datasets')}
        >
          {generateDataSetName()}
          <ChevronRightIcon class="size-4 text-muted-foreground" />
        </Button>
      </div>

      <!-- Date range -->
      <div class="space-y-1.5">
        <Label class="text-sm font-medium">Date range from (UTC)</Label>
        <Popover.Root bind:open={datePickerOpen} onOpenChange={handleDatePickerClose}>
          <Popover.Trigger>
            {#snippet child({ props })}
              <Button variant="outline" class="h-10 w-full justify-between font-normal" {...props}>
                {dateInputValue}
                <CalendarIcon class="size-4 text-muted-foreground" />
              </Button>
            {/snippet}
          </Popover.Trigger>
          <Popover.Content class="w-auto p-3" align="start">
            <Input
              type="text"
              placeholder="DD/MM/YYYY"
              bind:value={dateInputValue}
              oninput={handleDateInput}
              class="mb-3 {dateInputError ? 'border-destructive ring-destructive/30 ring-2' : ''}"
            />
            <Calendar
              type="single"
              bind:value={calendarValue}
              minValue={minCalDate}
              maxValue={maxCalDate}
              captionLayout="dropdown"
              onValueChange={handleCalendarSelect}
            />
          </Popover.Content>
        </Popover.Root>
      </div>

      <!-- Status -->
      <div class="space-y-1.5">
        <Label class="text-sm font-medium">Include completed or incomplete inspections</Label>
        <Select.Root type="single" bind:value={selectedStatus}>
          <Select.Trigger class="w-full">
            {statusItems.find((i) => i.value === selectedStatus)?.label ?? 'Select...'}
          </Select.Trigger>
          <Select.Content>
            {#each statusItems as item (item.value)}
              <Select.Item value={item.value} label={item.label} />
            {/each}
          </Select.Content>
        </Select.Root>
      </div>

      <!-- Archived -->
      <div class="space-y-1.5">
        <Label class="text-sm font-medium">Include active or archived inspections</Label>
        <Select.Root type="single" bind:value={selectedArchived}>
          <Select.Trigger class="w-full">
            {archivedItems.find((i) => i.value === selectedArchived)?.label ?? 'Select...'}
          </Select.Trigger>
          <Select.Content>
            {#each archivedItems as item (item.value)}
              <Select.Item value={item.value} label={item.label} />
            {/each}
          </Select.Content>
        </Select.Root>
      </div>
    </section>

    <!-- Right: Export details -->
    <section
      class="w-[380px] shrink-0 space-y-4 overflow-y-auto rounded-lg border border-border p-5"
      style="max-height: calc(100vh - 160px);"
    >
      <h2 class="text-sm font-semibold">Export details</h2>

      <div class="space-y-1.5">
        <Label class="text-sm font-medium">Export data as</Label>
        <Select.Root
          type="single"
          bind:value={selectedExportFormat}
          onValueChange={handleExportFormatUpdate}
        >
          <Select.Trigger class="w-full">
            {dataExportFormatItems.find((i) => i.value === selectedExportFormat)?.label ??
              'Select...'}
          </Select.Trigger>
          <Select.Content>
            {#each dataExportFormatItems as item (item.value)}
              <Select.Item value={item.value} label={item.label} />
            {/each}
          </Select.Content>
        </Select.Root>
      </div>

      {#if ['mysql', 'postgres', 'sqlserver'].includes(selectedExportFormat)}
        <div class="space-y-3">
          <Label class="text-sm font-medium">Database details</Label>

          <div class="space-y-1">
            <Label class="text-xs text-muted-foreground">Host address</Label>
            <Input bind:value={dbHost} class={dbHostShowError ? 'border-destructive' : ''} />
            {#if dbHostShowError}<p class="text-xs text-destructive">Host cannot be empty</p>{/if}
          </div>

          <div class="space-y-1">
            <Label class="text-xs text-muted-foreground">Host port</Label>
            <Input
              placeholder={dbPortPlaceholder}
              maxlength={5}
              bind:value={dbPort}
              oninput={handlePortInput}
              class={dbPortShowError ? 'border-destructive' : ''}
            />
            {#if dbPortShowError}<p class="text-xs text-destructive">
                Please enter a valid port number
              </p>{/if}
          </div>

          <div class="space-y-1">
            <Label class="text-xs text-muted-foreground">Username</Label>
            <Input bind:value={dbUser} class={dbUserShowError ? 'border-destructive' : ''} />
            {#if dbUserShowError}<p class="text-xs text-destructive">
                Username cannot be empty
              </p>{/if}
          </div>

          <div class="space-y-1">
            <Label class="text-xs text-muted-foreground">Password</Label>
            <Input
              type="password"
              bind:value={dbPassword}
              class={dbPasswordShowError ? 'border-destructive' : ''}
            />
            {#if dbPasswordShowError}<p class="text-xs text-destructive">
                Password cannot be empty
              </p>{/if}
          </div>

          <div class="space-y-1">
            <Label class="text-xs text-muted-foreground">Database name</Label>
            <Input bind:value={dbName} class={dbNameShowError ? 'border-destructive' : ''} />
            {#if dbNameShowError}<p class="text-xs text-destructive">Name cannot be empty</p>{/if}
          </div>

          <Separator />
        </div>
      {/if}

      {#if selectedExportFormat === 'reports'}
        <div class="space-y-1.5">
          <Label class="text-sm font-medium">Report format</Label>
          <Select.Root type="single" bind:value={selectedReportFormat}>
            <Select.Trigger class="w-full">
              {reportFormatItems.find((i) => i.value === selectedReportFormat)?.label ??
                'Select...'}
            </Select.Trigger>
            <Select.Content>
              {#each reportFormatItems as item (item.value)}
                <Select.Item value={item.value} label={item.label} />
              {/each}
            </Select.Content>
          </Select.Root>
        </div>
      {/if}

      {#if selectedExportFormat}
        <div class="space-y-1.5">
          <Label class="text-sm font-medium">Folder location</Label>
          <Button
            variant="outline"
            class="h-10 w-full justify-between text-sm font-normal {exportLocationError
              ? 'border-destructive text-destructive'
              : ''}"
            disabled={build === 'windows'}
            onclick={openFolderDialog}
            title={($shadowConfig as any)['Export']['Path']}
          >
            <span class="truncate">{trimPath(($shadowConfig as any)['Export']['Path'])}</span>
            <FolderIcon class="size-4 shrink-0 text-muted-foreground" />
          </Button>
          {#if build === 'windows'}
            <p class="text-xs text-muted-foreground">
              To change folder location on Windows, move the executable to the new export folder
            </p>
          {/if}
        </div>
      {/if}

      <div class="space-y-1.5">
        <Label class="text-sm font-medium">Export time zone</Label>
        <Select.Root type="single" bind:value={selectedTimeZone}>
          <Select.Trigger class="w-full">
            {selectedTimeZone || 'Select...'}
          </Select.Trigger>
          <Select.Content>
            <Select.Item value="UTC" label="UTC" />
          </Select.Content>
        </Select.Root>
      </div>

      <div class="flex items-center gap-2">
        <Checkbox id="media" bind:checked={($shadowConfig as any)['Export']['Media']} />
        <Label for="media" class="cursor-pointer text-sm">Inspection media</Label>
      </div>

      <div class="flex gap-2 pt-2">
        <Button variant="outline" class="flex-1" onclick={handleSaveAndClose}>Save and close</Button
        >
        <Button class="flex-1" onclick={handleSaveAndExport}>Save and export</Button>
      </div>
    </section>
  </div>
</div>
