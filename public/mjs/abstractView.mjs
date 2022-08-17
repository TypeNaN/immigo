'use strict'

export default class {
  constructor(params, query) {
    this.params = params
    this.query = query
  }

  setTitle = (title) => {
    document.title = title
  }

  checkAuthState = () => {
    const userAuth = JSON.parse(sessionStorage.getItem('userAuth')) || { exprire: 0, token: null }
    if (userAuth == undefined) window.location.href = '/signin'
    else if (userAuth.token == null) window.location.href = '/signin'
    else if (userAuth.exprire - (Math.floor(Date.now()) / 1000) < 0) window.location.href = '/signin'
    else this.getRefreshToken(userAuth)
  }

  getRefreshToken = (userAuth) => {
    if (userAuth['token']) {
      fetch('/api/user/refresh', {
        method: 'GET',
        headers: { 'Authorization': 'Bearer ' + userAuth.token }
      }).then(async (res) => {
        if (res.status == 200) {
          const data = JSON.parse(await res.text())
          if (data.token) return sessionStorage.setItem('userAuth', JSON.stringify(data))
        }
        console.error(e)
      }).catch((e) => console.error(e))
    }
  }
}