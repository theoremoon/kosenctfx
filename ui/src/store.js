const store = {
  // eslint-disable-next-line no-unused-vars
  install: function(Vue, options) {
    Vue.prototype.$store = new Vue({
      data() {
        return {
          teamname: null,
          teamid: null,
          ctfStart: null,
          ctfEnd: null,
          ctfName: null,
          challenges: [],
          ranking: []
        };
      }
    });
  }
};
export default store;
