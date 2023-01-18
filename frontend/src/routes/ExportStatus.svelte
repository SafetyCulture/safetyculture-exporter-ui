<script>
    import './common.css';
    import {ReadExportStatus} from "../../wailsjs/go/main/App.js"
    import Status from "./Status.svelte";
    import {exportCompleted, feedsToExport} from "../lib/utils.js";
    import {shadowConfig} from "../lib/store.js";

    if ($shadowConfig["Export"]["Tables"] !== null && $shadowConfig["Export"]["Tables"].length > 0) {
        $feedsToExport = $shadowConfig["Export"]["Tables"]
    }

    ReadExportStatus();
</script>

<div>
    <section class="top-nav">
        <div>
            <img class="logo" id="status-page-logo" src="../images/logo.png" alt="SafetyCulture logo"/>
        </div>
    </section>
    <section>
        <div class="status-nav">
            {#if $exportCompleted}
                <img id="status-completed" src='/images/completed.png' alt="export completed icon">
            {:else}
                <img id="status-in-progress" src='/images/in-progress.png' alt="export in progress icon">
            {/if}
            <div class="status-title">Export Status</div>
        </div>
    </section>
    <section>
        <table class="status-table">
            <thead>
                <tr>
                    <th class="status-table-th status-col-1">Export Item</th>
                    <th class="status-table-th status-col-2">Progress</th>
                    <th class="status-table-th status-col-3">Remaining</th>
                </tr>
            </thead>
            <tbody>
                {#each $feedsToExport as feed}
                    <tr><Status name={feed}></Status></tr>
                {/each}
            </tbody>
        </table>
    </section>
</div>

<style>
    .logo {
        width: 60%;
        height: 60%;
        padding-top: 20px;
        padding-left: 20px;
        padding-bottom: 20px;
    }

    .status-nav {
        display: flex;
        align-content: center;
    }
    .status-title {
        color: #1D2330;
        padding-left: 8px;
        padding-right: 8px;
        padding-bottom: 40px;
        text-justify: auto;
        font-size: 18px;
        font-weight: 600;
    }
</style>