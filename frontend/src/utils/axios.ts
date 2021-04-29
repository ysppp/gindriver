import axios from 'axios'

axios.interceptors.request.use(function (config) {
  // 这里的config包含每次请求的内容
  const token = window.sessionStorage.getItem('jwt')
  if (token) {
    // 添加headers
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
}, function (error) {
  // 对请求错误做些什么
  return Promise.reject(error);
});

export default axios