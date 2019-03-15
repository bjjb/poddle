export default class Audio extends HTMLAudioElement {
  constructor() {
    super()
  }
}
customElements.define('audio-player', Audio, { extends: 'audio' })
