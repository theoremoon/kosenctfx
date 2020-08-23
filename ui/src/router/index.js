import Vue from "vue";
import VueRouter from "vue-router";
import Home from "../views/Home.vue";
import Login from "../views/Login.vue";
import Register from "../views/Register.vue";
import Team from "../views/Team.vue";
import User from "../views/User.vue";
import ResetRequest from "../views/ResetRequest.vue";
import Reset from "../views/Reset.vue";
import Challenges from "../views/Challenges.vue";
import Ranking from "../views/Ranking.vue";

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
    path: "/user/:id",
    name: "User",
    component: User
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
  },
  {
    path: "/challenges",
    name: "Challenges",
    component: Challenges
  },
  {
    path: "/ranking",
    name: "Ranking",
    component: Ranking
  }
];

const router = new VueRouter({
  routes
});

export default router;
