import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import eventHub from "./eventHub";
import store from "./store";
import "./assets/tailwind.css";

Vue.config.productionTip = false;
Vue.use(eventHub);
Vue.use(store);

new Vue({
  router,
  render: h => h(App)
}).$mount("#app");
