import Search from './search.js'
import { Main, Form, Input, Result, Results } from './search.js'

// Sets up the test page
const { expect, assert } = chai

const results = [
  {
    collectionId: 123,
  },
]

describe('search', () => {

  describe('Main', () => {
    let main

    before(() => {
      main = new Main()
      main.form.submit = sinon.fake()
    })

    it('is a HTML element', () => expect(main).to.be.instanceOf(HTMLElement))

    it('has a term', () => expect(main.term).to.eq(''))

    describe('setting the term', () => {
      before(() => main.term = 'foo')
      it('submits the form', () => {
        assert(main.form.submit.calledOnce, '.submit was not called')
      })
    })
  })
})
