<script>
    export let label = '';
    export let onClick
    export let clazz = '';
    export let error = false
    export let type = 'purple'

    let computedClass = ''
    switch (type) {
        case 'active':
            computedClass = 'button-purple'
            break
        case 'active2':
            computedClass = 'button-white'
            break
    }
</script>

<button class="button {computedClass} {clazz}" class:button-error={error} on:click={onClick}>{label}</button>

<style>
    .button {
        padding: 14px 16px;
        border-width: 1px;
        border-radius: 12px;
        font-weight: 500;
        font-size: 1rem;
        line-height: 1.5rem;
    }

    .button:hover {
        cursor: pointer;
    }

    .button-white {
        background: #FFFFFF;
        color: #675DF4;
        border-color: #BFC5D4;
    }

    .button-purple {
        background: #675DF4;
        color: #FFFFFF;
        border-color: #FFFFFF;
    }

    .button-error {
        background-color: #9b3d41;
        color: #FFFFFF;
        border-color: #FFFFFF;
    }
</style>
