<template>
  <div class="flex h-[100dvh] overflow-hidden">
    <Sidebar :sidebarOpen="sidebarOpen" @close-sidebar="sidebarOpen = false" />

    <div class="relative flex flex-col flex-1 overflow-y-auto overflow-x-hidden">
      <Header :sidebarOpen="sidebarOpen" @toggle-sidebar="sidebarOpen = !sidebarOpen" />

      <main class="grow">
        <div class="px-2 sm:px-4 lg:px-6 py-4 w-full max-w-full">
          <div class="mb-4">
            <h2 class="text-xl md:text-2xl text-gray-800 dark:text-gray-100 font-bold">
              {{ t('sidebar.module') }} > IBSM
            </h2>
          </div>

          <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">

            <!-- ========== 좌측: TapBox 상태 (col-5) ========== -->
            <div class="lg:col-span-5 meter-card p-5 flex flex-col">
              <div class="meter-card-header !px-0 !py-0 mb-4">
                <h3 class="meter-card-title meter-accent-blue">TapBox Status</h3>
                <div class="flex flex-wrap items-center gap-2 text-xs">
                  <span v-for="s in statusDefs" :key="s.key" class="flex items-center gap-1">
                    <span class="w-2 h-2 rounded-full" :style="{ backgroundColor: s.color }"></span>
                    <span class="text-gray-500 dark:text-gray-400">{{ s.label }}</span>
                  </span>
                </div>
              </div>

              <div class="grid grid-cols-6 gap-1 overflow-y-auto flex-1 content-start p-2" style="max-height: calc(100vh - 220px);">
                <button
                  v-for="idx in 32"
                  :key="idx"
                  @click="isTapboxConfigured(idx) && onTapboxClick(idx)"
                  class="flex flex-col items-center justify-center rounded-lg border-2 py-2 px-1 transition-all"
                  :class="[
                    isTapboxConfigured(idx)
                      ? [getTapboxClass(idx), 'hover:scale-105 hover:shadow-md cursor-pointer']
                      : 'bg-gray-50 dark:bg-gray-800 text-gray-300 dark:text-gray-600 border-gray-200 dark:border-gray-700 cursor-default opacity-50',
                    selectedTapbox === idx ? 'ring-2 ring-offset-1 ring-blue-500 dark:ring-offset-gray-800' : ''
                  ]"
                  :style="isTapboxConfigured(idx) ? { borderColor: getTapboxColor(idx) } : {}"
                >
                  <span class="text-base font-bold leading-none">{{ idx }}</span>
                  <span class="text-[10px] font-semibold tracking-tight mt-1">
                    {{ isTapboxConfigured(idx) ? getTapboxStatus(idx) : 'NoDevice' }}
                  </span>
                </button>
              </div>
            </div>

            <!-- ========== 우측 (col-7) ========== -->
            <div v-if="selectedTapbox && tapBoxData" class="lg:col-span-7">
              <div class="meter-card p-5">
                <!-- 공통 헤더 -->
                <div class="meter-card-header !px-0 !py-0 mb-3">
                  <h3 class="meter-card-title meter-accent-emerald">TapBox #{{ selectedTapbox }}{{ tapBoxData.can_id ? ' - ' + tapBoxData.can_id : '' }}</h3>
                  <span class="text-xs text-gray-400 dark:text-gray-500 font-mono">
                    {{ tapBoxData.timestamp ? formatTimestamp(tapBoxData.timestamp) : '--' }}
                  </span>
                </div>

                <!-- Status / Temp / Freq -->
                <div class="grid grid-cols-3 gap-2 mb-4">
                  <div class="flex flex-col items-center py-2 px-1 rounded-lg border"
                       :class="getTapboxClass(selectedTapbox)"
                       :style="{ borderColor: getTapboxColor(selectedTapbox) }">
                    <span class="stat-label">Status</span>
                    <span class="text-sm font-mono font-bold tabular-nums">{{ getTapboxStatus(selectedTapbox) }}</span>
                  </div>
                  <div class="stat-box">
                    <span class="stat-label">Temperature</span>
                    <span class="stat-value">{{ tapBoxData.temp?.toFixed(1) ?? '-' }} °C</span>
                  </div>
                  <div class="stat-box">
                    <span class="stat-label">Frequency</span>
                    <span class="stat-value">{{ tapBoxData.freq?.toFixed(2) ?? '-' }} Hz</span>
                  </div>
                </div>

                <!-- CB 탭 -->
                <div v-if="tapBoxData.cb?.length" class="flex border-b border-gray-200 dark:border-gray-700 mb-4">
                  <button v-for="cb in tapBoxData.cb" :key="cb.index"
                    class="px-4 py-2 text-sm font-medium border-b-2 transition-colors"
                    :class="activeCbIndex === cb.index
                      ? 'border-violet-500 text-violet-600 dark:text-violet-400'
                      : 'border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300'"
                    @click="activeCbIndex = cb.index">
                    CB{{ cb.index }} · {{ getCbtypeLabel(cb.cb_type) }}
                  </button>
                </div>

                <!-- 선택된 CB 내용 -->
                <div v-if="activeCb">
                  <!-- 측정 테이블 -->
                  <div class="overflow-x-auto">
                    <table class="w-full text-sm min-w-[400px]">
                      <thead>
                        <tr class="border-b-2 border-gray-200 dark:border-gray-600">
                          <th class="text-left py-2 px-3 font-semibold text-gray-500 dark:text-gray-400 w-28">Items</th>
                          <th class="text-right py-2 px-3 font-semibold text-gray-500 dark:text-gray-400">Total</th>
                          <th v-if="showL1(activeCb.cb_type)" class="text-right py-2 px-3 font-semibold text-gray-500 dark:text-gray-400">L1</th>
                          <th v-if="showL2(activeCb.cb_type)" class="text-right py-2 px-3 font-semibold text-gray-500 dark:text-gray-400">L2</th>
                          <th v-if="showL3(activeCb.cb_type)" class="text-right py-2 px-3 font-semibold text-gray-500 dark:text-gray-400">L3</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr v-for="(row, i) in tableRows" :key="row.item"
                            :class="i % 2 === 0 ? 'bg-gray-50 dark:bg-gray-700/30' : 'bg-white dark:bg-gray-800'">
                          <td class="py-2.5 px-3 font-semibold text-gray-700 dark:text-gray-200">{{ row.item }}</td>
                          <td class="py-2.5 px-3 text-right font-mono text-gray-900 dark:text-gray-100 tabular-nums">{{ row.total ?? '-' }}</td>
                          <td v-if="showL1(activeCb.cb_type)" class="py-2.5 px-3 text-right font-mono text-gray-900 dark:text-gray-100 tabular-nums">{{ row.l1 ?? '-' }}</td>
                          <td v-if="showL2(activeCb.cb_type)" class="py-2.5 px-3 text-right font-mono text-gray-900 dark:text-gray-100 tabular-nums">{{ row.l2 ?? '-' }}</td>
                          <td v-if="showL3(activeCb.cb_type)" class="py-2.5 px-3 text-right font-mono text-gray-900 dark:text-gray-100 tabular-nums">{{ row.l3 ?? '-' }}</td>
                        </tr>
                      </tbody>
                    </table>
                  </div>

                  <!-- 요약 스탯 (테이블 아래) -->
                  <div class="grid grid-cols-5 gap-2 mt-4">
                    <div class="stat-box">
                      <span class="stat-label">kWh</span>
                      <span class="stat-value">{{ fmt(activeCb.kwh_import) }}</span>
                    </div>
                    <div class="stat-box">
                      <span class="stat-label">V ub(%)</span>
                      <span class="stat-value">{{ fmt(activeCb.volt_unbal) }}</span>
                    </div>
                    <div class="stat-box">
                      <span class="stat-label">I ub(%)</span>
                      <span class="stat-value">{{ fmt(activeCb.curr_unbal) }}</span>
                    </div>
                    <div class="stat-box">
                      <span class="stat-label">THD(%)</span>
                      <span class="stat-value">{{ fmt(activeCb.pthd) }}</span>
                    </div>
                    <div class="stat-box">
                      <span class="stat-label">Ig(A)</span>
                      <span class="stat-value">{{ fmt(activeCb.ig) }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

          </div>
        </div>
      </main>
      <Footer />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from 'axios'
import { useSetupStore } from '@/store/setup'
import { useAuthStore } from '@/store/auth'
import Sidebar from '../common/SideBar.vue'
import Header from '../common/Header.vue'
import Footer from '../common/Footer.vue'

const { t } = useI18n()
const sidebarOpen = ref(false)
const setupStore = useSetupStore()
const authStore = useAuthStore()

const tapBoxData = ref(null)
const tapBoxAllData = ref({})  // { tapboxIdx: tapboxData }
const selectedTapbox = ref(null)
const activeCbIndex = ref(1)
let pollInterval = null

// ===== TapBox 설정 =====
function isTapboxConfigured(idx) {
  const tapboxs = setupStore.ibsmTapboxs || []
  return idx <= tapboxs.length
}

// ===== device_type 비트 파싱 =====
const mtypeLabelMap = { 0: 'Not Used', 1: 'iBSM10', 2: 'iBSM20', 3: 'iBSM30' }
const phaseLabelMap = { 0: '3P', 1: '1P', 2: '1P+Z' }
const cbtypeLabelMap = { 0: 'N/A', 1: '1PL1', 2: '1PL2', 3: '1PL3', 4: '3P3W', 5: '3P4W', 6: '1PL1Z', 7: '1PL2Z', 8: '1PL3Z' }

const parsedDeviceType = computed(() => {
  const dt = tapBoxData.value?.device_type || 0
  const mtype  = (dt >> 8) & 0x03
  const stype  = (dt >> 6) & 0x03
  const phase  = (dt >> 4) & 0x03
  const cbtype = dt & 0x0F
  return {
    mtype, stype, phase, cbtype,
    mtypeLabel: mtypeLabelMap[mtype] || 'Unknown',
    phaseLabel: phaseLabelMap[phase] || 'Unknown',
    cbtypeLabel: cbtypeLabelMap[cbtype] || 'Unknown',
  }
})

// ===== Timestamp =====
function formatTimestamp(ts) {
  if (!ts) return ''
  const d = new Date(ts * 1000)
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

// ===== Status =====
const statusDefs = [
  { key: 'OFF',   label: 'OFF',   color: '#9ca3af' },
  { key: 'ON',    label: 'ON',    color: '#22c55e' },
  { key: 'TRIP',  label: 'TRIP',  color: '#ef4444' },
  { key: 'OPR',   label: 'OPR',   color: '#3b82f6' },
  { key: 'OC',    label: 'OC',    color: '#f97316' },
  { key: 'OCG',   label: 'OCG',   color: '#a855f7' },
  { key: 'ELD',   label: 'ELD',   color: '#eab308' },
  { key: 'SAG',   label: 'SAG',   color: '#06b6d4' },
  { key: 'SWELL', label: 'SWELL', color: '#ec4899' },
]
const statusColorMap = Object.fromEntries(statusDefs.map(s => [s.key, s.color]))

function getTapboxStatus(idx) {
  if (!isTapboxConfigured(idx)) return 'OFF'
  const data = tapBoxAllData.value[idx]
  if (!data) return 'OFF'
  if (data.status) {
    const st = data.status
    if (st & 0x03F00000) return 'TRIP'   // bit 20-25: CB Trip
    if (st & 0xFC000000) return 'OCG'    // bit 26-31: OCR_TRIP
    if (st & 0x000F0000) return 'ELD'    // bit 16-19: ELD (누설전류)
    if (st & 0x0000FC00) return 'OC'     // bit 10-15: OCR CB (과전류)
    if (st & 0x00000002) return 'SAG'    // bit 1: SAG
    if (st & 0x00000004) return 'SWELL'  // bit 2: SWELL
    if (st & 0x000003F0) return 'OPR'    // bit 4-9: LOAD CB (작동)
  }
  return 'ON'
}
function getTapboxColor(idx) { return statusColorMap[getTapboxStatus(idx)] || '#9ca3af' }
function getTapboxClass(idx) {
  switch (getTapboxStatus(idx)) {
    case 'ON':    return 'bg-green-50 dark:bg-green-900/20 text-green-700 dark:text-green-300'
    case 'TRIP':  return 'bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-300 animate-pulse'
    case 'OPR':   return 'bg-blue-50 dark:bg-blue-900/20 text-blue-700 dark:text-blue-300'
    case 'OC':    return 'bg-orange-50 dark:bg-orange-900/20 text-orange-700 dark:text-orange-300'
    case 'OCG':   return 'bg-purple-50 dark:bg-purple-900/20 text-purple-700 dark:text-purple-300'
    case 'ELD':   return 'bg-yellow-50 dark:bg-yellow-900/20 text-yellow-700 dark:text-yellow-300'
    case 'SAG':   return 'bg-cyan-50 dark:bg-cyan-900/20 text-cyan-700 dark:text-cyan-300'
    case 'SWELL': return 'bg-pink-50 dark:bg-pink-900/20 text-pink-700 dark:text-pink-300'
    default:      return 'bg-white dark:bg-gray-800 text-gray-400 dark:text-gray-500'
  }
}

// ===== CB Type 컬럼 =====
function showL1(cbType) { return [1, 4, 5, 6].includes(cbType) }
function showL2(cbType) { return [2, 4, 5, 7].includes(cbType) }
function showL3(cbType) { return [3, 4, 5, 8].includes(cbType) }
function getCbtypeLabel(cbType) { return cbtypeLabelMap[cbType] || 'Unknown' }
const fmt = (v) => v != null ? v.toFixed(2) : '-'

const activeCb = computed(() => tapBoxData.value?.cb?.find(c => c.index === activeCbIndex.value) || null)

const tableRows = computed(() => {
  const cb = activeCb.value
  if (!cb) return []
  const ct = cb.cb_type
  const l1 = showL1(ct), l2 = showL2(ct), l3 = showL3(ct)
  return [
    { item: 'V',   total: null,            l1: l1 ? fmt(cb.v[0]) : null,     l2: l2 ? fmt(cb.v[1]) : null,     l3: l3 ? fmt(cb.v[2]) : null },
    { item: 'Vll', total: null,            l1: l1 ? fmt(cb.vline[0]) : null, l2: l2 ? fmt(cb.vline[1]) : null, l3: l3 ? fmt(cb.vline[2]) : null },
    { item: 'I',   total: null,            l1: l1 ? fmt(cb.i[0]) : null,     l2: l2 ? fmt(cb.i[1]) : null,     l3: l3 ? fmt(cb.i[2]) : null },
    { item: 'W',   total: fmt(cb.p_total), l1: l1 ? fmt(cb.p[0]) : null,     l2: l2 ? fmt(cb.p[1]) : null,     l3: l3 ? fmt(cb.p[2]) : null },
    { item: 'var', total: fmt(cb.q_total), l1: null, l2: null, l3: null },
    { item: 'VA',  total: fmt(cb.s_total), l1: null, l2: null, l3: null },
    { item: 'PF',  total: fmt(cb.pf_total),l1: null, l2: null, l3: null },
  ]
})

// ===== Fetch =====
async function fetchTapBoxAll() {
  try {
    const response = await axios.get('/api/getTapBoxAll')
    if (response.data.success) {
      const map = {}
      response.data.data.forEach(item => {
        map[item.id] = item.data
      })
      tapBoxAllData.value = map
    }
  } catch (error) {
    console.error('[IBSM] fetchAll error:', error)
  }
}

async function fetchTapBox(id) {
  try {
    const response = await axios.get(`/api/getTapBox/${id}`)
    if (response.data.success) {
      tapBoxData.value = response.data.data
      // 현재 선택된 CB가 없거나 데이터에 없으면 첫번째 CB 선택
      if (tapBoxData.value.cb?.length) {
        const exists = tapBoxData.value.cb.some(c => c.index === activeCbIndex.value)
        if (!exists) {
          activeCbIndex.value = tapBoxData.value.cb[0].index
        }
      }
    }
  } catch (error) {
    console.error('[IBSM] fetch error:', error)
  }
}

function onTapboxClick(idx) {
  selectedTapbox.value = selectedTapbox.value === idx ? null : idx
  tapBoxData.value = null
  if (selectedTapbox.value) {
    fetchTapBox(selectedTapbox.value)
    if (pollInterval) clearInterval(pollInterval)
    pollInterval = setInterval(() => fetchTapBox(selectedTapbox.value), 5000)
  } else {
    if (pollInterval) clearInterval(pollInterval)
  }
}

onMounted(async () => {
  if (!setupStore.ibsmTapboxs?.length && authStore.getOpMode === 'device1') {
    await setupStore.fetchIbsm()
  }
  await fetchTapBoxAll()
  if (setupStore.ibsmTapboxs?.length) {
    selectedTapbox.value = 1
    await fetchTapBox(1)
  }
  pollInterval = setInterval(async () => {
    await fetchTapBoxAll()
    if (selectedTapbox.value) {
      await fetchTapBox(selectedTapbox.value)
    }
  }, 5000)
})

onUnmounted(() => {
  if (pollInterval) clearInterval(pollInterval)
})
</script>

<style scoped>
@import '../../css/meter-card.css';

.info-badge {
  @apply px-2 py-1 rounded bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300;
}

.stat-box {
  @apply flex flex-col items-center py-2 px-1 rounded-lg;
  @apply bg-gray-50 dark:bg-gray-700/40 border border-gray-100 dark:border-gray-700/60;
}
.stat-label {
  @apply text-[10px] font-bold text-gray-400 dark:text-gray-500 uppercase mb-1;
}
.stat-value {
  @apply text-sm font-mono font-bold text-gray-800 dark:text-gray-100 tabular-nums;
}
</style>
