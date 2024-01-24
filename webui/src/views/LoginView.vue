<script>
export default {
    data: function () {
        return {
            error: false,
            loading: false,
            some_data: null,
        }
    },
    methods: {
        async initialize() {

            this.$user_state.current_view = this.$views.LOGIN;
            this.$user_state.username = null;

        },
        async login() {

            let username = document.getElementById("login-form").value;

            // check username regex

            if (!username.match('^.{3,25}$')) {
                alert("Invalid username");
                return;
            }

            /*
            try {
                const jsonData = { key: 'value' }; // Replace this with your JSON object

                // Using Axios to send a POST request with JSON data
                const response = await axios.post('your-api-endpoint', jsonString, );

                console.log('Response:', response.data);
            } 
            */
            try {
                const jsonString = JSON.stringify(username);
                
                let response = await this.$axios.post("/session", 
                    username, {
                        headers: {
                            'Content-Type': 'application/json',
                            'accept': 'application/json',
                        },
                    }
                );
                console.log('Response:', response.data);
            } catch (error) {
                this.errormsg = error.toString();
                console.error('Error sending request:', error);
            }
            //check if the response is 201

            if (response.status == 201) {
                // get the response body               
                this.$user_state.headers.Authorization = response.data
                // new user created
                alert("Welcome to the community " + username + "!");

            } else if (response.status == 200) {
                // user already exists
                this.$user_state.headers.Authorization = response.data
                // pop up welcome message
                alert("Welcome back " + username + "!");

            } else {
                this.error = true;
                this.$user_state.headers.Authorization = null
                return;
            }
            this.$user_state.username = username
            this.isAuthenticated = true;
            localStorage.setItem("user_state", JSON.stringify(this.$user_state));
            this.error = false;
            this.$router.push("/home");

        }
    },
    mounted() {
        this.initialize()
    }
}
</script>

<template>
    <div class="container text-center pt-3 pb-2 border-bottom">
        <h2>Login</h2>
    </div>


    <div class="h-75 d-flex align-items-center justify-content-center">
        <form class="border border-dark p-5 rounded shadow-lg">
            <!-- Username input -->
            <div class="form-outline mb-4">
                <input type="text" id="login-form" class="form-control" pattern="^.{3,25}$" />
                <label class="form-label" for="login-form">Username</label>
            </div>

            <!-- Submit button -->
            <button type="button" class="btn btn-primary btn-block mb-4" @click="login()">Sign in</button>

        </form>

    </div>


</template>

<style>

</style>