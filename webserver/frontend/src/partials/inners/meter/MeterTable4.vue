<template>
    <div class="overflow-x-auto py-4">
      <table class="table-auto w-full dark:text-white">
        <!-- Table header -->
        <thead class="text-xs uppercase text-gray-400 bg-gray-50 dark:bg-gray-300 dark:text-gray-200 dark:bg-opacity-50 rounded-sm">
          <tr>
            <th class="p-2">
              <div class="font-bold text-left">{{ t('meter.Table.th_tilte') }}</div>
            </th>
            <th class="p-2">
              <div class="font-bold text-left">{{ t('meter.Table.th_subtitle') }}</div>
            </th>
            <th class="p-2">
              <div class="font-bold text-center">{{ t('meter.Table.th_value') }}</div>
            </th>
            <th class="p-2">
              <div class="font-bold text-center">{{ t('meter.Table.th_max') }}</div>
            </th>
            <th class="p-2">
              <div class="font-bold text-center">{{ t('meter.Table.th_maxTime') }}</div>
            </th>
          </tr>
        </thead>
        <!-- Table body -->
        <tbody class="text-sm font-medium divide-y divide-gray-100 dark:divide-gray-700/60">
          <template v-for="(group, groupIdx) in datlist" :key="groupIdx">
            <tr v-for="(item, index) in group.data" :key="item.id ?? index">
              <!-- TITLE 은 첫 줄에서만 rowspan으로 출력 -->
              <td v-if="index === 0" :rowspan="group.data.length" class="p-2 font-bold text-left align-top">
                {{t(`meter.Table.${group.subTitle}`)   }}
              </td>
              <td class="p-2 text-left">{{ item.subTitle }}</td>
              <td class="p-2 text-center font-bold">{{ (item.value / 1000).toFixed(2) }} {{ item.unit }}</td>
              <td class="p-2 text-center">{{ (item.max / 1000).toFixed(2) }} <span v-if="item.max !== '-'"> {{ item.unit }}</span></td>        
              <td class="p-2 text-center">{{ item.maxTime.split('.')[0] }} </td>           
            </tr>
          </template>
        </tbody>
      </table>
    </div>
  </template>
  

<script>
import { watchEffect, ref } from 'vue';
import { useI18n } from 'vue-i18n'  // ✅ 추가
export default {
name: 'DashboardCard07',
props: {
  channel: String,
  data: Object, // ✅ props.data의 타입을 명시
},
setup(props){
  const channel = ref(props.channel);
  const datlist = ref([]);
  const { t } = useI18n();

  watchEffect(() => {
    if (!props.data) return;
    Object.assign(datlist.value, props.data);
    //console.log('meterTable3', datlist.value);
  });
  

  return {
    channel, // ✅ props.channel 직접 반환
    datlist,
    t,
  };
}
}
</script>