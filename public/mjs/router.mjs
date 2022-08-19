'use strict'

import dashboard from './dashboard.mjs'
import signin from './signin.mjs'
import signup from './signup.mjs'
import signout from './signout.mjs'
import upimgs from './imageUpload.mjs'
import page404 from './404.mjs'

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
  const routes = [
    { path: '/', view: dashboard },
    { path: '/dashboard', view: dashboard },
    { path: '/signin', view: signin },
    { path: '/signup', view: signup },
    { path: '/signout', view: signout },
    { path: '/upimgs', view: upimgs },
    { path: '/:another', view: page404 }
  ]
  const potentialMatches = routes.map((route) => ({ route: route, result: window.location.pathname.match(pathToRegex(route.path)) }))
  const match = potentialMatches.find((potentialMatch) => potentialMatch.result !== null) || { route: routes[0], result: [window.location.pathname] }
  const view = new match.route.view(getParams(match), getUrlquery())
  return view.render()
}

export default router