<template>
    <div class="grow dark:text-white">
      <!-- Panel body -->
      <div class="p-6 space-y-6">
        <!-- Plans -->
        <section>
          <div class="grid grid-cols-12 gap-6">
  
            <!-- Alarm Setup -->
            <div
              class="relative col-span-full xl:col-span-12 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg"
            >
              <div
                class="absolute top-0 left-0 right-0 h-0.5 bg-stone-500"
                aria-hidden="true"
              ></div>
              <div
                class="px-5 pt-5 pb-6 border-b border-gray-200 dark:border-gray-700/60"
              >
                <header class="flex items-center mb-2">
                  <div class="w-6 h-6 rounded-full shrink-0 bg-stone-500 mr-3">

                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      class="w-6 h-6 text-white"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    >
                      <!-- XY 축 -->
                      <path d="M3 3v18h18" />
                      <!-- 트렌드 꺾은선 그래프 -->
                      <polyline points="5 15 9 10 13 13 17 7 21 11" />
                    </svg>
                  </div>
                  <h3
                    class="text-lg text-gray-800 dark:text-gray-100 font-semibold"
                  >
                    Trend Setup
                  </h3>
                </header>
              </div>
              <!-- ✅ 전체 여백 조정 -->
              <div class="flex px-4 py-3 space-x-4">
                <!-- ✅ 알람 설정 섹션 -->
                <div class="w-full xl:w-1/2">
                  <div
                    class="text-xs text-gray-800 dark:text-white font-semibold uppercase mb-4"
                  >
                    Main Channel
                  </div>
                  <div class="overflow-x-auto">
                    <table
                      class="min-w-1/2 border-collapse border border-gray-300 text-sm"
                    >
                      <thead>
                        <tr class="bg-gray-100 dark:bg-gray-700 dark:text-white">
                          <th class="border p-2 w-[5%] text-center">Index</th>
                          <th class="border p-2 w-[20%] text-center">
                            Parameter
                          </th>
                          <th class="border p-2 w-[10%]">Period</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr
                          v-for="n in 16"
                          :key="n"
                          class="hover:bg-gray-50 dark:hover:bg-gray-800"
                        >
                          <!-- Left Column: Alarm 1~16 -->
                          <td class="border p-2 text-center dark:text-white">{{ n }}</td>
                          <td class="border p-2">
                            <select
                              v-model.number="inputDict.alarm[n.toString()][0]"
                              class="form-select w-full"
                            >
                              <option
                                v-for="(param, index) in parameterOptions"
                                :key="index"
                                :value="index"
                              >
                                {{ param }}
                              </option>
                            </select>
                          </td>
                          <td class="border p-2">
                            <select
                              v-model.number="inputDict.alarm[n.toString()][1]"
                              class="form-select w-full"
                            >
                              <option value="0">&lt;</option>
                              <option value="1">&gt;</option>
                            </select>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>
                <div class="w-full xl:w-1/2">
                  <div
                    class="text-xs text-gray-800 dark:text-white font-semibold uppercase mb-4"
                  >
                    Sub Channel
                  </div>
                  <div class="overflow-x-auto">
                    <table
                      class="min-w-1/2 border-collapse border border-gray-300 text-sm"
                    >
                      <thead>
                        <tr class="bg-gray-100 dark:bg-gray-700 dark:text-white">
                          <th class="border p-2 w-[5%] text-center">Index</th>
                          <th class="border p-2 w-[20%] text-center">
                            Parameter
                          </th>
                          <th class="border p-2 w-[10%]">Period</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr
                          v-for="n in 16"
                          :key="n"
                          class="hover:bg-gray-50 dark:hover:bg-gray-800"
                        >
                          <!-- Left Column: Alarm 1~16 -->
                          <td class="border p-2 text-center dark:text-white">{{ n }}</td>
                          <td class="border p-2">
                            <select
                              v-model.number="inputDict.alarm[n.toString()][0]"
                              class="form-select w-full"
                            >
                              <option
                                v-for="(param, index) in parameterOptions"
                                :key="index"
                                :value="index"
                              >
                                {{ param }}
                              </option>
                            </select>
                          </td>
                          <td class="border p-2">
                            <select
                              v-model.number="inputDict.alarm[n.toString()][1]"
                              class="form-select w-full"
                            >
                              <option value="0">&lt;</option>
                              <option value="1">&gt;</option>
                            </select>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>
      </div>
  
      <!-- Panel footer -->
      <footer>
        <div
          class="flex px-6 py-5 border-t border-gray-200 dark:border-gray-700/60"
        >
          <!-- ✅ 왼쪽 정렬 -->
          <div class="flex gap-3">
            <button
              class="btn bg-gray-900 text-gray-100 hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white"
              @click.prevent="savefile"
            >
              Save Changes
            </button>
          </div>
        </div>
      </footer>
    </div>
  </template>
  
  <script>
  import { ref, onMounted, watch, inject, computed } from "vue";
  import axios from "axios";
  import ModalBasic from "../../components/ModalBasic.vue";
  import flatPickr from "vue-flatpickr-component";
  import "flatpickr/dist/flatpickr.css";
  import { useSetupStore } from "@/store/setup"; // ✅ Pinia Store 사용
  export default {
    name: "PlansPanel",
    props: ["channel"],
    components: {
      ModalBasic,
      flatPickr,
    },
    setup(props) {
      const setupStore = useSetupStore();
      const channel = ref(props.channel);
      const errorMessage = ref(""); // ✅ 에러 메시지 변수 추가
      const feedbackModalOpen = ref(false);
      const toggle1 = ref("YES");
      const message = ref("Select upload setting file");
      const channel_inputDict = inject("channel_inputDict");
      const parameterOptions = [
        "Temp.",
        "Freq.",
        "U1",
        "U2",
        "U3",
        "U~",
        "U12",
        "U23",
        "U31",
        "Upp~",
        "Uu",
        "Uo",
        "I1",
        "I2",
        "I3",
        "I~",
        "Itotal",
        "In",
        "P1",
        "P2",
        "P3",
        "Ptotal",
        "Q1",
        "Q2",
        "Q3",
        "Qtotal",
        "D1",
        "D2",
        "D3",
        "Dtotal",
        "S1",
        "S2",
        "S3",
        "Stotal",
        "PF1",
        "PF2",
        "PF3",
        "PFtotal"
      ];

  
      // ✅ 데이터 저장
      const save = async () => {
        inputDict.value.useFuction.ftp = showFTP.value ? 1 : 0; // ✅ Boolean → 숫자 변환
        inputDict.value.useFuction.sntp = showSNTP.value ? 1 : 0;
  
        try {
          const response = await axios.post(
            "/setting/save/Trend",
            inputDict.value,
            {
              headers: { "Content-Type": "application/json" },
              withCredentials: true,
            }
          );
  
          if (response.data && response.data["status"] === "1") {
            alert("✅ Data saved successfully!");
          } else {
            alert("❌ Data save failed!");
          }
        } catch (err) {
          console.log("Error occurred while saving:", err);
        }
      };

      const dateConfig = {
        dateFormat: "Y-m-d", // YYYY-MM-DD 포맷
        defaultDate: new Date(), // 기본값
        onChange: (selectedDates, dateStr) => {
          console.log("Selected Date:", dateStr);
        },
      };
  
      const triggerTimePicker = (e) => {
        e.target.showPicker();
      };
      return {
        inputDict,
        save,
        channel,
        errorMessage,
        feedbackModalOpen,
        message,
        toggle1,
        dateConfig,
        triggerTimePicker,
        parameterOptions,
      };
    },
  };
  </script>
  