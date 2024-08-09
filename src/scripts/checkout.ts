document.body.addEventListener('htmx:beforeSwap', (e) => {
  // make it on htmx-request or smth i need to extend the indicator
  if (e.detail.isError) return
  if (e.detail.pathInfo.requestPath !== "/htmx/checkout") return

  showStripe(e.detail.serverResponse)
  e.detail.shouldSwap = false
})

async function showStripe(clientSecret: string) {
  const loadStripe = (await import('@stripe/stripe-js')).loadStripe

  const stripe = await loadStripe('pk_test_51Plzla09srk6l0GXkHcOPFrRM2z2zUMhs8VraClycfkq9oUBeCzLKehc6RfTAbaIQQ6fPKDLAAeNBNeFBizV0vlt00cDaxIkDF')
  const checkout = await stripe?.initEmbeddedCheckout({clientSecret})

  checkout?.mount('#stripe')
}

