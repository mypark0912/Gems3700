<template>
  <div class="setting-card">
    <div class="card-accent bg-sky-500"></div>
    <div class="card-header">
      <div class="card-icon bg-sky-500">
        <svg class="w-3.5 h-3.5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
          <path d="M13 10V3L4 14h7v7l9-11h-7z" />
        </svg>
      </div>
      <h3 class="card-title">PT</h3>
    </div>

    <div class="card-body">
      <!-- Line Frequency -->
      <div class="field">
        <label class="field-label">Line Frequency (Hz)</label>
        <select
            :value="inputDict.ptInfo.linefrequency"
            @change="updateNestedField('ptInfo', 'linefrequency', Number($event.target.value))"
            class="field-select"
          >
            <option :value="0">50Hz</option>
            <option :value="1">60Hz</option>
          </select>
      </div>

      <!-- Wiring Mode / Rated Voltage -->
      <div class="grid grid-cols-2 gap-3">
        <div class="field">
          <label class="field-label">Wiring Mode</label>
          <select
            :value="inputDict.ptInfo.wiringmode"
            @change="updateNestedField('ptInfo', 'wiringmode', $event.target.value)"
            class="field-select"
          >
            <option :value="0">3P4W</option>
            <option :value="1">3P3W(2CT)</option>
            <option :value="2">3P3W(3CT)</option>
            <option :value="3">1P2W</option>
            <option :value="4">1P3W</option>
          </select>
        </div>
        <div class="field">
          <label class="field-label">Rated Voltage (V)</label>
          <input
            :value="inputDict.ptInfo.vnorminal"
            @input="updateNestedField('ptInfo', 'vnorminal', $event.target.value)"
            class="field-input" type="text" :maxlength="20"
          />
        </div>
      </div>

      <!-- Primary / Secondary -->
      <div class="grid grid-cols-2 gap-3">
        <div class="field">
          <label class="field-label">Primary</label>
          <input
            :value="inputDict.ptInfo.pt1"
            @input="updateNestedField('ptInfo', 'pt1', $event.target.value)"
            class="field-input" type="text" :maxlength="20"
          />
        </div>
        <div class="field">
          <label class="field-label">Secondary</label>
          <input
            :value="inputDict.ptInfo.pt2"
            @input="updateNestedField('ptInfo', 'pt2', $event.target.value)"
            class="field-input" type="text" :maxlength="20"
          />
        </div>
      </div>

      <!-- Dashboard Output -->
      <div class="field">
        <label class="field-label">{{ t('config.channelPanel.dashoutput') }}</label>
        <select
          :value="inputDict.ptInfo.dash"
          @change="updateNestedField('ptInfo', 'dash', Number($event.target.value))"
          class="field-select"
        >
          <option :value="0">Phase Voltage</option>
          <option :value="1">Line Voltage</option>
        </select>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
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
  @apply w-full h-10 px-3 text-sm rounded-lg border border-gray-200 dark:border-gray-600;
  @apply bg-gray-50 dark:bg-gray-700 text-gray-800 dark:text-gray-200;
  @apply focus:outline-none focus:ring-1 focus:ring-violet-400 focus:border-violet-400 transition-colors;
}
</style>
