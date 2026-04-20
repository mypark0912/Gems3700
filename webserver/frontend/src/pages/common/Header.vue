<template>
  <header
    class="sticky top-0 before:absolute before:inset-0 before:backdrop-blur-md max-lg:before:bg-white/90 dark:max-lg:before:bg-gray-800/90 before:-z-10 z-30"
    :class="[
      variant === 'v2' || variant === 'v3' ? 'before:bg-white after:absolute after:h-px after:inset-x-0 after:top-full after:bg-gray-200 dark:after:bg-gray-700/60 after:-z-10' : 'max-lg:shadow-sm lg:before:bg-gray-100/90 dark:lg:before:bg-gray-900/90',
      variant === 'v2' ? 'dark:before:bg-gray-800' : '',
      variant === 'v3' ? 'dark:before:bg-gray-900' : '',
    ]"
  >
    <div class="px-4 sm:px-6 lg:px-8">
      <div
        class="flex items-center justify-between h-16"
        :class="variant === 'v2' || variant === 'v3' ? '' : 'lg:border-b border-gray-200 dark:border-gray-700/60'"
      >

        <!-- Header: Left side -->
        <div class="flex">
<div class="flex items-center gap-1">
  <svg v-if="location !==''" 
    xmlns="http://www.w3.org/2000/svg"
    width="24"
    height="24"
    viewBox="0 0 16 16"
    fill="none"
    stroke-linecap="round"
    stroke-linejoin="round"
  >
    <!-- 바깥 물방울 모양 -->
    <path
      d="M11.771 11.771l-2.829 2.829a1.333 1.333 0 0 1 -1.885 0l-2.829 -2.829a5.333 5.333 0 1 1 7.543 0z"
      fill="#007aff"
    />
    <!-- 안쪽 원 -->
    <path
      d="M6 7.333a2 2 0 1 0 4 0a2 2 0 0 0 -4 0"
      stroke="#ffffff"
      stroke-width="1"
      fill="none"
    />
  </svg>
  <span class="text-sm leading-[32px]">{{ location }}</span>
</div>

          <!-- Hamburger button -->
          <button class="text-gray-500 hover:text-gray-600 dark:hover:text-gray-400 lg:hidden" @click.stop="$emit('toggle-sidebar')" aria-controls="sidebar" :aria-expanded="sidebarOpen">
            <span class="sr-only">Open sidebar</span>
            <svg class="w-6 h-6 fill-current" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <rect x="4" y="5" width="16" height="2" />
              <rect x="4" y="11" width="16" height="2" />
              <rect x="4" y="17" width="16" height="2" />
            </svg>
          </button>

        </div>

        <!-- Header: Right side -->
        <div class="flex items-center space-x-3">
          <!-- View Mode (fixed to design) -->
          <Notifications align="right" :status="sysIcon" :data="sysData" />
          <DeviceTime 
            align="right" 
            :status="timeStatus" 
            :device-time="deviceTime"
            @time-synced="onTimeSynced" @refresh-time="onTimeSynced"
          />
          <hr class="w-px h-6 bg-gray-200 dark:bg-gray-700/60 border-none" />
          <Help align="right" />
          <DropdownLanguage align="right" />
          <ThemeToggle />
          <!-- Divider -->
          <hr class="w-px h-6 bg-gray-200 dark:bg-gray-700/60 border-none" />
          <UserMenu align="right" />

        </div>

      </div>
    </div>
  </header>
</template>

<script>
import { computed, onMounted, ref, onUnmounted} from 'vue'
import ThemeToggle from '../header/ThemeToggle.vue'
import Help from '../header/DropdownHelp.vue'
import UserMenu from '../header/DropdownProfile.vue'
import DropdownLanguage from '../header/DropdownLanguage.vue'
import Notifications from '../header/DropdownNotifications.vue'
import DeviceTime from '../header/DropdownClock.vue'
import { useSetupStore } from "@/store/setup";
import { useAuthStore } from "@/store/auth";
import axios from 'axios'

export default {
  name: 'Header',
  props: [
    'sidebarOpen',
    'variant'
  ],
  components: {
    Help,
    UserMenu,
    ThemeToggle,
    DropdownLanguage,
    Notifications,
    DeviceTime,
  },
  setup() {
    const setupStore = useSetupStore();
    const authStore = useAuthStore();
    const viewMode = computed(() => authStore.getViewMode);
    const setViewMode = (mode) => authStore.setViewMode(mode);
    const sysIcon = ref(false);
    const sysData = ref(null);
    //const smartData = ref(null);
    const location = computed(()=> setupStore.getDevLocation);
    const channelSetting = computed(()=> setupStore.getChannelSetting);
    let updateInterval = null;
    const deviceTime = ref(null);
    const timeStatus = ref(true);  // true: 정상, false: 24시간 이상 차이

    const getSysStatus = async()=>{
      try {
        const response = await axios.get(`/api/getSystemStatus`);
        sysIcon.value = response.data.status;
        if (!response.data.status) {
          sysData.value = response.data.services;
          
        } else {
          // 정상일 때도 services 데이터 전달하려면
          sysData.value = response.data.services || {};
        }
        // if (channelSetting.value.MainDiagnosis || channelSetting.value.SubDiagnosis){
        //   if(response.data.smartStatus)
        //     smartData.value = response.data.smartStatus;
        // }
        //onsole.log(smartData.value);
      }catch (error) {
        console.log("데이터 가져오기 실패:", error);
      } 
    }

    const getDeviceTime = async () => {
      try {
        const now = new Date();
        const datetimeStr = now.getFullYear() + '-' +
          String(now.getMonth() + 1).padStart(2, '0') + '-' +
          String(now.getDate()).padStart(2, '0') + ' ' +
          String(now.getHours()).padStart(2, '0') + ':' +
          String(now.getMinutes()).padStart(2, '0') + ':' +
          String(now.getSeconds()).padStart(2, '0');
        const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
        
        const response = await axios.post('/config/checktime', {
          datetime_str: datetimeStr,
          timezone: timezone
        });
        if (response.data.success) {
          deviceTime.value = response.data.deviceTime;
          timeStatus.value = response.data.status;
        }
      } catch (error) {
        console.log("장비 시간 가져오기 실패:", error);
        timeStatus.value = false;
      }
    }

    // ✅ 시간 동기화 완료 핸들러
    const onTimeSynced = async () => {
      //console.log('시간 동기화 완료!');
      await getDeviceTime();  // 동기화 후 시간 다시 가져오기
    }

    onMounted(async () => {
      // 기존 인터벌 정리
      if (updateInterval) {
        clearInterval(updateInterval)
        updateInterval = null
      }
      
      // 초기 데이터 로드
      await getSysStatus();
      await getDeviceTime();
      
      // ✅ async 함수로 수정
      updateInterval = setInterval(async () => {
        await getSysStatus()
      }, 60 * 60 * 1000)  // 1시간마다
    })

    // ✅ cleanup 추가
    onUnmounted(() => {
      if (updateInterval) {
        clearInterval(updateInterval)
        updateInterval = null
      }
    })
    // ✅ langset이 null일 경우 기본 객체 설정
    // const langs = reactive(props.langset || {});
    // console.log("langs.value",langs.value);
    // ✅ props.langset이 변경될 때마다 langs 업데이트 (반응형 유지)
    // watch(() => setupStore.getDevLocation, (newVal) => {
    //   console.log(newVal);
    //   location.value = newVal;
    // },{immediate:true});

    return {
      location,
      sysIcon,
      sysData,
      deviceTime,      // ✅ 추가
      timeStatus,      // ✅ 추가
      onTimeSynced,    // ✅ 추가
      //smartData,
      channelSetting,
      viewMode,
      setViewMode,
    }
  }
}

</script>

<style scoped>
.view-mode-toggle {
  @apply flex items-center bg-gray-100 dark:bg-gray-700 rounded-lg p-0.5;
}
.view-mode-btn {
  @apply px-3 py-1 text-xs font-medium rounded-md transition-all duration-200;
  @apply text-gray-500 dark:text-gray-400;
}
.view-mode-btn.active {
  @apply text-white bg-blue-600 dark:bg-blue-500 shadow-sm;
}
</style>