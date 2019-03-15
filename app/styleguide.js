import App from './app.js'

export default class Styleguide extends App {
  constructor() {
    super()
  }
}

customElements.define('poddle-styleguide', Styleguide, { extends: 'body' })
