<template>
    <tr class="hover:bg-gray-50 dark:hover:text-gray-400">
    <td class="px-4 py-1">
        <div :style="{ paddingLeft: level * 20 + 'px' }" class="flex items-center text-sm" :class="{'font-bold': isParent}">
        <button
            v-if="item.children"
            @click="expanded = !expanded"
            class="mr-1 text-gray-500 hover:text-gray-800"
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
    <td class="px-4 py-1 text-sm">{{ item.AssemblyID }}</td>
    <td class="px-4 py-1 text-sm">{{ formatValue(item.Value) }}</td>
    <td class="px-4 py-1"><div class="text-sm rounded-full px-3 py-1 min-w-[100px] w-fit text-center" :class="mode !='PowerQuality'?getStatusColor2(getStatusText(item.Status)):getStatusColorPQ(getStatusTextPQ(item.Status))">
      {{ mode !='PowerQuality'? getStatusCText(getStatusText(item.Status)):getStatusCTextPQ(getStatusTextPQ(item.Status)) }}</div></td>
    <td class="px-4 py-1 text-sm"> {{  item.children ?  descStr:'-' }}</td>
    </tr>
    <template v-if="item.children && expanded">
        <TreeRow
        v-for="child in item.children"
        :key="child.id"
        :item="child"
        :level="level + 1"
        :expanded="false"
        :mode="mode"
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
  const descStr = ref('')
  const itemTitle = ref('')
  const { locale } = useI18n();


  onMounted(()=>{
    expanded.value = false; //props.expanded;
    mode.value = props.mode;
  });

watch(() => locale.value, (newLocale) => {
  if (props.item?.Descriptions?.[newLocale]) {
    if(props.item.Status > 1)
      descStr.value = props.item.Descriptions[newLocale];
    else
      descStr.value = "-";
  } else {
    descStr.value = "-";
  }
  if(props.item?.Titles?.[newLocale]){
    itemTitle.value = props.item.Titles[newLocale];
  }else{
    itemTitle.value = props.item.Title;
  }
}, { immediate: true });

  const formatValue = (value) => {
    return (value == null || value == NaN || value == 'NaN') ? '-' : parseFloat(value).toFixed(2);
  }

  // const printDescription = (status, value, dict) => {
  //   //console.log(locale.value);

  //   descStr.value = dict[locale.value];
  //   return 
  // }

  const getStatusColor2 = (status) => {
        switch (status) {
          case 'OK': return 'bg-green-500/20 text-green-700 font-semibold dark:bg-green-600/40 dark:text-green-300';
          case 'Warning': return 'bg-yellow-500/20 text-yellow-700 font-semibold dark:bg-yellow-600/40 dark:text-yellow-300';
          case 'Inspect': return 'bg-orange-500/20 text-orange-700 font-semibold dark:bg-orange-600/40 dark:text-orange-300';
          case 'Repair': return 'bg-red-500/20 text-red-700 font-semibold dark:bg-red-600/40 dark:text-red-300';
          default: return 'bg-gray-500/20 text-gray-700 font-semibold dark:bg-gray-600/40 dark:text-gray-300';
        }
      };

      const getStatusText = (status) => {
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

      const getStatusColorPQ = (status) => {
        switch (status) {
          case 'OK': return 'bg-green-500/20 text-green-700 font-semibold dark:bg-green-600/40 dark:text-green-300';
          case 'Low': return 'bg-yellow-500/20 text-yellow-700 font-semibold dark:bg-yellow-600/40 dark:text-yellow-300';
          case 'Medium': return 'bg-orange-500/20 text-orange-700 font-semibold dark:bg-orange-600/40 dark:text-orange-300';
          case 'High': return 'bg-red-500/20 text-red-700 font-semibold dark:bg-red-600/40 dark:text-red-300';
          default: return 'bg-gray-500/20 text-gray-700 font-semibold dark:bg-gray-600/40 dark:text-gray-300';
        }
      };

      const getStatusTextPQ = (status) => {
        switch (status) {
          case 1: return 'OK';
          case 2: return 'Low';
          case 3: return 'Medium';
          case 4: return 'High';
          default: return 'No Data';
        }
      }

      const getStatusCTextPQ = (status) => {
        switch (status) {
          case 'OK': return t('diagnosis.tabContext.pqfe1');
          case 'Low': return t('diagnosis.tabContext.pqfe2');
          case 'Medium': return t('diagnosis.tabContext.pqfe3');
          case 'High': return t('diagnosis.tabContext.pqfe4');
          default: return t('diagnosis.tabContext.pqfe0');
        }
      }
  
  </script>
  