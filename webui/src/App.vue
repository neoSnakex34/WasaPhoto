<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
export default {}
</script>

<template>

	<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
		<a class="navbar-brand d-flex justify-content-center col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WasaPHOTO</a>
		<button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
	</header>

	<div class="container-fluid">
		<div class="row">
			<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse" v-if="$route.name !== 'login'">
				<div class="position-sticky pt-3 sidebar-sticky">
					<h6 class="sidebar-heading d-flex justify-content-center align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>sections</span>
					</h6>
					<ul class="nav d-flex flex-column align-items-center">
						<li class="nav-item">
							<RouterLink to="/" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#list"/></svg>
								My Feed
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink to="/link1" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#layout"/></svg>
								My Profile	
							</RouterLink>
						</li>
					</ul>
				</div>
			</nav>
			<main :class="$route.name === 'login' ? 'w-80 mx-auto'  : 'col-md-9 ms-sm-auto col-lg-10 px-md-4'">
			<!-- <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4"> -->
				<RouterView />
			</main>
		</div>
	</div>
</template>

<style>
</style>
