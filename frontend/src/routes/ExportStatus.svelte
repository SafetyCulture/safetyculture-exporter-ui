<script lang="ts">
    import { CancelExport, ReadExportStatus, OpenDirectory } from "../../wailsjs/go/main/App.js"
    import Status from "./../components/Export/Status.svelte";
    import { shadowConfig, exportConfig } from "../lib/store";
    import { onMount } from "svelte";
    import { EventsOn } from "../../wailsjs/runtime/runtime.js";
    import { Button } from "$lib/components/ui/button";
    import * as Dialog from "$lib/components/ui/dialog";
    import * as Table from "$lib/components/ui/table";
    import { push } from "@keenmate/svelte-spa-router";


    let feedsToExport = ($exportConfig as any)['items'] as string[];
    let exportType = ($shadowConfig as any)["Session"]["ExportType"] as string;

    let cancelTriggered = $state(false);
    let exportCompleted = $state(false);

    onMount(() => {
        EventsOn("finished-export", (newValue: boolean) => {
            if (newValue === true) exportCompleted = true
        })
    })

    ReadExportStatus();
</script>

<div class="px-8 pt-6">
    <section class="flex items-center justify-between">
        <h1 class="text-xl font-semibold">Export status</h1>
        <div class="flex items-center gap-2">
            {#if cancelTriggered}
                <img src='/images/warning-red.svg' alt="cancelled"/>
                <span class="text-sm">Export cancelled</span>
            {:else if exportCompleted}
                <img src='/images/complete.svg' alt="complete"/>
                <span class="text-sm">Export complete</span>
            {:else}
                <img src='/images/in-progress.svg' alt="in progress"/>
                <span class="text-sm">In progress</span>
            {/if}

            {#if !exportCompleted}
                <Button variant="destructive" onclick={() => { cancelTriggered = true; CancelExport(); }}>Cancel export</Button>
            {:else if !cancelTriggered}
                {#if exportType === "csv" || exportType === "reports"}
                    <Button variant="outline" onclick={() => OpenDirectory(($shadowConfig as any)["Export"]["Path"])}>Open export folder</Button>
                {/if}
                <Button onclick={() => push("/config")}>Done</Button>
            {:else}
                {#if exportType === "reports"}
                    <Button variant="outline" onclick={() => OpenDirectory(($shadowConfig as any)["Export"]["Path"])}>Open export folder</Button>
                {/if}
                <Button onclick={() => push("/config")}>Go Back</Button>
            {/if}
        </div>
    </section>

    {#if cancelTriggered && !exportCompleted}
        <Dialog.Root open={true}>
            <Dialog.Content class="sm:max-w-md">
                <Dialog.Header>
                    <Dialog.Title>Cancelling export...</Dialog.Title>
                    <Dialog.Description>Please wait while the export is being cancelled.</Dialog.Description>
                </Dialog.Header>
            </Dialog.Content>
        </Dialog.Root>
    {/if}

    <div class="mt-4 max-h-[calc(100vh-200px)] overflow-y-auto rounded-lg border border-border">
        <Table.Root>
            <Table.Header>
                <Table.Row>
                    <Table.Head class="w-[40%]">Export item</Table.Head>
                    <Table.Head class="w-[20%]">Status</Table.Head>
                    <Table.Head></Table.Head>
                </Table.Row>
            </Table.Header>
            <Table.Body>
                {#each feedsToExport as feed}
                    <Table.Row><Status name={feed} cancelled={cancelTriggered}/></Table.Row>
                {/each}
            </Table.Body>
        </Table.Root>
    </div>
</div>
