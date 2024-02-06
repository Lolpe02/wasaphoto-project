import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import HomeView from '../views/HomeView.vue'
import ProfileView from '../views/ProfileView.vue'
import GodModeView from '../views/GodModeView.vue'



const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		// {path: '/', component: LoginView},
		{path: '/login', component: LoginView},
		{path: '/', component: HomeView, meta: { requiresAuth: true }},
		{path: '/GodMode', component: GodModeView},
		{path: '/profile/:username', component: ProfileView, meta: { requiresAuth: true }},
		// {path: '/some/:id/link', component: HomeView, meta: { requiresAuth: true }},
	]
})

export default router
