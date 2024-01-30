import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import HomeView from '../views/HomeView.vue'
import ProfileView from '../views/ProfileView.vue'



const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		//{path: '/', component: LoginView},
		{path: '/login', component: LoginView},
		{path: '/', component: HomeView, meta: { requiresAuth: true }, props: { msg: 'Hi, here\'s your daily content' }},
		{path: '/link2', component: HomeView, meta: { requiresAuth: true }},
		{path: '/profile/:username', component: ProfileView, meta: { requiresAuth: true }},
		{path: '/some/:id/link', component: HomeView, meta: { requiresAuth: true }},
	]
})

export default router
/*
import {createRouter, createWebHashHistory} from 'vue-router'

import StreamView from '../views/StreamView.vue'


const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: LoginView},
		{path: '/login', component: LoginView},
		{path: '/#/', component: LoginView},
		{path: '/stream/:username', component: StreamView},
		
		// allow GET requests to /photos/... to be handled by the backend
		{path: '/photos/.*', redirect: '/'},
	]
})

export default router
*/