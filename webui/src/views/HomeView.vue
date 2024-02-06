<script>
export default {
	data: function () {
		return {
			errormsg: null,
			loading: false,
			feed: [],
			modal: null,
			fetching: false,
			there_are_more_posts: true,
		}
	},
	methods: {

		async ToProfile() {

			if (this.$user_state.username == null) {
				return
			}

			this.$router.push("/profile/" + this.$user_state.username)// ;
		},

		async DeletePost(post_data) {

			this.refresh();
		},

		async refresh() {
			if (this.$user_state.username == null) {
				console.log("Empty username, redirecting to login")
				this.$router.push("/login");
				return;
			}
			this.loading = true;
			this.errormsg = null;
			this.$user_state.current_view = this.$views.STREAM;

			const mod = bootstrap.Modal.getOrCreateInstance(document.getElementById('exampleModal'))
			document.body.appendChild(mod._element)
			this.modal = mod

			await this.$axios.get("/Users/me/myStream", {
				headers: {
					'Content-Type': 'application/json',
					'accept': 'application/json',
					'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
				}
			}).then((response) => {
				if (response === undefined || response.data == null) {
					console.log("undefined response, empty array");
					this.there_are_more_posts = false;
					return;
				}
				if (response.data.length == 0) {
					this.there_are_more_posts = false;
				} else {
					this.there_are_more_posts = true;
					this.feed.push(...response.data);
				}
			}).catch((error) => {
				console.log(error);
				alert("Error loading stream");
				return;
			});

			window.onscroll = async () => {
				if (this.fetching) {
					return;
				}
				if (window.innerHeight + window.scrollY >= document.body.offsetHeight && this.$user_state.current_view == this.$views.STREAM) {
					this.loading = true;
					this.fetching = true;
					let newbatch = await this.$axios.get("/Users/me/myStream", {
						headers: {
							'Content-Type': 'application/json',
							'accept': 'application/json',
							'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
						}
					})
					if (newbatch == null || newbatch === undefined) {
						this.there_are_more_posts = false;
						return;
					}
					newbatch = newbatch.data;

					if (!this.there_are_more_posts) {
						return;
					}

					if (newbatch == null || newbatch.length == 0) {
						this.there_are_more_posts = false;
					}

					let newposts = newbatch.filter((item) => this.feed.indexOf(item) < 0);
					this.feed.push(...newposts);

					await new Promise(resolve =>
						setTimeout(resolve, 5000)).then(() => {
							console.log("waiting for 5 seconds")
							this.fetching = false;
						});

					this.loading = false;
				}
			};
			this.loading = false;
		},

		async UploadPhoto() {

			// Manually set the submit button to be waiting

			document.getElementById("submit-button").innerHTML = "Uploading...";
			document.getElementById("submit-button").classList.add("disabled");

			const image = document.getElementById("fileInput").files[0];

			if (!image) {
				alert("Please select a file to upload");
				return;
			}

			let reader = new FileReader();
			let data = null;

			reader.onload = async () => {
				data = reader.result;
			}

			reader.onerror = function (error) {
				console.log('Error: ', error);
				alert("Error uploading photo")
				return
			};

			reader.readAsDataURL(image);

			// Wait for the reader to finish reading the file

			while (data == null) {
				await new Promise(r => setTimeout(r, 1000));
				console.log("waiting for reader to finish")
			}

			// const filename = image.name;
			// get format
			// const format = filename.split('.').pop();

			const caption = document.getElementById("captionInput").value;
			// strip data:image/png;base64, from the beginning of the string

			data = data.substring(22);

			// Create a new instance of FormData
			let formData = new FormData();

			// Append your data
			formData.append('photo', image); // Assuming `data` is your file object
			formData.append('description', caption);

			// Send it with axios
			await this.$axios.post("/Images/", formData, {
				headers: {
					"Content-Type": "multipart/form-data",
					"Authorization": "Bearer " + this.$user_state.headers.Authorization,
					"accept": "application/json"
				}
			}).catch(error => {
				if (error.response == undefined) {
					alert("undefined response");
					return
				}
				console.error(error);
			}).then(response => {
				console.log(response);
				if (response.status == 201) {

					// manually restyle and rename the submit button

					const submit_button = document.getElementById("submit-button");

					submit_button.classList.remove("btn-primary");
					submit_button.classList.remove("disabled");
					submit_button.classList.add("btn-success");
					submit_button.innerHTML = "Success!";

					setTimeout(() => {

						// reset the button

						submit_button.classList.remove("btn-success");
						submit_button.classList.add("btn-primary");
						submit_button.innerHTML = "Submit";
					}, 3000);

					// Clear the form

					document.getElementById("fileInput").value = "";
					document.getElementById("captionInput").value = "";
				}
				else {
					alert("Error uploading photo")
				}
			});
			if (document.getElementById("submit-button").innerHTML == "Uploading...") {
				setTimeout(() => {
					// reset the button
					submit_button.classList.remove("btn-success");
					submit_button.classList.add("btn-primary");
					submit_button.innerHTML = "Submit";
				}, 3000);
			}

		}
	},
	mounted() {
		this.refresh()
	},
	beforeMount() { // this is a hack to make sure the user is logged in, beforeCreate is not working
		if (this.$user_state.username == null) {
			console.log("Empty username before mount, redirecting to login")
			this.$router.push("/login");
		}
	},
	beforeCreate() {
		if (this.$user_state.username == null) {
			console.log("Empty username before create, redirecting to login")
			this.$router.push("/login");
		}
	},
	updated() {
		const mod = bootstrap.Modal.getOrCreateInstance(document.getElementById('exampleModal'))

		document.body.appendChild(mod._element)
		document.body.removeChild(this.modal._element)

		this.modal = mod
	}
}
</script>

<template>
	<div class="container">
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
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
					<button type="button" class="btn btn-primary" data-bs-toggle="modal" data-bs-target="#exampleModal"
						data-backdrop="false">
						New Post
					</button>
				</div>
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

	</div>
	<div v-if="feed.length == 0" class="container d-flex justify-content-center align-items-center">
		<h1>No posts to show</h1>
	</div>
	<div v-else >
		<Stream :posts="feed" :key="feed.length" @delete-post="DeletePost"></Stream>
	</div>
	<!-- Modal -->
	<div class="modal fade" id="exampleModal" tabindex="-1" aria-labelledby="exampleModalLabel" aria-hidden="true">
		<div class="modal-dialog modal-dialog-centered modal-dialog-scrollable">
			<div class="modal-content">
				<div class="modal-header">
					<h5 class="modal-title" id="exampleModalLabel">Create a Post</h5>
					<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
				</div>
				<div class="modal-body">
					<!-- Image Input -->
					<div class="mb-3">
						<div class="row g-3 align-items-center">
							<form id="formFile">
								<label for="fileInput" class="form-label">Upload Image</label>
								<input class="form-control" type="file" id="fileInput"
									accept="image/jpeg, img/png, img/gif">

								<label for="captionInput" class="form-label">Caption</label>
								<textarea class="form-control" type="text" id="captionInput" rows="6"></textarea>
							</form>
						</div>
					</div>

				</div>
				<div class="modal-footer">
					<button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
					<button type="button" class="btn btn-primary" id="submit-button" @click="UploadPhoto()">
						Submit
					</button>
				</div>
			</div>
		</div>
	</div>
</template>

<style></style>
