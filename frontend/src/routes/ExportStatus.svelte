<script lang="ts">
    import { CancelExport, ReadExportStatus, OpenDirectory } from "../../wailsjs/go/main/App.js"
    import Status from "./../components/Export/Status.svelte";
    import { shadowConfig, exportConfig } from "../lib/store";
    import { onMount } from "svelte";
    import { EventsOn, Quit } from "../../wailsjs/runtime/runtime.js";
    import Button from "../components/Button.svelte";
    import { push } from "@keenmate/svelte-spa-router";
    import Overlay from "../components/Overlay.svelte";
    import StatusBar from "../components/StatusBar.svelte";

    let feedsToExport = ($exportConfig as any)['items'] as string[];
    let exportType = ($shadowConfig as any)["Session"]["ExportType"] as string;

    let cancelTriggered = $state(false);
    let exportCompleted = $state(false);

    onMount(() => {
        EventsOn("finished-export", (newValue: boolean) => {
            if (newValue === true) {
                exportCompleted = true
            }
        })
    })

    function handleCancel() {
        cancelTriggered = true
        CancelExport()
    }

    function handleClose() {
        Quit()
    }

    function goBack() {
        push("/config")
    }

    function openExportFolder() {
        OpenDirectory(($shadowConfig as any)["Export"]["Path"])
    }

    ReadExportStatus();
</script>

<div class="h-full bg-bg-light pt-8 pr-8 pl-8">
    <section class="flex items-center justify-between">
        <div class="flex items-center">
            <div class="my-4 font-semibold text-2xl">Export status</div>
        </div>
        <div class="flex items-center">
            <div class="inline">
                {#if cancelTriggered}
                    <img src='/images/warning-red.svg' alt="export cancelled icon">
                {:else if exportCompleted}
                    <img src='/images/complete.svg' alt="export completed icon">
                {:else}
                    <img src='/images/in-progress.svg' alt="export in progress icon">
                {/if}
            </div>
            <div class="inline flex items-center px-2 pr-4 text-sm">
                {#if cancelTriggered}
                    Export cancelled
                {:else if exportCompleted}
                    Export complete
                {:else}
                    In progress
                {/if}
            </div>

            {#if !exportCompleted}
                <Button label="Cancel export" type="active-red" onClick={handleCancel}/>
            {:else}
                {#if !cancelTriggered}
                    {#if exportType === "csv" || exportType === "reports"}
                        <Button label="Open export folder" type="active-white" onClick={openExportFolder}/>
                    {/if}
                    <Button label="Close" clazz="ml-2" type="active-purple" onClick={handleClose}/>
                {:else}
                    {#if exportType === "reports"}
                        <Button label="Open export folder" type="active-white" onClick={openExportFolder}/>
                    {/if}
                    <Button label="Go Back" type="active-purple" onClick={goBack}/>
                {/if}
            {/if}
        </div>
    </section>

    <div id="overlay-cancel-export">
        {#if cancelTriggered && exportCompleted === false}
            <Overlay>Cancelling export...</Overlay>
        {/if}
    </div>

    <div class="mt-4 h-[calc(100vh-200px)] overflow-y-scroll rounded-lg bg-white px-4 py-5">
        <table class="w-full border-collapse">
            <thead>
                <tr class="text-sm font-medium text-text-weak">
                    <th class="w-[40%] text-left">Export item</th>
                    <th class="w-[20%] text-left">Status</th>
                    <th class="text-left">&nbsp;</th>
                </tr>
            </thead>
            <tbody>
            {#each feedsToExport as feed}
                <tr class="border-b border-[#EEF1F7] text-sm"><Status name={feed} cancelled={cancelTriggered}></Status></tr>
            {/each}
            </tbody>
        </table>
    </div>
</div>

<StatusBar/>
