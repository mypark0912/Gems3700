<template>
    <table class="table-auto w-full dark:text-white">
    <!-- Table header -->
     <colgroup>
      <col style="width: 20%;" /> <!-- Title -->
      <col v-if="mode === 'Status' || mode === 'PowerQuality'" style="width: 10%;" /> <!-- Assembly -->
      <col style="width: 10%;" /> <!-- Value / Count -->
      <col style="width: 10%;" /> <!-- Status -->
      <col v-if="mode !== 'Event'" style="width: 50%;" /> <!-- Description -->
    </colgroup>
        <thead class="text-xs uppercase text-gray-400 dark:text-white bg-gray-50 dark:bg-gray-700 dark:bg-opacity-50 rounded-sm">
        <tr>
            <th class="p-2">
            <div class="font-bold text-left">{{ t('diagnosis.tabContext.th_title') }}</div>
            </th>
            <th v-if="mode === 'Status' || mode === 'PowerQuality'" class="p-2">
            <div class="font-bold text-left">{{ t('diagnosis.tabContext.th_assembly') }}</div>
            </th>
            <th class="p-2">
            <div class="font-bold text-left">
              <span v-if="mode === 'Status' || mode === 'PowerQuality'">{{ t('diagnosis.tabContext.th_value') }}</span>
              <span v-else>Count</span>
            </div>
            </th>
            <th class="p-2">
            <div class="font-bold text-left">{{ t('diagnosis.tabContext.th_status') }}</div>
            </th>
            <th v-if="mode !='Event'" class="p-2">
            <div class="font-bold text-left">{{ t('diagnosis.tabContext.th_description') }}</div>
            </th>
        </tr>
        </thead>
        <!-- Table body -->
        <tbody>
          <template v-for="item in items">
            <TreeRow
              v-if="mode === 'Status' || mode === 'PowerQuality'"
              :key="item.id"
              :item="item"
              :level="0"
              :mode="mode"
              :expanded="false"
            />
            <TreeRow_FaultEvent
              v-else 
              :key="item.id"
              :item="item"
              :level="0"
              :mode="mode"
              :expanded="false"
            />
          </template>
        </tbody>
    </table>
  </template>
    <script>
    import { ref } from 'vue'
    import TreeRow from './TreeRow.vue';
    import TreeRow_FaultEvent from './TreeRow_FaultEvent.vue';
    import { useI18n } from 'vue-i18n'  // ✅ 추가
    export default {
      name: 'Diagnosis_TreeTable2',
      components: {
          TreeRow,
          TreeRow_FaultEvent,
      },
      props: {
         data: {
          type: Array,
          default: () => [] // ✅ required 제거
        },
        channel:{
          type: String,
          default:''
        },
        mode:{
          type: String,
          default:''
        }
      },
      setup(props){
        const items = ref(props.data);
        const mode = ref(props.mode);
        const { t } = useI18n();

        return {
          items,
          mode,
          t,
        }
      }
    }
    </script>