import { StatusNotFound, StatusTooManyRequests } from "./http"

document.body.addEventListener('htmx:beforeSwap', (e) => {
  if (!e.detail.isError) return
  if (!e.detail.pathInfo.requestPath.startsWith("/address/") && !e.detail.pathInfo.requestPath.startsWith("/htmx/address/")) return

  switch (e.detail.xhr.status) {
  case StatusNotFound:
  case StatusTooManyRequests:
    break
  default:
    return
  }

  e.detail.isError = false
  e.detail.shouldSwap = true
})
