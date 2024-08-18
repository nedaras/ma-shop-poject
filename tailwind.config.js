/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["**/*.templ"],
  theme: {
    extend: {
      boxShadow: {
        'flipped': '0 -1px 3px 0px rgb(0 0 0 / 0.1), 0 -1px 2px -1px rgb(0 0 0 / 0.1)',
      },
    },  
  },
  plugins: [
    function ({ addVariant }) {
      addVariant('indicator', ['.htmx-request&', '.htmx-request &']) // todo: update to hx-indicator
    },
    function ({ addVariant }) {
      addVariant('hx-added', ['.htmx-added&', '.htmx-added&'])
    },
    function ({ addVariant }) {
      addVariant('hx-swapping', ['.htmx-swapping&', '.htmx-swapping&'])
    },
  ],
}

