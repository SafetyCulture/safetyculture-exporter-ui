<script>
    import {writable} from "svelte/store";
    import {onMount} from "svelte";
    import {EventsOn} from "../../../wailsjs/runtime/runtime.js";
    import {markExportCompleted} from "../../lib/utils.js";
    import Pill from "./Pill.svelte";

    const statusQueued = "Queued"
    const statusComplete = "Complete"
    const statusInProgress = "In Progress"

    export let name;
    export let status = writable(statusQueued);
    export let remaining = writable(-1)
    export let completed = writable(false);

    let pillLabel = {value: statusQueued, class: 'pill queued'};


    function formatName(str) {
        return str.toLowerCase().replace(/_/g, ' ').replace(/(^|\s)\S/g, (L) => L.toUpperCase());
    }

    onMount(() => {
        EventsOn("update-"+name, (newValue) => {
            console.log('EVENT' + name)
            remaining.set(newValue);
            status.set(statusInProgress);
            if (newValue === 0) {
                status.set(statusComplete);
                completed.set(true)
                markExportCompleted();
            }
        })
    })
</script>

<td class="status-col-1">{formatName(name)}</td>
<td class="status-col-2">
    <Pill name={pillLabel.value} />
</td>
<td class="status-col-3">{($remaining === -1 || $remaining === 0) ? "-" : $remaining}</td>

<style>
    td {
        padding-top: 16px;
        padding-bottom: 16px;
    }
</style>
