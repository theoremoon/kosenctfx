import Vue from "vue";
import VueRouter from "vue-router";
import Home from "../views/Home.vue";
import Login from "../views/Login.vue";
import Register from "../views/Register.vue";
import Team from "../views/Team.vue";
import ResetRequest from "../views/ResetRequest.vue";
import Reset from "../views/Reset.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "Home",
    component: Home
  },
  {
    path: "/login",
    name: "Login",
    component: Login
  },
  {
    path: "/register",
    name: "Register",
    component: Register
  },
  {
    path: "/team/:id",
    name: "Team",
    component: Team
  },
  {
    path: "/reset-request",
    name: "ResetRequest",
    component: ResetRequest
  },
  {
    path: "/reset",
    name: "Reset",
    component: Reset
  }
];

const router = new VueRouter({
  routes
});

export default router;
