<template>
  <ModalBasic
    id="ipsm-module-modal"
    :modalOpen="open"
    @close-modal="$emit('close')"
    title="Module Select (IPSM72)"
  >
    <div class="px-6 py-5 space-y-4 dark:text-white">
      <div class="grid grid-cols-2 gap-4">
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium">IPSM72 #1</span>
          <label class="relative inline-flex items-center cursor-pointer">
            <input type="checkbox" class="sr-only peer"
              :checked="ipsm72.ipsm72_1 === 1"
              @change="ipsm72.ipsm72_1 = $event.target.checked ? 1 : 0" />
            <div class="w-9 h-5 bg-gray-200 peer-focus:outline-none rounded-full peer dark:bg-gray-600 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-rose-500"></div>
          </label>
        </div>
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium">IPSM72 #2</span>
          <label class="relative inline-flex items-center cursor-pointer">
            <input type="checkbox" class="sr-only peer"
              :checked="ipsm72.ipsm72_2 === 1"
              @change="ipsm72.ipsm72_2 = $event.target.checked ? 1 : 0" />
            <div class="w-9 h-5 bg-gray-200 peer-focus:outline-none rounded-full peer dark:bg-gray-600 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-rose-500"></div>
          </label>
        </div>
      </div>
      <div class="grid grid-cols-2 gap-4">
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium">DI60 #1</span>
          <label class="relative inline-flex items-center cursor-pointer">
            <input type="checkbox" class="sr-only peer"
              :checked="ipsm72.di60_1 === 1"
              @change="ipsm72.di60_1 = $event.target.checked ? 1 : 0" />
            <div class="w-9 h-5 bg-gray-200 peer-focus:outline-none rounded-full peer dark:bg-gray-600 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-rose-500"></div>
          </label>
        </div>
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium">DI60 #2</span>
          <label class="relative inline-flex items-center cursor-pointer">
            <input type="checkbox" class="sr-only peer"
              :checked="ipsm72.di60_2 === 1"
              @change="ipsm72.di60_2 = $event.target.checked ? 1 : 0" />
            <div class="w-9 h-5 bg-gray-200 peer-focus:outline-none rounded-full peer dark:bg-gray-600 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-rose-500"></div>
          </label>
        </div>
      </div>

      <div class="flex justify-end pt-2">
        <button
          class="btn h-9 px-4 bg-gray-900 text-sm text-gray-100 hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white rounded-lg"
          @click="$emit('close')"
        >
          Close
        </button>
      </div>
    </div>
  </ModalBasic>
</template>

<script setup>
import { computed } from 'vue'
import ModalBasic from '../../../pages/common/ModalBasic.vue'

const props = defineProps({
  open: { type: Boolean, default: false },
  setupDict: { type: Object, required: true }
})

defineEmits(['close'])

const ipsm72 = computed(() => props.setupDict.ipsm72 ?? {})
</script>
