<template>
  <div v-if="ready">
    <div class="text-sm font-semibold text-gray-700 mb-2">{{ title }}</div>
    <VChart :option="chartOption" style="height: 300px;" />
  </div>
  <div v-else class="h-[300px] flex items-center justify-center text-gray-400 text-sm">
    Loading chart...
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { use } from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, LegendComponent, GridComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'

use([LineChart, TitleComponent, TooltipComponent, LegendComponent, GridComponent, CanvasRenderer])

const props = defineProps({
  label: {
    type: Array,
    required: true
  },
  data: {
    type: Array, // [[V1], [V2], [V3]] or [[Vfft]]
    required: true
  },
  title: {
    type: String,
    default: ''
  },
  legend: {
    type: Array,
    required: true
  },
  mode: {
    type: String, // '3phase' or '1phase'
    required: true
  }
})

const chartOption = ref({})

const ready = computed(() => {
  const labelLen = props.label.length
  if (!labelLen) return false
  if (props.mode === '3phase') {
    return props.data.length === 3 && props.data.every(arr => arr.length === labelLen)
  } else {
    return props.data.length === 1 && props.data[0].length === labelLen
  }
})

watch([() => props.label, () => props.data, () => props.mode], ([t, v, mode]) => {
  if (!t.length || !v.length) return

  const series =
    mode === '3phase'
      ? v.map((arr, i) => ({
          name: props.legend[i] ?? `S${i + 1}`,
          type: 'line',
          showSymbol: false,
          data: t.map((x, idx) => [x, arr[idx]])
        }))
      : [
          {
            name: props.legend[0] ?? 'S1',
            type: 'line',
            showSymbol: false,
            data: t.map((x, idx) => [x, v[0][idx]])
          }
        ]

  chartOption.value = {
    title: { text: props.title, left: 'center', top: 10 },
    tooltip: { trigger: 'axis' },
    legend: { data: props.legend, top: 30 },
    grid: { top: 60, left: 40, right: 30, bottom: 40 },
    xAxis: { type: 'value', name: props.mode === '3phase' ? 'Time (s)' : 'Frequency (Hz)' },
    yAxis: { type: 'value', name: 'Amplitude' },
    series,
    dataZoom: [
      { type: 'inside', throttle: 20 },
      { type: 'slider', height: 20, bottom: 10 }
    ]
  }
}, { immediate: true })
</script>
