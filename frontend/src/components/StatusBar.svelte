<script lang="ts">
    import {
        OpenDirectory,
        GetSettingDir
    } from "../../wailsjs/go/main/App.js"

    import { latestVersion } from '../lib/store';

    import { BrowserOpenURL } from "../../wailsjs/runtime/runtime.js";
    import { isNullOrEmptyObject } from "../lib/utils";

    async function openFolderDialog() {
        OpenDirectory(await GetSettingDir())
    }

    function openURL(url: string) {
        BrowserOpenURL(url)
    }

    let currentYear = new Date().getFullYear();
</script>

<div class="fixed bottom-0 flex w-full justify-between bg-bg-subtle px-4 py-3.5 text-foreground">
    <div>
        {#if !isNullOrEmptyObject($latestVersion)}
            <span>Current version: {$latestVersion['current']}</span>
            {#if $latestVersion['current'] !== $latestVersion['latest'] && $latestVersion['latest'] !== ''}
                {#if $latestVersion['download_url'] !== ''}
                <button class="ml-4 cursor-pointer text-primary" onclick={() => openURL($latestVersion['download_url'])}>Latest version available: {$latestVersion['latest']}</button>
                {:else}
                <span class="ml-4">Latest version: {$latestVersion['latest']}</span>
                {/if}
            {/if}
        {/if}
    </div>
    <div>
        <button class="cursor-pointer text-primary" onclick={openFolderDialog}>Open logs</button>
        <span class="ml-4 text-sm">Copyright &copy; {currentYear}</span>
    </div>
</div>
