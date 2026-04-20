<template>
  <Teleport to="body">
    <div v-if="open" class="fixed inset-0 bg-black/40 backdrop-blur-sm z-40" @click="$emit('close')"></div>
  </Teleport>
  <Teleport to="body">
    <div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="$emit('close')">
      <div class="bg-white dark:bg-gray-800 rounded-xl shadow-2xl border border-gray-200 dark:border-gray-700 w-[92vw] max-w-[1100px] flex flex-col" @click.stop>

        <!-- Header -->
        <div class="px-5 py-2.5 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between flex-shrink-0">
          <div class="flex items-center gap-3">
            <h3 class="text-sm font-extrabold text-gray-800 dark:text-gray-100">Data View</h3>
            <span class="bg-violet-500 text-white px-2.5 py-0.5 rounded-full text-[11px] font-bold">iBSM #{{ canid }}</span>
            <span class="text-[10px] font-semibold px-2 py-0.5 rounded bg-gray-100 dark:bg-gray-700 text-gray-500 dark:text-gray-400">{{ cbtypeLabel }}</span>
          </div>
          <div class="flex items-center gap-2">
            <button @click="fetchData" class="px-3 py-1 text-[11px] font-semibold rounded-md border border-violet-500 text-violet-600 dark:text-violet-400 bg-white dark:bg-gray-700 hover:bg-violet-50 dark:hover:bg-violet-900/20" :disabled="loading">
              {{ loading ? 'Loading...' : 'Refresh' }}
            </button>
            <button @click="$emit('close')" class="w-6 h-6 rounded-full bg-gray-100 dark:bg-gray-700 text-gray-500 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-gray-600 flex items-center justify-center text-xs flex-shrink-0">✕</button>
          </div>
        </div>

        <!-- Device Info (compact) -->
        <div v-if="data" class="px-5 py-2 border-b border-gray-100 dark:border-gray-700/50 flex flex-wrap items-center gap-x-5 gap-y-1 text-[11px] bg-gray-50/50 dark:bg-gray-700/10 flex-shrink-0">
          <span><b class="text-gray-400 mr-1">CAN ID</b> {{ data.canid }}</span>
          <span><b class="text-gray-400 mr-1">DevType</b> {{ data.devtype }}</span>
          <span><b class="text-gray-400 mr-1">FW</b> v{{ data.fwver }}</span>
          <span><b class="text-gray-400 mr-1">Status</b> <span :class="data.status===1?'text-emerald-600 dark:text-emerald-400':'text-red-500'">{{ data.status===1?'Normal':'Error' }}</span></span>
          <span><b class="text-gray-400 mr-1">Temp</b> {{ fmt(data.temp) }} °C</span>
          <span><b class="text-gray-400 mr-1">Freq</b> {{ fmt(data.freq) }} Hz</span>
          <span><b class="text-gray-400 mr-1">Time</b> {{ fmtTs(data.ts) }}</span>
        </div>

        <!-- CB Tabs -->
        <div v-if="data" class="flex border-b border-gray-200 dark:border-gray-700 flex-shrink-0">
          <button v-for="ci in cbcount" :key="ci"
            class="px-5 py-2 text-xs font-medium border-b-2 transition-colors"
            :class="activeTab === ci - 1
              ? 'border-violet-500 text-violet-600 dark:text-violet-400'
              : 'border-transparent text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300'"
            @click="activeTab = ci - 1">
            CB{{ ci }}
          </button>
        </div>

        <!-- Body -->
        <div class="flex-1 p-4 flex-shrink-0">
          <div v-if="loading" class="flex items-center justify-center h-32 text-gray-400 text-sm">데이터 요청 중...</div>
          <div v-else-if="error" class="flex items-center justify-center h-32 text-red-400 text-sm">{{ error }}</div>
          <div v-else-if="!data" class="flex items-center justify-center h-32 text-gray-400 text-sm">데이터 없음</div>

          <!-- CB Data Table -->
          <template v-else-if="currentCb">
            <!-- 3상 -->
            <table v-if="is3p(currentCb.cbtype)" class="dv-table">
              <thead>
                <tr>
                  <th class="w-[140px]">항목</th>
                  <th>L1</th>
                  <th>L2</th>
                  <th>L3</th>
                  <th>합계 / 기타</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td class="dv-label">전압 (V)</td>
                  <td>{{ fmt(currentCb.vrms?.[0]) }}</td>
                  <td>{{ fmt(currentCb.vrms?.[1]) }}</td>
                  <td>{{ fmt(currentCb.vrms?.[2]) }}</td>
                  <td class="text-gray-400">—</td>
                </tr>
                <tr>
                  <td class="dv-label">전류 (A)</td>
                  <td>{{ fmt(currentCb.irms?.[0]) }}</td>
                  <td>{{ fmt(currentCb.irms?.[1]) }}</td>
                  <td>{{ fmt(currentCb.irms?.[2]) }}</td>
                  <td class="text-gray-400">—</td>
                </tr>
                <tr>
                  <td class="dv-label">상간전압 (V)</td>
                  <td>{{ fmt(currentCb.vpp?.[0]) }}</td>
                  <td>{{ fmt(currentCb.vpp?.[1]) }}</td>
                  <td>{{ fmt(currentCb.vpp?.[2]) }}</td>
                  <td class="text-gray-400">—</td>
                </tr>
                <tr>
                  <td class="dv-label">유효전력 (W)</td>
                  <td>{{ fmt(currentCb.phwatt?.[0]) }}</td>
                  <td>{{ fmt(currentCb.phwatt?.[1]) }}</td>
                  <td>{{ fmt(currentCb.phwatt?.[2]) }}</td>
                  <td class="font-semibold">{{ fmt(currentCb.power?.[0]) }}</td>
                </tr>
                <tr>
                  <td class="dv-label">무효전력 (var)</td>
                  <td colspan="3" class="text-center text-gray-400">—</td>
                  <td class="font-semibold">{{ fmt(currentCb.power?.[1]) }}</td>
                </tr>
                <tr>
                  <td class="dv-label">피상전력 (VA)</td>
                  <td colspan="3" class="text-center text-gray-400">—</td>
                  <td class="font-semibold">{{ fmt(currentCb.power?.[2]) }}</td>
                </tr>
                <tr>
                  <td class="dv-label">역율 (PF)</td>
                  <td colspan="3" class="text-center text-gray-400">—</td>
                  <td class="font-semibold">{{ fmt(currentCb.pf) }}</td>
                </tr>
                <tr>
                  <td class="dv-label">고조파 (%)</td>
                  <td colspan="3" class="text-center text-gray-400">—</td>
                  <td>{{ fmt(currentCb.pth ?? currentCb.pthd) }}</td>
                </tr>
                <tr>
                  <td class="dv-label">전압 불평형 (%)</td>
                  <td colspan="3" class="text-center text-gray-400">—</td>
                  <td>{{ fmt(currentCb.vunbal) }}</td>
                </tr>
                <tr>
                  <td class="dv-label">전류 불평형 (%)</td>
                  <td colspan="3" class="text-center text-gray-400">—</td>
                  <td>{{ fmt(currentCb.iunbal) }}</td>
                </tr>
                <tr>
                  <td class="dv-label">누설전류 (mA)</td>
                  <td colspan="3" class="text-center text-gray-400">—</td>
                  <td>{{ fmt(currentCb.ig) }}</td>
                </tr>
                <tr class="bg-violet-50/50 dark:bg-violet-900/10">
                  <td class="dv-label">유효전력량 (kWh)</td>
                  <td colspan="3" class="text-center text-gray-400">—</td>
                  <td class="font-semibold">{{ currentCb.kwh }}</td>
                </tr>
                <tr class="bg-violet-50/50 dark:bg-violet-900/10">
                  <td class="dv-label">무효전력량 (kvarh)</td>
                  <td colspan="3" class="text-center text-gray-400">—</td>
                  <td>{{ currentCb.kvarh }}</td>
                </tr>
                <tr class="bg-violet-50/50 dark:bg-violet-900/10">
                  <td class="dv-label">피상전력량 (kVAh)</td>
                  <td colspan="3" class="text-center text-gray-400">—</td>
                  <td>{{ currentCb.kVAh }}</td>
                </tr>
                <tr class="bg-emerald-50/50 dark:bg-emerald-900/10">
                  <td class="dv-label">현월 전력량 (kWh)</td>
                  <td colspan="3" class="text-center text-gray-400">—</td>
                  <td class="font-semibold">{{ currentCb.kwh_tm }}</td>
                </tr>
                <tr class="bg-emerald-50/50 dark:bg-emerald-900/10">
                  <td class="dv-label">전월 전력량 (kWh)</td>
                  <td colspan="3" class="text-center text-gray-400">—</td>
                  <td>{{ currentCb.kwh_lm }}</td>
                </tr>
              </tbody>
            </table>

            <!-- 1상 -->
            <table v-else class="dv-table">
              <thead>
                <tr>
                  <th class="w-[180px]">항목</th>
                  <th>값</th>
                </tr>
              </thead>
              <tbody>
                <tr><td class="dv-label">전압 (V)</td><td>{{ fmt(currentCb.vrms?.[0]) }}</td></tr>
                <tr><td class="dv-label">전류 (A)</td><td>{{ fmt(currentCb.irms?.[0]) }}</td></tr>
                <tr><td class="dv-label">누설전류 (mA)</td><td>{{ fmt(currentCb.ig) }}</td></tr>
                <tr><td class="dv-label">유효전력 (W)</td><td class="font-semibold">{{ fmt(currentCb.power?.[0]) }}</td></tr>
                <tr><td class="dv-label">무효전력 (var)</td><td>{{ fmt(currentCb.power?.[1]) }}</td></tr>
                <tr><td class="dv-label">피상전력 (VA)</td><td>{{ fmt(currentCb.power?.[2]) }}</td></tr>
                <tr><td class="dv-label">역율 (PF)</td><td>{{ fmt(currentCb.pf) }}</td></tr>
                <tr><td class="dv-label">고조파 (%)</td><td>{{ fmt(currentCb.pthd) }}</td></tr>
                <tr class="bg-violet-50/50 dark:bg-violet-900/10"><td class="dv-label">유효전력량 (kWh)</td><td class="font-semibold">{{ currentCb.kwh }}</td></tr>
                <tr class="bg-violet-50/50 dark:bg-violet-900/10"><td class="dv-label">무효전력량 (kvarh)</td><td>{{ currentCb.kvarh }}</td></tr>
                <tr class="bg-violet-50/50 dark:bg-violet-900/10"><td class="dv-label">피상전력량 (kVAh)</td><td>{{ currentCb.kVAh }}</td></tr>
                <tr class="bg-emerald-50/50 dark:bg-emerald-900/10"><td class="dv-label">현월 전력량 (kWh)</td><td class="font-semibold">{{ currentCb.kwh_tm }}</td></tr>
                <tr class="bg-emerald-50/50 dark:bg-emerald-900/10"><td class="dv-label">전월 전력량 (kWh)</td><td>{{ currentCb.kwh_lm }}</td></tr>
              </tbody>
            </table>
          </template>

          <div v-else-if="data" class="flex items-center justify-center h-32 text-gray-400 text-sm">해당 CB 데이터 없음</div>
        </div>

        <!-- Footer -->
        <div class="px-5 py-2 border-t border-gray-200 dark:border-gray-700 flex justify-end bg-gray-50 dark:bg-gray-800/50 flex-shrink-0">
          <button @click="$emit('close')" class="px-5 py-1.5 text-xs font-semibold border border-gray-300 dark:border-gray-600 rounded-md text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700">닫기</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import axios from 'axios'

const props = defineProps({
  open: { type: Boolean, default: false },
  canid: { type: String, default: '' },
  cbtype: { type: Number, default: 0 },
  cbcount: { type: Number, default: 1 },
})

defineEmits(['close'])

const loading = ref(false)
const error = ref('')
const data = ref(null)
const activeTab = ref(0)

const cbtypeLabel = computed(() => {
  const m = { 1: '1P', 4: '3P3W', 5: '3P4W', 6: '1P+Z' }
  return m[props.cbtype] || '—'
})

const currentCb = computed(() => {
  const list = data.value?.cblist
  if (!list) return null
  const cb = list[activeTab.value]
  return (cb && cb.cbtype !== undefined) ? cb : null
})

function is3p(cbtype) { return cbtype === 4 || cbtype === 5 }

function fmt(v) {
  if (v == null || v === undefined) return '—'
  if (typeof v === 'number') return Number.isInteger(v) ? v : v.toFixed(2)
  return v
}

function fmtTs(ts) {
  if (!ts) return '—'
  return new Date(ts * 1000).toLocaleString()
}

// 더미 CB 생성
function makeDummyCb3p(idx) {
  return {
    cbtype: 5, pth: 1.2 + idx * 0.3,
    vrms: [208.40 + idx, 208.60 + idx, 208.78 + idx],
    irms: [12.3 + idx, 11.8 + idx, 12.1 + idx], ig: 0.02,
    power: [7650 + idx * 100, 1230 + idx * 50, 7750 + idx * 100],
    pf: 0.98, kwh: 8299 + idx * 100, kvarh: 4228 + idx * 50, kVAh: 10079 + idx * 120,
    phwatt: [2550 + idx * 30, 2540 + idx * 30, 2560 + idx * 30],
    kwh_tm: 5174 + idx * 80, kwh_lm: 3124 + idx * 60,
    vunbal: 0.09, iunbal: 0.12, vpp: [361.14 + idx, 361.47 + idx, 361.29 + idx],
  }
}

function makeDummyCb1p(idx) {
  return {
    cbtype: 1, pthd: 0.8 + idx * 0.2,
    vrms: [209.53 + idx * 0.5], irms: [5.2 + idx * 0.3], ig: 0.01,
    power: [1080 + idx * 50, 320 + idx * 20, 1130 + idx * 55],
    pf: 0.96, kwh: 8322 + idx * 80, kvarh: 4249 + idx * 40, kVAh: 10110 + idx * 90,
    kwh_tm: 5179 + idx * 60, kwh_lm: 3143 + idx * 45,
  }
}

function makeDummy() {
  const is3 = props.cbtype === 4 || props.cbtype === 5
  const cblist = Array.from({ length: props.cbcount }, (_, i) => is3 ? makeDummyCb3p(i) : makeDummyCb1p(i))
  return {
    ts: Math.floor(Date.now() / 1000),
    canid: Number(props.canid),
    devtype: is3 ? 3 : 259,
    fwver: 60, status: 1,
    temp: 25.97 + Math.random() * 2,
    freq: 60.05,
    cblist,
  }
}

async function fetchData() {
  if (!props.canid) return
  loading.value = true
  error.value = ''
  data.value = null
  activeTab.value = 0

  try {
    const resp = await axios.post('/ibsm/read', {
      seq: 1, cmd: 'read_req', group: 0, canid: Number(props.canid),
    })
    if (resp.data?.values) {
      data.value = resp.data.values
    } else {
      throw new Error('Invalid response')
    }
  } catch (e) {
    console.warn('API 실패, 더미 데이터 사용:', e.message)
    data.value = makeDummy()
  } finally {
    loading.value = false
  }
}

watch(() => props.open, (v) => {
  if (v) fetchData()
})
</script>

<style scoped>
.dv-table {
  @apply w-full text-xs border-collapse;
}
.dv-table thead tr {
  @apply bg-gray-50 dark:bg-gray-700/60 border-b border-gray-200 dark:border-gray-600;
}
.dv-table th {
  @apply px-3 py-2 text-[10px] font-bold text-gray-400 dark:text-gray-500 uppercase tracking-wider text-center;
}
.dv-table th:first-child {
  @apply text-left;
}
.dv-table td {
  @apply px-3 py-1.5 text-center text-gray-700 dark:text-gray-200 tabular-nums;
}
.dv-table tbody tr {
  @apply border-b border-gray-100 dark:border-gray-700/30;
}
.dv-label {
  @apply text-left text-gray-500 dark:text-gray-400 font-medium;
}
</style>
