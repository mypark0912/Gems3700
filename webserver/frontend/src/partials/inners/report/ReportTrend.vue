<template>
  <div 
    class="flex flex-col col-span-full sm:col-span-12 xl:col-span-12 shadow rounded-xl px-5 py-1"
    :class="isPdfMode ? 'bg-white' : 'bg-white dark:bg-gray-700 text-white'"
  >
    <LineChart
        :chart-data="option.lineData"
        :chart-labels="option.lineLabels"
        :title="option.lineTitle"
        :mode="option.mode || 'DiagnosisDetail'"
    />
  </div>
</template>

<script setup>
import LineChart from "../../../charts/connect/ReportTrendChart_Echart.vue";
import { ref, inject } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  data: Object,
})

const isPdfMode = inject('isPdfMode', false)

const option = ref(props.data);
</script>