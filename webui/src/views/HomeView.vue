<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: [],
		}
	},
	methods: {
		async Logout() {
			this.$user_state.username = null;
			this.$user_state.headers.Authorization = null;
			console.log("Logging out")
			this.$router.push("/login");
		},

		

		async ToProfile() {

			if (this.$user_state.username == null) {
				return
			}

			this.$router.push("/profile/" + this.$user_state.username)// ;
		},

		/* async ToStream() {

			if (this.$user_state.username == null) {
				return
			}

			this.$router.push("/stream/" + this.$user_state.username);
		},	*/
		
		async refresh() {
			if (this.$user_state.username == null) {
				console.log("Empty username, redirecting to login")
				this.$router.push("/login");
				return;
			}
			this.loading = true;
			this.errormsg = null;
			let response = await this.$axios.get("/Users/me/myStream",
				{
					headers: {
						'Content-Type': 'application/json',
						'accept': 'application/json',
						'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
					},
				});
			this.some_data = response.data;
			this.loading = false;
		},
	},
	mounted() {
		this.refresh()
	},
	beforeMount() { // this is a hack to make sure the user is logged in, beforeCreate is not working
		if (this.$user_state.username == null || this.$user_state.username == undefined) {
			console.log("Empty username before mount, redirecting to login")
			this.$router.push("/login");
		}
	},
	beforeCreate() { 
		if (this.$user_state.username == null || this.$user_state.username == undefined) {
			console.log("Empty username before create, redirecting to login")
			this.$router.push("/login");
		}
	},
}
</script>

<template>
	<div>
		<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Home page</h1>
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
						Refresh
					</button>
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="ToProfile">
						Your Profile
					</button>
				</div>
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-primary" @click="newItem">
						New
					</button>
				</div>
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
	
</template>

<style>
</style>
