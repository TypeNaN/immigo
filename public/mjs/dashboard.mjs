'use strict'

import abstractView from './abstractView.mjs'
import navBar from './navigationBar.mjs'
import imageAll from './imageAll.mjs'

export default class extends abstractView {
  constructor(params, query) {
    super(params, query)
    this.setTitle('Dashboard')
  }
  
  render = async () => {
    console.log('This is Dashboard')

    await this.checkAuthState()

    if (this.user.email === null) await this.checkAuthVerify(this.user.token)
    
    const container = document.getElementById('main-container')
    container.textContent = ''
    
    const nav = new navBar()
    const head = document.createElement('div')
    head.id = 'main-head'
    head.appendChild(nav.root)
    container.appendChild(head)

    const body = document.createElement('div')
    body.id = 'main-body'
    const preview = document.createElement('div')
    preview.id = 'main-preview'
    body.appendChild(preview)
    container.appendChild(body)

    fetch('/api/image/all', {
      method: 'GET',
      headers: { 'Authorization': 'Bearer ' + this.user.token },
    }).then(async (res) => {
      const data = JSON.parse(await res.text())
      return data
    }).then((data) => {
      new imageAll(preview, data)
    })
  }
}