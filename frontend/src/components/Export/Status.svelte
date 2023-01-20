<script>
    import {onMount} from "svelte";
    import {EventsOn} from "../../../wailsjs/runtime/runtime.js";
    import Pill from "./Pill.svelte";

    const statusQueued = "Queued"
    const statusFailed = "Failed"
    const statusComplete = "Complete"
    const statusInProgress = "In Progress"

    export let name = '';
    export let status = statusQueued

    let remaining = 0
    let pillType = "neutral"

    function formatExportItemName(str) {
        return str.toLowerCase().replace(/_/g, ' ').replace(/(^|\s)\S/g, (L) => L.toUpperCase());
    }

    onMount(() => {
        EventsOn("update-"+name, (newValue) => {
            if (newValue['started'] === true && newValue['finished'] === false) {
                status = statusInProgress
                remaining = newValue['remaining']
            }

            if (newValue['started'] === true && newValue['finished'] === true) {
                if (newValue['has_error'] === false) {
                    status = statusComplete
                    remaining = 0
                } else {
                    status = statusFailed
                    remaining = 0
                }
            }

            switch (status) {
                case 'Complete':
                    pillType = 'success'
                    break
                case 'In Progress':
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
<td class="status-col-3">{(remaining === -1 || remaining === 0) ? "-" : remaining}</td>

<style>
    td {
        padding-top: 16px;
        padding-bottom: 16px;
    }
</style>
