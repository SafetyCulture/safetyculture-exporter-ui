<script>
    import ValidatableInput from "./ValidatableInput.svelte";

    export let label = ''
    export let placeholder = ''
    export let pattern = ''
    export let error = false
    export let errorMsg = ''
    export let value = '';
    export let maxlength;
</script>

<div class="form-label">{label}</div>
<ValidatableInput placeholder={placeholder} pattern={pattern} error={error} errorMsg={errorMsg} maxlength={maxlength} bind:value={value}/>

<style>
    .form-label {
        margin-top: 8px;
        font-weight: 500;
        font-size: 0.8rem;
        line-height: 1rem;
        color: #545F70;
    }
</style>
