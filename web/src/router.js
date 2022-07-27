import Vue from 'vue'
import Router from 'vue-router'
import ListMocks from "./components/ListMocks.vue";
import MockDetails from "./components/MockDetails.vue"
import Help from "./components/Help.vue"

Vue.use(Router)

export default new Router({
    routes: [
        {
            path: '/',
            name: "ListMocks",
            component: ListMocks,
        },
        {
            path: '/details/:theKey',
            name: "MockDetails",
            component: MockDetails,
            props:true
        },
        {
            path: '/new',
            name: "NewMock",
            component: MockDetails,
        },
        {
            path: '/help',
            name: "Help",
            component: Help,
        },
    ]
});