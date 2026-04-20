<template>
  <div v-if="ready">
    <div class="text-sm font-semibold text-gray-700 dark:text-gray-100 mb-2">{{ title }}</div>
    <VChart :option="chartOption" style="height: 300px;" />
  </div>
  <div v-else class="h-[300px] flex items-center justify-center text-gray-400 text-sm">
    {{ noData ? 'No data' : 'Loading chart...' }}
  </div>
</template>

<script setup>
import { ref, watch, computed, onMounted, onUnmounted } from 'vue'
import { use } from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, LegendComponent, GridComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'

use([LineChart, TitleComponent, TooltipComponent, LegendComponent, GridComponent, CanvasRenderer])
const isDark = ref(false)
onMounted(() => {
  // Tailwind dark 클래스 감지
  isDark.value = document.documentElement.classList.contains('dark')
  
  // 클래스 변경 감시
  const observer = new MutationObserver(() => {
    isDark.value = document.documentElement.classList.contains('dark')
  })
  observer.observe(document.documentElement, { attributes: true, attributeFilter: ['class'] })
  
  onUnmounted(() => observer.disconnect())
})
const props = defineProps({
  label: {
    type: Array,
    required: true
  },
  data: {
    type: Array,
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
    type: String,
    required: true
  },
  noData: {
    type: Boolean,
    default: false
  }
})

const chartOption = ref({})

const ready = computed(() => {
  if (props.noData) return false
  const labelLen = props.label.length
  if (!labelLen) return false
  if (props.mode === '3phase') {
    return props.data.length === 3 && props.data.every(arr => arr.length === labelLen)
  } else {
    return props.data.length === 1 && props.data[0].length === labelLen
  }
})

watch([() => props.label, () => props.data, () => props.mode, isDark], ([t, v, mode]) => {
  if (!t.length || !v.length) return
  const dark = isDark.value
  const textColor = dark ? '#ffffff' : '#374151'  // 순수 흰색으로 변경
  const lineColor = dark ? '#9ca3af' : '#d1d5db'  // 더 밝은 회색
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
  textStyle: { color: textColor },
  title: { 
    text: props.title, 
    left: 'center', 
    top: 10,
    textStyle: { color: textColor }
  },
  tooltip: { trigger: 'axis' },
  legend: { 
    data: props.legend, 
    top: 30,
    textStyle: { color: textColor }
  },
  grid: { top: 60, left: 50, right: 40, bottom: 60 },
  xAxis: { 
    type: 'value', 
    name: props.mode === '3phase' ? 'Time (s)' : 'Frequency (Hz)',
    nameTextStyle: { color: textColor },
    axisLabel: { color: textColor },
    axisLine: { lineStyle: { color: lineColor } },
    splitLine: { lineStyle: { color: dark ? '#4b5563' : '#e5e7eb' } }
  },
  yAxis: { 
    type: 'value', 
    name: 'Amplitude',
    nameTextStyle: { color: textColor },
    axisLabel: { color: textColor },
    axisLine: { lineStyle: { color: lineColor } },
    splitLine: { lineStyle: { color: dark ? '#4b5563' : '#e5e7eb' } }
  },
  series,
  dataZoom: [
    { type: 'inside', throttle: 20 },
    { 
      type: 'slider', 
      height: 20, 
      bottom: 5, 
      textStyle: { color: textColor },
      borderColor: lineColor,
      fillerColor: dark ? 'rgba(100,100,100,0.3)' : 'rgba(200,200,200,0.3)'
    }
  ]
}
}, { immediate: true })
</script>