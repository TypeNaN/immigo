'use strict'

import abstractView from './abstractView.mjs'
import navBar from './navigationBar.mjs'

export default class extends abstractView {
  constructor(params, query) {
    super(params, query)
    this.setTitle('Upload images')

    this.Mime = ['jpeg', 'jpg', 'png', 'gif', 'webp']
    this.maxFileLarge = 10 // Limit 10MB
    this.maxImgUpload = 3
    this.imgData = {
      uploads: new DataTransfer(),
      previews: []
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

    const resmsg = document.getElementById('resmsg')
    const imageholder = document.getElementById('imageholder')

    container.ondragover = () => {
      imageholder.className = 'ondragover'
      return false
    }
    container.ondragend = () => {
      imageholder.className = 'ondragend'
      return false
    }
    container.ondragleave = () => {
      imageholder.className = 'ondragleave'
      return false
    }
    container.ondrop = (event) => {
      imageholder.className = 'ondrop'
      event.preventDefault && event.preventDefault()
      this.readImageFile(event.dataTransfer.files)
    }
    const imageinputfile = document.getElementById('selectfile')
    imageinputfile.onchange = (event) => {
      event.preventDefault && event.preventDefault()
      this.readImageFile(imageinputfile.files)
    }
    const submit = document.getElementById('submit')
    const cancel = document.getElementById('cancel')

    const clear = (e) => {
      this.imgData.uploads.clearData()
      this.imgData.previews = []
      document.getElementById('name').value = ''
      document.getElementById('description').value = ''
      document.getElementById('preview-upload-images').textContent = ''
    }
      
    cancel.onclick = (e) => {
      clear()
      window.location.href = '/'
    }

    const upload = async (e) => {
      await this.checkAuthState()
      resmsg.className = 'error'
      resmsg.textContent = 'ข้อมูลไม่ถูกต้อง'
      if (this.imgData.previews.length < 1) return
      if (document.getElementById('name').value == '') return
      if (document.getElementById('description').value == '') return
      resmsg.className = 'sending'
      resmsg.textContent = 'กำลังส่งข้อมูล'

      var target;
      if (!e) var e = window.event
      if (e.target) target = e.target
      else if (e.srcElement) target = e.srcElement
      if (target.nodeType == 3) target = target.parentNode // defeat Safari bug
      if (target) target.onclick = null

      const formData = new FormData()
      for (let i = 0; i < this.imgData.uploads.files.length; i++) {
        formData.append('myFile[]', this.imgData.uploads.files[i], this.imgData.uploads.files[i].name)
      }
      formData.append('name', this.remove_spacails(document.getElementById('name').value.trim()))
      formData.append('description', document.getElementById('description').value.trim())

      fetch('/api/image/upload', {
        method: 'POST',
        headers: { 'Authorization': 'Bearer ' + this.user.token },
        body: formData
      }).then(async (res) => {
        const data = JSON.parse(await res.text())
        if (res.status == 200) {
          resmsg.className = 'success'
          resmsg.textContent = 'บันทึกไฟล์เสร็จสิ้น'
          clear()
          if (data['refresh'] && data['refresh']['token']) return this.setUserAuth(data.refresh)
        }
        resmsg.className = 'error'
        resmsg.textContent = data.message
        return target.onclick = upload
      }).catch((e) => {
        console.error(e)
        resmsg.className = 'error'
        resmsg.innerHTML = 'ดูเหมือนทาง server จะเกิดปัญหา<br>โปรดลองใหม่ในภายหลัง'
        return target.onclick = upload
      })
    }

    submit.onclick = upload
  }

  readImageFile = async (files) => {
    const viewUpImages = document.getElementById('preview-upload-images')
    for (let i = 0; i < files.length; i++) {
      if (this.imgData.previews.length >= this.maxImgUpload) break
      if (files[i].size / 1000000 > this.maxFileLarge) continue
      const mime = files[i].type.split('image/')[1]
      if (!mime) continue
      if (this.Mime.indexOf(mime) < 0) continue
      const image = new Image()
      image.src = this.createObjectURL(files[i])
      image.onload = (e) => {
        let name = files[i].name
        const splited = name.split('.')
        const indexExt = splited.length - 1
        name = name.replace('.' + splited[indexExt], '').trim()
        document.getElementById('name').value = name
        const imageCrop = this.resizeCrop(e.target, 320, 320).toDataURL('image/png', 50)
        this.imgData.uploads.items.add(files[i])
        this.imgData.previews.push(imageCrop)
        const im = document.createElement('div')
        im.setAttribute('class', 'preview-upload-images')
        im.style.backgroundImage = 'url("' + imageCrop + '")'
        viewUpImages.appendChild(im)
        const x = document.createElement('div')
        x.setAttribute('class', 'delete-upload-image')
        x.innerText = 'X'
        im.appendChild(x)
        x.onclick = () => {
          let index = this.imgData.previews.indexOf(imageCrop)
          if (index > -1) {
            this.imgData.uploads.items.remove(files[i])
            this.imgData.previews.splice(index, 1)
            im.parentNode.removeChild(im)
            document.getElementById('name').value = ''
          }
        }
        window.URL.revokeObjectURL(image.src)
      }
    }
  }

  createObjectURL = (i) => (window.URL || window.webkitURL || window.mozURL || window.msURL).createObjectURL(i)

  resizeCrop = (src, width, height) => {
    const crop = width == 0 || height == 0
    if (src.width <= width && height == 0) {
      width = src.width
      height = src.height
    }
    if (src.width > width && height == 0) height = src.height * (width / src.width)
    const xscale = width / src.width
    const yscale = height / src.height
    const scale = crop ? Math.min(xscale, yscale) : Math.max(xscale, yscale);
    const canvas = document.createElement('canvas')
    canvas.width = width ? width : Math.round(src.width * scale)
    canvas.height = height ? height : Math.round(src.height * scale)
    canvas.getContext("2d").scale(scale, scale)
    canvas.getContext("2d").drawImage(src,
      - ((src.width * .5) + ((canvas.width * -.5) / scale)),
      - ((src.height * .5) + ((canvas.height * -.5) / scale))
    )
    return canvas
  }
}