<script>
    import {
        ReadVersion, GetLatestVersion, SelectSettingsDirectory
    } from "../../wailsjs/go/main/App.js"

    import {latestVersion} from '../lib/store.js';

    import {push} from "svelte-spa-router";
    import {BrowserOpenURL} from "../../wailsjs/runtime/runtime.js";

    function gotoWelcome() {
        push("/welcome")
    }

    function gotoConfig() {
        push("/config")
    }

    function gotoTables() {
        push("/config/tables")
    }

    function gotoTemplates() {
        push("/config/templates")
    }

    function gotoStatus() {
        push("/exportStatus")
    }

    let version = ""
    ReadVersion().then(it => {
        version = it
    })

    if ($latestVersion === null || $latestVersion !== {}) {
        GetLatestVersion().then(result => {
            latestVersion.set(result);
        })
    }

    function openFolderDialog() {
        SelectSettingsDirectory()
    }

    function openURL(url) {
        BrowserOpenURL(url)
    }
</script>

<div class="bar">
    <div>
        <span>Current version: {version}</span>
        {#if $latestVersion !== {} && $latestVersion["Version"] !== version}
            <span class="latest m-left-16 block-link" on:click={openURL($latestVersion['DownloadURL'])} on:keydown={openURL($latestVersion['DownloadURL'])}>Latest version available: {$latestVersion['Version']}</span>
        {/if}
    </div>
    {#if version === 'v0.0.0-dev'}
        <div class="p-left-16">
            <button on:click={gotoWelcome}>WELCOME</button>
            <button on:click={gotoConfig}>CONFIG</button>
            <button on:click={gotoTemplates}>TEMPLATES</button>
            <button on:click={gotoTables}>TABLES</button>
            <button on:click={gotoStatus}>STATUS</button>
        </div>
    {/if}
    <div class="block-link" on:click={openFolderDialog} on:keypress={openFolderDialog}>Open logs</div>
</div>

<style>
    .bar {
        padding: 5px 10px;
        background-color: #675DF4;
        color: whitesmoke;
        display: flex;
        justify-content: space-between;
    }

    .latest {
        color: #81E8F2;
    }

    .block-link:hover {
        color: yellow;
    }
</style>
