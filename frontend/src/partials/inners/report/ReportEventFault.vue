<template>
  <div class="flex flex-col col-span-full sm:col-span-12 xl:col-span-12 bg-white shadow rounded-xl px-5 py-1 dark:bg-gray-700 gap-2">
    <span class="text-sm font-semibold p-2"> Event & Faults </span>
    <Diagnosis_TreeTable v-if="eventData.length > 0" :channel="channel" :data="eventData" :mode="'Event'" />
    
  </div>
</template>

<script setup>
import Diagnosis_TreeTable from '../diagnosis/Diagnosis_TreeTable2.vue'
import {ref} from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  data: Object,
})

const eventData = ref(props.data);
</script>