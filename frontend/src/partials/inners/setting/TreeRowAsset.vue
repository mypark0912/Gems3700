<template>
<tr class="hover:bg-gray-50 dark:hover:text-gray-400">
  <!-- 제목 셀: 고정 너비 주기 -->
  <td class="px-4 py-1 w-[200px]">  <!-- ✅ 너비 고정 -->
    <div :style="{ paddingLeft: level * 20 + 'px' }" class="flex items-center text-sm" :class="{ 'font-bold': isParent, 'text-violet-500':item.connectedCh == channel }">
      <button
        v-if="item.children"
        @click="expanded = !expanded"
        class="mr-1 text-gray-500 hover:text-gray-800"
      >
        <component :is="expanded ? ChevronDownIcon : ChevronRightIcon" class="w-4 h-4" />
      </button>
      <span v-else class="w-4"></span>
      {{ item.Title || item.Name }}
      <div v-if="item.Connected">
        <span v-if="item.connectedCh == channel" class="mr-1 text-green-500"> ● - {{ channel }}</span> 
        <span v-else> - {{ item.connectedCh }}</span>
      </div>
    </div>
  </td>

  <!-- 체크박스 셀: 가운데 정렬 -->
  <td class="px-4 py-1 text-center w-[48px]">  <!-- ✅ 가운데 정렬 + 고정폭 -->
    <template v-if="!isParent">
        <input
          type="radio"
          class="form-radio text-violet-500 focus:ring-violet-500"
          :value="item.ID"
          :checked="isChecked"
          @change="onCheckChange"
        />
      </template>
  </td>
</tr>


    <template v-if="item.children && expanded">
        <TreeRowAsset
        v-for="child in item.children"
        :key="child.ID"
        :item="child"
        :level="level + 1"
        :expanded="false"
        :checked-id="checkedId"
        :checked-name="checkedName"
        :channel="channel"
        @check-change="emitCheckChange"
        />
  </template>
</template>
  
  <script setup>
  import { ref, computed, onMounted,watch } from 'vue'
  import { ChevronRightIcon, ChevronDownIcon } from '@heroicons/vue/20/solid'
  
  const props = defineProps({
    item: { type: Object, required: true },
    level: { type: Number, default: 0 },
    expanded : { type: Boolean, default: false},
    parentChecked: { type: Boolean, default: false } ,
    checkedId: { type: String, default: "" },
    checkedName: { type: String, default: "" },
    channel : { type:String, default:''},
    expandedId: { type: String, default: '' },
  })
  const emit = defineEmits(['check-change']);
  const isChecked = computed(() => props.checkedId === props.item.ID);
  const checked = ref(false)
  const channel = ref(props.channel);
  
  watch(() => props.parentChecked, (newVal) => {
  if (newVal) {
    checked.value = false
  }
})
  const isParent = computed(() => props.item.children || props.item.isParent)

  const expanded = ref(false);

  function onCheckChange() {
    // const childrenIds = collectChildrenIds(props.item);
    // const childrenNames = collectChildrenNames(props.item);
    // checked.value = event.target.checked
    // emit('check-change', {
    //   id: props.item.ID,
    //   checked: event.target.checked,
    //   connected : props.item.Connected,
    //   name:  props.item.Name,
    //   type:  props.item.Type,
    //   children: childrenIds,
    //   childrenNames: childrenNames,
    // })
     emit('check-change', {
    id: props.item.ID,
    name: props.item.Name,
    type: props.item.Type,
    connected: props.item.Connected,
    connectedCh : props.item.connectedCh,
  });
}
// function collectChildrenIds(item) {
//   let ids = [];
//   if (item.children) {
//     item.children.forEach(child => {
//       ids.push(child.ID);
//       ids = ids.concat(collectChildrenIds(child));
//     });
//   }
//   return ids;
// }
// function collectChildrenNames(item) {
//   let names = [];
//   if (item.children) {
//     item.children.forEach(child => {
//       names.push(child.Name);
//       names = names.concat(collectChildrenNames(child));
//     });
//   }
//   return names;
// }
// 자식 컴포넌트에서 전달받은 check-change 다시 상위로 emit
function emitCheckChange(payload) {
  emit('check-change', payload)
}

  onMounted(()=>{
    //expanded.value = false; //props.expanded;
    if (props.item.isExpanded) {
      expanded.value = true;
    }
  });
  
  </script>
  