import {createRouter, createWebHistory} from 'vue-router'
import {getUser} from "@/utils";

const routes = [
    {
        path: '/',
        name: 'Home',
        component: () => import("@/views/Home.vue"),
        children: [
            {
                path: '/dashboard',
                name: 'Dashboard',
                component: () => import("@/views/dashboard/Dashboard.vue")
            },
            {
                path: '/storage',
                name: 'Storage',
                component: () => import("@/views/storage/Storage.vue")
            }
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

    if (to.name === undefined) {
        next({name: 'Login'})
    }

    if (!isAuth && to.name !== 'Login') next({ name: 'Login' })
    else if(isAuth && (to.name == 'Login')) next({ name: 'Home'})
    else next()


})

export default router