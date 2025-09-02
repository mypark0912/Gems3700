<template>
  <tr class="hover:bg-gray-50 dark:hover:text-gray-400">
    <!-- 첫 번째 컬럼 (70%) - Item -->
    <td class="px-4 py-1 overflow-hidden">
      <div 
        :style="{ paddingLeft: level * 20 + 'px' }" 
        class="flex items-center text-sm truncate" 
        :class="{'font-bold': isParent}"
      >
        <button
          v-if="item.children"
          @click="expanded = !expanded"
          class="mr-1 text-gray-500 hover:text-gray-800 flex-shrink-0"
        >
          <component
            :is="expanded ? ChevronDownIcon : ChevronRightIcon"
            class="w-4 h-4"
          />
        </button>
        <span v-else class="w-5 flex-shrink-0"></span>
        <span class="truncate">{{ item.Title }}</span>
      </div>
    </td>
    
    <!-- 두 번째 컬럼 (20%) - Assembly -->
    <td class="px-4 py-1 overflow-hidden">
      <span class="text-sm truncate block">{{ item.AssemblyID }}</span>
    </td>
    
    <!-- 세 번째 컬럼 (10%) - Select -->
    <td class="px-4 py-1 text-center">
      <input
        type="checkbox"
        class="form-checkbox text-violet-500 focus:ring-violet-500"
        :value="item.ID" 
        :checked="isChecked"
        @change="onCheckChange"
      />
    </td>
  </tr>
  
  <!-- 자식 요소들 -->
  <template v-if="item.children && expanded">
    <TreeRowTrend
      v-for="child in item.children"
      :key="child.id"
      :item="child"
      :level="level + 1"
      :expanded="false"
      :checked-ids="checkedIds"
      :checked-names="checkedNames"
      @check-change="emitCheckChange"
    />
  </template>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { ChevronRightIcon, ChevronDownIcon } from '@heroicons/vue/20/solid'

const props = defineProps({
  item: { type: Object, required: true },
  level: { type: Number, default: 0 },
  expanded: { type: Boolean, default: false},
  parentChecked: { type: Boolean, default: false },
  checkedIds: { type: Array, default: () => [] },
  checkedNames: { type: Array, default: () => [] }
})

const emit = defineEmits(['check-change']);
const isChecked = computed(() => props.checkedIds.includes(props.item.ID));
const checked = ref(false)

watch(() => props.parentChecked, (newVal) => {
  if (newVal) {
    checked.value = false
  }
})

const isParent = computed(() => props.item.children || props.item.isParent)
const expanded = ref(false);

function onCheckChange(event) {
  const childrenIds = collectChildrenIds(props.item);
  const childrenNames = collectChildrenNames(props.item);
  checked.value = event.target.checked
  emit('check-change', {
    id: props.item.ID,
    checked: event.target.checked,
    name: props.item.Name,
    children: childrenIds,
    childrenNames: childrenNames,
  })
}

function collectChildrenIds(item) {
  let ids = [];
  if (item.children) {
    item.children.forEach(child => {
      ids.push(child.ID);
      ids = ids.concat(collectChildrenIds(child));
    });
  }
  return ids;
}

function collectChildrenNames(item) {
  let names = [];
  if (item.children) {
    item.children.forEach(child => {
      names.push(child.Name);
      names = names.concat(collectChildrenNames(child));
    });
  }
  return names;
}

function emitCheckChange(payload) {
  emit('check-change', payload)
}

onMounted(() => {
  expanded.value = false;
});
</script>