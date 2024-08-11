type Product = {
  tid: string,
  mid: string,
  amount: string
  size: string
}

// @ts-ignore
window.getProducts = () => {
  const element = document.getElementById('products') as HTMLUListElement || undefined
  if (!element) return undefined

  const products: Product[] = []
  element.querySelectorAll('section').forEach(section => {
    products.push({
      tid: section.getAttribute("aria-tid") || "",
      mid: section.getAttribute("aria-mid") || "",
      amount: section.getAttribute("aria-amount") || "",
      size: section.getAttribute("aria-size") || "",
    })
  })

  if (products.length == 0) {
    return undefined
  }

  return JSON.stringify(products) 
}
