'use strict'

import abstractView from './abstractView.mjs'

export default class extends abstractView {
  constructor(params, query) {
    super(params, query)
    this.setTitle('Dashboard')
  }
  
  render = async () => {
    console.log('This is Dashboard')
    const container = document.getElementById('main-container')
    container.textContent = 'Hello World!'
  }
}