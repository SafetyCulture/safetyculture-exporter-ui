<script lang="ts">
    let {
        onClick = null,
        value = '',
        clazz = '',
        trimLongWords = false,
        trimValue = 30,
        error = false,
    }: {
        onClick?: (() => void) | null;
        value?: string;
        clazz?: string;
        trimLongWords?: boolean;
        trimValue?: number;
        error?: boolean;
    } = $props();

    let shouldShowToolTip = $state(false);

    function getTrimmedFolderLocation(): string {
        const len = value.length
        if (value.length <= trimValue) {
            return value
        }
        return " ... ".concat(value.substring(len - trimValue, len))
    }

    function executeOnClick() {
        if (onClick === null) return
        onClick()
    }

    function handleHover() {
        shouldShowToolTip = value.length > trimValue
    }

    function handleOut() {
        shouldShowToolTip = false;
    }
</script>

<button
    class="relative flex w-full items-center justify-between rounded-lg border px-4 py-3 text-base leading-4 text-text-weak {error ? 'border-danger-dark bg-danger-dark/15' : 'border-border'} {clazz}"
    style:cursor={onClick === null ? 'default' : 'pointer'}
    onclick={executeOnClick}
    onmouseover={handleHover}
    onfocus={handleHover}
    onmouseout={handleOut}
    onblur={handleOut}
>
    <div class="break-words">{trimLongWords === true ? getTrimmedFolderLocation() : value}</div>
    <img src="../images/folder.svg" alt="folder icon">

    {#if shouldShowToolTip === true}
    <div class="absolute top-[30px] right-[-10px] z-10 w-[300px] break-words rounded-lg bg-[#bac9de] px-2.5 py-1.5">{value}</div>
    {/if}
</button>
