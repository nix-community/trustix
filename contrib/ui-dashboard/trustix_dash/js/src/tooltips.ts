import tippy from 'tippy.js'
// import 'tippy.js/dist/tippy.css'

(() => {
  document.querySelectorAll("*[x-data-tooltip]").forEach((elem) => {
    tippy(elem, {
      content: elem.getAttribute("x-data-tooltip"),
      placement: "top-start",
      arrow: false,
    })
  })
})()
