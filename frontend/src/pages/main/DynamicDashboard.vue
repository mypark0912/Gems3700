<template>
  <div class="flex h-[100dvh] overflow-hidden">
    <!-- Sidebar -->
    <Sidebar :sidebarOpen="sidebarOpen" @close-sidebar="sidebarOpen = false" :channel="channel" />

    <!-- Content area -->
    <div class="relative flex flex-col flex-1 overflow-y-auto overflow-x-hidden">
      <!-- Site header -->
      <Header :sidebarOpen="sidebarOpen" @toggle-sidebar="sidebarOpen = !sidebarOpen" />

      <main class="grow">
        <div class="px-2 sm:px-4 lg:px-6 py-4 w-full max-w-full">
          <!-- Dashboard actions -->
          <div v-if="opMode != 'device2'" class="sm:flex sm:justify-between sm:items-center mb-4">
            <div class="mb-4 sm:mb-0">
              <h2 class="text-xl md:text-2xl text-gray-800 dark:text-gray-100 font-bold ml-1">
                {{ t('dashboard.sitemap') }}
              </h2>
            </div>
          </div>

          <!-- 로딩 상태 표시 -->
          <div v-if="isLoading" class="flex justify-center items-center h-64">
            <div class="text-lg text-gray-600">데이터 로딩 중...</div>
          </div>

          <!-- Dynamic Dashboard Layout - 데이터 준비 완료 후에만 렌더링 -->
          <template v-else-if="isDataReady">
            <IntegratedLayout v-if="dashboardLayout === 'DualChannelLayout' && opMode === 'device2'"
              :main-data="meterDictMain"
              :sub-data="meterDictSub"
              :channel-state="ChannelState"
              :channel="channel"
            />
            <SingleChannelLayout v-else-if="dashboardLayout === 'SingleChannelLayout'"
              :main-data="meterDictMain"
              :sub-data="meterDictSub"
              :channel-state="ChannelState"
              :channel="channel"
            />
            <DualChannelLayout v-else-if="dashboardLayout === 'DualChannelLayout'"
              :main-data="meterDictMain"
              :sub-data="meterDictSub"
              :channel-state="ChannelState"
              :channel="channel"
            />
            <GemsLayout v-else               
            :main-data="meterDictMain"
              :channel-state="ChannelState"
              :channel="channel"
              />
          </template>

          <!-- 데이터 로딩 실패 또는 조건 불만족 -->
          <div v-else class="flex justify-center items-center h-64">
            <div class="text-lg text-gray-500">Dashboard Preparation in Progress...</div>
          </div>
        </div>
      </main>
      <Footer />
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed, onUnmounted, provide, watch } from 'vue'
import Sidebar from '../common/SideBarGems.vue'
import Header from '../common/Header.vue'
import SingleChannelLayout from '../layouts/SingleChannelLayout.vue'
import DualChannelLayout from '../layouts/DualChannelLayout.vue' 
import IntegratedLayout from '../layouts/IntegratedLayout.vue'  
import GemsLayout from '../layouts/GemsLayout2.vue'
import Footer from "../common/Footer.vue"
import axios from 'axios'
import { useSetupStore } from '@/store/setup'
import { useAuthStore } from '@/store/auth'
import { useI18n } from 'vue-i18n'


export default {
  name: 'Dashboard',
  props: ['user'],
  components: {
    Sidebar,
    Header,
    Footer,
    SingleChannelLayout,
    DualChannelLayout,
    IntegratedLayout,
    GemsLayout,
  },
  setup(props) {
    const sidebarOpen = ref(false)
    const user = ref(props.user)
    const setupStore = useSetupStore()
    const authStore = useAuthStore()
    const meterDictMain = ref({})
    const meterDictSub = ref({})
    const { t } = useI18n()
    const channel = ref('main')
    
    // 로딩 상태 관리
    const isLoading = ref(true)
    const isMainDataLoaded = ref(false)
    const isSubDataLoaded = ref(false)
    
    const ChannelState = computed(() => setupStore.getChannelSetting)
    const setup = computed(() => setupStore.getSetup)
    const opMode = computed(() => authStore.getOpMode)
    let updateInterval = null

    // 동적 레이아웃 선택
    const dashboardLayout = computed(() => {
      return 'GemsLayout'
      // if (opMode.value === 'device2') {
      //   return 'DualChannelLayout'
      // } else {
      //   return ChannelState.value?.SubEnable ? 'DualChannelLayout' : 'SingleChannelLayout'
      // }
    })

    // 데이터 준비 완료 상태 확인
    const isDataReady = computed(() => {
      if (!setup.value || !opMode.value || opMode.value === '') {
        return false
      }

      if (opMode.value === 'device2') {
        return isMainDataLoaded.value
      } else {
        const needsSubData = ChannelState.value?.SubEnable
        return isMainDataLoaded.value && (needsSubData ? isSubDataLoaded.value : true)
      }
    })

    // 데이터 페칭 함수
    const fetch9000Data = async (ch) => {
      if (opMode.value === '') {
        //console.log(`❌ ${ch} 채널: opMode가 빈 문자열`)
        return false
      }

      // console.log(`🚀 ${ch} 채널 데이터 페칭 시작`, {
      //   opMode: opMode.value,
      //   url: `/api/getMeterRedisNew/${ch}/${opMode.value}`
      // })
      
      try {
        const response = await axios.get(`/api/getMeterRedisNew/${ch}/${opMode.value}`)
        
        // console.log(`📡 ${ch} 채널 응답 받음:`, {
        //   status: response.status,
        //   success: response.data?.success,
        //   dataExists: !!response.data?.data,
        //   channel: response.data?.channel
        // })
        
        if (response.data.success) {
          //console.log(response.data);
          if (ch === 'Main') {
              meterDictMain.value = { ...meterDictMain.value, ...response.data.data }
              if (!ChannelState.value?.SubEnable && meterDictMain.value.P4) {
                meterDictMain.value.P4 = (meterDictMain.value.P4 / 1000).toFixed(2)
              }
              isMainDataLoaded.value = true
              //console.log(`✅ ${ch} 채널 데이터 로딩 완료 (Main)`)
            } else {
              meterDictSub.value = { ...meterDictSub.value, ...response.data.data }
              if (meterDictSub.value.P4) {
                meterDictSub.value.P4 = (meterDictSub.value.P4 / 1000).toFixed(2)
              }
              isSubDataLoaded.value = true
              //console.log(`✅ ${ch} 채널 데이터 로딩 완료 (Sub)`)
            }
          // if (opMode.value === 'device2') {
          //   meterDictMain.value = { ...meterDictMain.value, ...response.data.data }
          //   if (meterDictMain.value.P4) {
          //     meterDictMain.value.P4 = (meterDictMain.value.P4 / 1000).toFixed(2)
          //   }
          //   channel.value = response.data.channel === 'Main' ? 'main' : 'sub'
          //   isMainDataLoaded.value = true
          //   console.log(`✅ ${ch} 채널 데이터 로딩 완료 (device2)`)
          // } else {
          //   if (ch === 'Main') {
          //     meterDictMain.value = { ...meterDictMain.value, ...response.data.data }
          //     if (!ChannelState.value?.SubEnable && meterDictMain.value.P4) {
          //       meterDictMain.value.P4 = (meterDictMain.value.P4 / 1000).toFixed(2)
          //     }
          //     isMainDataLoaded.value = true
          //     console.log(`✅ ${ch} 채널 데이터 로딩 완료 (Main)`)
          //   } else {
          //     meterDictSub.value = { ...meterDictSub.value, ...response.data.data }
          //     if (meterDictSub.value.P4) {
          //       meterDictSub.value.P4 = (meterDictSub.value.P4 / 1000).toFixed(2)
          //     }
          //     isSubDataLoaded.value = true
          //     console.log(`✅ ${ch} 채널 데이터 로딩 완료 (Sub)`)
          //   }
          // }
          
          // console.log(`🎯 현재 로딩 상태:`, {
          //   isMainDataLoaded: isMainDataLoaded.value,
          //   isSubDataLoaded: isSubDataLoaded.value,
          //   isDataReady: isDataReady.value
          // })
          
          return true
        } else {
          console.log(`❌ ${ch} 채널: 서버에서 success=false 응답`)
          return false
        }
      } catch (error) {
        console.log(`❌ ${ch} 채널 데이터 가져오기 실패:`, error)
        return false
      }
    }

    // 초기 데이터 로딩
    const loadInitialData = async () => {
      if (!setup.value || !opMode.value || opMode.value === '') {
        console.log('❌ 초기 조건 불만족으로 데이터 로딩 중단')
        return
      }

      console.log('🔄 초기 데이터 로딩 시작')
      isLoading.value = true
      
      try {
        if (opMode.value === 'device2') {
          await fetch9000Data('Main')
        } else {
          await fetch9000Data('Main')
          if (ChannelState.value?.SubEnable) {
            await fetch9000Data('Sub')
          } else {
            isSubDataLoaded.value = true // Sub 데이터가 필요없으면 완료로 표시
          }
        }
      } catch (error) {
        console.error('초기 데이터 로딩 실패:', error)
      } finally {
        isLoading.value = false
      }
    }

    // 주기적 데이터 업데이트 설정
    const setupPeriodicUpdate = () => {
      if (updateInterval) {
        clearInterval(updateInterval)
      }

      updateInterval = setInterval(async () => {
        if (!setup.value || !opMode.value) return

          await fetch9000Data('Main')
          if (ChannelState.value?.SubEnable) {
            await fetch9000Data('Sub')
          }

        // if (opMode.value === 'device2') {
        //   await fetch9000Data('Main')
        // } else {
        //   await fetch9000Data('Main')
        //   if (ChannelState.value?.SubEnable) {
        //     await fetch9000Data('Sub')
        //   }
        // }
      }, 10000)
    }

    // Provide 데이터
    provide('meterDictMain', meterDictMain)
    provide('meterDictSub', meterDictSub)

    // 채널 상태 변화 감시
    watch(() => ChannelState.value, async (newVal) => {
      console.log('🔄 ChannelState 변경됨:', newVal)
      
      if (!setup.value) {
        if (updateInterval) {
          clearInterval(updateInterval)
          updateInterval = null
        }
        return
      }

      // 데이터 상태 초기화
      isMainDataLoaded.value = false
      isSubDataLoaded.value = false
      
      // 데이터 다시 로딩
      await loadInitialData()
      
      // 주기적 업데이트 재설정
      setupPeriodicUpdate()
    }, { immediate: true })

    // setup 값 변화 감시
    watch(() => setup.value, async (newSetup) => {
      console.log('🔄 Setup 변경됨:', newSetup)
      
      if (newSetup && opMode.value && opMode.value !== '') {
        await loadInitialData()
        setupPeriodicUpdate()
      }
    })

    onMounted(async () => {
      console.log('=== Dashboard onMounted 시작 ===')
      
      if (user.value != null) {
        authStore.setUser(user.value)
        authStore.setLogin(true)
      }
      
      //console.log('checkSetting 호출 전')
      await setupStore.checkSetting()
      //console.log('checkSetting 호출 후')
      
      // console.log('현재 상태:', {
      //   setup: setup.value,
      //   opMode: opMode.value,
      //   ChannelState: ChannelState.value
      // })
      
      // 초기 데이터 로딩은 watch에서 처리됨 (immediate: true)
      console.log('=== Dashboard onMounted 완료 ===')
    })

    onUnmounted(() => {
      if (updateInterval) {
        clearInterval(updateInterval)
      }
    })

    return {
      sidebarOpen,
      meterDictMain,
      meterDictSub,
      setup,
      user,
      ChannelState,
      dashboardLayout,
      t,
      opMode,
      channel,
      isLoading,
      isDataReady,
    }
  }
}
</script>