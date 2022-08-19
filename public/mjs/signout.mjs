'use strict'

import router from './router.mjs'
import abstractView from './abstractView.mjs'

export default class extends abstractView {
  constructor(params, query) {
    super(params, query)
    this.setTitle('Signout')
  }
  
  render = () => {
    document.getElementById('main-container').innerHTML = 'Signout'
    this.user = {
      email: null,
      exprire: 0,
      token: null
    }
    sessionStorage.removeItem('userAuth')

    router('/signin')

    // ต้องส่ง request ไปบอก server ให้รู้ว่า signout ด้วย
  }
}