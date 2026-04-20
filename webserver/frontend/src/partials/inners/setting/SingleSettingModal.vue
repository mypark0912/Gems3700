<template>
  <Teleport to="body">
    <div v-if="open" class="fixed inset-0 bg-black/40 backdrop-blur-sm z-40" @click="$emit('close')"></div>
  </Teleport>
  <Teleport to="body">
    <div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="$emit('close')">
      <div class="bg-white dark:bg-gray-800 rounded-xl shadow-2xl border border-gray-200 dark:border-gray-700 w-[90vw] max-w-[1080px] h-[84vh] max-h-[780px] flex flex-col" @click.stop>
        <!-- Header -->
        <div class="px-5 py-3 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between flex-shrink-0">
          <div class="flex items-center gap-3 flex-wrap">
            <h3 class="text-[15px] font-extrabold text-gray-800 dark:text-gray-100 flex items-center gap-2">
              설정
              <span class="bg-violet-500 text-white px-3 py-0.5 rounded-full text-xs font-bold">
                No.{{ rowIndex + 1 }} — iBSM #{{ canid }}
              </span>
            </h3>
            <div class="flex gap-1.5 flex-wrap">
              <button class="modal-action-btn">Read</button>
              <button class="modal-action-btn primary">Write</button>
              <button class="modal-action-btn danger">Factory Reset</button>
            </div>
          </div>
          <button @click="$emit('close')" class="w-7 h-7 rounded-full bg-gray-100 dark:bg-gray-700 text-gray-500 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-gray-600 flex items-center justify-center text-sm flex-shrink-0">✕</button>
        </div>
        <!-- Body -->
        <div class="flex-1 overflow-y-auto p-5">
          <!-- Read-only info -->
          <div class="bg-blue-50/70 dark:bg-blue-900/20 border border-blue-200/60 dark:border-blue-700/40 rounded-lg p-4 mb-4">
            <div class="flex items-center gap-2 text-[11px] font-bold text-blue-500 dark:text-blue-400 uppercase tracking-wider mb-3 pb-2 border-b border-blue-200/60 dark:border-blue-700/40">
              {{ touName || '장비 정보' }} <span class="bg-blue-500 text-white text-[9px] px-1.5 py-0.5 rounded normal-case tracking-normal">읽기 전용</span>
            </div>
            <div class="grid grid-cols-3 gap-2">
              <div class="ro-field"><span class="ro-lbl">Firmware Ver</span><div class="ro-val">{{ readData.fw }}</div></div>
              <div class="ro-field"><span class="ro-lbl">Meter Type</span><div class="ro-val">{{ readData.meter }}</div></div>
              <div class="ro-field"><span class="ro-lbl">Frequency (Hz)</span><div class="ro-val">{{ readData.freq }}</div></div>
              <div class="ro-field"><span class="ro-lbl">Status</span><div class="ro-val">{{ readData.status }}</div></div>
              <div class="ro-field"><span class="ro-lbl">Temperature (°C)</span><div class="ro-val">{{ readData.temp }}</div></div>
              <div class="ro-field"><span class="ro-lbl">Timestamp</span><div class="ro-val">{{ readData.ts }}</div></div>
            </div>
          </div>
          <!-- Editable settings -->
          <div class="text-[11px] font-bold text-gray-600 dark:text-gray-300 uppercase tracking-wider mb-3 pb-2 border-b-2 border-violet-100 dark:border-violet-900/40">
            설정값
          </div>
          <div class="grid grid-cols-[200px_1fr_1fr] gap-px bg-gray-200 dark:bg-gray-700 border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
            <!-- Col 1: PT -->
            <div class="bg-white dark:bg-gray-800 p-3">
              <div class="text-[11px] font-bold text-violet-600 dark:text-violet-400 pb-2 mb-3 border-b-2 border-violet-100 dark:border-violet-900/40">PT</div>
              <div class="sf-narrow"><span class="sf-lbl">Phase Mode</span><select class="sf-input"><option>3P</option><option>1P</option></select></div>
              <div class="sf-narrow"><span class="sf-lbl">Frequency</span><select class="sf-input"><option>60</option><option>50</option></select></div>
              <div class="sf-narrow"><span class="sf-lbl">PT1 (V)</span><input type="number" class="sf-input" placeholder="0" min="0"></div>
              <div class="sf-narrow"><span class="sf-lbl">PT2 (V)</span><input type="number" class="sf-input" placeholder="0" min="0"></div>
              <div class="sf-narrow"><span class="sf-lbl">ZCT Ratio</span><select class="sf-input"><option>200:1.5</option><option>100:1.5</option><option>50:1.5</option></select></div>
            </div>
            <!-- Col 2: CT -->
            <div class="bg-white dark:bg-gray-800 p-3">
              <div class="text-[11px] font-bold text-violet-600 dark:text-violet-400 pb-2 mb-3 border-b-2 border-violet-100 dark:border-violet-900/40">CT</div>
              <div class="grid grid-cols-2 gap-x-3">
                <div class="sf-narrow"><span class="sf-lbl">CT Ratio (A)</span><input type="number" class="sf-input" placeholder="0" min="0"></div>
                <div class="sf-narrow"><span class="sf-lbl">I1 (mA)</span><input type="number" class="sf-input" placeholder="0" min="0"></div>
                <div class="sf-narrow"><span class="sf-lbl">I2 (mA)</span><input type="number" class="sf-input" placeholder="0" min="0"></div>
                <div class="sf-narrow"><span class="sf-lbl">I3 (mA)</span><input type="number" class="sf-input" placeholder="0" min="0"></div>
                <div class="sf-narrow"><span class="sf-lbl">I4 (mA)</span><input type="number" class="sf-input" placeholder="0" min="0"></div>
                <div class="sf-narrow"><span class="sf-lbl">I5 (mA)</span><input type="number" class="sf-input" placeholder="0" min="0"></div>
                <div class="sf-narrow"><span class="sf-lbl">I6 (mA)</span><input type="number" class="sf-input" placeholder="0" min="0"></div>
                <div class="sf-narrow"><span class="sf-lbl">Ig1 (mA)</span><input type="number" class="sf-input" placeholder="0" min="0"></div>
                <div class="sf-narrow"><span class="sf-lbl">Ig2 (mA)</span><input type="number" class="sf-input" placeholder="0" min="0"></div>
                <div class="sf-narrow"><span class="sf-lbl">Ig3 (mA)</span><input type="number" class="sf-input" placeholder="0" min="0"></div>
                <div class="sf-narrow"><span class="sf-lbl">Ig4 (mA)</span><input type="number" class="sf-input" placeholder="0" min="0"></div>
              </div>
            </div>
            <!-- Col 3: Event -->
            <div class="bg-white dark:bg-gray-800 p-3">
              <div class="text-[11px] font-bold text-violet-600 dark:text-violet-400 pb-2 mb-3 border-b-2 border-violet-100 dark:border-violet-900/40">Event</div>
              <div class="grid grid-cols-2 gap-x-3">
                <div class="sf-narrow"><span class="sf-lbl">OCR 1 (A)</span><input type="number" class="sf-input" placeholder="0" min="0"></div>
                <div class="sf-narrow"><span class="sf-lbl">OCR 2 (A)</span><input type="number" class="sf-input" placeholder="0" min="0"></div>
                <div class="sf-narrow"><span class="sf-lbl">OCR 3 (A)</span><input type="number" class="sf-input" placeholder="0" min="0"></div>
                <div class="sf-narrow"><span class="sf-lbl">OCR 4 (A)</span><input type="number" class="sf-input" placeholder="0" min="0"></div>
                <div class="sf-narrow"><span class="sf-lbl">OCR 5 (A)</span><input type="number" class="sf-input" placeholder="0" min="0"></div>
                <div class="sf-narrow"><span class="sf-lbl">OCR 6 (A)</span><input type="number" class="sf-input" placeholder="0" min="0"></div>
                <div class="sf-narrow"><span class="sf-lbl">ELD 1 (mA)</span><input type="number" class="sf-input" placeholder="0" min="0"></div>
                <div class="sf-narrow"><span class="sf-lbl">ELD 2 (mA)</span><input type="number" class="sf-input" placeholder="0" min="0"></div>
                <div class="sf-narrow"><span class="sf-lbl">ELD 3 (mA)</span><input type="number" class="sf-input" placeholder="0" min="0"></div>
                <div class="sf-narrow"><span class="sf-lbl">ELD 4 (mA)</span><input type="number" class="sf-input" placeholder="0" min="0"></div>
                <div class="sf-narrow"><span class="sf-lbl">Sag (%)</span><input type="number" class="sf-input" placeholder="0" min="0" max="100"></div>
                <div class="sf-narrow"><span class="sf-lbl">Swell (%)</span><input type="number" class="sf-input" placeholder="0" min="0" max="100"></div>
                <div class="sf-narrow"><span class="sf-lbl">Hold Timer</span><input type="number" class="sf-input" placeholder="0" min="0"></div>
              </div>
            </div>
          </div>
        </div>
        <!-- Footer -->
        <div class="px-5 py-3 border-t border-gray-200 dark:border-gray-700 flex justify-end gap-2 bg-gray-50 dark:bg-gray-800/50 flex-shrink-0">
          <button class="modal-action-btn warning">ROM Save</button>
          <button @click="$emit('close')" class="px-5 py-2 text-xs font-semibold border border-gray-300 dark:border-gray-600 rounded-md text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700">닫기</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  open: { type: Boolean, default: false },
  rowIndex: { type: Number, default: 0 },
  canid: { type: String, default: '' },
  touName: { type: String, default: '' },
})

defineEmits(['close'])

const readData = ref({ fw: '—', meter: '—', freq: '—', status: '—', temp: '—', ts: '—' })

watch(() => props.open, (v) => {
  if (v) {
    readData.value = { fw: '—', meter: '—', freq: '—', status: '—', temp: '—', ts: '—' }
  }
})
</script>

<style scoped>
.modal-action-btn {
  @apply px-3 py-1.5 text-[11px] font-medium rounded-md border border-gray-300 dark:border-gray-600;
  @apply bg-white dark:bg-gray-700 text-gray-600 dark:text-gray-300;
  @apply hover:bg-gray-50 dark:hover:bg-gray-600 transition-colors cursor-pointer;
}
.modal-action-btn.primary { @apply bg-violet-500 text-white border-violet-500 hover:bg-violet-600; }
.modal-action-btn.warning { @apply border-amber-500 text-amber-600 dark:text-amber-400 hover:bg-amber-50 dark:hover:bg-amber-900/20; }
.modal-action-btn.danger { @apply border-red-500 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20; }

.ro-field { @apply flex flex-col gap-1; }
.ro-lbl { @apply text-[10px] text-blue-400 dark:text-blue-500 font-semibold uppercase tracking-wider; }
.ro-val { @apply px-2.5 h-[30px] leading-[30px] bg-white dark:bg-gray-700 border border-blue-200/60 dark:border-blue-700/40 rounded text-[11px] text-blue-600 dark:text-blue-400 font-semibold; }

.sf-narrow { @apply flex items-center gap-1.5 mb-[5px]; }
.sf-lbl { @apply text-[10px] text-gray-500 dark:text-gray-400 whitespace-nowrap min-w-[70px]; }
.sf-input {
  @apply px-1.5 py-[3px] border border-gray-200 dark:border-gray-600 rounded text-[11px] text-gray-700 dark:text-gray-200;
  @apply bg-gray-50 dark:bg-gray-700/50 w-full min-w-0;
  @apply focus:outline-none focus:border-violet-500 focus:bg-white dark:focus:bg-gray-700;
}
</style>
