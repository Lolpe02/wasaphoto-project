import {createApp, reactive} from 'vue'
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';
import ErrorMsg from './components/ErrorMsg.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'
import Stream from './components/Stream.vue'
import Photo from './components/Photo.vue'
import PhotoPost from './components/PhotoPost.vue'
import Modal from './components/Modal.vue'
import Comment from './components/Comment.vue'
import LikeCounter from './components/LikeCounter.vue'
import CommentWriter from './components/CommentWriter.vue'
import requests from './views/requests.vue'

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
    FEED: "register",
    PROFILE: "profile",
}

app.config.globalProperties.$axios = axios;
app.config.globalProperties.$views = views;
app.config.globalProperties.$user_state = reactive(state);
app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);
app.component("requests", requests);
app.component("Modal", Modal);
app.component("Stream", Stream);
app.component("Photo", Photo);
app.component("PhotoPost", PhotoPost);
app.component("Comment", Comment);
app.component("LikeCounter", LikeCounter);
app.component("CommentWriter", CommentWriter);

app.use(router)
app.mount('#app')

/*



import './assets/dashboard.css'
import './assets/main.css'

const app = createApp(App)





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