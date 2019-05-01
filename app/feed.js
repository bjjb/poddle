import { get } from './api.js'

function parse(xml) {
  const doc = new DOMParser().parseFromString(xml, 'application/xml')
  const channel = doc.querySelector('rss>channel')
  let title, summary, description, owner, image, items
  title = channel.querySelector('title').textContent
  summary = channel.querySelector('summary').textContent
  description = channel.querySelector('title').textContent
  owner = channel.querySelector('owner name').textContent
  try {
    image = channel.querySelector('image url').textContent
  } catch (e) {
    image = channel.querySelector('image').getAttribute('href')
  }
  items = Array.from(channel.querySelectorAll('item')).map((item) => {
    let title, pubDate, guid, enclosure
    title = item.querySelector('title').textContent
    pubDate = new Date(item.querySelector('pubDate').textContent)
    guid = item.querySelector('guid')
    if (guid) guid = guid.textContent
    enclosure = item.querySelector('enclosure').getAttribute('url')
    return { title, pubDate, guid, enclosure }
  })
  return { title, summary, description, owner, image, items }
}

async function fetch(url) {
  if (url) {
    const result = await get(url).then(r => r.text()).then(xml => parse(xml))
    return result
  }
}

class Item extends HTMLAnchorElement {
  constructor({ title, pubDate, guid, enclosure }) {
    super()

    this.classList.add('feed', 'item')

    let e = document.createElement('span')
    e.classList.add('item', 'title')
    e.innerHTML = title
    this.append(e)

    e = document.createElement('time')
    e.classList.add('item', 'date')
    e.setAttribute('datetime', pubDate.toISOString())
    e.textContent = pubDate.toDateString()
    this.append(e)

    this.id = guid
    this.href = enclosure

    this.addEventListener('click', this)
  }

  handleEvent(event) {
    const { target, type } = event
    switch (type) {
      case 'click':
        const detail = this.href
        event.preventDefault()
        dispatchEvent(new CustomEvent('feed:item:clicked', { detail }))
        break
      default:
        console.warn('unknown event: %o', event)
        break
    }
  }
}

customElements.define('feed-item', Item, { extends: 'a' })

class Items extends HTMLOListElement {
  constructor() {
    super()
    this.classList.add('feed', 'items')
  }

  clear() {
    this.innerHTML = ''
  }

  add(...items) {
    items.forEach((item) => {
      const li = document.createElement('li')
      li.append(new Item(item))
      this.append(li)
    })
  }
}

customElements.define('feed-items', Items, { extends: 'ol' })

export class Nav extends HTMLElement {
  constructor() {
    super()
    this.back = new Nav.Button('back')
    this.append(this.back)
  }
}

Nav.Button = class extends HTMLElement {
  constructor(name) {
    super()
    this.name = name
    this.classList.add(name, 'button', 'fas', 'fa-arrow-circle-left')
    this.addEventListener('click', this)
  }

  handleEvent(event) {
    const { type, target } = event
    switch(type) {
      case 'click':
        const detail = this.name
        dispatchEvent(new CustomEvent(`feed:back:clicked`), { detail })
        break
      default:
        console.error(`unexpected event: ${type}: %o`, event)
    }
  }
}

customElements.define('feed-nav', Nav, { extends: 'nav' })
customElements.define('feed-navbutton', Nav.Button, { extends: 'i' })

export class Header extends HTMLElement {
  constructor() {
    super()
    this.classList.add('feed', 'info')

    this.name = document.createElement('h1')
    this.name.classList.add('feed', 'title')
    this.append(this.name)

    this.summary = document.createElement('p')
    this.summary.classList.add('feed', 'summary')
    this.append(this.summary)

    this.description = document.createElement('p')
    this.description.classList.add('feed', 'description')
    this.append(this.description)
    
    this.owner = document.createElement('span')
    this.owner.classList.add('feed', 'owner')
    this.append(this.owner)
  }
}

customElements.define('feed-header', Header, { extends: 'header' })

export default class Main extends HTMLElement {
  constructor(fetch = Main.fetch) {
    super()

    this.classList.add('feed')

    this.nav = new Nav()
    this.append(this.nav)

    this.header = new Header()
    this.append(this.header)

    this.list = new Items()
    this.append(this.list)

    this.fetch = fetch.bind(this)
  }

  async load(fireEvent = true) {
    const properties = await this.fetch(this.url)
    Object.assign(this, properties)
  }

  set url(url) {
    if (url)
      this.setAttribute('url', url)
    else
      this.removeAttribute('url')
  }

  get url() {
    return this.getAttribute('url')
  }

  get title() {
    return this.header.name
  }

  set title(title) {
    this.title.innerText = title
  }

  get summary() {
    return this.header.summary
  }

  set summary(summary) {
    this.summary.innerText = summary
  }

  get description() {
    return this.header.description
  }

  set description(description) {
    this.description.innerText = description
  }

  get owner() {
    return this.header.owner
  }

  set owner(owner) {
    this.owner.innerText = owner
  }

  get image() {
    return this.getAttribute('image')
  }

  set image(url) {
    if (url)
      this.setAttribute('image', url)
    else
      this.removeAttribute('image')
  }

  set items(items) {
    this.list.clear()
    this.list.add(...items)
  }

  get location() {
    const uri = new URL(location.toString())
    uri.pathname = '/'
    uri.searchParams.set('feed', this.url)
    return uri.toString()
  }

  attributeChangedCallback(name, oldValue, newValue) {
    switch (name) {
      case 'url':
        if (oldValue != newValue)
          this.load()
        break
      case 'image':
        const url = newValue
        if (url)
          // TODO find a more CSS-ey way to do this
          this.header.style.backgroundImage =
            `radial-gradient(rgba(0, 0, 0, 0.8),
             rgba(0, 0, 0, 0.5)), url("${url}")`
        else
          delete(this.header.style.backgroundImage)
        break
      default:
        console.warn('unknown attribute %s (%s→%s)', name, oldValue, newValue)
        break
    }
  }

  static get observedAttributes() {
    return ['url', 'image']
  }

  static fetch(url) {
    return fetch(url)
  }
}

customElements.define('rss-feed', Main, { extends: 'main' })
