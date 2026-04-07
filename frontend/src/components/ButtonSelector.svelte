<script lang="ts">
    let {
        label = '',
        title = '',
        onClick,
    }: {
        label?: string;
        title?: string;
        onClick: () => void;
    } = $props();
</script>

{#if label !== ''}
    <div class="mt-4 text-sm font-medium leading-4">{label}</div>
{/if}
<button class="mt-1 flex w-full cursor-pointer items-center justify-between rounded-lg border border-border px-4 py-3 text-base leading-4" onclick={onClick}>
    <div>{title}</div>
    <div class="flex items-center">
        <img src="../images/arrow-right-compact.svg" alt="right arrow icon">
    </div>
</button>
