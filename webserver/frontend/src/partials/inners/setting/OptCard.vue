<template>
  <div class="setting-card">
    <div class="card-accent bg-amber-500"></div>
    <div class="card-header">
      <div class="card-icon bg-amber-500">
        <svg class="w-3.5 h-3.5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
          <path d="M12 20V10"/><path d="M18 20V4"/><path d="M6 20v-4"/>
        </svg>
      </div>
      <h3 class="card-title">Measurement Option</h3>
    </div>

    <div class="card-body">
      <div class="space-y-3">
        <div class="field">
          <label class="field-label">PF Sign</label>
          <select
            :value="inputDict.opt.pf_sign"
            @change="updateNestedField('opt', 'pf_sign', Number($event.target.value))"
            class="field-select"
          >
            <option :value="0">IEC</option>
            <option :value="1">IEEE</option>
          </select>
        </div>
        <div class="field">
          <label class="field-label">VA Type</label>
          <select
            :value="inputDict.opt.va_type"
            @change="updateNestedField('opt', 'va_type', Number($event.target.value))"
            class="field-select"
          >
            <option :value="0">Arithmetic</option>
            <option :value="1">Vector</option>
          </select>
        </div>
        <div class="field">
          <label class="field-label">Unbalance</label>
          <select
            :value="inputDict.opt.unbalance"
            @change="updateNestedField('opt', 'unbalance', Number($event.target.value))"
            class="field-select"
          >
            <option :value="0">NEMA</option>
            <option :value="1">IEC</option>
          </select>
        </div>
        <div class="field">
          <label class="field-label">I Load (TDD)</label>
          <input
            :value="inputDict.opt.iload"
            @input="updateNestedField('opt', 'iload', parseFloat($event.target.value) || 0)"
            class="field-input" type="number" step="0.1" min="0"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

const inputDict = inject('channel_inputDict')
const updateNestedField = inject('updateNestedField')
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
</style>
