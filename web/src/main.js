import { createApp } from 'vue';
import { createPinia } from 'pinia';
import App from './App.vue';
import ArcoVue from '@arco-design/web-vue';
import { Message, Modal, Notification } from '@arco-design/web-vue';
import '@arco-themes/vue-go-admin/css/arco.css';
import '@/utils/messageCopy'; // 错误弹条追加「复制」按钮
import router from './router/';
import { parseTime } from '@/utils/parseTime';

// Directive
import permission from '@/directive/permission/permission';

// 引入 Arco 图标库
import * as ArcoIconModules from '@arco-design/web-vue/es/icon';
// 引入 RemixIcon 图标库
import * as RemixIcon from "@remixicon/vue";

// 恢复上次保存的主题(暗色),避免刷新/登录后回到亮色
if (localStorage.getItem('wiki_theme') === 'dark') {
  document.body.setAttribute('arco-theme', 'dark');
}

console.log(import.meta.env);

// Initialize the Pinia instance
const pinia = createPinia();
const app = createApp(App);

app.directive('has', permission.checkPermission);

// 挂载全局变量
app.config.globalProperties.message = Message;
app.config.globalProperties.modal = Modal;
app.config.globalProperties.notification = Notification;
app.config.globalProperties.parseTime = parseTime;

// 挂载全局图标
for(const name in ArcoIconModules){
  app.component(name,ArcoIconModules[name])
}

for(const name in RemixIcon){
  app.component(name,RemixIcon[name])
}

app.use(ArcoVue);
app.use(router);
app.use(pinia);
app.mount('#app');

