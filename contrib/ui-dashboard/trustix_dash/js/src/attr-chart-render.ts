import * as c3 from 'c3'

(() => {

  const selectData = (tr: HTMLTableRowElement, type: string, extra?: string): HTMLElement => {
    let selector = `td[x-data-type='${type}']`
    if (typeof(extra) !== "undefined") {
      selector += ` ${extra}`
    }
    return tr.querySelector(selector)
  }

  const clearGlow = (rootElem: HTMLElement) => {
    rootElem.querySelectorAll(".glow-border").forEach(elem => {
      elem.classList.remove("glow-border")
    })
  }

  Array.from(document.querySelectorAll("div[x-data-attr]")).map((div) => {
    const attr = div.getAttribute("x-data-attr")

    const table = div.querySelector("table")
    const tbody = table.querySelector("tbody")

    const drvs = Array.from(tbody.querySelectorAll("tr")).map((tr) => {
      const a: HTMLElement = selectData(tr, "drv", "> a")
      return {
        drv: a.innerText,
        url: a.getAttribute("href"),
        pct_reproduced: parseFloat(selectData(tr, "pct_reproduced").innerText),
        num_outputs: parseInt(selectData(tr, "num_outputs").innerText),
        num_reproduced: parseInt(selectData(tr, "num_reproduced").innerText),
      }
    })

    const chartDiv = document.createElement("div");
    chartDiv.setAttribute("class", "index-line-chart")
    div.insertBefore(chartDiv, table)

    c3.generate({
      bindto: chartDiv,
      axis: {
        y: {
          padding: {top: 3, bottom: 0},
          max: 100,
          min: 0,
        },
      },
      data: {
        columns: [
          [attr, ...drvs.map(drv => drv.pct_reproduced)]
        ],
        onclick: (d) => {
          const drv = drvs[d.index]
          window.location.assign(drv.url)
        },
        onmouseout: () => {
          clearGlow(tbody)
        },
        onmouseover: (d) => {
          const drv = drvs[d.index]
          clearGlow(tbody)
          tbody.querySelector(`tr[x-data-drv='${drv.drv}']`).classList.add("glow-border")
        },
      },
    });
  })

})()
