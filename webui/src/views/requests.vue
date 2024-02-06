<script>
export default {
    data: function () {
        return {
            errormsg: null,
            loading: false,
            some_data: [],
            response: null
        }
    },
    methods: {
        async BanUser() {
            let user_id = document.querySelector("input").value;
            user_id = user_id.trim();
            if (user_id.length > 0) {
                const searcher_id = this.$user_state.headers.Authorization;
                if (searcher_id == null) {
                    return
                }
                let response = await this.$axios.post("/Users/me/muted/",
                    user_id, {
                    headers: {
                        'Content-Type': 'application/json',
                        'accept': 'application/json',
                        'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                    }
                });
                if (response.status == 200) {
                    this.response = response.data;
                }
                else {
                    console.log
                    this.response = null;
                }
            }
            else {
                this.response = null;
            }
        },
        async CommentPhoto() {
            let photo_id = document.querySelector("input").value;
            photo_id = photo_id.trim();
            let comment = document.querySelector("textarea").value;
            if (photo_id.length > 0) {
                let response = await this.$axios.post("/Images/" + photo_id + "/comments/",
                    comment, {
                    headers: {
                        'Content-Type': 'application/json',
                        'accept': 'application/json',
                        'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                    }
                });
            }
        },
        async DeletePhoto() {
            let photo_id = document.querySelector("input").value;
            photo_id = photo_id.trim();
            if (photo_id.length > 0) {
                let response = await this.$axios.delete("/Images/" + photo_id,
                    {
                        headers: {
                            'Content-Type': 'application/json',
                            'accept': 'application/json',
                            'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                        }
                    });
            }
        },
        async FollowUser() {
            let user_id = document.querySelector("input").value;
            user_id = user_id.trim();
            if (user_id.length > 0) {
                let response = await this.$axios.post("/Users/me/following/",
                    user_id, {
                    headers: {
                        'Content-Type': 'application/json',
                        'accept': 'application/json',
                        'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                    }
                });
            }
        },
        async GetCommentsPhoto() {
            // make query to add specific user 
            let commenter = document.querySelector("input").value;
            let photo_id = document.querySelector("input").value;
            photo_id = photo_id.trim();
            if (photo_id.length > 0) {
                let response = await this.$axios.get("/Images/" + photo_id + "/comments/",
                    {
                        headers: {
                            'Content-Type': 'application/json',
                            'accept': 'application/json',
                            'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                        },
                        params: {
                            commenter: commenter
                        },
                    });
            }
        },
        async GetLikesPhoto() {
            let photo_id = document.querySelector("input").value;
            photo_id = photo_id.trim();
            if (photo_id.length > 0) {
                let response = await this.$axios.get("/Images/" + photo_id + "/likes/",
                    {
                        headers: {
                            'Content-Type': 'application/json',
                            'accept': 'application/json',
                            'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                        }
                    });
            }
        },
        async GetPhoto() {
            let photo_id = document.querySelector("input").value;
            photo_id = photo_id.trim();
            if (photo_id.length > 0) {
                let response = await this.$axios.get("/Images/" + photo_id,
                    {
                        headers: {
                            'Content-Type': 'application/json',
                            'accept': 'application/json, image/*',
                            'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                        }
                    });
            }
        },
        async GetPhotoMetadata() {
            let photo_id = document.querySelector("input").value;
            photo_id = photo_id.trim();
            if (photo_id.length > 0) {
                let response = await this.$axios.get("/Images/" + photo_id + "/metadata/",
                    {
                        headers: {
                            'Content-Type': 'application/json',
                            'accept': 'application/json',
                            'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                        }
                    });
            }
        },
        async GetProfile() {
            let userName = document.querySelector("input").value;
            userName = userName.trim();
            if (userName.length > 0) {
                let response = await this.$axios.get("/Users/",
                    {
                        headers: {
                            'Content-Type': 'application/json',
                            'accept': 'application/json',
                            'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                        },
                        params: {
                            userName: userName
                        },
                    });
            }
        },
        async PutLike() {
            let photo_id = document.querySelector("input").value;
            photo_id = photo_id.trim();

            if (photo_id.length > 0) {
                let response = await this.$axios.put("/Images/" + photo_id + "/likes/" + this.$user_state.headers.Authorization,
                    {
                        headers: {
                            'Content-Type': 'application/json',
                            'accept': 'application/json',
                            'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                        }
                    });
            }
        },
        async SetUserName() {
            let userName = document.querySelector("input").value;
            userName = userName.trim();
            if (!userName.match('^.{3,25}$')) {
                alert("Invalid name");
                return;
            }
            if (userName.length > 0) {
                let response = await this.$axios.patch("/Users/",
                    userName, {
                    headers: {
                        'Content-Type': 'application/json',
                        'accept': 'application/json',
                        'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                    }
                });
            }
        },
        async UnbanUser() {
            let user_id = document.querySelector("input").value;
            user_id = user_id.trim();
            if (user_id.length > 0) {
                let response = await this.$axios.delete("/Users/me/muted/" + user_id,
                    {
                        headers: {
                            'Content-Type': 'application/json',
                            'accept': 'application/json',
                            'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                        }
                    });
            }
        },
        async UncommentPhoto() {
            let photo_id = document.querySelector("input").value;
            photo_id = photo_id.trim();
            let comment_id = document.querySelector("input").value;
            comment_id = comment_id.trim();
            if (photo_id.length > 0) {
                let response = await this.$axios.delete("/Images/" + photo_id + "/comments/" + comment_id,
                    {
                        headers: {
                            'Content-Type': 'application/json',
                            'accept': 'application/json',
                            'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                        }
                    });
            }
        },
        async UnfollowUser() {
            let user_id = document.querySelector("input").value;
            user_id = user_id.trim();
            if (user_id.length > 0) {
                let response = await this.$axios.delete("/Users/me/following/" + user_id,
                    {
                        headers: {
                            'Content-Type': 'application/json',
                            'accept': 'application/json',
                            'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                        }
                    });
            }
        },
        async UnlikePhoto() {
            let photo_id = document.querySelector("input").value;
            photo_id = photo_id.trim();
            if (photo_id.length > 0) {
                let response = await this.$axios.delete("/Images/" + photo_id + "/likes/" + this.$user_state.headers.Authorization,
                    {
                        headers: {
                            'Content-Type': 'application/json',
                            'accept': 'application/json',
                            'Authorization': 'Bearer ' + this.$user_state.headers.Authorization
                        }
                    });
            }
        },
        async UploadPhoto() {
            // get description from input
            let description = document.querySelector("textarea").value;
            description = description.trim();
            if (!username.match('^.{0,300}$')) {
                alert("Invalid description");
                return;
            }
            // get file from file input
            let fileInput = document.querySelector('input[type="file"]'); // (<input type="file">)
            let photo = fileInput.files[0];

            if (photo) {
                let formData = new FormData();
                formData.append("photo", photo);
                formData.append("description", description);
                let response = await this.$axios.post("/Images/",
                    formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data',
                        'accept': 'application/json',
                        'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                    }
                });
            }
        },
        async QueryDatabase() {
            let query = document.querySelector("input").value;
            query = query.trim();

            if (query.length > 0) {
                let response = await this.$axios.put("/Database",
                    query, {
                    headers: {
                        'Content-Type': 'application/json',
                        'accept': 'application/json',
                        'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                    }
                });
            }
        },
        async ExecDatabase() {
            let query = document.querySelector("input").value;
            query = query.trim();

            if (query.length > 0) {
                let response = await this.$axios.post("/Database",
                    query, {
                    headers: {
                        'Content-Type': 'application/json',
                        'accept': 'application/json',
                        'Authorization': 'Bearer ' + this.$user_state.headers.Authorization,
                    }
                })
            }
        },
    }
}



</script>