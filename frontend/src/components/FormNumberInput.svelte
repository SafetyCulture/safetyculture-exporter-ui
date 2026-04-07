<script lang="ts">
    import ValidatableNumberInput from "./ValidatableNumberInput.svelte";

    let {
        label = '',
        placeholder = '',
        error = false,
        errorMsg = '',
        value = $bindable(''),
        maxlength = 5,
    }: {
        label?: string;
        placeholder?: string;
        error?: boolean;
        errorMsg?: string;
        value?: string;
        maxlength?: number;
    } = $props();
</script>

<div class="mt-2 text-xs font-medium leading-4 text-text-weak">{label}</div>
<ValidatableNumberInput {placeholder} {error} {errorMsg} {maxlength} bind:value />
