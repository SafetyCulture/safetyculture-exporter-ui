<script>
    import {push} from 'svelte-spa-router'
    import {ReadVersion, GetSettings, ValidateApiKey} from "../../wailsjs/go/main/App.js"
    import {shadowConfig, latestVersion, emptyStores, appUpdateAttempted} from '../lib/store.js';
    import {isNullOrEmptyObject} from "../lib/utils.js";
    
    async function processVersion() {
        return await ReadVersion().then(result => {
            if (result != null) {
                latestVersion.set(result);
            } else {
                latestVersion.set({})
            }
        }).then(() => {
            return shouldForceUpdate() === true; 
        }) .catch(() => {
            return false
        }) 
    }
    
    async function processGetSettings() {
        return await GetSettings().then(result => {
            shadowConfig.set(result);
        }).then(() => {
            if ("AccessToken" in $shadowConfig) {
                const token = $shadowConfig["AccessToken"]
                if (token.length === 0) {
                    return '/welcome'
                } else {
                    // check if it can auth
                    return ValidateApiKey(token).then(res => {
                        if (res === false) {
                            return '/welcome'
                        } else {
                            return '/config'
                        }
                    })
                }
            } else {
                return '/welcome'
            }
        })
    }

    function shouldForceUpdate() {
        console.log(JSON.stringify($latestVersion))
        console.log(JSON.stringify($appUpdateAttempted))
        // console.log(!isNullOrEmptyObject($latestVersion))
        // console.log($appUpdateAttempted === false)
        // console.log($latestVersion['current'] !== 'v0.0.0-dev')
        // console.log($latestVersion['download_url'] !== '')
        
        return !isNullOrEmptyObject($latestVersion)
            && $appUpdateAttempted === false
            && $latestVersion["should_update"] === true 
            && $latestVersion['current'] !== 'v0.0.0-dev'
            && $latestVersion['download_url'] !== ''
    }

    emptyStores();
    processVersion().then(shouldUpdate => {
        if (shouldUpdate === true) {
            push('/update')
        } else {
            processGetSettings().then(page => {
                console.log(page)
                push(page)    
            })
        }
    })

</script>
