document.body.addEventListener('htmx:beforeSwap', (e) => {
  if (e.detail.isError) return
  if (!e.detail.pathInfo.requestPath.startsWith('/htmx/product')) return
  const placeholder = document.getElementById('placeholder')
  if (!placeholder) return
  if (e.detail.serverResponse != "") return


  clear(placeholder)
  e.detail.shouldSwap = false
})

// todo: add like attr wait for img to load or smth
document.body.addEventListener('htmx:afterSwap', (e) => {
  if (e.detail.failed) return
  if (!e.detail.pathInfo.requestPath.startsWith('/htmx/add_to_bag')) return

  const placeholder = document.getElementById('placeholder')
  if (!placeholder) return

  const images = placeholder.getElementsByTagName('img')
  if (images.length == 0) return

  let i = 0
  for (const image of images) {
    image.onload = () => {
      if (images.length != ++i) return
      const close = document.getElementById('placeholder-close')

      placeholder.style.transform = 'translateY(0%)'
      placeholder.onclick = clickHandler(
        setTimeout(() => clear(placeholder), 2000)
      )

      close && (close.onclick = () => clear(placeholder))
    }
  }
})

// todo: if we like clicked make it if we scroll close it
function clickHandler(timeout: number) {
  return () => {
    clearTimeout(timeout)
  }
}

function clear(placeholder: HTMLElement) {
  placeholder.style.transform = ''
  setTimeout(() => {
    const div = placeholder.querySelector('div')
    div && (div.innerHTML = '')
  }, 200)
}
