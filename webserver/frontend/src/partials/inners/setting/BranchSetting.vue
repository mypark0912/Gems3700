<template>
  <div class="grow dark:text-white">
    <div class="p-6 space-y-6">

    <!-- Action Buttons -->
    <div class="flex items-center flex-wrap gap-2">
      <button @click="onScan" class="cmd-btn" :disabled="!props.isSetupMode">Scan</button>
      <button @click="onClear" class="cmd-btn" :disabled="!props.isSetupMode">Clear</button>
      <button @click="onAutoFill" class="cmd-btn" :disabled="!props.isSetupMode || !scanDone">Auto Fill</button>
      <div class="w-px h-5 bg-gray-200 dark:bg-gray-600 mx-1"></div>
      <button @click="batchModalOpen = true" class="cmd-btn" :disabled="!props.isSetupMode || !hasRegistered">일괄 설정</button>
      <div class="w-px h-5 bg-gray-200 dark:bg-gray-600 mx-1"></div>
      <button @click="rs485TestModalOpen = true" class="cmd-btn" :disabled="!props.isSetupMode">통신 테스트</button>
      <div class="flex-1"></div>
      <button @click="fwUpgradeModalOpen = true" class="cmd-btn" :disabled="!props.isSetupMode">FW Upgrade</button>
    </div>

    <!-- Table -->
    <div class="rounded-xl border border-gray-200 dark:border-gray-700/60 overflow-hidden">
      <div class="overflow-auto bg-white dark:bg-gray-800" style="max-height: calc(100vh - 260px);">
        <table class="line-table">
          <thead><tr>
            <th class="w-10">No.</th>
            <th class="w-[130px]">iBSM ID</th>
            <th class="w-[160px]">TOU Name</th>
            <th class="w-20">Type</th>
            <th class="w-20">CB Count</th>
            <th>CB Setting</th>
            <th class="w-[50px]">설정</th>
            <th class="w-[50px]">데이터</th>
            <th class="w-[50px]">보정</th>
          </tr></thead>
          <tbody>
            <tr v-for="(row, i) in rows" :key="i"
                :class="i % 2 === 0 ? '' : 'bg-gray-50/50 dark:bg-gray-700/10'">
              <!-- No -->
              <td class="text-center text-xs text-gray-500 dark:text-gray-400 font-medium">{{ i + 1 }}</td>
              <!-- iBSM ID (combo) -->
              <td>
                <select v-model="row.id" @change="onIdChange(i)" class="id-sel">
                  <option :value="null">Select</option>
                  <option v-for="dev in getAvailableDevices(i)" :key="dev.canid" :value="String(dev.canid)">
                    {{ dev.canid }}
                  </option>
                </select>
              </td>
              <!-- TOU Name -->
              <td>
                <input v-if="row.id" type="text" v-model="row.name" maxlength="20"
                  @input="onNameChange(i)"
                  class="tou-input" placeholder="TOU Name">
                <span v-else class="text-gray-300 dark:text-gray-600 text-xs">-</span>
              </td>
              <!-- Type -->
              <td class="text-center">
                <span v-if="row.id" class="type-badge">{{ row.type }}</span>
                <span v-else class="text-gray-300 dark:text-gray-600 text-xs">-</span>
              </td>
              <!-- CB Count -->
              <td class="text-center">
                <select v-if="row.id" v-model="row.cbCount" @change="onCbCountChange(i)" class="cb-sel">
                  <option v-for="n in row.maxCb" :key="n" :value="n">{{ n }}</option>
                </select>
                <span v-else class="text-gray-300 dark:text-gray-600 text-xs">-</span>
              </td>
              <!-- CB Setting (cbtype로 결정된 stype 옵션 — 혼용 없음) -->
              <td>
                <div v-if="row.id && row.cbCount > 0" class="flex flex-wrap items-center gap-2 py-0.5">
                  <!-- Per-CB stype 선택 -->
                  <template v-for="ci in row.cbCount" :key="ci - 1">
                    <span v-if="ci > 1" class="text-gray-300 dark:text-gray-600 mx-0.5">│</span>
                    <div class="flex items-center gap-0.5">
                      <span class="chip chip-label">CB{{ ci }}</span>
                      <span v-for="opt in getStypeOptions(row.cbtype)" :key="opt.value"
                        class="chip"
                        :class="row.stype[ci - 1] === opt.value ? 'chip-active' : 'chip-inactive'"
                        @click="setStype(i, ci - 1, opt.value)">
                        {{ opt.label }}
                      </span>
                    </div>
                  </template>
                </div>
                <span v-else class="text-gray-300 dark:text-gray-600 text-xs">-</span>
              </td>
              <!-- 설정 button -->
              <td class="text-center">
                <button v-if="row.id" class="act-btn act-violet" :disabled="!props.isSetupMode" @click="openSingleModal(i)" title="설정"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"/><circle cx="12" cy="12" r="3"/></svg></button>
              </td>
              <!-- 데이터 확인 button -->
              <td class="text-center">
                <button v-if="row.id" class="act-btn act-sky" :disabled="!props.isSetupMode" @click="openDataModal(i)" title="데이터 확인"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M3 3v18h18"/><path d="M18 17V9"/><path d="M13 17V5"/><path d="M8 17v-3"/></svg></button>
              </td>
              <!-- 보정 button -->
              <td class="text-center">
                <button v-if="row.id" class="act-btn act-pink" :disabled="!props.isSetupMode" @click="openWattCorrModal(i)" title="전력량 보정"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M21.174 6.812a1 1 0 0 0-3.986-3.987L3.842 16.174a2 2 0 0 0-.5.83l-1.321 4.352a.5.5 0 0 0 .623.622l4.353-1.32a2 2 0 0 0 .83-.497z"/></svg></button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 단독 설정 Modal -->
    <SingleSettingModal
      :open="singleModalOpen"
      :rowIndex="singleIdx"
      :canid="rows[singleIdx]?.id || ''"
      :touName="rows[singleIdx]?.name || ''"
      @close="singleModalOpen = false"
    />

    <!-- 배치 설정 Modal -->
    <BatchSettingModal
      :open="batchModalOpen"
      :registeredRows="registeredRows"
      @close="batchModalOpen = false"
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
                <div class="text-xs font-bold" :class="slot.id ? 'text-gray-700 dark:text-gray-200' : 'text-gray-300 dark:text-gray-600'">
                  {{ idx + 1 }}
                </div>
                <div v-if="slot.id" class="text-[10px] text-gray-400 dark:text-gray-500 mt-0.5">
                  {{ slot.id }}
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

    <!-- 데이터 확인 Modal -->
    <DataViewModal
      :open="canTestModalOpen"
      :canid="rows[dataModalIdx]?.id || ''"
      :cbtype="rows[dataModalIdx]?.cbtype || 0"
      :cbcount="rows[dataModalIdx]?.cbCount || 1"
      @close="canTestModalOpen = false"
    />

    <!-- FW Upgrade Modal -->
    <FwUpgradeModal
      :open="fwUpgradeModalOpen"
      :registeredRows="registeredRows"
      @close="fwUpgradeModalOpen = false"
    />

    <!-- Wattage Correction Modal -->
    <WattCorrModal
      :open="wattCorrModalOpen"
      :canid="rows[wattCorrIdx]?.id || ''"
      :touName="rows[wattCorrIdx]?.name || ''"
      :cbtype="rows[wattCorrIdx]?.cbtype || 0"
      :cbcount="rows[wattCorrIdx]?.cbCount || 1"
      @close="wattCorrModalOpen = false"
    />

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
            <!-- 스캔 중 -->
            <div v-if="scanLoading" class="flex items-center justify-center py-8">
              <svg class="animate-spin h-6 w-6 text-violet-500 mr-3" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none"/>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
              </svg>
              <span class="text-sm text-gray-500 dark:text-gray-400">Scanning...</span>
            </div>
            <!-- 스캔 완료 -->
            <template v-else>
              <!-- 기존 장비 -->
              <div class="mb-4">
                <h4 class="text-sm font-semibold text-gray-500 dark:text-gray-400 mb-2">기존 등록 장비 ({{ existingDeviceIds.length }})</h4>
                <div v-if="existingDeviceIds.length" class="flex flex-wrap gap-1.5">
                  <span v-for="id in existingDeviceIds" :key="'e-'+id"
                    class="px-2 py-0.5 text-xs rounded bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300">
                    {{ id }}
                  </span>
                </div>
                <p v-else class="text-xs text-gray-400">없음</p>
              </div>
              <!-- 스캔 장비 -->
              <div class="mb-5">
                <h4 class="text-sm font-semibold text-violet-500 dark:text-violet-400 mb-2">스캔 장비 ({{ pendingScanDevices.length }})</h4>
                <div v-if="pendingScanDevices.length" class="flex flex-wrap gap-1.5">
                  <span v-for="dev in pendingScanDevices" :key="'s-'+dev.canid"
                    class="px-2 py-0.5 text-xs rounded"
                    :class="existingDeviceIds.includes(dev.canid) ? 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300' : 'bg-violet-100 dark:bg-violet-900/40 text-violet-700 dark:text-violet-300'">
                    {{ dev.canid }}
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
import { ref, reactive, computed, inject, watch, onMounted, onUnmounted } from 'vue'
import axios from 'axios'
import DataViewModal from './DataViewModal.vue'
import SingleSettingModal from './SingleSettingModal.vue'
import BatchSettingModal from './BatchSettingModal.vue'
import FwUpgradeModal from './FwUpgradeModal.vue'
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

  cmdWs.onopen = () => {
    wsReady.value = true
  }

  cmdWs.onmessage = (e) => {
    const data = JSON.parse(e.data)
    console.log('ws response:', data)
    if (data.error) {
      console.error('ws error:', data.error)
      return
    }
    if (data.cmd === 'ping_resp') {
      handlePingResp(data)
    }
  }

  cmdWs.onclose = () => {
    cmdWs = null
    wsReady.value = false
    if (shouldReconnect) {
      setTimeout(connectCmdWs, 1000)
    }
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

// isSetupMode(Remote) 전환 시 WebSocket 연결/해제
watch(() => props.isSetupMode, (remote) => {
  if (remote) {
    connectCmdWs()
  } else {
    disconnectCmdWs()
  }
}, { immediate: true })

onUnmounted(() => {
  disconnectCmdWs()
})

const MAX_ROWS = 32

// ── Type definitions (ibsm_map.xlsx 기준) ──
// mtype: 0=None, 1=M10(CT1), 2=M20(CT3), 3=M30(CT6)
// stype (per-CB): 0=not used, 1=1R, 2=1S, 3=1T, 4=3P3W, 5=3P4W, 6=1RZ, 7=1SZ, 8=1TZ
// cbtype (대표값): 1, 4, 5, 6

const mtypeLabel = { 0: 'None', 1: 'M10', 2: 'M20', 3: 'M30' }

const stypeLabel = {
  0: 'N/A', 1: '1PL1', 2: '1PL2', 3: '1PL3',
  4: '3P3W', 5: '3P4W',
  6: '1PL1Z', 7: '1PL2Z', 8: '1PL3Z',
}

const scanDone = ref(false)

// ── devicetype 비트 레이아웃 (16bit) ──
// bit 15: Main/Branch | 14: ST(0)/MM(1) | 13-12: reserved
// bit 11-8: phase mode (0=3P, 1=1P, 2=1P+Z)
// bit 7-4:  stype (기본 stype, 0=not used)
// bit 3-0:  mtype (0=None, 1=M10, 2=M20, 3=M30)
// phase=0(3P) 일 때는 3P3W/3P4W 를 사용자가 선택. 이 때 cbtype 은 4 로 두고 CB 옵션에서 둘 다 노출.
function parseDeviceType(dt) {
  const mtype = dt & 0xF
  const stype = (dt >> 4) & 0xF
  const phase = (dt >> 8) & 0xF
  let cbtype
  if (phase === 0) cbtype = 4          // 3P — 3P3W/3P4W 선택 가능
  else if (phase === 1) cbtype = 1     // 1P
  else if (phase === 2) cbtype = 6     // 1P+Z
  else cbtype = getCbtypeFromStype(stype)
  return { mtype, stype, cbtype }
}

const scannedDevices = ref([])

// maxCb: mtype + cbtype 으로 결정
// M20(CT3): 단상→3, 3상→1 | M30(CT6): 단상→6, 3상→2
function getMaxCb(mtype, cbtype) {
  if (mtype === 1) return 1
  if (mtype === 2) {
    if (cbtype === 4 || cbtype === 5) return 1   // 3상: 3CT=1그룹
    if (cbtype === 6) return 2                    // 1P+Z: 최대 2
    return 3                                      // 단상: 3
  }
  if (mtype === 3) {
    if (cbtype === 4 || cbtype === 5) return 2   // 3상: 6CT=2그룹
    if (cbtype === 6) return 4                    // 1P+Z: 최대 4
    return 6                                      // 단상: 6
  }
  return 0
}

// cbtype별 stype 선택 옵션
// cbtype=4 (phase=0, 3P) 은 3P3W/3P4W 중 사용자가 선택
function getStypeOptions(cbtype) {
  if (cbtype === 1) return [{ value: 1, label: '1PL1' }, { value: 2, label: '1PL2' }, { value: 3, label: '1PL3' }]
  if (cbtype === 4) return [{ value: 4, label: '3P3W' }, { value: 5, label: '3P4W' }]
  if (cbtype === 5) return [{ value: 5, label: '3P4W' }]
  if (cbtype === 6) return [{ value: 6, label: '1PL1Z' }, { value: 7, label: '1PL2Z' }, { value: 8, label: '1PL3Z' }]
  return []
}

// cbtype 라벨 (테이블 Type 컬럼 보조 표시)
const cbtypeMode = { 1: '1P', 4: '3P', 5: '3P', 6: '1P+Z' }

function getCbtypeFromStype(stype) {
  if (stype >= 1 && stype <= 3) return 1
  if (stype === 4) return 4
  if (stype === 5) return 5
  if (stype >= 6 && stype <= 8) return 6
  return 0
}

function createEmptyRow() {
  return {
    id: null,
    name: '',
    type: '',
    mtype: 0,
    cbtype: 0,
    cbCount: 1,
    maxCb: 0,
    stype: [0, 0, 0, 0, 0, 0],
    main: 0,
  }
}

const rows = ref(Array.from({ length: MAX_ROWS }, () => reactive(createEmptyRow())))

const hasRegistered = computed(() => rows.value.some(r => r.id))

const registeredRows = computed(() => {
  return rows.value
    .map((r, i) => ({ ...r, rowIdx: i }))
    .filter(r => r.id)
})

// ibsmSet → rows 로드
function loadFromIbsmSet(data) {
  const tapboxs = data?.tapboxs || []
  for (let i = 0; i < MAX_ROWS; i++) {
    if (i < tapboxs.length && tapboxs[i].CANid) {
      const tb = tapboxs[i]
      const mtype = tb.mtype || 0
      const cbtype = tb.cbtype || 0
      const cbcount = tb.cbcount || 1
      const maxCb = getMaxCb(mtype, cbtype)
      const stype = Array.isArray(tb.stype) ? [...tb.stype, ...Array(6).fill(0)].slice(0, 6) : [0, 0, 0, 0, 0, 0]

      // 기존 등록된 장비를 scannedDevices에 추가
      const devicetype = (cbtype << 8) | (mtype)
      if (!scannedDevices.value.find(d => String(d.canid) === String(tb.CANid))) {
        scannedDevices.value.push({ canid: String(tb.CANid), devicetype })
      }

      Object.assign(rows.value[i], {
        id: String(tb.CANid),
        name: tb.name || '',
        type: mtypeLabel[mtype] || 'Unknown',
        mtype,
        cbtype,
        cbCount: cbcount,
        maxCb,
        stype,
        main: tb.main || 0,
      })
    } else {
      Object.assign(rows.value[i], createEmptyRow())
    }
  }
}

let isSyncing = false

// rows → ibsmSet 동기화
function syncToIbsmSet() {
  isSyncing = true
  const tapboxs = []
  for (let i = 0; i < MAX_ROWS; i++) {
    const row = rows.value[i]
    if (!row.id) continue
    tapboxs.push({
      CANid: row.id,
      index: tapboxs.length,
      mtype: row.mtype,
      cbtype: row.cbtype,
      cbcount: row.cbCount,
      stype: [...row.stype],
      main: 0,
      name: row.name || '',
    })
  }
  setupDict.value.ibsm = { channel: 'ibsm', Enable: 1, tapboxs }
  setTimeout(() => { isSyncing = false }, 0)
}

// setupDict.ibsm 변경 감지 → rows 로드
watch(() => setupDict.value.ibsm, (newVal) => {
  if (isSyncing) return
  if (newVal?.tapboxs?.length) {
    loadFromIbsmSet(newVal)
  }
}, { immediate: true, deep: true })

function getUsedIds() {
  return rows.value.filter(r => r.id).map(r => String(r.id))
}

function getAvailableDevices(rowIdx) {
  const usedIds = getUsedIds()
  const currentId = rows.value[rowIdx].id
  return scannedDevices.value.filter(dev => {
    const idStr = String(dev.canid)
    return idStr === currentId || !usedIds.includes(idStr)
  })
}

function onIdChange(idx) {
  const row = rows.value[idx]
  if (!row.id) {
    Object.assign(row, createEmptyRow())
    syncToIbsmSet()
    return
  }
  const dev = scannedDevices.value.find(d => String(d.canid) === row.id)
  if (!dev) return
  const { mtype, stype, cbtype } = parseDeviceType(dev.devicetype)
  const maxCb = getMaxCb(mtype, cbtype)
  Object.assign(row, {
    name: '',
    type: mtypeLabel[mtype] || 'Unknown',
    mtype,
    cbtype,
    maxCb,
    cbCount: 1,
    stype: Array(6).fill(stype),
    main: 0,
  })
  syncToIbsmSet()
}

function onNameChange(idx) {
  syncToIbsmSet()
}

function onCbCountChange(idx) {
  syncToIbsmSet()
}

function setStype(rowIdx, cbIdx, value) {
  rows.value[rowIdx].stype.splice(cbIdx, 1, value)
  syncToIbsmSet()
}

// ── Single Modal ──
const singleModalOpen = ref(false)
const singleIdx = ref(0)

function openSingleModal(idx) {
  singleIdx.value = idx
  singleModalOpen.value = true
}

// ── Data View Modal ──
const dataModalIdx = ref(0)

function openDataModal(idx) {
  dataModalIdx.value = idx
  canTestModalOpen.value = true
}

// ── Wattage Correction Modal ──
const wattCorrIdx = ref(0)

function openWattCorrModal(idx) {
  wattCorrIdx.value = idx
  wattCorrModalOpen.value = true
}

// ── Batch Modal ──
const batchModalOpen = ref(false)

// ── Tool Modals ──
const rs485TestModalOpen = ref(false)
const canTestModalOpen = ref(false)
const fwUpgradeModalOpen = ref(false)
const wattCorrModalOpen = ref(false)

// ── RS-485 Communication Test ──
const commTestRunning = ref(false)
const commTestResults = ref({})  // { [slotIdx]: 'waiting' | 'online' | 'error' }
const commSelected = ref({})     // { [slotIdx]: true }

const commSelectedCount = computed(() => {
  return Object.keys(commSelected.value).filter(k => commSelected.value[k] && rows.value[k]?.id).length
})

const commSlots = computed(() => {
  return Array.from({ length: MAX_ROWS }, (_, i) => {
    const row = rows.value[i]
    const result = commTestResults.value[i]
    let status = 'rs485-nodev'
    if (row.id) {
      status = result ? `rs485-${result}` : 'rs485-waiting'
    }
    return { id: row.id || null, name: row.name || '', status, selected: !!commSelected.value[i] && !!row.id }
  })
})

function toggleCommSelect(idx) {
  if (!rows.value[idx]?.id) return
  commSelected.value = { ...commSelected.value, [idx]: !commSelected.value[idx] }
}

function commSelectAll() {
  const sel = {}
  rows.value.forEach((r, i) => { if (r.id) sel[i] = true })
  commSelected.value = sel
}

function commDeselectAll() {
  commSelected.value = {}
}

function runCommTest() {
  commTestRunning.value = true
  commTestResults.value = {}
  const targets = rows.value
    .map((r, i) => ({ idx: i, id: r.id }))
    .filter(r => r.id && commSelected.value[r.idx])

  // 선택된 장비에 대해 waiting 상태 설정
  targets.forEach(t => {
    commTestResults.value = { ...commTestResults.value, [t.idx]: 'waiting' }
  })

  // 각 장비에 ping_req 전송
  commPingTargets.value = targets
  targets.forEach(t => {
    sendCmd('ping_req', { canid: Number(t.id) })
  })
}

// ping 응답 처리용
const commPingTargets = ref([])

function handlePingResp(data) {
  const canid = String(data.canid)
  const target = commPingTargets.value.find(t => rows.value[t.idx]?.id === canid)
  if (!target) return
  commTestResults.value = {
    ...commTestResults.value,
    [target.idx]: data.res === 'ok' ? 'online' : 'error'
  }
  // 전부 응답 왔는지 확인
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

const existingDeviceIds = computed(() =>
  scannedDevices.value.map(d => d.canid)
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
  try {
    const res = await axios.get('/setting/scanDevices')
    console.log('[scanDevices] raw response:', res.data)
    if (res.data.success) {
      const tapboxes = res.data.tapboxes || []
      tapboxes.forEach((tb, i) => {
        console.log(
          `[scanDevices] #${i} canid=${tb.canid} type=${tb.type} (0x${Number(tb.type).toString(16)}) ` +
          `→ mtype=${tb.type & 0xF}, subtype=${(tb.type >> 4) & 0xF}, phasemode=${(tb.type >> 8) & 0xF}`
        )
      })
      pendingScanDevices.value = tapboxes.map(tb => ({
        canid: String(tb.canid),
        devicetype: tb.type,
      }))
      console.log('[scanDevices] pendingScanDevices:', pendingScanDevices.value)
    } else {
      console.error('scanDevices error:', res.data.error)
    }
  } catch (e) {
    console.error('scanDevices fetch error:', e)
  } finally {
    scanLoading.value = false
  }
}

function onClear() {
  openConfirm('Clear', 'Are you sure you want to clear the device list?', () => {
    for (let i = 0; i < MAX_ROWS; i++) {
      Object.assign(rows.value[i], createEmptyRow())
    }
    scanDone.value = false
    syncToIbsmSet()
  })
}

function onAutoFill() {
  if (!scannedDevices.value.length) return
  // canid 오름차순 정렬
  const sorted = [...scannedDevices.value].sort((a, b) => String(a.canid).localeCompare(String(b.canid)))
  for (const dev of sorted) {
    if (getUsedIds().includes(String(dev.canid))) continue
    const emptyRow = rows.value.find(r => !r.id)
    if (!emptyRow) break
    const { mtype, stype, cbtype } = parseDeviceType(dev.devicetype)
    const maxCb = getMaxCb(mtype, cbtype)
    Object.assign(emptyRow, {
      id: String(dev.canid),
      name: '',
      type: mtypeLabel[mtype] || 'Unknown',
      mtype,
      cbtype,
      maxCb,
      cbCount: maxCb,
      stype: Array(6).fill(stype),
      main: 0,
    })
  }
  syncToIbsmSet()
}
</script>

<style scoped>
.cmd-btn {
  @apply px-3 py-1.5 text-xs font-semibold rounded-md border border-gray-300 dark:border-gray-600;
  @apply bg-white dark:bg-gray-700 text-gray-600 dark:text-gray-200;
  @apply hover:bg-gray-50 dark:hover:bg-gray-600 transition-colors;
  @apply disabled:opacity-40 disabled:cursor-not-allowed;
}

.cmd-btn-accent {
  @apply px-3 py-1.5 text-xs font-semibold rounded-md border border-gray-400 dark:border-gray-500;
  @apply bg-gray-800 dark:bg-gray-200 text-white dark:text-gray-800;
  @apply hover:bg-gray-700 dark:hover:bg-gray-300 transition-colors;
}

.batch-btn {
  @apply px-4 py-1.5 text-xs font-semibold rounded-md border-[1.5px] border-violet-500;
  @apply bg-white dark:bg-gray-700 text-violet-600 dark:text-violet-400;
  @apply hover:bg-violet-50 dark:hover:bg-violet-900/20 transition-colors;
  @apply disabled:opacity-40 disabled:cursor-not-allowed disabled:hover:bg-white dark:disabled:hover:bg-gray-700;
  @apply flex items-center gap-1.5;
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
  @apply text-xs px-2 py-1 min-w-[50px] rounded border border-gray-200 dark:border-gray-600;
  @apply bg-transparent text-gray-700 dark:text-gray-200;
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

/* Action buttons (설정/데이터/보정) */
.act-btn {
  @apply w-8 h-8 rounded-lg flex items-center justify-center mx-auto transition-all;
  @apply disabled:bg-gray-100 dark:disabled:bg-gray-800 disabled:text-gray-300 dark:disabled:text-gray-600;
  @apply disabled:cursor-not-allowed disabled:shadow-none;
}
.act-violet {
  @apply bg-violet-500 text-white shadow-sm shadow-violet-200 dark:shadow-violet-900/40;
  @apply hover:bg-violet-600 hover:shadow-md hover:shadow-violet-300 dark:hover:shadow-violet-800/40;
}
.act-sky {
  @apply bg-sky-500 text-white shadow-sm shadow-sky-200 dark:shadow-sky-900/40;
  @apply hover:bg-sky-600 hover:shadow-md hover:shadow-sky-300 dark:hover:shadow-sky-800/40;
}
.act-pink {
  @apply bg-pink-500 text-white shadow-sm shadow-pink-200 dark:shadow-pink-900/40;
  @apply hover:bg-pink-600 hover:shadow-md hover:shadow-pink-300 dark:hover:shadow-pink-800/40;
}

/* Mode badge (devicetype에서 결정, 읽기 전용) */
.mode-badge {
  @apply px-1.5 py-0.5 text-[10px] font-bold rounded;
  @apply bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400;
  @apply border border-indigo-200 dark:border-indigo-700/40;
}

/* Chips for CB Setting */
.chip {
  @apply inline-flex items-center px-2 py-1 rounded text-[11px] font-semibold cursor-pointer transition-all;
}
.chip-label {
  @apply bg-transparent text-gray-500 dark:text-gray-400 font-bold cursor-default px-0 py-0 mr-0.5;
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

/* RS-485 legend badges */
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
