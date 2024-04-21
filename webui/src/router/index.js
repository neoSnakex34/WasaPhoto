import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import Login from '../views/LoginView.vue'
import ProfileView from '../views/ProfileView.vue'
import SearchView from '../views/SearchView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/login', name: 'login', component: Login},
		{path: '/profile', name: 'profile',  component: ProfileView},
		{path: '/search', name: 'search', component: SearchView},
		{path: '/home', component: HomeView},
		{path: '/link1', component: HomeView},
		{path: '/link2', component: HomeView},
		{path: '/', redirect: "/login"},
		
		
	]
})

export default router
