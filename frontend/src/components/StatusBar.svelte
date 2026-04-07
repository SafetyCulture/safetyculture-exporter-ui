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

<footer class="fixed bottom-0 left-0 flex w-full items-center justify-between border-t border-border bg-background px-4 py-2.5 text-xs text-muted-foreground">
    <div class="flex items-center gap-3">
        {#if !isNullOrEmptyObject($latestVersion)}
            <span>v{$latestVersion['current']}</span>
            {#if $latestVersion['current'] !== $latestVersion['latest'] && $latestVersion['latest'] !== ''}
                <span class="text-border">|</span>
                {#if $latestVersion['download_url'] !== ''}
                    <button class="cursor-pointer text-xs text-primary hover:underline" onclick={() => openURL($latestVersion['download_url'])}>
                        Update available: {$latestVersion['latest']}
                    </button>
                {:else}
                    <span>Latest: {$latestVersion['latest']}</span>
                {/if}
            {/if}
        {/if}
    </div>
    <div class="flex items-center gap-3">
        <button class="cursor-pointer text-xs text-primary hover:underline" onclick={openFolderDialog}>Open logs</button>
        <span class="text-border">|</span>
        <span>&copy; {currentYear}</span>
    </div>
</footer>
