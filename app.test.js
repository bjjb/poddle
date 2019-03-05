import { expect } from './tests.js'
import App from './app.js'

describe('App', () => {
  it('is cool', () => {
    expect(new App()).to.be.instanceOf(HTMLElement)
  })
})
