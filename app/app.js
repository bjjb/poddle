import Search from './search.js'
import Feed   from './feed.js'
import Audio  from './audio.js'

export default class App extends HTMLBodyElement {
  constructor() {
    super()
    this.classList.add('poddle', 'app')

    const url = new URL(location)

    this.search = new Search(App.iTunesPodcastSearch)
    this.search.hidden = true
    this.append(this.search)
    
    this.feed = new Feed()
    this.feed.hidden = true
    this.append(this.feed)

    this.audio = new Audio()
    this.audio.hidden = true
    this.append(this.audio)

    addEventListener('popstate', this)
    addEventListener('search:result:clicked', this)
    addEventListener('feed:item:clicked', this)
    addEventListener('feed:back:clicked', this)
    addEventListener('search:done', this)
  }

  connectedCallback() {
    this.route(location)
  }

  handleEvent(event) {
    const { type, detail } = event
    switch (type) {
      case 'popstate':
        this.route(location)
        break
      case 'search:done':
        history.pushState(location.toString(), 'Search', this.search.location)
        break
      case 'search:result:clicked':
        this.search.hidden = true
        this.feed.hidden = false
        this.feed.url = detail
        history.pushState({ url: detail }, 'Podcast', this.feed.location)
        this.feed.load()
        break
      case 'feed:item:clicked':
        this.audio.hidden = false
        this.audio.url = detail
        break
      case 'feed:back:clicked':
        history.back()
        break
      default:
        console.warn('unexpected event: %o', event)
        break
    }
  }

  route(location) {
    const { searchParams } = new URL(location)

    this.search.term = searchParams.get('search') || ''
    this.search.submit(false)
    this.feed.url = searchParams.get('feed') || ''
    this.feed.load(false)

    if (searchParams.has('feed')) {
      this.feed.hidden = false
      this.search.hidden = true
    } 
    else {
      this.feed.hidden = true
      this.search.hidden = false
    }
  }

  static async iTunesPodcastSearch(term) {
    const url = new URL('https://itunes.apple.com/search')
    const parse = (result) => {
      const {
        collectionId,
        artworkUrl600,
        collectionName,
        artistName,
        feedUrl
      } = result
      return {
        id: collectionId,
        artwork: artworkUrl600,
        primary: collectionName,
        secondary: artistName,
        href: feedUrl
      }
    }
    url.searchParams.set('entity', 'podcast')
    if (!term)
      throw new Error(`invalid search term: ${term}`)
    url.searchParams.set('term', term)
    const headers = { Accept: 'application/json' }
    return await fetch(url, { headers })
                 .then(r => r.json())
                 .then(r => r.results)
                 .then(r => r.map(parse))
  }
}

customElements.define('poddle-app', App, { extends: 'body' })
