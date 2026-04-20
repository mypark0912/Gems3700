<template>
  <Teleport to="body">
    <div v-if="open" class="fixed inset-0 bg-black/40 backdrop-blur-sm z-40" @click="$emit('close')"></div>
  </Teleport>
  <Teleport to="body">
    <div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="$emit('close')">
      <div class="bg-white dark:bg-gray-800 rounded-xl shadow-2xl border border-gray-200 dark:border-gray-700 w-[80vw] max-w-[860px] flex flex-col" @click.stop>

        <!-- Header -->
        <div class="px-5 py-2.5 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between flex-shrink-0">
          <div class="flex items-center gap-3">
            <h3 class="text-sm font-extrabold text-gray-800 dark:text-gray-100">Wattage Correction</h3>
            <span class="bg-violet-500 text-white px-2.5 py-0.5 rounded-full text-[11px] font-bold">iBSM #{{ canid }}</span>
          </div>
          <button @click="$emit('close')" class="w-6 h-6 rounded-full bg-gray-100 dark:bg-gray-700 text-gray-500 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-gray-600 flex items-center justify-center text-xs flex-shrink-0">✕</button>
        </div>

        <!-- Toolbar -->
        <div class="px-5 py-2 border-b border-gray-100 dark:border-gray-700/50 flex items-center justify-between flex-shrink-0">
          <button @click="reload" class="wc-btn">Reload</button>
          <div class="flex items-center gap-2">
            <button @click="onWrite" class="wc-btn">Write</button>
            <button @click="onClearWh" class="wc-btn">Clear Wh</button>
            <button @click="onReboot" class="wc-btn">Reboot</button>
          </div>
        </div>

        <!-- Device Info (grid) -->
        <div class="px-5 py-2.5 border-b border-gray-100 dark:border-gray-700/50 flex-shrink-0">
          <div class="grid grid-cols-4 gap-x-4 gap-y-1.5 text-[11px]">
            <div class="wc-info-field"><span class="wc-info-lbl">TOU Name</span><span class="wc-info-box">{{ touName || '—' }}</span></div>
            <div class="wc-info-field"><span class="wc-info-lbl">Firmware Ver</span><span class="wc-info-box">{{ devInfo.fwVer }}</span></div>
            <div class="wc-info-field"><span class="wc-info-lbl">Meter Type</span><span class="wc-info-box">{{ devInfo.meterType }}</span></div>
            <div class="wc-info-field">
              <span class="wc-info-lbl">Status</span>
              <span class="wc-info-box" :class="devInfo.status === 'ONLINE' ? '!bg-emerald-50 !text-emerald-600 !border-emerald-200 dark:!bg-emerald-900/20 dark:!text-emerald-400 dark:!border-emerald-700' : ''">{{ devInfo.status }}</span>
            </div>
            <div class="wc-info-field"><span class="wc-info-lbl">Temperature</span><span class="wc-info-box">{{ devInfo.temp }}</span></div>
            <div class="wc-info-field"><span class="wc-info-lbl">Timestamp</span><span class="wc-info-box">{{ devInfo.timestamp }}</span></div>
            <div class="wc-info-field"><span class="wc-info-lbl">Frequency</span><span class="wc-info-box">{{ devInfo.freq }}</span></div>
          </div>
        </div>

        <!-- CB Wattage Table -->
        <div class="px-5 py-3">
          <table class="wc-table">
            <thead>
              <tr>
                <th class="w-[60px]">CB No.</th>
                <th>kWh</th>
                <th>kVARh</th>
                <th>kVAh</th>
                <th>This Month kWh</th>
                <th>Last Month kWh</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(cb, ci) in cbRows" :key="ci">
                <td class="text-center font-medium text-gray-500 text-xs">{{ ci + 1 }}</td>
                <td><input type="number" v-model.number="cb.kwh" class="wc-input" step="0.1"></td>
                <td><input type="number" v-model.number="cb.kvarh" class="wc-input" step="0.1"></td>
                <td><input type="number" v-model.number="cb.kvah" class="wc-input" step="0.1"></td>
                <td><input type="number" v-model.number="cb.kwh_tm" class="wc-input" step="0.1"></td>
                <td><input type="number" v-model.number="cb.kwh_lm" class="wc-input" step="0.1"></td>
              </tr>
            </tbody>
          </table>
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
import { ref, reactive, watch } from 'vue'

const props = defineProps({
  open: { type: Boolean, default: false },
  canid: { type: String, default: '' },
  touName: { type: String, default: '' },
  cbtype: { type: Number, default: 0 },
  cbcount: { type: Number, default: 1 },
})

defineEmits(['close'])

const devInfo = ref({ fwVer: '—', meterType: '—', status: '—', temp: '—', timestamp: '—', freq: '—' })
const cbRows = ref([])

function initData() {
  // 더미 장비 정보
  const is3p = props.cbtype === 4 || props.cbtype === 5
  devInfo.value = {
    fwVer: '00.40h / 0.64',
    meterType: is3p ? 'iBSM30/Direct/3P' : 'iBSM20/Direct/1P',
    status: 'ONLINE',
    temp: `${(12 + Math.random() * 5).toFixed(2)} °C`,
    timestamp: new Date().toISOString().replace('T', ' ').slice(0, 19),
    freq: '60.00',
  }
  cbRows.value = Array.from({ length: props.cbcount }, () => reactive({
    kwh: 0.0, kvarh: 0.0, kvah: 0.0, kwh_tm: 0.0, kwh_lm: 0.0,
  }))
}

function reload() { initData() }

function onWrite() {
  alert(`Write: iBSM #${props.canid} — ${cbRows.value.length} CB(s) 전력량 저장`)
}

function onClearWh() {
  cbRows.value.forEach(cb => { cb.kwh = 0; cb.kvarh = 0; cb.kvah = 0; cb.kwh_tm = 0; cb.kwh_lm = 0 })
}

function onReboot() {
  alert(`Reboot: iBSM #${props.canid}`)
}

watch(() => props.open, (v) => {
  if (v) initData()
})
</script>

<style scoped>
.wc-btn {
  @apply px-3 py-1.5 text-[11px] font-semibold rounded-md border border-gray-300 dark:border-gray-600;
  @apply bg-white dark:bg-gray-700 text-gray-600 dark:text-gray-200;
  @apply hover:bg-gray-50 dark:hover:bg-gray-600 transition-colors;
}
.wc-info-field { @apply flex flex-col gap-0.5; }
.wc-info-lbl { @apply text-[9px] text-gray-400 dark:text-gray-500 font-semibold uppercase tracking-wider; }
.wc-info-box {
  @apply px-2 h-[24px] leading-[24px] bg-gray-50 dark:bg-gray-700/50 border border-gray-200 dark:border-gray-600 rounded text-[11px] text-gray-700 dark:text-gray-200 font-medium truncate;
}
.wc-table { @apply w-full text-xs border-collapse; }
.wc-table thead tr { @apply bg-gray-100 dark:bg-gray-700/60 border-b border-gray-200 dark:border-gray-600; }
.wc-table th { @apply px-2 py-2 text-[10px] font-bold text-gray-500 dark:text-gray-400 uppercase tracking-wider text-center; }
.wc-table td { @apply px-1.5 py-1 border-b border-gray-100 dark:border-gray-700/30; }
.wc-input {
  @apply w-full px-1.5 py-[3px] text-[11px] text-right border border-gray-200 dark:border-gray-600 rounded;
  @apply bg-white dark:bg-gray-700/50 text-gray-700 dark:text-gray-200 tabular-nums;
  @apply focus:outline-none focus:border-violet-500;
}
</style>
