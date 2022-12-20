<script>
    import './common.css';
    import {shadowConfig, templateCache} from "../lib/store.js";
    import {GetTemplates} from "../../wailsjs/go/main/App.js"
    import {push} from "svelte-spa-router";

    let searchFilter = ""

    function trim(org) {
        if (org.length > 80) {
            return org.substring(0, 80).concat(" ...")
        }
        return org
    }

    function gotoConfig() {
        push("/config")
    }

    if (Array.isArray($templateCache) && $templateCache.length === 0) {
        GetTemplates().then((result) => {
            templateCache.set(result)
        })
    }
</script>

<div class="template-filter-page p-48">
    <section class="top-nav">
        <div class="nav-left">
            <div class="h1">Export Configuration</div>
        </div>
        <div class="nav-right">
            <button class="button button-white border-round-12" on:click={gotoConfig}>Done</button>
        </div>
    </section>

    <section class="top-nav m-top-16">
        <div class="nav-left">
            <div class="h2">Choose a template</div>
        </div>
        <div class="nav-right">
            <input class="input search" placeholder="Search" bind:value={searchFilter}/>
        </div>
    </section>

    <section class="m-top-16">
        <div class="table-header text-gray-2">
            <div class="table-row flex-spaced p-horiz-8">
                <div class="nav-left">
                    <input type="checkbox" class="checkbox-purple"/>
                    <div class="m-left-32">Template</div>
                </div>
                <div class="nav-right">
                    <div>Last Modified</div>
                    <div class="m-left-8"><img src="../images/arrow-down.png" alt="down" width="12" height="12"/></div>
                </div>
            </div>
        </div>
        <div class="table-body text-gray-2 m-top-8">
        {#each $templateCache as { id, name, modified_at }, i}
            {#if (searchFilter.length > 2 && name.includes(searchFilter)) || searchFilter.length <= 2}
                <div class="table-row flex-spaced p-horiz-8">
                    <div class="nav-left">
                        <input type="checkbox" class="checkbox-purple" bind:group={$shadowConfig["Export"]["TemplateIds"]} value="{id}"/>
                        <img class="m-left-32" src="../images/template-icon.png" alt="template" width="28" height="28"/>
                        <div class="m-left-8">{trim(name)}</div>
                    </div>
                    <div class="nav-right">
                        <div>{modified_at}</div>
                    </div>
                </div>
            {/if}
        {/each}
        </div>
    </section>
</div>

<style>
    .table-header {
        background-color: #DBDFEB;
    }

    .table-header > .table-row {
        height: 36px;
    }

    .table-body > .table-row {
        height: 52px;
    }

    .table-row {
        display: flex;
        align-items: center;
    }

    input.search {

    }
</style>
