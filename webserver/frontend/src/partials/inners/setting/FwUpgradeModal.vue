<template>
  <Teleport to="body">
    <div v-if="open" class="fixed inset-0 bg-black/40 backdrop-blur-sm z-40" @click="$emit('close')"></div>
  </Teleport>
  <Teleport to="body">
    <div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="$emit('close')">
      <div class="bg-white dark:bg-gray-800 rounded-xl shadow-2xl border border-gray-200 dark:border-gray-700 w-[90vw] max-w-[1100px] flex flex-col" @click.stop>

        <!-- Header -->
        <div class="px-5 py-3 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between flex-shrink-0">
          <h3 class="text-sm font-extrabold text-gray-800 dark:text-gray-100">Firmware Upgrade</h3>
          <button @click="$emit('close')" class="w-7 h-7 rounded-full bg-gray-100 dark:bg-gray-700 text-gray-500 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-gray-600 flex items-center justify-center text-sm flex-shrink-0">✕</button>
        </div>

        <!-- Toolbar -->
        <div class="px-5 py-2.5 border-b border-gray-100 dark:border-gray-700/50 flex items-center justify-between flex-shrink-0">
          <button @click="reload" class="fw-btn">Reload</button>
        </div>

        <!-- Firmware File -->
        <div class="px-5 py-4 border-b border-gray-100 dark:border-gray-700/50 flex items-center gap-4 flex-shrink-0">
          <span class="text-xs font-bold text-gray-700 dark:text-gray-200 whitespace-nowrap">Firmware File</span>
          <label class="fw-file-label">
            파일 선택
            <input type="file" accept=".bin,.hex,.fw" class="hidden" @change="onFileSelect">
          </label>
          <span class="text-xs text-gray-500 dark:text-gray-400 truncate">{{ selectedFile ? selectedFile.name : '선택한 파일 없음' }}</span>
        </div>

        <!-- Device List -->
        <div class="px-5 pt-3 pb-1 flex-shrink-0">
          <span class="text-xs font-bold text-gray-700 dark:text-gray-200">Device List</span>
        </div>
        <div class="flex-1 px-5 pb-4 overflow-auto" style="max-height: calc(70vh - 200px);">
          <table class="fw-table">
            <thead>
              <tr>
                <th class="w-10"><input type="checkbox" v-model="allChecked" @change="toggleAll" class="accent-violet-500 w-3.5 h-3.5 cursor-pointer"></th>
                <th>Device</th>
                <th>CAN ID</th>
                <th>TOU Name</th>
                <th>Current Firmware Ver</th>
                <th class="w-[100px]">Status</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(dev, idx) in deviceListStable" :key="dev.id">
                <td class="text-center">
                  <input type="checkbox" v-model="checked[dev.id]" class="accent-violet-500 w-3.5 h-3.5 cursor-pointer">
                </td>
                <td>No.{{ dev.rowIdx + 1 }}</td>
                <td class="font-mono">{{ dev.id }}</td>
                <td>{{ dev.name || '—' }}</td>
                <td class="font-mono">{{ dev.fwVer }}</td>
                <td>
                  <span class="fw-status" :class="dev.upgradeStatus">{{ statusLabel(dev.upgradeStatus) }}</span>
                </td>
              </tr>
              <tr v-if="deviceListStable.length === 0">
                <td colspan="6" class="text-center text-gray-400 dark:text-gray-500 py-8">등록된 장비가 없습니다</td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Progress -->
        <div v-if="upgrading || upgradeDone" class="mx-5 mb-3 p-3 bg-violet-50 dark:bg-violet-900/20 rounded-lg flex-shrink-0">
          <div class="text-[11px] font-semibold text-violet-600 dark:text-violet-400 mb-1.5">{{ progressLabel }}</div>
          <div class="bg-violet-200 dark:bg-violet-800 rounded h-2 overflow-hidden">
            <div class="bg-gradient-to-r from-violet-500 to-violet-400 h-full rounded transition-all duration-300" :style="{ width: progressPct + '%' }"></div>
          </div>
        </div>

        <!-- Footer -->
        <div class="px-5 py-2.5 border-t border-gray-200 dark:border-gray-700 flex justify-end bg-gray-50 dark:bg-gray-800/50 flex-shrink-0">
          <button @click="startUpgrade" class="fw-btn-primary" :disabled="!selectedFile || checkedCount === 0 || upgrading">
            {{ upgrading ? 'Upgrading...' : 'Firmware Upgrade' }}
          </button>
          <button @click="$emit('close')" class="px-5 py-1.5 text-xs font-semibold border border-gray-300 dark:border-gray-600 rounded-md text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700">닫기</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  open: { type: Boolean, default: false },
  registeredRows: { type: Array, default: () => [] },
})

defineEmits(['close'])

const selectedFile = ref(null)
const checked = ref({})
const allChecked = ref(false)
const upgrading = ref(false)
const upgradeDone = ref(false)
const progressLabel = ref('')
const progressPct = ref(0)
const upgradeResults = ref({}) // { [id]: 'waiting' | 'uploading' | 'success' | 'fail' }

const checkedCount = computed(() => {
  return (props.registeredRows || []).filter(r => checked.value[r.id]).length
})

function onFileSelect(e) {
  const file = e.target.files?.[0]
  selectedFile.value = file || null
}

function toggleAll() {
  for (const r of (props.registeredRows || [])) {
    checked.value[r.id] = allChecked.value
  }
}

function statusLabel(status) {
  const m = { idle: '—', waiting: 'Waiting', uploading: 'Uploading', success: 'Success', fail: 'Failed' }
  return m[status] || '—'
}

function reload() {
  upgradeResults.value = {}
  upgradeDone.value = false
  progressLabel.value = ''
  progressPct.value = 0
}

function startUpgrade() {
  if (!selectedFile.value || checkedCount.value === 0) return
  upgrading.value = true
  upgradeDone.value = false
  upgradeResults.value = {}

  const targets = (props.registeredRows || []).filter(r => checked.value[r.id])
  targets.forEach(r => { upgradeResults.value[r.id] = 'waiting' })

  let i = 0
  function next() {
    if (i >= targets.length) {
      upgrading.value = false
      upgradeDone.value = true
      progressLabel.value = `완료 — ${targets.length}개 장비 업그레이드 완료`
      progressPct.value = 100
      return
    }
    const dev = targets[i]
    upgradeResults.value = { ...upgradeResults.value, [dev.id]: 'uploading' }
    progressLabel.value = `업로드 중... iBSM #${dev.id} (${i + 1}/${targets.length})`
    progressPct.value = Math.round(((i + 0.5) / targets.length) * 100)

    setTimeout(() => {
      // 더미: 90% 성공
      const result = Math.random() > 0.1 ? 'success' : 'fail'
      upgradeResults.value = { ...upgradeResults.value, [dev.id]: result }
      progressPct.value = Math.round(((i + 1) / targets.length) * 100)
      i++
      setTimeout(next, 200)
    }, 600)
  }
  next()
}

// 더미 FW 버전 고정 (deviceListStable computed에서 random 방지)
const fwVersions = ref({})
watch(() => props.open, (v) => {
  if (v) {
    checked.value = {}
    allChecked.value = false
    selectedFile.value = null
    upgradeResults.value = {}
    upgrading.value = false
    upgradeDone.value = false
    progressLabel.value = ''
    progressPct.value = 0
    // 더미 FW 버전 생성
    const versions = {}
    ;(props.registeredRows || []).forEach(r => {
      versions[r.id] = `v${60 + Math.floor(Math.random() * 5)}.${Math.floor(Math.random() * 10)}`
    })
    fwVersions.value = versions
  }
})

// fwVer를 고정 더미값으로 재정의
const deviceListStable = computed(() => {
  return (props.registeredRows || []).map(r => ({
    ...r,
    fwVer: fwVersions.value[r.id] || '—',
    upgradeStatus: upgradeResults.value[r.id] || 'idle',
  }))
})
</script>

<style scoped>
.fw-btn {
  @apply px-3 py-1.5 text-[11px] font-semibold rounded-md border border-gray-300 dark:border-gray-600;
  @apply bg-white dark:bg-gray-700 text-gray-600 dark:text-gray-200;
  @apply hover:bg-gray-50 dark:hover:bg-gray-600 transition-colors;
}
.fw-btn-primary {
  @apply px-4 py-1.5 text-[11px] font-semibold rounded-md border border-violet-500;
  @apply bg-white dark:bg-gray-700 text-violet-600 dark:text-violet-400;
  @apply hover:bg-violet-50 dark:hover:bg-violet-900/20 transition-colors;
  @apply disabled:opacity-40 disabled:cursor-not-allowed;
}
.fw-file-label {
  @apply px-3 py-1.5 text-[11px] font-semibold rounded-md border border-gray-300 dark:border-gray-600;
  @apply bg-white dark:bg-gray-700 text-gray-600 dark:text-gray-200 cursor-pointer;
  @apply hover:bg-gray-50 dark:hover:bg-gray-600 transition-colors;
}
.fw-table { @apply w-full text-xs border-collapse; }
.fw-table thead tr { @apply bg-gray-100 dark:bg-gray-700/60 border-b border-gray-200 dark:border-gray-600; }
.fw-table th { @apply px-3 py-2.5 text-[10px] font-bold text-gray-500 dark:text-gray-400 uppercase tracking-wider text-left; }
.fw-table td { @apply px-3 py-2 text-gray-700 dark:text-gray-200 border-b border-gray-100 dark:border-gray-700/30; }
.fw-table tbody tr:hover { @apply bg-gray-50/50 dark:bg-gray-700/20; }

.fw-status { @apply text-[10px] font-bold px-2 py-0.5 rounded; }
.fw-status.idle { @apply text-gray-400; }
.fw-status.waiting { @apply bg-gray-100 dark:bg-gray-700 text-gray-500; }
.fw-status.uploading { @apply bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400; }
.fw-status.success { @apply bg-emerald-100 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400; }
.fw-status.fail { @apply bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400; }
</style>
