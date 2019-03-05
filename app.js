import route from './router.js'

export default class App extends HTMLElement {

  static route(location) {
    const [ controller, action, params ] = route(location)
  }

  constructor() {
    super()
    this.home = { index: () => console.log('Home') }
    this.search = { index: (p) => console.log('Search: params=%o', p) }
    this.albums = {
      index: () => console.log('Albums'),
      show: ({ id }) => console.log('Albums, id=%d', id)
    },
    this.artists = {
      index: () => console.log('Artists'),
      show: ({ id }) => console.log('Artists, id=%d', id)
    }
    this[404] = { index: () => { console.log('Not found!') }}
  }

  connectedCallback() {
    this.route(location)
  }

  route(location) {
    const [ c, a, params ] = route(location)
    const controller = this[c]
    if (!controller) throw new Error(`no such controller: ${c}`)
    const action = controller[a]
    if (!action) throw new Error(`no such action: ${c}#${a}`)
    action(params)
  }
}

customElements.define('ceol-app', App)
