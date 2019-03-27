import Search from './search.js'

// Sets up the test page
const { expect, assert } = chai

const results = [
  {
    collectionId: 123,
  },
]

describe('search', () => {
  let search, fake
  beforeEach(() => {
    fake = sinon.fake()
    search = new Search(fake)
  })
  it('has the class "search" (for styling)', () => {
    expect(Array.from(search.classList)).to.include('search')
  })
  it('has a <form> with a (required) <input type="search">', () => {
    const form = search.querySelector('form')
    expect(form).to.be.instanceOf(HTMLFormElement)
    expect(search.form).to.eq(form)
    const input = search.querySelector('input[type="search"]')
    expect(input).to.be.instanceOf(HTMLInputElement)
    expect(form.input).to.eq(input)
    assert(form.input.required, 'input should be required!')
  })
  it('has a <ol> for results', () => {
    const results = search.querySelector('ol')
    expect(results).to.be.instanceOf(HTMLOListElement)
    expect(results).to.eq(search.results)
  })
  it('can have its term reflected by the input value', () => {
    search.term = 'foo'
    expect(search.form.input.value).to.eq('foo')
    search.form.input.value = 'bar'
    expect(search.term).to.eq('bar')
  })
  it('has a location function, which returns ?search=<term>', () => {
    search.term = 'boo'
    const l = new URL(search.location)
    expect(l.pathname).to.eq(location.pathname)
    expect(l.searchParams.get('search')).to.eq('boo')
  })
  it('handles submit events on the form', () => {
    const s = new Search()
    const event = new CustomEvent('submit')
    event.preventDefault = sinon.fake()
    const handleEvent = s.handleEvent
    s.submit = sinon.fake()
    s.form.dispatchEvent(event)
    assert(event.preventDefault.calledOnce, 'preventDefault was not called!')
    assert(s.submit.calledOnce, 'handler was not called!')
  })
})
