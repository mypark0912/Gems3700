// 📁 stores/auth.js
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'
import router from '@/router'
import { useSetupStore } from './setup'

export const useAuthStore = defineStore('auth', () => {
  const setupStore = useSetupStore();
  const user = ref(localStorage.getItem('user') || '')
  const userRole = ref(localStorage.getItem('userRole') || '')
  const logined = ref(localStorage.getItem('logined') === 'true')
  const sessionChecked = ref(false);
  const installed = ref(0);
  const opMode = ref(localStorage.getItem('opMode') || '')
  const langDefault = ref(localStorage.getItem('langDefault') || '')
  const viewMode = ref('design')

  const initialState = {
    user: '',
    userRole: '',
    logined: false,
    sessionChecked: false,
    installed: 0,
    opMode: '',
    langDefault: ''
  }

  function setUser(newUser, newRole) {
    user.value = newUser
    userRole.value = newRole
  }

  function setLogin(status) {
    logined.value = status
  }

  function setSessionChecked(status) {
    sessionChecked.value = status
  }

  function setInstall(status) {
    installed.value = status
  }

  function setOpMode(mode) {
    opMode.value = mode
  }

  function setLangDefault(lang){
    langDefault.value = lang
  }

  function setViewMode(mode) {
    viewMode.value = mode
    localStorage.setItem('viewMode', mode)
  }

  function toggleViewMode() {
    const newMode = viewMode.value === 'classic' ? 'design' : 'classic'
    setViewMode(newMode)
  }

  // 수정된 login 함수 - 에러 코드도 반환
  async function login(data) {
    try {
      const response = await axios.post('/auth/checkLogins', data, {
        headers: { 'Content-Type': 'application/json' },
        withCredentials: true,
      })
      
      const passOK = response.data?.passOK
      const success = parseInt(passOK)
      
      // passOK 값에 따른 처리
      // 0: 초기화 필요, 1: 성공, 2: API 로그인 실패, 3: API 계정 없음, 4: 비번 오류, 5: 계정 오류
      if (success === 1) {
        // 로그인 성공
        setUser(data.account, response.data.data.userRole)
        setLogin(true)
        localStorage.setItem('user', data.account)
        localStorage.setItem('userRole', response.data.data.userRole)
        localStorage.setItem('logined', 'true')
        localStorage.setItem('opMode', response.data.mode)
        setOpMode(response.data.mode);
        
        return {
          passOK: true, 
          mode: response.data.mode,
          errorCode: null
        }
      } else {
        // 로그인 실패
        setUser('', '')
        setLogin(false)
        localStorage.clear()
        setOpMode('')
        
        return {
          passOK: false, 
          mode: '',
          errorCode: passOK // 에러 코드 전달 (4: 비번오류, 5: 계정오류)
        }
      }
    } catch (err) {
      console.error('Login Error:', err)
      setUser('', '')
      setOpMode('')
      setLogin(false)
      localStorage.clear()
      
      return {
        passOK: false, 
        mode: '',
        errorCode: 'network' // 네트워크 오류
      }
    }
  }

  async function logout() {
    try {
      await axios.get('/auth/logout', { withCredentials: true })
      setUser('', '')
      setLogin(false)
      setSessionChecked(false)
      localStorage.clear()
      resetAndLogout()
      router.push('/signin')
    } catch (err) {
      console.error('Logout Error:', err)
    }
  }

  async function checkSession() {
    try {
      const response = await axios.get('/auth/checkSession', { withCredentials: true })
      if (response.data?.loggedIn) {
        setUser(response.data.username, response.data.userRole)
        setLogin(true)
        localStorage.setItem('user', response.data.username)
        localStorage.setItem('userRole', response.data.role)
        localStorage.setItem('logined', 'true')
        setOpMode(response.data.mode);
      } else {
        setUser('', '')
        setLogin(false)
        localStorage.clear()
        setOpMode('')
      }
    } catch (err) {
      console.error('Session Check Error:', err)
      setUser('', '')
      setLogin(false)
      localStorage.clear()
      setOpMode('')
    }
  }

  async function checkInstalled() {
    try {
      const response = await axios.get('/auth/checkInstall')
      //console.log(response.data.result);
      setInstall(response.data.result)
      if (response.data.result > 0) 
        setupStore.setCalib(response.data.calibration)
    } catch (err) {
      console.error('checkInstalled Error:', err)
    }
  }

  async function checkRemoteUser(username) {
    try {
      const response = await axios.get(`/auth/checkRemote/${username}`)
      
      if (response.data.success) {
        // data.account가 아닌 response.data.data.account 사용해야 함
        const account = response.data.data.account || response.data.account
        
        setUser(account, response.data.data.userRole)
        setLogin(true)
        localStorage.setItem('user', account)
        localStorage.setItem('userRole', response.data.data.userRole)
        localStorage.setItem('logined', 'true')
        localStorage.setItem('opMode', response.data.data.mode)
        setOpMode(response.data.data.mode)
        
        // 사용자 정보 반환
        return {
          success: true,
          user: account,
          userRole: response.data.data.userRole
        }
      }else
        return { success: false }
    } catch (err) {
      console.error('checkRemoteUser Error:', err)
      return { success: false, error: err }
    }
  }

  function $reset() {
    user.value = initialState.user
    userRole.value = initialState.userRole
    logined.value = initialState.logined
    sessionChecked.value = initialState.sessionChecked
    installed.value = initialState.installed
    opMode.value = initialState.opMode
    langDefault.value = initialState.langDefault
    
    // localStorage도 함께 초기화
    localStorage.clear()
    
    console.log('✅ AuthStore reset complete')
  }

  async function resetAndLogout() {
    try {
      // 1. 서버에 로그아웃 요청
      await axios.get('/auth/logout', { withCredentials: true })
      
      // 2. AuthStore 초기화
      $reset()
      
      // 3. SetupStore도 초기화 (있다면)
      if (setupStore.$reset) {
        setupStore.$reset()
      }
      
      // 4. 모든 저장소 초기화
      localStorage.clear()
      sessionStorage.clear()
      
      // 5. 로그인 페이지로 이동
      router.push('/signin')
      
      console.log('🔄 Complete reset and logout')
      
    } catch (err) {
      console.error('Reset and Logout Error:', err)
      // 에러가 발생해도 클라이언트 데이터는 초기화
      $reset()
      if (setupStore.$reset) {
        setupStore.$reset()
      }
      localStorage.clear()
      sessionStorage.clear()
      router.push('/signin')
    }
  }

  const getLogin = computed(() => logined.value)
  const getUser = computed(() => user.value)
  const getUserRole = computed(() => userRole.value)
  const getInstalled = computed(() => installed.value)
  const getOpMode = computed(() => opMode.value)
  const getLangDefault = computed(()=> langDefault.value)
  const getViewMode = computed(() => viewMode.value)

  return {
    user,
    userRole,
    logined,
    sessionChecked,
    installed,
    opMode,
    langDefault,
    setUser,
    setLogin,
    setSessionChecked,
    setInstall,
    login,
    logout,
    checkSession,
    checkInstalled,
    checkRemoteUser,
    setOpMode,
    setLangDefault,
    getLogin,
    getUser,
    getUserRole,
    getInstalled,
    getOpMode,
    getLangDefault,
    getViewMode,
    viewMode,
    setViewMode,
    toggleViewMode,
    resetAndLogout,
    $reset,
  }
})