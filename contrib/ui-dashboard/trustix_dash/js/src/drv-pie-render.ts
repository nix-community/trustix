import * as c3 from 'c3';
import 'c3/c3.css'

export function renderDrvCharts(): void {

  document.querySelectorAll("form.unreproduced").forEach((form) => {

    const chartDiv = document.createElement("div");
    chartDiv.setAttribute("class", "drv-pie-chart")

    form.insertBefore(chartDiv, form.firstChild)

    const columns: Array<[string, ...c3.PrimitiveArray]> = []

    form.querySelectorAll("div.drv-output").forEach(output => {
      const outputHash = output.getAttribute("x-data-outputhash").toString()
      const numLogs = output.querySelectorAll("ul.log-names > li.log-name").length
      columns.push([outputHash, numLogs])
    })

    c3.generate({
      bindto: chartDiv,
      padding: {
        left: 0,
      },
      data: {
        type : 'pie',
        onclick: (d) => {
          const elem: HTMLInputElement = form.querySelector(`input[type=checkbox][name='output_hash'][value='${d.id}']`)
          elem.click()
        },
        columns: columns,
        onmouseout: () => {
          document.querySelectorAll("div.drv-output.glow-border").forEach(elem => {
            elem.classList.remove("glow-border")
          })
        },
        onmouseover: (d) => {
          document.querySelectorAll("div.drv-output.glow-border").forEach(elem => {
            elem.classList.remove("glow-border")
          })
          document.querySelector(`div.drv-output[x-data-outputhash='${d.id}']`).classList.add("glow-border")
        },
      }
    });

  })
}
