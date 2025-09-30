<template>
  <div class="w-full">
    <!-- 전체 요약 헤더 -->
    <div class="mb-4 p-4 bg-gradient-to-r from-blue-50 to-indigo-50 dark:from-gray-800 dark:to-gray-900 rounded-lg">
      <div class="flex justify-between items-center">
        <div>
          <h3 class="text-lg font-semibold text-gray-800 dark:text-white">IPSM72 #1(72채널)</h3>
          <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">실시간 모니터링</p>
        </div>
        <div class="flex gap-4">
          <div class="text-center">
            <div class="text-2xl font-bold text-green-600">{{ activeCount }}</div>
            <div class="text-xs text-gray-500">정상</div>
          </div>
          <div class="text-center">
            <div class="text-2xl font-bold text-yellow-600">{{ warningCount }}</div>
            <div class="text-xs text-gray-500">경고</div>
          </div>
          <div class="text-center">
            <div class="text-2xl font-bold text-red-600">{{ errorCount }}</div>
            <div class="text-xs text-gray-500">오류</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 그룹별 아코디언 -->
    <div class="space-y-3">
      <div v-for="group in channelGroups" :key="group.id" 
           class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
        
        <!-- 그룹 헤더 -->
        <button
          @click="toggleGroup(group.id)"
          class="w-full px-4 py-3 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
        >
          <div class="flex items-center gap-3">
            <div :class="getGroupStatusClass(group.status)" 
                 class="w-3 h-3 rounded-full"></div>
            <div class="text-left">
              <h4 class="font-medium text-gray-800 dark:text-white">
                {{ group.name }} ({{ group.type }})
              </h4>
              <p class="text-xs text-gray-500 dark:text-gray-400">
                {{ group.channels.length }}개 채널
              </p>
            </div>
          </div>
          
          <div class="flex items-center gap-4">
            <!-- 그룹 요약 정보 -->
            <div class="flex gap-3 text-xs">
              <span class="px-2 py-1 bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400 rounded">
                정상: {{ group.normalCount }}
              </span>
              <span v-if="group.warningCount > 0" 
                    class="px-2 py-1 bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400 rounded">
                경고: {{ group.warningCount }}
              </span>
              <span v-if="group.errorCount > 0"
                    class="px-2 py-1 bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400 rounded">
                오류: {{ group.errorCount }}
              </span>
            </div>
            
            <!-- 화살표 아이콘 -->
            <svg class="w-5 h-5 text-gray-400 transition-transform duration-200" 
                 :class="{ 'rotate-180': openGroups.includes(group.id) }"
                 fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
            </svg>
          </div>
        </button>

        <!-- 채널 상세 (확장 시) -->
        <div v-show="openGroups.includes(group.id)" 
             class="border-t border-gray-200 dark:border-gray-700">
          <div class="p-4 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3">
            <div v-for="channel in group.channels" :key="channel.id"
                 class="bg-gray-50 dark:bg-gray-900 rounded-lg p-3 border border-gray-200 dark:border-gray-700">
              
              <!-- 채널 헤더 -->
              <div class="flex items-center justify-between mb-2">
                <div class="flex items-center gap-2">
                  <div :class="getChannelStatusClass(channel.status)" 
                       class="w-2 h-2 rounded-full"></div>
                  <span class="font-medium text-sm text-gray-800 dark:text-white">
                    {{ channel.name }}
                  </span>
                </div>
                <span class="text-xs px-2 py-0.5 rounded"
                      :class="channel.online ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400' : 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-400'">
                  {{ channel.online ? 'ONLINE' : 'OFFLINE' }}
                </span>
              </div>

              <!-- 채널 데이터 -->
              <div class="space-y-2">
                <!-- 전류 -->
                <div class="flex justify-between items-center">
                  <span class="text-xs text-gray-500">전류</span>
                  <div class="text-sm font-medium">
                    <span v-if="channel.type === '3P4W'" class="space-x-1">
                      <span class="text-blue-600">{{ channel.current.L1 }}A</span>
                      <span class="text-yellow-600">{{ channel.current.L2 }}A</span>
                      <span class="text-red-600">{{ channel.current.L3 }}A</span>
                    </span>
                    <span v-else class="text-gray-800 dark:text-white">
                      {{ channel.current }}A
                    </span>
                  </div>
                </div>

                <!-- 전력 -->
                <div class="flex justify-between items-center">
                  <span class="text-xs text-gray-500">전력</span>
                  <span class="text-sm font-medium text-gray-800 dark:text-white">
                    {{ channel.power }} kW
                  </span>
                </div>

                <!-- 추가 정보 -->
                <div class="grid grid-cols-2 gap-2 pt-2 border-t border-gray-200 dark:border-gray-700">
                  <div>
                    <span class="text-xs text-gray-500">THD</span>
                    <span class="block text-sm font-medium text-gray-800 dark:text-white">
                      {{ channel.thd }}%
                    </span>
                  </div>
                  <div>
                    <span class="text-xs text-gray-500">PF</span>
                    <span class="block text-sm font-medium text-gray-800 dark:text-white">
                      {{ channel.powerFactor }}
                    </span>
                  </div>
                </div>

                <!-- 에너지 -->
                <div class="pt-2 border-t border-gray-200 dark:border-gray-700">
                  <span class="text-xs text-gray-500">총 에너지</span>
                  <span class="block text-sm font-medium text-gray-800 dark:text-white">
                    {{ channel.energy }} kWh
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'

export default {
  name: 'SubDashboard',
  props: {
    data: Object,
    channel: String
  },
  setup(props) {
    const openGroups = ref([])
    
    // 더미 데이터 생성 (실제로는 props.data 사용)
    const channelGroups = ref([
      {
        id: 'group1',
        name: '3P4W',
        type: '3P4W',
        status: 'normal',
        normalCount: 1,
        warningCount: 0,
        errorCount: 0,
        channels: [
          {
            id: 'CT1-3',
            name: 'CT1-3',
            type: '3P4W',
            online: true,
            status: 'normal',
            current: { L1: '0.98', L2: '0.98', L3: '0.98' },
            power: '0.65',
            thd: '0.00',
            powerFactor: '100.0',
            energy: '239.75'
          }
        ]
      },
      {
        id: 'group2',
        name: '3P3W',
        type: '3P3W',
        status: 'warning',
        normalCount: 0,
        warningCount: 1,
        errorCount: 0,
        channels: [
          {
            id: 'CT4-6',
            name: 'CT4-6',
            type: '3P4W',
            online: true,
            status: 'warning',
            current: { L1: '0.98', L2: '0.98', L3: '0.98' },
            power: '0.64',
            thd: '0.15',
            powerFactor: '98.9',
            energy: '152.77'
          }
        ]
      },
      {
        id: 'group3',
        name: '1P2W',
        type: 'L1/L2/L3',
        status: 'normal',
        normalCount: 6,
        warningCount: 0,
        errorCount: 0,
        channels: [
          {
            id: 'CT7',
            name: 'CT7 (L1)',
            type: 'L1',
            online: true,
            status: 'normal',
            current: '0.98',
            power: '0.21',
            thd: '0.26',
            powerFactor: '98.3',
            energy: '53.98'
          },
          {
            id: 'CT8',
            name: 'CT8 (L2)',
            type: 'L2',
            online: true,
            status: 'normal',
            current: '0.98',
            power: '0.22',
            thd: '0.18',
            powerFactor: '100.0',
            energy: '55.52'
          },
          {
            id: 'CT9',
            name: 'CT9 (L3)',
            type: 'L3',
            online: true,
            status: 'normal',
            current: '0.98',
            power: '0.21',
            thd: '0.18',
            powerFactor: '100.0',
            energy: '54.91'
          },
          {
            id: 'CT10',
            name: 'CT10 (L1)',
            type: 'L1',
            online: true,
            status: 'normal',
            current: '0.98',
            power: '0.21',
            thd: '0.00',
            powerFactor: '98.3',
            energy: '49.87'
          },
          {
            id: 'CT11',
            name: 'CT11 (L2)',
            type: 'L2',
            online: true,
            status: 'normal',
            current: '0.98',
            power: '0.21',
            thd: '0.00',
            powerFactor: '98.3',
            energy: '52.13'
          },
          {
            id: 'CT12',
            name: 'CT12 (L3)',
            type: 'L3',
            online: true,
            status: 'normal',
            current: '0.98',
            power: '0.21',
            thd: '0.00',
            powerFactor: '98.3',
            energy: '52.42'
          }
        ]
      }
    ])

    // 전체 카운트 계산
    const activeCount = computed(() => {
      return channelGroups.value.reduce((sum, group) => sum + group.normalCount, 0)
    })

    const warningCount = computed(() => {
      return channelGroups.value.reduce((sum, group) => sum + group.warningCount, 0)
    })

    const errorCount = computed(() => {
      return channelGroups.value.reduce((sum, group) => sum + group.errorCount, 0)
    })

    // 그룹 토글
    const toggleGroup = (groupId) => {
      const index = openGroups.value.indexOf(groupId)
      if (index > -1) {
        openGroups.value.splice(index, 1)
      } else {
        openGroups.value.push(groupId)
      }
    }

    // 상태 클래스
    const getGroupStatusClass = (status) => {
      switch(status) {
        case 'normal': return 'bg-green-500'
        case 'warning': return 'bg-yellow-500'
        case 'error': return 'bg-red-500'
        default: return 'bg-gray-400'
      }
    }

    const getChannelStatusClass = (status) => {
      switch(status) {
        case 'normal': return 'bg-green-500'
        case 'warning': return 'bg-yellow-500'
        case 'error': return 'bg-red-500'
        default: return 'bg-gray-400'
      }
    }

    onMounted(() => {
      // 첫 번째 그룹 자동 열기
      openGroups.value.push('group1')
    })

    return {
      openGroups,
      channelGroups,
      activeCount,
      warningCount,
      errorCount,
      toggleGroup,
      getGroupStatusClass,
      getChannelStatusClass
    }
  }
}
</script>

<style scoped>
/* 애니메이션 효과 */
.rotate-180 {
  transform: rotate(180deg);
}

/* 스크롤바 스타일 */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  @apply bg-gray-100 dark:bg-gray-800 rounded;
}

::-webkit-scrollbar-thumb {
  @apply bg-gray-400 dark:bg-gray-600 rounded;
}

::-webkit-scrollbar-thumb:hover {
  @apply bg-gray-500 dark:bg-gray-500;
}
</style>