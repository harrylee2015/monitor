import Vue from "vue"
import Vuex from "vuex"
import { constantRouterMap, asyncRouterMap } from '../router'
Vue.use(Vuex)

export default new Vuex.Store({
    state: {
        // HOST: 'http://192.168.0.110:8080/',
        HOST: '/proxy/',
        defaultActive: "",
        headerName: '',
        permissionData: {
            button:{}
        },
    },
    mutations: {
        updateUid(state, payload) {
            state.uid = payload;
        },
        updateDefaultActive(state, payload) {
            state.defaultActive = payload;
        },
        updateHeaderName(state, payload) {
            state.headerName = payload;
        },
        updatePermissionData(state, payload) {
            state.permissionData = payload;
        }
    }
})