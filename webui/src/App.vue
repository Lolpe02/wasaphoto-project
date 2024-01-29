<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>

	<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
		<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" >MyWASAPhoto</a>
		 <!--href="#/login"-->
		<div class="col-md-10 col-sm-0 text-light text-truncate d-inline-block">
			<h5 class="">
				{{
					$user_state.username == null ? "Not Logged In" : "Logged in as " + $user_state.username
				}}
			</h5>
		</div>
	</header>
	
	<div class="container-fluid">
		<div class="row">
			<div v-if = "this.$user_state.headers.Authorization != null" class="container-fluid">
				
				<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
						<!--<button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#Links"
							aria-controls="Links" aria-expanded="false" aria-label="Toggle navigation">
							<span class="navbar-toggler-icon"></span>
						</button>-->
						<div class="position-sticky pt-3 sidebar-sticky">
							<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
								<span>General</span>
							</h6>
							<ul class="nav flex-column">
								<li class="nav-item">
									<RouterLink to="/" class="nav-link">
										<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
										Home
									</RouterLink>
								</li>
								<li class="nav-item">
									<RouterLink to="'/profile/' + 'this.user_state.username'" class="nav-link">
										<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#layout"/></svg>
										Profile
									</RouterLink>
								</li>
								<li class="nav-item">
									<RouterLink to="/home" class="nav-link">
										<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#key"/></svg>
										Menu item 2
									</RouterLink>
								</li>
							</ul>

							<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
								<span>Secondary menu</span>
							</h6>
							<ul class="nav flex-column">
								<li class="nav-item">
									<RouterLink :to="'/some/' + 'variable_here' + '/path'" class="nav-link">
										<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#file-text"/></svg>
										Item 1
									</RouterLink>
								</li>
							</ul>
						</div>
					
				</nav>
			</div>
			<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
				<RouterView />
			</main>	
		</div> 
	</div>
</template>

<style>
</style>





