// mb we can use css variables to get colors or sum 
// i want mutiple script arrhhh
// can i fix that title thingy here with ignoreTitle stufff uk uk
document.addEventListener('DOMContentLoaded', () => {
  document.body.addEventListener('htmx:beforeSwap', (e) => {
    if (e.detail.pathInfo.requestPath !== '/htmx/search') return
    switch (e.detail.xhr.status) {
    case 404:
    case 500:
      break
    default:
      return
    }
    e.detail.isError = false
    e.detail.shouldSwap = true
  })

  const selector = document.getElementById('country_selector') as HTMLSelectElement
  const code = document.getElementById('country_code') as HTMLDivElement

  if (!selector || !code) return;

  const defaultColor = selector.style.color
  const defaultInput = code.innerText

  function getCountryCode(country: string): string {
    switch (country) {
      case "AL": return "+355"
      case "LT": return "+370"
      case "LV": return "+371"
      case "EE": return "+372"
      case "MD": return "+373"
      case "RS": return "+381"
      case "ME": return "+382"
      case "XK": return "+383"
      case "BA": return "+387"
      case "MK": return "+389"
      case "LI": return "+423"
      default: return ""
    }
  }

  function update() {
    const disabled = selector.options[selector.selectedIndex].disabled
    const country = getCountryCode(selector.options[selector.selectedIndex].value)
    selector.style.color = disabled ? defaultColor : 'black'
    code.style.color = country ? 'black' : defaultColor
    code.innerText = country ? country : defaultInput 
  }

  selector.addEventListener('change', update)
  update()
})

// we need to support mutiple script files not only main.ts chek on vite
