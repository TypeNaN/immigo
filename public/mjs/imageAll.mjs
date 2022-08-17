export default class {
  constructor(container, data) {
    
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
    }
  }
}

