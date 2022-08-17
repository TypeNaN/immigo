'use strict'

import abstractView from './abstractView.mjs'
import navBar from './navigationBar.mjs'

export default class extends abstractView {
  constructor(params, query) {
    super(params, query)
    this.setTitle('Upload images')

    this.imgData = {
      uploads: new DataTransfer(),
      images: []
    }
  }

  render = async () => {
    
    await this.checkAuthState()

    if (this.user.email === null) await this.checkAuthVerify(this.user.token)

    const html = `
      <form>
        <center>สูงสุด 3 ภาพ</center>
        <div id="preview-upload-images"></div>
        <div id="imageholder">
        <input id="selectfile" type="file" name="myFile" multiple required />
        </div>
        <div class="inputbox">
        <input type="text" name="name" id="name" required />
        <label for="name">ชื่อ:</label>
        </div>
        <div class="inputbox">
        <input type="text" name="description" id="description" required />
        <label for="description">คำอธิบาย:</label>
        </div>
        <input type="button" id="submit" class="btn form ok two-one" value="upload" />
        <input type="button" id="cancel" class="btn form cancel two-two" value="cancel" />
        <center id="resmsg">&nbsp;</center>
      </form>

    `
    
    const container = document.getElementById('main-container')
    container.textContent = ''

    const nav = new navBar()
    const head = document.createElement('div')
    head.id = 'main-head'
    head.appendChild(nav.root)
    container.appendChild(head)

    container.innerHTML += html
  }
}