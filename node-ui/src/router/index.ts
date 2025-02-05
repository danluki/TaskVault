import {createRouter, createWebHistory} from 'vue-router'
import {getUser} from "../utils";

const routes = [
    {
        path: '/',
        name: 'Dashboard',
        component: () => import("@/views/Home.vue"),
        children: [

        ],
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import("@/views/auth/Login.vue")
    }
]

const router = createRouter({
    history: createWebHistory("/"),
    routes
})

router.beforeEach((to, _, next) => {
    let isAuth = false;
    let user = getUser();

    if (user) {
        isAuth = true
    }

    if (!isAuth && to.name !== 'Login') next({ name: 'Login' })
    else if(isAuth && (to.name == 'Login')) next({ name: 'Dashboard'})
    else next()
})

export default router