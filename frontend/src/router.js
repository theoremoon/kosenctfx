import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from "@/views/Home";
import Challenges from "@/views/Challenges";
import Ranking from "@/views/Ranking";
import Team from "@/views/Team";
import AdminHome from "@/views/admin/Home"

Vue.use(VueRouter)

const routes = [
    {
        name: "Home",
        path: "/",
        component: Home,
    },
    {
        name: "Challenges",
        path: "/challenges",
        component: Challenges,
    },
    {
        name: "Ranking",
        path: "/ranking",
        component: Ranking,
    },
    {
        name: "Team",
        path: "/team",
        component: Team,
    },
    {
        name: "Admin",
        path: "/admin",
        component: AdminHome
    }
]

const router = new VueRouter({routes})
export default router;