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
            comments: [],
            likesList: [],
            commentNumber: 0
        }
    },

    emits: ["delete-post"],

    methods: {

        async initialize() {

            this.photo_id = this.post_data;
            await this.$axios.get("/Images/" + this.post_data + "/metadata/",
                {
                    headers: {
                        'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                        'Content-Type': 'application/json',
                        'accept': 'application/json',
                    }
                }).then(response => {
                    if (response == undefined || response.data == null) {
                        alert("undefined response");
                        return
                    }
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
                }).catch((error) => {
                    if (error.response.status == 403) {
                        alert("You are not authorized to view this photo");
                    } else if (error.response.status == 404) {
                        alert("Photo not found");
                    } else {
                        alert("Error: " + error.response.data);
                    }
                    return;
                });

            // Fetch likes
            await this.$axios.get("/Images/" + this.post_data + "/likes/", {
                headers: {
                    'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                    'Content-Type': 'application/json',
                    'accept': 'application/json',
                }
            }).then((response) => {
                if (response == undefined || response.data == null) {
                    this.likes = 0;
                    this.have_i_liked_this = false;
                    return;
                } else {
                    this.likes = response.data.length;
                    this.have_i_liked_this = response.data.includes(this.$user_state.username);
                    this.likesList = response.data;
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
            });

            // Fetch comments
            await this.$axios.get("/Images/" + this.post_data + "/comments/", {
                headers: {
                    'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                    'Content-Type': 'application/json',
                    'accept': 'application/json'
                }
            }).then((response) => {
                if (response == undefined || response.data == null) {
                    this.comments = [];
                    this.commentNumber = 0;
                    return;
                } else {
                    this.comments = response.data;
                    this.commentNumber = response.data.length;
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
            });
        },
        /*async ToCommentManager() {

            // Jump to the comment writer

            let comment_writer = document.getElementById("comment-writer");

            comment_writer.scrollIntoView({
                behavior: "smooth",
                block: "start"
            });


        },*/

        async AddComment(text) {

            // Update the frontend, then update the state on the server
            let creation_time = new Date().toISOString();

            await this.$axios.post("/Images/" + this.photo_id + "/comments/",
                text, {
                headers: {
                    'Content-Type': 'application/json',
                    'accept': 'application/json',
                    'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                }
            }).then((response) => {
                if (response == undefined || response.data == null) {
                    alert("undefined response: comment inserted but not visible");
                    return;
                }
                let comm_obj = {
                    commentId: response.data,
                    creator: this.$user_state.username,
                    content: text,
                    date: creation_time
                }
                this.comments.push(comm_obj);
                this.commentNumber++;
            }).catch((error) => {
                if (error.response.status == 403) {
                    alert("You are not authorized to comment");
                } else if (error.response.status == 404) {
                    alert("Photo not found, cant comment on it");
                } else {
                    alert("Error: " + error.response.data);
                }
                return;
            });
        },

        async DeletePost() {
            // Update the state on the server
            await this.$axios.delete("/Images/" + this.post_data, {
                headers: {
                    'accept': 'application/json',
                    'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                }
            }).then((response) => {
                if (response === undefined || response.data == null) {
                    alert("undefined response");
                    
                }
                // Remove the post from the stream
                this.$emit("delete-post", this.post_data);
            }).catch((error) => {
                console.log("Error: ", error.response);
                if (error.response.status == 403) {
                    alert("You are not authorized to delete this photo");
                } else if (error.response.status == 404) {
                    alert("Photo not found, cant delete it");
                } else {
                    alert("Error: " + error.response.data);
                }
                return;
            });
        },

        async Like() {
            await this.$axios.put("/Images/" + this.photo_id + "/likes/" + this.$user_state.headers.Authorization, null, {
                headers: {
                    "Authorization": "Bearer " + this.$user_state.headers.Authorization,
                    "accept": "application/json",
                    "Content-Type": "application/json"
                }
            }).catch((error) => {
                if (error.response.status == 401) {
                    alert("You are not authorized to like this photo");
                } else if (error.response.status == 404) {
                    alert("Photo not found, cant like it");
                } else {
                    alert("Error: " + error.response.data);
                }
                return;
            }).then((response) => {
                if (response === undefined || response.data == null) {
                    alert("undefined response");
                    return
                }
                if (response.status == 200) {
                    alert("You have already liked this photo");
                    return;
                } else if (response.status == 201) {
                    this.likes++;
                    this.likesList.push(this.$user_state.username);
                }
                this.have_i_liked_this = true;

            })
        },

        async Unlike() {

            // Update the state on the server
            await this.$axios.delete("/Images/" + this.photo_id + "/likes/" + this.$user_state.headers.Authorization, {
                headers: {
                    'Content-Type': 'application/json',
                    'accept': 'application/json',
                    'Authorization': 'Bearer ' + this.$user_state.headers.Authorization
                }
            }).catch((error) => {
                if (error.response.status == 403) {
                    alert("You are not authorized to unlike this photo");
                } else if (error.response.status == 404) {
                    alert("Photo not found, cant unlike it");
                } else {
                    alert("Error: " + error.response.data);
                }
                return;
            }).then((response) => {
                if (response.status != 200) {
                    alert("Error: " + response.data);
                    return;
                }

                this.have_i_liked_this = false;
                this.likes--;
                this.likesList = this.likesList.filter((like) => like != this.$user_state.username);
            });
        },

        async DeleteComment(comment) {

            // Update the state on the server
            await this.$axios.delete("/Images/" + this.photo_id + "/comments/" + comment.commentId, {
                headers: {
                    'Content-Type': 'application/json',
                    'accept': 'application/json',
                    'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                }
            }).then((response) => {
                if (response === undefined || response.data == null) {
                    alert("undefined response");
                    return
                }
                if (response.status != 200) {
                    alert("Error: " + response.data);
                    return;
                }

                // Remove the comment from the array
                this.comments = this.comments.filter((c) => c.commentId != comment.commentId);
                this.commentNumber--;
            }).catch((error) => {

                if (error.response.status == 403) {
                    alert("You are not authorized to delete this comment", error.response.data);
                } else if (error.response.status == 404) {
                    alert("Comment not found, cant delete it");
                } else {
                    alert("Error: " + error.response.data);
                }
                return;
            });
        },
    },

    mounted() {

        this.initialize();

    }
}
</script>

<template>
    <!-- Bordered Wrapper :id="ccc"-->

    <div class="rounded p-2 m-4 border shadow-lg" style="width: 500px; ">
        <div class="row align-content-between my-2 justify-content-center">
            <div class="col">
                <i class="bi-person-circle mx-2" style="font-size: 2em"></i>
                <span v-if="!is_your_post" class="col font-weight-bold h1">
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
            <Photo :src="post_data" :alt="description" :style="{ width: '100%', cursor: 'pointer' }"></Photo>
        </div>

        <!-- Divider <hr class="mt-1 mb-4"> -->



        <!-- Caption -->

        <div class="row">
            <div class="col-12 post-box" style="width: 500px">
                <i class="bi-person-circle mx-1"></i>
                <span class="font-weight-bold h1 post-box" style="font-size: 1.6em;
                margin-right: 5px; word-break: break-all; overflow-wrap: break-word;">
                    {{ description }}</span>
            </div>
        </div>

        <!-- Divider -->

        <hr class="my-4">

        <!-- Comments -->
        <div class="row mt-3 align-content-start justify-content-between">
            <div class="col-auto d-flex align-items-center pb-2">
                <button type="button" class="btn btn-primary" data-bs-toggle="modal"
                    :data-bs-target="'#commList' + post_data">Comments {{ commentNumber }}</button>
            </div>

            <div class="col-auto " style="width: max-content;">
                <LikeManager class="v-center" :likes_count="this.likes" :liked="this.have_i_liked_this" @like="Like"
                    @unlike="Unlike">
                </LikeManager>
            </div>


            <div class="col-auto d-flex align-items-center pb-2">
                <button type="button" class="btn btn-primary" data-bs-toggle="modal"
                    :data-bs-target="'#likeList' + post_data">Likes</button>
            </div>

        </div>


        <!-- CommentManager -->


        <Modal :id="'commList' + post_data">
            <!---->
            <template v-slot:header>

                <div class="col-16">
                    <CommentManager id="commentwriter" :photo_id="this.photo_id"
                        :author_name="this.creatorname" @comment="AddComment">
                    </CommentManager>
                </div>

            </template>
            <template v-slot:body>
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
            </template>
        </Modal>
        <Modal :id="'likeList' + post_data">
            <!---->
            <template v-slot:header>
                <div class="col-auto d-flex align-items-center pb-2">
                    <i class="bi-chat"> Who liked this post</i>

                </div>
            </template>
            <template v-slot:body>
                <div class="row">
                    <div v-if="likes == 0" class="col-12 align-content-center w-100"><!-- Center the text -->
                        <span class="h5 mx-1 font-weight-bold align-middle text-muted text-center">No Likes yet.</span>
                    </div>
                    <div v-else class="col-12 my-3">
                        <ul class="list-group">
                            <li v-for="(like, index) in this.likesList" :key="index" class="list-group">{{ like }}</li>
                        </ul>
                    </div>
                </div>
            </template>
        </Modal>
    </div>
</template>

<style>
.v-center {

    display: inline-block;
    vertical-align: middle;
    line-height: normal;
    align-items: center;

}

.comment-button {

    padding: 0.35rem 0.5rem;
    font-size: 0.8rem;
    width: 120px;

}

.post-box {
    word-wrap: break-word;
    overflow-wrap: break-word;
    white-space: pre-line;
}
</style>
