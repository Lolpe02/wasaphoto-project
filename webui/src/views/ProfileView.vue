<script>
export default {
    data: function () {
        return {
            search_results: null,
            followers: 0,
            following: 0,
            posts: 0,
            is_me: false,
            is_banned: false,
            is_following: false,
            searchedUser: "",
            username: null,
            isListVisible: true,
            toorfrom: false,
            has_banned_you: false,
            subscription: null,
            their_id: -1,
            photos: [], // list of IDs, pairs of ("hash", SHA256 hash of the photo)
            followerList: [],
            followingList: [],
        };
    },
    methods: {
        async refresh() {
            console.log("Refreshing profile");
            // Redirect to login if not logged in
            if (this.$user_state.username == null || this.$user_state.headers.Authorization == null) {
                this.$router.push("/login");
                return;
            }
            this.username = this.$route.params.username;
            this.is_me = this.$route.params.username == this.$user_state.username;
            this.$user_state.current_view = this.$views.PROFILE;
            // 
            await this.$axios.get("/Users/profile", {
                headers: {
                    "Authorization": 'Bearer ' + this.$user_state.headers.Authorization,
                    "accept": "application/json",
                    "Content-Type": "application/json",
                },
                params: {
                    "userName": this.$route.params.username
                }
            }).then(response => {
                if (response == undefined || response.data == null) {
                    console.log("Error: undefined response getting profile");
                    return;
                }
                if (response.data.posted == null) {
                    this.posts = 0;
                } else {
                    this.photos = response.data.posted;
                    this.posts = this.photos.length;
                    this.subscription = this.FormatDate(response.data.date);
                    this.their_id = response.data.userId;
                    let their_name = response.data.userName;
                    if (their_name != this.$route.params.username) {
                        alert("WHAT HAPPENED? CRITICAL ERROR! BOMBING YOUR HOUSE NOW!");
                        this.$router.push("/login");
                        return;
                    }
                }
            }).catch(err => {
                this.search_results = null;
                if (err.response.status == 400) {
                    alert("Error getting profile: " + err.response.data);
                } else if (err.response.status == 404) {
                    alert("Profile not found");
                    console.log("Users not found");
                } else if (err.response.status == 403) {
                    console.log("banned by user");
                    this.has_banned_you = true;
                } else {
                    alert("Error: " + err.response.data);
                }
                return;
            });
            await this.$axios.get("/Users/me/following/", {
                headers: {
                    "Authorization": 'Bearer ' + this.$user_state.headers.Authorization,
                    "accept": "application/json",
                    "Content-Type": "application/json",
                },
                params: {
                    "userName": this.$route.params.username
                }
            }).then(response => {
                if (response == undefined || response.data == null) {
                    this.followingList = [];
                    this.following = 0;
                } else {
                    this.followingList = response.data;
                    this.following = this.followingList.length;
                }
            }).catch((err) => {
                if (err.response.status == 404) {
                    console.log("Users not found");
                    this.search_results = null;
                }
                else if (err.response.status == 403) {
                    console.log("banned by user");
                    this.has_banned_you = true;
                    this.search_results = null;
                } else {
                    alert("Error getting following: " + err.response.data);
                    return;
                }
            });
            await this.$axios.get("/Users/me/followers/", {
                headers: {
                    "Authorization": 'Bearer ' + this.$user_state.headers.Authorization,
                    "accept": "application/json",
                    "Content-Type": "application/json",
                },
                params: {
                    "userName": this.$route.params.username
                }
            }).then(response => {
                if (response == undefined || response.data == null) {
                    this.followerList = [];
                    this.followers = 0;
                } else {
                    this.followerList = response.data;
                    this.followers = this.followerList.length;
                    this.is_following = this.followerList.includes(this.$user_state.username);
                }
            }).catch((err) => {
                if (err.response.status == 404) {
                    console.log("Users not found");
                    this.search_results = null;
                }
                else if (err.response.status == 403) {
                    console.log("banned by user");
                    this.has_banned_you = true;
                    this.search_results = null;
                } else {
                    alert("Error getting followers: " + err.response.data);
                    return;
                }
            });
        },
        SeeProfile(name, id) {
            this.search_results = null;
            this.searchedUser = name
            this.their_id = id;
            this.$router.push("/profile/" + this.searchedUser);
        },
        async PerformSearch() {
            let search = document.querySelector("input").value;
            search = search.trim();
            if (search.length > 2 && search.match("^.{3,25}$")) {

                await this.$axios.get("/Users/", {
                    headers: {
                        "Authorization": 'Bearer ' + this.$user_state.headers.Authorization,
                        "accept": "application/json",
                        "Content-Type": "application/json",
                    },
                    params: {
                        "userName": search
                    }
                }).catch(err => {

                    if (err.response.status == 404) {
                        console.log("Users not found");
                        this.search_results = null;
                    }
                    else if (err.response.status == 403) {
                        console.log("banned by user");
                        this.search_results = null;
                    } else {
                        alert("Error: " + err.response.data);
                        return;
                    }
                }).then(response => {
                    if (response == undefined || response.data == null) {
                        alert("undefined response");
                        return
                    }
                    if (response.status != 200) {
                        alert("Error: " + response.data);
                        return;
                    }
                    this.search_results = response.data;
                });
            }
            else {
                this.search_results = null;
            }
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
            console.log("Changing name", this.$user_state.headers.Authorization);
            const new_name = prompt("Set a new name", "Who will you be?");
            if (new_name == null || new_name == "") {
                return;
            }
            if (!new_name.match("^.{3,25}$")) {
                alert("Invalid username");
                return;
            }
            await this.$axios.patch("/Users/", new_name, {
                headers: {
                    "Authorization": "Bearer " + this.$user_state.headers.Authorization,
                    "accept": "application/json",
                    "Content-Type": "application/json",
                }
            }).then(response => {
                if (response == undefined || response.data == null) {
                    alert("undefined response");
                    return
                }
                this.$user_state.username = new_name;
                this.username = new_name;
                this.$router.push("/profile/" + new_name);
            }).catch((error) => {
                if (error.response.status == 404) {
                    alert("profile not found");
                    return;
                }
                else {
                    alert("Error: " + error.response.data);
                    return;
                }
            });
        },
        async Follow() {
            console.log("Following user", this.their_id);
            await this.$axios.post("/Users/me/following/",
                this.their_id, {
                headers: {
                    'Content-Type': 'application/json',
                    'accept': 'application/json',
                    'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                }
            }).then(response => {
                if (response == undefined || response.data == null) {
                    alert("undefined response");
                    return
                }
                if (response.status == 201) {
                    alert("Now following user");
                }
                else if (response.status == 200) {
                    alert("Already followed user");
                    return;
                }
                else {
                    alert("Error: " + response.data);
                    return;
                }
                this.is_following = true;
                this.followers += 1;
            }).catch(error => {
                if (error.response.status == 401) {
                    alert("Error: " + error.response.data);
                    this.$router.push("/login");
                    return;
                } else if (error.response.status == 403) {
                    this.has_banned_you = true;
                    alert("Banned by " + this.username);
                    return;
                } else if (error.response.status == 404) {
                    alert("Error: " + error.response.data);
                    return;
                } else {
                    alert("Error: " + error.response.data);
                    return;
                }
            });

        },
        async Unfollow() {
            if (!this.is_following) {
                return;
            }
            let response = await this.$axios.delete("/Users/me/following/" + this.their_id, {
                headers: {
                    'Content-Type': 'application/json',
                    'accept': 'application/json',
                    'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                }
            });
            if (response.status != 200) {
                alert("Error: " + response.data);
                return;
            }
            this.is_following = false;
            this.followers -= 1;
        },
        async Ban() {
            let response = await this.$axios.post("/Users/me/muted/", this.their_id, {
                headers: {
                    'Content-Type': 'application/json',
                    'accept': 'application/json',
                    'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                }
            });
            switch (response.status) {
                case 200:
                    alert("Already banned user");

                    break;
                case 201:
                    alert("Now banned user");
                    break;
                default:
                    alert("Error: " + response.data);
                    return;
            }
            if (this.followingList.includes(this.$user_state.username)) {
                // theyre following you,
                this.following--;
                this.followingList = this.followingList.filter((value) => value != this.$user_state.username); 
            };      
            this.is_banned = true;
        },
        async UnBan() {
            let response = await this.$axios.delete("/Users/me/muted/" + this.their_id, {
                headers: {
                    'Content-Type': 'application/json',
                    'accept': 'application/json',
                    'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                }
            });
            if (response.status != 200) {
                alert("Error: " + response.data);
                return;
            }
            this.is_banned = false;
        },

    },
    mounted() {
        this.refresh();
    },
    watch: {
        $route(to, from) {
            if (to.path == "/login" || to.path == "/") {
                return;
            }
            this.is_me = to.params.username == this.$user_state.username;
            this.searchedUser = "";
            this.search_results = null;
            if (to.params.username != from.params.username) {
                this.has_banned_you = false,
                    this.is_following = false,
                    this.is_banned = false,
                    this.followerList = [],
                    this.followingList = [],
                    this.username = null
                this.photos = [];
                this.refresh();
            };
        }
    }
}
</script>

<template>
    <div class="container">
        <div class="align-items-center text-center h-80">
            <div class="col-md-12 col-sm-6">
                <form class="nav form-inline my-2 my-md-0" :class="{
                    disabled: $user_state.username == null, 'd-none': $user_state.username == null
                }">
                    <input class="form-control" id="SearchBox" v-model="this.searchedUser" type="text"
                        placeholder="Search for your friends (if you have them)" aria-label="Search"
                        @input="PerformSearch()" v-on:keyup.enter="PerformSearch()">
                    <!-- Results -->
                    <ul class="list-group custom-select w-25 dropdown mt-5 position-absolute">

                        <li class=" list-group-item align-middle" v-if="search_results"
                            v-for="(userId, userName) in search_results" :key="userName"
                            @click="SeeProfile(userName, userId)">
                            <i class="bi-person-circle m-1 fa-lg" style="font-size: 1.0rem;"></i>
                            {{ userName }}
                        </li>
                    </ul>

                </form>
            </div>
            <div class="container text-center pt-3 pb-2 border-bottom">
                <div class="row w-80 my-3">
                    <h2 class="col-3 text-break d-inline-block" style="vertical-align: bottom;">
                        <i class="bi-person-circle mx-1"></i> {{ this.username }}'s profile.
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
                                    <button type="button" class="btn btn-primary" @click="toorfrom = true"
                                        data-bs-toggle="modal" data-bs-target="#ciao">List</button>
                                    <div class="col-12">
                                        <h5>Followers: {{ followers }}</h5>
                                    </div>


                                </div>
                            </div>
                            <div class="col-3">
                                <div class="row border p-1 pt-2 rounded me-1 shadow-sm">
                                    <!--<button @click= "ToggleVisibility(true)" >Show list</button>-->
                                    <button type="button" class="btn btn-primary" @click="toorfrom = false"
                                        data-bs-toggle="modal" data-bs-target="#ciao">List</button>

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
                <div v-if="this.is_me" class="row w-100">
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
                    <p class="mb-0"></p>
                </div>
            </div>
        </div>
    </div>
    <div v-else>
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
                <li v-for="(follower, index) in this.followerList || []" :key="index">{{ index + 1 }} - {{ follower }}</li>
            </ul>
        </template>
        <template v-else v-slot:body>
            <ul>
                <li v-for="(follower, index) in this.followingList || []" :key="index">{{index+1}} - {{ follower }}</li>
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
}

</style>