<template>
  <div class="premium-dashboard-card">
    <!-- 헤더 -->
    <div class="card-header">
      <header class="header-content">
        <h2 class="card-title">IO 모듈 상태</h2>        
      </header>
    </div>

    <!-- 컨텐츠 -->
    <div class="card-content">
      <div class="main-layout">
        <!-- 첫 번째 줄: DI -->
        <div class="row-section">
          <div class="section di-section">
            <div class="section-header">
              <span class="section-title">DI</span>
              <div class="di-controls">
                <span class="label">Di:</span>
                <button class="btn-toggle active">ON</button>
                <button class="btn-toggle btn-toggle-off">OFF</button>
                <span class="label">Pi:</span>
                <span class="pi-count">{{ piCount }}</span>
              </div>
            </div>
            <div class="di-grid">
              <div v-for="di in diData" :key="di.id" class="di-item" :class="getDIClass(di)">
                <span class="di-num">{{ di.id }}</span>
                <span class="di-val">{{ di.type === 'PI' ? di.count : di.status }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 두 번째 줄: DO + DC -->
        <div class="row-section dual-row">
          <!-- DO -->
          <div class="section compact-section do-section">
            <div class="section-header-sm">
              <span class="section-title-sm">DO</span>
              <div class="do-labels">
                <span class="do-type-label alarm">Alarm</span>
                <span class="do-type-label event">Event</span>
                <span class="do-type-label output">Output</span>
              </div>
            </div>
            <div class="do-grid">
              <div v-for="item in doData" :key="item.id" class="do-item" :class="getDOClass(item)">
                <span class="do-num">{{ item.id }}</span>
                <span class="do-val">{{ item.status }}</span>
              </div>
            </div>
          </div>

          <!-- DC -->
          <div class="section compact-section dc-section">
            <div class="section-header-sm">
              <span class="section-title-sm">DC</span>
            </div>
            <div class="dc-compact">
              <div class="dc-row">
                <span class="dc-label-sm">DC V</span>
                <span class="dc-value-sm">{{ dcData.voltage.toFixed(3) }} V</span>
              </div>
              <div class="dc-row">
                <span class="dc-label-sm">DC I</span>
                <span class="dc-value-sm">{{ dcData.current.toFixed(3) }} A</span>
              </div>
              <div class="dc-row">
                <span class="dc-label-sm">DC Ibat</span>
                <span class="dc-value-sm">{{ dcData.battery.toFixed(3) }} A</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 세 번째 줄: AI + Temperature -->
        <div class="row-section dual-row">
          <div class="section compact-section">
            <div class="section-header-sm">
              <span class="section-title-sm">AI</span>
            </div>
            <div class="values-grid">
              <div v-for="item in aiData" :key="item.id" class="value-cell">
                <span class="cell-num">{{ item.id }}</span>
                <span class="cell-value">{{ item.value.toFixed(3) }}</span>
              </div>
            </div>
          </div>

          <div class="section compact-section">
            <div class="section-header-sm">
              <span class="section-title-sm">Temperature</span>
            </div>
            <div class="values-grid">
              <div v-for="item in tempData" :key="item.id" class="value-cell">
                <span class="cell-num">{{ item.id }}</span>
                <span class="cell-value">{{ item.value.toFixed(3) }} °C</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue'

export default {
  name: 'ModuleStatusCard',
  props: {
    channel: String,
    data: Object,
  },
  setup() {
    const piCount = ref(3)

    // DI 24개 고정 - PI 샘플 3개 추가
    const diData = ref([
      { id: 1, type: 'DI', status: 'OFF' },
      { id: 2, type: 'DI', status: 'ON' },
      { id: 3, type: 'DI', status: 'ON' },
      { id: 4, type: 'DI', status: 'OFF' },
      { id: 5, type: 'DI', status: 'ON' },
      { id: 6, type: 'DI', status: 'ON' },
      { id: 7, type: 'DI', status: 'OFF' },
      { id: 8, type: 'DI', status: 'ON' },
      { id: 9, type: 'DI', status: 'ON' },
      { id: 10, type: 'DI', status: 'ON' },
      { id: 11, type: 'DI', status: 'OFF' },
      { id: 12, type: 'DI', status: 'OFF' },
      { id: 13, type: 'DI', status: 'OFF' },
      { id: 14, type: 'DI', status: 'OFF' },
      { id: 15, type: 'DI', status: 'ON' },
      { id: 16, type: 'DI', status: 'ON' },
      { id: 17, type: 'PI', status: 'PI', count: 21 },
      { id: 18, type: 'DI', status: 'OFF' },
      { id: 19, type: 'DI', status: 'ON' },
      { id: 20, type: 'DI', status: 'ON' },
      { id: 21, type: 'DI', status: 'ON' },
      { id: 22, type: 'PI', status: 'PI', count: 15 },
      { id: 23, type: 'PI', status: 'PI', count: 8 },
      { id: 24, type: 'DI', status: 'ON' },
    ])

    // DI 클래스 결정
    const getDIClass = (di) => {
      if (di.type === 'PI') return 'di-pi'
      return di.status === 'ON' ? 'di-on' : 'di-off'
    }

    // DO 데이터 - 타입별 색상 구분
    const doData = ref([
      { id: 1, type: 'Alarm', status: 'OFF' },
      { id: 2, type: 'Event', status: 'OFF' },
      { id: 3, type: 'Output', status: 'ON' },
    ])

    // DO 클래스 결정 (타입별 색상)
    const getDOClass = (doItem) => {
      const typeClass = {
        'Alarm': 'do-alarm',
        'Event': 'do-event',
        'Output': 'do-output'
      }[doItem.type] || 'do-alarm'
      
      return typeClass
    }

    const aiData = ref([
      { id: 1, value: 8.685 },
      { id: 2, value: 8.199 },
      { id: 3, value: 0.000 },
      { id: 4, value: 0.000 },
    ])

    const tempData = ref([
      { id: 1, value: 28.262 },
      { id: 2, value: 25.258 },
      { id: 3, value: 0.000 },
      { id: 4, value: 0.000 },
    ])

    const dcData = ref({
      voltage: 0.000,
      current: 0.000,
      battery: 0.000,
    })

    return {
      piCount,
      diData,
      getDIClass,
      doData,
      getDOClass,
      aiData,
      tempData,
      dcData,
    }
  }
}
</script>

<style scoped>
.premium-dashboard-card {
  @apply flex flex-col col-span-full sm:col-span-6 xl:col-span-6;
  @apply bg-gradient-to-br from-white to-gray-50 dark:from-gray-800 dark:to-gray-900;
  @apply shadow-lg rounded-xl border border-gray-200/50 dark:border-gray-700/50;
  @apply backdrop-blur-sm;
  @apply transition-all duration-300 hover:shadow-xl;
}

/* 헤더 섹션 */
.card-header {
  @apply p-3 border-b border-gray-200/50 dark:border-gray-700/50;
  @apply bg-gradient-to-r from-blue-100 to-purple-100 dark:from-blue-900/20 dark:to-purple-900/20;
  @apply rounded-t-xl;
}

.header-content {
  @apply flex justify-between items-center;
}

.card-title {
  @apply text-lg font-bold text-gray-900 dark:text-gray-100;
  @apply bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent;
}

.card-content {
  @apply p-2;
}

/* 메인 3줄 레이아웃 */
.main-layout {
  @apply flex flex-col gap-2;
}

.row-section {
  @apply w-full;
}

.dual-row {
  @apply grid grid-cols-2 gap-2;
}

/* 섹션 공통 */
.section {
  @apply bg-gray-50 dark:bg-gray-800/50 rounded-lg p-2;
  @apply border border-gray-200 dark:border-gray-700;
}

.compact-section {
  @apply bg-gray-50 dark:bg-gray-800/50 rounded p-1.5;
  @apply border border-gray-200 dark:border-gray-700;
}

.section-header {
  @apply flex justify-between items-center mb-1.5;
}

.section-header-sm {
  @apply flex justify-between items-center mb-1;
}

.section-title {
  @apply text-sm font-bold text-gray-800 dark:text-white;
}

.section-title-sm {
  @apply text-xs font-bold text-gray-800 dark:text-white;
}

/* DI 섹션 */
.di-section {
  @apply flex flex-col;
}

.di-controls {
  @apply flex items-center gap-1;
}

.label {
  @apply text-xs text-gray-600 dark:text-gray-400;
}

.btn-toggle {
  @apply px-2 py-0.5 rounded text-xs font-medium;
  @apply bg-gray-200 dark:bg-gray-700 text-gray-600 dark:text-gray-300;
  @apply transition-colors;
}

.btn-toggle.active {
  @apply bg-green-500 text-white;
}

.btn-toggle-off {
  @apply bg-gray-400 text-white;
}

.pi-count {
  @apply text-xs font-bold px-1.5 py-0.5 rounded;
  @apply bg-blue-400 text-white;
}

.di-grid {
  @apply grid grid-cols-12 gap-1;
}

.di-item {
  @apply flex flex-col items-center justify-center;
  @apply rounded p-1;
}

.di-item.di-on {
  @apply bg-green-500 dark:bg-green-600;
}

.di-item.di-off {
  @apply bg-gray-400 dark:bg-gray-500;
}

.di-item.di-pi {
  @apply bg-blue-400 dark:bg-blue-500;
}

.di-num {
  @apply text-xs text-white font-medium;
}

.di-val {
  @apply text-xs text-white font-bold;
}

/* DO 섹션 */
.do-section {
  @apply flex flex-col;
}

.do-labels {
  @apply flex items-center gap-1;
}

.do-type-label {
  @apply text-xs font-medium px-1.5 py-0.5 rounded;
}

.do-type-label.alarm {
  @apply bg-blue-400 text-white;
}

.do-type-label.event {
  @apply bg-green-500 text-white;
}

.do-type-label.output {
  @apply bg-gray-400 text-white;
}

.do-grid {
  @apply grid grid-cols-3 gap-1;
}

.do-item {
  @apply flex flex-col items-center justify-center;
  @apply rounded p-1;
}

.do-item.do-alarm {
  @apply bg-blue-400 dark:bg-blue-500;
}

.do-item.do-event {
  @apply bg-green-500 dark:bg-green-600;
}

.do-item.do-output {
  @apply bg-gray-400 dark:bg-gray-500;
}

.do-num {
  @apply text-xs text-white font-medium;
}

.do-val {
  @apply text-xs text-white font-bold;
}

/* DC 섹션 */
.dc-section {
  @apply flex flex-col;
}

.dc-compact {
  @apply flex gap-1;
}

.dc-row {
  @apply flex flex-col items-center flex-1;
  @apply bg-white dark:bg-gray-700 rounded px-1.5 py-1;
  @apply border border-gray-200 dark:border-gray-600;
}

.dc-label-sm {
  @apply text-xs text-gray-600 dark:text-gray-400 mb-0.5;
}

.dc-value-sm {
  @apply text-xs font-bold text-gray-800 dark:text-white;
}

/* AI & Temperature - 2x2 그리드 (좌우 배치) */
.values-grid {
  @apply grid grid-cols-2 gap-1;
}

.value-cell {
  @apply flex items-center justify-between;
  @apply bg-white dark:bg-gray-700 rounded px-2 py-1;
  @apply border border-gray-200 dark:border-gray-600;
}

.cell-num {
  @apply text-xs text-gray-600 dark:text-gray-400;
}

.cell-value {
  @apply text-xs font-bold text-gray-800 dark:text-white;
}

/* 반응형 */
@media (max-width: 1024px) {
  .di-grid {
    @apply grid-cols-8;
  }
  .dual-row {
    @apply grid-cols-1;
  }
}

@media (max-width: 768px) {
  .di-grid {
    @apply grid-cols-6;
  }
}
</style>