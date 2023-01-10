<script>
    import './common.css';
    import {ReadExportStatus} from "../../wailsjs/go/main/App.js"

    let data = {};
    let counter = 0;

    function pollAPI() {
        setTimeout(function() {
            ReadExportStatus().then((result) => {
                data = result;
                console.log(data)
            })
            counter++;
            pollAPI();
        }, 2000);
    }

    pollAPI();
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