<script>
    import {shadowConfig} from '../lib/store.js';

    const d = JSON.stringify($shadowConfig)
</script>

<div class="outer-container">
    <div class="inner-container">{d}</div>
</div>

<style>
    .outer-container {
        background-color: #E9EEF6;
        display: flex;
        height: 100vh;
        padding: 20px;
        overflow: scroll;
    }

    .inner-container {
        background-color: white;
        height: 100%;
    }
</style>
