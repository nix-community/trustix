import { renderAttrCharts } from "./attr-chart-render"
import { renderDrvCharts } from "./drv-pie-render"
import { initSuggestions } from "./suggestions"
import { initTooltips } from "./tooltips"

function init() {
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
