<template>
  <div class="flex flex-col col-span-full sm:col-span-12 xl:col-span-12 bg-white shadow rounded-xl px-5 py-4 dark:bg-gray-700 gap-2">
    <div class="flex items-center space-x-1">
        <span class="text-sm font-semibold mr-2">{{ Item.title }}</span>
        <span
          class="text-sm rounded-full px-3 py-1 min-w-[100px] w-fit text-center whitespace-nowrap transition-all"
          :class="Item.css"
        >
          {{ Item.status }}
        </span>
      </div>
      <div>
        <template v-for="(chunk, chunkIndex) in chunkedItems" :key="chunkIndex">
          <div class="flex items-start">
            <template v-for="child in chunk" :key="child.id || child.Title">
              <span class="text-sm mr-2">{{ child.Title }} : {{ child.Value }}</span>
            </template>
          </div>
        </template>
      </div>
      <span class="text-sm font-semibold">{{ Item.description }}</span>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  data: Object,
})

const { t, locale } = useI18n()
const stData = ref(props.data);
const Item = ref({});
const childItems = ref([]);
const stText = [        
        t('dashboard.diagnosis.st0'),
        t('dashboard.diagnosis.st1'),
        t('dashboard.diagnosis.st2'),
        t('dashboard.diagnosis.st3'),
        t('dashboard.diagnosis.st4')
    ];
const css = [
    'bg-gray-500/20 text-gray-700 font-semibold',
    'bg-green-500/20 text-green-700 font-semibold',
    'bg-yellow-500/20 text-yellow-700 font-semibold',
    'bg-orange-500/20 text-orange-700 font-semibold',
    'bg-red-500/20 text-red-700 font-semibold'
]

const chunkedItems = computed(() => {
  const chunks = [];
  for (let i = 0; i < childItems.value.length; i += 5) {
    chunks.push(childItems.value.slice(i, i + 5));
  }
  return chunks;
});

watch(() => [locale.value, props.data], ([newLocale, newData]) => {
  // 데이터가 변경될 때마다 childItems 초기화
  childItems.value = [];
  
  if (newData && newData.Item) {
    Item.value = {
      title: newData.Item.Titles[newLocale],
      status: stText[newData.Item.Status],
      description: newData.Item.Descriptions[newLocale],
      css: css[newData.Item.Status]
    };

    // Child 데이터 처리 - 구조 변경 반영
    if (newData.Child && newData.Child.length > 0) {
      for (let i = 0; i < newData.Child.length; i++) {
        const child = newData.Child[i];
        // Assembly와 Title이 있는지 확인
        if (child.Assembly && child.Title) {
          const childTitle = '[' + child.Assembly + ']' + child.Title;
          const childValue = typeof child.Value === 'number' ? child.Value.toFixed(2) : child.Value;
          childItems.value.push({
            Title: childTitle,
            Value: childValue
          });
        } else {
          // 이전 구조와의 호환성을 위한 fallback
          const childTitle = child.Title || 'Unknown';
          const childValue = typeof child.Value === 'number' ? child.Value.toFixed(2) : child.Value;
          childItems.value.push({
            Title: childTitle,
            Value: childValue
          });
        }
      }
    }
  }
}, { immediate: true, deep: true });
</script>