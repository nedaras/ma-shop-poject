type HTMXBeforeSwap = {
  xhr: XMLHttpRequest
  shouldSwap: boolean
  isError: boolean
}

type HTMXBeforeSwapEvent = CustomEvent<HTMXBeforeSwap>;

interface HTMLElementEventMap {
    'htmx:beforeSwap': HTMXBeforeSwapEvent
}
