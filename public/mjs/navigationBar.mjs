'use strict'

export default class {
  constructor() {
    this.root = document.createElement('ul')
    this.root.id = 'main-nav'
    this.dashboard = this.add('Dashboard', '/dashboard')
    this.uploads = this.add('Uploads', '/upimgs')
    this.signout = this.add('Signout', '/signout')
  }

  add = (name, href) => {
    const menu = document.createElement('li')
    const a = document.createElement('a')
    a.textContent = name
    a.href = href
    menu.appendChild(a)
    this.root.appendChild(menu)
  }
}