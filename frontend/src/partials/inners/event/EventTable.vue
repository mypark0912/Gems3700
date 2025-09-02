<template>
  <div class="bg-white dark:bg-gray-800 shadow-sm rounded-xl relative">
    <div>
      <!-- Table -->
      <div class="overflow-x-auto">
        <div class="flex justify-start items-center mb-4 px-2 gap-x-4 mt-2 text-sm">
          <label class="flex items-center">{{ t("event.search.title") }} :</label>
          <label class="flex items-center space-x-2">
              <input
                type="checkbox"
                class="form-checkbox text-violet-500 focus:ring-violet-500"
                v-model="period"
              />
              <span>{{ t("event.search.period") }}</span>
            </label>
            <label class="flex items-center space-x-2">
              <input
                type="checkbox"
                class="form-checkbox text-violet-500 focus:ring-violet-500"
                v-model="parameter"
              />
              <span>{{ t("event.search.parameter") }}</span>
            </label>
        </div>
        <div class="flex justify-start items-center mb-4 px-2 gap-x-4 mt-2 text-sm">
          <label class="flex items-center">
            <span class="mr-2" style="white-space: nowrap"
              >{{ t("event.date.StartDate") }}:</span
            >

            <flat-pickr
              v-model="startDate"
              :config="dateConfig"
              class="form-input w-full p-2 border border-gray-300 rounded-md text-gray-700"
            />
          </label>
          <label class="flex items-center">
            <span class="mr-2" style="white-space: nowrap"
              >{{ t("event.date.EndDate") }}:</span
            >

            <flat-pickr
              v-model="endDate"
              :config="dateConfig"
              class="form-input w-full p-2 border border-gray-300 rounded-md text-gray-700"
            />
          </label>

          <DropdownSelect v-if="options.length > 0" :options="options" />
          <button
            @click="find"
            class="btn bg-gray-900 text-gray-100 hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white"
          >
            <svg
              class="fill-current shrink-0 xs:hidden"
              width="16"
              height="16"
              viewBox="0 0 16 16"
            >
              <path
                d="M15 7H9V1c0-.6-.4-1-1-1S7 .4 7 1v6H1c-.6 0-1 .4-1 1s.4 1 1 1h6v6c0 .6.4 1 1 1s1-.4 1-1V9h6c.6 0 1-.4 1-1s-.4-1-1-1z"
              />
            </svg>
            <span class="max-xs:sr-only">Search</span>
          </button>
        </div>
        
        <!-- 테이블 컨테이너에 최소 높이 적용 -->
        <div class="table-container">
          <table
            v-if="mode == 'Event'"
            class="table-auto w-full dark:text-white divide-y divide-gray-100 dark:divide-gray-700/60"
          >
            <!-- Table header -->
            <thead
              class="text-xs uppercase text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-900/20 border-t border-gray-100 dark:border-gray-700/60"
            >
              <tr>
                <th class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
                  <div class="font-semibold text-left">
                    {{ t("event.eventTable.Type") }}
                  </div>
                </th>
                <th class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
                  <div class="font-semibold text-left">
                    {{ t("event.eventTable.StartTime") }}
                  </div>
                </th>
                <th class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
                  <div class="font-semibold text-left">
                    {{ t("event.eventTable.Duration") }}
                  </div>
                </th>
                <th class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
                  <div class="font-semibold text-left">
                    {{ t("event.eventTable.EndTime") }} 
                  </div>
                </th>
                <th class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
                  <div class="font-semibold text-left">
                    {{ t("event.eventTable.Phase") }}
                  </div>
                </th>
                <th class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
                  <div class="font-semibold text-left">
                    {{ t("event.eventTable.Level") }}
                  </div>
                </th>
              </tr>
            </thead>
            <!-- Table body -->
            <EventTableItem :channel="channel" :mode="mode" />
          </table>
          <table
            v-else
            class="table-auto w-full dark:text-white divide-y divide-gray-100 dark:divide-gray-700/60"
          >
            <!-- Table header -->
            <thead
              class="text-xs uppercase text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-900/20 border-t border-gray-100 dark:border-gray-700/60"
            >
              <tr>
                <th class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
                  <div class="font-semibold text-left">
                    {{ t("event.alarmTable.Timestamp") }}
                  </div>
                </th>
                <th class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
                  <div class="font-semibold text-left">
                    {{ t("event.alarmTable.AlarmChannel") }}
                  </div>
                </th>
                <th class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
                  <div class="font-semibold text-left">
                    {{ t("event.alarmTable.State") }}
                  </div>
                </th>
                <th class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
                  <div class="font-semibold text-left">
                    {{ t("event.alarmTable.Value") }}
                  </div>
                </th>
              </tr>
            </thead>
            <!-- Table body -->
            <EventTableItem :channel="channel" :mode="mode" />
          </table>
          
  
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { inject, ref, provide, onMounted, computed } from "vue";
import { useI18n } from "vue-i18n";
import EventTableItem from "./EventTableItem.vue";
import Datepicker from "../../../components/Datepicker.vue";
import flatPickr from "vue-flatpickr-component";
import axios  from 'axios'
import DropdownSelect from "../../../components/DropdownSelect.vue";

export default {
  name: "EventTable",
  components: {
    EventTableItem,
    DropdownSelect,
    Datepicker,
    flatPickr,
  },
  props: ["channel", "mode"],
  setup(props) {
    const { t } = useI18n();
    const sidebarOpen = ref(false);
    const channel = ref(props.channel);
    const mode = ref(props.mode);
    const options = ref([]);
    const selectedOptions = ref(0);
    const startDate = ref(new Date());
    const endDate = ref(new Date());
    const param = inject('param');
    const start = inject('start');
    const end = inject('end');
    const period = inject('period');
    const parameter = inject('parameter');
    const search = inject('search');
    const eventData = inject('eventData');
    
    // 빈 상태 표시 여부를 계산
    const showEmptyState = computed(() => {
      return eventData.value && eventData.value.length === 0;
    });
    
    onMounted(() => {
      if (mode.value == "Event") {
        options.value = [
        { id: 0, label: "All" },
          { id: 1, label: "SAG" },
          { id: 2, label: "SWELL" },
          { id: 3, label: "SHORT INTERRUPT" },
          { id: 4, label: "LONG INTERRUPT" },
          { id: 5, label: "OVER CURRENT" },
          { id: 6, label: "UNDER CURRENT" },
          { id: 7, label: "VOLTAGE TRANSIENT" },
          { id: 8, label: "CURRENT TRANSIENT" },
        ];
      } else {
        options.value = [];
      }
    });
    
    const dateConfig = {
      dateFormat: "Y-m-d",
      defaultDate: new Date(),
      onChange: (selectedDates, dateStr) => {
        console.log("Selected Date:", dateStr);
      },
    };
    
    provide("selectedOptions", selectedOptions);

    const find = () => {
      start.value = '';
      end.value='';
      param.value = -1;
      if(period.value){
        start.value = formatToFluxTime(startDate.value, 0);
        end.value = formatToFluxTime(endDate.value, 1);
        if(parameter.value){
          param.value = selectedOptions.value;
        }
      }else{
        if(parameter.value){
          param.value = selectedOptions.value;
        }
      }
      search.value = true;
    };

    const formatToFluxTime = (date, soe) => {
      if (typeof date === "string") {
        date = new Date(date);
      }
      if (!(date instanceof Date) || isNaN(date)) {
        throw new Error("Invalid date");
      }

      const pad = (num, size = 2) => String(num).padStart(size, "0");

      const year = date.getUTCFullYear();
      const month = pad(date.getUTCMonth() + 1);
      const day = pad(date.getUTCDate());
      let hours, minutes, seconds, milliseconds;

      if (soe === 0) {
        hours = "00";
        minutes = "00";
        seconds = "00";
        milliseconds = "000";
      } else {
        hours = "23";
        minutes = "59";
        seconds = "59";
        milliseconds = "999";
      }

      return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}.${milliseconds}Z`;
    };

    const fetchParamData = async() => {
      try {
        const response = await axios.get(`/api/getAlarmParms/${channel.value}`);
        if(response.data.success){
          const result = response.data.data;
          options.value.push({ id: 0, label: "All" });
          //console.log(result);
          for (let i = 0 ; i < result.length; i++){
            options.value.push({id:result[i].id,label: result[i].label })
          }
        }else{
          //alert('No Data');
        }
      } catch (error) {
        console.error("데이터를 가져오는 중 오류 발생:", error);
      }
    };

    onMounted(async()=>{
      if(mode.value != 'Event'){
          await fetchParamData();
      }
    });

    return {
      t,
      sidebarOpen,
      options,
      fetchParamData,
      channel,
      mode,
      find,
      dateConfig,
      startDate,
      endDate,
      period,
      parameter,
      start,
      param,
      end,
      search,
      showEmptyState,
    };
  },
};
</script>

<style scoped>
/* 전체 래퍼에 최소 높이 설정 */
.bg-white.dark\:bg-gray-800 {
  min-height: 700px;
  display: flex;
  flex-direction: column;
}

/* 내부 div도 flex로 설정하여 공간 채우기 */
.bg-white.dark\:bg-gray-800 > div {
  flex: 1;
  display: flex;
  flex-direction: column;
}

/* overflow-x-auto 영역도 flex로 설정 */
.overflow-x-auto {
  flex: 1;
  display: flex;
  flex-direction: column;
}

/* 테이블 컨테이너 - 최소 높이 제거, flex 사용 */
.table-container {
  /* min-height 제거 */
  position: relative;
  display: flex;
  flex-direction: column;
  flex: 1;
}

/* 테이블 자체 - 최소 높이 제거 */
.table-container table {
  /* min-height 제거 */
  width: 100%;
}

/* 테이블이 늘어나지 않도록 설정 */
.table-container table tbody {
  height: auto;
}

/* 빈 상태 메시지 스타일 */
.empty-state {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  text-align: center;
  padding: 2rem;
}

/* 다크 모드에서도 잘 보이도록 */
.dark .empty-state {
  color: #9ca3af;
}
</style>
