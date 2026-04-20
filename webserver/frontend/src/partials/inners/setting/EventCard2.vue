<template>
  <div
    class="relative bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg"
  >
    <div
      class="absolute top-0 left-0 right-0 h-0.5 bg-purple-500"
      aria-hidden="true"
    ></div>
    <div
      class="px-5 pt-5 pb-6 border-b border-gray-200 dark:border-gray-700/60"
    >
      <header class="flex items-center">
        <div class="w-6 h-6 rounded-full shrink-0 bg-purple-500 mr-3">
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
          {{ t("config.eventSetup") }}
        </h3>
      </header>
    </div>
    <div class="px-4 py-3 space-y-4">
      <div class="space-y-5">

        <!-- Sag -->
        <div>
          <div class="flex items-center gap-2 mb-3">
            <div class="w-1 h-4 bg-purple-500 rounded-full"></div>
            <span class="text-xs text-gray-800 dark:text-gray-100 font-semibold uppercase tracking-wider">Sag</span>
          </div>
          <div class="grid grid-cols-4 gap-3">
            <div>
              <label class="block text-sm font-medium mb-1">Action</label>
              <select v-model.number="inputDict.eventInfo.sag_action" class="form-select w-full">
                <option :value="0">None</option>
                <option :value="1">Event</option>
                <option :value="2">Capture</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Level (%)</label>
              <input v-model="inputDict.eventInfo.sag_level" class="form-input w-full" type="number" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">HoldOff time (ms)</label>
              <input v-model="inputDict.eventInfo.sag_holdofftime" class="form-input w-full" type="number" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Hysteresis</label>
              <input v-model="inputDict.eventInfo.sag_hysteresis" class="form-input w-full" type="number" />
            </div>
          </div>
        </div>

        <!-- Swell -->
        <div>
          <div class="flex items-center gap-2 mb-3">
            <div class="w-1 h-4 bg-purple-500 rounded-full"></div>
            <span class="text-xs text-gray-800 dark:text-gray-100 font-semibold uppercase tracking-wider">Swell</span>
          </div>
          <div class="grid grid-cols-4 gap-3">
            <div>
              <label class="block text-sm font-medium mb-1">Action</label>
              <select v-model.number="inputDict.eventInfo.swell_action" class="form-select w-full">
                <option :value="0">None</option>
                <option :value="1">Event</option>
                <option :value="2">Capture</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Level (%)</label>
              <input v-model="inputDict.eventInfo.swell_level" class="form-input w-full" type="number" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">HoldOff time (ms)</label>
              <input v-model="inputDict.eventInfo.swell_holdofftime" class="form-input w-full" type="number" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Hysteresis</label>
              <input v-model="inputDict.eventInfo.swell_hysteresis" class="form-input w-full" type="number" />
            </div>
          </div>
        </div>

        <!-- Interruption -->
        <div>
          <div class="flex items-center gap-2 mb-3">
            <div class="w-1 h-4 bg-purple-500 rounded-full"></div>
            <span class="text-xs text-gray-800 dark:text-gray-100 font-semibold uppercase tracking-wider">Interruption</span>
          </div>
          <div class="grid grid-cols-4 gap-3">
            <div>
              <label class="block text-sm font-medium mb-1">Action</label>
              <select v-model.number="inputDict.eventInfo.inter_action" class="form-select w-full">
                <option :value="0">None</option>
                <option :value="1">Event</option>
                <option :value="2">Capture</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Level (%)</label>
              <input v-model="inputDict.eventInfo.inter_level" class="form-input w-full" type="number" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">HoldOff time (ms)</label>
              <input v-model="inputDict.eventInfo.inter_holdofftime" class="form-input w-full" type="number" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Hysteresis</label>
              <input v-model="inputDict.eventInfo.inter_hysteresis" class="form-input w-full" type="number" />
            </div>
          </div>
        </div>

        <!-- RVC -->
        <div>
          <div class="flex items-center gap-2 mb-3">
            <div class="w-1 h-4 bg-purple-500 rounded-full"></div>
            <span class="text-xs text-gray-800 dark:text-gray-100 font-semibold uppercase tracking-wider">RVC</span>
          </div>
          <div class="grid grid-cols-4 gap-3">
            <div>
              <label class="block text-sm font-medium mb-1">Action</label>
              <select v-model.number="inputDict.eventInfo.rvc_action" class="form-select w-full">
                <option :value="0">None</option>
                <option :value="1">Event</option>
                <option :value="2">Capture</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Level (%)</label>
              <input v-model="inputDict.eventInfo.rvc_level" class="form-input w-full" type="number" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">HoldOff time (ms)</label>
              <input v-model="inputDict.eventInfo.rvc_holdofftime" class="form-input w-full" type="number" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Hysteresis</label>
              <input v-model="inputDict.eventInfo.rvc_hysteresis" class="form-input w-full" type="number" />
            </div>
          </div>
        </div>

        <!-- MSV -->
        <div>
          <div class="flex items-center gap-2 mb-3">
            <div class="w-1 h-4 bg-purple-500 rounded-full"></div>
            <span class="text-xs text-gray-800 dark:text-gray-100 font-semibold uppercase tracking-wider">MSV</span>
          </div>
          <div class="grid grid-cols-3 gap-3">
            <div>
              <label class="block text-sm font-medium mb-1">Level</label>
              <input v-model="inputDict.eventInfo.msv_level" class="form-input w-full" type="number" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Carrier Frequency</label>
              <input v-model="inputDict.eventInfo.msv_carrierFrequency" class="form-input w-full" type="number" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Record Length</label>
              <input v-model="inputDict.eventInfo.msv_recordLength" class="form-input w-full" type="number" />
            </div>
          </div>
        </div>

        <!-- Flicker -->
        <div>
          <div class="flex items-center gap-2 mb-3">
            <div class="w-1 h-4 bg-purple-500 rounded-full"></div>
            <span class="text-xs text-gray-800 dark:text-gray-100 font-semibold uppercase tracking-wider">Flicker</span>
          </div>
          <div class="grid grid-cols-4 gap-3">
            <div>
              <label class="block text-sm font-medium mb-1">Model</label>
              <select v-model.number="inputDict.eventInfo.flicker_model" class="form-select w-full">
                <option :value="0">230V</option>
                <option :value="1">120V</option>
              </select>
            </div>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup>
import { inject } from "vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();
const inputDict = inject("channel_inputDict");
</script>
