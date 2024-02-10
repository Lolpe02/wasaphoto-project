<script>
export default {
    data() {
        return {
            query: '',
            password: '',
            mode: 'Query',
            response: '',
        };
    },
    methods: {
        toggleMode() {
            this.mode = this.mode === 'Exec' ? 'Query' : 'Exec';
        },
        async submitForm() {
            // Here you would handle form submission, e.g., using Axios to make an API request
            // For demonstration purposes, just show the entered username and password
            switch (this.mode) {
                case 'Exec':
                    await this.$axios.post("/Database/",
                        this.query, {
                        headers: {
                            "Authorization": 'Bearer ' + this.password,
                            "accept": "application/json",
                            "Content-Type": "application/json",
                        },
                    }).then(response => {
                        if (response == undefined || response.data == null) {
                            this.response = "undefined response";
                            return;
                        } else {
                            this.response = response.data;
                        }
                    }).catch((err) => {
                        if (err.response.status == 403) {
                            alert("wrong password");

                        } else {
                            alert("Error querying " + err.response.data);

                        }
                        this.password = '';
                        this.response = '';
                        return;
                    });
                    break;
                case 'Query':
                    await this.$axios.put("/Database/",
                        this.query, {
                        headers: {
                            "Authorization": 'Bearer ' + this.password,
                            "accept": "application/json",
                            "Content-Type": "application/json",
                        },
                    }).then(response => {
                        if (response == undefined || response.data == null) {
                            this.response = [];
                            return;
                        } else {
                            this.response = response.data;
                        }
                    }).catch((err) => {
                        if (err.response.status == 403) {
                            alert("wrong password");

                        } else {
                            alert("Error querying " + err.response.data);

                        }
                        this.password = '';
                        this.response = '';
                        return;
                    });
            }

        },
    },
    beforeMount() { // this is a hack to make sure the user is logged in, beforeCreate is not working
		if (this.$user_state.username == null) {
			console.log("Empty username before mount, redirecting to login")
			this.$router.push("/login");
		}
	},
};
</script>
  
  
<template>
    <div class="form-container">
        <form @submit.prevent="submitForm" class="form">
            <div class="form-group">
                <label for="yourquery" class="form-label">Insert query</label>
                <input type="text" id="yourquery" v-model="query" class="form-input" />
            </div>
            <div class="form-group">
                <label for="password" class="form-label">Password:</label>
                <input type="password" id="password" v-model="password" class="form-input" />
            </div>
            <div class="form-group">
                <div class="button-group">
                    <button type="button" @click="toggleMode" class="form-button">{{ mode }}</button>
                    <button type="submit" class="form-button1">Submit</button>
                </div>
            </div>
        </form>
        <div class="response-container">
            <p class="response">{{ response }}</p>
        </div>
    </div>
</template>

<style scoped>
.form-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    margin-top: 50px;
}

.form {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 10px;
}

.form-group {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 5px;
}

.form-label {
    font-weight: bold;
}

.form-input {
    padding: 5px;
    border: 1px solid #ccc;
    border-radius: 5px;
}

.form-button {
    padding: 5px 10px;
    background-color: #007bff;
    color: #fff;
    border: none;
    border-radius: 5px;
    cursor: pointer;
}
.form-button1 {
    width: 10px auto;
    padding: 15px 20px;
    background-color: #ff1500;
    color: #fff;
    border: none;
    border-radius: 10px;
    cursor: pointer;
}
.response-container {
    margin-top: 20px;
}

.response {
    font-weight: bold;
}
</style>
  

  
