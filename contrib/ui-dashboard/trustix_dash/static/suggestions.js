(() => {
  const suggestDiv = document.querySelector("form#search-form > div#search-suggest")
  const elem = document.querySelector("form#search-form > input[type=text]")
  const form = document.querySelector("form#search-form")

  suggestDiv.style["min-width"] = `${elem.offsetWidth}px`

  const template = ejs.compile(`
    <ul>
      <% suggestions.forEach(function(suggestion){ %>
        <li x-data-attr="<%= suggestion %>"><a href="#"><%= suggestion %></a></li>
      <% }); %>
    </ul>
  `)

  elem.addEventListener('blur', (event) => {
    // Prevent disappearing before we have a chance to react to the click event
    setTimeout(() => {
      suggestDiv.hidden = true
    }, 200);
  })

  let searchTimeout = setTimeout(() => null, 0)

  elem.addEventListener("keyup", (e) => {
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

          Array.from(suggestDiv.querySelectorAll("ul > li")).map((suggestElem) => {
            suggestElem.addEventListener("click", (e) => {
              elem.value = suggestElem.getAttribute("x-data-attr")
              form.submit()
            })
          })
        })
    }, 200)
  })

})()
