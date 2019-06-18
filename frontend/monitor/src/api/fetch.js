import axios from 'axios'
import store from '../store'


axios.defaults.headers['Content-Type'] = 'application/json';
const HOST = store.state.HOST;

export default {
    postFetch(url, data){
        return new Promise((resolve,reject)=>{
            axios({
                method:"post",
                url:HOST + url,
                data:data
            }).then((res)=>{
                resolve(res)
            }).catch((error)=>{
                reject(error)
            })
        })
    },
    getFetch(url){
        return new Promise((resolve, reject) => {
            axios({
                method:"get",
                url:HOST + url
            }).then((res)=>{
                resolve(res)
            }).catch((error)=>{
                reject(error)
            })
        })
    },
    putFetch(url, data){
        return new Promise((resolve, reject) => {
            axios({
                method:"put",
                url:HOST + url,
                data:data
            }).then((res)=>{
                resolve(res)
            }).catch((error)=>{
                reject(error)
            })
        })
    },
    deleteFetch(url){
        return new Promise((resolve, reject) => {
            axios({
                method:"delete",
                url:HOST + url,
            }).then((res)=>{
                resolve(res)
            }).catch((error)=>{
                reject(error)
            })
        })
    },
}