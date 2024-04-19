import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import Login from '../views/LoginView.vue'
import ProfileView from '../views/ProfileView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/login', name: 'login', component: Login},
		{path: '/profile', name: 'profile',  component: ProfileView},
		{path: '/home', component: HomeView},
		
		{path: '/link1', component: HomeView},
		{path: '/link2', component: HomeView},
		
	]
})

export default router
