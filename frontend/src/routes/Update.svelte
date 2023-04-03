<script>
    import {appUpdateAttempted, latestVersion} from '../lib/store.js';
    import Button from "../components/Button.svelte";
    import {TriggerUpdate} from "../../wailsjs/go/main/App.js";
    import {push} from "svelte-spa-router";
    import {Quit} from "../../wailsjs/runtime/runtime.js";
    
    // in order to block future attempts within the same session
    $appUpdateAttempted = true

    let isUpdating = true
    let cancelActive = false 
    let restartActive = false

    TriggerUpdate($latestVersion['url']).then(result => {
        if (result === true) {
            isUpdating = false
            cancelActive = false
            restartActive = true
        } else {
            isUpdating = false
            cancelActive = true
            restartActive = false   
        }
    }).catch(e => {
        isUpdating = false
        cancelActive = true
        restartActive = false
    })

    function cancelHandler() {
        push("/")
    }

    function restartHandler() {
        Quit()
    }
    
</script>
<div class="update-page">
    <img id="update-page-logo" class="p-top-32" src="../images/logo.svg" alt="SafetyCulture logo"/>
    <div class="h1">SafetyCulture Exporter Updater</div>
    
    <div class="middle">
        {#if isUpdating === true}
            <img id="spinner" src="../images/spinning.gif" alt="loading"/>
        {:else}
            <img id="complete" src="../images/complete.svg" alt="ok"/>
        {/if}

        <div class="h3 p-top-64">
            Please wait until we update your application from {$latestVersion['current']} to {$latestVersion['latest']}
        </div>

        <div class="p-top-32">
            <Button label="Cancel" type="active-purple" active={cancelActive} onClick={cancelHandler}/>
            <Button label="Restart" type="active-purple" active={restartActive} onClick={restartHandler}/>
        </div>  
    </div>
</div>

<style>
    .update-page {
        display: flex;
        flex-direction: column;
        align-items: center;
        
        height: 100vh;
    }

    .update-page .h1 {
        font-size: 1.8rem;
    }

    #update-page-logo {
        width: 150px;
    }
    
    img#spinner {
        width: 200px;
        height: 200px;
    }
    
    img#complete {
        width: 200px;
        height: 200px;
    }
    
    .middle {
        flex-grow: 1;

        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
    }
</style>