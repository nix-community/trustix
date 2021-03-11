import * as ejs from 'ejs'

(() => {

  const suggestDiv: HTMLDivElement = document.querySelector("form#search-form > div#search-suggest")
  const elem: HTMLInputElement = document.querySelector("form#search-form > input[type=text]")
  const form: HTMLFormElement = document.querySelector("form#search-form")

  suggestDiv.style.setProperty("min-width", `${elem.offsetWidth}px`)

  const template = ejs.compile(`
    <ul>
      <% suggestions.forEach(function(suggestion){ %>
        <li x-data-attr="<%= suggestion %>"><a href="#"><%= suggestion %></a></li>
      <% }); %>
    </ul>
  `)

  elem.addEventListener('blur', () => {
    // Prevent disappearing before we have a chance to react to the click event
    setTimeout(() => {
      suggestDiv.hidden = true
    }, 200);
  })

  let searchTimeout = setTimeout((): void => null, 0)

  elem.addEventListener("keyup", (): void => {
    const value = elem.value
    if (value.length < 3) {
      suggestDiv.hidden = true
      clearTimeout(searchTimeout)
      return
    }

    clearTimeout(searchTimeout)
    searchTimeout = setTimeout(() => {

      fetch(`/suggest/${value}`)
        .then(resp => resp.json())
        .then(suggestions => {
          if (suggestions.length <= 0) {
            return
          }

          suggestDiv.innerHTML = template({
            suggestions: suggestions,
          })

          suggestDiv.hidden = false

          suggestDiv.querySelectorAll("ul > li").forEach((suggestElem) => {
            suggestElem.addEventListener("click", () => {
              elem.value = suggestElem.getAttribute("x-data-attr")
              form.submit()
            })
          })
        })
    }, 200)
  })

})()
