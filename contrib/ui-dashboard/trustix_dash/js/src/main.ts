import { renderAttrCharts } from "./attr-chart-render"
import { renderDrvCharts } from "./drv-pie-render"
import { initSuggestions } from "./suggestions"
import { initTooltips } from "./tooltips"

document.addEventListener("DOMContentLoaded", () => {
  renderAttrCharts()
  renderDrvCharts()
  initSuggestions()
  initTooltips()
})
