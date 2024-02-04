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
            likeslist: [],
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
            }).catch((error) => {
                if (error.response.status == 403) {
                    alert("You are not authorized to comment");
                } else if (error.response.status == 404) {
                    alert("Photo not found, cant comment on it");
                } else {
                    alert("Error: " + error.response.data);
                }
                return;
            }).then((response) => {
                if (response.data == null || response == undefined) {
                    alert("Error: " + response.data);
                    return;
                } else {
                    let comm_obj = {
                        commentId: response.data,
                        creator: this.$user_state.username,
                        content: text,
                        date: creation_time
                    }
                    this.comments.push(comm_obj);
                }
            });
        },

        async DeletePost() {

            // Update the state on the server
            let response = await this.$axios.delete("/Images/" + this.photo_id,
                    {
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

            // Remove the post from the stream
            this.$emit("delete-post", this.post_data);

        },

        async Like() {
            await this.$axios.put("/Images/" + this.photo_id + "/likes/" + this.$user_state.headers.Authorization, null, {
                headers: {
                    "Authorization": "Bearer " + this.$user_state.headers.Authorization,   
                    "accept": "application/json",
                    "Content-Type": "application/json"
                }
            }).catch((error) => {
                if (error.response.status == 403) {
                    alert("You are not authorized to like this photo");
                } else if (error.response.status == 404) {
                    alert("Photo not found, cant like it");
                } else {
                    alert("Error: " + error.response.data);
                }
                return;
            }).then((response) => {
                if (response.status != 200 && response.status != 201) {
                    alert("200 or 201: " + response.data);
                    return;
                } else if (response.status == 200) {
                    alert("You have already liked this photo");
                    return;
                } else {
                    this.have_i_liked_this = true;
                    this.likes++;
                }
            });
        },

        async Unlike() {

            // Update the state on the server
            let response = await this.$axios.delete("/Images/" + this.photo_id + "/likes/" + this.$user_state.headers.Authorization, {
                    headers: {
                    'Content-Type': 'application/json',
                    'accept': 'application/json',
                    'Authorization': 'Bearer ' + this.$user_state.headers.Authorization
                    }
                });

            if (response.status != 200) {
                return;
            }

            this.have_i_liked_this = false;
            this.likes--;
        },

        async DeleteComment(comment) {

            // Update the state on the server
            let response = await this.$axios.delete("/Images/" + this.photo_id + "/comments/" + comment.commentId,
                    {
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

            // Remove the comment from the array
            this.comments = this.comments.filter((c) => c.commentId != comment.commentId);
        },
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
            <Photo :src="post_data" :alt="description" :style="{ width: '100%', cursor: 'pointer' }"></Photo>
        </div>

        <!-- Divider <hr class="mt-1 mb-4"> -->

        

        <!-- Caption -->

        <div class="row">
            <div class="col-12" style="width: 500px;">
                <i class="bi-person-circle mx-1" style="font-size: 1.5em"></i>
                <span class="font-weight-bold h1" style="margin-right: 5px; width: 500px; overflow-y: auto;"> {{ description }}</span>
            </div>
        </div>

        <!-- Divider -->

        <hr class="my-4">

        <!-- Comments -->
        <div class="row mt-3 align-content-start justify-content-between">
            <div class="col-auto d-flex align-items-center pb-2">
                <button type="button" class="btn btn-primary" data-bs-toggle="modal" data-bs-target="#commList">Comments</button>
            </div>
            
            <div class="col-auto " style="width: max-content;">
                <LikeCounter class="v-center" :likes_count="this.likes" :liked="this.have_i_liked_this" @like="Like"
                    @unlike="Unlike">
                </LikeCounter>
            </div>
            

            <div class="col-auto d-flex align-items-center pb-2">
                <button type="button" class="btn btn-primary" data-bs-toggle="modal" data-bs-target="#likeList">Likes</button>
            </div>

        </div>
            
        
        <!-- CommentWriter -->
        
        
        <Modal id="commList" >
            <!---->
            <template v-slot:header>
                
                    <div class="col-16">
                        <CommentWriter id="comment-writer" :photo_id="this.photo_id" :author_name="this.creatorname" @comment="AddComment">
                        </CommentWriter>
                    </div>
                
            </template>
            <template  v-slot:body>
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
        <Modal id="likeList">
            <!---->
            <template v-slot:header>
                <div class="col-auto d-flex align-items-center pb-2">
                    <i class="bi-chat"> Who liked this post</i>
                    
                </div>
            </template>
            <template  v-slot:body>
                <div class="row">
                    <div v-if="likes == 0" class="col-12 align-content-center w-100"><!-- Center the text -->
                        <span class="h5 mx-1 font-weight-bold align-middle text-muted text-center">No Likes yet.</span>
                    </div>
                    <div v-else class="col-12 my-3">
                        <Comment v-for="comment in comments" :comment="comment" :key="comment.commentId"
                            @delete="DeleteComment">
                        </Comment>
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

}

.comment-button {

    padding: 0.35rem 0.5rem;
    font-size: 0.8rem;
    width: 120px;

}
</style>
