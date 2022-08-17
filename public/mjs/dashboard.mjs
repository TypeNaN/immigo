'use strict'

import abstractView from './abstractView.mjs'
import navBar from './navigationBar.mjs'

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

  }
}