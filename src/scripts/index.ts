import { StatusInternalServerError, StatusNotFound } from "./http"

document.body.addEventListener('htmx:beforeSwap', (e) => {
  if (e.detail.pathInfo.requestPath !== '/htmx/search') return
  switch (e.detail.xhr.status) {
  case StatusNotFound:
  case StatusInternalServerError:
    break
  default:
    return
  }
  e.detail.isError = false
  e.detail.shouldSwap = true
})

document.body.addEventListener('htmx:afterSwap', (e) => {
  if (e.detail.pathInfo.requestPath !== '/htmx/search') return

  const search = document.getElementById('search') as HTMLInputElement
  const label = search?.nextElementSibling as HTMLLabelElement

  if (!search) return
  if (!label) return

  function listener() {
    label.remove()
    search.removeEventListener('input', listener)
  }

  search.addEventListener('input', listener)
})
