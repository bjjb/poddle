// Sets up the test page
const { expect, assert } = chai

import App from './app.js'

describe('App', function() {
  let app

  before(() => app = new App())

  it('is a HTML element', () => expect(app).to.be.instanceOf(HTMLElement))

  it('has a search main element which is hidden', () => {
    expect(app.search).to.be.instanceOf(HTMLElement)
    expect(app.search.hidden).to.be.true
  })

  it('has a feed main element which is hidden', () => {
    expect(app.feed).to.be.instanceOf(HTMLElement)
    expect(app.feed.hidden).to.be.true
  })

  it('has an audio element which is hidden', () => {
    expect(app.audio).to.be.instanceOf(HTMLElement)
    expect(app.audio.hidden).to.be.true
  })

  describe('.route', () => {
    describe('http://example.com?search=x', () => {
      before(() => app.route('http://example.com?search=x'))
      it('sets the search term to x')
      it('hides the feed')
      it('shows the search')
      it('submits the search')
    })
    describe('?feed=x', () => {
      before(() => app.route('http://example.com?feed=x'))
      it('sets the feed url to x')
      it('hides the search')
      it('shows the feed')
      it('loads the feed')
    })
    describe('?feed=x&search=y', () => {
      before(() => app.route('http://example.com?feed=x&search=y'))
      it('sets the feed url to x')
      it('sets the search term to y')
      it('hides the feed')
      it('loads the feed')
      it('submits the search')
      it('shows the search')
    })
  })
})
