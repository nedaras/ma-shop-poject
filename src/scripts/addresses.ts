document.body.addEventListener('htmx:afterSwap', (e) => {
  if (e.detail.failed) return
  if (!e.detail.pathInfo.requestPath.startsWith("/htmx/address/")) return

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
