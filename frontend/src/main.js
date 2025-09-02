import { createApp } from 'vue';
import { createPinia } from 'pinia'; // ✅ Pinia 추가
import router from './router2';
import App from './App.vue';
import i18n from './i18n'
//import { useAuthStore } from '@/store/auth'; // ✅ Pinia Store 사용
// import vSelect from 'vue-select';
// import 'vue-select/dist/vue-select.css';
import './css/style.css';

const pinia = createPinia(); // ✅ Pinia 인스턴스 생성

const app = createApp(App);
app.use(pinia); // ✅ Vuex 대신 Pinia 등록
app.use(router);
app.use(i18n);
app.mount('#app');

// const app = createApp(App);
// const pinia = createPinia();

// app.use(pinia);  // ✅ Pinia 먼저 연결
// app.use(router);
// app.use(i18n);

// // ✅ Pinia가 연결된 후에 store 사용 가능
// import { useSetupStore } from '@/store/setup';
// const setupStore = useSetupStore();

// setupStore.checkSetting().then(() => {
//   app.mount('#app'); // ✅ store 세팅 완료 후 mount
// });
