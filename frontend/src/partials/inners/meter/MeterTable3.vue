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
              <div class="font-bold text-center">maxTime</div>
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
              <td class="p-2 text-center font-bold">{{ item.value }} {{ item.unit }}</td>
              <td class="p-2 text-center">{{ item.max }} <span v-if="item.max !== '-'"> {{ item.unit }}</span></td>        
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
  // ✅ `watch()`에서 props.data의 변화를 감지하고 UI 강제 갱신
  // watch(
  //   () => props.data,
  //   (newData) => {
  //     if (newData && Object.keys(newData).length > 0) {
  //       if(channel.value == 'Main')
  //          Object.assign(data.value, newData); // ✅ 기존 데이터 유지하며 새로운 데이터 적용
  //       else{
  //         Object.assign(meters.value, newData); 
  //       }
  //     }
  //   },
  //   { immediate: true } // ✅ 즉시 실행 (deep: true 제거)
  // );

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