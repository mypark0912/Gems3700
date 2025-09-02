<template>
  <div class="flex flex-col col-span-full sm:col-span-12 xl:col-span-12 bg-white shadow rounded-xl px-5 py-1 dark:bg-gray-700">
    <LineChart
        :chart-data="option.lineData"
        :chart-labels="option.lineLabels"
        :title = "option.lineTitle"
    />
    
  </div>
</template>

<script setup>
import LineChart from "../../../charts/connect/ReportTrendChart_Echart.vue";
import {ref} from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  data: Object,
})

const option = ref(props.data);
</script>