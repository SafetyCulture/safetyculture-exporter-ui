<script>
    import {
        SelectSettingsDirectory
    } from "../../wailsjs/go/main/App.js"

    import {latestVersion} from '../lib/store.js';

    import {BrowserOpenURL} from "../../wailsjs/runtime/runtime.js";
    import {isNullOrEmptyObject} from "../lib/utils.js";

    function openFolderDialog() {
        SelectSettingsDirectory()
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
                <span class="latest m-left-16 block-link" on:click={openURL($latestVersion['download_url'])} on:keydown={openURL($latestVersion['download_url'])}>Latest version available: {$latestVersion['latest']}</span>
            {/if}
        {/if}
    </div>
    <div class="block-link" on:click={openFolderDialog} on:keypress={openFolderDialog}>Open logs</div>
</div>

<style>
    .bar {
        position: absolute;
        padding: 5px 10px;
        background-color: #675DF4;
        color: whitesmoke;
        display: flex;
        justify-content: space-between;
        height: auto;
        width: 100%;
        bottom: 0;
    }

    .latest {
        color: #81E8F2;
    }

    .block-link:hover {
        color: yellow;
    }
</style>
