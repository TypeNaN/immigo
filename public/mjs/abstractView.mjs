'use strict'

export default class {
  constructor(params, query) {
    this.params = params
    this.query = query
  }

  setTitle = (title) => {
    document.title = title
  }
}