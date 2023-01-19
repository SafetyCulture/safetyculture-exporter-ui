<script>
    import {onMount} from "svelte";
    import {EventsOn} from "../../../wailsjs/runtime/runtime.js";
    import Pill from "./Pill.svelte";

    const statusQueued = "Queued"
    const statusComplete = "Complete"
    const statusInProgress = "In Progress"

    export let name = '';
    export let status = statusQueued

    let remaining = 0

    function formatExportItemName(str) {
        return str.toLowerCase().replace(/_/g, ' ').replace(/(^|\s)\S/g, (L) => L.toUpperCase());
    }

    onMount(() => {
        EventsOn("update-"+name, (newValue) => {
            console.debug('RECEIVED EVENT > ' + name + " with value: " + newValue)
            remaining = newValue.remaining
            status = statusInProgress
            if (remaining === 0) {
                status = statusComplete
            }
        })
    })
</script>

<td class="status-col-1">{formatExportItemName(name)}</td>
<td class="status-col-2">
    <Pill name={status} />
</td>
<td class="status-col-3">{(remaining === -1 || remaining === 0) ? "-" : remaining}</td>

<style>
    td {
        padding-top: 16px;
        padding-bottom: 16px;
    }
</style>
