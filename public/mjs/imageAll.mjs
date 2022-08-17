export default class {
  constructor(obj, container, data) {
    
    this.data = data

    for (var [key, value] of Object.entries(this.data)) {
      const imcontainer = document.createElement('div')
      imcontainer.id = key
      imcontainer.className = 'imcontainer'
      container.appendChild(imcontainer)

      const imbox = document.createElement('div')
      imbox.className = 'imbox'
      imcontainer.appendChild(imbox)

      const im = document.createElement('div')
      im.className = 'im'
      im.style.backgroundImage = 'url("' + value.path + '")'
      imbox.appendChild(im)

      const imdetail = document.createElement('div')
      imdetail.className = 'imdetail'
      imcontainer.appendChild(imdetail)

      const imname = document.createElement('div')
      imname.className = 'imname'
      imname.textContent = value.name
      imname.contentEditable = true
      imdetail.appendChild(imname)

      const imddesc = document.createElement('div')
      imddesc.className = 'imddesc'
      imddesc.textContent = value.description
      imddesc.contentEditable = true
      imdetail.appendChild(imddesc)

      const imfile = document.createElement('div')
      imfile.className = 'imfile'
      imfile.textContent = value.file
      imdetail.appendChild(imfile)

      const x = document.createElement('div')
      x.setAttribute('class', 'delete-upload-image')
      x.innerText = 'X'
      im.appendChild(x)
      x.onclick = async (e) => {
        imcontainer.className = 'imcontainer remove'
        await this.waitfor(1500)
        imcontainer.parentNode.removeChild(imcontainer)

        fetch('/api/image/remove', {
          method: 'POST',
          headers: { 'Authorization': 'Bearer ' + obj.user.token },
          body: JSON.stringify({ id: imcontainer.id })
        }).then(async (res) => {
          const data = JSON.parse(await res.text())
          if (res.status == 200) {
            if (data['refresh'] && data['refresh']['token']) return obj.setUserAuth(data.refresh)
          }
        }).catch((e) => console.error(e))
      }

      const edit = (e) => {
        if (e.key == 'Enter' || e.key == 'Escape') {
          if (e.key == 'Escape') if (e.target['tempData']) e.target.textContent = e.target['tempData']
          if (e.key == 'Enter') {
            e.target['tempData'] = e.target.textContent
            
            fetch('/api/image/edit', {
              method: 'POST',
              headers: { 'Authorization': 'Bearer ' + obj.user.token },
              body: JSON.stringify({
                id: imcontainer.id,
                name: e.target.textContent,
                description: imddesc.textContent
              })
            }).then(async (res) => {
              const data = JSON.parse(await res.text())
              if (res.status == 200) {
                if (data['refresh'] && data['refresh']['token']) {
                  obj.setUserAuth(data.refresh)
                  console.log(obj.user.token);
                  return
                }
              }
            }).catch((e) => console.error(e))
          }
          e.target.blur()
        }
      }
      imname.onblur = (e) => { console.log(e) }
      imname.onfocus = (e) => e.target['tempData'] = e.target.childNodes[0].textContent
      imname.onkeydown = edit

      imddesc.onblur = (e) => { console.log(e) }
      imddesc.onfocus = (e) => e.target['tempData'] = e.target.childNodes[0].textContent
      imddesc.onkeydown = edit
    }
  }

  waitfor = (millisecond) => {
    return new Promise(resolve => {
      setTimeout(() => {
        resolve();
      }, millisecond);
    })
  }
}

