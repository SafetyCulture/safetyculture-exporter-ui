<script>
    import ValidatablePassword from "./ValidatablePassword.svelte";

    export let label = ''
    export let placeholder = ''
    export let error = false
    export let value = '';
</script>

<div class="form-label">{label}</div>
<ValidatablePassword placeholder={placeholder} error={error} bind:value={value}/>

<style>
    .form-label {
        margin-top: 8px;
        font-weight: 500;
        font-size: 0.8rem;
        line-height: 1rem;
        color: #545F70;
    }
</style>
