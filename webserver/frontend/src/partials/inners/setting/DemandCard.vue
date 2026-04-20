<template>
  <div class="setting-card">
    <div class="card-accent bg-indigo-500"></div>
    <div class="card-header">
      <div class="card-icon bg-indigo-500">
        <svg class="w-3.5 h-3.5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
          <path d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
      </div>
      <h3 class="card-title">Demand</h3>
    </div>

    <div class="card-body">
      <!-- Target / Interval -->
      <div class="space-y-3">
        <div class="field">
          <label class="field-label">Target (W)</label>
          <input
            :value="inputDict.demand.target"
            @input="updateNestedField('demand', 'target', Number($event.target.value))"
            class="field-input" type="number" :maxlength="20"
          />
        </div>
        <div class="field">
          <label class="field-label">{{ t('config.plansPanel.demand.interval') }}</label>
          <select
            :value="inputDict.demand.demand_interval"
            @change="updateNestedField('demand', 'demand_interval', Number($event.target.value))"
            class="field-select"
          >
            <option :value="0">1 {{ t('config.plansPanel.demand.minutes') }}</option>
            <option :value="1">5 {{ t('config.plansPanel.demand.minutes') }}</option>
            <option :value="2">15 {{ t('config.plansPanel.demand.minutes') }}</option>
            <option :value="3">30 {{ t('config.plansPanel.demand.minutes') }}</option>
            <option :value="4">1 {{ t('config.plansPanel.demand.hours') }}</option>
          </select>
        </div>
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
  @apply w-full h-9 px-3 text-sm rounded-lg border border-gray-200 dark:border-gray-600;
  @apply bg-gray-50 dark:bg-gray-700 text-gray-800 dark:text-gray-200;
  @apply focus:outline-none focus:ring-1 focus:ring-violet-400 focus:border-violet-400 transition-colors;
}
</style>
