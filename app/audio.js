export default class Audio extends HTMLAudioElement {
  constructor() {
    super()
    this.controls = true
    this.preload = true
    this.classList.add('audio', 'player')
  }

  set url(uri) {
    this.setAttribute('src', Audio.url(uri))
    this.play()
  }

  static url(uri) {
    const url = new URL(location)
    url.pathname = '/convert'
    url.search = new URLSearchParams({ uri })
    return url.toString()
  }
}
customElements.define('audio-player', Audio, { extends: 'audio' })
