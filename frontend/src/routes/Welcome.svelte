<script>
	import './common.css';

	import {push} from 'svelte-spa-router'
	import {ValidateApiKey} from "../../wailsjs/go/main/App.js"
	import {shadowConfig, templateCache} from '../lib/store.js';
	import ValidatableInput from "../components/ValidatableInput.svelte";
	import Button from "../components/Button.svelte";
	import StatusBar from "../components/StatusBar.svelte";
	import {BrowserOpenURL} from "../../wailsjs/runtime/runtime.js";


	let isValid = false;
	let buttonLabel = "Verify"
	let displayBadApiKeyErr = false
	let displayConnectionErr = false
	let displayValidationError = false
	let tries = 1
	let buttonActive = true
	let accessToken = $shadowConfig["AccessToken"]

	function openURL() {
		BrowserOpenURL("https://app.safetyculture.com/account/api-tokens")
	}

	function validate() {
		if (buttonActive === false) {
			return;
		}
		tries++
		isValid = false
		if (accessToken.length === 0) {
			displayValidationError = true
			return
		}

		function checkErr(errMsg) {
			if (errMsg === "connection error") {
				displayConnectionErr = true
				displayValidationError = true
			} else {
				displayBadApiKeyErr = true
				displayValidationError = true
			}
		}

		buttonActive = false
		ValidateApiKey(accessToken).then((result) => {
			if (result !== "") {
				buttonLabel = "Try again"
				checkErr(result);
			} else {
				displayValidationError = false

				if ($shadowConfig["AccessToken"] !== accessToken) {
					$shadowConfig["Export"]["TemplateIds"] = []
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
<div class="welcome-page">
	<section class="welcome-left-side">
		<section class="welcome-header">
			<img id="welcome-page-logo" src="../images/logo.svg" alt="SafetyCulture logo"/>
			<div class="h1">Welcome to SafetyCulture Exporter</div>
		</section>
		<section class="token-validation">
			<div class="token-validation-text">Generate an API token from your SafetyCulture <span class="link" on:click={openURL} on:keypress={openURL}>user profile</span>.</div>
			<ValidatableInput placeholder="Enter API Token here" error={displayValidationError} bind:value={accessToken}/>

			{#if displayBadApiKeyErr}
				<div class="error-block">
					<div class="error-block-title">Unable to verify token</div>
					<div class="error-block-body">It looks like your token may be expired after 30 days of inactivity. Please generate a new token and try again.</div>
				</div>
			{/if}

			{#if displayConnectionErr}
				<div class="error-block">
					<div class="error-block-title">Connection error</div>
					<div class="error-block-body">It looks like you are not connected to the internet or behind a firewall. Please check your connection and try again.</div>
				</div>
			{/if}

			<Button label={buttonLabel} type="active-purple" active={buttonActive} error={displayValidationError} clazz="m-top-8" onClick={validate}/>
		</section>

		<section class="storage-info">
			<div class="note border-round-8">
				<div>
					<img src="../images/warning.svg" alt="alert icon" width="24" height="24">
				</div>
				<div>
					<div class="note-title">Before you continue</div>
					<div class="note-body">All files (apart from SQL) you export will be stored in the same place on your computer or server as the SafetyCulture Exporter. If you want to change where your files get exported, please move the SafetyCulture Exporter file itself to that place.</div>
				</div>
			</div>
		</section>
	</section>
	<section class="welcome-right-side">
        <div class="right-side-bg bg-light border-round-12-left ">
            <div class="right-image">
                <img class="border-round-12-left" src="../images/token.svg" alt="example generating token">
            </div>
            <div class="img-caption text-weak">Example API token</div>
        </div>
	</section>
</div>

<StatusBar/>

<style>
	.welcome-page {
		display: flex;
	}

	.welcome-left-side {
		width: 50%;
		padding-left: 24px;
		padding-top: 48px;
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
		width: 95%;
	}

	.note .note-title {
		font-weight: bold;
		padding-bottom: 8px;
	}

	.note .note-body {
		line-height: 1.5rem;
	}

    .welcome-right-side {
        width: 50%;
        padding-top: 48px;
    }

    .right-side-bg {
        padding: 50px 0;
    }

    .right-image {
        display: flex;
        flex-direction: row-reverse;
    }

	.right-image img {
		height: 90%;
		width: 90%;
    }

    .img-caption {
        padding-top: 8px;
        text-align: center;
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
