'use strict'

import abstractView from './abstractView.mjs'

export default class extends abstractView {
  constructor(params, query) {
    super(params, query)
    this.setTitle('Signin')
  }
  
  render = () => {
    const html = `
      <div id="signin-containner">
        <div id="signin-description">โปรดลงชื่อเช้าสู่ระบบ หรือ <span><a href="/signup" />สมัครสมาชิก</a></span></div>
        <div class="inputbox">
          <input type="email" name="email" id="email" required />
          <label for="email">อีเมล</label>
        </div>
        <div class="inputbox">
          <input type="password" name="password" id="password" pattern=".{6,}" required />
          <label for="password">รหัสผ่าน</label>
        </div>
        <input type="submit" id="submit" class="btn form ok one-one" value="เข้าระบบ" />
        <center id="resmsg">&nbsp;</center>
      </div>
    `
    document.getElementById('main-container').innerHTML = html
    const resmsg = document.getElementById('resmsg')

    const signin = (e) => {
      resmsg.className = 'error'
      resmsg.textContent = 'ข้อมูลไม่ถูกต้อง'
      if (document.getElementById('email').value == '') return
      if (document.getElementById('password').value == '') return
      resmsg.className = 'sending'
      resmsg.textContent = 'กำลังเข้าสู่ระบบ'
      var target;
      if (!e) var e = window.event
      if (e.target) target = e.target
      else if (e.srcElement) target = e.srcElement
      if (target.nodeType == 3) target = target.parentNode // defeat Safari bug
      if (target) target.onclick = null

      const user = {
        Email: document.getElementById('email').value.trim(),
        Password: document.getElementById('password').value.trim()
      }

      fetch('/api/user/signin', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(user)
      }).then(async (res) => {
        if (res.status == 200) {
          const data = JSON.parse(await res.text())
          if (data['token']) {
            resmsg.className = 'success'
            resmsg.textContent = 'ยืนยันเข้าระบบ'
            sessionStorage.setItem('userAuth', JSON.stringify(data))
            return window.location.href = '/dashboard'
          }
        }
        resmsg.className = 'error'
        resmsg.textContent = 'โปรดตรวจสอบให้มั่นใจว่าพิมพ์ถูกต้อง'
        return target.onclick = signin
      }).catch((e) => {
        console.error(e)
        return target.onclick = signin
      })
    }
    const submit = document.getElementById('submit')
    submit.onclick = signin
  }
}