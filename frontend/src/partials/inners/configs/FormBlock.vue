<template>
  <div v-if="forms" class="flex flex-col gap-2 text-xs">
    <!--h2>{{ title }} Calibration</h2-->
    <template v-for="(row, key) in forms" :key="key">
      <div class="flex items-center gap-2 whitespace-nowrap">
        <template v-if="row.lbl">
          <label class="w-24 text-xs">{{ row.lbl }}</label>
          <input v-model="row.text" type="text" class="border px-2 py-1 rounded w-28 text-xs" />
          <button v-if="row.btn0" class="btn btn-primary w-24 h-9">{{ row.btn0 }}</button>
          <button v-if="row.btn1" class="btn btn-green w-24 h-9">{{ row.btn1 }}</button>
          <button v-if="row.btn2" class="btn btn-red w-24 h-9">{{ row.btn2 }}</button>
        </template>
        <template v-else>
          <button v-if="row.btn1" class="btn btn-green w-32 h-9">{{ row.btn1 }}</button>
          <button v-if="row.btn2" class="btn btn-red w-32 h-9">{{ row.btn2 }}</button>
          <button v-if="row.btn3" class="btn btn-green w-32 h-9">{{ row.btn2 }}</button>
          <button v-if="row.btn4" class="btn btn-red w-32 h-9">{{ row.btn2 }}</button>
        </template>
      </div>
    </template>
  </div>
</template>

<script setup>
const props = defineProps({
  forms: Object,
  title:String,
})
</script>

<style scoped>
.btn {
  min-width: 80px;
  padding: 6px 12px;
  font-size: 0.75rem; /* 🔵 글자 작게 */
  border-radius: 4px;
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
