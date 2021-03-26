import Vue from 'vue'
import App from './App.vue'
import store from './store'

Vue.config.productionTip = false

Vue.mixin({
  methods: {
    // this function returns an object that is suitable to be put into
    // the fetch api as options
    httpOptions: function (m, obj) {
      // this object will be returned
      var options = {
        method: m !== undefined ? m : 'GET',
        mode: 'cors',
        cache: 'no-cache',
        credentials: 'same-origin',
        headers: {
          'Content-Type': 'application/json'
        },
        redirect: 'follow'
      }
      if (obj !== undefined) {
        options.body = JSON.stringify(obj)
      }
      return options
    }
  }
})

new Vue({
  store,
  render: h => h(App)
}).$mount('#app')
