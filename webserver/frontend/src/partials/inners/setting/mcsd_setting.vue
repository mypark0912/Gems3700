<template>
  <div class="grow dark:text-white">
    <div class="p-6 space-y-6">

    <!-- Action Buttons -->
    <div class="flex items-center flex-wrap gap-2">
      <button @click="onScan" class="cmd-btn" :disabled="!props.isSetupMode">Scan</button>
      <button @click="onClear" class="cmd-btn" :disabled="!props.isSetupMode">Clear</button>
      <button @click="onAutoFill" class="cmd-btn" :disabled="!props.isSetupMode || !scanDone">Auto Fill</button>
      <div class="w-px h-5 bg-gray-200 dark:bg-gray-600 mx-1"></div>
      <button @click="rs485TestModalOpen = true" class="cmd-btn" :disabled="!props.isSetupMode">통신 테스트</button>
      <div class="flex-1"></div>
    </div>

    <!-- Table -->
    <div class="rounded-xl border border-gray-200 dark:border-gray-700/60 overflow-hidden">
      <div class="overflow-auto bg-white dark:bg-gray-800" style="max-height: calc(100vh - 260px);">
        <table class="line-table">
          <thead><tr>
            <th class="w-10">No.</th>
            <th class="w-[150px]">Serial Number</th>
            <th class="w-[140px]">TOU Name</th>
            <th class="w-[100px]">Modbus ID</th>
            <th class="w-[200px]">System Type</th>
            <th class="w-[170px]">CT1</th>
            <th class="w-[170px]">CT2</th>
            <th class="w-[170px]">CT3</th>
            <!-- <th class="w-[50px]">데이터</th> -->
            <!-- <th class="w-[50px]">보정</th> -->
          </tr></thead>
          <tbody>
            <tr v-for="(row, i) in rows" :key="i"
                :class="i % 2 === 0 ? '' : 'bg-gray-50/50 dark:bg-gray-700/10'">
              <!-- No -->
              <td class="text-center text-xs text-gray-500 dark:text-gray-400 font-medium">{{ i + 1 }}</td>
              <!-- Serial Number -->
              <td>
                <select v-model="row.serialNumber" @change="onSerialChange(i)" class="id-sel">
                  <option :value="null">Select</option>
                  <option v-for="dev in getAvailableDevices(i)" :key="dev.serialNumber" :value="dev.serialNumber">
                    {{ dev.serialNumber }}
                  </option>
                </select>
              </td>
              <!-- TOU Name -->
              <td>
                <input v-if="row.serialNumber" type="text" v-model="row.touName" maxlength="20"
                  class="tou-input" placeholder="TOU Name">
                <span v-else class="text-gray-300 dark:text-gray-600 text-xs">-</span>
              </td>
              <!-- Modbus ID (combobox, unique) -->
              <td class="text-center">
                <select v-if="row.serialNumber" v-model.number="row.modbusId" class="cb-sel">
                  <option :value="0">-</option>
                  <option v-for="mid in getAvailableModbusIds(i)" :key="mid" :value="mid">{{ mid }}</option>
                </select>
                <span v-else class="text-gray-300 dark:text-gray-600 text-xs">-</span>
              </td>
              <!-- System Type (category 내 2개 버튼그룹) -->
              <td>
                <div v-if="row.serialNumber" class="flex items-center justify-center gap-1">
                  <span v-for="opt in getSystemTypeOptions(row.systemType)" :key="opt.value"
                    class="chip"
                    :class="row.systemType === opt.value ? 'chip-active' : 'chip-inactive'"
                    @click="setSystemType(i, opt.value)">
                    {{ opt.label }}
                  </span>
                </div>
                <span v-else class="text-gray-300 dark:text-gray-600 text-xs">-</span>
              </td>
              <!-- CT1 (Normal / Reverse button group) -->
              <td>
                <div v-if="row.serialNumber" class="flex items-center justify-center gap-1">
                  <span class="chip" :class="row.ct1 === 1 ? 'chip-active' : 'chip-inactive'"
                    @click="setCt(i, 'ct1', 1)">Normal</span>
                  <span class="chip" :class="row.ct1 === -1 ? 'chip-active' : 'chip-inactive'"
                    @click="setCt(i, 'ct1', -1)">Reverse</span>
                </div>
                <span v-else class="text-center text-gray-300 dark:text-gray-600 text-xs block">-</span>
              </td>
              <!-- CT2 -->
              <td>
                <div v-if="row.serialNumber" class="flex items-center justify-center gap-1">
                  <span class="chip" :class="row.ct2 === 1 ? 'chip-active' : 'chip-inactive'"
                    @click="setCt(i, 'ct2', 1)">Normal</span>
                  <span class="chip" :class="row.ct2 === -1 ? 'chip-active' : 'chip-inactive'"
                    @click="setCt(i, 'ct2', -1)">Reverse</span>
                </div>
                <span v-else class="text-center text-gray-300 dark:text-gray-600 text-xs block">-</span>
              </td>
              <!-- CT3 -->
              <td>
                <div v-if="row.serialNumber" class="flex items-center justify-center gap-1">
                  <span class="chip" :class="row.ct3 === 1 ? 'chip-active' : 'chip-inactive'"
                    @click="setCt(i, 'ct3', 1)">Normal</span>
                  <span class="chip" :class="row.ct3 === -1 ? 'chip-active' : 'chip-inactive'"
                    @click="setCt(i, 'ct3', -1)">Reverse</span>
                </div>
                <span v-else class="text-center text-gray-300 dark:text-gray-600 text-xs block">-</span>
              </td>
              <!-- 데이터 확인 button -->
              <!-- <td class="text-center">
                <button v-if="row.serialNumber" class="act-btn act-sky" :disabled="!props.isSetupMode" @click="openDataModal(i)" title="데이터 확인"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M3 3v18h18"/><path d="M18 17V9"/><path d="M13 17V5"/><path d="M8 17v-3"/></svg></button>
              </td> -->
              <!-- 보정 button -->
              <!-- <td class="text-center">
                <button v-if="row.serialNumber" class="act-btn act-pink" :disabled="!props.isSetupMode" @click="openWattCorrModal(i)" title="전력량 보정"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M21.174 6.812a1 1 0 0 0-3.986-3.987L3.842 16.174a2 2 0 0 0-.5.83l-1.321 4.352a.5.5 0 0 0 .623.622l4.353-1.32a2 2 0 0 0 .83-.497z"/></svg></button>
              </td> -->
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 데이터 확인 Modal -->
    <DataViewModal
      :open="canTestModalOpen"
      :canid="rows[dataModalIdx]?.serialNumber || ''"
      :cbtype="0"
      :cbcount="1"
      @close="canTestModalOpen = false"
    />

    <!-- Wattage Correction Modal -->
    <WattCorrModal
      :open="wattCorrModalOpen"
      :canid="rows[wattCorrIdx]?.serialNumber || ''"
      :touName="rows[wattCorrIdx]?.touName || ''"
      :cbtype="0"
      :cbcount="1"
      @close="wattCorrModalOpen = false"
    />

    <!-- ══════════════════════════════════════════
         Modal – RS-485 Communication Test
    ══════════════════════════════════════════ -->
    <Teleport to="body">
      <div v-if="rs485TestModalOpen" class="fixed inset-0 bg-black/40 backdrop-blur-sm z-40" @click="rs485TestModalOpen = false"></div>
    </Teleport>
    <Teleport to="body">
      <div v-if="rs485TestModalOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="rs485TestModalOpen = false">
        <div class="bg-white dark:bg-gray-800 rounded-xl shadow-2xl border border-gray-200 dark:border-gray-700 w-[90vw] max-w-[1100px] flex flex-col" @click.stop>
          <!-- Header -->
          <div class="px-5 py-3 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between flex-shrink-0">
            <h3 class="text-[15px] font-extrabold text-gray-800 dark:text-gray-100">RS-485 Communication Test</h3>
            <button @click="rs485TestModalOpen = false" class="w-7 h-7 rounded-full bg-gray-100 dark:bg-gray-700 text-gray-500 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-gray-600 flex items-center justify-center text-sm flex-shrink-0">✕</button>
          </div>
          <!-- Toolbar -->
          <div class="px-5 py-2.5 border-b border-gray-100 dark:border-gray-700/50 flex items-center justify-between flex-shrink-0">
            <div class="flex items-center gap-2">
              <button @click="runCommTest" class="cmd-btn" :disabled="commTestRunning || commSelectedCount === 0">
                {{ commTestRunning ? 'Testing...' : `Start (${commSelectedCount})` }}
              </button>
              <div class="w-px h-5 bg-gray-200 dark:bg-gray-600 mx-1"></div>
              <button @click="commSelectAll" class="cmd-btn text-[10px]">전체 선택</button>
              <button @click="commDeselectAll" class="cmd-btn text-[10px]">선택 해제</button>
            </div>
            <div class="flex items-center gap-2">
              <span class="rs485-legend rs485-waiting">WAITING</span>
              <span class="rs485-legend rs485-online">ONLINE</span>
              <span class="rs485-legend rs485-error">COMM ERROR</span>
              <span class="rs485-legend rs485-nodev">NO DEVICE</span>
            </div>
          </div>
          <!-- Grid -->
          <div class="p-5 overflow-y-auto" style="max-height: calc(80vh - 140px);">
            <div class="grid grid-cols-8 gap-2.5">
              <div v-for="(slot, idx) in commSlots" :key="idx"
                class="rs485-cell cursor-pointer"
                :class="[slot.status, { 'ring-2 ring-violet-400 ring-offset-1 dark:ring-offset-gray-800': slot.selected }]"
                @click="toggleCommSelect(idx)">
                <div class="text-xs font-bold" :class="slot.serialNumber ? 'text-gray-700 dark:text-gray-200' : 'text-gray-300 dark:text-gray-600'">
                  {{ idx + 1 }}
                </div>
                <div v-if="slot.serialNumber" class="text-[10px] text-gray-400 dark:text-gray-500 mt-0.5">
                  {{ slot.serialNumber }}
                </div>
              </div>
            </div>
          </div>
          <!-- Footer -->
          <div class="px-5 py-3 border-t border-gray-200 dark:border-gray-700 flex justify-end bg-gray-50 dark:bg-gray-800/50 flex-shrink-0">
            <button @click="rs485TestModalOpen = false" class="px-5 py-2 text-xs font-semibold border border-gray-300 dark:border-gray-600 rounded-md text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700">닫기</button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Scan Result Modal -->
    <Teleport to="body">
      <div v-if="showScanResult" class="fixed inset-0 bg-black/30 z-40"></div>
    </Teleport>
    <Teleport to="body">
      <div v-if="showScanResult" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="bg-white dark:bg-gray-800 rounded-xl shadow-2xl border border-gray-200 dark:border-gray-700 w-full max-w-[500px]" @click.stop>
          <div class="px-5 py-4 border-b border-gray-200 dark:border-gray-700">
            <h3 class="text-base font-bold text-gray-800 dark:text-gray-100">Scan Result</h3>
          </div>
          <div class="p-5 max-h-[400px] overflow-y-auto">
            <div v-if="scanLoading" class="flex items-center justify-center py-8">
              <svg class="animate-spin h-6 w-6 text-violet-500 mr-3" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none"/>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
              </svg>
              <span class="text-sm text-gray-500 dark:text-gray-400">Scanning...</span>
            </div>
            <template v-else>
              <div class="mb-4">
                <h4 class="text-sm font-semibold text-gray-500 dark:text-gray-400 mb-2">기존 등록 장비 ({{ existingSerials.length }})</h4>
                <div v-if="existingSerials.length" class="flex flex-wrap gap-1.5">
                  <span v-for="sn in existingSerials" :key="'e-'+sn"
                    class="px-2 py-0.5 text-xs rounded bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300">
                    {{ sn }}
                  </span>
                </div>
                <p v-else class="text-xs text-gray-400">없음</p>
              </div>
              <div class="mb-5">
                <h4 class="text-sm font-semibold text-violet-500 dark:text-violet-400 mb-2">스캔 장비 ({{ pendingScanDevices.length }})</h4>
                <div v-if="pendingScanDevices.length" class="flex flex-wrap gap-1.5">
                  <span v-for="dev in pendingScanDevices" :key="'s-'+dev.serialNumber"
                    class="px-2 py-0.5 text-xs rounded"
                    :class="existingSerials.includes(dev.serialNumber) ? 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300' : 'bg-violet-100 dark:bg-violet-900/40 text-violet-700 dark:text-violet-300'">
                    {{ dev.serialNumber }}
                  </span>
                </div>
                <p v-else class="text-xs text-gray-400">장비를 찾지 못했습니다</p>
              </div>
              <div class="flex items-center justify-end gap-3">
                <button @click="applyScanResult" :disabled="!pendingScanDevices.length" class="h-9 px-5 rounded-md text-sm font-semibold bg-violet-500 hover:bg-violet-600 text-white transition-all disabled:opacity-40">적용</button>
                <button @click="showScanResult = false" class="h-9 px-5 rounded-md text-sm font-semibold border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 transition-all">취소</button>
              </div>
            </template>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Confirm Dialog -->
    <Teleport to="body">
      <div v-if="showConfirmDialog" class="fixed inset-0 bg-black/30 z-40" @click="showConfirmDialog = false"></div>
    </Teleport>
    <Teleport to="body">
      <div v-if="showConfirmDialog" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="bg-white dark:bg-gray-800 rounded-xl shadow-2xl border border-gray-200 dark:border-gray-700 w-full max-w-[420px]" @click.stop>
          <div class="px-5 py-4 border-b border-gray-200 dark:border-gray-700">
            <h3 class="text-base font-bold text-gray-800 dark:text-gray-100">{{ confirmTitle }}</h3>
          </div>
          <div class="p-5">
            <p class="text-sm text-gray-600 dark:text-gray-300 mb-5 whitespace-pre-line">{{ confirmMsg }}</p>
            <div class="flex items-center justify-end gap-3">
              <button @click="confirmAction" class="h-9 px-5 rounded-md text-sm font-semibold bg-violet-500 hover:bg-violet-600 text-white transition-all">Yes</button>
              <button @click="showConfirmDialog = false" class="h-9 px-5 rounded-md text-sm font-semibold border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 transition-all">No</button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch, inject, onUnmounted } from 'vue'
import DataViewModal from './DataViewModal.vue'
import WattCorrModal from './WattCorrModal.vue'

const props = defineProps({
  isSetupMode: { type: Boolean, default: false }
})

const setupDict = inject('setupDict')

// ── WebSocket (RS-485 command) ──
let cmdWs = null
let shouldReconnect = false
const wsReady = ref(false)

function connectCmdWs() {
  if (cmdWs) return
  shouldReconnect = true
  const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
  cmdWs = new WebSocket(`${protocol}//${location.host}/setting/ws`)

  cmdWs.onopen = () => { wsReady.value = true }
  cmdWs.onmessage = (e) => {
    const data = JSON.parse(e.data)
    if (data.error) { console.error('ws error:', data.error); return }
    if (data.cmd === 'ping_resp') handlePingResp(data)
  }
  cmdWs.onclose = () => {
    cmdWs = null
    wsReady.value = false
    if (shouldReconnect) setTimeout(connectCmdWs, 1000)
  }
}

function disconnectCmdWs() {
  shouldReconnect = false
  if (cmdWs) {
    cmdWs.onclose = null
    cmdWs.close()
    cmdWs = null
    wsReady.value = false
  }
}

function sendCmd(cmd, extra = {}) {
  if (!cmdWs || cmdWs.readyState !== WebSocket.OPEN) {
    console.error('ws not connected')
    return
  }
  cmdWs.send(JSON.stringify({ cmd, ...extra }))
}

watch(() => props.isSetupMode, (remote) => {
  if (remote) connectCmdWs()
  else disconnectCmdWs()
}, { immediate: true })

onUnmounted(() => { disconnectCmdWs() })

const MAX_ROWS = 32

// ── System Type 정의 ──
// 1: 1P2W, 2: 3P3W, 3: 3P4W, 4: 1P3W
const systemTypeLabel = {
  0: '',
  1: '1P2W',
  2: '3P3W',
  3: '3P4W',
  4: '1P3W',
}

// systemType 카테고리별 옵션 (category 내에서 사용자가 선택)
// 1P 카테고리(scan이 1 또는 4): 1P2W(1), 1P3W(4)
// 3P 카테고리(scan이 2 또는 3): 3P3W(2), 3P4W(3)
function getSystemTypeOptions(st) {
  if (st === 1 || st === 4) return [
    { value: 1, label: '1P2W' },
    { value: 4, label: '1P3W' },
  ]
  if (st === 2 || st === 3) return [
    { value: 2, label: '3P3W' },
    { value: 3, label: '3P4W' },
  ]
  return []
}

// ── 더미 데이터 ──
const DUMMY_DEVICES = [
  { serialNumber: 'MCSD-0001', systemType: 1, modbusId: 1,  ct1: 1,  ct2: 1,  ct3: 1  },
  { serialNumber: 'MCSD-0002', systemType: 2, modbusId: 2,  ct1: 1,  ct2: -1, ct3: 1  },
  { serialNumber: 'MCSD-0003', systemType: 3, modbusId: 3,  ct1: 1,  ct2: 1,  ct3: -1 },
  { serialNumber: 'MCSD-0004', systemType: 4, modbusId: 4,  ct1: -1, ct2: 1,  ct3: 1  },
  { serialNumber: 'MCSD-0005', systemType: 2, modbusId: 5,  ct1: 1,  ct2: 1,  ct3: 1  },
  { serialNumber: 'MCSD-0006', systemType: 3, modbusId: 6,  ct1: -1, ct2: -1, ct3: -1 },
  { serialNumber: 'MCSD-0007', systemType: 1, modbusId: 7,  ct1: 1,  ct2: 1,  ct3: 1  },
  { serialNumber: 'MCSD-0008', systemType: 4, modbusId: 8,  ct1: 1,  ct2: -1, ct3: 1  },
  { serialNumber: 'MCSD-0009', systemType: 3, modbusId: 9,  ct1: 1,  ct2: 1,  ct3: 1  },
  { serialNumber: 'MCSD-0010', systemType: 2, modbusId: 10, ct1: -1, ct2: 1,  ct3: -1 },
  { serialNumber: 'MCSD-0011', systemType: 1, modbusId: 11, ct1: 1,  ct2: 1,  ct3: 1  },
  { serialNumber: 'MCSD-0012', systemType: 3, modbusId: 12, ct1: 1,  ct2: -1, ct3: 1  },
]

const scannedDevices = ref([])
const scanDone = ref(false)

function createEmptyRow() {
  return {
    serialNumber: null,
    touName: '',
    systemType: 0,
    modbusId: 0,
    ct1: 1,
    ct2: 1,
    ct3: 1,
  }
}

const MODBUS_ID_MIN = 1
const MODBUS_ID_MAX = 247

function getUsedModbusIds() {
  return rows.value.filter(r => r.serialNumber && r.modbusId > 0).map(r => r.modbusId)
}

function getAvailableModbusIds(rowIdx) {
  const used = getUsedModbusIds()
  const current = rows.value[rowIdx].modbusId
  const result = []
  for (let id = MODBUS_ID_MIN; id <= MODBUS_ID_MAX; id++) {
    if (id === current || !used.includes(id)) result.push(id)
  }
  return result
}

const rows = ref(Array.from({ length: MAX_ROWS }, () => reactive(createEmptyRow())))

// ── setupDict.mcs ↔ rows 양방향 동기화 ──
let isSyncing = false

function loadFromMcsSet(data) {
  const feeders = data?.feeders || []
  for (let i = 0; i < MAX_ROWS; i++) {
    const f = feeders[i]
    if (f && f.serialNumber) {
      // 스캔 목록에도 반영
      if (!scannedDevices.value.find(d => d.serialNumber === f.serialNumber)) {
        scannedDevices.value.push({
          serialNumber: f.serialNumber,
          systemType: f.systemType || 0,
          modbusId: f.modbusId || 0,
          ct1: f.ct1 ?? 1,
          ct2: f.ct2 ?? 1,
          ct3: f.ct3 ?? 1,
        })
      }
      Object.assign(rows.value[i], {
        serialNumber: f.serialNumber,
        touName: f.touName || '',
        systemType: f.systemType || 0,
        modbusId: f.modbusId || 0,
        ct1: f.ct1 ?? 1,
        ct2: f.ct2 ?? 1,
        ct3: f.ct3 ?? 1,
      })
    } else {
      Object.assign(rows.value[i], createEmptyRow())
    }
  }
}

function syncToMcsSet() {
  if (!setupDict?.value) return
  isSyncing = true
  const feeders = rows.value.map(r => ({
    serialNumber: r.serialNumber,
    touName: r.touName || '',
    systemType: r.systemType || 0,
    modbusId: r.modbusId || 0,
    ct1: r.ct1 ?? 1,
    ct2: r.ct2 ?? 1,
    ct3: r.ct3 ?? 1,
  }))
  setupDict.value.mcs = { channel: 'mcs', Enable: 1, feeders }
  setTimeout(() => { isSyncing = false }, 0)
}

// setupDict.mcs 외부 변경 감지 → rows 로드
watch(() => setupDict?.value?.mcs, (newVal) => {
  if (isSyncing) return
  if (newVal?.feeders?.length) loadFromMcsSet(newVal)
}, { immediate: true, deep: true })

// rows 변경 → setupDict.mcs 반영
watch(rows, () => {
  if (isSyncing) return
  syncToMcsSet()
}, { deep: true })

function getUsedSerials() {
  return rows.value.filter(r => r.serialNumber).map(r => r.serialNumber)
}

function getAvailableDevices(rowIdx) {
  const used = getUsedSerials()
  const current = rows.value[rowIdx].serialNumber
  return scannedDevices.value.filter(dev => dev.serialNumber === current || !used.includes(dev.serialNumber))
}

function onSerialChange(idx) {
  const row = rows.value[idx]
  if (!row.serialNumber) {
    Object.assign(row, createEmptyRow())
    return
  }
  const dev = scannedDevices.value.find(d => d.serialNumber === row.serialNumber)
  if (!dev) return
  // 중복되지 않는 modbusId 자동 할당 (이미 사용 중이면 0으로 두고 사용자가 선택)
  const used = getUsedModbusIds().filter(m => m !== row.modbusId)
  const autoModbus = used.includes(dev.modbusId) ? 0 : dev.modbusId
  Object.assign(row, {
    touName: row.touName || `Feeder #${idx + 1}`,
    systemType: dev.systemType,
    modbusId: autoModbus,
    ct1: dev.ct1,
    ct2: dev.ct2,
    ct3: dev.ct3,
  })
}

function setSystemType(idx, value) {
  rows.value[idx].systemType = value
}

function setCt(idx, key, value) {
  rows.value[idx][key] = value
}

// ── Data View / Wattage Correction Modals ──
const canTestModalOpen = ref(false)
const wattCorrModalOpen = ref(false)
const dataModalIdx = ref(0)
const wattCorrIdx = ref(0)

function openDataModal(idx) {
  dataModalIdx.value = idx
  canTestModalOpen.value = true
}

function openWattCorrModal(idx) {
  wattCorrIdx.value = idx
  wattCorrModalOpen.value = true
}

// ── Tool Modals ──
const rs485TestModalOpen = ref(false)

// ── RS-485 Communication Test ──
const commTestRunning = ref(false)
const commTestResults = ref({})
const commSelected = ref({})

const commSelectedCount = computed(() => {
  return Object.keys(commSelected.value).filter(k => commSelected.value[k] && rows.value[k]?.serialNumber).length
})

const commSlots = computed(() => {
  return Array.from({ length: MAX_ROWS }, (_, i) => {
    const row = rows.value[i]
    const result = commTestResults.value[i]
    let status = 'rs485-nodev'
    if (row.serialNumber) {
      status = result ? `rs485-${result}` : 'rs485-waiting'
    }
    return {
      serialNumber: row.serialNumber || null,
      status,
      selected: !!commSelected.value[i] && !!row.serialNumber,
    }
  })
})

function toggleCommSelect(idx) {
  if (!rows.value[idx]?.serialNumber) return
  commSelected.value = { ...commSelected.value, [idx]: !commSelected.value[idx] }
}

function commSelectAll() {
  const sel = {}
  rows.value.forEach((r, i) => { if (r.serialNumber) sel[i] = true })
  commSelected.value = sel
}

function commDeselectAll() {
  commSelected.value = {}
}

const commPingTargets = ref([])

function runCommTest() {
  commTestRunning.value = true
  commTestResults.value = {}
  const targets = rows.value
    .map((r, i) => ({ idx: i, serialNumber: r.serialNumber, modbusId: r.modbusId }))
    .filter(r => r.serialNumber && commSelected.value[r.idx])

  targets.forEach(t => {
    commTestResults.value = { ...commTestResults.value, [t.idx]: 'waiting' }
  })

  commPingTargets.value = targets
  targets.forEach(t => {
    sendCmd('ping_req', { modbusId: Number(t.modbusId), serialNumber: t.serialNumber })
  })
}

function handlePingResp(data) {
  const sn = data.serialNumber
  const target = commPingTargets.value.find(t => rows.value[t.idx]?.serialNumber === sn)
  if (!target) return
  commTestResults.value = {
    ...commTestResults.value,
    [target.idx]: data.res === 'ok' ? 'online' : 'error'
  }
  const allDone = commPingTargets.value.every(t =>
    commTestResults.value[t.idx] === 'online' || commTestResults.value[t.idx] === 'error'
  )
  if (allDone) commTestRunning.value = false
}

// ── Confirm Dialog ──
const showConfirmDialog = ref(false)
const confirmTitle = ref('')
const confirmMsg = ref('')
let confirmCb = null

function openConfirm(title, msg, cb) {
  confirmTitle.value = title
  confirmMsg.value = msg
  confirmCb = cb
  showConfirmDialog.value = true
}
function confirmAction() {
  showConfirmDialog.value = false
  if (confirmCb) confirmCb()
}

// ── Scan Result Modal ──
const showScanResult = ref(false)
const scanLoading = ref(false)
const pendingScanDevices = ref([])

const existingSerials = computed(() =>
  scannedDevices.value.map(d => d.serialNumber)
)

function applyScanResult() {
  scannedDevices.value = [...pendingScanDevices.value]
  scanDone.value = true
  showScanResult.value = false
}

async function onScan() {
  pendingScanDevices.value = []
  scanLoading.value = true
  showScanResult.value = true
  // 더미 데이터 기반 스캔 시뮬레이션 (800ms 지연)
  await new Promise(r => setTimeout(r, 800))
  pendingScanDevices.value = DUMMY_DEVICES.map(d => ({ ...d }))
  scanLoading.value = false
}

function onClear() {
  openConfirm('Clear', 'Are you sure you want to clear the device list?', () => {
    for (let i = 0; i < MAX_ROWS; i++) {
      Object.assign(rows.value[i], createEmptyRow())
    }
  })
}

function onAutoFill() {
  if (!scannedDevices.value.length) return
  const sorted = [...scannedDevices.value].sort((a, b) =>
    String(a.modbusId).localeCompare(String(b.modbusId), undefined, { numeric: true })
  )
  for (const dev of sorted) {
    if (getUsedSerials().includes(dev.serialNumber)) continue
    const emptyIdx = rows.value.findIndex(r => !r.serialNumber)
    if (emptyIdx < 0) break
    const emptyRow = rows.value[emptyIdx]
    Object.assign(emptyRow, {
      serialNumber: dev.serialNumber,
      touName: `Feeder #${emptyIdx + 1}`,
      systemType: dev.systemType,
      modbusId: emptyIdx + 1,
      ct1: dev.ct1,
      ct2: dev.ct2,
      ct3: dev.ct3,
    })
  }
}
</script>

<style scoped>
.cmd-btn {
  @apply px-3 py-1.5 text-xs font-semibold rounded-md border border-gray-300 dark:border-gray-600;
  @apply bg-white dark:bg-gray-700 text-gray-600 dark:text-gray-200;
  @apply hover:bg-gray-50 dark:hover:bg-gray-600 transition-colors;
  @apply disabled:opacity-40 disabled:cursor-not-allowed;
}

.line-table { @apply w-full text-[13px]; }
.line-table thead { @apply sticky top-0 z-10; }
.line-table thead tr { @apply bg-gray-50 dark:bg-gray-700/80 border-b border-gray-200 dark:border-gray-600; }
.line-table th { @apply px-2 py-2 text-[11px] font-bold text-gray-400 dark:text-gray-500 uppercase tracking-wider text-center whitespace-nowrap; }
.line-table td { @apply px-2 py-[7px] whitespace-nowrap; }
.line-table tbody tr { @apply border-b border-gray-100 dark:border-gray-700/30; }

.id-sel {
  @apply w-full text-xs px-1.5 py-1 rounded border border-gray-200 dark:border-gray-600;
  @apply bg-transparent text-gray-700 dark:text-gray-200;
  @apply focus:outline-none focus:ring-1 focus:ring-violet-400;
}

.cb-sel {
  @apply text-xs px-2 py-1 min-w-[60px] rounded border border-gray-200 dark:border-gray-600;
  @apply bg-transparent text-gray-700 dark:text-gray-200;
  @apply focus:outline-none focus:ring-1 focus:ring-violet-400;
}

.tou-input {
  @apply w-full text-xs px-1.5 py-1 rounded border border-gray-200 dark:border-gray-600;
  @apply bg-transparent text-gray-700 dark:text-gray-200;
  @apply focus:outline-none focus:ring-1 focus:ring-violet-400;
}

.type-badge {
  @apply inline-block px-2 py-0.5 rounded text-[10px] font-semibold;
  @apply bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300;
}

/* Action buttons (데이터/보정) */
.act-btn {
  @apply w-8 h-8 rounded-lg flex items-center justify-center mx-auto transition-all;
  @apply disabled:bg-gray-100 dark:disabled:bg-gray-800 disabled:text-gray-300 dark:disabled:text-gray-600;
  @apply disabled:cursor-not-allowed disabled:shadow-none;
}
.act-sky {
  @apply bg-sky-500 text-white shadow-sm shadow-sky-200 dark:shadow-sky-900/40;
  @apply hover:bg-sky-600 hover:shadow-md hover:shadow-sky-300 dark:hover:shadow-sky-800/40;
}
.act-pink {
  @apply bg-pink-500 text-white shadow-sm shadow-pink-200 dark:shadow-pink-900/40;
  @apply hover:bg-pink-600 hover:shadow-md hover:shadow-pink-300 dark:hover:shadow-pink-800/40;
}

/* Chips (button group) */
.chip {
  @apply inline-flex items-center justify-center min-w-[36px] px-2 py-1 rounded text-[11px] font-semibold cursor-pointer transition-all;
}
.chip-active {
  @apply bg-violet-500 text-white;
}
.chip-inactive {
  @apply bg-gray-100 dark:bg-gray-700 text-gray-400 dark:text-gray-500;
  @apply hover:bg-gray-200 dark:hover:bg-gray-600 hover:text-gray-600 dark:hover:text-gray-300;
}

/* 통신 테스트 grid */
.rs485-cell {
  @apply rounded-lg border-2 px-3 py-3 text-center min-h-[56px] flex flex-col items-center justify-center transition-all;
}
.rs485-nodev {
  @apply border-dashed border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800/50;
}
.rs485-waiting {
  @apply border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-700/30;
}
.rs485-online {
  @apply border-emerald-400 bg-emerald-50 dark:bg-emerald-900/20 dark:border-emerald-600;
}
.rs485-error {
  @apply border-red-400 bg-red-50 dark:bg-red-900/20 dark:border-red-600;
}

.rs485-legend {
  @apply px-3 py-1 text-[11px] font-bold rounded-md border;
}
.rs485-legend.rs485-waiting {
  @apply border-gray-300 dark:border-gray-600 text-gray-500 dark:text-gray-400 bg-white dark:bg-gray-800;
}
.rs485-legend.rs485-online {
  @apply border-emerald-500 text-white bg-emerald-500;
}
.rs485-legend.rs485-error {
  @apply border-red-500 text-white bg-red-500;
}
.rs485-legend.rs485-nodev {
  @apply border-gray-300 dark:border-gray-600 text-gray-400 dark:text-gray-500 bg-white dark:bg-gray-800;
}
</style>
