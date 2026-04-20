<template>
  <div class="flex h-[100dvh] overflow-hidden">
    <Sidebar :sidebarOpen="sidebarOpen" @close-sidebar="sidebarOpen = false" />

    <div class="relative flex flex-col flex-1 overflow-y-auto overflow-x-hidden">
      <Header :sidebarOpen="sidebarOpen" @toggle-sidebar="sidebarOpen = !sidebarOpen" />

      <main class="grow">
        <div class="px-2 sm:px-4 lg:px-6 py-4 w-full max-w-full">
          <div class="mb-4">
            <h2 class="text-xl md:text-2xl text-gray-800 dark:text-gray-100 font-bold">
              Branch > iPSM #{{ id }}
            </h2>
          </div>

          <!-- 메인 2컬럼: 좌측 장비상태(5) + 우측 벡터도/상세표(7) -->
          <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">

            <!-- ========== 좌측: 장비 상태 (col-5 ≈ 40%) ========== -->
            <div class="lg:col-span-5 bg-white dark:bg-gray-800 shadow-sm rounded-xl border border-gray-200 dark:border-gray-700/60 p-5 flex flex-col">
              <div class="flex items-center justify-between mb-4 flex-wrap gap-2">
                <span class="inline-block px-3 py-1 text-sm font-semibold text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-500/10 rounded">
                  Channel Status
                </span>
                <div class="flex flex-wrap items-center gap-2 text-xs">
                  <span v-for="s in statusDefs" :key="s.key" class="flex items-center gap-1">
                    <span class="w-2 h-2 rounded-full" :style="{ backgroundColor: s.color }"></span>
                    <span class="text-gray-500 dark:text-gray-400">{{ s.label }}</span>
                  </span>
                </div>
              </div>

              <div class="grid grid-cols-6 gap-1 overflow-y-auto flex-1 content-start p-2" style="max-height: calc(100vh - 220px);">
                <button
                  v-for="idx in 72"
                  :key="idx"
                  @click="onDeviceClick(idx)"
                  class="flex flex-col items-center justify-center rounded-lg border-2 py-2 px-1 transition-all hover:scale-105 hover:shadow-md cursor-pointer"
                  :class="[getDeviceClass(idx), selectedDevice === idx ? 'ring-2 ring-offset-1 ring-blue-500 dark:ring-offset-gray-800' : '']"
                  :style="getDeviceBorderStyle(idx)"
                >
                  <span class="text-base font-bold leading-none">{{ idx }}</span>
                  <span class="text-[10px] font-semibold tracking-tight mt-1">{{ getDeviceStatus(idx) }}</span>
                </button>
              </div>
            </div>

            <!-- ========== 우측 (col-7 ≈ 60%) ========== -->
            <div class="lg:col-span-7 flex flex-col gap-6">

              <!-- 우측 상단: 벡터도 -->
              <div class="bg-white dark:bg-gray-800 shadow-sm rounded-xl border border-gray-200 dark:border-gray-700/60 p-5 flex flex-col">
                <div class="flex items-center justify-between mb-3">
                  <span class="inline-block px-3 py-1 text-sm font-semibold text-indigo-600 dark:text-indigo-400 bg-indigo-50 dark:bg-indigo-500/10 rounded">
                    Phasor Diagram
                  </span>
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

              <!-- 우측 하단: 상세 데이터 테이블 -->
              <transition name="slide-fade">
                <div v-if="selectedDevice" class="bg-white dark:bg-gray-800 shadow-sm rounded-xl border border-gray-200 dark:border-gray-700/60 p-5">
                  <div class="flex items-center justify-between gap-3 mb-4 flex-wrap">
                    <span class="inline-block px-3 py-1 text-sm font-semibold text-indigo-600 dark:text-indigo-400 bg-indigo-50 dark:bg-indigo-500/10 rounded">
                      Measurement
                    </span>
                     <span v-if="selectedDevice" class="text-xs px-2 py-1 rounded-full font-semibold border-2"
                        :class="getDeviceClass(selectedDevice)"
                        :style="getDeviceBorderStyle(selectedDevice)">
                    CH {{ selectedDevice }} · {{ deviceInfo.ptType }}
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
                        <tr v-for="(row, i) in selectedDeviceData" :key="row.item"
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
import { watch, ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import Sidebar from '../common/SideBar.vue'
import Header from '../common/Header.vue'
import Footer from '../common/Footer.vue'

const { t } = useI18n()
const sidebarOpen = ref(false)
const phasorCanvasRef = ref(null)
const phasorContainer = ref(null)

const props = defineProps({
  id: { type: String, required: true },
})

const deviceInfo = ref({
  psupply: 'powers',
  driveType: '직입기동',
  location: '분전반 Smart',
  ptType: '3P4W',
  ratedVoltage: '220 V',
  ratedFreq: '60 Hz',
  ratedCurrent: '150 A',
})

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

// 72채널 장비 생성
const devices = ref(
  Array.from({ length: 72 }, (_, i) => {
    const num = i + 1
    // 다양한 상태 샘플링
    if (num >= 65 && num <= 66) return { id: num, status: 'OFF' }
    if (num === 67) return { id: num, status: 'TRIP' }
    if (num === 68) return { id: num, status: 'OPR' }
    if (num === 69) return { id: num, status: 'OC' }
    if (num === 70) return { id: num, status: 'OCG' }
    if (num === 71) return { id: num, status: 'OCG+' }
    if (num === 72) return { id: num, status: 'BAD' }
    return { id: num, status: 'ON' }
  })
)

const selectedDevice = ref(1)

const getDeviceStatus = (idx) => devices.value[idx - 1]?.status || 'OFF'
const getDeviceColor = (idx) => statusColorMap[getDeviceStatus(idx)] || '#9ca3af'

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

const getDeviceBorderStyle = (idx) => ({ borderColor: getDeviceColor(idx) })

const selectedDeviceData = computed(() => {
  if (!selectedDevice.value) return []
  return [
    { item: 'V',        total: '220.00', l1: '220.00', l2: '220.00', l3: '220.00' },
    { item: 'I',        total: '10.00',  l1: '3.33',   l2: '3.33',   l3: '3.33' },
    { item: 'W',        total: '2200.0', l1: '733.3',  l2: '733.3',  l3: '733.3' },
    { item: 'var',      total: '100.0',  l1: '33.3',   l2: '33.3',   l3: '33.3' },
    { item: 'VA',       total: '2200.0', l1: '733.3',  l2: '733.3',  l3: '733.3' },
    { item: 'PF(%)',    total: '0.95',   l1: '0.95',   l2: '0.95',   l3: '0.95' },
    { item: 'kWh',      total: '100.5',  l1: null,     l2: null,     l3: null },
    { item: 'I ub(%)',  total: '0.00',   l1: null,     l2: null,     l3: null },
    { item: 'I THD(%)', total: null,     l1: '5.00',   l2: '5.00',   l3: '5.00' },
    { item: 'Ig(mA)',   total: '0.123',  l1: null,     l2: null,     l3: null },
  ]
})

// ===== Phasor Diagram =====
const phasorColors = ['#6b7280', '#f97316', '#0ea5e9', '#1e293b', '#dc2626', '#2563eb']

const phasorData = ref({
  degree:    [0, 240, 120, 350, 230, 110],
  magnitude: [220, 220, 220, 3.33, 3.33, 3.33],
  maxlist:   [220, 220, 220, 5, 5, 5],
  texts:     ['U1', 'U2', 'U3', 'I1', 'I2', 'I3'],
})

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

const onDeviceClick = (idx) => {
  selectedDevice.value = selectedDevice.value === idx ? null : idx
  if (selectedDevice.value) {
    loadDeviceDetail(idx)
    nextTick(() => drawPhasor())
  }
}

function loadDeviceDetail(deviceIdx) {
  console.log(`[Branch] loading detail for device #${deviceIdx}`)
}

let resizeTimer = null

function onResize() {
  if (resizeTimer) clearTimeout(resizeTimer)
  resizeTimer = setTimeout(() => drawPhasor(), 150)
}

onMounted(async () => {
  console.log('[Branch] mounted, id:', props.id)
  loadBranchData(props.id)
  await nextTick()
  drawPhasor()
  window.addEventListener('resize', onResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', onResize)
  if (resizeTimer) clearTimeout(resizeTimer)
})

watch(() => props.id, (newId) => {
  selectedDevice.value = 1
  loadBranchData(newId)
  nextTick(() => drawPhasor())
})

function loadBranchData(branchId) {
  console.log(`[Branch] loading data for branch #${branchId}`)
}
</script>

<style scoped>
.slide-fade-enter-active { transition: all 0.3s ease-out; }
.slide-fade-leave-active { transition: all 0.2s ease-in; }
.slide-fade-enter-from, .slide-fade-leave-to { transform: translateY(-10px); opacity: 0; }
</style>