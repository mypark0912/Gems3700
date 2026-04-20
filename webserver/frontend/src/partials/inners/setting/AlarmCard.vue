<template>
  <!-- Alarm Setup -->
  <div
    class="relative bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg"
  >
    <div
      class="absolute top-0 left-0 right-0 h-0.5 bg-yellow-500"
      aria-hidden="true"
    ></div>
    <div
      class="px-5 pt-5 pb-6 border-b border-gray-200 dark:border-gray-700/60"
    >
      <header class="flex items-center">
        <div class="flex items-center">
          <div class="w-6 h-6 rounded-full shrink-0 bg-yellow-500 mr-3">
            <svg
              class="w-6 h-6 fill-current text-white"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                d="M12 2C9.5 2 7.5 4 7.5 6.5V10c0 1.2-.4 2.4-1.2 3.3L5 14.5c-1 1 0 2.5 1.3 2.5h11.4c1.3 0 2.3-1.5 1.3-2.5l-1.3-1.2c-.8-.9-1.2-2.1-1.2-3.3V6.5C16.5 4 14.5 2 12 2z"
                stroke="currentColor"
                stroke-width="2"
                fill="none"
              />
              <path
                d="M10 19a2 2 0 0 0 4 0"
                stroke="currentColor"
                stroke-width="2"
                fill="none"
              />
            </svg>
          </div>
          <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold">
            {{t("config.channelPanel.alarmCard.title") }}
          </h3>
        </div>
      </header>
    </div>
    <!-- 전체 여백 조정 -->
    <div class="px-4 py-3 space-y-4">
      <!-- 알람 설정 섹션 -->
      <div>
        <div class="overflow-x-auto">
          <table
            class="min-w-full border-collapse border border-gray-300 text-sm"
          >
            <thead>
              <tr class="bg-gray-100 dark:bg-gray-700 dark:text-white">
                <th class="border p-2 w-[3%] text-center">
                   {{t("config.channelPanel.alarmCard.title") }}</th>
                <th class="border p-2 w-[14%] text-center">
                  {{t("config.channelPanel.alarmCard.parameter") }}</th>
                <th class="border p-2 w-[8%]">
                   {{t("config.channelPanel.alarmCard.conditions") }}</th>
                <th class="border p-2 w-[7%]">
                  {{t("config.channelPanel.alarmCard.value") }}</th>
                <th class="border p-2 w-[6%]">
                   {{t("config.channelPanel.alarmCard.hysteresis") }}</th>
                <th class="border p-2 w-[8%]">Action</th>
                <th class="border p-2 w-[6%]">Delay</th>

                <th class="border p-2 w-[3%] text-center">
                   {{t("config.channelPanel.alarmCard.title") }}</th>
                <th class="border p-2 w-[14%] text-center">
                   {{t("config.channelPanel.alarmCard.parameter") }}</th>
                <th class="border p-2 w-[8%]">
                   {{t("config.channelPanel.alarmCard.conditions") }}</th>
                <th class="border p-2 w-[7%]">
                  {{t("config.channelPanel.alarmCard.value") }}</th>
                <th class="border p-2 w-[6%]">
                    {{t("config.channelPanel.alarmCard.hysteresis") }}</th>
                <th class="border p-2 w-[8%]">Action</th>
                <th class="border p-2 w-[6%]">Delay</th>
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
                    v-model.number="inputDict.alarm[n.toString()].chan"
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
                    v-model.number="inputDict.alarm[n.toString()].cond"
                    class="form-select w-full"
                  >
                    <option value="0">Under</option>
                    <option value="1">Equal</option>
                    <option value="2">Over</option>
                  </select>
                </td>
                <td class="border p-2">
                  <input
                    v-model.number="inputDict.alarm[n.toString()].threshold"
                    class="form-input w-full"
                    type="number"
                  />
                </td>
                <td class="border p-2">
                  <input
                    v-model.number="inputDict.alarm[n.toString()].hysteresis"
                    class="form-input w-full"
                    type="number"
                  />
                </td>
                <td class="border p-2">
                  <select
                    v-model.number="inputDict.alarm[n.toString()].action"
                    class="form-select w-full"
                  >
                    <option value="0">Alarm 1</option>
                    <option value="1">Alarm 2</option>
                  </select>
                </td>
                <td class="border p-2">
                  <input
                    v-model.number="inputDict.alarm[n.toString()].delay"
                    class="form-input w-full"
                    type="number"
                  />
                </td>

                <!-- Right Column: Alarm 17~32 -->
                <td class="border p-2 text-center dark:text-white">
                  {{ n + 16 }}
                </td>
                <td class="border p-2">
                  <select
                    v-model.number="inputDict.alarm[(n + 16).toString()].chan"
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
                    v-model.number="inputDict.alarm[(n + 16).toString()].cond"
                    class="form-select w-full"
                  >
                    <option value="0">Under</option>
                    <option value="1">Equal</option>
                    <option value="2">Over</option>
                  </select>
                </td>
                <td class="border p-2">
                  <input
                    v-model.number="inputDict.alarm[(n + 16).toString()].threshold"
                    class="form-input w-full"
                    type="number"
                  />
                </td>
                <td class="border p-2">
                  <input
                    v-model.number="inputDict.alarm[(n + 16).toString()].hysteresis"
                    class="form-input w-full"
                    type="number"
                  />
                </td>
                <td class="border p-2">
                  <select
                    v-model.number="inputDict.alarm[(n + 16).toString()].action"
                    class="form-select w-full"
                  >
                    <option value="0">Alarm 1</option>
                    <option value="1">Alarm 2</option>
                  </select>
                </td>
                <td class="border p-2">
                  <input
                    v-model.number="inputDict.alarm[(n + 16).toString()].delay"
                    class="form-input w-full"
                    type="number"
                  />
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject } from "vue";
import { useI18n } from "vue-i18n";

const inputDict = inject("channel_inputDict");
const { t } = useI18n();

defineProps({
  parameterOptions: {
    type: Array,
    required: true
  }
});
</script>
