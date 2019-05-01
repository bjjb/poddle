export default class Audio extends HTMLAudioElement {
  constructor() {
    super()
    this.controls = true
    this.preload = true
    this.classList.add('audio', 'player')
    Audio.events.forEach((type) => {
      this.addEventListener(type, this)
    })
  }

  set url(uri) {
    this.setAttribute('src', Audio.url(uri))
    this.play()
  }

  handleEvent(event) {
    console.log(event)
  }

  static url(uri) {
    const url = new URL(location)
    url.pathname = '/convert'
    url.search = new URLSearchParams({ uri })
    return url.toString()
  }
}

Audio.events = [
  'loadstart',
  'progress',
  'suspend',
  'abort',
  'error',
  'emptied',
  'stalled',
  'loadmetadata',
  'loadeddata',
  'canplay',
  'canplaythrough',
  'playing',
  'waiting',
  'seeking',
  'seeked',
  'ended',
  'durationchange',
  'timeupdate',
  'play',
  'pause',
  'ratechange',
  'resize',
  'volumechange'
]
customElements.define('audio-player', Audio, { extends: 'audio' })
