<script>
export default {
    data: function () {
        return {
            error: false,
            loading: false,
            some_data: null,
            isAuthenticated: false,
            $user_state: this.$user_state,
        }
    },
    methods: {
        async initialize() {

            this.$user_state.current_view = this.$views.LOGIN;
            this.$user_state.username = null;
            this.$user_state.headers.Authorization = null;

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
            await this.$axios.post("/session",
                username, {
                headers: {
                    'Content-Type': 'application/json',
                    'accept': 'application/json',
                },
            }).then((response) => {
                if (response === undefined || response.data == null) {
                    alert("undefined response");
                    return
                }
                 if (response.status == 201) {
                    // new user created
                    alert("Welcome to the community " + username + "!");
                } else if (response.status == 200) {
                    // user already exists
                    alert("Welcome back " + username + "!");
                }
                this.$user_state.headers.Authorization = response.data
                this.$user_state.username = username
                this.isAuthenticated = true;
                // localStorage.setItem("userToken", JSON.stringify(response.data));
                this.error = false;
                this.$router.push("/");
            }).catch((error) => {
                if (error.response.status == 400) {
                    alert("Invalid username");
                    return;
                }
                this.error = true;
                this.$user_state.headers.Authorization = null
                console.log("Error logging in");
                this.initialize();
                return;
            });
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

<style></style>