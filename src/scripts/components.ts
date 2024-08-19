import { StatusOK } from "./http"

export class AddressCheckbox extends HTMLElement {

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
    const defaultAddress = localStorage.getItem('default_address')
    if (defaultAddress == null) return false

    return defaultAddress === this.getAttribute('address')
  }

  isChecked() {
    return (this.querySelector('input') as HTMLInputElement).checked
  }

}

export class AddressRadio extends HTMLElement {

  constructor() {
    super()

    const checked = () => {
      const defaultAddress = localStorage.getItem('default_address')
      return defaultAddress ? defaultAddress === this.getAttribute('address') : false
    }

    this.innerHTML = `<label class="block">
      <input
        type="radio" 
        id="address-${this.getAttribute("address")}"
        class="peer appearance-none absolute outline-none"
        name="address"
        value="${this.getAttribute("address")}"
        ${checked() ? "checked" : ""}
      >
        ${this.innerHTML}
    </label>`
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
