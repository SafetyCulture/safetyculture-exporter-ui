<script>
  import Popover from './Popover.svelte'
  import { dayjs } from './lib/date-utils'
  import { contextKey, setup } from './lib/context'
  import { createEventDispatcher, setContext, getContext } from 'svelte'
  import { CalendarStyle } from './calendar-style.js'
  import { createViewContext } from './lib/view-context.js'
  import Toolbar from './Toolbar.svelte'
  import View from './view/View.svelte'

  export let range = false
  export let placeholder = 'Choose Date'
  export let format = 'DD / MM / YYYY'
  export let start = dayjs().subtract(10, 'year')
  export let end = dayjs()
  export let trigger = null
  export let selectableCallback = null
  export let styling = new CalendarStyle()
  export let selected
  export let closeOnFocusLoss = true
  export let time = false
  export let morning = 7
  export let night = 19
  export let minuteStep = 5
  export let continueText = 'Continue'

  const dispatch = createEventDispatcher()

  const startContextKey = {}
  const endContextKey = {}

  const config = {
    start: dayjs(start),
    end: dayjs(end),
    isRangePicker: range,
    isTimePicker: time,
    closeOnFocusLoss,
    format,
    morning,
    night,
    selectableCallback,
    minuteStep: parseInt(minuteStep)
  }

  setContext(contextKey, setup(selected, config))
  const {
    selectedStartDate,
    selectedEndDate,
    isOpen,
    isClosing,
    highlighted,
    formatter,
    isDateChosen,
    isSelectingFirstDate
  } = getContext(contextKey)

  setContext(startContextKey, createViewContext(true, getContext(contextKey)))

  if (config.isRangePicker) {
    setContext(endContextKey, createViewContext(false, getContext(contextKey)))
  }

  let popover

  function initialisePicker () {
    highlighted.set($selectedStartDate)
    dispatch('open')
  }

  function setRangeValue () {
    selected = [ $selectedStartDate, $selectedEndDate ]
    dispatch('range-selected', {
      from: $selectedStartDate.toDate(),
      to: $selectedEndDate.toDate()
    })
  }

  function setDateValue () {
    selected = $selectedStartDate.toDate()
    dispatch('date-selected', {
      date: $selectedStartDate.toDate()
    })
  }

  function swapDatesIfRequired () {
    if (!config.isRangePicker) { return }
    const from = $selectedStartDate
    const to = $selectedEndDate
    if (to.isBefore(from)) {
      selectedStartDate.set(to)
      selectedEndDate.set(from)
    }
  }

  function addDate (e) {
    const { date } = e.detail
    if ($isSelectingFirstDate) {
      selectedStartDate.set(date)
    } else {
      selectedEndDate.set(date)
    }
    swapDatesIfRequired()
    config.isRangePicker && isSelectingFirstDate.update(v => !v)
  }

  function close () {
    swapDatesIfRequired()
    popover.close()
  }

  $: {
    if ($isDateChosen) {
      config.isRangePicker ? setRangeValue() : setDateValue()
      dispatch('change')
    }
  }

</script>

<style>
  .body {
    display: flex;
    justify-content: space-between;

    padding: 12px 12px;
    font-weight: 400;
    font-size: 1rem;
    line-height: 1rem;
    /*display: inline-block;*/
    /*text-align: center;*/
    /*overflow: visible;*/
    /*width: var(--datepicker-width);*/
    cursor: pointer;
  }

  .calendar-button {
    border: none;
    /*background-color: #FFFFFF;*/
    background-color: red;
    /*display: block;*/
    /*text-align: center;*/
    /*width: var(--button-width);*/
    /*text-decoration: none;*/
    height: 100%;
    width: 100%;

    /*color: var(--button-text-color);*/
    /*border-radius: 7px;*/
    /*box-shadow: 0px 0px 3px rgba(0, 0, 0, 0.1);*/
  }

  .button-right {
    display: flex;
    align-items: center;
  }

  *,
  *:before,
  *:after {
    box-sizing: inherit;
  }

  .contents {
    min-width: 320px;
    width: 100%;
    display: flex;
    flex-direction: column;
  }

  .view {
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  @media (min-width: 680px) {
    .view {
      flex-direction: row;
      justify-content: center;
    }
  }


</style>

<div class="body"
     class:open={$isOpen} class:closing={$isClosing}
     style={styling.toWrapperStyle()}
     on:click={() => {
       console.log('cliked');
       dispatch('open');
     }}
>
  <Popover
    {trigger}
    bind:this={popover}
    on:opened={initialisePicker}
    on:closed={() => dispatch('close')}>
    <div slot="trigger">
      <div formatted={$formatter}>
        {#if !trigger}
          <button class="calendar-button" type="button">
            {#if $isDateChosen}
              {$formatter.formattedCombined}
            {:else}
              {placeholder}
            {/if}
          </button>
        {/if}
      </div>
    </div>
    <div class="contents" slot="contents" class:is-range-picker={config.isRangePicker}>
      <div class="view">
        <View
          viewContextKey={startContextKey}
          on:chosen={addDate}
        />
        {#if config.isRangePicker}
        <View
          viewContextKey={endContextKey}
          on:chosen={addDate}
        />
        {/if}
      </div>
      <Toolbar continueText={continueText} on:close={close} />
    </div>
  </Popover>
  <div class="button-right">
    <img src="../images/calendar-icon.svg" alt="calendar icon">
  </div>
</div>
