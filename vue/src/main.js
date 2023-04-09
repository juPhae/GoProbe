import Vue from 'vue'
import App from './App.vue'
import router from './router' // 新增的路由器

Vue.config.productionTip = false

import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';

Vue.use(ElementUI);

new Vue({
  router, // 注册路由器
  render: h => h(App),
}).$mount('#app')
