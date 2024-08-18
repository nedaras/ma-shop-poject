import AddressCheckbox from './AddressCheckbox'

import './index'
import './bag'
import './sneaker'
import './address'
import './checkout'

// @ts-ignore
htmx.config.selfRequestsOnly = false 
customElements.define("address-checkbox", AddressCheckbox)
