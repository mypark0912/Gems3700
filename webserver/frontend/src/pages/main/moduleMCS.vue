<template>
  <div class="flex h-[100dvh] overflow-hidden">
    <Sidebar :sidebarOpen="sidebarOpen" @close-sidebar="sidebarOpen = false" />

    <div class="relative flex flex-col flex-1 overflow-y-auto overflow-x-hidden">
      <Header :sidebarOpen="sidebarOpen" @toggle-sidebar="sidebarOpen = !sidebarOpen" />

      <main class="grow">
        <div class="px-2 sm:px-4 lg:px-6 py-4 w-full max-w-full">
          <!-- 타이틀 + 탭 -->
          <div class="mb-4 flex items-center gap-3 flex-wrap">
            <h2 class="text-xl md:text-2xl text-gray-800 dark:text-gray-100 font-bold">
              {{ t('sidebar.module') }} > MCS
            </h2>
            <span class="text-xl md:text-2xl text-gray-400 dark:text-gray-500 font-bold">&gt;</span>
            <div class="flex items-center gap-1">
              <button
                class="px-3 py-1 text-lg md:text-xl font-bold rounded-md transition-colors"
                :class="activeTab === 'sv'
                  ? 'text-violet-600 dark:text-violet-400 bg-violet-50 dark:bg-violet-900/20'
                  : 'text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300'"
                @click="activeTab = 'sv'"
              >System / Voltage</button>
              <span class="text-gray-300 dark:text-gray-600">·</span>
              <button
                class="px-3 py-1 text-lg md:text-xl font-bold rounded-md transition-colors"
                :class="activeTab === 'feeder'
                  ? 'text-violet-600 dark:text-violet-400 bg-violet-50 dark:bg-violet-900/20'
                  : 'text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300'"
                @click="activeTab = 'feeder'"
              >Feeder</button>
            </div>
          </div>

          <!-- ══════════════════════════════════════════
               System / Voltage 탭 (상전압 · 선간전압 분리)
          ══════════════════════════════════════════ -->
          <div v-if="activeTab === 'sv'" class="meter-card p-5">
            <div class="meter-card-header !px-0 !py-0 mb-4">
              <h3 class="meter-card-title meter-accent-blue">System / Voltage</h3>
              <span class="text-xs font-mono text-gray-500 dark:text-gray-400">{{ formattedTimestamp }}</span>
            </div>

            <!-- Frequency 인라인 -->
            <div class="flex flex-wrap items-center gap-x-8 gap-y-2 pb-4 mb-5 border-b border-gray-200 dark:border-gray-700/60">
              <div class="flex items-baseline gap-2">
                <span class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider">Frequency</span>
                <span class="text-2xl font-mono font-bold text-gray-800 dark:text-gray-100 tabular-nums">
                  {{ voltageData?.freq != null ? voltageData.freq.toFixed(2) : '-' }}
                  <span class="text-sm font-semibold text-gray-400 dark:text-gray-500">Hz</span>
                </span>
              </div>
            </div>

            <!-- Voltage (상전압 · 선간전압) -->
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <!-- 상전압 -->
              <div class="voltage-section">
                <div class="flex items-center gap-2 mb-3">
                  <div class="w-1 h-4 bg-indigo-500 rounded-full"></div>
                  <span class="text-sm font-semibold text-gray-700 dark:text-gray-200 uppercase tracking-wider">상전압 (Phase)</span>
                </div>

                <div class="grid grid-cols-3 gap-3">
                  <div v-for="(lbl, i) in ['L1','L2','L3']" :key="'p'+i" class="voltage-cell">
                    <span class="voltage-label">{{ lbl }}</span>
                    <div class="flex items-baseline gap-1.5">
                      <span class="voltage-value">{{ fmt(voltageData?.u?.[i]) }}</span>
                      <span class="voltage-unit">V</span>
                    </div>
                    <div class="voltage-thd">
                      <span class="thd-label">THD</span>
                      <span class="thd-value">{{ fmt(voltageData?.thd_u?.[i]) }}<span class="thd-unit">%</span></span>
                    </div>
                  </div>
                </div>

                <div class="voltage-summary">
                  <div class="summary-item">
                    <span class="summary-label">평균</span>
                    <span class="summary-value">{{ fmt(voltageData?.avg_u) }}<span class="summary-unit">V</span></span>
                  </div>
                  <div class="summary-divider"></div>
                  <div class="summary-item">
                    <span class="summary-label">평균 THD</span>
                    <span class="summary-value">{{ fmt(voltageData?.avg_thd_u) }}<span class="summary-unit">%</span></span>
                  </div>
                </div>
              </div>

              <!-- 선간전압 -->
              <div class="voltage-section">
                <div class="flex items-center gap-2 mb-3">
                  <div class="w-1 h-4 bg-emerald-500 rounded-full"></div>
                  <span class="text-sm font-semibold text-gray-700 dark:text-gray-200 uppercase tracking-wider">선간전압 (Line)</span>
                </div>

                <div class="grid grid-cols-3 gap-3">
                  <div v-for="(lbl, i) in ['L1-L2','L2-L3','L3-L1']" :key="'l'+i" class="voltage-cell">
                    <span class="voltage-label">{{ lbl }}</span>
                    <div class="flex items-baseline gap-1.5">
                      <span class="voltage-value">{{ fmt(voltageData?.ull?.[i]) }}</span>
                      <span class="voltage-unit">V</span>
                    </div>
                    <div class="voltage-thd">
                      <span class="thd-label">THD</span>
                      <span class="thd-value">{{ fmt(voltageData?.thd_ull?.[i]) }}<span class="thd-unit">%</span></span>
                    </div>
                  </div>
                </div>

                <div class="voltage-summary">
                  <div class="summary-item">
                    <span class="summary-label">평균</span>
                    <span class="summary-value">{{ fmt(voltageData?.avg_ull) }}<span class="summary-unit">V</span></span>
                  </div>
                  <div class="summary-divider"></div>
                  <div class="summary-item">
                    <span class="summary-label">평균 THD</span>
                    <span class="summary-value">{{ fmt(voltageData?.avg_thd_ull) }}<span class="summary-unit">%</span></span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- ══════════════════════════════════════════
               Feeder 탭 (상단 채널선택 전체폭 + 하단 상세)
          ══════════════════════════════════════════ -->
          <div v-else-if="activeTab === 'feeder'" class="flex flex-col gap-6">

            <!-- 1단: Channel Status (col-12, 전체폭) -->
            <div class="meter-card p-5">
              <div class="meter-card-header !px-0 !py-0 mb-4">
                <h3 class="meter-card-title meter-accent-blue">Channel Status</h3>
                <div class="flex flex-wrap items-center gap-2 text-xs">
                  <span v-for="s in statusDefs" :key="s.key" class="flex items-center gap-1">
                    <span class="w-2 h-2 rounded-full" :style="{ backgroundColor: s.color }"></span>
                    <span class="text-gray-500 dark:text-gray-400">{{ s.label }}</span>
                  </span>
                </div>
              </div>

              <div class="grid grid-cols-[repeat(16,minmax(0,1fr))] gap-1.5 p-1">
                <button
                  v-for="idx in MAX_FEEDERS"
                  :key="idx"
                  @click="onFeederClick(idx)"
                  class="flex flex-col items-center justify-center rounded-lg border-2 py-2 px-1 transition-all hover:scale-105 hover:shadow-md cursor-pointer"
                  :class="[getDeviceClass(idx), selectedDevice === idx ? 'ring-2 ring-offset-1 ring-blue-500 dark:ring-offset-gray-800' : '']"
                  :style="getDeviceBorderStyle(idx)"
                >
                  <span class="text-sm font-bold leading-none">{{ idx }}</span>
                  <span class="text-[10px] font-semibold tracking-tight mt-1">{{ getDeviceStatus(idx) }}</span>
                </button>
              </div>
            </div>

            <!-- 2단: Measurement + Energy (6-6) -->
            <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
              <!-- 상세 데이터 테이블 -->
              <transition name="slide-fade">
                <div v-if="selectedDevice" class="lg:col-span-6 meter-card p-5">
                  <div class="meter-card-header !px-0 !py-0 mb-4">
                    <h3 class="meter-card-title meter-accent-emerald">Measurement</h3>
                    <span class="text-xs px-2 py-1 rounded-full font-semibold border-2"
                          :class="getDeviceClass(selectedDevice)"
                          :style="getDeviceBorderStyle(selectedDevice)">
                      CH {{ selectedDevice }}
                    </span>
                  </div>

                  <div class="overflow-x-auto">
                    <table class="w-full text-sm min-w-[400px]">
                      <thead>
                        <tr class="border-b-2 border-gray-200 dark:border-gray-600">
                          <th class="text-left py-2 px-3 font-semibold text-gray-500 dark:text-gray-400 w-28">Items</th>
                          <th class="text-right py-2 px-3 font-semibold text-gray-500 dark:text-gray-400">Total</th>
                          <th class="text-right py-2 px-3 font-semibold text-gray-500 dark:text-gray-400">L1</th>
                          <th class="text-right py-2 px-3 font-semibold text-gray-500 dark:text-gray-400">L2</th>
                          <th class="text-right py-2 px-3 font-semibold text-gray-500 dark:text-gray-400">L3</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr v-for="(row, i) in measurementRows" :key="row.item"
                            :class="i % 2 === 0 ? 'bg-gray-50 dark:bg-gray-700/30' : 'bg-white dark:bg-gray-800'">
                          <td class="py-2.5 px-3 font-semibold text-gray-700 dark:text-gray-200">{{ row.item }}</td>
                          <td class="py-2.5 px-3 text-right font-mono text-gray-900 dark:text-gray-100 tabular-nums">{{ row.total ?? '-' }}</td>
                          <td class="py-2.5 px-3 text-right font-mono text-gray-900 dark:text-gray-100 tabular-nums">{{ row.l1 ?? '-' }}</td>
                          <td class="py-2.5 px-3 text-right font-mono text-gray-900 dark:text-gray-100 tabular-nums">{{ row.l2 ?? '-' }}</td>
                          <td class="py-2.5 px-3 text-right font-mono text-gray-900 dark:text-gray-100 tabular-nums">{{ row.l3 ?? '-' }}</td>
                        </tr>
                      </tbody>
                    </table>
                  </div>

                  <!-- 요약 스탯 -->
                  <div v-if="stats" class="grid grid-cols-5 gap-2 mt-4">
                    <div class="stat-box">
                      <span class="stat-label">In(A)</span>
                      <span class="stat-value">{{ stats.in }}</span>
                    </div>
                    <div class="stat-box">
                      <span class="stat-label">Avg I(A)</span>
                      <span class="stat-value">{{ stats.avgI }}</span>
                    </div>
                    <div class="stat-box">
                      <span class="stat-label">Avg THD(%)</span>
                      <span class="stat-value">{{ stats.avgThd }}</span>
                    </div>
                    <div class="stat-box">
                      <span class="stat-label">Ah</span>
                      <span class="stat-value">{{ stats.ah }}</span>
                    </div>
                    <div class="stat-box">
                      <span class="stat-label">Total kWh</span>
                      <span class="stat-value">{{ stats.totalKwh }}</span>
                    </div>
                  </div>
                </div>
              </transition>

              <!-- Energy 카드 -->
              <transition name="slide-fade">
                <div v-if="selectedDevice && feederData" class="lg:col-span-6 meter-card p-5">
                <div class="meter-card-header !px-0 !py-0 mb-4">
                  <h3 class="meter-card-title meter-accent-blue">Energy</h3>
                  <span v-if="selectedDevice" class="text-xs px-2 py-1 rounded-full font-semibold border-2"
                        :class="getDeviceClass(selectedDevice)"
                        :style="getDeviceBorderStyle(selectedDevice)">
                    CH {{ selectedDevice }} · {{ getDeviceStatus(selectedDevice) }}
                  </span>
                </div>
                <!-- Active / Reactive / Apparent 총합 -->
                <div class="grid grid-cols-1 md:grid-cols-3 gap-3 mb-4">
                  <div class="energy-group">
                    <div class="energy-title">Active (kWh)</div>
                    <div class="energy-row"><span>Import</span><span>{{ fmt(feederData.total_import_kwh) }}</span></div>
                    <div class="energy-row"><span>Export</span><span>{{ fmt(feederData.total_export_kwh) }}</span></div>
                    <div class="energy-row"><span>Total (I+E)</span><span>{{ fmt(feederData.total_kwh) }}</span></div>
                    <div class="energy-row"><span>Cal Sum (I-E)</span><span>{{ fmt(feederData.total_kwh_cal) }}</span></div>
                  </div>
                  <div class="energy-group">
                    <div class="energy-title">Reactive (kVArh)</div>
                    <div class="energy-row"><span>Import</span><span>{{ fmt(feederData.total_import_kvarh) }}</span></div>
                    <div class="energy-row"><span>Export</span><span>{{ fmt(feederData.total_export_kvarh) }}</span></div>
                    <div class="energy-row"><span>Total (I+E)</span><span>{{ fmt(feederData.total_kvarh) }}</span></div>
                    <div class="energy-row"><span>Cal Sum (I-E)</span><span>{{ fmt(feederData.total_kvarh_cal) }}</span></div>
                  </div>
                  <div class="energy-group">
                    <div class="energy-title">Apparent (kVAh)</div>
                    <div class="energy-row"><span>Total</span><span>{{ fmt(feederData.total_kvah) }}</span></div>
                    <div class="energy-row"><span>Ah</span><span>{{ fmt(feederData.ah) }}</span></div>
                  </div>
                </div>
                <!-- Phase(L1/L2/L3)를 컬럼으로 전치한 상세 -->
                <div class="overflow-x-auto">
                  <table class="w-full text-sm">
                    <thead>
                      <tr class="border-b-2 border-gray-200 dark:border-gray-600">
                        <th class="text-left py-2 px-3 font-semibold text-gray-500 dark:text-gray-400 w-40">Item</th>
                        <th class="text-right py-2 px-3 font-semibold text-gray-500 dark:text-gray-400">L1</th>
                        <th class="text-right py-2 px-3 font-semibold text-gray-500 dark:text-gray-400">L2</th>
                        <th class="text-right py-2 px-3 font-semibold text-gray-500 dark:text-gray-400">L3</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="(row, i) in energyRows" :key="row.item"
                          :class="i % 2 === 0 ? 'bg-gray-50 dark:bg-gray-700/30' : 'bg-white dark:bg-gray-800'">
                        <td class="py-2 px-3 font-semibold text-gray-700 dark:text-gray-200">{{ row.item }}</td>
                        <td class="py-2 px-3 text-right font-mono tabular-nums">{{ row.l1 }}</td>
                        <td class="py-2 px-3 text-right font-mono tabular-nums">{{ row.l2 }}</td>
                        <td class="py-2 px-3 text-right font-mono tabular-nums">{{ row.l3 }}</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </transition>
            </div>
            <!-- /2단 grid (Measurement + Energy) -->

            <!-- 3단: Phasor + Demand (6-6) -->
            <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
              <!-- Phasor Diagram -->
              <div class="lg:col-span-6 meter-card p-5 flex flex-col">
                <div class="meter-card-header !px-0 !py-0 mb-3">
                  <h3 class="meter-card-title meter-accent-indigo">Phasor Diagram</h3>
                  <span v-if="selectedDevice" class="text-xs px-2 py-1 rounded-full font-semibold border-2"
                        :class="getDeviceClass(selectedDevice)"
                        :style="getDeviceBorderStyle(selectedDevice)">
                    CH {{ selectedDevice }} · {{ getDeviceStatus(selectedDevice) }}
                  </span>
                </div>
                <div ref="phasorContainer" class="flex justify-center items-center" style="height: 320px;">
                  <canvas ref="phasorCanvasRef" class="max-w-full h-auto"></canvas>
                </div>
              </div>

              <!-- Demand 카드 -->
              <transition name="slide-fade">
                <div v-if="selectedDevice && feederData" class="lg:col-span-6 meter-card p-5">
                  <div class="meter-card-header !px-0 !py-0 mb-4">
                    <h3 class="meter-card-title meter-accent-indigo">Demand</h3>
                    <span v-if="selectedDevice" class="text-xs px-2 py-1 rounded-full font-semibold border-2"
                          :class="getDeviceClass(selectedDevice)"
                          :style="getDeviceBorderStyle(selectedDevice)">
                      CH {{ selectedDevice }} · {{ getDeviceStatus(selectedDevice) }}
                    </span>
                  </div>
                  <div class="overflow-x-auto">
                    <table class="w-full text-sm">
                      <thead>
                        <tr class="border-b-2 border-gray-200 dark:border-gray-600">
                          <th class="text-left py-2 px-3 font-semibold text-gray-500 dark:text-gray-400 w-28">Item</th>
                          <th class="text-right py-2 px-3 font-semibold text-gray-500 dark:text-gray-400">Present</th>
                          <th class="text-right py-2 px-3 font-semibold text-gray-500 dark:text-gray-400">Max</th>
                          <th class="text-right py-2 px-3 font-semibold text-gray-500 dark:text-gray-400">Unit</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr v-for="(row, i) in demandRows" :key="row.item"
                            :class="i % 2 === 0 ? 'bg-gray-50 dark:bg-gray-700/30' : 'bg-white dark:bg-gray-800'">
                          <td class="py-2 px-3 font-semibold text-gray-700 dark:text-gray-200">{{ row.item }}</td>
                          <td class="py-2 px-3 text-right font-mono tabular-nums">{{ row.present ?? '-' }}</td>
                          <td class="py-2 px-3 text-right font-mono tabular-nums">{{ row.max ?? '-' }}</td>
                          <td class="py-2 px-3 text-right font-mono text-gray-400 dark:text-gray-500">{{ row.unit }}</td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>
              </transition>
            </div>

          </div>
        </div>
      </main>
      <Footer />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import Sidebar from '../common/SideBar.vue'
import Header from '../common/Header.vue'
import Footer from '../common/Footer.vue'

const { t } = useI18n()
const sidebarOpen = ref(false)

const MAX_FEEDERS = 32
const activeTab = ref('feeder') // 'sv' | 'feeder'

// ── System / Voltage 데이터 ──
// voltageData 스키마: { u, avg_u, ull, avg_ull, thd_u, avg_thd_u, thd_ull, avg_thd_ull, freq }
const voltageData = ref(null)
const svTimestamp = ref(null)

const formattedTimestamp = computed(() => svTimestamp.value ? formatTimestamp(svTimestamp.value) : '--')

function formatTimestamp(ts) {
  if (!ts) return ''
  const d = new Date(ts * 1000)
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

// ── Feeder 상태 / 선택 ──
const devices = ref(
  Array.from({ length: MAX_FEEDERS }, (_, i) => ({ id: i + 1, status: 'OFF' }))
)
const selectedDevice = ref(1)
const feederData = ref(null)

// ── 상태 정의 ──
const statusDefs = [
  { key: 'OFF',  label: 'OFF',  color: '#9ca3af' },
  { key: 'ON',   label: 'ON',   color: '#22c55e' },
  { key: 'TRIP', label: 'TRIP', color: '#ef4444' },
  { key: 'OPR',  label: 'OPR',  color: '#3b82f6' },
  { key: 'OC',   label: 'OC',   color: '#f97316' },
  { key: 'OCG',  label: 'OCG',  color: '#a855f7' },
  { key: 'OCG+', label: 'OCG+', color: '#92400e' },
  { key: 'BAD',  label: 'BAD',  color: '#1f2937' },
]
const statusColorMap = Object.fromEntries(statusDefs.map(s => [s.key, s.color]))

function mapFeederStatus(code) {
  if (code == null) return 'OFF'
  return code > 0 ? 'ON' : 'OFF'
}

const getDeviceStatus = (idx) => devices.value[idx - 1]?.status || 'OFF'
const getDeviceColor = (idx) => statusColorMap[getDeviceStatus(idx)] || '#9ca3af'
const getDeviceBorderStyle = (idx) => ({ borderColor: getDeviceColor(idx) })
const getDeviceClass = (idx) => {
  switch (getDeviceStatus(idx)) {
    case 'ON':   return 'bg-green-50 dark:bg-green-900/20 text-green-700 dark:text-green-300'
    case 'TRIP': return 'bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-300 animate-pulse'
    case 'OPR':  return 'bg-blue-50 dark:bg-blue-900/20 text-blue-700 dark:text-blue-300'
    case 'OC':   return 'bg-orange-50 dark:bg-orange-900/20 text-orange-700 dark:text-orange-300'
    case 'OCG':  return 'bg-purple-50 dark:bg-purple-900/20 text-purple-700 dark:text-purple-300'
    case 'OCG+': return 'bg-amber-50 dark:bg-amber-900/20 text-amber-800 dark:text-amber-400'
    case 'BAD':  return 'bg-gray-100 dark:bg-gray-900 text-gray-900 dark:text-gray-300'
    default:     return 'bg-white dark:bg-gray-800 text-gray-400 dark:text-gray-500'
  }
}

function onFeederClick(idx) {
  selectedDevice.value = selectedDevice.value === idx ? null : idx
  if (selectedDevice.value) {
    fetchFeeder(selectedDevice.value)
  }
}

const fmt = (v) => v != null ? v.toFixed(2) : '-'

// ── Feeder 측정 테이블 ──
// feederData 스키마: i:[I1,I2,I3], avg_i, total_i, in, thd_i:[...], avg_thd_i,
//                   p:[...], p_total, s:[...], s_total, q:[...], q_total,
//                   pf:[...], pf_total, ph:[...], ph_total, ...energy/demand
const measurementRows = computed(() => {
  const fd = feederData.value
  const vd = voltageData.value
  if (!selectedDevice.value || !fd) return []
  return [
    { item: 'V (L-N)', total: fmt(vd?.avg_u),     l1: fmt(vd?.u?.[0]),    l2: fmt(vd?.u?.[1]),    l3: fmt(vd?.u?.[2]) },
    { item: 'V (L-L)', total: fmt(vd?.avg_ull),   l1: fmt(vd?.ull?.[0]),  l2: fmt(vd?.ull?.[1]),  l3: fmt(vd?.ull?.[2]) },
    { item: 'I',       total: fmt(fd.total_i),    l1: fmt(fd.i?.[0]),     l2: fmt(fd.i?.[1]),     l3: fmt(fd.i?.[2]) },
    { item: 'W',       total: fmt(fd.p_total),    l1: fmt(fd.p?.[0]),     l2: fmt(fd.p?.[1]),     l3: fmt(fd.p?.[2]) },
    { item: 'var',     total: fmt(fd.q_total),    l1: fmt(fd.q?.[0]),     l2: fmt(fd.q?.[1]),     l3: fmt(fd.q?.[2]) },
    { item: 'VA',      total: fmt(fd.s_total),    l1: fmt(fd.s?.[0]),     l2: fmt(fd.s?.[1]),     l3: fmt(fd.s?.[2]) },
    { item: 'PF',      total: fmt(fd.pf_total),   l1: fmt(fd.pf?.[0]),    l2: fmt(fd.pf?.[1]),    l3: fmt(fd.pf?.[2]) },
    { item: 'PH (°)',  total: fmt(fd.ph_total),   l1: fmt(fd.ph?.[0]),    l2: fmt(fd.ph?.[1]),    l3: fmt(fd.ph?.[2]) },
    { item: 'I THD(%)',total: fmt(fd.avg_thd_i),  l1: fmt(fd.thd_i?.[0]), l2: fmt(fd.thd_i?.[1]), l3: fmt(fd.thd_i?.[2]) },
  ]
})

const stats = computed(() => {
  const fd = feederData.value
  if (!selectedDevice.value || !fd) return null
  return {
    in: fmt(fd.in),
    avgI: fmt(fd.avg_i),
    avgThd: fmt(fd.avg_thd_i),
    ah: fmt(fd.ah),
    totalKwh: fmt(fd.total_kwh),
  }
})

// ── Energy 상세 테이블 (행: 항목, 컬럼: L1/L2/L3) ──
const energyRows = computed(() => {
  const fd = feederData.value
  if (!fd) return []
  return [
    { item: 'Import kWh',   l1: fmt(fd.l1_import_kwh),   l2: fmt(fd.l2_import_kwh),   l3: fmt(fd.l3_import_kwh) },
    { item: 'Export kWh',   l1: fmt(fd.l1_export_kwh),   l2: fmt(fd.l2_export_kwh),   l3: fmt(fd.l3_export_kwh) },
    { item: 'Total kWh',    l1: fmt(fd.l1_total_kwh),    l2: fmt(fd.l2_total_kwh),    l3: fmt(fd.l3_total_kwh) },
    { item: 'Import kVArh', l1: fmt(fd.l1_import_kvarh), l2: fmt(fd.l2_import_kvarh), l3: fmt(fd.l3_import_kvarh) },
    { item: 'Export kVArh', l1: fmt(fd.l1_export_kvarh), l2: fmt(fd.l2_export_kvarh), l3: fmt(fd.l3_export_kvarh) },
    { item: 'Total kVArh',  l1: fmt(fd.l1_total_kvarh),  l2: fmt(fd.l2_total_kvarh),  l3: fmt(fd.l3_total_kvarh) },
  ]
})

// ── Demand 테이블 ──
const demandRows = computed(() => {
  const fd = feederData.value
  if (!fd) return []
  return [
    { item: 'I1 Demand',    present: fmt(fd.i1_demand),     max: fmt(fd.max_i1_demand),     unit: 'A' },
    { item: 'I2 Demand',    present: fmt(fd.i2_demand),     max: fmt(fd.max_i2_demand),     unit: 'A' },
    { item: 'I3 Demand',    present: fmt(fd.i3_demand),     max: fmt(fd.max_i3_demand),     unit: 'A' },
    { item: 'Total W',      present: fmt(fd.total_w_demand),max: fmt(fd.max_total_w_demand),unit: 'W' },
    { item: 'Total VA',     present: fmt(fd.total_va_demand),max:fmt(fd.max_total_va_demand),unit:'VA' },
    { item: 'Total VAr',    present: fmt(fd.total_var_demand),max:fmt(fd.max_total_var_demand),unit:'VAr' },
  ]
})

// ── Phasor Diagram ──
const phasorCanvasRef = ref(null)
const phasorContainer = ref(null)
const phasorColors = ['#6b7280', '#f97316', '#0ea5e9', '#1e293b', '#dc2626', '#2563eb']
const phasorData = ref({
  degree:    [0, 0, 0, 0, 0, 0],
  magnitude: [0, 0, 0, 0, 0, 0],
  maxlist:   [1, 1, 1, 1, 1, 1],
  texts:     ['U1', 'U2', 'U3', 'I1', 'I2', 'I3'],
})

// 페이저 데이터 구성
// - U 각도: 표준 3상 기준 U1=0°, U2=-120°, U3=+120° (전압 각도 스펙 미제공)
// - I 각도: feederData.ph[0..2] (PH1/PH2/PH3)
function applyPhasor() {
  const fd = feederData.value
  const vd = voltageData.value
  const u = vd?.u || [0, 0, 0]
  const i = fd?.i || [0, 0, 0]
  const uAngles = [0, -120, 120]
  const iAngles = fd?.ph || [0, 0, 0]
  const uMax = Math.max(...u, 1)
  const iMax = Math.max(...i, 1)
  phasorData.value = {
    degree:    [uAngles[0], uAngles[1], uAngles[2], iAngles[0], iAngles[1], iAngles[2]],
    magnitude: [u[0], u[1], u[2], i[0], i[1], i[2]],
    maxlist:   [uMax, uMax, uMax, iMax, iMax, iMax],
    texts:     ['U1', 'U2', 'U3', 'I1', 'I2', 'I3'],
  }
}

function drawPhasor() {
  const canvas = phasorCanvasRef.value
  const container = phasorContainer.value
  if (!canvas || !container) return

  const dpr = window.devicePixelRatio || 1
  const rect = container.getBoundingClientRect()
  const displaySize = Math.min(rect.width, rect.height, 460)

  canvas.width = displaySize * dpr
  canvas.height = displaySize * dpr
  canvas.style.width = displaySize + 'px'
  canvas.style.height = displaySize + 'px'

  const ctx = canvas.getContext('2d')
  if (!ctx) return

  ctx.scale(dpr, dpr)

  const size = displaySize
  const center = size / 2
  const radius = center * 0.78

  ctx.clearRect(0, 0, size, size)
  ctx.translate(center, center)

  ctx.beginPath()
  ctx.arc(0, 0, radius, 0, 2 * Math.PI)
  ctx.fillStyle = '#f3f4f6'
  ctx.fill()

  ctx.setLineDash([4, 3])
  ctx.strokeStyle = '#d1d5db'
  ctx.lineWidth = 1
  for (let i = 1; i <= 3; i++) {
    ctx.beginPath()
    ctx.arc(0, 0, radius * i / 3, 0, 2 * Math.PI)
    ctx.stroke()
  }
  for (let i = 0; i < 12; i++) {
    const angle = i * Math.PI / 6
    ctx.beginPath()
    ctx.moveTo(0, 0)
    ctx.lineTo(Math.cos(angle) * radius, -Math.sin(angle) * radius)
    ctx.stroke()
  }
  ctx.setLineDash([])

  ctx.font = `${radius * 0.08}px sans-serif`
  ctx.fillStyle = '#8b5cf6'
  ctx.textAlign = 'center'
  ctx.textBaseline = 'middle'
  ;[0, 60, 120, 180, 240, 300].forEach((deg) => {
    const rad = (90 - deg) * Math.PI / 180
    ctx.fillText(`${deg}°`, Math.cos(rad) * radius * 1.20, -Math.sin(rad) * radius * 1.20)
  })

  const d = phasorData.value
  for (let i = 0; i < 6; i++) {
    const angleRad = (90 - d.degree[i]) * Math.PI / 180
    const ratio = d.maxlist[i] === 0 ? 0 : d.magnitude[i] / d.maxlist[i]
    const len = i < 3 ? radius * ratio : radius * 0.55 * ratio
    if (len <= 0) continue

    const endX = Math.cos(angleRad) * len
    const endY = -Math.sin(angleRad) * len

    ctx.save()
    ctx.strokeStyle = phasorColors[i]
    ctx.fillStyle = phasorColors[i]
    ctx.lineWidth = 3
    ctx.lineCap = 'round'

    ctx.beginPath()
    ctx.moveTo(0, 0)
    ctx.lineTo(endX, endY)
    ctx.stroke()

    const headLen = 10
    const angle = Math.atan2(-endY, endX)
    ctx.beginPath()
    ctx.moveTo(endX, endY)
    ctx.lineTo(endX - headLen * Math.cos(angle - 0.4), endY + headLen * Math.sin(angle - 0.4))
    ctx.moveTo(endX, endY)
    ctx.lineTo(endX - headLen * Math.cos(angle + 0.4), endY + headLen * Math.sin(angle + 0.4))
    ctx.stroke()

    ctx.font = 'bold 12px sans-serif'
    const perpX = Math.sin(angleRad) * 16
    const perpY = Math.cos(angleRad) * 16
    const extX = Math.cos(angleRad) * 8
    const extY = -Math.sin(angleRad) * 8
    ctx.fillText(d.texts[i], endX + extX + perpX, endY + extY + perpY)
    ctx.restore()
  }

  ctx.beginPath()
  ctx.arc(0, 0, 3, 0, 2 * Math.PI)
  ctx.fillStyle = '#333'
  ctx.fill()
}

let resizeTimer = null
function onResize() {
  if (resizeTimer) clearTimeout(resizeTimer)
  resizeTimer = setTimeout(() => drawPhasor(), 150)
}

watch([feederData, voltageData], () => {
  const fd = feederData.value
  if (fd) {
    const sel = selectedDevice.value
    if (sel >= 1 && sel <= devices.value.length) {
      devices.value[sel - 1].status = mapFeederStatus(fd.feeder_status)
    }
  }
  applyPhasor()
  nextTick(() => drawPhasor())
})

// ── Fetch (API 엔드포인트는 데이터 스펙 확정 후 연결) ──
let pollInterval = null

// ═══════════════════════════════════════════
// 더미 데이터 (API 스펙 확정 후 교체 예정)
// ═══════════════════════════════════════════
const rnd = (min, max, dec = 2) => +(Math.random() * (max - min) + min).toFixed(dec)

function makeDummyFeeder(idx) {
  const baseI = 20 + (idx % 5) * 5
  const baseP = baseI * 220
  return {
    feeder_status: idx % 7 === 0 ? 0 : 1,
    i:       [rnd(baseI, baseI + 5), rnd(baseI, baseI + 5), rnd(baseI, baseI + 5)],
    avg_i:   rnd(baseI, baseI + 5),
    total_i: rnd(baseI * 3, baseI * 3 + 10),
    in:      rnd(0, 2),
    thd_i:    [rnd(1, 5), rnd(1, 5), rnd(1, 5)],
    avg_thd_i: rnd(1, 5),
    p:       [rnd(baseP, baseP + 200), rnd(baseP, baseP + 200), rnd(baseP, baseP + 200)],
    p_total: rnd(baseP * 3, baseP * 3 + 500),
    s:       [rnd(baseP, baseP + 200), rnd(baseP, baseP + 200), rnd(baseP, baseP + 200)],
    s_total: rnd(baseP * 3, baseP * 3 + 500),
    q:       [rnd(100, 300), rnd(100, 300), rnd(100, 300)],
    q_total: rnd(300, 900),
    pf:      [rnd(0.85, 1), rnd(0.85, 1), rnd(0.85, 1)],
    pf_total:rnd(0.85, 1),
    ph:      [rnd(-10, 10), rnd(110, 130), rnd(-130, -110)],
    ph_total:rnd(-10, 10),
    // Energy
    total_import_kwh: rnd(1000, 5000),
    total_export_kwh: rnd(100, 500),
    total_kwh:        rnd(1100, 5500),
    total_kwh_cal:    rnd(900, 4500),
    total_import_kvarh: rnd(200, 800),
    total_export_kvarh: rnd(50, 200),
    total_kvarh:        rnd(250, 1000),
    total_kvarh_cal:    rnd(150, 600),
    total_kvah:         rnd(1200, 6000),
    ah:                 rnd(500, 2000),
    l1_import_kwh: rnd(300, 1500), l2_import_kwh: rnd(300, 1500), l3_import_kwh: rnd(300, 1500),
    l1_export_kwh: rnd(30, 150),   l2_export_kwh: rnd(30, 150),   l3_export_kwh: rnd(30, 150),
    l1_total_kwh:  rnd(330, 1650), l2_total_kwh:  rnd(330, 1650), l3_total_kwh:  rnd(330, 1650),
    l1_import_kvarh: rnd(60, 250), l2_import_kvarh: rnd(60, 250), l3_import_kvarh: rnd(60, 250),
    l1_export_kvarh: rnd(15, 60),  l2_export_kvarh: rnd(15, 60),  l3_export_kvarh: rnd(15, 60),
    l1_total_kvarh:  rnd(75, 310), l2_total_kvarh:  rnd(75, 310), l3_total_kvarh:  rnd(75, 310),
    // Demand
    i1_demand: rnd(baseI, baseI + 5), i2_demand: rnd(baseI, baseI + 5), i3_demand: rnd(baseI, baseI + 5),
    max_i1_demand: rnd(baseI + 5, baseI + 10), max_i2_demand: rnd(baseI + 5, baseI + 10), max_i3_demand: rnd(baseI + 5, baseI + 10),
    total_w_demand:     rnd(baseP * 3, baseP * 3 + 500),
    max_total_w_demand: rnd(baseP * 3 + 500, baseP * 3 + 1000),
    total_va_demand:    rnd(baseP * 3, baseP * 3 + 500),
    max_total_va_demand:rnd(baseP * 3 + 500, baseP * 3 + 1000),
    total_var_demand:    rnd(300, 900),
    max_total_var_demand:rnd(900, 1500),
  }
}

function makeDummyVoltage() {
  return {
    u:       [rnd(218, 222), rnd(218, 222), rnd(218, 222)],
    avg_u:    rnd(218, 222),
    ull:     [rnd(378, 382), rnd(378, 382), rnd(378, 382)],
    avg_ull:  rnd(378, 382),
    thd_u:    [rnd(1, 3), rnd(1, 3), rnd(1, 3)],
    avg_thd_u: rnd(1, 3),
    thd_ull:  [rnd(1, 3), rnd(1, 3), rnd(1, 3)],
    avg_thd_ull: rnd(1, 3),
    freq: rnd(59.9, 60.1),
  }
}

async function fetchFeederAll() {
  // 더미: 첫 20개 피더만 구성됨
  const map = {}
  for (let i = 1; i <= 20; i++) map[i] = makeDummyFeeder(i)
  // 상태를 devices에도 반영
  devices.value.forEach((d, idx) => {
    const n = idx + 1
    d.status = map[n] ? mapFeederStatus(map[n].feeder_status) : 'OFF'
  })
}

async function fetchFeeder(id) {
  feederData.value = makeDummyFeeder(id)
}

async function fetchSystemVoltage() {
  voltageData.value = makeDummyVoltage()
  svTimestamp.value = Math.floor(Date.now() / 1000)
}

function stopPolling() {
  if (pollInterval) { clearInterval(pollInterval); pollInterval = null }
}

function startPollingForTab(tab) {
  stopPolling()
  if (tab === 'sv') {
    pollInterval = setInterval(fetchSystemVoltage, 5000)
  } else if (tab === 'feeder') {
    // Feeder 탭도 V/Vll 표시 및 페이저에 전압 사용 → voltage 함께 폴링
    pollInterval = setInterval(async () => {
      await Promise.all([fetchFeederAll(), fetchSystemVoltage()])
      if (selectedDevice.value) await fetchFeeder(selectedDevice.value)
    }, 5000)
  }
}

async function loadTabData(tab) {
  if (tab === 'sv') {
    await fetchSystemVoltage()
  } else if (tab === 'feeder') {
    await Promise.all([fetchFeederAll(), fetchSystemVoltage()])
    if (selectedDevice.value) await fetchFeeder(selectedDevice.value)
  }
}

onMounted(async () => {
  await loadTabData(activeTab.value)
  startPollingForTab(activeTab.value)
  await nextTick()
  drawPhasor()
  window.addEventListener('resize', onResize)
})

onUnmounted(() => {
  stopPolling()
  window.removeEventListener('resize', onResize)
  if (resizeTimer) clearTimeout(resizeTimer)
})

watch(activeTab, async (tab) => {
  await loadTabData(tab)
  startPollingForTab(tab)
  if (tab === 'feeder') {
    await nextTick()
    drawPhasor()
  }
})
</script>

<style scoped>
@import '../../css/meter-card.css';

.slide-fade-enter-active { transition: all 0.3s ease-out; }
.slide-fade-leave-active { transition: all 0.2s ease-in; }
.slide-fade-enter-from, .slide-fade-leave-to { transform: translateY(-10px); opacity: 0; }

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

.energy-group {
  @apply rounded-lg p-4;
  @apply bg-gray-50 dark:bg-gray-700/40 border border-gray-100 dark:border-gray-700/60;
}
.energy-title {
  @apply text-sm font-bold uppercase text-gray-600 dark:text-gray-300 mb-3 pb-2 border-b border-gray-200 dark:border-gray-600 tracking-wider;
}
.energy-row {
  @apply flex justify-between items-center py-1.5 text-base;
  @apply text-gray-800 dark:text-gray-100 font-mono font-semibold tabular-nums;
}
.energy-row span:first-child {
  @apply text-xs text-gray-500 dark:text-gray-400 font-sans font-medium;
}

/* Voltage 섹션 (System/Voltage 탭) */
.voltage-section {
  @apply rounded-xl p-4 bg-gray-50/60 dark:bg-gray-700/20 border border-gray-200/60 dark:border-gray-700/60;
}
.voltage-summary {
  @apply mt-3 pt-3 flex items-center justify-around gap-3 border-t border-gray-200 dark:border-gray-700/70;
}
.summary-item {
  @apply flex flex-col items-center flex-1 gap-0.5;
}
.summary-label {
  @apply text-[11px] font-bold text-gray-400 dark:text-gray-500 uppercase tracking-wider;
}
.summary-value {
  @apply text-lg font-mono font-bold text-gray-800 dark:text-gray-100 tabular-nums;
}
.summary-unit {
  @apply text-xs font-semibold text-gray-400 dark:text-gray-500 ml-1;
}
.summary-divider {
  @apply w-px h-8 bg-gray-200 dark:bg-gray-700/70;
}
.voltage-cell {
  @apply relative flex flex-col items-center py-4 px-3 rounded-lg;
  @apply bg-gradient-to-b from-gray-50 to-white dark:from-gray-700/30 dark:to-gray-800/30;
  @apply border border-gray-200/60 dark:border-gray-700/60;
}
.voltage-label {
  @apply text-xs font-bold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1;
}
.voltage-value {
  @apply text-xl md:text-2xl font-mono font-bold text-gray-900 dark:text-gray-100 tabular-nums leading-tight;
}
.voltage-unit {
  @apply text-sm font-semibold text-gray-400 dark:text-gray-500;
}
.voltage-thd {
  @apply flex items-center gap-2 mt-3 pt-3 w-full justify-center border-t border-gray-200/70 dark:border-gray-700/70;
}
.thd-label {
  @apply text-xs font-bold text-gray-400 dark:text-gray-500 uppercase tracking-wider;
}
.thd-value {
  @apply text-xl md:text-2xl font-mono font-bold text-indigo-600 dark:text-indigo-400 tabular-nums;
}
.thd-unit {
  @apply text-sm text-gray-400 dark:text-gray-500 ml-0.5 font-semibold;
}
</style>
