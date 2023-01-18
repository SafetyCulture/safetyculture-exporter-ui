<script>
    import {writable} from "svelte/store";
    import {onMount} from "svelte";
    import {EventsOn} from "../../wailsjs/runtime/runtime.js";
    import {markExportCompleted} from "../lib/utils.js";


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

    function handleStatusUpdate() {
        switch ($status) {
            case statusComplete:
                pillLabel = { value: $status, class: 'pill complete' };
                break;
            case statusInProgress:
                pillLabel = { value: $status, class: 'pill in-progress' };
                break;
            case statusQueued :
                pillLabel = { value: $status, class: 'pill queued' };
                break;
        }
    }


    onMount(() => {
        EventsOn("update-"+name, (newValue) => {
            remaining.set(newValue);
            status.set(statusInProgress);
            handleStatusUpdate()
            if (newValue === 0) {
                status.set(statusComplete);
                handleStatusUpdate()
                completed.set(true)
                markExportCompleted();
            }
        })
    })

</script>
        <td class="status-table-td status-col-1">{formatName(name)}</td>
        <td class="status-table-td status-col-2">
            <span class={pillLabel.class}>{pillLabel.value}</span>
        </td>
        <td class="status-table-td status-col-3">{($remaining == -1 || $remaining == 0) ? "-" : $remaining}</td>
<style>
    .pill {
        padding: 4px;
        border-radius: 16px;
        box-sizing: border-box;
    }

   .pill.complete {
        background-color: #E8FDF5;
        color: #005438;

    }

    .pill.in-progress {
        background-color: #E5FAFF;
        color: #0D75B5;
    }

    .pill.queued {
        background-color: #EEF1F7;
        color: #3F495A;
    }
</style>
