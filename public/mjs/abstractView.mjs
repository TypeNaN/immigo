'use strict'

import router from './router.mjs'

export default class {
  constructor(params, query) {
    this.params = params
    this.query = query
    this.user = {
      email: null,
      exprire: 0,
      token: null
    }
  }

  setTitle = (title) => document.title = title
  getUserAuth = () => JSON.parse(sessionStorage.getItem('userAuth')) || { exprire: 0, token: null }
  setUserAuth = (userAuth) => {
    this.user.exprire = userAuth.exprire
    this.user.token = userAuth.token
    sessionStorage.setItem('userAuth', JSON.stringify(userAuth))
  }

  checkAuthState = async () => {
    const userAuth = this.getUserAuth()
    if (userAuth === undefined) return router('/signin')
    else if (userAuth.token === null) return router('/signin')
    else if (userAuth.exprire - (Math.floor(Date.now()) / 1000) < 0) return router('/signin')
    else await this.getRefreshToken(userAuth.token)
  }

  getRefreshToken = async (token) => {
    if (token) {
      return await fetch('/api/user/refresh', {
        method: 'GET',
        headers: { 'Authorization': 'Bearer ' + token }
      }).then(async (res) => {
        const data = JSON.parse(await res.text())
        if (res.status == 200) if (data['token']) this.setUserAuth(data)
        return
      }).catch((e) => console.error(e))
    }
    return
  }

  checkAuthVerify = async (token) => {
    if (token) {
      return await fetch('/api/user/verify', {
        method: 'GET',
        headers: { 'Authorization': 'Bearer ' + token }
      }).then(async (res) => {
        const data = JSON.parse(await res.text())
        if (res.status == 200) {
          this.user.email = data.Info.Email
          return true
        }
      }).catch((e) => console.error(e))
    }
    return false
  }

  remove_spacails = (data) => data.replace(/[\`~!@#$%^&*\(\)+=\[\]\{\};:\'\"\\|,.<>/?]/g, '')
}