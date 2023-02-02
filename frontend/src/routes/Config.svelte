<script>
	import './common.css';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import Select from 'svelte-select';
	import { DateInput } from 'date-picker-svelte'
	import {
		SaveSettings, SelectDirectory, ExportCSV, ExportSQL, ReadBuild,CheckDBConnection
	} from "../../wailsjs/go/main/App.js"
	import {latestVersion, shadowConfig} from '../lib/store.js';
	import {BrowserOpenURL, Quit} from "../../wailsjs/runtime/runtime.js";
	import {push} from "svelte-spa-router";
	import FormTextInput from "../components/FormTextInput.svelte";
	import Button from "../components/Button.svelte";
	import {isNullOrEmptyObject} from "../lib/utils.js";
	import Overlay from "../components/Overlay.svelte";
	import StatusBar from "../components/StatusBar.svelte";

	let build = ""
	ReadBuild().then(it => {
		build = it
	})

	const statusItems = [
		{value: "true", label: "Completed only"},
		{value: "false", label: "Incomplete only"},
		{value: "both", label: "Both - completed and incomplete"}
	];
	let selectedStatus = $shadowConfig["Export"]["Inspection"]["Completed"]

	const archivedItems = [
		{value: "true", label: "Archived Only"},
		{value: "false", label: "Unarchived Only"},
		{value: "both", label: "Both - archived and unarchived"}
	];
	let selectedArchived = $shadowConfig["Export"]["Inspection"]["Archived"]

	const dataExportFormatItems = [
		{value: "csv", label: "CSV"},
		{value: "mysql", label: "MySQL"},
		{value: "postgres", label: "Postgres"},
		{value: "sqlserver", label: "SQL Server"},
		// {value: "reports", label: "Reports"},
	];

	const POSTGRES_DIALECT = 'postgres';
	const SQLSERVER_DIALECT = 'sqlserver';
	const MYSQL_DIALECT = 'mysql';

	const connectionStrings = {
		postgres: 'postgresql://${dbUser}:${dbPassword}@${dbHost}:${dbPort}/${dbName}',
		sqlserver: 'sqlserver://${dbUser}:${dbPassword}@${dbHost}:${dbPort}?database=${dbName}',
		mysql: '${dbUser}:${dbPassword}@tcp(${dbHost}:${dbPort})/${dbName}?charset=utf8mb4&parseTime=True&loc=Local',
	};
	const dialects = {
		postgres: POSTGRES_DIALECT,
		sqlserver: SQLSERVER_DIALECT,
		mysql: MYSQL_DIALECT,
	};

	let dbHost = '', dbHostShowError = false
	let dbPort='', dbPortPlaceholder = "e.g. " + getDefaultSQLPort($shadowConfig['Db']['Dialect']), dbPortShowError = false
	let dbUser='', dbUserShowError = false
	let dbPassword='', dbPasswordShowError = false
	let dbName='', dbNameShowError = false
	let formError = false
	let dbError = false

	let selectedExportFormat = $shadowConfig["Session"]["ExportType"];

	const reportFormatItems = [
		{value: "PDF", label: "PDF"},
		{value: "WORD", label: "Word"},
		{value: "BOTH", label: "Both PDF and Word"},
	];
	let selectedReportFormat = readReportFormat()

	const timezoneItems = [
		{value: "UTC", label: "UTC"}
	];
	let selectedTimeZone = $shadowConfig["Export"]["TimeZone"]

	function readReportFormat() {
		if ($shadowConfig["Report"]["Format"] === ["PDF"]) {
			return "PDF"
		}
		if ($shadowConfig["Report"]["Format"] === ["WORD"]) {
			return "WORD"
		}
		if ($shadowConfig["Report"]["Format"].includes("PDF") && $shadowConfig["Report"]["Format"].includes("WORD")) {
			return "BOTH"
		}
		return "PDF"
	}

	function prepareReportFormatForSave() {
		switch (selectedReportFormat.value) {
			case "PDF":
				return ["PDF"]
			case "WORD":
				return ["WORD"]
			case "BOTH":
				return ["PDF", "WORD"]
			default:
				return ["PDF"]
		}
	}

	// DATE PICKER
	dayjs.extend(utc);
	dayjs.extend(timezone);
	const minDate = dayjs().add(-1, 'year').toDate()
	let date = convertStringToDate($shadowConfig["Export"]["ModifiedAfter"], selectedTimeZone);

	function generateTemplateName() {
		let num = $shadowConfig["Export"]["TemplateIds"].length
		switch (num) {
			case 0: return "All templates selected"
			case 1: return "1 template selected"
			default: return num + " templates selected"
		}
	}

	function generateTableName() {
		let num = $shadowConfig["Export"]["Tables"].length
		switch (num) {
			case 0: return "All tables selected"
			case 1: return "1 table selected"
			default: return num + " tables selected"
		}
	}

	function handleExportFormatUpdate(event) {
		selectedExportFormat = event.detail.value;
		if (['mysql', 'postgres', 'sqlserver'].includes(selectedExportFormat)) {
			dbPortPlaceholder = "e.g. " + getDefaultSQLPort(selectedExportFormat)
		}
	}

	function setConnString() {
		if (selectedExportFormat !== '') {
			const connectionString = connectionStrings[selectedExportFormat.value];
			const replacedConnectionString = connectionString.replace(/\${dbUser}/g, dbUser)
					.replace(/\${dbPassword}/g, dbPassword)
					.replace(/\${dbHost}/g, dbHost)
					.replace(/\${dbPort}/g, dbPort)
					.replace(/\${dbName}/g, dbName);

			$shadowConfig['Db']['ConnectionString'] = replacedConnectionString;
			$shadowConfig['Db']['Dialect'] = dialects[selectedExportFormat.value];
		}
	}

	function parseDbConnectionString() {
		const url = $shadowConfig["Db"]["ConnectionString"];

		function mapFields(dbStringMatch) {
			dbUser = dbStringMatch[1];
			dbPassword = dbStringMatch[2];
			dbHost = dbStringMatch[3];
			dbPort = dbStringMatch[4];
			dbName = dbStringMatch[5];
		}

		let dbStringMatch;

		// Parse SQL Server connection string
		dbStringMatch = url.match(/^sqlserver:\/\/(.+):(.+)@(.+):(\d+)\?database=(.+)$/);
		if (dbStringMatch) {
			mapFields(dbStringMatch);
			return;
		}

		// Parse MySQL connection string
		dbStringMatch = url.match(/^(.+):(.+)@tcp\((.+):(\d+)\)\/(.+)\?/);
		if (dbStringMatch) {
			mapFields(dbStringMatch);
			return;
		}

		// Parse Postgres connection string
		dbStringMatch = url.match(/^postgresql:\/\/(.+):(.+)@(.+):(\d+)\/(.+)$/);
		if (dbStringMatch) {
			mapFields(dbStringMatch);
			return;
		}
	}

	function saveConfiguration() {
		const validFormats = {
			mysql: true,
			postgres: true,
			sqlserver: true
		};

		if (selectedExportFormat.value in validFormats) {
			setConnString();
		}

		if (selectedTimeZone.value !== '') {
			$shadowConfig["Export"]["TimeZone"] = selectedTimeZone.value
			$shadowConfig["Export"]["ModifiedAfter"] = convertDateToString(date, selectedTimeZone.value)
		}

		$shadowConfig["Export"]["Inspection"]["Completed"] = selectedStatus.value
		$shadowConfig["Export"]["Inspection"]["Archived"] = selectedArchived.value
		$shadowConfig["Report"]["Format"] = prepareReportFormatForSave()
		$shadowConfig["Session"]["ExportType"] = selectedExportFormat.value

		if($shadowConfig !== {}) {
			return SaveSettings($shadowConfig)
		}
		return Promise.reject("empty configuration")
	}

	function convertDateToString(dt, tz) {
		return dayjs(dt).tz(tz).format()
	}

	function convertStringToDate(input, tz) {
		if (input === "") {
			return minDate
		}

		return dayjs(input).tz(tz).toDate()
	}

	function getDefaultSQLPort(flavour) {
		switch (flavour) {
			case 'mysql':
				return '3306'
			case 'postgres':
				return '5432'
			case 'sqlserver':
				return '1433'
			default:
				return '1234'
		}
	}

	function isValidPortNumber(input) {
		return Number.isFinite(+input) && input >= 1 && input <= 65535
	}

	function validateExport() {
		let hasError = false

		switch (selectedExportFormat.value) {
			case 'csv':
				break
			case 'mysql':
			case 'postgres':
			case 'sqlserver':
				if (dbHost.trim() === '') {
					dbHostShowError = true
					hasError = true
				} else {
					dbHostShowError = false
				}

				if (dbPort.trim() === '' || isValidPortNumber(dbPort) === false) {
					dbPortShowError = true
					hasError = true
				} else {
					dbPortShowError = false
				}

				if (dbUser.trim() === '') {
					dbUserShowError = true
					hasError = true
				} else {
					dbUserShowError = false
				}

				if (dbPassword.trim() === '') {
					dbPasswordShowError = true
					hasError = true
				} else {
					dbPasswordShowError = false
				}

				if (dbName.trim() === '') {
					dbNameShowError = true
					hasError = true
				} else {
					dbNameShowError = false
				}
				break
			case 'reports':
				break
			default:
				hasError = true
				break
		}

		return hasError
	}

	function handleSaveAndExport() {
		formError = validateExport()
		if (formError === true) {
			return
		}

		saveConfiguration().then(it => {
			switch (selectedExportFormat.value) {
				case 'csv':
					ExportCSV()
					push("/export/status")
					break
				case 'mysql':
				case 'postgres':
				case 'sqlserver':
					CheckDBConnection().then(() => {
						ExportSQL()
						push("/export/status")
					})
					.catch(() => {
						dbError = true
					})
					break
				case 'reports':
					console.debug('NOT SUPPORTED')
					break
			}
		}).catch(e => {
			console.debug('saveConfiguration err')
			console.debug(e)
		})
	}

	function handleSaveAndClose() {
		saveConfiguration().then(it => {
			Quit()
		})
	}

	function handleBackButton() {
		push("/welcome")
	}

	function handleSelectTemplates() {
		push("/config/templates")
	}

	function handleTables() {
		push("/config/tables")
	}

	function openFolderDialog() {
		if (build === 'windows' || build === '') {
			return
		}
		SelectDirectory($shadowConfig["Export"]["Path"]).then(result => {
			if (result !== "") {
				$shadowConfig["Export"]["Path"] = result
				$shadowConfig["Export"]["MediaPath"] = result + '/media/'
			}
		})
	}

	function openURL(url) {
		BrowserOpenURL(url)
	}

	function removeOverlay() {
		dbError = false
	}

	parseDbConnectionString();
</script>

{#if !isNullOrEmptyObject($latestVersion) && $latestVersion["should_update"] === true && $latestVersion['current'] !== 'v0.0.0-dev'}
	<Overlay>
		<div class="download-alert" on:click={openURL($latestVersion['download_url'])} on:keydown={openURL($latestVersion['download_url'])}>
			<div>This version is not longer supported</div>
			<div>Latest version is {$latestVersion['latest']}</div>
			<div>Please click here to download it</div>
		</div>
	</Overlay>
{/if}

{#if dbError === true}
	<Overlay>
		<div class="db-error" on:click={removeOverlay} on:keydown={removeOverlay}>
			<div>Error connecting to the database</div>
			<div>Please ensure the database details are correct</div>
			<div>Click <a on:click={removeOverlay} on:keydown={removeOverlay}>here</a> to go back</div>
		</div>
	</Overlay>
{/if}


<div class="config-page">
	<section class="top-nav">
		<div class="nav-left">
			<div class="block-link p-left-8" on:click={handleBackButton} on:keypress={handleBackButton}>
				<img src="../images/arrow-left.png" alt="back arrow icon" width="15" height="15">
			</div>
			<div class="h1 p-left-8">Export Configuration</div>
		</div>
		<div class="nav-right">
			<Button label="Save and close" type="active2" onClick={handleSaveAndClose}/>
			<Button label="Save and export" type="active" clazz="m-left-8" error={formError} onClick={handleSaveAndExport}/>
		</div>
	</section>
	<div class="config-body m-top-8">
		<section class="filters">
			<div class="filter-title">
				<div class="h3">Filters</div>
				<div class="text-weak m-top-8">Select which sets of data you want to export from your organization.</div>
			</div>
			<div class="label">Select templates</div>
				<div class="button-long selector border-weak border-round-8 block-link" on:click={handleSelectTemplates} on:keypress={handleSelectTemplates}>
					<div class="templates">{generateTemplateName()}</div>
					<div class="template-button-right">
						<img class="m-left-8" src="../images/arrow-right-compact.png" alt="right arrow icon" width="4" height="8">
					</div>
				</div>
			<div class="label">Select tables</div>
			<div class="button-long selector border-weak border-round-8 block-link" on:click={handleTables} on:keypress={handleTables}>
				<div class="templates">{generateTableName()}</div>
				<div class="template-button-right">
					<img class="m-left-8" src="../images/arrow-right-compact.png" alt="right arrow icon" width="4" height="8">
				</div>
			</div>
			<div class="label">Date range from</div>
			<div class="m-top-8">
				<DateInput max={new Date()} format="dd-MM-yyyy" bind:value={date} />
			</div>
			<div class="label">Include completed or incomplete inspections</div>
			<div class="border-weak border-round-8 m-top-4">
				<Select items={statusItems} clearable={false} showChevron={true} searchable={false} --border="0px" bind:value={selectedStatus} />
			</div>

			<div class="label">Include active or archived inspections</div>
			<div class="border-weak border-round-8 m-top-4">
				<Select items={archivedItems} clearable={false} showChevron={true} searchable={false} --border="0px" bind:value={selectedArchived}/>
			</div>
		</section>
		<section class="export-details border-round-8">
			<div class="h3">Export details</div>
			<div class="label">Data export format</div>
			<div class="border-weak border-round-8 m-top-4">
				<Select items={dataExportFormatItems} clearable={false} showChevron={true} searchable={false} on:change={handleExportFormatUpdate} --border="0px" bind:value={selectedExportFormat} />
			</div>
			{#if selectedExportFormat != null && ['mysql', 'postgres', 'sqlserver'].includes(selectedExportFormat.value)}
				<div>
					<div class="label">Database details:</div>
					<FormTextInput label="Host address" placeholder="e.g. localhost" error={dbHostShowError} bind:value={dbHost}/>
					<FormTextInput label="Host port" placeholder={dbPortPlaceholder} error={dbPortShowError} bind:value={dbPort}/>
					<FormTextInput label="Username" placeholder="e.g. john" error={dbUserShowError} bind:value={dbUser}/>
					<FormTextInput label="Password" placeholder="e.q. mySup3rS3cr3t" error={dbPasswordShowError} bind:value={dbPassword}/>
					<FormTextInput label="Database name" placeholder="e.g. safetyculture" error={dbNameShowError} bind:value={dbName}/>
					<hr>
				</div>
			{/if}

			{#if selectedExportFormat != null && selectedExportFormat.value === 'reports' }
				<div class="label">Report format</div>
				<div class="border-weak border-round-8 m-top-4">
					<Select items={reportFormatItems} clearable={false} showChevron={true} searchable={false} --border="0px" bind:value={selectedReportFormat}/>
				</div>
			{/if}

			<div class="label">Folder location</div>
			<div id="folder" class="button-long selector border-weak border-round-8" on:click={openFolderDialog} on:keypress={openFolderDialog}>
				<div class="text-weak" >{$shadowConfig["Export"]["Path"]}</div>
				<img class="{build === 'windows' ? '' : 'cursor-pointer'}" src="../images/folder.png" alt="folder icon" width="15" height="15">
			</div>
			{#if build === 'windows'}
				<div class="sub-label m-top-4">To change folder location on Windows, move the executable to the new export folder</div>
			{/if}

			<div class="label">Export timezone</div>
			<div class="border-weak border-round-8 m-top-4">
				<Select items={timezoneItems} clearable={false} showChevron={true} searchable={false} --border="0px" bind:value={selectedTimeZone}/>
			</div>
			<div class="label">Include:</div>
			<input type="checkbox" id="media" name="media" bind:checked={$shadowConfig["Export"]["Media"]}>
			<label class="text-size-medium" for="media">Media</label>
		</section>
	</div>
</div>

<StatusBar/>

<style>
	.config-page {
		padding-top: var(--main-gutter-top);
		padding-left: var(--main-gutter-left);
		padding-right: var(--main-gutter-right);
	}
	.config-body {
		display: flex;
		justify-content: space-between;
	}

	.filters {
		width: 55%;
	}

	.template-button-right {
		display: flex;
		align-items: center;
	}

	.export-details {
		width: 380px;
		height: 600px;
		background-color: #E9EEF6;
		padding: 16px;
		overflow-y: auto;
	}

	.download-alert {
		text-align: center;
		cursor: pointer;
	}

	.db-error {
		text-align: center;
		cursor: pointer;
	}
</style>
