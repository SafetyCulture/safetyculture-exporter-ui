<script lang="ts">
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import {
		CheckDBConnection,
		ExportCSV,
		ExportReports,
		ExportSQL,
		ExportSQLite,
		ReadBuild,
		SaveSettings,
		SelectDirectory, ValidateExportDirectory
	} from "../../wailsjs/go/main/App.js"
	import { exportConfig, shadowConfig } from '../lib/store';
	import { allTables } from "../lib/utils";
	import { Quit } from "../../wailsjs/runtime/runtime.js";
	import { push } from "@keenmate/svelte-spa-router";
	import FormTextInput from "../components/FormTextInput.svelte";
	import Button from "../components/Button.svelte";
	import Overlay from "../components/Overlay.svelte";
	import StatusBar from "../components/StatusBar.svelte";
	import FormPassword from "../components/FormPassword.svelte";
	import FormNumberInput from "../components/FormNumberInput.svelte";
	import ButtonSelector from "../components/ButtonSelector.svelte";
	import FolderPicker from "../components/FolderPicker.svelte";

	let build = $state("");
	ReadBuild().then((it: string) => {
		build = it
	})

	const statusItems = [
		{ value: "true", label: "Completed only" },
		{ value: "false", label: "Incomplete only" },
		{ value: "both", label: "Both - completed and incomplete" }
	];
	let selectedStatus = $state(($shadowConfig as any)["Export"]["Inspection"]["Completed"] as string);

	const archivedItems = [
		{ value: "true", label: "Archived only" },
		{ value: "false", label: "Active only" },
		{ value: "both", label: "Both - active and archived" }
	];
	let selectedArchived = $state(($shadowConfig as any)["Export"]["Inspection"]["Archived"] as string);

	const dataExportFormatItems = [
		{ value: "csv", label: "CSV" },
		{ value: "mysql", label: "MySQL" },
		{ value: "postgres", label: "Postgres" },
		{ value: "sqlserver", label: "SQL Server" },
		{ value: "sqlite", label: "SQLite" },
		{ value: "reports", label: "Reports" },
	];

	const POSTGRES_DIALECT = 'postgres';
	const SQLSERVER_DIALECT = 'sqlserver';
	const MYSQL_DIALECT = 'mysql';

	const connectionStrings: Record<string, string> = {
		postgres: 'postgresql://${dbUser}:${dbPassword}@${dbHost}:${dbPort}/${dbName}',
		sqlserver: 'sqlserver://${dbUser}:${dbPassword}@${dbHost}:${dbPort}?database=${dbName}',
		mysql: '${dbUser}:${dbPassword}@tcp(${dbHost}:${dbPort})/${dbName}?charset=utf8mb4&parseTime=True&loc=Local',
	};
	const dialects: Record<string, string> = {
		postgres: POSTGRES_DIALECT,
		sqlserver: SQLSERVER_DIALECT,
		mysql: MYSQL_DIALECT,
	};

	let dbHost = $state(''), dbHostShowError = $state(false), dbHostErrMsg = 'Host cannot be empty'
	let dbPort = $state(''), dbPortPlaceholder = $state("e.g. " + getDefaultSQLPort(($shadowConfig as any)['Db']['Dialect'])), dbPortShowError = $state(false), dbPortErrMsg = 'Please enter a valid port number'
	let dbUser = $state(''), dbUserShowError = $state(false), dbUserErrMsg = 'Username cannot be empty'
	let dbPassword = $state(''), dbPasswordShowError = $state(false), dbPasswordErrMsg = 'Password cannot be empty'
	let dbName = $state(''), dbNameShowError = $state(false), dbNameErrMsg = 'Name cannot be empty'
	let formError = $state(false)
	let dbError = $state(false)
	let showBanner = $state(false)
	let exportLocationError = $state(false)

	let selectedExportFormat = $state(($shadowConfig as any)["Session"]["ExportType"] as string);

	const reportFormatItems = [
		{ value: "PDF", label: "PDF" },
		{ value: "WORD", label: "Word" },
		{ value: "BOTH", label: "Both PDF and Word" },
	];
	let selectedReportFormat = $state(readReportFormat());

	const timezoneItems = [
		{ value: "UTC", label: "UTC" }
	];
	let selectedTimeZone = $state(($shadowConfig as any)["Export"]["TimeZone"] as string);

	function readReportFormat(): string {
		const fmt = ($shadowConfig as any)["Report"]["Format"];
		if (Array.isArray(fmt) && fmt.includes("PDF") && fmt.includes("WORD")) {
			return "BOTH"
		}
		if (Array.isArray(fmt) && fmt.includes("WORD")) {
			return "WORD"
		}
		return "PDF"
	}

	function prepareReportFormatForSave(): string[] {
		switch (selectedReportFormat) {
			case "PDF": return ["PDF"]
			case "WORD": return ["WORD"]
			case "BOTH": return ["PDF", "WORD"]
			default: return ["PDF"]
		}
	}

	// DATE PICKER
	dayjs.extend(utc);
	dayjs.extend(timezone);
	const minDate = dayjs().add(-1, 'year').format('YYYY-MM-DD');
	const today = dayjs().format('YYYY-MM-DD');
	let dateValue = $state(convertToInputDate(($shadowConfig as any)["Export"]["ModifiedAfter"], selectedTimeZone));

	function convertToInputDate(input: string, tz: string): string {
		if (input === "" || input === "0001-01-01T00:00:00Z") {
			return minDate
		}
		return dayjs(input).tz(tz).format('YYYY-MM-DD')
	}

	function generateTemplateName(): string {
		let num = ($shadowConfig as any)["Export"]["TemplateIds"].length
		switch (num) {
			case 0: return "All templates selected"
			case 1: return "1 template selected"
			default: return num + " templates selected"
		}
	}

	function generateDataSetName(): string {
		let num = ($shadowConfig as any)["Export"]["Tables"].length
		switch (num) {
			case 0: return "All data sets selected"
			case 1: return "1 data set selected"
			default: return num + " data sets selected"
		}
	}

	function handleExportFormatUpdate(e: Event) {
		const target = e.target as HTMLSelectElement;
		selectedExportFormat = target.value;
		if (['mysql', 'postgres', 'sqlserver'].includes(selectedExportFormat)) {
			dbPortPlaceholder = "e.g. " + getDefaultSQLPort(selectedExportFormat)
			dbPort = getDefaultSQLPort(selectedExportFormat)
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
		const url = ($shadowConfig as any)["Db"]["ConnectionString"];
		function mapFields(dbStringMatch: RegExpMatchArray) {
			dbUser = dbStringMatch[1];
			dbPassword = dbStringMatch[2];
			dbHost = dbStringMatch[3];
			dbPort = dbStringMatch[4];
			dbName = dbStringMatch[5];
		}

		let dbStringMatch;

		dbStringMatch = url.match(/^sqlserver:\/\/(.+):(.+)@(.+):(\d+)\?database=(.+)$/);
		if (dbStringMatch) { mapFields(dbStringMatch); return; }

		dbStringMatch = url.match(/^(.+):(.+)@tcp\((.+):(\d+)\)\/(.+)\?/);
		if (dbStringMatch) { mapFields(dbStringMatch); return; }

		dbStringMatch = url.match(/^postgresql:\/\/(.+):(.+)@(.+):(\d+)\/(.+)$/);
		if (dbStringMatch) { mapFields(dbStringMatch); return; }

		if (dbPort === '') {
			dbPort = getDefaultSQLPort(($shadowConfig as any)['Db']['Dialect'])
		}
	}

	function saveConfiguration(): Promise<void> {
		const validFormats: Record<string, boolean> = { mysql: true, postgres: true, sqlserver: true };

		if (selectedExportFormat in validFormats) {
			setConnString();
		}

		if (selectedTimeZone !== '') {
			($shadowConfig as any)["Export"]["TimeZone"] = selectedTimeZone;
			($shadowConfig as any)["Export"]["ModifiedAfter"] = dayjs(dateValue).tz(selectedTimeZone).format();
		}

		($shadowConfig as any)["Export"]["Inspection"]["Completed"] = selectedStatus;
		($shadowConfig as any)["Export"]["Inspection"]["Archived"] = selectedArchived;
		($shadowConfig as any)["Report"]["Format"] = prepareReportFormatForSave();
		($shadowConfig as any)["Session"]["ExportType"] = selectedExportFormat;

		if (($shadowConfig as any)["Export"]["Tables"].length > 0) {
			if (($shadowConfig as any)["Export"]["Media"] === true && !($shadowConfig as any)["Export"]["Tables"].includes("inspection_items")) {
				($shadowConfig as any)["Export"]["Tables"].push("inspection_items")
			}

			if (($shadowConfig as any)["Export"]["Tables"].includes("inspection_items") && !($shadowConfig as any)["Export"]["Tables"].includes("inspections")) {
				($shadowConfig as any)["Export"]["Tables"].push("inspections")
			}
		}

		if ($shadowConfig !== null && Object.keys($shadowConfig).length > 0) {
			return SaveSettings($shadowConfig)
		}
		return Promise.reject("empty configuration")
	}

	function getDefaultSQLPort(flavour: string): string {
		switch (flavour) {
			case 'mysql': return '3306'
			case 'postgres': return '5432'
			case 'sqlserver': return '1433'
			default: return ''
		}
	}

	function isValidPortNumber(input: string): boolean {
		return Number.isFinite(+input) && +input >= 1 && +input <= 65535
	}

	function validateExport(): boolean {
		let hasError = false

		switch (selectedExportFormat) {
			case 'csv':
			case 'sqlite':
			case 'reports':
				break
			case 'mysql':
			case 'postgres':
			case 'sqlserver':
				if (dbHost.trim() === '') { dbHostShowError = true; hasError = true } else { dbHostShowError = false }
				if (dbPort.trim() === '' || isValidPortNumber(dbPort) === false) { dbPortShowError = true; hasError = true } else { dbPortShowError = false }
				if (dbUser.trim() === '') { dbUserShowError = true; hasError = true } else { dbUserShowError = false }
				if (dbPassword.trim() === '') { dbPasswordShowError = true; hasError = true } else { dbPasswordShowError = false }
				if (dbName.trim() === '') { dbNameShowError = true; hasError = true } else { dbNameShowError = false }
				break
			default:
				hasError = true
				break
		}

		return hasError
	}

	function handleSaveAndExport() {
		showBanner = true
		formError = validateExport()
		if (formError === true) return

		saveConfiguration().then(() => {
			switch (selectedExportFormat) {
				case 'csv':
					ValidateExportDirectory().then((result: boolean) => {
						if (result === true) {
							showBanner = false; exportLocationError = false;
							($exportConfig as any)['items'] = getFeedsForExport()
							ExportCSV(); push("/export/status")
						} else { showBanner = false; exportLocationError = true }
					}).catch(() => { exportLocationError = true })
					break
				case 'mysql':
				case 'postgres':
				case 'sqlserver':
					CheckDBConnection().then(() => {
						($exportConfig as any)['items'] = getFeedsForExport()
						showBanner = false; ExportSQL(); push("/export/status")
					}).catch(() => { dbError = true })
					break
				case 'reports':
					ValidateExportDirectory().then((result: boolean) => {
						if (result === true) {
							showBanner = false; exportLocationError = false;
							($exportConfig as any)['items'] = ['inspections', 'reports']
							ExportReports(); push("/export/status")
						} else { showBanner = false; exportLocationError = true }
					}).catch(() => { exportLocationError = true })
					break
				case 'sqlite':
					ValidateExportDirectory().then((result: boolean) => {
						if (result === true) {
							showBanner = false; exportLocationError = false;
							($exportConfig as any)['items'] = getFeedsForExport()
							ExportSQLite(); push("/export/status")
						} else { showBanner = false; exportLocationError = true }
					}).catch(() => { exportLocationError = true })
					break
			}
		}).catch((e: unknown) => {
			console.debug('saveConfiguration err', e)
		})
	}

	function getFeedsForExport(): string[] {
		let feedsToExport: string[] = []
		const tables = ($shadowConfig as any)["Export"]["Tables"];
		if (tables !== null && tables.length > 0) {
			feedsToExport = Array.from(tables)
		}
		if (feedsToExport.length === 0) {
			feedsToExport = Array.from(allTables)
		}
		if (($shadowConfig as any)["Export"]["Media"] === true && !feedsToExport.includes("media")) {
			feedsToExport.push("media")
		}
		return feedsToExport
	}

	function handleSaveAndClose() {
		saveConfiguration().then(() => { Quit() })
	}

	function handleBackButton() { push("/welcome") }
	function handleSelectTemplates() { push("/config/templates") }
	function handleTables() { push("/config/datasets") }

	function openFolderDialog() {
		if (build === 'windows' || build === '') return
		SelectDirectory(($shadowConfig as any)["Export"]["Path"]).then((result: string) => {
			if (result !== "") {
				($shadowConfig as any)["Export"]["Path"] = result;
				($shadowConfig as any)["Export"]["MediaPath"] = result + '/media/'
			}
		})
	}

	function removeOverlay() {
		dbError = false
		showBanner = false
	}

	parseDbConnectionString();
</script>

{#if formError === false && showBanner === true}
	<Overlay>This might take a while ...</Overlay>
{/if}

{#if dbError === true}
	<Overlay>
		<button class="cursor-pointer text-center" onclick={removeOverlay}>
			<div>Error connecting to the database</div>
			<div>Please ensure the database details are correct</div>
			<div>Click here to go back</div>
		</button>
	</Overlay>
{/if}

<div class="px-8 pt-8">
	<section class="flex items-center justify-between">
		<div class="flex items-center">
			<button class="cursor-pointer" onclick={handleBackButton}>
				<img src="../images/back.svg" alt="back arrow icon">
			</button>
			<div class="my-4 pl-4 font-semibold text-2xl">Export configuration</div>
		</div>
		<div class="flex items-center">
			<Button label="Save and close" type="active-white" onClick={handleSaveAndClose}/>
			<Button label="Save and export" type="active-purple" clazz="ml-2" error={formError} onClick={handleSaveAndExport}/>
		</div>
	</section>

	<div class="mt-2 flex justify-between">
		<section class="w-[55%]">
			<div>
				<div class="text-base font-semibold">Filters</div>
				<div class="mt-2 text-text-weak">Select which sets of data you want to export from your organization.</div>
			</div>
			<ButtonSelector label="Select templates" title={generateTemplateName()} onClick={handleSelectTemplates}/>
			<ButtonSelector label="Select data sets" title={generateDataSetName()} onClick={handleTables}/>

			<div class="mt-4 text-sm font-medium leading-4">Date range from (UTC)</div>
			<div class="mt-1">
				<input type="date" class="w-full rounded-lg border border-border px-4 py-2.5 text-base" min={minDate} max={today} bind:value={dateValue}/>
			</div>

			<div class="mt-4 text-sm font-medium leading-4">Include completed or incomplete inspections</div>
			<div class="mt-1">
				<select class="w-full rounded-lg border border-border bg-white px-4 py-2.5 text-base" bind:value={selectedStatus}>
					{#each statusItems as item}
						<option value={item.value}>{item.label}</option>
					{/each}
				</select>
			</div>

			<div class="mt-4 text-sm font-medium leading-4">Include active or archived inspections</div>
			<div class="mt-1">
				<select class="w-full rounded-lg border border-border bg-white px-4 py-2.5 text-base" bind:value={selectedArchived}>
					{#each archivedItems as item}
						<option value={item.value}>{item.label}</option>
					{/each}
				</select>
			</div>
		</section>

		<section class="h-[590px] w-[380px] overflow-y-auto rounded-lg bg-bg-light p-4">
			<div class="text-base font-semibold">Export details</div>

			<div class="mt-4 text-sm font-medium leading-4">Export data as:</div>
			<div class="mt-1">
				<select class="w-full rounded-lg border border-border bg-white px-4 py-2.5 text-base" bind:value={selectedExportFormat} onchange={handleExportFormatUpdate}>
					{#each dataExportFormatItems as item}
						<option value={item.value}>{item.label}</option>
					{/each}
				</select>
			</div>

			{#if selectedExportFormat != null && ['mysql', 'postgres', 'sqlserver'].includes(selectedExportFormat)}
				<div>
					<div class="mt-4 text-sm font-medium leading-4">Database details:</div>
					<FormTextInput label="Host address" error={dbHostShowError} errorMsg={dbHostErrMsg} bind:value={dbHost}/>
					<FormNumberInput label="Host port" placeholder={dbPortPlaceholder} maxlength={5} error={dbPortShowError} errorMsg={dbPortErrMsg} bind:value={dbPort} />
					<FormTextInput label="Username" error={dbUserShowError} errorMsg={dbUserErrMsg} bind:value={dbUser}/>
					<FormPassword label="Password" error={dbPasswordShowError} errorMsg={dbPasswordErrMsg} bind:value={dbPassword}/>
					<FormTextInput label="Database name" error={dbNameShowError} errorMsg={dbNameErrMsg} bind:value={dbName}/>
					<hr class="my-2 border-border">
				</div>
			{/if}

			{#if selectedExportFormat != null && selectedExportFormat === 'reports'}
				<div class="mt-4 text-sm font-medium leading-4">Report format</div>
				<div class="mt-1">
					<select class="w-full rounded-lg border border-border bg-white px-4 py-2.5 text-base" bind:value={selectedReportFormat}>
						{#each reportFormatItems as item}
							<option value={item.value}>{item.label}</option>
						{/each}
					</select>
				</div>
			{/if}

			{#if selectedExportFormat != null}
				<div class="mt-4 text-sm font-medium leading-4">Folder location</div>
				<div class="mt-1">
					<FolderPicker value={($shadowConfig as any)["Export"]["Path"]} trimLongWords={true} onClick={build === 'windows' ? null : openFolderDialog} error={exportLocationError}/>
				</div>
			{/if}
			{#if build === 'windows'}
				<div class="mt-1 text-xs font-medium text-text-weak">To change folder location on Windows, move the executable to the new export folder</div>
			{/if}

			<div class="mt-4 text-sm font-medium leading-4">Export time zone</div>
			<div class="mt-1">
				<select class="w-full rounded-lg border border-border bg-white px-4 py-2.5 text-base" bind:value={selectedTimeZone}>
					{#each timezoneItems as item}
						<option value={item.value}>{item.label}</option>
					{/each}
				</select>
			</div>

			<div class="mt-4 text-sm font-medium leading-4">Include:</div>
			<div class="mt-1 flex items-center gap-2">
				<input type="checkbox" id="media" name="media" class="size-5 accent-primary" bind:checked={($shadowConfig as any)["Export"]["Media"]}/>
				<label class="text-sm" for="media">Inspection media</label>
			</div>
		</section>
	</div>
</div>

<StatusBar/>
