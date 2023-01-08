<script>
    import './common.css';
    // import SvelteTable from "svelte-table";
    import {ReadExportStatus} from "../../wailsjs/go/main/App.js"

    let data = {};
    let counter = 0;

    function pollAPI() {
        // if (data !== null && data["ExportCompleted"]) {
        //     return;
        // }

        setTimeout(function() {
            ReadExportStatus().then((result) => {
                data = result;
                console.log(data)
            })
            // Update DOM with API response

            counter++;
            pollAPI();
        }, 2000); // 2000 milliseconds = 2 seconds
    }

    pollAPI();

    // const COLUMNS = {
    //     feed_name: {
    //         key: "FeedName",
    //         title: "Feed",
    //         value: v => v.feed_name
    //     },
    //     debug_string: {
    //         key: "DebugString",
    //         title: "Status",
    //         value: v => v.debug_string
    //     }
    // }

    // $: cols = selectedCols.map(key => COLUMNS[key]);
</script>

<div class="">
    <section class="top-nav">
        <div class = "nav-left">
            <div class="h1 p-left-8">Export Status</div>
        </div>
        <div class="row">
        </div>
    </section>
</div>