<script>
    import './common.css';
    import {ReadExportStatus} from "../../wailsjs/go/main/App.js"
    import Status from "./../components/Export/Status.svelte";
    import {shadowConfig} from "../lib/store.js";
    import {onMount} from "svelte";
    import {EventsOn} from "../../wailsjs/runtime/runtime.js";
    import {allTables} from "../lib/utils.js";

    let feedsToExport = []
    if ($shadowConfig["Export"]["Tables"] !== null && $shadowConfig["Export"]["Tables"].length > 0) {
        feedsToExport = $shadowConfig["Export"]["Tables"]
    }
    if (feedsToExport === []) {
        feedsToExport = allTables
    }

    console.log(feedsToExport)

    let exportCompleted = false

    onMount(() => {
        EventsOn("finished-export", (newValue) => {
            console.debug('RECEIVED EVENT > ' + 'finished-export' + " with value: " + newValue)
            if (newValue === true) {
                exportCompleted = true
            }
        })
    })

    ReadExportStatus();
</script>

<div class="status-page">
    <img class="logo" src="../images/logo.png" alt="SafetyCulture logo"/>

    <div class="progress-title m-top-32">
        <div class="inline">
            {#if exportCompleted}
                <img id="status-completed" src='/images/completed.png' alt="export completed icon">
            {:else}
                <img id="status-in-progress" src='/images/in-progress.png' alt="export in progress icon">
            {/if}
        </div>
        <div class="inline status-title p-left-8">Export Status</div>
    </div>

    <div class="progress-body m-top-16">
        <table class="status-table">
            <thead>
                <tr class="text-weak">
                    <th class="status-col-1">Export item</th>
                    <th class="status-col-2">Progress</th>
                    <th class="status-col-3">Remaining</th>
                </tr>
            </thead>
            <tbody>
            {#each feedsToExport as feed}
                <tr><Status name={feed}></Status></tr>
            {/each}
            </tbody>
        </table>
    </div>
</div>

<style>
    .status-page {
        padding: 20px;
        background-color: #E9EEF6;
        height: 100%;
    }

    .logo {
        width: 150px;
    }

    .progress-title {
        display: flex;
        align-items: center;
        justify-content: start;
    }

    .progress-title .inline {
        display: flex;
    }

    .status-title {
        font-size: 20px;
        font-weight: 600;
    }

    .progress-body {
        background-color: white;
        height: calc( 100% - 100px );
        padding: 20px 16px;
    }

    .status-table {
        width: 100%;
        border-collapse: collapse;
        overflow: hidden;
    }

    .status-table th {
        font-size: 14px;
        font-weight: 500;
    }

    .status-table tr {
        font-size: 14px;
        border-bottom: 1px solid #EEF1F7;
    }
</style>
