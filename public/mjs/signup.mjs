'use strict'

import router from './router.mjs'
import abstractView from './abstractView.mjs'

export default class extends abstractView {
  constructor(params, query) {
    super(params, query)
    this.setTitle('Signup')
  }

  render = async () => {
    const html = `
      <div id="signup-containner">
        <div id="signup-description">โปรดกรอกข้อมูลเพื่อสมัครสมาชิก</div>
        <div class="inputbox">
          <input type="text" name="displayname" id="displayname" required />
          <label for="displayname">ชื่อแสดง</label>
        </div>
        <div class="inputbox">
          <input type="text" name="username" id="username" required />
          <label for="username">ชื่อผู้ใช้</label>
        </div>
        <div class="inputbox">
          <input type="email" name="email" id="email" required />
          <label for="email">อีเมล</label>
        </div>
        <div class="inputbox">
          <input type="password" name="password1" id="password1" pattern=".{6,}" required />
          <label for="password1">รหัสผ่าน</label>
        </div>
        <div class="inputbox">
          <input type="password" name="password2" id="password2" pattern=".{6,}" required />
          <label for="password2">ทวนรหัส</label>
        </div>
        <input type="submit" id="submit" class="btn form ok one-one" value="สมัคร" />
        <center id="resmsg">&nbsp;</center>
      </div>
    `
    document.getElementById('main-container').innerHTML = html

    const signup = (e) => {
      resmsg.className = 'error'
      resmsg.textContent = 'ข้อมูลไม่ถูกต้อง'
      if (document.getElementById('displayname').value == '') return
      if (document.getElementById('username').value == '') return
      if (document.getElementById('email').value == '') return
      if (document.getElementById('password1').value == '') return
      if (document.getElementById('password2').value != document.getElementById('password1').value) return
      resmsg.className = 'sending'
      resmsg.textContent = 'กำลังเข้าส่งข้อมูลการสมัคร'
      var target;
      if (!e) var e = window.event
      if (e.target) target = e.target
      else if (e.srcElement) target = e.srcElement
      if (target.nodeType == 3) target = target.parentNode // defeat Safari bug
      if (target) target.onclick = null

      const user = {
        DisplayName: document.getElementById('displayname').value.trim(),
        UserName: document.getElementById('username').value.trim().replace(/[^a-zA-Z0-9_]/g, ''),
        Email: document.getElementById('email').value.trim(),
        Password: document.getElementById('password2').value.trim()
      }
      fetch('/api/user/signup', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(user)
      }).then(async (res) => {
        const data = JSON.parse(await res.text())
        if (res.status == 200) {
          if (data.message == 'successfully') {
            resmsg.className = 'success'
            resmsg.textContent = 'เข้าสู่การลงชื่อเข้าใช้งาน'
            router('/signin')
          }
        }
        resmsg.className = 'error'
        // resmsg.textContent = 'โปรดตรวจสอบให้มั่นใจว่ากรอกถูกต้องครบถ้วน'
        resmsg.textContent = data.message
        return target.onclick = signup
      }).catch((e) => {
        console.error(e)
        resmsg.className = 'error'
        resmsg.innerHTML = 'ดูเหมือนทาง server จะเกิดปัญหา<br>โปรดลองใหม่ในภายหลัง'
        return target.onclick = signup
      })
    }
    const submit = document.getElementById('submit')
    submit.onclick = signup
  }
}