type HTMXBeforeSwap = {
  xhr: XMLHttpRequest
  shouldSwap: boolean
  isError: boolean
  pathInfo: {
    finalRequestPath: string
    requestPath: string
    responsePath: string
  }
}

type HTMXBeforeSwapEvent = CustomEvent<HTMXBeforeSwap>;

interface HTMLElementEventMap {
    'htmx:beforeSwap': HTMXBeforeSwapEvent
}
