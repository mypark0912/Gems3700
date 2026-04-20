<template>
  <div class="setting-card">
    <div class="card-accent bg-green-500"></div>
    <div class="card-header">
      <div class="card-icon bg-green-500">
        <svg class="w-3.5 h-3.5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M3 12c1.5-2 3-2 4-2s2 .5 3 2 3 2 4 2 2-.5 3-2 3-2 4-2 2 .5 3 2" />
          <circle cx="3" cy="12" r="1.5" fill="currentColor" />
          <circle cx="21" cy="12" r="1.5" fill="currentColor" />
          <path d="M5 12v7m14-7v7" />
        </svg>
      </div>
      <h3 class="card-title">CT</h3>
    </div>

    <div class="card-body">
      <!-- CT Directions -->
      <div class="field">
        <label class="field-label">CT Directions</label>
        <div class="space-y-2 mt-1">
          <div v-for="(ct, idx) in ['CT1', 'CT2', 'CT3']" :key="ct" class="flex items-center gap-3">
            <span class="text-sm font-medium text-gray-500 dark:text-gray-400 w-10">{{ ct }}</span>
            <div class="btn-group">
              <button
                v-for="(opt, oi) in ctDirectionOptions"
                :key="opt.value"
                @click.prevent="setDirection(idx, opt.value)"
                :class="[
                  'btn-toggle',
                  oi === 0 ? 'rounded-l-lg' : '',
                  oi === ctDirectionOptions.length - 1 ? 'rounded-r-lg' : '',
                  inputDict.ctInfo.direction[idx] === opt.value
                    ? 'btn-toggle-active'
                    : 'btn-toggle-inactive',
                ]"
              >
                {{ opt.label }}
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Starting / Rated Current -->
      <div class="grid grid-cols-2 gap-3">
        <div class="field">
          <label class="field-label">Starting Current (mA)</label>
          <input
            :value="inputDict.ctInfo.startingcurrent"
            @input="updateNestedField('ctInfo', 'startingcurrent', $event.target.value)"
            class="field-input" type="text" :maxlength="20"
          />
        </div>
        <div class="field">
          <label class="field-label">Rated Current (A)</label>
          <input
            :value="inputDict.ctInfo.inorminal"
            @input="updateNestedField('ctInfo', 'inorminal', $event.target.value)"
            class="field-input" type="text" :maxlength="20"
          />
        </div>
      </div>

      <!-- Primary / Secondary -->
      <div class="grid grid-cols-2 gap-3">
        <div class="field">
          <label class="field-label">Primary</label>
          <input
            :value="inputDict.ctInfo.ct1"
            @input="updateNestedField('ctInfo', 'ct1', $event.target.value)"
            class="field-input" type="text" :maxlength="20"
          />
        </div>
        <div class="field">
          <label class="field-label">Secondary</label>
          <select
            :value="inputDict.ctInfo.ct2"
            @change="updateNestedField('ctInfo', 'ct2', $event.target.value)"
            class="field-select"
          >
            <option :value="0">5A</option>
            <option :value="1">100mA</option>
            <option :value="2">333mV</option>
            <option :value="3">Rogowski</option>
          </select>
        </div>
      </div>

      <!-- ZCT Scale / Type -->
      <div class="grid grid-cols-2 gap-3">
        <div class="field">
          <label class="field-label">ZCT Scale</label>
          <input
            :value="inputDict.ctInfo.zctscale"
            @input="updateNestedField('ctInfo', 'zctscale', $event.target.value)"
            class="field-input" type="text" :maxlength="20"
          />
        </div>
        <div class="field">
          <label class="field-label">ZCT Type</label>
          <select
            :value="inputDict.ctInfo.zcttpye"
            @change="updateNestedField('ctInfo', 'zcttpye', $event.target.value)"
            class="field-select"
          >
            <option :value="0">None</option>
            <option :value="1">200mA:100mV</option>
            <option :value="2">200mA:1.5mA</option>
            <option :value="3">200mV:0.1mA</option>
          </select>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

const inputDict = inject('channel_inputDict')
const updateNestedField = inject('updateNestedField')
const updateArrayField = inject('updateArrayField')

const ctDirectionOptions = [
  { value: 0, label: 'Positive' },
  { value: 1, label: 'Negative' },
]

const setDirection = (index, value) => {
  updateArrayField('ctInfo', 'direction', index, value)
}
</script>

<style scoped>
.setting-card {
  @apply relative bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg overflow-hidden flex flex-col;
}
.card-accent {
  @apply absolute top-0 left-0 right-0 h-0.5;
}
.card-header {
  @apply flex items-center gap-2.5 px-5 py-4 border-b border-gray-200 dark:border-gray-700/60;
}
.card-icon {
  @apply w-6 h-6 rounded-full flex items-center justify-center shrink-0;
}
.card-title {
  @apply text-base font-bold text-gray-800 dark:text-gray-100;
}
.card-body {
  @apply px-5 py-4 flex-1 space-y-1;
}
.field {
  @apply mb-3;
}
.field-label {
  @apply block text-sm font-medium text-gray-600 dark:text-gray-400 mb-1.5;
}
.field-input {
  @apply w-full h-9 px-3 text-sm rounded-lg border border-gray-200 dark:border-gray-600;
  @apply bg-gray-50 dark:bg-gray-700 text-gray-800 dark:text-gray-200;
  @apply focus:outline-none focus:ring-1 focus:ring-violet-400 focus:border-violet-400 transition-colors;
}
.field-select {
  @apply w-full h-9 px-3 text-sm rounded-lg border border-gray-200 dark:border-gray-600;
  @apply bg-gray-50 dark:bg-gray-700 text-gray-800 dark:text-gray-200;
  @apply focus:outline-none focus:ring-1 focus:ring-violet-400 focus:border-violet-400 transition-colors;
}
.btn-group {
  @apply flex flex-1 rounded-lg overflow-hidden border border-gray-200 dark:border-gray-600;
}
.btn-toggle {
  @apply flex-1 py-1.5 text-sm font-semibold transition-colors border-0;
}
.btn-toggle-active {
  @apply bg-violet-500 text-white;
}
.btn-toggle-inactive {
  @apply bg-white dark:bg-gray-700 text-gray-500 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-600;
}
</style>
