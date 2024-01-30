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

		async PerformSearch() {

			let search = document.querySelector("input").value;

			search = search.trim();

			if (search.length > 0) {
				// query ./users for results

				const searcher_id = this.$user_state.headers.Authorization;

				if (searcher_id == null) {
					return
				}

				const header = {
					"Authorization": searcher_id,
					"user_name": this.$user_state.username
				}

				let response = await this.$axios.get("/users", {
					params: {
						"search_term": search
					},
					headers: header
				});

				if (response.status == 200) {
					this.search_results = response.data;
				} else {
					console.log
					this.search_results = null;
				}
			}
			else {
				this.search_results = null;
			}
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
			console.log("Refreshing", this.$user_state.headers.Authorization)
			let response = await this.$axios.get("/Users/me/myStream",
				{
					headers: {
						'Content-Type': 'application/json',
						'accept': 'application/json',
						'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
					},
				});
			this.some_data = response.data;
			console.log('Response: ', response.data);
			this.loading = false;
		},
	},
	mounted() {
		this.refresh()
	},
	beforeMount() { // this is a hack to make sure the user is logged in, beforeCreate is not working
		if (this.$user_state.username == null || this.$user_state.username == undefined) {
			console.log("Empty username, redirecting to login")
			this.$router.push("/login");
		}
	},
	beforeCreate() { 
		if (this.$user_state.username == null || this.$user_state.username == undefined) {
			console.log("Empty username, redirecting to login")
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
