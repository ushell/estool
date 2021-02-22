import 'core-js/stable';
import 'regenerator-runtime/runtime';
import Vue from 'vue';
import App from './App.vue';

import ViewUI from 'view-design';
// 引入全局样式
import 'view-design/dist/styles/iview.css';
Vue.use(ViewUI);

Vue.config.productionTip = false;
Vue.config.devtools = true;
//
// import * as Wails from '@wailsapp/runtime';
//
// Wails.Init(() => {
// 	new Vue({
// 		render: h => h(App)
// 	}).$mount('#app');
// });

// debug
new Vue({
	render: h => h(App)
}).$mount('#app');