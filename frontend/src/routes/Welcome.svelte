<script lang="ts">
	import { push } from '@keenmate/svelte-spa-router'
	import { ValidateApiKey } from "../../wailsjs/go/main/App.js"
	import { shadowConfig, templateCache } from '../lib/store';
	import ValidatableInput from "../components/ValidatableInput.svelte";
	import Button from "../components/Button.svelte";
	import StatusBar from "../components/StatusBar.svelte";
	import { BrowserOpenURL } from "../../wailsjs/runtime/runtime.js";

	let isValid = $state(false);
	let buttonLabel = $state("Verify");
	let displayBadApiKeyErr = $state(false);
	let displayConnectionErr = $state(false);
	let displayValidationError = $state(false);
	let tries = $state(1);
	let buttonActive = $state(true);
	let accessToken = $state($shadowConfig["AccessToken"] as string ?? '');

	function openURL() {
		BrowserOpenURL("https://app.safetyculture.com/account/api-tokens")
	}

	function validate() {
		tries++
		isValid = false
		if (accessToken.length === 0) {
			displayValidationError = true
			return
		}

		function checkErr(errMsg: string) {
			if (errMsg === "connection error") {
				displayConnectionErr = true
				displayValidationError = true
			} else {
				displayBadApiKeyErr = true
				displayValidationError = true
			}
		}

		buttonActive = false
		ValidateApiKey(accessToken).then((result: string) => {
			if (result !== "") {
				buttonLabel = "Try again"
				checkErr(result);
			} else {
				displayValidationError = false

				if ($shadowConfig["AccessToken"] !== accessToken) {
					($shadowConfig as any)["Export"] = ($shadowConfig as any)["Export"] ?? {};
					($shadowConfig as any)["Export"]["TemplateIds"] = []
				}
				templateCache.set([]);
				$shadowConfig["AccessToken"] = accessToken
				push("/config")
			}
		}).finally(() => {
			buttonActive = true
		})
	}
</script>

<div class="flex h-full flex-col items-center pb-15">
	<img class="w-[150px] pt-8" src="../images/logo.svg" alt="SafetyCulture logo"/>
	<div class="my-4 font-semibold text-[1.8rem]">Welcome to SafetyCulture Exporter</div>
	<img class="w-[600px]" src="../images/welcome.png" alt="welcome"/>
	<div class="mb-2 pt-4 text-base leading-6">Generate an API token from your <button class="cursor-pointer text-primary" onclick={openURL}>SafetyCulture account</button>.</div>

	<div class="flex">
		<div class="w-[360px]">
			<ValidatableInput placeholder="Enter API Token here" error={displayValidationError} bind:value={accessToken}/>
		</div>

		<div class="pl-2">
			<Button label={buttonLabel} type="active-purple" active={buttonActive} error={displayValidationError} onClick={validate}/>
		</div>
	</div>

	{#if displayBadApiKeyErr}
		<div class="text-sm text-foreground">
			<div class="text-danger-dark">Unable to verify token</div>
			<div>It looks like your token may be expired after 30 days of inactivity. Please generate a new token and try again.</div>
		</div>
	{/if}

	{#if displayConnectionErr}
		<div class="text-sm text-foreground">
			<div class="text-danger-dark">Connection error</div>
			<div>It looks like you are not connected to the internet or behind a firewall. Please check your connection and try again.</div>
		</div>
	{/if}

	<section class="m-4">
		<div class="flex gap-4 rounded-lg bg-bg-light p-4 text-sm text-text-gray">
			<div>
				<img src="../images/warning.svg" alt="alert icon">
			</div>
			<div>
				<div class="pb-2 font-bold">Important note</div>
				<div class="leading-6">All files (apart from SQL) you export will be stored in the same place on your computer or server as the SafetyCulture Exporter. If you want to change where your files get exported, please move the SafetyCulture Exporter file itself to that place.</div>
			</div>
		</div>
	</section>
</div>

<StatusBar/>
