import { AddressCheckbox, AddressRadio } from './components'

import './index'
import './bag'
import './sneaker'
import './address'
import './checkout'

// @ts-ignore
htmx.config.selfRequestsOnly = false 

customElements.define("address-checkbox", AddressCheckbox)
customElements.define("address-radio", AddressRadio)
