<script>
</script>

<div class="backfill">
    <div class="overlay">
        <slot></slot>
    </div>
</div>

<style>
    .backfill {
        background-color: rgba(0, 0, 0, 0.75);
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        width: 100%;
        height: 100%;
        z-index: 9999;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: wait;
    }

    .overlay {
        position: absolute;
        font-size: 1.5rem;
        background-color: #d08f17;
        padding: 10px 20px;
    }
</style>
