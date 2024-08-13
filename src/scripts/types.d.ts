type HTMXBeforeSwap = {
  xhr: XMLHttpRequest
  shouldSwap: boolean
  isError: boolean
  serverResponse: string
  elt: HTMLElement
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

type HTMXAfterOnLoad = {
  elt: HTMLElement,
  xhr: XMLHttpRequest,
  pathInfo: {
    finalRequestPath: string
    requestPath: string
    responsePath: string
  }
}

interface HTMLElementEventMap {
    'htmx:beforeSwap': CustomEvent<HTMXBeforeSwap>
    'htmx:afterSwap': CustomEvent<HTMXAfterSwap>
    'htmx:historyRestore': CustomEvent<HTMXHistoryRestore>
    'htmx:afterOnLoad': CustomEvent<HTMXAfterOnLoad>
}
