<template>
  <div class="grid grid-cols-12 gap-x-4 gap-y-4">
    <!-- 2채널 레이아웃 - 메인과 서브 채널 모두 표시 -->
    
    <DashboardCard04 
      v-if="channelState.MainDiagnosis" 
      :channel="'sub'" 
    />



    <!-- 메인 채널 카드들 -->
    <DashboardCard_Meter_Integrated 
      v-if="channelState.MainEnable" 
      :channel="'main'" 
      :data="mainData" 
    />



    <!-- 서브 채널 카드들 -->
    <DashboardCard_Meter_Integrated 
      v-if="channelState.SubEnable" 
      :channel="'sub'" 
      :data="subData" 
    />
    

  </div>
</template>

<script>
import DashboardCard_Meter_Integrated from '../../partials/inners/dashboard/DashboardCard_Meter_Integrated.vue'
import DashboardCard04 from '../../partials/inners/dashboard/DashboardCard_New2.vue'

export default {
  name: 'DualChannelLayout',
  components: {
    DashboardCard_Meter_Integrated,
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