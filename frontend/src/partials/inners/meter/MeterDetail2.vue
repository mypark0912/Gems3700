<template>
    <div :class="widCss" class="bg-white dark:bg-gray-800 shadow-sm rounded-xl">
    <div class="card">
      <div class="premium-card-header">
        <div class="header-content">
          <div class="header-left">
            <h2 class="card-title">
              {{ cTitle }}
            </h2>    
          </div>  
        </div>
      </div> 
    </div>
      <div class="p-3 space-y-4">
        <MeterTable v-if="!mode" :data="data" :channel="channel" />
        <MeterTable3 v-else  :data="data" :channel="channel"/>
      </div>
    </div>
  </template>
    
    <script>
    import { watch, ref, computed } from 'vue';
    import MeterTable from './MeterTable2.vue'
    import { useI18n } from 'vue-i18n'
import MeterTable3 from './MeterTable3.vue';
    
    export default {
      name: 'MeterDetail',
      components:{
        MeterTable,
        MeterTable3,
      },
      props: {
        channel:String,
        title:String,
        data: Object,
      },
      setup(props){
        const channel = ref(props.channel);
        const data = ref([]);
        const { t } = useI18n();

        const mode = computed(()=> props.title.includes('Demand')?true:false);
        
        const cTitle = computed(() => {
          if (props.title === 'Meter') return t('meter.cardTitle.title_meter');
          if (props.title === 'Power') return t('meter.cardTitle.title_power');
          if (props.title === 'THD') return t('meter.cardTitle.title_thd');
          if (props.title === 'Demand') return t('meter.cardTitle.title_demand');
          if (props.title === 'Demand I') return t('meter.cardTitle.title_demand_i');
          return props.title;
        });

        const widCss = computed(() => {
          if (props.title === 'Meter') return 'col-span-full xl:col-span-7';
          if (props.title === 'THD' || props.title === 'Power') return 'col-span-full xl:col-span-6';
          if (props.title === 'Demand' || props.title === 'Demand I') return 'col-span-full xl:col-span-5';
          return 'col-span-full xl:col-span-12';
        });

        watch(
          () => props.data,
          (newData) => {
            if (newData && Array.isArray(newData) && newData.length > 0) {
              data.value = [...newData];            
            }
          },
          { immediate: true }
        );

        return {
          channel,
          data,
          widCss,
          t,
          cTitle,
          mode,
        };
      }
    }
    </script>

    <style>
    @import '../../../css/card-styles.css';
    </style>