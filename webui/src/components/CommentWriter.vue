<script>
export default {
    props: ['photo_id', 'author_name'],

    emits: ['comment'],

    data: function () {
        return {
            author: this.author_name,
            photo_id_: this.photo_id,
            curr_time: null,
            intervalId: null,
        }
    },

    methods: {
        async initialize() {

            this.format_date_now();
            this.intervalId = setInterval(this.format_date_now, 60000);
        },

        async format_date_now() {
            // Get date in RFC3339 format
            let now = new Date().toISOString();
            // format to dd month yyyy at hh:mm
            let hour = new Date().getHours();
            let date_split = now.split("T");
            let date = date_split[0].split("-");
            let time = date_split[1].split(":");
            time = hour + ":" + time[1];
            date = date[2] + "/" + date[1] + "/" + date[0] + " at " + time;
            this.curr_time = date;
        },

        async add_comment() {
            let text = document.getElementById("comment").value;
            this.$emit('comment', text);

            // Clear the text area

            document.getElementById("comment").value = "";
        },
        beforeDestroy() {
            clearInterval(this.intervalId);
        }
    },


    mounted() {
        this.initialize();
    }
}
</script>

<template>
    <div class="row mt-3 align-content-start justify-content-between">

        <p class="card-text col-8 d-flex">
            <textarea class="form-control p-1 textarea-width" id="comment" rows="3" placeholder="Leave a Comment!"></textarea>
        </p>
        <button type="button" class="btn mx-2 btn-primary col-2" @click="add_comment"
            style="font-size: 0.8em, font-style: italic;">Comment</button>
    </div>
</template>

<style>
.textarea-width {
    width: 500px;
}
</style>