document.body.addEventListener('htmx:beforeSwap', (e) => {
  if (e.detail.isError) return
  if (!e.detail.pathInfo.requestPath.startsWith('/htmx/product')) return
  if (e.detail.serverResponse !== "") return

  const placeholder = document.getElementById('placeholder')
  if (!placeholder) return

  clear(placeholder)
  e.detail.shouldSwap = false
})

document.body.addEventListener('htmx:beforeSwap', (e) => {
  if (e.detail.isError) return
  if (e.detail.pathInfo.requestPath !== '/bag') return

  const placeholder = document.getElementById('placeholder')
  if (!placeholder) return

  const div = placeholder.querySelector('div')
  if (!div) return

  placeholder.style.transitionDuration = '0ms'
  placeholder.style.transform = ''

  div.innerHTML = ''
})

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

      placeholder.style.transitionDuration = ''
      placeholder.style.transform = 'translateY(0%)'
      placeholder.onclick = clickHandler(
        setTimeout(() => clear(placeholder), 2000)
      )

      close && (close.onclick = () => clear(placeholder))
    }
  }
})

function clickHandler(timeout: number) {
  return () => {
    clearTimeout(timeout)
  }
}
// todo: we would need whole day to sit make that scroll fell right

function clear(placeholder: HTMLElement) {
  placeholder.style.transform = ''
  setTimeout(() => {
    const div = placeholder.querySelector('div')
    div && (div.innerHTML = '')
  }, 200)
}
