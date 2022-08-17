'use strict'

import dashboard from './dashboard.mjs'
import signin from './signin.mjs'

const pathToRegex = (path) => new RegExp(`^${path.replace(/[-\/\\^$*+?.()|[\]{}]/g, '\\$&').replace(/:(\w+)/g, '(?<$1>[^/]+)')}\/?$`)

const getParams = (match) => {
  const values = match.result.slice(1)
  if (values.length > 0) {
    const keys = Array.from(match.route.path.matchAll(/:(\w+)/g)).map((result) => result[1])
    return Object.fromEntries(keys.map((key, i) => [key, values[i]]))
  }
  return []
}

const getUrlquery = () => {
  const urlquery = {}
  window.location.href.replace(/[?&]+([^=&]+)=([^&]*)/gi, (m, key, value) => urlquery[key] = value)
  return urlquery
}

const router = () => {
  var routes = [
    { path: '/', view: dashboard },
    { path: '/dashboard', view: dashboard },
    { path: '/signin', view: signin },
  ]
  const potentialMatches = routes.map((route) => ({ route: route, result: window.location.pathname.match(pathToRegex(route.path)) }))
  let match = potentialMatches.find((potentialMatch) => potentialMatch.result !== null)
  if (!match) match = { route: routes[0], result: [window.location.pathname] }
  const view = new match.route.view(getParams(match), getUrlquery())
  return view.render()
}

export default router