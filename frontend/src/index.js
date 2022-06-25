import Vue from 'vue';
import App from './App';

// this imports the javascript part of botstrap, including jquery:
import 'bootstrap';

// also get fontawsome:
import '@fortawesome/fontawesome-free/css/all.css'

// this imports our styles, including the scss part of bootstrap:
import './styles.scss'



new Vue({
  el: '#app',
  render: h => h(App),
});

