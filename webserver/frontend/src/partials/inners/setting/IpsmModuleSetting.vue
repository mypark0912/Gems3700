<template>
  <div class="grow dark:text-white">
    <div class="p-6 space-y-4">

      <!-- 1행: 모듈 Enable 토글 + 구성하기 버튼 -->
      <div class="flex items-center gap-6">
        <!-- IPSM72 섹션 -->
        <div class="flex items-center gap-4">
          <span class="text-xs font-bold text-gray-400 dark:text-gray-500 uppercase tracking-wider">IPSM72</span>
          <div class="flex items-center gap-3">
            <label class="inline-flex items-center gap-2 cursor-pointer">
              <span class="text-sm font-medium">#1</span>
              <div class="relative inline-flex items-center">
                <input type="checkbox" class="sr-only peer"
                  :checked="enable1" @change="onToggleIpsm(0, $event.target.checked)" />
                <div class="w-9 h-5 bg-gray-200 peer-focus:outline-none rounded-full peer dark:bg-gray-600 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-rose-500"></div>
              </div>
            </label>
            <label class="inline-flex items-center gap-2 cursor-pointer">
              <span class="text-sm font-medium">#2</span>
              <div class="relative inline-flex items-center">
                <input type="checkbox" class="sr-only peer"
                  :checked="enable2" @change="onToggleIpsm(1, $event.target.checked)" />
                <div class="w-9 h-5 bg-gray-200 peer-focus:outline-none rounded-full peer dark:bg-gray-600 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-rose-500"></div>
              </div>
            </label>
          </div>
        </div>

        <div class="w-px h-5 bg-gray-200 dark:bg-gray-600"></div>

        <!-- DI60 섹션 -->
        <div class="flex items-center gap-4">
          <span class="text-xs font-bold text-gray-400 dark:text-gray-500 uppercase tracking-wider">DI60</span>
          <div class="flex items-center gap-4">
            <!-- DI60 #1 -->
            <div class="flex items-center gap-2" :class="hasAnyIpsm ? '' : 'opacity-40'">
              <span class="text-sm font-medium">#1</span>
              <label class="relative inline-flex items-center" :class="hasAnyIpsm ? 'cursor-pointer' : 'cursor-not-allowed'">
                <input type="checkbox" class="sr-only peer"
                  :checked="di60_1" :disabled="!hasAnyIpsm"
                  @change="di60_1 = $event.target.checked; sync()" />
                <div class="w-9 h-5 bg-gray-200 peer-focus:outline-none rounded-full peer dark:bg-gray-600 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-rose-500 peer-disabled:opacity-50"></div>
              </label>
              <select v-if="di60_1 && is2x2" v-model.number="di60Target1" class="target-sel" @change="onDi60TargetChange(0)">
                <option v-for="opt in ipsmOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
              </select>
            </div>
            <!-- DI60 #2 -->
            <div class="flex items-center gap-2" :class="hasAnyIpsm ? '' : 'opacity-40'">
              <span class="text-sm font-medium">#2</span>
              <label class="relative inline-flex items-center" :class="hasAnyIpsm ? 'cursor-pointer' : 'cursor-not-allowed'">
                <input type="checkbox" class="sr-only peer"
                  :checked="di60_2" :disabled="!hasAnyIpsm"
                  @change="di60_2 = $event.target.checked; sync()" />
                <div class="w-9 h-5 bg-gray-200 peer-focus:outline-none rounded-full peer dark:bg-gray-600 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-rose-500 peer-disabled:opacity-50"></div>
              </label>
              <select v-if="di60_2 && is2x2" v-model.number="di60Target2" class="target-sel" @change="onDi60TargetChange(1)">
                <option v-for="opt in ipsmOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
              </select>
            </div>
          </div>
        </div>

        <div class="w-px h-5 bg-gray-200 dark:bg-gray-600"></div>

        <!-- 구성하기 버튼 -->
        <button class="cmd-btn" :disabled="!hasAnyIpsm" @click="configOpen = !configOpen">
          {{ configOpen ? '닫기' : '구성하기' }}
        </button>
      </div>

      <!-- 아코디언 -->
      <template v-if="configOpen">
        <div v-for="ipsm in activeIpsmList" :key="ipsm.idx" class="rounded-xl border border-gray-200 dark:border-gray-700/60 overflow-hidden">
          <!-- 헤더 -->
          <button
            class="w-full flex items-center justify-between px-4 py-3 bg-gray-50 dark:bg-gray-700/50 hover:bg-gray-100 dark:hover:bg-gray-700/80 transition-colors"
            @click="accordionOpen[ipsm.idx] = !accordionOpen[ipsm.idx]"
          >
            <span class="text-sm font-bold text-gray-700 dark:text-gray-200">{{ ipsm.label }}</span>
            <svg class="w-4 h-4 text-gray-400 transition-transform" :class="{ 'rotate-180': accordionOpen[ipsm.idx] }"
              xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
              <polyline points="6 9 12 15 18 9" />
            </svg>
          </button>

          <!-- 바디 -->
          <div v-if="accordionOpen[ipsm.idx]" class="px-4 py-3 space-y-3 border-t border-gray-200 dark:border-gray-700/60">
            <!-- 범위 선택 + Auto Fill / Clear -->
            <div class="flex items-center gap-2">
              <select v-model.number="selectedRange[ipsm.idx]" class="target-sel">
                <option :value="-1">전체 (1-72)</option>
                <option v-for="col in 6" :key="col" :value="col - 1">{{ (col - 1) * 12 + 1 }} - {{ col * 12 }}</option>
              </select>
              <button class="cmd-btn" :disabled="!activeDi60.length" @click="onAutoFill(ipsm.idx)">Auto Fill</button>
              <button class="cmd-btn" @click="onClear(ipsm.idx)">Clear</button>
            </div>

            <!-- 6단 테이블 -->
            <div class="grid grid-cols-6 gap-2">
              <div v-for="col in 6" :key="col" class="rounded-lg border border-gray-200 dark:border-gray-700/60 overflow-hidden">
                <table class="w-full text-[13px]">
                  <thead>
                    <tr class="bg-gray-50 dark:bg-gray-700/80 border-b border-gray-200 dark:border-gray-600">
                      <th class="px-2 py-2 text-[11px] font-bold text-gray-400 dark:text-gray-500 text-center w-10">No.</th>
                      <th class="px-2 py-2 text-[11px] font-bold text-gray-400 dark:text-gray-500 text-center">DI60</th>
                      <th class="px-2 py-2 text-[11px] font-bold text-gray-400 dark:text-gray-500 text-center">Pt</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="row in 12" :key="row"
                        class="border-b border-gray-100 dark:border-gray-700/30"
                        :class="row % 2 === 0 ? 'bg-gray-50/50 dark:bg-gray-700/10' : ''">
                      <td class="px-2 py-[5px] text-center text-xs font-medium text-gray-500 dark:text-gray-400">{{ (col - 1) * 12 + row }}</td>
                      <td class="px-2 py-[5px]">
                        <select
                          :value="mapping[ipsm.idx][(col - 1) * 12 + row - 1].mod"
                          class="di-sel"
                          @change="setMod(ipsm.idx, (col - 1) * 12 + row - 1, $event.target.value)">
                          <option :value="''">-</option>
                          <option v-for="opt in activeDi60" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
                        </select>
                      </td>
                      <td class="px-2 py-[5px]">
                        <select
                          v-if="mapping[ipsm.idx][(col - 1) * 12 + row - 1].mod != null"
                          :value="mapping[ipsm.idx][(col - 1) * 12 + row - 1].pt"
                          class="di-sel"
                          @change="setPt(ipsm.idx, (col - 1) * 12 + row - 1, $event.target.value)">
                          <option :value="''">-</option>
                          <option v-for="d in availablePoints(ipsm.idx, (col - 1) * 12 + row - 1)" :key="d" :value="d">{{ d }}</option>
                        </select>
                        <span v-else class="text-gray-300 dark:text-gray-600 text-xs px-2">-</span>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>
      </template>

    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'

const props = defineProps({
  setupDict: { type: Object, required: true }
})

// ── 로컬 반응형 상태 ──
const enable1 = ref(false)
const enable2 = ref(false)
const di60_1 = ref(false)
const di60_2 = ref(false)
const di60Target1 = ref(0)
const di60Target2 = ref(0)

const EMPTY = () => ({ mod: null, pt: null })
const EMPTY_72 = () => Array.from({ length: 72 }, EMPTY)

// 매핑 데이터 — reactive 배열 (인덱스 할당도 반응성 유지)
const mapping = reactive({
  0: EMPTY_72(),
  1: EMPTY_72(),
})

const configOpen = ref(false)
const accordionOpen = reactive({ 0: true, 1: false })
const selectedRange = reactive({ 0: -1, 1: -1 })

// ── computed ──
const hasAnyIpsm = computed(() => enable1.value || enable2.value)
const is2x2 = computed(() => enable1.value && enable2.value && di60_1.value && di60_2.value)

const ipsmOptions = computed(() => {
  const opts = []
  if (enable1.value) opts.push({ value: 0, label: 'IPSM72 #1' })
  if (enable2.value) opts.push({ value: 1, label: 'IPSM72 #2' })
  return opts
})

const activeIpsmList = computed(() => {
  const list = []
  if (enable1.value) list.push({ idx: 0, label: 'IPSM72 #1' })
  if (enable2.value) list.push({ idx: 1, label: 'IPSM72 #2' })
  return list
})

// 활성화된 DI60 전체 (테이블 콤보박스 + Auto Fill 버튼 활성화용)
const activeDi60 = computed(() => {
  const opts = []
  if (di60_1.value) opts.push({ value: 0, label: 'DI60 #1' })
  if (di60_2.value) opts.push({ value: 1, label: 'DI60 #2' })
  return opts
})

function isLinked(target, ipsmIdx) {
  return target === ipsmIdx
}

function getLinkedDi60(ipsmIdx) {
  const opts = []
  if (di60_1.value && isLinked(di60Target1.value, ipsmIdx)) opts.push({ value: 0, label: 'DI60 #1' })
  if (di60_2.value && isLinked(di60Target2.value, ipsmIdx)) opts.push({ value: 1, label: 'DI60 #2' })
  return opts
}

function availablePoints(ipsmIdx, currentIdx) {
  const arr = mapping[ipsmIdx]
  const mod = arr[currentIdx].mod
  if (mod == null) return []
  const used = new Set()
  for (let i = 0; i < 72; i++) {
    if (i === currentIdx) continue
    if (arr[i].mod === mod && arr[i].pt != null) used.add(arr[i].pt)
  }
  const list = []
  for (let d = 1; d <= 30; d++) {
    if (!used.has(d)) list.push(d)
  }
  return list
}

// ── 이벤트 핸들러 ──
function setMod(ipsmIdx, ptIdx, raw) {
  const mod = raw === '' ? null : Number(raw)
  mapping[ipsmIdx][ptIdx] = { mod, pt: null }
  sync()
}

function setPt(ipsmIdx, ptIdx, raw) {
  const pt = raw === '' ? null : Number(raw)
  mapping[ipsmIdx][ptIdx] = { ...mapping[ipsmIdx][ptIdx], pt }
  sync()
}

function onToggleIpsm(num, checked) {
  if (num === 0) enable1.value = checked; else enable2.value = checked
  if (!checked) {
    mapping[num] = EMPTY_72()
    const remaining = enable1.value ? 0 : enable2.value ? 1 : null
    if (di60Target1.value === num && remaining != null) di60Target1.value = remaining
    if (di60Target2.value === num && remaining != null) di60Target2.value = remaining
  }
  if (!enable1.value && !enable2.value) {
    di60_1.value = false
    di60_2.value = false
  }
  sync()
}

function onDi60TargetChange(di60Mod) {
  clearDi60(di60Mod)
  sync()
}

function clearDi60(di60Mod) {
  for (const key of [0, 1]) {
    const arr = mapping[key]
    for (let i = 0; i < 72; i++) {
      if (arr[i].mod === di60Mod) arr[i] = EMPTY()
    }
  }
}

// DI60 끌 때 매핑 정리
watch(di60_1, (v) => { if (!v) { clearDi60(0); sync() } })
watch(di60_2, (v) => { if (!v) { clearDi60(1); sync() } })

// ── Auto Fill / Clear ──
function getRange(ipsmIdx) {
  const sel = selectedRange[ipsmIdx]
  return sel === -1 ? { start: 0, end: 72 } : { start: sel * 12, end: sel * 12 + 12 }
}

function onAutoFill(ipsmIdx) {
  const arr = mapping[ipsmIdx]
  const { start, end } = getRange(ipsmIdx)
  const avail = []
  if (is2x2.value) {
    // 2:2: 타겟에 연결된 DI60만
    if (di60_1.value && isLinked(di60Target1.value, ipsmIdx)) {
      for (let i = 1; i <= 30; i++) avail.push({ mod: 0, pt: i })
    }
    if (di60_2.value && isLinked(di60Target2.value, ipsmIdx)) {
      for (let i = 1; i <= 30; i++) avail.push({ mod: 1, pt: i })
    }
  } else {
    // 그 외: 활성 DI60 전부
    if (di60_1.value) {
      for (let i = 1; i <= 30; i++) avail.push({ mod: 0, pt: i })
    }
    if (di60_2.value) {
      for (let i = 1; i <= 30; i++) avail.push({ mod: 1, pt: i })
    }
  }
  // 이미 사용된 포인트 제외
  const used = new Set()
  // 같은 IPSM72 내 범위 밖
  for (let i = 0; i < 72; i++) {
    if (i >= start && i < end) continue
    if (arr[i].mod != null && arr[i].pt != null) used.add(arr[i].mod + '_' + arr[i].pt)
  }
  // 2:2가 아니면 다른 IPSM72에서 사용 중인 포인트도 제외
  if (!is2x2.value) {
    const otherIdx = ipsmIdx === 0 ? 1 : 0
    const otherArr = mapping[otherIdx]
    if (otherArr) {
      for (let i = 0; i < 72; i++) {
        if (otherArr[i].mod != null && otherArr[i].pt != null) used.add(otherArr[i].mod + '_' + otherArr[i].pt)
      }
    }
  }
  const filtered = avail.filter(a => !used.has(a.mod + '_' + a.pt))
  let ai = 0
  for (let i = start; i < end; i++) {
    arr[i] = ai < filtered.length ? { ...filtered[ai++] } : EMPTY()
  }
  sync()
}

function onClear(ipsmIdx) {
  const arr = mapping[ipsmIdx]
  const { start, end } = getRange(ipsmIdx)
  for (let i = start; i < end; i++) {
    arr[i] = EMPTY()
  }
  sync()
}

// ── setupDict 동기화 ──
function sync() {
  const d = props.setupDict.ipsm72
  d.enable1 = enable1.value ? 1 : 0
  d.enable2 = enable2.value ? 1 : 0
  d.di60_1 = di60_1.value ? 1 : 0
  d.di60_2 = di60_2.value ? 1 : 0
  d.di60Target1 = di60Target1.value
  d.di60Target2 = di60Target2.value
  d.ipsm72_1 = mapping[0].map(m => ({ ...m }))
  d.ipsm72_2 = mapping[1].map(m => ({ ...m }))
}

// setupDict → 로컬 상태 로드
function loadFromSetupDict() {
  const d = props.setupDict.ipsm72
  if (!d) return
  enable1.value = d.enable1 === 1
  enable2.value = d.enable2 === 1
  di60_1.value = d.di60_1 === 1
  di60_2.value = d.di60_2 === 1
  di60Target1.value = d.di60Target1 ?? 0
  di60Target2.value = d.di60Target2 ?? 0
  if (Array.isArray(d.ipsm72_1)) {
    for (let i = 0; i < 72; i++) mapping[0][i] = d.ipsm72_1[i] ? { ...d.ipsm72_1[i] } : EMPTY()
  }
  if (Array.isArray(d.ipsm72_2)) {
    for (let i = 0; i < 72; i++) mapping[1][i] = d.ipsm72_2[i] ? { ...d.ipsm72_2[i] } : EMPTY()
  }
}

watch(() => props.setupDict.ipsm72, () => loadFromSetupDict(), { immediate: true })
</script>

<style scoped>
.cmd-btn {
  @apply px-3 py-1.5 text-xs font-semibold rounded-md border border-gray-300 dark:border-gray-600;
  @apply bg-white dark:bg-gray-700 text-gray-600 dark:text-gray-200;
  @apply hover:bg-gray-50 dark:hover:bg-gray-600 transition-colors;
  @apply disabled:opacity-40 disabled:cursor-not-allowed;
}

.target-sel {
  @apply text-xs px-3 py-1 rounded border border-gray-300 dark:border-gray-600 min-w-[120px];
  @apply bg-white dark:bg-gray-700 text-gray-700 dark:text-gray-200;
  @apply focus:outline-none focus:ring-1 focus:ring-violet-400;
}

.di-sel {
  @apply w-full text-xs px-2 py-1 rounded border border-gray-200 dark:border-gray-600;
  @apply bg-transparent text-gray-700 dark:text-gray-200;
  @apply focus:outline-none focus:ring-1 focus:ring-violet-400;
}
</style>
