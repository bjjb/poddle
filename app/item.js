import { convert } from './api.js'

export async function load(uri) {
  return await convert(uri)
}

export default class Item extends HTMLAudioElement {
  constructor(uri) {
    super()
    this.uri = uri
  }

  set uri(uri) {
    if (uri)
      this.setAttribute('uri', uri)
    else
      this.removeAttribute('uri')
  }

  get uri() {
    return this.getAttribute('uri')
  }

  attributeChangedCallback(name, oldValue, newValue) {
    switch (name) {
      case 'uri':
        if (newValue)
          dispatchEvent(new CustomEvent('poddle:item:clicked', { detail }))
        else {
          this.removeAttribute('src')
        }
        break
      default:
        console.warn('unknown attribute %s (%s→%s)', name, oldValue, newValue)
        break
    }
  }

  static get observedElements() {
    return ['uri']
  }

  static async load(uri) {
    return await load(uri)
  }
}

customElements.define('item-audio', Item, { extends: 'audio' })
