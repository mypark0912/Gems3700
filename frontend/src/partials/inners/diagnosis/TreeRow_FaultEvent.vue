<template>
    <tr class="hover:bg-gray-50 dark:hover:bg-gray-800 hover:text-gray-900 dark:hover:text-white transition-colors duration-200">
    <td class="px-4 py-1">
        <div :style="{ paddingLeft: level * 20 + 'px' }" class="flex items-center text-sm" :class="{'font-bold': isParent}">
        <button
            v-if="item.children"
            @click="expanded = !expanded"
            class="mr-1 text-gray-500 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-200 transition-colors duration-200"
        >
            <component
            :is="expanded ? ChevronDownIcon : ChevronRightIcon"
            class="w-4 h-4"
            />
        </button>
        <span v-else class="w-5"></span>
        {{ itemTitle }}
        </div>
    </td>
    <td class="px-4 py-1 text-sm"> {{ item.Value }} </td>
    <td class="px-4 py-1"><div class="text-sm rounded-full px-3 py-1 min-w-[100px] w-fit text-center" :class="getStatusColor2(getStatusText(item.Status))">
      {{ getStatusCText(getStatusText(item.Status)) }}</div></td>
    <td v-if="mode !='Event'" class="px-4 py-1 text-sm"> {{ descStr }} </td>
    </tr>
    <template v-if="item.children && expanded">
        <TreeRow
        v-for="child in item.children"
        :key="child.id"
        :item="child"
        :level="level + 1"
        :expanded="false"
        />
  </template>
</template>
  
  <script setup>
  import { ref, computed, onMounted, watch } from 'vue'
  import { ChevronRightIcon, ChevronDownIcon } from '@heroicons/vue/20/solid'
  import { useI18n } from 'vue-i18n'  // ✅ 추가
  const props = defineProps({
    item: { type: Object, required: true },
    level: { type: Number, default: 0 },
    expanded : { type: Boolean, default: false},
    mode : { type: String, default:''},
  })

  const isParent = computed(() => props.item.children || props.item.isParent)

  const expanded = ref(false);
  const mode = ref('');
  const { t } = useI18n();
  const { locale } = useI18n();
  const descStr = ref('')
  const itemTitle = ref('')
  onMounted(()=>{
    expanded.value = false; //props.expanded;
    mode.value = props.mode;
  });

  watch(() => locale.value, (newLocale) => {
  if (props.item?.Descriptions?.[newLocale]) {
    descStr.value = props.item.Descriptions[newLocale];
  } else {
    descStr.value = "-";
  }
  if(props.item?.Titles?.[newLocale]){
    itemTitle.value = props.item.Titles[newLocale];
  }else{
    itemTitle.value = props.item.Title;
  }
}, { immediate: true });

  // const printDescription = (status, value) => {
  //   return status ?  value:'-';
  // }
  const getStatusColor2 = (status) => {
        switch (status) {
          case 'OK': return 'bg-green-500/20 text-green-700 font-semibold dark:bg-green-600/40 dark:text-green-300';
          case 'Low': return 'bg-yellow-500/20 text-yellow-700 font-semibold dark:bg-yellow-600/40 dark:text-yellow-300';
          case 'Medium': return 'bg-orange-500/20 text-orange-700 font-semibold dark:bg-orange-600/40 dark:text-orange-300';
          case 'High': return 'bg-red-500/20 text-red-700 font-semibold dark:bg-red-600/40 dark:text-red-300';
          default: return 'bg-gray-500/20 text-gray-700 font-semibold dark:bg-gray-600/40 dark:text-gray-300';
        }
      }; 

      const getStatusCText = (status) => {
        switch (status) {
          case 'OK': return t('diagnosis.tabContext.pqfe1');
          case 'Low': return t('diagnosis.tabContext.pqfe2');
          case 'Medium': return t('diagnosis.tabContext.pqfe3');
          case 'High': return t('diagnosis.tabContext.pqfe4');
          default: return t('diagnosis.tabContext.pqfe0');
        }
      }; 

      const getStatusText = (status) => {
        switch (status) {
          case 1: return 'OK';
          case 2: return 'Low';
          case 3: return 'Medium';
          case 4: return 'High';
          default: return 'No Data';
        }
      }


//   const formatValue = (value) => {
//     return (value == null || value == NaN || value == 'NaN') ? '-' : parseFloat(value).toFixed(2);
//   }

  // const getStatusColorFault = (status) => {
  //       switch (status) {
  //         case 'None': return 'bg-green-500/20 text-green-700 font-semibold';
  //         case 'Fault': return 'bg-red-500/20 text-red-700 font-semibold';
  //         default: return 'bg-gray-500/20 text-gray-700 font-semibold';
  //       }
  //     };

  //     const getStatusTextFault = (status) => {
  //       switch (status) {
  //         case 0: return 'None';
  //         case 1: return 'Fault';
  //         default: return 'No Data';
  //       }
  //     }

  //     const getStatusColorEvent = (status) => {
  //       switch (status) {
  //         case 'None': return 'bg-green-500/20 text-green-700 font-semibold';
  //         case 'Occured': return 'bg-red-500/20 text-red-700 font-semibold';
  //         default: return 'bg-gray-500/20 text-gray-700 font-semibold';
  //       }
  //     };

  //     const getStatusTextEvent = (status) => {
  //       switch (status) {
  //         case 0: return 'None';
  //         case 1: return 'Occured';
  //         default: return 'No Data';
  //       }
  //     }
  
  </script>