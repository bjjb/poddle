import { get } from './api.js'

/**
 * Searches on iTunes Podcasts for 'term', and returns the results.
 */
export async function podcasts(term) {
  if (!term.trim()) throw new Error('missing search term')
  const uri = new URL('https://itunes.apple.com/search?entity=podcast')
  uri.searchParams.set('term', term)
  return await get(uri).then(r => r.json()).then(r => r.results)
}

/**
 * Custom <input type='search'> for a search Form input, which automatically
 * enables/disables itself on online/offline.
 */
class Input extends HTMLInputElement {
  constructor(placeholder = 'Search...') {
    super()

    this.type = 'search'
    this.name = 'q'
    this.placeholder = placeholder

    this.classList.add('search', 'input')

    this.setAttribute('required', true)
    this.setAttribute('autofocus', true)

    addEventListener('online', this)
    addEventListener('offline', this)
  }

  handleEvent(event) {
    const { target, type } = event
    let { detail } = event
    switch (type) {
      case 'online':
      case 'offline':
        this.disabled = !navigator.onLine
        break
      default:
        throw new Error(`unexpected event '${type}'`)
    }
  }
}

customElements.define('search-input', Input, { extends: 'input' })

/**
 * A custom <ol> for search results which provides an interface for replacing
 * or adding results.
 */
class Results extends HTMLOListElement {
  constructor() {
    super()
    this.classList.add('search', 'results')
  }

  /**
   * Clears this list.
   */
  clear() {
    this.innerHTML = ''
  }

  /**
   * Clears this list, and adds the results (which should be an array).
   */
  set(results) {
    this.clear()
    this.add(...results)
  }

  /**
   * Adds the result(s) to this list, but constructing a Result element, 
   * putting it in a <li> and appending it to this element. Does *not* clear
   * the results first.
   */
  add(...results) {
    results.forEach((result) => {
      const a = new Result(result)
      const li = document.createElement('li')
      li.classList.add('search', 'result')
      li.append(a)
      this.append(li)
    })
  }
}

customElements.define('search-results', Results, { extends: 'ol' })

/**
 * A custom <a> to hold a search result item, which should have the following
 * attributes:
 *
 * - id::         a unique identifier
 * - artwork::    an image URL
 * - primary::    main result text
 * - secondary::  additional result text
 * - href::       the URL of the actual item's resource
 *
 * These can be styled with classes coresponding to the properties. A Result
 * captures clicks and fires a custom `search:result:clicked` event on the
 * document instead.
 */
class Result extends HTMLAnchorElement {
  constructor({ id, primary, secondary, href, artwork }) {
    super()

    this.classList.add('search', 'result')

    this.artwork = document.createElement('img')
    this.artwork.classList.add('artwork')
    this.artwork.src = artwork
    this.append(this.artwork)

    this.primary = document.createElement('span')
    this.primary.classList.add('primary')
    this.primary.innerText = primary
    this.append(this.primary)

    this.secondary = document.createElement('span')
    this.secondary.classList.add('secondary')
    this.secondary.innerText = secondary
    this.append(this.secondary)

    this.href = href
    this.id = id

    this.addEventListener('click', this)
  }

  handleEvent(event) {
    const { type } = event
    switch (type) {
      case 'click':
        event.preventDefault()
        const detail = this.href
        dispatchEvent(new CustomEvent('search:result:clicked', { detail }))
        break
      default:
        console.warn('%o: unexpected event %o', this, event)
        break
    }
  }
}

customElements.define('search-result', Result, { extends: 'a' })

/**
 * A custom search Form. Has a custom input, and listens for submits events,
 * which it captures, and issues a `search:submitted` action on the document
 * instead.
 */
class Form extends HTMLFormElement {
  constructor() {
    super()
    this.innerHTML = ''
    this.classList.add('search', 'form')
    this.input = new Input()
    this.append(this.input)
  }

  get term() {
    return this.input.value
  }

  set term(term) {
    this.input.value = term || ''
  }
}

customElements.define('search-form', Form, { extends: 'form' })

export default class Main extends HTMLElement {
  constructor(search = Main.search) {
    super()
    this.classList.add('search')
    if (typeof(search) != 'function')
      throw new Error('invalid search function')
    this.search = search
    this.form = new Form()
    this.form.addEventListener('submit', this)
    this.results = new Results()
    this.term = this.form.input.value
    this.append(this.form)
    this.append(this.results)
  }

  set term(term) {
    if (term)
      this.setAttribute('term', term)
    else
      this.removeAttribute('term')
  }

  get term() {
    return this.form.input.value
  }

  get location() {
    const url = new URL(location)
    url.searchParams.set('search', this.term)
    return url.toString()
  }

  handleEvent(event) {
    const { type, detail } = event
    switch(type) {
      case 'submit':
        event.preventDefault()
        this.submit()
        break
      default:
        console.error('unexpected event: %o', event)
    }
  }

  attributeChangedCallback(name, oldValue, newValue) {
    switch(name) {
      case 'term':
        this.form.input.value = newValue || ''
        break
      default:
        throw new Error(`unknown attribute ${name} ${oldValue}→${newValue}`)
    }
  }

  async submit(fireEvent = true) {
    if (!this.term) {
      this.results.clear()
      return
    }
    this.classList.add('working')
    this.results.set(await this.search(this.term))
    const detail = this.term
    if (fireEvent)
      dispatchEvent(new CustomEvent('search:done', { detail }))
    this.classList.remove('working')
  }

  static get observedAttributes() {
    return ['term']
  }

  static async search(term) {
    throw new Error(`you need to implement Main.search as an async function
      which takes a search term and returns an array of objects like
      [ { id, primary, secondary, artwork, href }... ]`)
  }
}

customElements.define('custom-search', Main, { extends: 'main' })
