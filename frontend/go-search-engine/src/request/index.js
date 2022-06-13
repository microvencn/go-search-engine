import axios from 'axios'


// create an axios instance
const service = axios.create({
    baseURL: 'http://localhost:8000/api/v1',
    // withCredentials: true,
    timeout: 5000 // request timeout
})

// request interceptor
service.interceptors.request.use(
    // config => {
    //     if (store.getters.token) {
    //         config.headers['Authorization'] = 'Bearer ' + getToken()
    //     }
    //     config.headers['X-Requested-With'] = 'XMLHttpRequest'
    //     return config
    // },
    // error => {
    //     // do something with request error
    //     console.log(error) // for debug
    //     return Promise.reject(error)
    // }
)

// response interceptor
service.interceptors.response.use(
    response => {
        return response.data
    },
    error => {
        // this.$message.error("请求失败请重试")
        return Promise.reject(error)
    }
)

export default service

export const picRequest = axios.create({
    // baseURL: 'http://43.138.211.125:9999',
    baseURL: 'http://110.42.237.111:9999',
    withCredentials: true,
    timeout: 5000 // request timeout
})

