<script lang="ts">
    import { latestVersion } from '../lib/store';
    import Button from "../components/Button.svelte";
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

    function cancelHandler() {
        push("/welcome")
    }

    function restartHandler() {
        Quit()
    }

    function openURL(url: string) {
        if (url !== '') {
            BrowserOpenURL(url)
        }
    }
</script>

<div class="flex h-screen flex-col items-center">
    <img class="w-[150px] pt-8" src="../images/logo.svg" alt="SafetyCulture logo"/>
    <div class="my-4 font-semibold text-[1.8rem]">SafetyCulture Exporter Updater</div>

    <div class="flex grow flex-col items-center justify-center">
        {#if updateStatus === 'updating'}
            <img class="size-[200px]" src="../images/spinning.gif" alt="loading"/>
        {/if}
        {#if updateStatus === 'success'}
            <img class="size-[200px]" src="../images/complete.svg" alt="ok"/>
        {/if}
        {#if updateStatus === 'failed'}
            <img class="size-[200px]" src="../images/warning.svg" alt="ok"/>
        {/if}

        <div class="pt-16 text-base font-semibold">{updateMessage}</div>
        {#if updateStatus === 'failed'}
            {#if $latestVersion['os'] === 'darwin'}
                <div class="pt-2">SafetyCulture Exporter must be moved into the Applications folder in order for the auto-update to work</div>
            {/if}
            <button class="cursor-pointer pt-2 text-sm text-[#0d75b5]" onclick={() => openURL($latestVersion['download_url'])}>
                You can manually download and install the Exporter
            </button>
        {/if}

        <div class="flex gap-2 pt-8">
            <Button label="Cancel" type="active-purple" active={cancelActive} onClick={cancelHandler}/>
            <Button label="Restart" type="active-purple" active={restartActive} onClick={restartHandler}/>
        </div>
    </div>
</div>
