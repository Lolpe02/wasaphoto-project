<script>
export default {

    props: {
        posts: {
            type: Array
        }
    },

    watch: {
        posts: function (new_posts) {
            this.posts_ = new_posts;
        }
    },

    data: function () {
        return {
            posts_: this.posts
        }
    },

    emits: ["delete-post"],

    methods: {

        async initialize() {

        },

        async RemovePost(post_data) {

            // Bubble up the event, let the parent handle it
            // Not optimal design but its what we are stuck with

            this.$emit("delete-post", post_data);
        }

    },
    mounted() {

        this.initialize();

    }
}
</script>

<template>
    <div class="container d-flex justify-content-center align-items-center" style="display: flex; flex-direction: column; align-content: center; ">
        <div v-for="post in posts_" :key="post" class="m-1">
            <WASAPost :post_data="post" @delete-post="RemovePost" />
        </div>
    </div>
</template>

<style>


</style>