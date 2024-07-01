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

type HTMXAfterSwap = {
  xhr: XMLHttpRequest
  pathInfo: {
    finalRequestPath: string
    requestPath: string
    responsePath: string
  }
  requestConfig: {
    elt: HTMLElement
  }
}

type HTMXBeforeSwapEvent = CustomEvent<HTMXBeforeSwap>;
type HTMXAfterSwapEvent = CustomEvent<HTMXAfterSwap>;

interface HTMLElementEventMap {
    'htmx:beforeSwap': HTMXBeforeSwapEvent
    'htmx:afterSwap': HTMXAfterSwapEvent
}
