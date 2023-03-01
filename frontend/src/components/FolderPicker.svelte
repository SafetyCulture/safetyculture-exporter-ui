<script>
    export let label = ''
    export let onClick
    export let value = ''
    export let clazz = ''
    export let trimLongWords = false
    export let trimValue = 30

    let shouldShowToolTip = value.length > trimValue

    function getTrimmedFolderLocation() {
        const len = value.length

        // don't trim if not necessary
        if (value.length <= trimValue) {
            return value
        }
        return " ... ".concat(value.substring(len - trimValue, len))
    }

    function executeOnClick() {
        if (onClick === null) {
            return
        }

        onClick()
    }
</script>

<div class="container {clazz}" style:cursor={onClick === null ? 'default' : 'pointer'}
     on:click={executeOnClick} on:keypress={executeOnClick} >
    <div class="value">{trimLongWords === true ? getTrimmedFolderLocation() : value}</div>
    <img src="../images/folder.svg" alt="folder icon">

    {#if shouldShowToolTip === true}
    <div class="tooltip">{value}</div>
    {/if}
</div>

<style>
    .container {
        color: #545F70;
        cursor: default;
        border: 1px solid #BFC5D4;
        border-radius: 8px;

        display: flex;
        justify-content: space-between;
        position: relative;

        padding: 12px 11px 12px 16px;
        font-weight: 400;
        font-size: 1rem;
        line-height: 1rem;
        align-items: center;
    }

    .container:hover .tooltip {
        visibility: visible;
    }

    .value {
        word-wrap: break-word;
    }

    .tooltip {
        background-color: #bac9de;
        border-radius: 8px;
        position: absolute;
        z-index: 1;

        width: 300px;
        padding: 5px 10px;

        right: -10px;
        top: 30px;

        visibility: hidden;
    }
</style>
