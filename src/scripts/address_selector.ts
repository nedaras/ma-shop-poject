document.body.addEventListener('htmx:afterSwap', (e) => {
  if (e.detail.failed) return
  if (e.detail.pathInfo.requestPath != '/htmx/checkout') return

  const placeholder = document.getElementById('placeholder')
  if (!placeholder) return

  console.log("we will do some crazy stuff.")

})
