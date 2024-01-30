<script>
//import Modal from '../components/Modal.vue';
export default {
    
    data: function () {
        return {
            followers: 1,
            following: 0,
            posts: 0,
            is_me: false,
            is_banned: false,
            is_following: false,
            username: null,
            isListVisible: true,
            toorfrom: false,
            has_banned_you: false,
            subscription: null,
            their_id: null,
            photos: [], // list of IDs, pairs of ("hash", SHA256 hash of the photo)
            followerList: [],
            followingList: [],
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
            }).then(response => {

                if (response == undefined || response.data == null) {
                    console.log("Error: undefined response gettig profile");
                    return
                }

                if (response.status != 200) {
                    alert("Error: " + response.data);
                    return
                }
                console.log(response.data);
                if (response.data["followed"] == undefined || response.data["followed"] == null) {
                    this.following = 0;
                } else {
                    this.followingList = response.data["followed"];
                    this.following = response.data["followed"].length;
                }
                if (response.data["following"] == undefined || response.data["following"] == null) {
                    this.followers = 0;
                } else {
                    this.followerList = response.data["following"];
                    this.followers = response.data["following"].length;
                }
                if (response.data["posted"] == undefined || response.data["posted"] == null) {
                    this.posts = 0;
                } else {
                    this.photos = response.data["posted"];
                    this.posts = this.photos.length;
                }
                this.subscription = this.FormatDate(response.data["date"]);     
                this.their_id = response.data["userId"];
                if (this.their_id == this.$route.params.username) {
                    alert("WHAT HAPPENED?");
                    return
                }

                this.is_following = true;
                this.has_banned_you = false;
            });
            
        },
        
        async ToggleVisibility(listype) {
            this.isListVisible = !this.isListVisible;
            console.log(this.isListVisible),
            this.toorfrom = listype;
        },

        FormatDate(datestamp) {
            let datetime = datestamp.split("T");
            let date = datetime[0].split("-");
            let time = datetime[1].split(":");
            time = time[0] + ":" + time[1];
            return date[2] + "/" + date[1] + "/" + date[0] + " at " + time;
		},

        async DeletePost(post_data) {

            this.refresh();
        },

        async ChangeName() {

            const new_name = prompt("Set a new name", "Who will you be?");

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
            let response = await this.$axios.post("/Users/me/following/",
                this.$user_state.username, {
                headers: {
                'Content-Type': 'application/json',
                'accept': 'application/json',
                'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                }
            });

            if (response.status == 404) {
                alert("Error: " + response.statusText);
                return
            } else if (response.status == 403) {
                alert("Error: " + response.statusText);
                $router.push("/login");
                return
            } else if (response.status == 201) {
                alert("Now following user");
            } else if (response.status == 200) {
                alert("Already followed user");
            } else {
                alert("Error: " + response.statusText);
                return
            }

            this.is_following = true;
            this.followers += 1;
        },

        async Unfollow() {

            if (!this.is_following) {
                return
            }

            let response = await this.$axios.delete("/Users/me/following/" + user_id,
                    {
                    headers: {
                    'Content-Type': 'application/json',
                    'accept': 'application/json',
                    'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                    }
                });

            if (response.status != 200) {

                alert("Error: " + response.data);
                return
            }

            this.is_following = false;
            this.followers -= 1;
        },

        async Ban() {

            let response = await this.$axios.post("/Users/me/muted/", 
                user_id, {
				headers: {
					'Content-Type': 'application/json',
					'accept': 'application/json',
					'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
				}});

            if (response.status == 200 ) {
                alert("Already banned user");
            } else if (response.status == 201) {
                alert("Now banned user");
            } else {
                alert("Error: " + response.data);
                return
            }

            this.is_banned = true;
        },

        async UnBan() {

            let response = await this.$axios.delete("/Users/me/muted/" + user_id,
                    {
                    headers: {
                    'Content-Type': 'application/json',
                    'accept': 'application/json',
                    'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                    }
                });

            if (response.status != 200) {
                alert("Error: " + response.data);
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
                <div class="row w-80 my-3">
                    <h2 class="col-3 text-break d-inline-block" style="vertical-align: bottom;">
                        <i class="bi-person-circle mx-1"></i>  {{$user_state.username}}'s profile.
                    </h2>
                    <div class="col-9" style="align-items: center; vertical-align: middle;">
                        <div class="row">
                            <div class="col-3">
                                <div class="row border p-1 pt-2 rounded me-1 shadow-sm">
                                    <div class="col-12">
                                        <h5>Posts: {{ posts }}</h5>
                                    </div>
                                </div>
                            </div>
                            <div class="col-3">
                                <div class="row border p-1 pt-2 rounded me-1 shadow-sm">
                                    <!-- @click= "ToggleVisibility(false)"this.isListVisible data-toggle="modal"=<button  >Show list</button> !this.isListVisible v-b-modal= "FollowL"-->
                                    <button type="button" class="btn btn-primary" @click= "toorfrom=true" data-bs-toggle="modal" data-bs-target="#ciao">List</button>
                                    <div class="col-12">
                                        <h5>Followers: {{ followers }}</h5>
                                    </div>
                                    
        
                                </div>
                            </div>
                            <div class="col-3">
                                <div class="row border p-1 pt-2 rounded me-1 shadow-sm">
                                    <!--<button @click= "ToggleVisibility(true)" >Show list</button>-->
                                    <button type="button" class="btn btn-primary" @click= "toorfrom=false" data-bs-toggle="modal" data-bs-target="#ciao">List</button>

                                    <div class="col-12">
                                        <h5>Following: {{ following }}</h5>
                                    </div>
                                </div>
                            </div>
                            <div class="col-3">
                                <div class="row border p-1 pt-2 rounded me-1 shadow-sm">
                                    <div class="col-12">
                                        <h5>Subscription: {{ subscription }}</h5>
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
        <Stream :posts="photos" @delete-post="DeletePost" :key="photos.length"></Stream><!---->
    </div>
    <!--  v-if="isListVisible" -->
    <Modal id="ciao">
            <!---->
            <template v-if="toorfrom" v-slot:header>
                People following this user
            </template>
            <template v-else v-slot:header>
                People this user is following
            </template>
            <template v-if="toorfrom" v-slot:body>
                <ul>
                    <li v-for="(follower, index) in this.followerList" :key="index">{{index+1}} - {{ follower }}</li>
                </ul>
            </template>
            <template v-else v-slot:body>
                <ul>
                    <li v-for="(follower, index) in this.followingList" :key="index">{{index+1}} - {{ follower }}</li>
                </ul>
            </template>
        </Modal>
    
</template>

<style>
.fade-enter-active,
.fade-leave-active {
    transition: opacity cubic-bezier(0.4, 0, 0.2, 1) 0.1s
}
.fade-enter,
.fade-leave-to {
    opacity: 0
};
</style>