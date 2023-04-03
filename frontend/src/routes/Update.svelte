<script>
    import {latestVersion} from '../lib/store.js';
    import {TriggerUpdate} from "../../wailsjs/go/main/App.js";

    // TriggerUpdate("https://github.com/MickStanciu/SC-Mock/releases/download/v1.0.0/exporter-darwin-amd64.zip").then(result => {
    //     console.log(result)
    // })
</script>
<div style="background-color: #9b3d41">
    FORCE UPDATE {JSON.stringify($latestVersion)}
</div>