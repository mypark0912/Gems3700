import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/store/auth'; // ✅ Pinia Store 사용
import { nextTick } from 'vue'
import Dashboard from './pages/main/DynamicDashboard.vue'
import MasterDashboard from './pages/main/MasterDashboard.vue'
import Meter from './pages/main/Meter.vue'
import PowerQ from './pages/main/PowerQ.vue'
import Diagnosis from './pages/main/Diagnosis2.vue'
import Report2 from './pages/main/Report2.vue'
import Event from './pages/main/Event.vue'
import Setting from './pages/main/Setting.vue'
import Signin from './pages/main/Signin.vue'
import Trend from './pages/main/Trend.vue'
import Signup from './pages/main/Signup.vue'
import Calibrate from './pages/main/Calibrate2.vue'
import GemsSubDashboard from './pages/main/GemsSubDashboard.vue'
const routerHistory = createWebHistory('/')

const router = createRouter({
  history: routerHistory,
  routes: [
    {
      path: "/",
      name: "home",  // ✅ "/" 경로를 별도의 네임드 라우트로 정의
      component: Dashboard, // 기본적으로 Dashboard 로드
    },
    {
      path: "/signin",
      name: "signin",
      component: Signin
    },
    {
      path: "/signup",
      name: "signup",
      component: Signup
    },
    {
      path: '/dashboard',
      name : 'dashboard',
      component: Dashboard
    },
    {
      path: '/dashboard/:user',
      name : 'dashboardUser',
      component: Dashboard,
      props:true
    }, 
    {
      path: "/master",
      name: "master",  // ✅ "/" 경로를 별도의 네임드 라우트로 정의
      component: MasterDashboard, // 기본적으로 Dashboard 로드
    },
    {
      path: '/meter/:channel',
      component: Meter,
      name : 'Meter',
      props: true  // 🔹 params를 props로 자동 전달
    },
    {
      path: '/Sub/:channel',
      component: GemsSubDashboard,
      name : 'GemsSubDashboard',
      props: true  // 🔹 params를 props로 자동 전달
    },
    {
      path: '/powerq/:channel',
      component: PowerQ,
      name : 'PowerQ',
      props: true  // 🔹 params를 props로 자동 전달
    },
    {
      path: '/diagnosis/:channel',
      component: Diagnosis,
      name : 'Diagnosis',
      props: true  // 🔹 params를 props로 자동 전달
    },
    {
      path: '/event/:channel',
      component: Event,
      name : 'Event',
      props: true  // 🔹 params를 props로 자동 전달
    },
    {
      path: '/report/:channel',
      component: Report2,
      name : 'Report',
      props: true  // 🔹 params를 props로 자동 전달
    },
    {
      path: '/trend/:channel',
      component: Trend,
      name : 'Trend',
      props: true  // 🔹 params를 props로 자동 전달
    },
    {
      path: '/settings/:mode/:channel',
      component: Setting,
      name : 'Setting',
      props: true  // 🔹 params를 props로 자동 전달
    },
    {
      path: '/config/:channel',
      component: Calibrate,
      name : 'Calibrate',
      props: true  // 🔹 params를 props로 자동 전달
    },
  ]
})

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore(); // ✅ Pinia store 사용

  if (to.name === 'dashboardUser') {
    // checkRemoteUser 호출하여 로그인 처리
    await authStore.checkRemoteUser(to.params.user);
    
    const isAuthenticated = authStore.getLogin;
    if (!isAuthenticated) {
      next("/signin");
    } else {
      next(`/dashboard`);
    }
    return; // 여기서 처리 완료
  }

  
  await authStore.checkSession(); // ✅ Vuex dispatch → Pinia action 사용
  const isAuthenticated = authStore.getLogin; // ✅ Vuex getters → Pinia computed 사용
  const opMode = authStore.getOpMode;
  // console.log(to.path);
  nextTick();
  if (to.path === "/signin" && isAuthenticated) {
    next("/dashboard");  // ✅ 로그인 상태면 대시보드로 이동
  } 
  else if (!isAuthenticated && to.path !== "/signin" && to.path !== "/signup") {
    next("/signin"); // ✅ 비로그인 상태면 로그인 페이지로 이동 (단, 회원가입 페이지는 예외)
  } 
  else if (to.path === "/") {
    //next(isAuthenticated ? "/dashboard" : "/signin");
    if(opMode !== 'server')
      next(isAuthenticated ? "/dashboard" : "/signin"); // ✅ "/" 요청 시 로그인 상태에 따라 리디렉트
    else
      next(isAuthenticated ? "/master" : "/signin"); // ✅ "/" 요청 시 로그인 상태에 따라 리디렉트
  } 
  else {
    next(); // ✅ 정상적으로 요청 진행
  }
});

export default router
