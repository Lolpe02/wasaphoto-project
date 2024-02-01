<script>


export default {
    props: {
        post_data: {
            type: Number
        }
    },
    components: {

    },
    data: function () {
        return {
            datetime: null,
            is_your_post: false,
            photo_id: null,
            description: null,
            have_i_liked_this: false,
            creatorname: null,
            likes: 0,
            comments: []
        }
    },

    emits: ["delete-post"],

    methods: {

        async initialize() {
            /*
                let photo_jhjid = document.querySelector("input").value;
                photo_id = photo_id.trim();
                if (photo_id.length > 0) 
                {
                    let response = await this.$axios.get("/Images/" + photo_id,
                        {
                        headers: {
                        'Content-Type': 'application/json',
                        'accept': 'application/json, image/*',
                        'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                        }
                        });
                }*/
            this.photo_id = this.post_data;
            await this.$axios.get("/Images/" + this.post_data + "/metadata/",
                {
                headers: {
                'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                'Content-Type': 'application/json',
                'accept': 'application/json',
                }
            }).catch((error) => {
                if (error.response.status == 403) {
                    alert("You are not authorized to view this photo");
                    
                } else if (error.response.status == 404) {
                    alert("Photo not found");
                } else {
                    alert("Error: " + error.response.data);
                }
                return;
            }).then(response => {
                this.datetime = response.data.date;
                this.creatorname = response.data.creatorName;
                this.description = response.data.description;
                
                // format to dd month yyyy at hh:mm
                this.datetime = this.datetime.split("T");
                let date = this.datetime[0].split("-");
                let time = this.datetime[1].split(":");
                time = time[0] + ":" + time[1];
                this.datetime = date[2] + "/" + date[1] + "/" + date[0] + " at " + time;
                this.is_your_post = this.creatorname == this.$user_state.username;
            });                  
            
            // Fetch likes
            await this.$axios.get("/Images/" + this.post_data + "/likes/", {
                headers: {
                    'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                    'Content-Type': 'application/json',
                    'accept' : 'application/json',
                }
            }).catch((error) => {
                if (error.response.status == 403) {
                    alert("You are not authorized to view this photo");
                } else if (error.response.status == 404) {
                    alert("Photo not found");
                } else {
                    alert("Error: " + error.response.data);
                }
                return;
            }).then((response) => {
                if (response.data == null || response == undefined) {
                    this.likes = 0;
                    this.have_i_liked_this = false;
                    return;
                } else {
                    this.likes = response.data.length;
                    this.have_i_liked_this = response.data.includes(this.$user_state.username);
                }
            });

            // Fetch comments
            await this.$axios.get("/Images/" + this.post_data + "/comments/", {
                headers: {
                    'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                    'Content-Type': 'application/json',
                    'accept': 'application/json'
                }
            }).catch((error) => {
                if (error.response.status == 403) {
                    alert("You are not authorized to view this photo");
                } else if (error.response.status == 404) {
                    alert("Photo not found");
                } else {
                    alert("Error: " + error.response.data);
                }
                return;
            }).then((response) => {
                if (response.data == null || response == undefined) {
                    this.comments = [];
                    return;
                } else {
                    this.comments = response.data;
                }
            });
        },

        async ToCommentWriter() {

            // Jump to the comment writer

            let comment_writer = document.getElementById("comment-writer");

            comment_writer.scrollIntoView({
                behavior: "smooth",
                block: "start"
            });


        },
/*
        async AddComment(text) {

            // Update the frontend, then update the state on the server

            const creation_time = new Date().toISOString();

            console.log("Creation time: " + creation_time);

            const to_hash = creation_time + this.$user_state.username + this.photo_id;

            // SHA256 hash the comment ID

            const comment_id = this.$hasher(to_hash);

            let comm_obj = {
                comment_id: {
                    "hash": comment_id
                },
                author: {
                    "username-string": this.$user_state.username
                },
                body: text,
                "creation-time": creation_time,
                parent_post: {
                    "hash": this.photo_id
                }
            }

            console.log(comm_obj)

            // Add to the comments array, triggering a re-render
            this.comments.push(comm_obj);

            // Update the state on the server

            console.log("Request Path: " + "/users/" + this.post_data.author_name["username-string"] + "/profile/photos/" + this.photo_id + "/comments");

            let response = await this.$axios.put("/users/" + this.post_data.author_name["username-string"] + "/profile/photos/" + this.photo_id + "/comments/" + comment_id,
                comm_obj,
                {
                    headers: {
                        "Authorization": this.$user_state.headers.Authorization,
                        "commenter_name": this.$user_state.username
                    }
                });


        },

        async DeletePost() {

            // Update the state on the server

            console.log("Request Path: " + "/users/" + this.post_data.author_name["username-string"] + "/profile/photos/" + this.photo_id);

            let response = await this.$axios.delete("/users/" + this.post_data.author_name["username-string"] + "/profile/photos/" + this.photo_id, {
                headers: {
                    "Authorization": this.$user_state.headers.Authorization
                }
            });

            if (response.statusText != "No Content") {
                alert("Error: " + response.statusText);
                return;
            }

            // Remove the post from the stream
            this.$emit("delete-post", this.post_data);

        },

        async Like() {

            // Update the state on the server

            console.log("Request Path: " + "/users/" + this.post_data.author_name["username-string"] + "/profile/photos/" + this.photo_id + "/likes/" + this.$user_state.headers.Authorization);

            console.log(this.$user_state.headers.Authorization)

            let response = await this.$axios.put("/users/" + this.post_data.author_name["username-string"] + "/profile/photos/" + this.photo_id + "/likes/" + this.$user_state.headers.Authorization, {}, {
                headers: {
                    "Authorization": this.$user_state.headers.Authorization
                }
            });

            if (response.statusText != "No Content") {
                alert("Error: " + response.statusText);
                return;
            }

            // Only for consistency, the component does this internally.
            this.have_i_liked_this = true;
            this.likes++;
        },

        async Unlike() {

            // Update the state on the server

            console.log("Request Path: " + "/users/" + this.post_data.author_name["username-string"] + "/profile/photos/" + this.photo_id + "/likes/" + this.$user_state.headers.Authorization);

            console.log(this.$user_state.headers.Authorization)

            let response = await this.$axios.delete("/users/" + this.post_data.author_name["username-string"] + "/profile/photos/" + this.photo_id + "/likes/" + this.$user_state.headers.Authorization, {
                headers: {
                    "Authorization": this.$user_state.headers.Authorization
                }
            });

            if (response.statusText != "No Content") {
                alert("Error: " + response.statusText);
                return;
            }

            this.have_i_liked_this = false;
            this.likes--;
        },

        async DeleteComment(comment) {

            // Update the state on the server

            console.log("Request Path: " + "/users/" + this.post_data.author_name["username-string"] + "/profile/photos/" + this.photo_id + "/comments/" + comment.comment_id["hash"]);

            let response = await this.$axios.delete("/users/" + this.post_data.author_name["username-string"] + "/profile/photos/" + this.photo_id + "/comments/" + comment.comment_id["hash"], {
                headers: {
                    "Authorization": this.$user_state.headers.Authorization
                }
            });

            if (response.statusText != "No Content") {
                alert("Error: " + response.statusText);
                return;
            }

            // Remove the comment from the array
            this.comments = this.comments.filter((c) => c.comment_id["hash"] != comment.comment_id["hash"]);

        },
*/

    },

    mounted() {

        this.initialize();

    }
}
</script>

<template>


    <!-- Bordered Wrapper -->

    <div class="rounded p-2 m-2 border shadow-lg" style="width: 500px;">
        <div class="row align-content-between my-2">
            <div  class="col">
                <i class="bi-person-circle mx-2" style="font-size: 2em"></i>
                <span class="col font-weight-bold h1" v-if="!is_your_post">
                    {{ creatorname }}
                </span>
            </div>

            <!-- Right-aligned datetime -->
            <div class="col-auto ">
                <span class="text-muted v-center" style="font-size: 0.8em, font-style: italic;">
                    {{ datetime }}
                </span>

            </div>

            <!-- Delete Button -->

            <div class="col-auto" v-if="is_your_post">
                <button class="btn btn-danger v-center" @click="DeletePost">
                    <i class="bi-trash">X</i>
                </button>
            </div>
        </div>

        <div class="row">
            <Photo :src="post_data" :alt="description" :style="{ width: '100%' }" @click=" width= 2*width ">
            </Photo>
        </div>

        <!-- Divider <hr class="mt-1 mb-4"> -->

        

        <!-- Caption -->

        <div class="row">
            <div class="col-12">
                <i class="bi-person-circle mx-1" style="font-size: 1.5em"></i>
                <span class="font-weight-bold h1" style="margin-right: 5px;">{{ description }}</span>
            </div>
        </div>

        <!-- Divider -->

        <hr class="my-4">

        <!-- Comments -->
        <div class="row mt-3 align-content-start">
            <div class="col-auto">
                <LikeCounter class="v-center" :likes_count="this.likes" :liked="this.have_i_liked_this" @like="Like"
                    @unlike="Unlike"></LikeCounter>
            </div>
            <div class="col-auto d-flex align-items-center pb-2">
                <button class="btn bt-sm comment-button btn-outline-primary v-center">
                    <i class="bi-chat" @click="ToCommentWriter">Comment</i>
                </button>
            </div>

        </div>

        <div class="row">
            <div v-if="comments.length == 0" class="col-12 align-content-center w-100"><!-- Center the text -->
                <span class="h5 mx-1 font-weight-bold align-middle text-muted text-center">No comments yet.</span>
            </div>
            <div v-else class="col-12">
                <span class="h4 mx-1 font-weight-bold align-middle mb-2 text-start">Comments: </span>
            </div>
            <div class="col-12 my-3">
                <Comment v-for="comment in comments" :comment="comment" :key="comment.commentId"
                    @delete="DeleteComment">
                </Comment>
            </div>

        </div>

        <!-- CommentWriter -->

        <div class="row my-0">
            <div class="col-12">
                <CommentWriter id="comment-writer" :photo_id="this.post_data" :author_name="creatorname" @comment="AddComment">
                </CommentWriter>
            </div>
        </div>


    </div>


</template>

<style>
.v-center {

    display: inline-block;
    vertical-align: middle;
    line-height: normal;

}

.comment-button {

    padding: 0.35rem 0.5rem;
    font-size: 0.8rem;
    width: 120px;

}
</style>
