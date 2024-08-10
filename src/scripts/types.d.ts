type HTMXBeforeSwap = {
  xhr: XMLHttpRequest
  shouldSwap: boolean
  isError: boolean
  serverResponse: string
  pathInfo: {
    finalRequestPath: string
    requestPath: string
    responsePath: string
  }
}

type HTMXAfterSwap = {
  xhr: XMLHttpRequest
  failed: boolean
  pathInfo: {
    finalRequestPath: string
    requestPath: string
    responsePath: string
  }
  requestConfig: {
    elt: HTMLElement
  }
}

type HTMXHistoryRestore = {
  path: string
}

type HTMXBeforeSwapEvent = CustomEvent<HTMXBeforeSwap>
type HTMXAfterSwapEvent = CustomEvent<HTMXAfterSwap>
type HTMXHistoryRestoreEvent = CustomEvent<HTMXHistoryRestore>

interface HTMLElementEventMap {
    'htmx:beforeSwap': HTMXBeforeSwapEvent
    'htmx:afterSwap': HTMXAfterSwapEvent
    'htmx:historyRestore': HTMXHistoryRestoreEvent 
}
