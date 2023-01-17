<script>
    export let placeholder = ''
    export let error = false
    export let value = '';
</script>

<input
        class="form-input" class:input-error={error}
        type="text"
        placeholder={placeholder}
        bind:value={value}
/>

<style>
    .input-error {
        border-color: #9b3d41!important;
    }

    .input-error::placeholder { /* Chrome, Firefox, Opera, Safari 10.1+ */
        color: #9b3d41;
        opacity: 1; /* Firefox */
    }

    .input-error:-ms-input-placeholder { /* Internet Explorer 10-11 */
        color: #9b3d41;
    }

    .input-error::-ms-input-placeholder { /* Microsoft Edge */
        color: #9b3d41;
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
