<script>
    import './common.css';
    import dayjs from 'dayjs';
    import {shadowConfig, templateCache} from "../lib/store.js";
    import {GetTemplates} from "../../wailsjs/go/main/App.js"
    import {push} from "svelte-spa-router";
    import {trim} from "../lib/utils.js";
    import Button from "../components/Button.svelte";
    import Overlay from "../components/Overlay.svelte";
    import StatusBar from "../components/StatusBar.svelte";

    let searchFilter = ""
    let isChecked = false
    let templatesLoaded = false

    if (Array.isArray($templateCache)) {
        if ($templateCache.length === 0) {
            GetTemplates().then((result) => {
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

    function checkAllSelected() {
        if($shadowConfig["Export"]["TemplateIds"].length === 0) {
            $shadowConfig["Export"]["TemplateIds"] = $templateCache.map(e => e.id)
            isChecked = true;
        }
    }

    function toggleHeaderCheckbox() {
        if (isChecked) {
            isChecked = false;
        }
    }

    function toggleBodyCheckboxes() {
        const checkboxes = document.querySelectorAll('.table-body input[type="checkbox"]');
        for (const checkbox of checkboxes) {
            checkbox.checked = !isChecked;
        }
    }

    function handleDone() {
        const checkboxes = document.querySelectorAll('.table-body input[type="checkbox"]');
        let selectedTemplates = [];

        for (const checkbox of checkboxes) {
            if (checkbox.checked) {
                selectedTemplates.push(checkbox.__value)
            }
        }

        if ($templateCache.length === selectedTemplates.length) {
            $shadowConfig["Export"]["TemplateIds"] = []
        } else {
            $shadowConfig["Export"]["TemplateIds"] = selectedTemplates
        }

        push("/config")
    }

    // will mark as selected all templates in the UI that were chosen
    checkAllSelected()
</script>

{#if templatesLoaded === false}
<Overlay>
    <div class="p-top-32">
        <img src="../images/loading.gif" alt="loading"/>
    </div>
    <div class="p-top-8 p-bottom-48 loading-message">Please wait while we processing your request ...</div>
</Overlay>
{/if}

<div class="template-filter-page">
    <div class="top-nav">
        <div class="nav-left">
            <div class="h1 p-left-8">Template selection</div>
        </div>
        <div class="nav-right">
            <Button label="Done" type="active-white" onClick={handleDone}/>
        </div>
    </div>

    <div class="top-nav m-top-16">
        <div class="nav-left">
            <div class="h2">Select templates</div>
        </div>
        <div class="nav-right">
            <input class="input search" placeholder="Search" bind:value={searchFilter}/>
        </div>
    </div>

    <div class="templates-body m-top-16">
        <div class="table-header text-gray-2" >
            <div class="table-row flex-spaced p-horiz-8">
                <div class="nav-left">
                    <input type="checkbox" class="checkbox-purple" on:click="{toggleBodyCheckboxes}" bind:checked={isChecked}/>
                    <div class="m-left-32">Template</div>
                </div>
                <div class="nav-right">
                    <div class="m-right-8">Last modified</div>
                    <img src="../images/arrow-down.svg" alt="down"/>
                </div>
            </div>
        </div>
        <div class="table-body text-gray-2 m-top-8">
        {#each $templateCache as { id, name, modified_at }, i}
        <div class="table-row flex-spaced p-horiz-8 m-right-8" class:hide={searchFilter.length >= 2 && !name.toLowerCase().includes(searchFilter.toLowerCase())}>
            <div class="nav-left">
                <input type="checkbox" class="checkbox-purple" on:click={toggleHeaderCheckbox} bind:group={$shadowConfig["Export"]["TemplateIds"]} value="{id}"/>
                <img class="m-left-32" src="../images/template-icon.svg" alt="template"/>
                <div class="m-left-8">{trim(name)}</div>
            </div>
            <div class="nav-right">
                <div>{modified_at}</div>
            </div>
        </div>
        {/each}
        </div>
    </div>
</div>

<StatusBar/>

<style>
    .loading-message {
        text-align: center;
        font-size: 1.2rem;
    }

    .hide {
        display: none!important;
    }

    .template-filter-page {
        padding-top: var(--main-gutter-top);
        padding-left: var(--main-gutter-left);
        padding-right: var(--main-gutter-right);
    }

    .templates-body {
        overflow: hidden;
    }

    .table-header {
        background-color: #DBDFEB;
    }

    .table-header > .table-row {
        height: 36px;
    }

    .table-body {
        max-height: calc(100vh - 300px);
        overflow-y: scroll;
    }

    .table-body > .table-row {
        height: 52px;
    }

    .table-row {
        display: flex;
        align-items: center;
    }
</style>
