<template>
    <div class="col-span-full xl:col-span-12 bg-white dark:bg-gray-800 mt-1 p-2">
        <BarChart v-if="mode == 'Status'" :data="chartData" width="450" height="360" />
        <BarChart_PQ v-else-if="mode == 'PowerQuality'" :data="chartData" width="450" height="360" />
        <BarChart_FaultEvent v-else-if="mode == 'Fault'" :data="chartData" width="450" height="360" />
        <BarChart_Event_Claude v-else :data="chartData" width="450" height="360" />
    </div>
</template>

<script>
import { ref, watch, computed } from 'vue'
import BarChart from '../../../charts/connect/BatteryCharger.vue'
import BarChart_PQ from '../../../charts/connect/BarChart_PQ_Claude.vue'
import BarChart_FaultEvent from '../../../charts/connect/BarChart_FaultEvent.vue'
import BarChart_Event_Claude from '../../../charts/connect/BarChart_Event_Claude.vue'

// Import utilities
import { tailwindConfig } from '../../../utils/Utils'

export default {
  name: 'Diagnosis_BarChart2',
  props: {
    data: {
      type: Object,
      default: () => {} 
    },
    channel: {
      type: String,
      default: ''
    },
    mode: {
      type: String,
      default: ''
    }
  },
  components: {
    BarChart,
    BarChart_PQ,
    BarChart_FaultEvent,
    BarChart_Event_Claude  // 👈 추가된 컴포넌트
  },
  setup(props) {
    const channel = ref('');
    channel.value = props.channel;
    const mode = ref(props.mode);
    
    const chartData = ref({
      labels: [],
      datasets: [
        {
          label: 'Direct',
          data: [],
          backgroundColor: tailwindConfig().theme.colors.sky[500],
          hoverBackgroundColor: tailwindConfig().theme.colors.sky[600],
          barPercentage: 0.7,
          categoryPercentage: 0.7,
          borderRadius: 4,
        },
      ],
      titles: []
    })

    const isDataValid = computed(() => {
      return props.data && 
             props.data.Names && 
             Array.isArray(props.data.Names) &&
             props.data.Values && 
             Array.isArray(props.data.Values) &&
             props.data.Names.length > 0 &&
             props.data.Values.length > 0
    })

    // ✅ 안전한 watch
    watch(() => props.data, (newData) => {
      //console.log('Props data changed:', newData) // 디버깅용
      
      if (isDataValid.value) {
        chartData.value = {
          labels: newData.Names,
          datasets: [
            {
              label: 'Direct',
              data: newData.Values,
              backgroundColor: tailwindConfig().theme.colors.sky[500],
              hoverBackgroundColor: tailwindConfig().theme.colors.sky[600],
              barPercentage: 0.7,
              categoryPercentage: 0.7,
              borderRadius: 4,
            },
          ],
          titles: newData.Titles || []
        }
        console.log('Chart data updated:', chartData.value) // 디버깅용
      }
    }, { immediate: true, deep: true })

    return {
      chartData,
      channel,
      mode,
    }    
  }
}
</script>