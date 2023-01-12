<script>
    import {
        ReadVersion
    } from "../../wailsjs/go/main/App.js"

    import {push} from "svelte-spa-router";

    function gotoWelcome() {
        push("/welcome")
    }

    function gotoConfig() {
        push("/config")
    }

    function gotoTables() {
        push("/config/tables")
    }

    function gotoTemplates() {
        push("/config/templates")
    }

    function gotoStatus() {
        push("/exportStatus")
    }

    function gotoDebug() {
        push("/debug")
    }

    let version = ""
    ReadVersion().then(v => {
        version = v
    })

</script>

<div style="padding: 5px; background-color: #b25656">
    <span style="padding: 5px; color: whitesmoke">{version}</span>&nbsp;
    <button on:click={gotoWelcome}>WELCOME</button>
    <button on:click={gotoConfig}>CONFIG</button>
    <button on:click={gotoTemplates}>TEMPLATES</button>
    <button on:click={gotoTables}>TABLES</button>
    <button on:click={gotoStatus}>STATUS</button>
    <button on:click={gotoDebug}>DEBUG</button>
</div>
