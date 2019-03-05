export default function route(location) {
    const url = new URL(location)
    const { pathname, searchParams } = url
    const r = /^(?:\/ceol)?\/(|search|artists|albums|podcasts|tracks|episodes)\/?(?:\/(\d+))?$/
    const m = pathname.match(r)

    let action = 'index'

    if (!m) return ['404', action, {}]
    let [_, controller, id] = m

    if (id) searchParams.append('id', id)

    if (searchParams.has('id')) action = 'show'

    if (!controller) controller = 'home'

    return [controller, action, Object.fromEntries(searchParams)]
}
