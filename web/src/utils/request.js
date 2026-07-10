import axios from 'axios';
import { Message } from '@arco-design/web-vue';
import { useUserStore } from '../store/userInfo'

// create an axios instance
const service = axios.create({
  baseURL: import.meta.env.VITE_APP_BASE_URL,
  timeout: 65000, // 查询可能较慢,与后端 60s 查询超时对齐
});

// request interceptor
service.interceptors.request.use(
  (config) => {
    // Store 必须在拦截器内部导入，在外部导入会显示 Pinia 未初始化
    const store = useUserStore();
    // 设置请求头部 Authorization
    if (store.token) {
      config.headers['Authorization'] = 'Bearer ' + store.token;
      config.headers['Content-Type'] = 'application/json'
    }
    return config;
  },
  (error) => {
    console.error(error);
    return Promise.reject(error);
  }
);

// response interceptor
service.interceptors.response.use(
  (response) => {
    // 后端恒返回 HTTP 200,业务码在 body.code;成功/业务错误都从这里返回 body,
    // 由调用方判 code(便于把后端 msg 显示出来)。
    return response.data;
  },
  (error) => {
    // 主动中断(AbortController)不弹错误,交由调用方处理为「已中断」
    if (axios.isCancel(error)) {
      return Promise.reject({ canceled: true });
    }
    // 仅网络层错误(超时/断网/CORS 等)走这里。此时可能没有 response,
    // 不能直接读 error.response.data,否则拦截器自身会抛 TypeError 吞掉真实错误。
    const store = useUserStore();
    const resp = error.response;
    const body = resp && resp.data ? resp.data : null;

    if (body && body.code === 401) {
      Message.error({ content: 'Token 已过期, 请重新登陆', duration: 3000 });
      store.userLogout();
      window.location.href = '/login';
      return Promise.reject('登录已过期');
    }

    // 优先用后端返回的 msg;没有(纯网络/超时)则用 axios 的错误信息
    const msg = (body && body.msg) ? body.msg : (error.message || '请求失败');
    Message.error({ content: msg, duration: 4000 });
    return Promise.reject(msg);
  }
);

export default service;
