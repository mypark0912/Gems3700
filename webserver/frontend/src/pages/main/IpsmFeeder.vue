<template>
  <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">

    <!-- 좌측: 피더 상태 그리드 (col-5) -->
    <div class="lg:col-span-5 meter-card p-5 flex flex-col">
      <div class="meter-card-header !px-0 !py-0 mb-4">
        <h3 class="meter-card-title meter-accent-blue">Channel Status</h3>
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
          @click="handleClick(idx)"
          class="flex flex-col items-center justify-center rounded-lg border-2 py-2 px-1 transition-all hover:scale-105 hover:shadow-md cursor-pointer"
          :class="[getDeviceClass(idx), selectedDevice === idx ? 'ring-2 ring-offset-1 ring-blue-500 dark:ring-offset-gray-800' : '']"
          :style="getDeviceBorderStyle(idx)"
        >
          <span class="text-base font-bold leading-none">{{ idx }}</span>
          <span class="text-[10px] font-semibold tracking-tight mt-1">{{ getDeviceStatus(idx) }}</span>
        </button>
      </div>
    </div>

    <!-- 우측 (col-7): 벡터도 + 상세표 -->
    <div class="lg:col-span-7 flex flex-col gap-6">

      <!-- 벡터도 -->
      <div class="meter-card p-5 flex flex-col">
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

      <!-- 상세 데이터 테이블 -->
      <transition name="slide-fade">
        <div v-if="selectedDevice" class="meter-card p-5">
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
              <span class="stat-label">kWh</span>
              <span class="stat-value">{{ stats.kwh }}</span>
            </div>
            <div class="stat-box">
              <span class="stat-label">V ub(%)</span>
              <span class="stat-value">{{ stats.vub }}</span>
            </div>
            <div class="stat-box">
              <span class="stat-label">I ub(%)</span>
              <span class="stat-value">{{ stats.iub }}</span>
            </div>
            <div class="stat-box">
              <span class="stat-label">THD(%)</span>
              <span class="stat-value">{{ stats.thd }}</span>
            </div>
            <div class="stat-box">
              <span class="stat-label">Ig(A)</span>
              <span class="stat-value">{{ stats.ig }}</span>
            </div>
          </div>
        </div>
      </transition>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'

const props = defineProps({
  feederData: { type: Object, default: null },
  energyData: { type: Object, default: null },
  voltageData: { type: Object, default: null },
  selectedDevice: { type: Number, default: 1 },
})

const emit = defineEmits(['select'])

const phasorCanvasRef = ref(null)
const phasorContainer = ref(null)

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

const devices = ref(
  Array.from({ length: 72 }, (_, i) => ({ id: i + 1, status: 'OFF' }))
)

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

function handleClick(idx) {
  emit('select', idx)
}

const fmt = (v) => v != null ? v.toFixed(2) : '-'

const measurementRows = computed(() => {
  const fd = props.feederData
  const vd = props.voltageData
  if (!props.selectedDevice || !fd) return []
  return [
    { item: 'V',   total: null,           l1: fmt(fd.u?.[0]),  l2: fmt(fd.u?.[1]),  l3: fmt(fd.u?.[2]) },
    { item: 'Vll', total: null,           l1: fmt(vd?.ull?.[0]),l2: fmt(vd?.ull?.[1]),l3: fmt(vd?.ull?.[2]) },
    { item: 'I',   total: fmt(fd.i_total),l1: fmt(fd.i?.[0]),  l2: fmt(fd.i?.[1]),  l3: fmt(fd.i?.[2]) },
    { item: 'W',   total: fmt(fd.p_total),l1: fmt(fd.p?.[0]),  l2: fmt(fd.p?.[1]),  l3: fmt(fd.p?.[2]) },
    { item: 'var', total: fmt(fd.q_total),l1: fmt(fd.q?.[0]),  l2: fmt(fd.q?.[1]),  l3: fmt(fd.q?.[2]) },
    { item: 'VA',  total: fmt(fd.s_total),l1: fmt(fd.s?.[0]),  l2: fmt(fd.s?.[1]),  l3: fmt(fd.s?.[2]) },
    { item: 'PF',  total: fmt(fd.pf_total),l1: fmt(fd.pf?.[0]),l2: fmt(fd.pf?.[1]),l3: fmt(fd.pf?.[2]) },
  ]
})

const stats = computed(() => {
  const fd = props.feederData
  const vd = props.voltageData
  const ed = props.energyData
  if (!props.selectedDevice || !fd) return null
  const thdAvg = Array.isArray(fd.thd_i) ? fd.thd_i.reduce((a, b) => a + b, 0) / fd.thd_i.length : null
  return {
    kwh: fmt(ed?.total_energy?.import_kwh),
    vub: fmt(vd?.u_unbal),
    iub: fmt(fd.i_unbal),
    thd: fmt(thdAvg),
    ig: fmt(fd.ig),
  }
})

// ===== Phasor Diagram =====
const phasorColors = ['#6b7280', '#f97316', '#0ea5e9', '#1e293b', '#dc2626', '#2563eb']
const phasorData = ref({
  degree:    [0, 0, 0, 0, 0, 0],
  magnitude: [0, 0, 0, 0, 0, 0],
  maxlist:   [1, 1, 1, 1, 1, 1],
  texts:     ['U1', 'U2', 'U3', 'I1', 'I2', 'I3'],
})

function applyPhasorFromFeeder(fd) {
  if (!fd) return
  const u = fd.u || [0, 0, 0]
  const i = fd.i || [0, 0, 0]
  const au = fd.angle_u || [0, 0, 0]
  const ai = fd.angle_i || [0, 0, 0]
  const uMax = Math.max(...u, 1)
  const iMax = Math.max(...i, 1)
  phasorData.value = {
    degree:    [au[0], au[1], au[2], ai[0], ai[1], ai[2]],
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

watch(() => props.feederData, (fd) => {
  if (!fd) return
  const sel = props.selectedDevice
  if (sel >= 1 && sel <= devices.value.length) {
    devices.value[sel - 1].status = mapFeederStatus(fd.feeder_status)
  }
  applyPhasorFromFeeder(fd)
  nextTick(() => drawPhasor())
})

onMounted(async () => {
  await nextTick()
  drawPhasor()
  window.addEventListener('resize', onResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', onResize)
  if (resizeTimer) clearTimeout(resizeTimer)
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
</style>
