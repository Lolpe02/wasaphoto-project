<script>
export default {
    data: function () {
        return {
            followers: 0,
            following: 0,
            posts: 0,
            is_me: false,
            is_banned: false,
            is_following: false,
            username: null,
            has_banned_you: false,
            photos: [] // list of IDs, pairs of ("hash", SHA256 hash of the photo)
        }
    },
    methods: {
        async refresh() {

            this.username = this.$route.params.username;

            // Redirect to login if not logged in
            if (this.$user_state.username == null) {
                this.$router.push("/login");
                return
            }

            if (this.$route.params.username == this.$user_state.username) {
                this.is_me = true;
            }

            this.$user_state.current_view = this.$views.PROFILE;

            // this.has_banned_you = response.data.users.map(x => x["username-string"]).includes(this.$user_state.username);
            let response = await this.$axios.get("/Users/", {
                headers: {
                    "Authorization": 'Bearer ' + this.$user_state.headers.Authorization,
                    "accept": "application/json",
                    "Content-Type": "application/json",
                },
                params: {
                    "userName": this.$route.params.username
                }
            }).catch(err => {

                if (err.response.status == 404) {
                    alert("User not found");
                    this.$router.push("/");
                    return
                } else {
                    alert("Error: " + err.response.data);
                    return
                }
            }).then(res => {

                if (res == undefined) {
                    console.log("Error: undefined response");
                    return
                }

                if (res.statusText != "OK") {
                    alert("Error: " + res.statusText);
                    return
                }
                console.log(res, res.data);
                
            });
            if (response.data["followed"] != undefined) {
                this.following = response.data["followed"].length;

            }
            if (response.data["following"] != undefined) {
                this.followers = response.data["following"].length;
            }
            if (response.data["posted"] != undefined) {
                this.photos = response.data["posted"];
                this.posts = this.photos.length;
            }      
        },

        async DeletePost(post_data) {

            this.refresh();
        },

        async ChangeName() {

            const new_name = prompt("Change name", "New name");

            if (new_name == null || new_name == "") {
                return
            }


            if (!new_name.match("^.{3,25}$")) {
                alert("Invalid username");
                return;
            }

            const res = await this.$axios.patch("/Users/",
            new_name, {
                headers: {
                    "Authorization": this.$user_state.headers.Authorization,
                    "accept": "application/json",
                }
            }).catch(err => {

                if (err.response.status == 404) {
                    alert("either banned by user or not following");
                    return
                } else {
                    alert("Error: " + err.response.data);
                    return
                }

                return
            }).then(res => {

                if (res == undefined) {
                    console.log("Error: undefined response");
                    return
                }

                if (res.statusText != "OK") {
                    alert("Error: " + res.statusText);
                    return
                }

                this.$user_state.username = new_name;
                this.username = new_name;
                this.$router.push("/profile/" + new_name);
            });

        },

        async Follow() {

            const req_body = {
                "username-string": this.$user_state.username
            }

            const res = await this.$axios.put("/users/" + this.$user_state.username + "/following/" + this.username, req_body, {
                headers: this.$user_state.headers
            });

            if (res.statusText != "No Content") {

                alert("Error: " + res.statusText);
                return
            }

            this.is_following = true;
            this.followers += 1;
        },

        async Unfollow() {

            if (!this.is_following) {
                return
            }

            const req_body = {
                "username-string": this.$user_state.username
            }

            const res = await this.$axios.delete("/users/" + this.$user_state.username + "/following/" + this.username, {
                headers: this.$user_state.headers,
                data: req_body
            });

            if (res.statusText != "No Content") {

                alert("Error: " + res.statusText);
                return
            }

            this.is_following = false;
            this.followers -= 1;
        },

        async Ban() {

            const res = await this.$axios.put("/users/" + this.$user_state.username + "/bans/" + this.username, {}, {
                headers: this.$user_state.headers
            });

            if (res.statusText != "No Content") {

                alert("Error: " + res.statusText);
                return
            }

            this.is_banned = true;
        },

        async UnBan() {

            const res = await this.$axios.delete("/users/" + this.$user_state.username + "/bans/" + this.username, {
                headers: this.$user_state.headers
            });

            if (res.statusText != "No Content") {

                alert("Error: " + res.statusText);
                return
            }

            this.is_banned = false;

        },

    },

    mounted() {
        this.refresh()
    }
}
</script>

<template>
    <div class="container">
        <div class="align-items-center text-center h-100">
            <div class="container text-center pt-3 pb-2 border-bottom">
                <div class="row w-100 my-3">
                    <h2 class="col-3 text-break d-inline-block" style="vertical-align: middle;">
                        <i class="bi-person-circle mx-1"></i>  {{$user_state.username}}'s profile.
                    </h2>
                    <div class="col-9" style="align-items: center; vertical-align: middle;">
                        <div class="row">
                            <div class="col-4">
                                <div class="row border p-1 pt-2 rounded me-1 shadow-sm">
                                    <div class="col-12">
                                        <h5>Posts</h5>
                                    </div>
                                    <div class="col-12">
                                        <h5> {{ posts }}</h5>
                                    </div>
                                </div>
                            </div>
                            <div class="col-4">
                                <div class="row border p-1 pt-2 rounded me-1 shadow-sm">
                                    <div class="col-12">
                                        <h5>Followers</h5>
                                    </div>
                                    <div class="col-12">
                                        <h5>{{ followers }}</h5>
                                    </div>
                                </div>
                            </div>
                            <div class="col-4">
                                <div class="row border p-1 pt-2 rounded me-1 shadow-sm">
                                    <div class="col-12">
                                        <h5>Following</h5>
                                    </div>
                                    <div class="col-12">
                                        <h5>{{ following }}</h5>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div v-if="is_me" class="row w-100">
                    <div class="row w-100">
                        <div class="col-3">
                            <button class="btn btn-primary btn-md" type="button" @click="ChangeName()">
                                <i class="bi-pencil-square"></i>
                                Change Name
                            </button>
                        </div>
                    </div>
                </div>
                <div v-else>
                    <div class="row w-100 align-content-between my-1">
                        <!-- Follow Button -->
                        <div class="col">
                            <Transition name="fade" mode="out-in">
                                <div v-if="is_following && !has_banned_you">
                                    <button class="btn btn-warning btn-lg" type="button" @click="Unfollow()">
                                        <i class="bi-person-dash-fill"></i>
                                        Unfollow
                                    </button>
                                </div>
                                <div v-else-if="!is_following && !has_banned_you">
                                    <button class="btn btn-primary btn-lg" type="button" @click="Follow()">
                                        <i class="bi-person-plus-fill"></i>
                                        Follow
                                    </button>
                                </div>
                            </Transition>
                        </div>
                        <!-- Ban Button -->
                        <div class="col">
                            <Transition name="fade" mode="out-in">
                                <div v-if="is_banned && !has_banned_you">
                                    <button class="btn btn-success btn-lg" type="button" @click="UnBan()">
                                        <i class="bi-person-check-fill"></i>
                                        Unban
                                    </button>
                                </div>
                                <div v-else-if="!is_banned && !has_banned_you">
                                    <button class=" btn btn-danger btn-lg" type="button" @click="Ban()">
                                        <i class="bi-person-x-fill"></i>
                                        Ban
                                    </button>
                                </div>
                            </Transition>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
    <div v-if="has_banned_you" class="container">
        <div class="row">
            <div class="col-12">
                <div class="alert alert-danger" role="alert">
                    <h4 class="alert-heading">You have been banned by this user!</h4>
                    <p>Sorry, but you have been banned from this user's profile. You cannot view their posts or
                        interact with them.</p>
                    <hr>
                    <p class="mb-0">Try not to be so mean next time!</p>
                </div>
            </div>
        </div>
    </div>
    <div v-else class="container">
        <Stream :posts="photos" @delete-post="DeletePost" :key="photos.length"></Stream>
    </div>
</template>

<style>
.fade-enter-active,
.fade-leave-active {
    transition: opacity cubic-bezier(0.4, 0, 0.2, 1) 0.1s
}

.fade-enter,
.fade-leave-to {
    opacity: 0
}
</style>