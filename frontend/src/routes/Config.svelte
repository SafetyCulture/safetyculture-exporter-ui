<script>
	import './common.css';
	import { onMount } from 'svelte';
	import Select from 'svelte-select';
	import { Datepicker } from 'svelte-calendar';
	import {SaveSettings, SelectDirectory, GetUserHomeDirectory} from "../../wailsjs/go/main/App.js"
	import {shadowConfig} from '../lib/store.js';
	import {Quit} from "../../wailsjs/runtime/runtime.js";
	import {push} from "svelte-spa-router";

	const statusItems = [
		{value: "true", label: "Completed only"},
		{value: "false", label: "Incomplete only"},
		{value: "both", label: "Both - completed and incomplete"}
	];

	const archivedItems = [
		{value: "true", label: "Archived Only"},
		{value: "false", label: "Unarchived Only"},
		{value: "both", label: "Both - archived and unarchived"}
	];

	const dataExportFormatItems = [
		{value: "csv", label: "CSV"},
		{value: "sql", label: "SQL"}
	];

	const reportFormatItems = [
		{value: "PDF", label: "PDF"},
		{value: "WORD", label: "Word"},
		{value: "both", label: "Both - PDF and Word"},
	];

	const timezoneItems = [
		{value: "utc", label: "UTC"}
	];

	let selectedExportFormat;
	let date = new Date();
	let exportFolder = setExportFolder();

	function generateTemplateName() {
		if ($shadowConfig["Export"]["TemplateIds"].length === 0) {
			return "All templates selected"
		} else {
			return $shadowConfig["Export"]["TemplateIds"].length + " templates selected"
		}
	}

	function generateTableName() {
		if ($shadowConfig["Export"]["Tables"].length === 0) {
			return "All tables selected"
		} else {
			return $shadowConfig["Export"]["Tables"].length + " tables selected"
		}
	}

	function setExportFolder() {
		if ($shadowConfig["Export"]["Path"].length === 0) {
			GetUserHomeDirectory().then(result => {
				$shadowConfig["Export"]["Path"] = result
			})
		}
	}

	function handleDataExport(event) {
		selectedExportFormat = event.detail.value;
	}

	function saveConfiguration() {
		if($shadowConfig !== {}) {
			SaveSettings($shadowConfig)
		}
	}

	function handleSaveAndExport() {
		saveConfiguration()
		// TODO trigger export
	}

	function handleSaveAndClose() {
		saveConfiguration()
		const millisecondsToWait = 500;
		setTimeout(function() {
			Quit()
		}, millisecondsToWait);
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
			}
		})
	}


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
			<div class="label">Date range</div>
			<div class="sub-label text-weak">From:</div>
			<div class="button-long selector border-weak border-round-8">
				<Datepicker bind:value={date} />
			</div>
			<div class="label">Include inspections with the following status:</div>
			<select class="custom-select m-top-8" bind:value={$shadowConfig["Export"]["Inspection"]["Completed"]}>
				{#each statusItems as item}
				<option value={item.value}>{item.label}</option>
				{/each}
			</select>
			<div class="label">Include archived inspections?</div>
			<select class="custom-select m-top-8" bind:value={$shadowConfig["Export"]["Inspection"]["Archived"]}>
				{#each archivedItems as item}
					<option value={item.value}>{item.label}</option>
				{/each}
			</select>
		</section>
		<section class="export-details border-round-8">
			<div class="h3">Export details</div>
			<div class="label">Data export format</div>
			<div class="border-weak border-round-8 m-top-4">
				<Select
					items={dataExportFormatItems}
					isClearable={false}
					on:change={handleDataExport}
					bind:value={selectedExportFormat}
				>
				</Select>
			</div>
			{#if selectedExportFormat != null && selectedExportFormat.value == 'sql'}
				<div>
					<div class="label">Database details:</div>
					<div class="sub-label text-weak">Host Address</div>
					<input class="input" type="text">
					<div class="sub-label text-weak">Host Port</div>
					<input class="input" type="text">
					<div class="sub-label text-weak">Username</div>
					<input class="input" type="text">
					<div class="sub-label text-weak">Password</div>
					<input class="input" type="password">
					<div class="sub-label text-weak">Name</div>
					<input class="input" type="text">

				</div>
			{/if}
			<div class="label">Report format</div>
			<Select
				items={reportFormatItems}
				isClearable={false}
			>
			</Select>
			<div class="folder-title">
				<div class="label">Folder location</div>
				<div class="link text-size-small">Want to change location?</div>
			</div>
			<div class="button-long selector border-weak border-round-8" on:click={openFolderDialog}>
				<div class="text-weak">{$shadowConfig["Export"]["Path"]}</div>
				<img src="../images/folder.png" alt="folder icon" width="15" height="15">
			</div>
			<div class="label">Export timezone</div>
			<Select
				items={timezoneItems}
				isClearable={false}
			>
			</Select>
			<div class="label">Include:</div>
			<input type="checkbox" id="media" name="media" value="media">
			<label class="text-size-medium" for="media">Media</label>
		</section>
	</div>
</div>

<style>
	body {
		-ms-overflow-style: none; /* for Internet Explorer, Edge */
		scrollbar-width: none; /* for Firefox */
		overflow-y: hidden;
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

	select.custom-select {
		cursor: pointer;

		/* styling */
		background-color: #FFFFFF;
		background-image: url("../images/arrow-down-compact.png");
		background-position: right 16px center;
		background-repeat: no-repeat;
		background-size: 8px;
		border: 1px solid #BFC5D4;
		border-radius: 8px;
		display: inline-block;
		font: inherit;
		line-height: 1.5em;
		padding: 0.5em 3.5em 0.5em 1em;

		/* reset */
		width: 100%;
		-webkit-box-sizing: border-box;
		-moz-box-sizing: border-box;
		box-sizing: border-box;
		-webkit-appearance: none;
		-moz-appearance: none;
		appearance: none;
		outline: 0;
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
</style>
