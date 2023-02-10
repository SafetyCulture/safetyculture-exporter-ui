<script>
    import {
        OpenDirectory,
        GetSettingDir
    } from "../../wailsjs/go/main/App.js"

    import {latestVersion} from '../lib/store.js';

    import {BrowserOpenURL} from "../../wailsjs/runtime/runtime.js";
    import {isNullOrEmptyObject} from "../lib/utils.js";

    async function openFolderDialog() {
        OpenDirectory(await GetSettingDir())
    }

    function openURL(url) {
        BrowserOpenURL(url)
    }


</script>

<div class="bar">
    <div>
        {#if !isNullOrEmptyObject($latestVersion)}
            <span>Current version: {$latestVersion['current']}</span>
            {#if $latestVersion['current'] !== $latestVersion['latest'] && $latestVersion['latest'] !== ''}
                <span class="accent m-left-16 block-link" on:click={openURL($latestVersion['download_url'])} on:keydown={openURL($latestVersion['download_url'])}>Latest version available: {$latestVersion['latest']}</span>
            {/if}
        {/if}
    </div>
    <div class="accent block-link" on:click={openFolderDialog} on:keypress={openFolderDialog}>Open logs</div>
</div>

<style>
    .bar {
        position: absolute;
        padding: 14px 16px;
        background-color: #F8F9FC;
        color: #1D2330;
        display: flex;
        justify-content: space-between;
        height: auto;
        width: 100%;
        bottom: 0;
    }

    .accent {
        color: #4740D4;
    }
</style>
