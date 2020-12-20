import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import eventHub from "./eventHub";
import store from "./store";
import "./assets/tailwind.css";
import { library } from "@fortawesome/fontawesome-svg-core";
import { faFlag } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import countryFlag from "country-code-emoji";
import VueApollo from "vue-apollo";
import { apolloClient } from "./apollo";

library.add(faFlag);

Vue.component("font-awesome-icon", FontAwesomeIcon);

Vue.config.productionTip = false;
Vue.use(eventHub);
Vue.use(store);
Vue.use(VueApollo);

const apolloProvider = new VueApollo({
  defaultClient: apolloClient
});

Vue.filter("countryFlag", function(countryCode) {
  try {
    return countryFlag(countryCode);
  } catch (e) {
    return "";
  }
});

new Vue({
  router,
  render: h => h(App),
  apolloProvider
}).$mount("#app");
