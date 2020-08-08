const eventHub = {
  // eslint-disable-next-line no-unused-vars
  install: function(Vue, options) {
    Vue.prototype.$eventHub = new Vue();
  }
};
export default eventHub;
