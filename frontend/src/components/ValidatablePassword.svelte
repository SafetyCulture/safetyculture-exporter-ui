<script lang="ts">
    import { Input } from "$lib/components/ui/input";

    let {
        placeholder = '',
        error = false,
        errorMsg = '',
        value = $bindable(''),
    }: {
        placeholder?: string;
        error?: boolean;
        errorMsg?: string;
        value?: string;
    } = $props();
</script>

<div>
    <Input
        type="password"
        {placeholder}
        bind:value
        class="rounded-lg border px-4 py-2.5 text-base leading-6 {error ? 'border-danger-dark placeholder:text-danger-dark' : 'border-border'}"
    />

    {#if error && errorMsg}
        <div class="mt-1 text-xs text-danger-dark">{errorMsg}</div>
    {/if}
</div>
