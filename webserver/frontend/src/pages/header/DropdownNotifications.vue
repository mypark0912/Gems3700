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
      <span class="sr-only">Service Status</span>
      <svg class="fill-current text-gray-500/80 dark:text-gray-400/80" width="16" height="16" viewBox="0 0 16 16" xmlns="http://www.w3.org/2000/svg">
        <path d="M9 7.5a1 1 0 1 0-2 0v4a1 1 0 1 0 2 0v-4ZM9 4.5a1 1 0 1 1-2 0 1 1 0 0 1 2 0Z" />
        <path fill-rule="evenodd" d="M8 16A8 8 0 1 0 8 0a8 8 0 0 0 0 16Zm6-8A6 6 0 1 1 2 8a6 6 0 0 1 12 0Z" />
      </svg>
      <!-- ✅ :class 로 수정 -->
      <div 
        class="absolute top-0 right-0 w-2.5 h-2.5 border-2 border-gray-100 dark:border-gray-900 rounded-full" 
        :class="stIconCss"
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
        class="origin-top-right z-10 absolute top-full -mr-48 sm:mr-0 min-w-80 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 py-1.5 rounded-lg shadow-lg overflow-hidden mt-1" 
        :class="[
          align === 'right' ? 'right-0' : 'left-0',
          inactiveServicesCount === 0 ? 'min-w-60' : 'min-w-80'
        ]"
      >
        <div class="text-xs font-semibold text-gray-400 dark:text-gray-500 uppercase pt-1.5 pb-2 px-4">
          System Status
        </div>
        <ul
            ref="dropdown"
            @focusin="dropdownOpen = true"
            @focusout="dropdownOpen = false"
          >
            <li 
              v-for="service in inactiveServices" 
              :key="service.name"
              class="border-b border-gray-200 dark:border-gray-700/60 last:border-0"
            >
              <router-link 
                class="block py-2 px-4 hover:bg-gray-50 dark:hover:bg-gray-700/20" 
                to="#0" 
                @click="dropdownOpen = false"
              >
                <span class="block text-sm mb-2">
                  ❌ <span class="font-medium text-red-600 dark:text-red-400">Stopped {{ service.name }}</span>
                </span>
                <span class="block text-xs font-medium text-gray-400 dark:text-gray-500">
                  {{ getCurrentTime() }}
                </span>
              </router-link>
            </li>
            <li 
              v-if="inactiveServices.length === 0"
              class="border-b border-gray-200 dark:border-gray-700/60 last:border-0"
            >
              <div class="block py-2 px-4">
                <span class="block text-sm text-gray-500 dark:text-gray-400">
                  ✅ All Service is running
                </span>
              </div>
            </li>
          </ul>
          <!--template v-if="smartData">
            <div class="text-xs font-semibold text-gray-400 dark:text-gray-500 uppercase pt-1.5 pb-2 px-4">
              Smart Status
            </div>
            <ul>
              <li 
                v-for="(error, index) in smartData['RunTimeErrors']" 
                :key="index"
                class="border-b border-gray-200 dark:border-gray-700/60 last:border-0"
              >
                <div class="block py-2 px-4">
                  <span class="block text-sm mb-2">
                    ❌ <span class="font-medium text-red-600 dark:text-red-400">{{ error }}</span>
                  </span>
                  <span class="block text-xs font-medium text-gray-400 dark:text-gray-500">
                    {{ smartData['ServiceStartTime'] }}
                  </span>
                </div>
              </li>
              <li 
                v-if="smartData['State'] == 1"
                class="border-b border-gray-200 dark:border-gray-700/60 last:border-0"
              >
                <div class="block py-2 px-4">
                  <span class="block text-sm text-gray-500 dark:text-gray-400">
                    ✅ All Status is OK
                  </span>
                </div>
              </li>
              <li 
                v-else
                class="border-b border-gray-200 dark:border-gray-700/60 last:border-0"
              >
                <div class="block py-2 px-4">
                  <span class="block text-sm">
                    ❌ <span class="font-medium text-red-600 dark:text-red-400">{{ smartData['Message'] }}</span>
                  </span>
                </div>
              </li>
            </ul>
          </template-->
        <!--ul
          ref="dropdown"
          @focusin="dropdownOpen = true"
          @focusout="dropdownOpen = false"
        >
          <li 
            v-for="(isActive, serviceName) in sysData" 
            :key="serviceName"
            v-show="!isActive"
            class="border-b border-gray-200 dark:border-gray-700/60 last:border-0"
          >
            <router-link 
              class="block py-2 px-4 hover:bg-gray-50 dark:hover:bg-gray-700/20" 
              to="#0" 
              @click="dropdownOpen = false"
            >
              <span class="block text-sm mb-2">
                ❌ <span class="font-medium text-red-600 dark:text-red-400">Stopped {{ serviceName }}</span>
              </span>
              <span class="block text-xs font-medium text-gray-400 dark:text-gray-500">
                {{ getCurrentTime() }}
              </span>
            </router-link>
          </li>
      
          <li 
            v-if="inactiveServicesCount === 0"
            class="border-b border-gray-200 dark:border-gray-700/60 last:border-0"
          >
            <div class="block py-2 px-4">
              <span class="block text-sm text-gray-500 dark:text-gray-400 text-center">
                ✅ All Service is running
              </span>
            </div>
          </li>
        </ul-->
      </div> 
    </transition>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'

export default {
  name: 'DropdownNotifications',
  props: {
    align: String,
    status: Boolean,
    data: {
      type: Object,
      default: () => ({}),  // ✅ 함수로 변경,
    }
  },
  setup(props) {
    const dropdownOpen = ref(false)
    const trigger = ref(null)
    const dropdown = ref(null)
    
    // ✅ props.data를 직접 사용 (computed로)
    const sysData = computed(() => props.data)
    
    // ✅ 상태 아이콘 CSS
    const stIconCss = computed(() => {
      return props.status ? 'bg-green-500' : 'bg-red-500'
    })
    
    // ✅ 비활성 서비스 개수 계산
    const inactiveServicesCount = computed(() => {
      if (!props.data) return 0
      return Object.values(props.data).filter(isActive => !isActive).length
    })

    const inactiveServices = computed(() => {
    if (!props.data) return []
    return Object.entries(props.data)
      .filter(([name, isActive]) => !isActive)
      .map(([name, isActive]) => ({ name, isActive }))
  })
    
    // ✅ 현재 시간 포맷팅
    const getCurrentTime = () => {
      const now = new Date()
      return now.toLocaleString('en-US', { 
        month: 'short', 
        day: 'numeric',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      })
    }

    // close on click outside
    const clickHandler = ({ target }) => {
      if (!dropdownOpen.value || dropdown.value.contains(target) || trigger.value.contains(target)) return
      dropdownOpen.value = false
    }

    // close if the esc key is pressed
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
      stIconCss,
      sysData,
      inactiveServicesCount,
      getCurrentTime,
      inactiveServices,
    }
  }
}
</script>