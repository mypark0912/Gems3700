<template>
  <div class="grid grid-cols-12 gap-6">
    <!-- 1채널 전용 레이아웃 - 카드들이 더 크게 표시 -->
    <DashboardCard_Meter_Single :channel="channel" />
    
    <DashboardCard_PQ_Claude 
      :channel="computedChannel"  
    />
    
    <DashboardCard_kwh
      :channel="computedChannel" 
    />
    
  </div>
</template>

<script>
import DashboardCard_Meter_Single from '../../partials/inners/dashboard/DashboardCard_Meter_Single.vue'
import DashboardCard_PQ_Claude from '../../partials/inners/dashboard/DashboardCard_PQ.vue'
//import Dashboard_TransInfo from '../../partials/inners/dashboard/Dashboard_TransInfo.vue'
//import Dashboard_Single_Info from '../../partials/inners/dashboard/Dashboard_Single_Info.vue'
import Dashboard_Single from '../../partials/inners/dashboard/Dashboard_Single.vue'
import DashboardCard_kwh from '../../partials/inners/dashboard/DashboardCard_kwh_realtime.vue'
//import DashboardCard_Diagnosis from '../../partials/inners/dashboard/DashboardCard_Diagnosis.vue'
import { computed } from 'vue'

export default {
  name: 'SingleChannelLayout',
  components: {
    DashboardCard_Meter_Single,
    DashboardCard_PQ_Claude,
    DashboardCard_kwh,
    Dashboard_Single,
  },
  props: {
    channelState: {
      type: Object,
      required: true
    },
    channel: {
      type: String,
      required: true
    }
  },
  setup(props) {
    const computedChannel = computed(() => {
      console.log(props.channel);
      if (props.channel == 'Main' || props.channel == 'main')
        return 'Main'
      else
        return 'Sub'
    })
    console.log(computedChannel.value);
    //const mainData = inject('meterDictMain');
    //const subData = inject('meterDictSub');
    return {
      //opMode,
      //mainData,
      //subData,
      computedChannel,
    }
  }
}
</script>