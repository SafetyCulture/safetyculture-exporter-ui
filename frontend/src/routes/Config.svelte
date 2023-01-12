<script>
	import './common.css';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import Select from 'svelte-select';
	import { DateInput } from 'date-picker-svelte'
	import {
		SaveSettings,
		SelectDirectory,
		GetUserHomeDirectory,
		ReadExportStatus,
		ExportCSV,
		ReloadConfig, ExportSQL
	} from "../../wailsjs/go/main/App.js"
	import {shadowConfig} from '../lib/store.js';
	import {Quit} from "../../wailsjs/runtime/runtime.js";
	import {push} from "svelte-spa-router";

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
		{value: "reports", label: "Reports"},
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

	let dbHost='', dbPort, dbUser='', dbPassword='', dbName='';
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
	let date = convertStringToDate($shadowConfig["Export"]["ModifiedAfter"], selectedTimeZone);
	const minDate = dayjs().add(-1, 'year').toDate()


	function generateTemplateName() {
		let num = $shadowConfig["Export"]["TemplateIds"].length
		switch (num) {
			case 0: return "All templates selected"
			case 1: return "1 template selected"
			default: return num + " templates selected"
		}
	}

	function generateTableName() {
		if ($shadowConfig["Export"]["Tables"].length === 0) {
			return "All tables selected"
		} else {
			return $shadowConfig["Export"]["Tables"].length + " tables selected"
		}
	}

	function handleExportFormatUpdate(event) {
		selectedExportFormat = event.detail.value;
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
			SaveSettings($shadowConfig)
		}

		setTimeout(function() {
			ReloadConfig()
		}, 500);
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

	function handleSaveAndExport() {
		saveConfiguration()
		if (selectedExportFormat != null && selectedExportFormat.value === 'sql') {
			ExportSQL()
		} else {
			ExportCSV()
		}
		push("/exportStatus")

	}

	function handleSaveAndClose() {
		saveConfiguration()
		setTimeout(function() {
			// Quit()
		}, 500);
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
		SelectDirectory($shadowConfig["Export"]["Path"]).then(result => {
			if (result !== "") {
				$shadowConfig["Export"]["Path"] = result
				$shadowConfig["Export"]["MediaPath"] = result + '/media/'
			}
		})
	}

	parseDbConnectionString();
</script>

<div class="config-page p-48">
	<section class="top-nav">
		<div class="nav-left">
			<div class="block-link p-left-8" on:click={handleBackButton} on:keypress={handleBackButton}>
				<img src="../images/arrow-left.png" alt="back arrow icon" width="15" height="15">
			</div>
			<div class="h1 p-left-8">Export Configuration</div>
		</div>
		<div class="nav-right">
			<button class="button button-white border-round-12" on:click={handleSaveAndClose} on:keypress={handleSaveAndClose}>Save and Close</button>
			<button class="button button-purple m-left-8 border-round-12" on:click={handleSaveAndExport}>Save and Export</button>
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
			<div class="label">Date range From</div>
			<div class="m-top-8">
				<DateInput min={minDate} max={new Date()} format="dd-MM-yyyy" bind:value={date} />
			</div>
			<div class="label">Include inspections with the following status:</div>
			<div class="border-weak border-round-8 m-top-4">
				<Select items={statusItems} clearable={false} showChevron={true} searchable={false} --border="0px" bind:value={selectedStatus} />
			</div>

			<div class="label">Include archived inspections?</div>
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
					<div class="sub-label text-weak">Host Address</div>
					<input class="input" type="text" bind:value={dbHost}>
					<div class="sub-label text-weak">Host Port</div>
					<input class="input" type="text" bind:value={dbPort}>
					<div class="sub-label text-weak">Username</div>
					<input class="input" type="text" bind:value={dbUser}>
					<div class="sub-label text-weak">Password</div>
					<input class="input" type="password" bind:value={dbPassword}>
					<div class="sub-label text-weak">Name</div>
					<input class="input" type="text" bind:value={dbName}>
					<hr>
				</div>
			{/if}

			{#if selectedExportFormat != null && selectedExportFormat.value === 'reports' }
				<div class="label">Report format</div>
				<div class="border-weak border-round-8 m-top-4">
					<Select items={reportFormatItems} clearable={false} showChevron={true} searchable={false} --border="0px" bind:value={selectedReportFormat}/>
				</div>
			{/if}

			<div class="folder-title">
				<div class="label">Folder location</div>
			</div>

			<div id="folder" class="button-long selector border-weak border-round-8" on:click={openFolderDialog} on:keypress={openFolderDialog}>
				<div class="text-weak" >{$shadowConfig["Export"]["Path"]}</div>
				<img class="cursor-pointer" src="../images/folder.png" alt="folder icon" width="15" height="15">
			</div>
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

<style>
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

	.folder-title {
		display: flex;
		align-items: baseline;
		justify-content: space-between;
	}

	.cursor-pointer {
		cursor: pointer;
	}

</style>
