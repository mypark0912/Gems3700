<template>
    <div class="overflow-x-auto min-w-[500px] p-4 rounded-lg shadow-lg">
      <div class="text-xs font-bold text-black-400 dark:text-black-500 uppercase mb-1">Current</div>
      <div class="flex items-start">
        <div v-if="devType.p <= 4" class="text-xs font-medium text-blue-700 px-1.5 bg-blue-500/20 rounded-full">L{{ devType.p == 4?'1':devType.p }}</div>
        <div v-if="devType.p <= 4" class="text-lg font-bold text-gray-800 dark:text-gray-100 mr-2">{{ datlist.irms[0].toFixed(2) }} A</div>
        <div v-if="devType.p == 4" class="text-xs font-medium text-blue-700 px-1.5 bg-blue-500/20 rounded-full">L2</div>
        <div v-if="devType.p == 4" class="text-lg font-bold text-gray-800 dark:text-gray-100 mr-2">{{ datlist.irms[1].toFixed(2) }} A</div>
        <div v-if="devType.p == 4" class="text-xs font-medium text-blue-700 px-1.5 bg-blue-500/20 rounded-full">L3</div>
        <div v-if="devType.p == 4" class="text-lg font-bold text-gray-800 dark:text-gray-100 mr-2">{{ datlist.irms[2].toFixed(2) }} A</div>
        <div v-if="devType.g == 1" class="text-xs font-medium text-blue-700 px-1.5 bg-blue-500/20 rounded-full">Ig</div>
        <div v-if="devType.g == 1" class="text-lg font-bold text-gray-800 dark:text-gray-100 mr-2">{{ datlist.ig.toFixed(2) }} mA</div>
      </div>
      <br/>
      <div class="col-span-6 flex flex-row gap-8">
        <div>
          <div class="text-xs font-bold text-black-400 dark:text-black-500 uppercase mb-1">THD</div>
          <div class="flex items-start">
            <div class="text-lg font-bold text-gray-800 dark:text-gray-100 mr-2">{{ datlist.pthd.toFixed(2) }} %</div>
          </div>
        </div>
        <div>
          <div class="text-xs font-bold text-black-400 dark:text-black-500 uppercase mb-1">Active Energy</div>
          <div class="flex items-start">
            <div class="text-lg font-bold text-gray-800 dark:text-gray-100 mr-2">{{ datlist.kwh.toFixed(2) }} kWh</div>
          </div>
        </div>
        <div>
          <div class="text-xs font-bold text-black-400 dark:text-black-500 uppercase mb-1">Temperature</div>
          <div class="flex items-start">
            <div class="text-lg font-bold text-gray-800 dark:text-gray-100 mr-2">{{ temp.toFixed(2) }} ℃</div>
          </div>
        </div>
      </div>
    </div>
  </template>
  

<script>
import { watch, ref } from 'vue';

export default {
name: 'DashboardCard07',
props: {
  data: Object, // ✅ props.data의 타입을 명시
  temp: Number,
},
setup(props){
  const datlist = ref([]);
  const temp = ref(props.temp);
  const devType = ref({});
    watch(
      () => props.data,
      (newVal) => {
        if (newVal && typeof newVal === 'object') {
          //datlist.value = convertCblistToDatlist([newVal]); // ⬅ 여기 배열로 감싸기!!
          datlist.value = newVal;
          if(datlist.value.cbtype > 5){
            devType.value = { p: datlist.value.cbtype%3 + 1, g:1}
          }else if(datlist.value.cbtype == 4 || datlist.value.cbtype == 5){
            devType.value = { p: 4, g:1}
          }else{
            devType.value = { p: datlist.value.cbtype, g:0}
          }
          console.log(devType.value);
        }
      },
      { immediate: true, deep: true }
    );

  // watchEffect(() => {
  //   if (!props.data) return;
  //   Object.assign(dat.value, props.data);
  //   console.log(dat.value)
  //   datlist.value = convertCblistToDatlist(dat.value)
  //   console.log(datlist.value);
  // });
  

  return {
    datlist,
    devType,
    //dat,
  };
}
}
</script>