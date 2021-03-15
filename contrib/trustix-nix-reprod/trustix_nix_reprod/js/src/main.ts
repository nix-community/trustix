import { renderAttrCharts } from "./attr-chart-render"
import { renderDrvCharts } from "./drv-pie-render"
import { initSuggestions } from "./suggestions"
import { initTooltips } from "./tooltips"

function init() {
  // While it may look like we should load this CSS from base.jinja2
  // it's entirely dependent on the javascript loading.
  // Give noscript users a better experience by completely omiting the css
  // loading
  const link = document.createElement("link")
  link.href = "/js/bundle.css"
  link.type = "text/css"
  link.rel = "stylesheet"
  document.getElementsByTagName("head")[0].appendChild(link)

  renderAttrCharts()
  renderDrvCharts()
  initSuggestions()
  initTooltips()
}

if (document.readyState === "complete") {
  init()
} else {
  document.addEventListener("DOMContentLoaded", init)
}
