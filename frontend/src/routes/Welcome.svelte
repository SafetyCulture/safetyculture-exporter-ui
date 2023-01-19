<script>
	import './common.css';

	import {push} from 'svelte-spa-router'
	import {ValidateApiKey} from "../../wailsjs/go/main/App.js"
	import {shadowConfig, templateCache} from '../lib/store.js';
	import ValidatableInput from "../components/ValidatableInput.svelte";
	import Button from "../components/Button.svelte";


	let isValid = false;
	let buttonLabel = "Verify"
	let displayBadApiKeyErr = false
	let displayValidationError = false
	let tries = 1

	function validate() {
		tries++
		isValid = false
		if ($shadowConfig["AccessToken"].length === 0) {
			displayValidationError = true
			return
		}

		ValidateApiKey($shadowConfig["AccessToken"]).then((result) => {
			isValid = result
			if (isValid === false) {
				buttonLabel = "Try again"
				displayBadApiKeyErr = true
				displayValidationError = true
			} else {
				displayValidationError = false
				templateCache.set([]);
				push("/config")
			}
		})
	}
</script>
<div class="welcome-page">
	<section class="welcome-left-side">
		<section class="welcome-header">
			<img id="welcome-page-logo" src="../images/logo.png" alt="SafetyCulture logo"/>
			<div class="h1">Welcome to SafetyCulture Exporter</div>
		</section>
		<section class="token-validation">
			<div class="token-validation-text">Generate an API token from your SafetyCulture <span class="link">user profile</span>.</div>
			<ValidatableInput placeholder="Enter API Token here" error={displayValidationError} bind:value={$shadowConfig["AccessToken"]}/>

			{#if displayBadApiKeyErr}
				<div class="error-block">
					<div class="error-block-title">Unable to verify token</div>
					<div class="error-block-body">It looks like your token may be expired after 30 days of inactivity. Please generate a new token and try again.</div>
				</div>
			{/if}

			<Button label={buttonLabel} type="active" error={displayValidationError} clazz="m-top-8" onClick={validate}/>
		</section>

		<section class="storage-info">
			<div class="note border-round-8">
				<div>
					<img src="../images/round_exclamation_mark.png" alt="alert icon" width="20" height="20">
				</div>
				<div>
					<div class="note-title">Before you continue</div>
					<div class="note-body">All files (apart from SQL) you export will be stored in the same place on your computer or server as the SafetyCulture Exporter. If you want to change where your files get exported, please move the SafetyCulture Exporter file itself to that place.</div>
				</div>
			</div>
		</section>
	</section>
	<section class="welcome-right-side">
		<div class="right-image">
			<img src="../images/token_example.png" alt="example generating token">
		</div>
	</section>
</div>

<style>
	.welcome-page {
		display: flex;
	}

	.welcome-left-side {
		width: 50%;
		padding: 20px;
	}

	.welcome-right-side {
		width: 50%;
	}

	#welcome-page-logo {
		width: 150px;
	}

	div.token-validation-text {
		margin-bottom: 8px;
		color: #1D2330;
	}

	.token-validation-text {
		font-style: normal;
		font-weight: 400;
		font-size: 1rem;
		line-height: 1.5rem;
	}

	.note {
		background-color: #EEF1F7;
		padding: 16px;
		color: #3F495A;
		display: flex;
		flex-direction: row;
		gap: 16px;
		font-size: 0.9rem;
	}

	.note .note-title {
		font-weight: bold;
		padding-bottom: 8px;
	}

	.note .note-body {
		line-height: 1.5rem;
	}

	div.right-image {
		max-width: 100%;
		max-height: 100%;
	}

	.right-image img {
		height: 100%;
		width: 100%;
		object-fit: contain;
	}

	div.error-block {
		font-size: 0.8rem;
		color: #1D2330;
	}

	div.error-block .error-block-title {
		color: #A02228;
	}

	section.storage-info {
		margin-top: 24px;
	}
</style>
