<template>
    <div class="col-span-full xl:col-span-12 bg-white dark:bg-gray-800">
        <table class="table-auto w-full dark:text-white">
            <!-- Table header -->
            <thead class="text-xs uppercase text-gray-400 dark:text-gray-200 dark:bg-gray-300 bg-gray-50 dark:bg-opacity-50 rounded-sm">
              <tr>
                <th class="p-4">
                  <div class="font-bold text-left">{{ t('report.cardContext.th_title') }}</div>
                </th>
                <th class="p-4">
                  <div class="font-bold text-left">{{ t('report.cardContext.th_status') }}</div>
                </th>
                <th class="p-4">
                  <div class="font-bold text-left">{{ t('report.cardContext.th_description') }}</div>
                </th>
              </tr>
            </thead>
            <!-- Table body -->
            <tbody class="text-sm font-medium divide-y divide-gray-100 dark:divide-gray-700/60">
              <!-- Row -->
              <tr v-for="(data, index) in printdata" :key="index">
                <td class="px-4 py-2 text-left font-semibold">{{ data.Item }}</td>
                <td class="px-4 py-2 text-left">
                  <div class="text-sm rounded-full px-3 py-1 min-w-[100px] w-fit text-center" :class="getStatusColor2(data.Status)">{{ getStatusCText(data.Status) }}</div>
                </td>
                <td class="px-4 py-2 text-left font-semibold">{{ data.Status>1? data.Description:"-" }}</td>
              </tr>
            </tbody>
          </table>
    </div>
  </template>
  
  <script>
  import { ref, watch } from 'vue'
  import { useI18n } from 'vue-i18n'  // ✅ 추가
  export default {
    name: 'DashboardCard_table4',
    props: {
      data: {
        type: Array,
        default: () => [] // ✅ required 제거
      },
      channel:{
        type: String,
        default:''
      }
    },
    setup(props){
      const { t, locale } = useI18n();
      const tbdata = ref(props.data);
      const printdata = ref([]);

      for (let i = 0; i < tbdata.value.length; i++) {
        printdata.value.push({
          Item: tbdata.value[i].Item,
          Status: getStatusText(tbdata.value[i].Status),
          Description: tbdata.value[i].Description,
          Descriptions: tbdata.value[i].Descriptions
        });
      }
      //console.log(printdata.value);
  
      const getStatusColor = (status) => {
        switch (status) {
          case 'OK': return 'text-green-600 font-semibold dark:bg-green-600/40 dark:text-green-300';
          case 'Warning': return 'text-yellow-500 font-semibold dark:bg-yellow-600/40 dark:text-yellow-300';
          case 'Inspect': return 'text-orange-500 font-semibold dark:bg-orange-600/40 dark:text-orange-300';
          case 'Repair': return 'text-red-500 font-semibold dark:bg-red-600/40 dark:text-red-300';
          default: return 'text-gray-400 font-semibold dark:bg-gray-600/40 dark:text-gray-300';
        }
      };

      const getStatusColor2 = (status) => {
        switch (status) {
          case 'OK': return 'bg-green-500/20 text-green-700 font-semibold dark:bg-green-600/40 dark:text-green-300';
          case 'Warning': return 'bg-yellow-500/20 text-yellow-700 font-semibold dark:bg-yellow-600/40 dark:text-yellow-300';
          case 'Inspect': return 'bg-orange-500/20 text-orange-700 font-semibold dark:bg-orange-600/40 dark:text-orange-300';
          case 'Repair': return 'bg-red-500/20 text-red-700 font-semibold dark:bg-red-600/40 dark:text-red-300';
          default: return 'bg-gray-500/20 text-gray-700 font-semibold dark:bg-gray-600/40 dark:text-gray-300';
        }
      };
  
  // ✅ 상태 변환 함수
      function getStatusText(status) {
        switch (status) {
          case 1: return 'OK';
          case 2: return 'Warning';
          case 3: return 'Inspect';
          case 4: return 'Repair';
          default: return 'No Data';
        }
      }

      const getStatusCText = (status) => {
        switch (status) {
          case 'OK': return t('diagnosis.tabContext.st1');
          case 'Warning': return t('diagnosis.tabContext.st2');
          case 'Inspect': return t('diagnosis.tabContext.st3');
          case 'Repair': return t('diagnosis.tabContext.st4');
          default: return t('diagnosis.tabContext.st0');
        }
      };

      watch(() => locale.value, (newLocale) => {
        printdata.value.forEach(item => {
          item.Description = item.Descriptions[newLocale];
        });
      }, { immediate: true });
  
  
      return {
        tbdata,
        printdata,
        getStatusColor,
        getStatusColor2,
        t,
        locale,
        getStatusCText,
      }
    }
  }
  </script>