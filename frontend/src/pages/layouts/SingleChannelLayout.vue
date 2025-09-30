<template>
  <div class="grid grid-cols-12 gap-6">
    <!-- 1채널 전용 레이아웃 - 카드들이 더 크게 표시 -->
    <DashboardCard_idpm300_meter 
      v-if="Object.keys(mainData).length > 0 && channelState.MainEnable" 
      :channel="channel" 
      :data="mainData" 
    />
    
    <Dashboard_idpm300_power />

    <DashboardCard_PQ_Claude 
      v-if="Object.keys(mainData).length > 0 && channelState.MainEnable" 
      :channel="channel" 
      :data="mainData" 
    />
    
    <DashboardCard_kwh v-if="Object.keys(mainData).length > 0 && channelState.MainEnable"
      :channel="channel" 
    />
    
    <DashboardCard_ipsm72_summary />
  </div>
</template>

<script>
//import DashboardCard_Meter_Claude from '../../partials/inners/dashboard/DashboardCard_Meter_Ch1.vue'
import DashboardCard_idpm300_meter from '../../partials/inners/dashboard/DashboardCard_idpm300_meter.vue'
import Dashboard_idpm300_power from '../../partials/inners/dashboard/Dashbaord_idpm300_power.vue'
import DashboardCard_PQ_Claude from '../../partials/inners/dashboard/DashboardCard_PQ_Claude2.vue'
// import Dashboard_TransInfo_final from '../../partials/inners/dashboard/Dashboard_TransInfo_final.vue'
import Dashboard_Single_Info from '../../partials/inners/dashboard/Dashboard_Single_Info.vue'
import DashboardCard_kwh from '../../partials/inners/dashboard/DashboardCard_kwh_realtime.vue'
import DashboardCard_ipsm72_summary from '../../partials/inners/dashboard/DashboardCard_ipsm72_summary.vue'
//import DashboardCard_Diagnosis from '../../partials/inners/dashboard/DashboardCard_Diagnosis.vue'
import { useAuthStore } from '@/store/auth'
import { computed } from 'vue'

export default {
  name: 'SingleChannelLayout',
  components: {
    DashboardCard_idpm300_meter,
    DashboardCard_PQ_Claude,
    Dashboard_idpm300_power,
    DashboardCard_kwh,
    DashboardCard_ipsm72_summary,
    Dashboard_Single_Info,
  },
  props: {
    mainData: {
      type: Object,
      required: true
    },
    subData: {
      type: Object,
      default: () => ({})
    },
    channelState: {
      type: Object,
      required: true
    },
    channel: {
      type: String,
      required: true
    }
  },
  setup() {
    const authStore = useAuthStore()
    const opMode = computed(() => authStore.getOpMode);
    
    return {
      opMode
    }
  }
}
</script>