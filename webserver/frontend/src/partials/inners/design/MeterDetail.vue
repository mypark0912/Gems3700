<template>
    <div :class="widCss" class="meter-card">
      <div class="meter-card-header">
        <h3 class="meter-card-title" :class="accentClass">{{ cTitle }}</h3>
      </div>
      <div class="meter-card-body">
        <MeterTable v-if="mode === 0" :data="data" :channel="channel" />
        <MeterTable3 v-else-if="mode === 1"  :data="data" :channel="channel"/>
        <MeterTable4 v-else :data="data" :channel="channel"/>
      </div>
    </div>
  </template>

    <script>
    import { watch, ref, computed } from 'vue';
    import MeterTable from './MeterTable2.vue'
    import { useI18n } from 'vue-i18n'
import MeterTable3 from './MeterTable3.vue';
import MeterTable4 from './MeterTable4.vue';

    export default {
      name: 'MeterDetail',
      components:{
        MeterTable,
        MeterTable3,
        MeterTable4,
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

        const mode = computed(()=> {
          if (props.title == 'Demand'){
            return 2;
          }else if(props.title == 'Demand I'){
            return 1;
          }else{
            return 0;
          }
        });

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
          if (props.title === 'THD' || props.title === 'Power') return 'col-span-full xl:col-span-7';
          if (props.title === 'Demand' || props.title === 'Demand I') return 'col-span-full xl:col-span-5';
          return 'col-span-full xl:col-span-12';
        });

        const accentClass = computed(() => {
          if (props.title === 'Meter') return 'meter-accent-blue'
          if (props.title === 'THD') return 'meter-accent-violet'
          if (props.title === 'Power') return 'meter-accent-emerald'
          if (props.title === 'Demand') return 'meter-accent-amber'
          if (props.title === 'Demand I') return 'meter-accent-orange'
          return 'meter-accent-blue'
        })

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
          accentClass,
        };
      }
    }
    </script>

    <style scoped>
    @import '../../../css/meter-card.css';
    </style>
