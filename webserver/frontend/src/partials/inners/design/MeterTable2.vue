<template>
    <div class="overflow-x-auto py-1">
      <table class="table-auto w-full dark:text-white">
        <!-- Table header -->
        <thead class="text-sm font-semibold text-gray-500 dark:text-gray-400 bg-gray-50/80 dark:bg-gray-700/50">
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
              <div class="font-bold text-center">{{ t('meter.Table.th_min') }}</div>
            </th>
            <th class="p-2">
              <div class="font-bold text-center">{{ t('meter.Table.th_max') }}</div>
            </th>
            <th class="p-2">
              <div class="font-bold text-center">{{ t('meter.Table.th_minTime') }}</div>
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
              <td class="p-2 text-center font-bold">{{ item.value }} {{ item.unit }}</td>
              <td class="p-2 text-center">{{ item.minTime ? `${item.min} ${item.unit}` : '-' }}</td>
              <td class="p-2 text-center">{{ item.maxTime ? `${item.max} ${item.unit}` : '-' }}</td>
              <td class="p-2 text-center">{{item.minTime ? item.minTime.split('.')[0] : '-'}} </td>
              <td class="p-2 text-center">{{item.maxTime ? item.maxTime.split('.')[0] : '-'}} </td>
            </tr>
          </template>
        </tbody>
      </table>
    </div>
  </template>


<script>
import { watchEffect, ref } from 'vue';
import { useI18n } from 'vue-i18n'
export default {
name: 'MeterTable2',
props: {
  channel: String,
  data: Object,
},
setup(props){
  const channel = ref(props.channel);
  const datlist = ref([]);
  const { t } = useI18n();

  watchEffect(() => {
    if (!props.data) return;
    Object.assign(datlist.value, props.data);
  });

  return {
    channel,
    datlist,
    t,
  };
}
}
</script>
