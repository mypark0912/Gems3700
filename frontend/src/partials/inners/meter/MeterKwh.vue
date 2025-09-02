<template>
    <div class="col-span-full xl:col-span-3 bg-white dark:bg-gray-800 shadow-sm rounded-xl">
     
      <div class="card">
      <div class="premium-card-header">
        <div class="header-content">
          <div class="header-left">
            <h2 class="card-title">
             {{ t(`meter.cardContext.${data[index].subTitle}`)   }} {{ title }}
            </h2>    
          </div>  
        </div>
      </div>
  
    </div>
      <!-- <header class="px-5 py-4 border-b border-gray-100 dark:border-gray-700/60">
        <h2 class="font-semibold text-gray-800 dark:text-gray-100">{{ t(`meter.cardContext.${data[index].subTitle}`)   }} {{ title }}</h2>
      </header> -->
      <div class="p-3 space-y-4">
        <div class="flex items-start ml-4">
          <!--div class="text-xs font-bold text-black-400 dark:text-black-500 uppercase mb-1">{{ data[0].subTitle }}</div-->
          <div class="text-xl font-bold text-gray-800 dark:text-gray-100 mr-2">{{data[index].data[0].value.toFixed(2) }} {{ data[0].data[0].unit }}</div>
        </div>
      </div>
    </div>
  </template>
    
    <script>
    import { watch, ref, computed } from 'vue';
    import { useI18n } from 'vue-i18n'  // ✅ 추가
    //import AccordionTableItem from '../../components/AccordionTableItem.vue'
    export default {
      name: 'MeterDetail3',
      props: {
        channel:String,
        title:String,
        data: Object, // ✅ props.data의 타입을 명시
        mode:String,
      },
      setup(props){
        const channel = ref(props.channel);
        const title = computed(() => {
          if (props.title === 'Energy') return t('meter.cardTitle.title_energy');
          //return t('meter.cardTitle.title_thd');
        });
        const index = computed(()=> props.mode == 'export'? 0 : 1);
        const data = ref([]);
        const { t } = useI18n();
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
          channel, // ✅ props.channel 직접 반환
          data,
          title,
          t,
          index,
        };
      }
    }
    </script>
    <style>
      /* CSS 파일 import */
      @import '../../../css/card-styles.css';
    </style>