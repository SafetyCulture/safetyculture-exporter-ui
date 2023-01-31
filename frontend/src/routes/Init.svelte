<script>
    import {push} from 'svelte-spa-router'
    import {ReadVersion, GetSettings, ValidateApiKey} from "../../wailsjs/go/main/App.js"
    import {shadowConfig, latestVersion, emptyStores} from '../lib/store.js';

    emptyStores();

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

    ReadVersion().then(result => {
        if (result != null) {
            latestVersion.set(result);
        } else {
            latestVersion.set({})
        }
        console.debug(result)
    })

</script>

