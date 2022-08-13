import Vue from 'vue'
import App from './App.vue'
import router from './router'
import vuetify from './plugins/vuetify';
import confirm from 'vuetify-confirm'
import title from './mixins/title'
import '@mdi/font/css/materialdesignicons.css'

Vue.config.productionTip = false;
Vue.use(confirm, {
  vuetify,
  icon: '',
  property: '$confirm',
});

Vue.mixin(title)


new Vue({
  router,
  vuetify,
  render: h => h(App)
}).$mount('#app')
