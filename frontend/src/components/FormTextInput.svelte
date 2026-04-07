<script lang="ts">
    import ValidatableInput from "./ValidatableInput.svelte";

    let {
        label = '',
        placeholder = '',
        error = false,
        errorMsg = '',
        value = $bindable(''),
    }: {
        label?: string;
        placeholder?: string;
        error?: boolean;
        errorMsg?: string;
        value?: string;
    } = $props();
</script>

<div class="mt-2 text-xs font-medium leading-4 text-text-weak">{label}</div>
<ValidatableInput {placeholder} {error} {errorMsg} bind:value />
