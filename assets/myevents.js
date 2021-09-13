export let createMyEvents = () => ({
  events: {},
  emit(event, ...args) {
    (this.events[event] || []).forEach(i => i(...args))
  },
  on(event, cb) {
    (this.events[event] = this.events[event] || []).push(cb)
    return () =>
      (this.events[event] = (this.events[event] || []).filter(i => i !== cb))
  }
})
