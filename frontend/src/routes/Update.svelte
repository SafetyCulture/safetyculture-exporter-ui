<script lang="ts">
    import { latestVersion } from '../lib/store';
    import { Button } from "$lib/components/ui/button";
    import { TriggerUpdate } from "../../wailsjs/go/main/App.js";
    import { push } from "@keenmate/svelte-spa-router";
    import { BrowserOpenURL, Quit } from "../../wailsjs/runtime/runtime.js";

    let updateStatus = $state("updating");
    let updateMessage = $state('Please wait until we update your application from ' + $latestVersion['current'] + ' to ' + $latestVersion['latest']);
    let cancelActive = $state(false);
    let restartActive = $state(false);

    TriggerUpdate($latestVersion['download_url']).then((result: boolean) => {
        if (result === true) {
            updateStatus = "success"
            updateMessage = 'We have updated your application to the ' + $latestVersion['latest']
            cancelActive = false
            restartActive = true
        } else {
            updateStatus = "failed"
            updateMessage = 'There was an error updating to version ' + $latestVersion['latest']
            cancelActive = true
            restartActive = false
        }
    }).catch(() => {
        updateStatus = "failed"
        cancelActive = true
        restartActive = false
    })

    function openURL(url: string) {
        if (url !== '') BrowserOpenURL(url)
    }
</script>

<div class="flex h-screen flex-col items-center">
    <img class="w-[150px] pt-8" src="../images/logo.svg" alt="SafetyCulture logo"/>
    <h1 class="my-4 font-semibold text-[1.8rem]">SafetyCulture Exporter Updater</h1>

    <div class="flex grow flex-col items-center justify-center">
        {#if updateStatus === 'updating'}
            <img class="size-[200px]" src="../images/spinning.gif" alt="loading"/>
        {:else if updateStatus === 'success'}
            <img class="size-[200px]" src="../images/complete.svg" alt="ok"/>
        {:else}
            <img class="size-[200px]" src="../images/warning.svg" alt="warning"/>
        {/if}

        <p class="pt-16 text-base font-semibold">{updateMessage}</p>

        {#if updateStatus === 'failed'}
            {#if $latestVersion['os'] === 'darwin'}
                <p class="pt-2 text-sm text-muted-foreground">SafetyCulture Exporter must be moved into the Applications folder in order for the auto-update to work</p>
            {/if}
            <button class="cursor-pointer pt-2 text-sm text-primary hover:underline" onclick={() => openURL($latestVersion['download_url'])}>
                You can manually download and install the Exporter
            </button>
        {/if}

        <div class="flex gap-2 pt-8">
            <Button variant="default" onclick={() => push("/welcome")} disabled={!cancelActive}>Cancel</Button>
            <Button variant="default" onclick={() => Quit()} disabled={!restartActive}>Restart</Button>
        </div>
    </div>
</div>
