import { StatusNotFound } from "./http"

document.body.addEventListener('htmx:beforeSwap', (e) => {
  if (!e.detail.pathInfo.requestPath.startsWith("/htmx/address/")) return
  if (e.detail.xhr.status != StatusNotFound) return

  e.detail.isError = false
  e.detail.shouldSwap = true

})


document.body.addEventListener('htmx:beforeSwap', (e) => {
  if (!e.detail.pathInfo.requestPath.startsWith("/address/")) return
  if (e.detail.xhr.status != StatusNotFound) return

  e.detail.isError = false
  e.detail.shouldSwap = true

})

// todo: fix this shit
document.body.addEventListener('htmx:afterSwap', (e) => {
  if (e.detail.failed) return
  if (!e.detail.pathInfo.requestPath.startsWith("/address/")) return

  const selector = document.getElementById('country_selector') as HTMLSelectElement
  if (!selector) return

  const defaultColor = selector.style.color

  selector.onchange = () => {
    const disabled = selector.options[selector.selectedIndex].disabled
    selector.style.color = disabled ? defaultColor : 'black'
  }

  const disabled = selector.options[selector.selectedIndex].disabled
  selector.style.color = disabled ? defaultColor : 'black'
})

window.addEventListener('load', () => {
  const selector = document.getElementById('country_selector') as HTMLSelectElement
  if (!selector) return

  const defaultColor = selector.style.color

  selector.onchange = () => {
    const disabled = selector.options[selector.selectedIndex].disabled
    selector.style.color = disabled ? defaultColor : 'black'
  }

  const disabled = selector.options[selector.selectedIndex].disabled
  selector.style.color = disabled ? defaultColor : 'black'
})

