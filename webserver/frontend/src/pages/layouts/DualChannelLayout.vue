<template>
  <div class="grid grid-cols-12 gap-x-2 gap-y-6">
    <!-- 2채널 레이아웃 - 메인과 서브 채널 모두 표시 -->
    
    <!-- 메인 채널 카드들 -->
    <Dashboard_MeterCard 
      :channel="'main'" 
      
    />
    
    <!-- 서브 채널 카드들 -->
    <Dashboard_MeterCard 
      :channel="'sub'" 
    
    />
    
  </div>
</template>

<script>
import Dashboard_MeterCard from '../../partials/inners/dashboard/DashboardCard_Meter_Ch2.vue'
//import DashboardCard04 from '../../partials/inners/dashboard/DashboardCard_New.vue'

export default {
  name: 'DualChannelLayout',
  components: {
    Dashboard_MeterCard,
  },
  props: {
    channelState: {
      type: Object,
      required: true
    }
  }
}
</script>