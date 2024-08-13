import { StatusOK } from "./http"

export default class AddressCheckbox extends HTMLElement {

  constructor() {
    super()

    const checked = () => {
      return localStorage.getItem("default_address") == null || this.isDefaultAddress()
    }

    this.innerHTML = `<input id='address-checkbox' type=\"checkbox\"${checked() ? " checked" : ""}/>`

    document.body.removeEventListener('htmx:beforeSwap', afterOnLoad)
    document.body.addEventListener('htmx:beforeSwap', afterOnLoad)
  }
  
  isDefaultAddress() {
    const defaultAddress = localStorage.getItem("default_address")
    if (defaultAddress == null) return false

    return defaultAddress === this.getAttribute("address")
  }

  isChecked() {
    return (this.querySelector('input') as HTMLInputElement).checked
  }

}

function afterOnLoad(e: CustomEvent<HTMXBeforeSwap>) {
  if (e.detail.xhr.status != StatusOK) return

  const checkbox = e.detail.elt.querySelector('address-checkbox') as AddressCheckbox | null
  const address = checkbox?.getAttribute("address")

  if (checkbox && address) {
    if (checkbox.isChecked()) {
       localStorage.setItem("default_address", address)
       return
    }
    checkbox.isDefaultAddress() && localStorage.removeItem("default_address")
  }
}
