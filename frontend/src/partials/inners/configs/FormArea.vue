<template>
  <div v-if="forms" class="space-y-2">
    <template v-for="(row, rowKey) in forms" :key="rowKey">
      <div v-if="row.lbl" class="flex items-center gap-2 mb-1">
        <label class="w-20 font-medium whitespace-nowrap">{{ row.lbl }}</label>
        <input
          v-model="row.text"
          type="text"
          class="border px-2 py-1 rounded w-28"
        />
        <button class="btn btn-primary">{{ row.btn }}</button>
      </div>
      <div v-else class="flex gap-2 mb-1">
        <button class="btn btn-green">{{ row.btn1 }}</button>
        <button class="btn btn-red">{{ row.btn2 }}</button>
      </div>
    </template>
  </div>
</template>

<script setup>
const props = defineProps({
  forms: Object
})
</script>

<style scoped>
label {
  display: inline-block;
  min-width: 70px;
}
.btn {
  min-width: 90px;
  padding: 6px 12px;
  font-size: 0.875rem;
  border-radius: 4px;
  font-weight: 500;
  text-align: center;
}
.btn-primary {
  background-color: #3b82f6;
  color: white;
}
.btn-green {
  background-color: #22c55e;
  color: white;
}
.btn-red {
  background-color: #ef4444;
  color: white;
}
</style>