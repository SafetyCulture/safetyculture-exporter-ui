<script lang="ts">
    import dayjs from 'dayjs';
    import { shadowConfig, templateCache } from "../lib/store";
    import { GetTemplates } from "../../wailsjs/go/main/App.js"
    import { push } from "@keenmate/svelte-spa-router";
    import { trim } from "../lib/utils";
    import Button from "../components/Button.svelte";
    import Overlay from "../components/Overlay.svelte";
    import StatusBar from "../components/StatusBar.svelte";
    import SearchText from "../components/SearchText.svelte";

    let searchFilter = $state("");
    let isChecked = $state(false);
    let templatesLoaded = $state(false);

    if (Array.isArray($templateCache)) {
        if ($templateCache.length === 0) {
            GetTemplates().then((result: any[]) => {
                let niceFormat = result.map(elem => {
                    return {
                        id: elem.id,
                        name: elem.name.length > 90
                            ? `${elem.name.substring(0, 90)}…`
                            : elem.name,
                        modified_at: dayjs(elem.modified_at).format('DD-MMM-YYYY')
                    }
                }).slice(0, 3000)
                templatesLoaded = true
                templateCache.set(niceFormat)
                checkAllSelected()
            })
        } else {
            templatesLoaded = true
        }
    }

    let showEmptyFilter = $derived(searchFilter.length >= 2 && ($templateCache as any[])
        .filter((v: any) => v.name.toLowerCase().includes(searchFilter.toLowerCase()))
        .length === 0);

    function checkAllSelected() {
        if (($shadowConfig as any)["Export"]["TemplateIds"].length === 0) {
            ($shadowConfig as any)["Export"]["TemplateIds"] = ($templateCache as any[]).map((e: any) => e.id)
            isChecked = true;
        }
    }

    function toggleHeaderCheckbox() {
        if (isChecked) {
            isChecked = false;
        }
    }

    function toggleBodyCheckboxes() {
        const checkboxes = document.querySelectorAll('.table-body input[type="checkbox"]') as NodeListOf<HTMLInputElement>;
        for (const checkbox of checkboxes) {
            checkbox.checked = !isChecked;
        }
    }

    function handleDone() {
        if (showEmptyFilter === true) {
            push("/config")
            return
        }

        const checkboxes = document.querySelectorAll('.table-body input[type="checkbox"]') as NodeListOf<HTMLInputElement>;
        let selectedTemplates: string[] = [];

        for (const checkbox of checkboxes) {
            if (checkbox.checked) {
                selectedTemplates.push((checkbox as any).__value)
            }
        }

        if (($templateCache as any[]).length === selectedTemplates.length) {
            ($shadowConfig as any)["Export"]["TemplateIds"] = []
        } else {
            ($shadowConfig as any)["Export"]["TemplateIds"] = selectedTemplates
        }

        push("/config")
    }

    checkAllSelected()
</script>

{#if templatesLoaded === false}
<Overlay>
    <div class="pt-8">
        <img src="../images/loading.gif" alt="loading"/>
    </div>
    <div class="pb-12 pt-2 text-center text-lg">Please wait while we processing your request ...</div>
</Overlay>
{/if}

<div class="px-8 pt-8">
    <div class="flex items-center justify-between">
        <div class="flex items-center">
            <div class="my-4 font-semibold text-2xl">Template selection</div>
        </div>
        <div class="flex items-center">
            <Button label="Done" type="active-white" onClick={handleDone}/>
        </div>
    </div>

    <div class="mt-4 flex items-center justify-between">
        <div class="flex items-center">
            <div class="font-semibold text-xl">Select templates</div>
        </div>
        <div class="flex items-center">
            <SearchText placeholder="Search" bind:value={searchFilter}/>
        </div>
    </div>

    {#if showEmptyFilter === true}
        <div class="mt-4 flex h-[60vh] items-center justify-center">
            <div class="flex flex-col leading-6">
                <img class="mx-auto" src="../images/empty_page.svg" alt="empty page"/>
                <div class="pt-12">
                    <div>Your search - <span class="font-semibold">{searchFilter.length > 15 ? searchFilter.substring(0, 15).concat(" ...") : searchFilter}</span> - did not match any template names.</div>
                    <div class="pt-2">
                        Suggestions:<br/>
                        &#x2022; Make sure all the words are spelled correctly.<br/>
                        &#x2022; Try different keywords.<br/>
                    </div>
                </div>
            </div>
        </div>
    {:else}
        <div class="mt-4 overflow-hidden">
            <div class="bg-[#DBDFEB] text-text-gray">
                <div class="flex h-9 items-center justify-between px-2">
                    <div class="flex items-center">
                        <input type="checkbox" class="size-5 accent-primary" onclick={toggleBodyCheckboxes} bind:checked={isChecked}/>
                        <div class="ml-8">Template</div>
                    </div>
                    <div class="flex items-center">
                        <div class="mr-2">Last modified</div>
                        <img src="../images/arrow-down.svg" alt="down"/>
                    </div>
                </div>
            </div>
            <div class="table-body mt-2 max-h-[calc(100vh-300px)] overflow-y-scroll text-text-gray">
                {#each $templateCache as { id, name, modified_at }, i}
                    <div class="flex h-13 items-center justify-between px-2 pr-2" class:hidden={searchFilter.length >= 2 && !(name as string).toLowerCase().includes(searchFilter.toLowerCase())}>
                        <div class="flex items-center">
                            <input type="checkbox" class="size-5 accent-primary" onclick={toggleHeaderCheckbox} bind:group={($shadowConfig as any)["Export"]["TemplateIds"]} value={id}/>
                            <img class="ml-8" src="../images/template-icon.svg" alt="template"/>
                            <div class="ml-2">{trim(name as string, 90)}</div>
                        </div>
                        <div class="flex items-center">
                            <div>{modified_at}</div>
                        </div>
                    </div>
                {/each}
            </div>
        </div>
    {/if}
</div>

<StatusBar/>
