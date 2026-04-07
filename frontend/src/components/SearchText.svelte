<script lang="ts">
    import { Input } from "$lib/components/ui/input";

    let {
        clazz = '',
        placeholder = '',
        value = $bindable(''),
    }: {
        clazz?: string;
        placeholder?: string;
        value?: string;
    } = $props();

    function emptyFn() {
        value = '';
    }
</script>

<div class="flex w-full items-center rounded-lg border border-border px-2.5 py-2.5 text-base leading-6">
    <img src="../images/search.svg" alt="search"/>
    <Input class="border-none px-2 text-sm shadow-none outline-none focus-visible:ring-0 {clazz}" {placeholder} bind:value />
    <button class="cursor-pointer" onclick={emptyFn}>
        <img src="../images/close-circle.svg" alt="clear"/>
    </button>
</div>
