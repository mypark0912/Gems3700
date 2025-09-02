<template>
  <div class="relative w-full">
    <input
      type="text"
      v-model="search"
      @input="onSearch"
      @focus="isOpen = true"
      @blur="onBlur"
      class="text-xs w-full border rounded p-1 text-center"
      placeholder="Search bearing..."
    />

    <ul
      v-if="isOpen && filteredList.length"
      class="absolute bg-white border rounded shadow w-full z-50 max-h-40 overflow-auto"
    >
      <li
        v-for="(item, index) in filteredList"
        :key="index"
        @mousedown.prevent="select(item)"
        class="px-3 py-1 hover:bg-gray-100 cursor-pointer"
      >
        {{ item }}
      </li>
    </ul>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  modelValue: String,
  options: Array
})
const emit = defineEmits(['update:modelValue'])

const search = ref(props.modelValue || '')
const isOpen = ref(false)

const filteredList = computed(() =>
  props.options.filter(opt =>
    opt.toLowerCase().includes(search.value.toLowerCase())
  )
)

function select(item) {
  search.value = item
  emit('update:modelValue', item)
  isOpen.value = false
}

function onSearch() {
  isOpen.value = true
}

function onBlur() {
  setTimeout(() => {
    isOpen.value = false
  }, 150)
}
</script>
