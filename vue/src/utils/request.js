import axios from 'axios'

// 创建一个 axios 实例
const request = axios.create({
  baseURL: process.env.VUE_APP_BASE_API, // 设置基础URL
  timeout: 5000 // 设置超时时间
})

// 请求拦截器
request.interceptors.request.use(
  config => {
    // 在请求发送之前可以在这里对请求做一些处理，比如添加 token 等
    return config
  },
  error => {
    // 对请求错误做一些处理
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  response => {
    // 在响应数据返回之前可以在这里对响应做一些处理，比如判断响应状态码是否为200等
    return response.data
  },
  error => {
    // 对响应错误做一些处理
    return Promise.reject(error)
  }
)

// 导出请求函数
export function post(url, data) {
  return request({
    method: 'post',
    url,
    data
  })
}

export function get(url, params) {
  return request({
    method: 'get',
    url,
    params
  })
}
