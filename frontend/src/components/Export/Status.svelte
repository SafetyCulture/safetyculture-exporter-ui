<script lang="ts">
    import { onMount } from "svelte";
    import { EventsOn } from "../../../wailsjs/runtime/runtime.js";
    import Pill from "./Pill.svelte";

    const statusQueued = "Queued"
    const statusFailed = "Failed"
    const statusComplete = "Complete"
    const statusDownloading = "Downloading"
    const statusExporting = "Exporting"

    let {
        name = '',
        cancelled = false,
    }: {
        name?: string;
        cancelled?: boolean;
    } = $props();

    let status = $state(statusQueued);
    let counter = $state(-1);
    let counterDecremental = $state(true);
    let pillType = $state<'success' | 'info' | 'neutral' | 'error' | 'cancelled'>('neutral');

    function formatExportItemName(str: string): string {
        return str.toLowerCase().replace(/_/g, ' ').replace(/(^|\s)\S/g, (L) => L.toUpperCase());
    }

    onMount(() => {
        EventsOn("update-"+name, (newValue: Record<string, unknown>) => {
            counterDecremental = newValue['counter_decremental'] as boolean
            if (newValue['started'] === true && newValue['finished'] === false) {
                switch (newValue['stage']) {
                    case 'API_DOWNLOAD':
                        status = statusDownloading
                        break
                    case 'CSV_EXPORT':
                        status = statusExporting
                        break
                }
                counter = newValue['counter'] as number
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

<td class="w-[40%] py-4 text-left">{formatExportItemName(name)}</td>
<td class="w-[20%] py-4 text-left">
    {#if cancelled === true}
        <Pill name="Cancelled" type="cancelled"/>
    {:else}
        <Pill name={status} type={pillType}/>
    {/if}
</td>
<td class="py-4 text-left">
    {#if cancelled !== true}
        {(counter === -1 || counter === 0) ? "" : counter + " " + (counterDecremental === true ? "remaining" : "exported")}
    {/if}
</td>
