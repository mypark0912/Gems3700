<template>
        <div class="flex flex-wrap items-center -m-1.5">
          <div class="m-1.5">
            <!-- Start -->
            <div class="flex flex-wrap -space-x-px">
              <button
                v-for="option in options"
                :key="option.value"
                :value="option.value"
                @click.prevent="setSelectedOption(option.value)"
                :class="[
                        'btn border px-4 py-2 transition-colors duration-200 rounded-none first:rounded-l-lg last:rounded-r-lg',
                        selectedOption === option.value
                          ? 'bg-violet-500 text-white border-violet-500'  // ✅ 선택된 버튼: 보라색 배경 + 흰색 텍스트
                          : 'bg-white text-violet-500 border-gray-200 hover:bg-gray-100 dark:bg-gray-800 dark:text-white dark:hover:bg-gray-900' // ✅ 선택되지 않은 버튼: 흰 배경 + 보라 텍스트
                      ]">
                {{ option.label }}
              </button>
            </div>
            <!-- End -->
          </div>
        </div>
        <!-- Table -->
        <div class="overflow-x-auto">
          <table class="table-auto w-full dark:text-white">
            <!-- Table header -->
            <thead class="text-xs uppercase text-gray-400 dark:text-gray-500 bg-gray-50 dark:bg-gray-700 dark:bg-opacity-50 rounded-sm">
              <tr>
                <th class="p-2">
                  <div class="font-bold text-left">Parameter</div>
                </th>
                <th class="p-2">
                  <div class="font-bold text-center">L1</div>
                </th>
                <th class="p-2">
                  <div class="font-bold text-center">L2</div>
                </th>
                <th class="p-2">
                  <div class="font-bold text-center">L3</div>
                </th>
                <th class="p-2">
                  <div class="font-bold text-center">Compliance</div>
                </th>
                <th class="p-2">
                  <div class="font-bold text-center">Required Quality</div>
                </th>
              </tr>
            </thead>
            <!-- Table body -->
            <tbody class="text-sm font-medium divide-y divide-gray-100 dark:divide-gray-700/60">
              <!-- Row -->
              <tr v-for="(row, index) in data" :key="index">
                <td class="p-2 text-left">{{ row.parameter }}</td>
                <td class="p-2 text-center">-</td>
                <td class="p-2 text-center">-</td>
                <td class="p-2 text-center">-</td>
                <td class="p-2 text-center">-</td>
                <td class="p-2 text-center">{{ row.required }}</td>
              </tr>

            </tbody>
          </table>
  
        </div>
  </template>
  
  <script>
  import { ref } from 'vue'

  export default {
    name: 'ReportTable2',
    props: {
      data: {
        type: Array,
        required: true,
        default: () => []
      }
    },
    setup(props){
      const selectedOption = ref('this_week')

      const options = ref([
        { label: "This Week", value: "this_week" },
        { label: "Last Week", value: "last_week" },
      ]);

    const setSelectedOption = (value) => {
      selectedOption.value = value;
      console.log("선택된 옵션:", value);
      getData(value); // 데이터 가져오기 실행
    };

    const getData = (value) => {
      console.log(`📊 데이터를 가져옵니다: ${value}`);
      // 여기에 API 호출 또는 데이터 로직 추가
    };

    return { selectedOption, options, setSelectedOption };
    }
  }
  </script>