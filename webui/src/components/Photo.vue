<script>


export default {

    props: {
        src: {
            type: Number
        },
        alt: {
            type: String
        },
        style: {
            type: Object
        }
    },
    data: function () {
        return {
            is_loading: false,
            src_: null,
            alt_: null,
            style_: null,
        }
    },

    methods: {

        async initialize() {

            this.is_loading = true
            this.alt_ = this.alt;
            this.style_ = this.style;

            // Get the image data from the server

            await this.$axios.get("/Images/" + this.src, {
                responseType: "arraybuffer",
                headers : {
                    "Authorization": "Bearer " + this.$user_state.headers.Authorization,
                    "accept": "image/*",     
                    "Content-Type": "application/json"
                }   
            }).catch((error) => {
                if (error.response.status == 401) {
                    this.$router.push("/login");
                } else if (error.response.status == 404) {
                    alert("Image not found");
                }
                this.is_loading = false;
                this.src_ = null;
                return;
            }).then((response) => {
                const blob = new Blob([response.data], {
                    type: response.headers["Content-Type"]
                });
                console.log(response); //.headers["Content-Type"]
                this.src_ = URL.createObjectURL(blob)
                this.is_loading = false;
            });
        }

    },
    mounted() {

        // Here we need to perform a sequence of async operations through axios 
        // to get the image data from the server. We need to get the image data
        // likes, comments, etc...

        this.initialize();

    }
}
</script>

<template>

    <div v-if="!is_loading">
        <img :src="src_" :alt="(alt != null ? alt : 'WPimage')" class="shadow-sm rounded img-fluid opacity-100"
            :style="(style_ != null ? style_ : '')" />
    </div>
    <div v-else>
        <LoadingSpinner></LoadingSpinner>
    </div>

</template>

<style>

</style>