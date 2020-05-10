import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from "@/views/Home";
import Challenges from "@/views/Challenges";

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
]

const router = new VueRouter({routes})
export default router;