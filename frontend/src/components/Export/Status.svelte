<script>
    import {onMount} from "svelte";
    import {EventsOn} from "../../../wailsjs/runtime/runtime.js";
    import Pill from "./Pill.svelte";

    const statusQueued = "Queued"
    const statusFailed = "Failed"
    const statusComplete = "Complete"
    const statusDownloading = "Downloading"
    const statusExporting = "Exporting"

    export let name = '';
    export let status = statusQueued

    let counter = -1
    let counterDecremental = true
    let pillType = "neutral"

    function formatExportItemName(str) {
        return str.toLowerCase().replace(/_/g, ' ').replace(/(^|\s)\S/g, (L) => L.toUpperCase());
    }

    onMount(() => {
        EventsOn("update-"+name, (newValue) => {
						counterDecremental = newValue['counter_decremental']
            if (newValue['started'] === true && newValue['finished'] === false) {
                switch (newValue['stage']) {
                    case 'API_DOWNLOAD':
                        status = statusDownloading
                        break
                    case 'CSV_EXPORT':
                        status = statusExporting
                        break
                }
                counter = newValue['counter']
            }

            if (newValue['started'] === true && newValue['finished'] === true) {
                if (newValue['has_error'] === false) {
                    status = statusComplete
	                  counter = 0
                } else {
                    status = statusFailed
	                  counter = 0
                }
            }

            switch (status) {
                case 'Complete':
                    pillType = 'success'
                    break
                case 'Downloading':
                case 'Exporting':
                    pillType = 'info'
                    break
                case 'Queued':
                    pillType = 'neutral'
                    break
                case 'Failed':
                    pillType = 'error'
                    break
            }
        })
    })
</script>

<td class="status-col-1">{formatExportItemName(name)}</td>
<td class="status-col-2">
    <Pill name={status} type={pillType}/>
</td>
<td class="status-col-3">
		{(counter === -1 || counter === 0) ? "-" : counter + " " + (counterDecremental === true ? "remaining" : "downloaded")}
</td>

<style>
    td {
        padding-top: 16px;
        padding-bottom: 16px;
    }
</style>
