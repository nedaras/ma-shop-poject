// todo: copy goes status enums
document.body.addEventListener('htmx:beforeSwap', (e) => {
  if (!e.detail.pathInfo.requestPath.startsWith('/htmx/product')) return
  switch (e.detail.xhr.status) {
  case 404:
    break
  default:
    return
  }
  e.detail.isError = false
  e.detail.shouldSwap = true
})
