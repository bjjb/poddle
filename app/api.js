export async function get(uri) {
  const url = new URL(location)
  url.pathname = '/get'
  url.search = new URLSearchParams({ uri })
  return await fetch(url)
}

export async function convert(uri) {
  const url = new URL(location)
  url.pathname = '/convert'
  url.search = new URLSearchParams({ uri })
  return await fetch(url)
}
