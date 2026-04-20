<template>
  <div class="relative inline-flex">
    <button
      ref="trigger"
      class="w-8 h-8 flex items-center justify-center hover:bg-gray-100 lg:hover:bg-gray-200 dark:hover:bg-gray-700/50 dark:lg:hover:bg-gray-800 rounded-full"
      :class="{ 'bg-gray-200 dark:bg-gray-800': dropdownOpen }"
      aria-haspopup="true"
      @click.prevent="dropdownOpen = !dropdownOpen"
      :aria-expanded="dropdownOpen"
    >
      <span class="sr-only">Device Time</span>
      <!-- 시계 아이콘 -->
      <svg class="fill-current text-gray-500/80 dark:text-gray-400/80" width="16" height="16" viewBox="0 0 16 16" xmlns="http://www.w3.org/2000/svg">
        <path fill-rule="evenodd" d="M8 16A8 8 0 1 0 8 0a8 8 0 0 0 0 16Zm6-8A6 6 0 1 1 2 8a6 6 0 0 1 12 0Z" />
        <path d="M8 4a.75.75 0 0 1 .75.75v2.69l1.78 1.78a.75.75 0 1 1-1.06 1.06l-2-2A.75.75 0 0 1 7.25 8V4.75A.75.75 0 0 1 8 4Z" />
      </svg>
      <!-- 상태 표시 띵 -->
      <div 
        class="absolute top-0 right-0 w-2.5 h-2.5 border-2 border-gray-100 dark:border-gray-900 rounded-full" 
        :class="status ? 'bg-green-500' : 'bg-red-500'"
      ></div>
    </button>
    
    <transition
      enter-active-class="transition ease-out duration-200 transform"
      enter-from-class="opacity-0 -translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition ease-out duration-200"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div 
        v-show="dropdownOpen" 
        ref="dropdown"
        class="origin-top-right z-10 absolute top-full min-w-64 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 py-1.5 rounded-lg shadow-lg overflow-hidden mt-1" 
        :class="align === 'right' ? 'right-0' : 'left-0'"
      >
        <div class="text-xs font-semibold text-gray-400 dark:text-gray-500 uppercase pt-1.5 pb-2 px-4">
          Device Time
        </div>
        
        <div class="px-4 py-3">
          <!-- 장비 시간 + 새로고침 버튼 -->
          <div class="mb-4 flex items-center justify-between">
            <div class="text-sm font-medium text-gray-800 dark:text-gray-100">
              {{ formattedDeviceTime }}
            </div>
            <button
              @click="refreshTime"
              :disabled="refreshing"
              class="p-1.5 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-full transition-colors"
              title="Refresh"
            >
              <svg 
                class="w-4 h-4 text-gray-500 dark:text-gray-400"
                :class="{ 'animate-spin': refreshing }"
                fill="none" 
                stroke="currentColor" 
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
            </button>
          </div>
          
          <!-- 시간 동기화 버튼 -->
          <button
            @click="syncTime"
            :disabled="syncing"
            class="w-full py-2 px-4 bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white text-sm font-medium rounded-lg transition-colors duration-150 flex items-center justify-center gap-2"
          >
            <svg 
              v-if="syncing" 
              class="animate-spin h-4 w-4" 
              xmlns="http://www.w3.org/2000/svg" 
              fill="none" 
              viewBox="0 0 24 24"
            >
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <svg 
              v-else 
              class="w-4 h-4" 
              fill="none" 
              stroke="currentColor" 
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            {{ syncing ? 'Setting...' : 'Set Time' }}
          </button>
          
          <!-- 동기화 결과 메시지 -->
          <div v-if="syncMessage" class="mt-2 text-xs text-center" :class="syncSuccess ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'">
            {{ syncMessage }}
          </div>
        </div>
      </div> 
    </transition>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import axios from 'axios'

export default {
  name: 'DropdownDeviceTime',
  props: {
    align: {
      type: String,
      default: 'right'
    },
    status: {
      type: Boolean,
      default: true
    },
    deviceTime: {
      type: String,
      default: null
    }
  },
  emits: ['time-synced', 'sync-error', 'refresh-time'],
  setup(props, { emit }) {
    const dropdownOpen = ref(false)
    const trigger = ref(null)
    const dropdown = ref(null)
    
    const syncing = ref(false)
    const refreshing = ref(false)
    const syncMessage = ref('')
    const syncSuccess = ref(false)

    // 장비 시간 포맷팅
    const formattedDeviceTime = computed(() => {
      if (!props.deviceTime) return '--'
      const date = new Date(props.deviceTime)
      return date.toLocaleString('ko-KR', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false
      })
    })

    // 시간 새로고침
    const refreshTime = async () => {
      refreshing.value = true
      emit('refresh-time')
      setTimeout(() => {
        refreshing.value = false
      }, 500)
    }

    // 시간 동기화
    const syncTime = async () => {
      syncing.value = true
      syncMessage.value = ''
      
      try {
        const now = new Date()
        const datetimeStr = now.getFullYear() + '-' +
          String(now.getMonth() + 1).padStart(2, '0') + '-' +
          String(now.getDate()).padStart(2, '0') + ' ' +
          String(now.getHours()).padStart(2, '0') + ':' +
          String(now.getMinutes()).padStart(2, '0') + ':' +
          String(now.getSeconds()).padStart(2, '0')
        
        const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone
        
        const response = await axios.post('/config/calibrate/setSystemTime', {
          datetime_str: datetimeStr,
          timezone: timezone
        })
        
        if (response.data.success) {
          syncMessage.value = 'Success!'
          syncSuccess.value = true
          emit('time-synced')
          refreshTime()
        } else {
          syncMessage.value = 'Failed'
          syncSuccess.value = false
          emit('sync-error', response.data)
        }
      } catch (err) {
        syncMessage.value = 'Failed'
        syncSuccess.value = false
        emit('sync-error', err)
        console.error('Error syncing time:', err)
      } finally {
        syncing.value = false
        setTimeout(() => { syncMessage.value = '' }, 3000)
      }
    }

    // 클릭 외부 닫기
    const clickHandler = ({ target }) => {
      if (!dropdownOpen.value) return
      if (trigger.value?.contains(target)) return
      if (dropdown.value?.contains(target)) return
      dropdownOpen.value = false
    }

    // ESC 키 닫기
    const keyHandler = ({ keyCode }) => {
      if (!dropdownOpen.value || keyCode !== 27) return
      dropdownOpen.value = false
    }

    onMounted(() => {
      document.addEventListener('click', clickHandler)
      document.addEventListener('keydown', keyHandler)
    })

    onUnmounted(() => {
      document.removeEventListener('click', clickHandler)
      document.removeEventListener('keydown', keyHandler)
    })

    return {
      dropdownOpen,
      trigger,
      dropdown,
      syncing,
      refreshing,
      syncMessage,
      syncSuccess,
      formattedDeviceTime,
      syncTime,
      refreshTime,
    }
  }
}
</script>