<template>
  <div class="grid grid-cols-12 gap-x-2 gap-y-6">
    <!-- 2채널 레이아웃 - 메인과 서브 채널 모두 표시 -->
    
    <!-- 메인 채널 카드들 -->
    <Dashboard_MeterCard 
      v-if="Object.keys(mainData).length > 0 && channelState.MainEnable" 
      :channel="'main'" 
      :data="mainData" 
    />
    
    <DashboardCard04 
      v-if="channelState.MainDiagnosis" 
      :channel="'main'" 
    />

    <!-- 서브 채널 카드들 -->
    <Dashboard_MeterCard 
      v-if="Object.keys(subData).length > 0 && channelState.SubEnable" 
      :channel="'sub'" 
      :data="subData" 
    />
    
    <DashboardCard04 
      v-if="channelState.SubDiagnosis" 
      :channel="'sub'" 
    />
  </div>
</template>

<script>
import Dashboard_MeterCard from '../../partials/inners/dashboard/DashboardCard_Meter_Ch2.vue'
import DashboardCard04 from '../../partials/inners/dashboard/DashboardCard_New.vue'

export default {
  name: 'DualChannelLayout',
  components: {
    Dashboard_MeterCard,
    DashboardCard04,
  },
  props: {
    mainData: {
      type: Object,
      required: true
    },
    subData: {
      type: Object,
      required: true
    },
    channelState: {
      type: Object,
      required: true
    }
  }
}
</script>