<script>
    export let label = ''
    export let placeholder = ''
    export let error = false
    export let value = '';
</script>

<div class="form-label">{label}</div>
<input
        class="form-input" class:input-error={error}
        type="text"
        placeholder={placeholder}
        bind:value={value}
/>

<style>
    .input-error {
        border-color: #a02228!important;
    }

    .form-label {
        margin-top: 8px;
        font-weight: 500;
        font-size: 0.8rem;
        line-height: 1rem;
        color: #545F70;
    }

    .form-input {
        width: 100%;
        padding: 10px 16px;
        border-radius: 8px;
        border: 1px solid #BFC5D4;
        font-size: 1rem;
        line-height: 1.5rem;
    }
</style>
