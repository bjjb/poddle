import { expect } from './tests.js'
import route from './router.js'

describe('route', () => {
  Object.entries({
    'https://example.com/': ['home', 'index', {}],
    'https://example.com/ceol/': ['home', 'index', {}],
    'https://example.com/albums': ['albums', 'index', {}],
    'https://example.com/ceol/albums': ['albums', 'index', {}],
    'https://example.com/albums/': ['albums', 'index', {}],
    'https://example.com/albums/123': ['albums', 'show', { id: '123' }],
    'https://example.com/albums/123?foo=bar':
      ['albums', 'show', { id: '123', foo: 'bar' }],
    'https://example.com/artists': ['artists', 'index', {}],
    'https://example.com/artists/': ['artists', 'index', {}],
    'https://example.com/artists/123': ['artists', 'show', { id: '123' }],
    'https://example.com/ceol/artists/123': ['artists', 'show', { id: '123' }],
    'https://example.com/artists/123?foo=bar':
      ['artists', 'show', { id: '123', foo: 'bar' }],
    'https://example.com/tracks': ['tracks', 'index', {}],
    'https://example.com/tracks/': ['tracks', 'index', {}],
    'https://example.com/tracks/123': ['tracks', 'show', { id: '123' }],
    'https://example.com/tracks/123?foo=bar':
      ['tracks', 'show', { id: '123', foo: 'bar' }],
    'https://example.com/search?a=x&b=99':
      ['search', 'index', { a: 'x', b: '99' }],
    'https://example.com/no/such/path': ['404', 'index', {}]
  }).forEach(([r, [x, y, z]]) => {
    describe(r, () => {
      const [c, a, p] = route(r)
      it(`controller = ${x}`, () => { expect(x).to.eq(x) })
      it(`action = ${y}`, () => { expect(x).to.eq(x) })
      it(`params = ${new URLSearchParams(z).toString()}`, () => {
        expect(p).to.deep.eq(z)
      })
    })
  })
})
