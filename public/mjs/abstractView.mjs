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
  }
}