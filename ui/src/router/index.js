import Vue from "vue";
import VueRouter from "vue-router";
import Index from "../views/Index.vue";
import Login from "../views/Login.vue";
import Register from "../views/Register.vue";
import Team from "../views/Team.vue";
import Diploma from "../views/Diploma.vue";
import ResetRequest from "../views/ResetRequest.vue";
import Reset from "../views/Reset.vue";
import Challenges from "../views/Challenges.vue";
import Ranking from "../views/Ranking.vue";
import Admin from "../views/Admin.vue";
import AdminConfig from "../views/admin/Config.vue";
import AdminChallenges from "../views/admin/Challenges.vue";
import AdminSQL from "../views/admin/SQL.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "Index",
    component: Index
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
    path: "/team/:id/diploma",
    name: "Diploma",
    component: Diploma
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
    component: Challenges
  },
  {
    path: "/challenges/:id",
    name: "Challenges",
    component: Challenges
  },
  {
    path: "/ranking",
    name: "Ranking",
    component: Ranking
  },

  {
    path: "/admin",
    name: "Admin",
    component: Admin
  },
  {
    path: "/admin/",
    component: Admin,
    children: [
      {
        path: "config",
        component: AdminConfig
      },
      {
        path: "challenges",
        component: AdminChallenges
      },
      {
        path: "sql",
        component: AdminSQL
      }
    ]
  }
];

const router = new VueRouter({
  routes
});

export default router;
