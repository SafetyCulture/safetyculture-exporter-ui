<script>
    import {push} from 'svelte-spa-router'
    import {ReadVersion, GetSettings, ValidateApiKey} from "../../wailsjs/go/main/App.js"
    import {shadowConfig, latestVersion, emptyStores, appUpdateAttempted} from '../lib/store.js';
    import {isNullOrEmptyObject} from "../lib/utils.js";

    
    function nextPage() {
        emptyStores();

        ReadVersion().then(result => {
            if (result != null) {
                latestVersion.set(result);
            } else {
                latestVersion.set({})
            }
        })
        if (shouldForceUpdate() === true) {
            push('/update')
            return
        } 

        GetSettings().then(result => {
            shadowConfig.set(result);

            if ("AccessToken" in result) {
                const token = result["AccessToken"]

                // check if not empty
                if (token.length === 0) {
                    push("/welcome")
                } else {
                    // check if it can auth
                    ValidateApiKey(token).then(res => {
                        if (res === false) {
                            push("/welcome")
                        } else {
                            push("/config")
                        }
                    })
                }
            } else {
                push("/welcome")
            }
        })
    }

    function shouldForceUpdate() {
        return !isNullOrEmptyObject($latestVersion)
            && $appUpdateAttempted === false
            && $latestVersion["should_update"] === true 
            && $latestVersion['current'] !== 'v0.0.0-dev'
            && $latestVersion['url'] !== ''
    }
    
    nextPage()
</script>
