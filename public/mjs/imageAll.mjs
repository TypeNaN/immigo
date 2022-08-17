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

      const edit = (e) => {
        if (e.key == 'Enter' || e.key == 'Escape') {
          if (e.key == 'Escape') if (e.target['tempData']) e.target.textContent = e.target['tempData']
          if (e.key == 'Enter') {
            e.target['tempData'] = e.target.textContent

            console.log(obj.user.token);
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
            }).then((data) => {
              console.log('images', data)
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
}

