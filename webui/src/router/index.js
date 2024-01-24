import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import HomeView from '../views/HomeView.vue'



const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: LoginView},
		{path: '/login', component: LoginView},
		{path: '/home', component: HomeView, meta: { requiresAuth: true }, props: { msg: 'Hi, here\'s your daily content' }},
		{path: '/link2', component: HomeView, meta: { requiresAuth: true }},
		{path: '/some/:id/link', component: HomeView, meta: { requiresAuth: true }},
	]
})

router.beforeEach((to, next) => {
	// Check if the route requires authentication
	if (to.matched.some((record) => record.meta.requiresAuth)) {
		// Check if the user is authenticated
		const isAuthenticated = localStorage.getItem('Authorization') != null ? true : false;

		if (!isAuthenticated) {
			// If not authenticated, redirect to the login page
			next('/');
		} else {
			// If authenticated, proceed to the requested route
			next();
		}
	} else {
		// If the route doesn't require authentication, proceed
		next();
	}
});
export default router
/*
import {createRouter, createWebHashHistory} from 'vue-router'

import StreamView from '../views/StreamView.vue'
import ProfileView from '../views/ProfileView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: LoginView},
		{path: '/login', component: LoginView},
		{path: '/#/', component: LoginView},
		{path: '/stream/:username', component: StreamView},
		{path: '/profile/:username', component: ProfileView},
		// allow GET requests to /photos/... to be handled by the backend
		{path: '/photos/.*', redirect: '/'},
	]
})

export default router
*/