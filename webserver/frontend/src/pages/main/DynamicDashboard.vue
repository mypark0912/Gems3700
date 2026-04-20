<template>
    <div class="flex h-[100dvh] overflow-hidden">
      <!-- Sidebar -->
      <Sidebar :sidebarOpen="sidebarOpen" @close-sidebar="sidebarOpen = false" />
  
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
            <!-- <div v-if="realtimeStore.isLoading" class="flex justify-center items-center h-64">
              <div class="text-lg text-gray-600">{{ t('dashboard.loading') }}...</div>
            </div> -->
            <div v-if="realtimeStore.isLoading" class="col-span-12 flex justify-center items-center h-64">
            <div class="flex flex-col items-center gap-3">
              <svg class="animate-spin h-10 w-10 text-blue-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <span class="text-gray-600 dark:text-gray-400">{{ t('dashboard.loading') }}...</span>
            </div>
          </div>
            <!-- 에러 상태 표시 -->
            <div v-else-if="realtimeStore.error" class="flex justify-center items-center h-64">
              <div class="text-lg text-red-600">
                Data loading failed: {{ realtimeStore.error.message }}
                <button @click="retryLoading" class="ml-4 px-4 py-2 bg-blue-500 text-white rounded">
                  Retry
                </button>
              </div>
            </div>
  
            <!-- Dynamic Dashboard Layout - 데이터 준비 완료 후에만 렌더링 -->
            <template v-else-if="isDataReady">
              <idpmSingleChannelLayout
                v-if="viewMode === 'classic'"
                :channel-state="ChannelState"
                :channel="channel"
              />
              <idpmLayout
                v-else
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
  import { ref, onMounted, computed, onUnmounted, watch } from 'vue'
  import { storeToRefs } from 'pinia'
  import Sidebar from '../common/SideBar.vue'
  import Header from '../common/Header.vue'
  import SingleChannelLayout from '../layouts/SingleChannelLayout.vue'
  import idpmSingleChannelLayout from '../layouts/idpmSingleChannelLayout.vue'
  import idpmLayout from '../layouts/idpmLayout.vue'
  import DualChannelLayout from '../layouts/DualChannelLayout.vue' 
  import IntegratedLayout from '../layouts/IntegratedLayout3.vue'  
  import Footer from "../common/Footer.vue"
  import { useSetupStore } from '@/store/setup'
  import { useAuthStore } from '@/store/auth'
  import { useRealtimeStore } from '@/store/realtime'
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
      idpmSingleChannelLayout,
      idpmLayout,
    },
    setup(props) {
      const sidebarOpen = ref(false)
      const user = ref(props.user)
      const { t } = useI18n()
      const channel = ref('main')
      
      // Store 인스턴스
      const setupStore = useSetupStore()
      const authStore = useAuthStore()
      const realtimeStore = useRealtimeStore()
      const isInitialized = ref(false)
      // Computed from stores
      const ChannelState = computed(() => setupStore.getChannelSetting)
      const setup = computed(() => setupStore.getSetup)
      const opMode = computed(() => authStore.getOpMode)
      const viewMode = computed(() => authStore.getViewMode)
     
      // 데이터 준비 완료 상태 확인
      const isDataReady = computed(() => {
        if (!setup.value) return false
        return realtimeStore.isDataReady(ChannelState.value, opMode.value)
      })
  
      // 동적 레이아웃 선택
      const dashboardLayout = computed(() => {
        if (opMode.value === 'device2') {
          return 'DualChannelLayout'
        } else {
          if (ChannelState.value?.SubEnable && ChannelState.value?.MainEnable)
            return 'DualChannelLayout'
          else{
            if (ChannelState.value?.MainEnable)
              channel.value = 'main'
            else
              channel.value = 'sub'
            
            return 'SingleChannelLayout'
          }

          //return ChannelState.value?.SubEnable ? 'DualChannelLayout' : 'SingleChannelLayout'
        }
      })
  
      // 재시도 함수
      const retryLoading = () => {
        realtimeStore.startPolling(opMode.value, ChannelState.value)
      }
  
      // 채널 상태 변화 감시
      // watch(() => ChannelState.value, async (newVal) => {
      //   console.log('🔄 ChannelState 변경됨:', newVal)
        
      //   if (!setup.value) {
      //     realtimeStore.stopPolling()
      //     return
      //   }
  
      //   // Store에서 채널 상태 변경 처리
      //   await realtimeStore.onChannelStateChange(opMode.value, newVal)
      // }, { immediate: true })
      watch(() => ChannelState.value, async (newVal) => {
        console.log('🔄 ChannelState 변경됨:', newVal)
        
        if (!setup.value || !isInitialized.value) {
          return
        }

        // 실제 변경 시에만 재시작
        await realtimeStore.onChannelStateChange(opMode.value, newVal)
      })
  
      // setup 값 변화 감시
      // watch(() => setup.value, async (newSetup) => {
      //   console.log('🔄 Setup 변경됨:', newSetup)
        
      //   if (newSetup && opMode.value && opMode.value !== '') {
      //     realtimeStore.startPolling(opMode.value, ChannelState.value)
      //   }
      // })
      watch(() => setup.value, async (newSetup) => {
        console.log('🔄 Setup 변경됨:', newSetup)
        
        if (newSetup && opMode.value && opMode.value !== '' && !isInitialized.value) {
          isInitialized.value = true
          realtimeStore.startPolling(opMode.value, ChannelState.value)
        }
      });
  
      onMounted(async () => {
        console.log('=== Dashboard onMounted 시작 ===')
        
        if (user.value != null) {
          authStore.setUser(user.value)
          authStore.setLogin(true)
        }
        
        await setupStore.checkSetting()
        
        // 초기 데이터 로딩은 watch에서 처리됨 (immediate: true)
        if (opMode.value && opMode.value !== '') {
          isInitialized.value = true
          realtimeStore.startPolling(opMode.value, ChannelState.value)
        }
        console.log('=== Dashboard onMounted 완료 ===')
      })
  
      onUnmounted(() => {
        // 컴포넌트 언마운트 시 폴링 중지
        realtimeStore.stopPolling()
      })
  
      return {
        sidebarOpen,
        setup,
        user,
        ChannelState,
        dashboardLayout,
        t,
        opMode,
        channel,
        realtimeStore,
        isDataReady,
        retryLoading,
        isInitialized,
        idpmSingleChannelLayout,
        idpmLayout,
        viewMode,
      }
    }
  }
  </script>