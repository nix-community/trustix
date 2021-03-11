import tippy from 'tippy.js'

(() => {
  Array.from(document.querySelectorAll("*[x-data-tooltip]")).map((elem) => {
    tippy(elem, {
      content: elem.getAttribute("x-data-tooltip"),
      placement: "top-start",
      arrow: false,
    })
  })
})()
