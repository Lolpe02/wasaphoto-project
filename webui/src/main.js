import {createApp, reactive} from 'vue'
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';
import ErrorMsg from './components/ErrorMsg.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'
import Stream from './components/Stream.vue'
import Photo from './components/Photo.vue'
import WASAPost from './components/WASAPost.vue'
import Modal from './components/Modal.vue'
import Comment from './components/Comment.vue'
import LikeManager from './components/LikeManager.vue'
import CommentManager from './components/CommentManager.vue'

import './assets/dashboard.css'
import './assets/main.css'

const app = createApp(App)

var state = {
    headers: {
        Authorization: null,
        accept: "application/json",
    },
    username: null,
    viewing: null
}
const views = {
    LOGIN: "login",
    PROFILE: "profile",
    STREAM: "stream",
}

app.config.globalProperties.$axios = axios;
app.config.globalProperties.$views = views;
app.config.globalProperties.$user_state = reactive(state || {}); // Initialize user_state with an empty object if it is undefined
app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);
app.component("Modal", Modal);
app.component("Stream", Stream);
app.component("Photo", Photo);
app.component("WASAPost", WASAPost);
app.component("Comment", Comment);
app.component("LikeManager", LikeManager);
app.component("CommentManager", CommentManager);

app.use(router)
app.mount('#app')

/*









const months = [
    "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"
]


app.config.globalProperties.$months = months;

app.config.globalProperties.$hasher = (password) => {
    return sha256(password);
}

app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);


app.use(router)
app.mount('#app')
*/