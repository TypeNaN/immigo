'use strict'

import router from './router.mjs'

export default class {
  constructor(container) {
    this.root = document.createElement('ul')
    this.root.id = 'main-nav'
    this.add('Dashboard', '/dashboard')
    this.add('Uploads', '/upimgs')
    this.add('Signout', '/signout')
  }
  

  add = (name, href) => {
    const menu = document.createElement('li')
    menu.className = 'anchor'
    menu.textContent = name
    menu.onclick = (e) => router(href)
    this.root.appendChild(menu)
  }
}