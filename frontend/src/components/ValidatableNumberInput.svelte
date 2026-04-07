<script lang="ts">
    import { Input } from "$lib/components/ui/input";

    let {
        placeholder = '',
        error = false,
        errorMsg = '',
        value = $bindable(''),
        maxlength = 5,
    }: {
        placeholder?: string;
        error?: boolean;
        errorMsg?: string;
        value?: string;
        maxlength?: number;
    } = $props();

    function handleInput(e: Event) {
        const target = e.target as HTMLInputElement;
        value = target.value.replace(/[^0-9]/g, '');
    }
</script>

<div>
    <Input
        type="text"
        {placeholder}
        pattern="^([0-9]{'{'}1,5{'}'})$"
        {maxlength}
        bind:value
        oninput={handleInput}
        class="w-[95%] rounded-lg border px-4 py-2.5 text-base leading-6 {error ? 'border-danger-dark placeholder:text-danger-dark' : 'border-border'}"
    />

    {#if error && errorMsg}
        <div class="mt-1 text-xs text-danger-dark">{errorMsg}</div>
    {/if}
</div>
