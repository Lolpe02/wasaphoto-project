<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>

<script>
export default {
	
	computed: {
        user_State() {
            return this.$user_state || {
                headers: {
                    Authorization: null,
                    accept: "application/json",
                },
                username: null,
                viewing: null
            };
        }
    },
	methods: {
		Logout() {
			console.log("Logging out", this.user_State, this.$user_state);
			this.user_State.username = null;
			this.user_State.headers.Authorization = null;
			if (this.$user_state) {
				this.$user_state.username = null;
				this.$user_state.headers.Authorization = null;
			}
			console.log("Logging out", this.user_State, this.$user_state);
			// this.$router.push("/");
		},
	},
};

</script>

<template>
	<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">

		<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6">MyWASAPhoto</a>
		<!--  <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#sidebarMenu"
			aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>   href="#/login"-->
		<div class="col-md-10 col-sm-0 text-light text-truncate d-inline-block">
			<h5 class="">
				{{ user_State.username == null ? "Not Logged In" : "Logged in as " + user_State.username }}
			</h5>
		</div>
	</header>

	<div class="container-fluid">
		<div class="row">
			<div v-if="user_State.headers.Authorization != null"
				class="container-fluid">

				<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
					<!---->

					<div class="position-sticky pt-3 sidebar-sticky">
						<h6
							class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
							<span>General</span>
						</h6>
						<ul class="nav flex-column">
							<li class="nav-item">
								<RouterLink to="/" class="nav-link">
									<svg class="feather">
										<use href="/feather-sprite-v4.29.0.svg#home" />
									</svg>
									Home
								</RouterLink>
							</li>
							<li class="nav-item">
								<RouterLink :to="'/profile/' + user_State.username" class="nav-link">
									<svg class="feather">
										<use href="/feather-sprite-v4.29.0.svg#layout" />
									</svg>
									Profile
								</RouterLink>
							</li>
							<li class="nav-item" @click="Logout"> <!---->
								<RouterLink to="/login" class="nav-link">
									<svg class="feather">
										<use href="/feather-sprite-v4.29.0.svg#key" />
									</svg>
									Logout
								</RouterLink>
							</li>
						</ul>

						<h6
							class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
							<span>Secondary menu</span>
						</h6>
						<ul class="nav flex-column">
							<li class="nav-item">
								<RouterLink to="/GodMode" class="nav-link">
									<svg class="feather">
										<use href="/feather-sprite-v4.29.0.svg#file-text" />
									</svg>
									SUPER USER
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

<style></style>





