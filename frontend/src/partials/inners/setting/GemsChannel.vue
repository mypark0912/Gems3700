<template>
    <div class="grow dark:text-white">
      <!-- Panel body -->
      <div class="p-6 space-y-6">
        <!-- Plans -->
        <section>
          <div class="mb-4">
            <div class="flex items-center space-x-4 text-sm">
              <!-- Enable Checkbox -->
              <label v-if="false" class="flex items-center space-x-2">
                <input
                  type="checkbox"
                  class="form-checkbox text-violet-500 focus:ring-violet-500"
                  :checked="getInputDict().Enable === 1"
                  @change="updateField('Enable', $event.target.checked ? 1 : 0)"
                />
                <span class="text-sm font-normal">
                  {{ t("config.channelPanel.selEnable") }}
                </span>
              </label>
  
              <!-- Power Quality Checkbox -->
              <label v-if="false" class="flex items-center space-x-2">
                <input
                  type="checkbox"
                  class="form-checkbox text-violet-500 focus:ring-violet-500"
                  :checked="getInputDict().PowerQuality === 1"
                  @change="
                    updateField('PowerQuality', $event.target.checked ? 1 : 0)
                  "
                />
                <span class="text-sm font-normal">
                  {{ t("config.channelPanel.selPQ") }}
                </span>
              </label>
            </div>
          </div>
  
          <!--CT,PT,Demand,Samlping -->
          <div>
            <div v-if="getInputDict().Enable" class="grid grid-cols-12 gap-6">
              <!-- CT -->
              <div
                class="relative col-span-full xl:col-span-3 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg flex flex-col flex-1"
              >
                <div
                  class="absolute top-0 left-0 right-0 h-0.5 bg-green-500"
                  aria-hidden="true"
                ></div>
                <div
                  class="px-5 pt-5 pb-6 border-b border-gray-200 dark:border-gray-700/60"
                >
                  <header class="flex items-center mb-2">
                    <div class="w-6 h-6 rounded-full shrink-0 bg-green-500 mr-3">
                      <svg
                        class="w-6 h-6 fill-current text-white"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                      >
                        <path
                          d="M3 12c1.5-2 3-2 4-2s2 .5 3 2 3 2 4 2 2-.5 3-2 3-2 4-2 2 .5 3 2"
                          stroke="currentColor"
                          stroke-width="2"
                          fill="none"
                          stroke-linecap="round"
                          stroke-linejoin="round"
                        />
                        <circle cx="3" cy="12" r="1.5" fill="currentColor" />
                        <circle cx="21" cy="12" r="1.5" fill="currentColor" />
                        <path
                          d="M5 12v7m14-7v7"
                          stroke="currentColor"
                          stroke-width="2"
                          fill="none"
                          stroke-linecap="round"
                        />
                      </svg>
                    </div>
                    <h3
                      class="text-lg text-gray-800 dark:text-gray-100 font-semibold"
                    >
                      CT
                    </h3>
                  </header>
                </div>
                <div class="px-4 py-4 space-y-3">
                  <div>
                    <label
                      class="block text-sm font-medium mb-2"
                      for="device-model"
                      >CT Directions
                    </label>
                    <div class="flex items-center w-full">
                      <label
                        class="flex items-center justify-start w-3/12 text-sm font-medium text-gray-700 dark:text-gray-300"
                      >
                        CT1
                      </label>
                      <div class="flex w-9/12">
                        <button
                          v-for="(option, index) in ct_direction[0]"
                          :key="option.value"
                          :value="option.value"
                          @click.prevent="setct_direction(option.value, 0)"
                          :class="[
                            'btn-direction border px-4 py-2 transition-colors duration-200 flex-1',
                            index === 0 ? 'rounded-l-lg' : 'rounded-none',
                            index === ct_direction[0].length - 1
                              ? 'rounded-r-lg'
                              : 'rounded-none',
                            getInputDict().ctInfo.direction[0] === option.value
                              ? 'bg-violet-500 text-white border-violet-500'
                              : 'bg-white text-violet-500 border-gray-200 hover:bg-gray-100 dark:bg-gray-800 dark:text-white dark:hover:bg-gray-900',
                          ]"
                        >
                          {{ option.label }}
                        </button>
                      </div>
                    </div>
                  </div>
                  <div>
                    <div class="flex items-center w-full">
                      <label
                        class="flex items-center justify-start w-3/12 text-sm font-medium text-gray-700 dark:text-gray-300"
                      >
                        CT2
                      </label>
                      <div class="flex w-9/12">
                        <button
                          v-for="(option, index) in ct_direction[1]"
                          :key="option.value"
                          :value="option.value"
                          @click.prevent="setct_direction(option.value, 1)"
                          :class="[
                            'btn-direction border px-4 py-2 transition-colors duration-200 flex-1',
                            index === 0 ? 'rounded-l-lg' : 'rounded-none',
                            index === ct_direction[1].length - 1
                              ? 'rounded-r-lg'
                              : 'rounded-none',
                            getInputDict().ctInfo.direction[1] === option.value
                              ? 'bg-violet-500 text-white border-violet-500'
                              : 'bg-white text-violet-500 border-gray-200 hover:bg-gray-100 dark:bg-gray-800 dark:text-white dark:hover:bg-gray-900',
                          ]"
                        >
                          {{ option.label }}
                        </button>
                      </div>
                    </div>
                  </div>
  
                  <div>
                    <div class="flex items-center w-full">
                      <label
                        class="flex items-center justify-start w-3/12 text-sm font-medium text-gray-700 dark:text-gray-300"
                      >
                        CT3
                      </label>
                      <div class="flex w-9/12">
                        <button
                          v-for="(option, index) in ct_direction[2]"
                          :key="option.value"
                          :value="option.value"
                          @click.prevent="setct_direction(option.value, 2)"
                          :class="[
                            'btn-direction border px-4 py-2 transition-colors duration-200 flex-1',
                            index === 0 ? 'rounded-l-lg' : 'rounded-none',
                            index === ct_direction[2].length - 1
                              ? 'rounded-r-lg'
                              : 'rounded-none',
                            getInputDict().ctInfo.direction[2] === option.value
                              ? 'bg-violet-500 text-white border-violet-500'
                              : 'bg-white text-violet-500 border-gray-200 hover:bg-gray-100 dark:bg-gray-800 dark:text-white dark:hover:bg-gray-900',
                          ]"
                        >
                          {{ option.label }}
                        </button>
                      </div>
                    </div>
                  </div>
                  <div class="grid grid-cols-2 gap-4">
                    <div>
                      <label
                        class="block text-sm font-medium mb-2"
                        for="serial-number"
                        >Starting Current (mA)</label
                      >
                      <input
                        :value="getInputDict().ctInfo.startingcurrent"
                        @input="
                          updateNestedField(
                            'ctInfo',
                            'startingcurrent',
                            $event.target.value
                          )
                        "
                        id="serial-number"
                        class="form-input w-full"
                        type="text"
                      />
                    </div>
  
                    <div>
                      <label
                        class="block text-sm font-medium mb-2"
                        for="mac-address"
                        >Rated Current (A)</label
                      >
                      <input
                        :value="getInputDict().ctInfo.inorminal"
                        @input="
                          updateNestedField(
                            'ctInfo',
                            'inorminal',
                            $event.target.value
                          )
                        "
                        id="mac-address"
                        class="form-input w-full"
                        type="text"
                      />
                    </div>
                  </div>
  
                  <div class="grid grid-cols-2 gap-4">
                    <div>
                      <label
                        class="block text-sm font-medium mb-2"
                        for="fw-version"
                        >Primary</label
                      >
                      <input
                        :value="getInputDict().ctInfo.ct1"
                        @input="
                          updateNestedField('ctInfo', 'ct1', $event.target.value)
                        "
                        id="fw-version"
                        class="form-input w-full"
                        type="text"
                      />
                    </div>
  
                    <div>
                      <label
                        class="block text-sm font-medium mb-2"
                        for="build-date"
                        >Secondary</label
                      >
                      <select
                        :value="getInputDict().ctInfo.ct2"
                        @change="
                          updateNestedField('ctInfo', 'ct2', $event.target.value)
                        "
                        id="use-dhcp"
                        class="form-select w-full bg-white border-gray-300 rounded-md shadow-sm focus:ring-violet-500 focus:border-violet-500"
                      >
                        <option :value="0">5A</option>
                        <option :value="1">100mA</option>
                        <option :value="2">333mV</option>
                        <option :value="3">Rogowski</option>
                      </select>
                    </div>
                  </div>
  
                  <div class="grid grid-cols-2 gap-4">
                    <div>
                      <label
                        class="block text-sm font-medium mb-2"
                        for="mac-address"
                        >ZCT Scale</label
                      >
                      <input
                        :value="getInputDict().ctInfo.zctscale"
                        @input="
                          updateNestedField(
                            'ctInfo',
                            'zctscale',
                            $event.target.value
                          )
                        "
                        id="mac-address"
                        class="form-input w-full"
                        type="text"
                      />
                    </div>
  
                    <div>
                      <label
                        class="block text-sm font-medium mb-2"
                        for="mac-address"
                        >ZCT Type</label
                      >
                      <select
                        :value="getInputDict().ctInfo.zcttpye"
                        @change="
                          updateNestedField(
                            'ctInfo',
                            'zcttpye',
                            $event.target.value
                          )
                        "
                        id="use-dhcp"
                        class="form-select w-full bg-white border-gray-300 rounded-md shadow-sm focus:ring-violet-500 focus:border-violet-500"
                      >
                        <option :value="0">A_Type</option>
                        <option :value="1">B_Type</option>
                        <option :value="2">C_Type</option>
                      </select>
                    </div>
                  </div>
                </div>
              </div>
              <!-- PT -->
              <div v-if="getInputDict().PowerQuality == 1 || getInputDict().PowerQuality == true"
                class="relative col-span-full xl:col-span-3 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg"
              >
                <div
                  class="absolute top-0 left-0 right-0 h-0.5 bg-sky-500"
                  aria-hidden="true"
                ></div>
                <div
                  class="px-5 pt-5 pb-6 border-b border-gray-200 dark:border-gray-700/60"
                >
                  <header class="flex items-center mb-2">
                    <div class="w-6 h-6 rounded-full shrink-0 bg-sky-500 mr-3">
                      <svg
                        class="w-6 h-6 fill-current text-white"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                      >
                        <path
                          d="M3 12c1.5-2 3-2 4-2s2 .5 3 2 3 2 4 2 2-.5 3-2 3-2 4-2 2 .5 3 2"
                          stroke="currentColor"
                          stroke-width="2"
                          fill="none"
                          stroke-linecap="round"
                          stroke-linejoin="round"
                        />
                        <circle cx="3" cy="12" r="1.5" fill="currentColor" />
                        <circle cx="21" cy="12" r="1.5" fill="currentColor" />
                        <path
                          d="M5 12v7m14-7v7"
                          stroke="currentColor"
                          stroke-width="2"
                          fill="none"
                          stroke-linecap="round"
                        />
                      </svg>
                    </div>
                    <h3
                      class="text-lg text-gray-800 dark:text-gray-100 font-semibold"
                    >
                      PT
                    </h3>
                  </header>
                </div>
                <div class="px-4 py-4 space-y-3">
                  <div>
                    <label
                      class="block text-sm font-medium mb-2"
                      for="serial-number"
                      >Line Frequency</label
                    >
                    <div class="flex w-full">
                      <button
                        v-for="(option, index) in options"
                        :key="option.value"
                        :value="option.value"
                        @click.prevent="setbtnOption(option.value)"
                        :class="[
                          'btn border px-4 py-2 transition-colors duration-200 flex-1',
                          index === 0 ? 'rounded-l-lg' : 'rounded-none',
                          index === options.length - 1
                            ? 'rounded-r-lg'
                            : 'rounded-none',
                          getInputDict().ptInfo.linefrequency === option.value
                            ? 'bg-violet-500 text-white border-violet-500'
                            : 'bg-white text-violet-500 border-gray-200 hover:bg-gray-100 dark:bg-gray-800 dark:text-white dark:hover:bg-gray-900',
                        ]"
                      >
                        {{ option.label }}
                      </button>
                    </div>
                  </div>
  
                  <div class="grid grid-cols-2 gap-4">
                    <div>
                      <label
                        class="block text-sm font-medium mb-2"
                        for="device-model"
                        >Wiring Mode</label
                      >
                      <select
                        :value="getInputDict().ptInfo.wiringmode"
                        @change="
                          updateNestedField(
                            'ptInfo',
                            'wiringmode',
                            $event.target.value
                          )
                        "
                        id="ctdirection"
                        class="form-select w-full bg-white border-gray-300 rounded-md shadow-sm focus:ring-violet-500 focus:border-violet-500"
                      >
                        <option :value="0">3P4W</option>
                        <option :value="1">3P3W</option>
                      </select>
                    </div>
                    <div>
                      <label
                        class="block text-sm font-medium mb-2"
                        for="mac-address"
                        >Rated Voltage (V)</label
                      >
                      <input
                        :value="getInputDict().ptInfo.vnorminal"
                        @input="
                          updateNestedField(
                            'ptInfo',
                            'vnorminal',
                            $event.target.value
                          )
                        "
                        id="mac-address"
                        class="form-input w-full"
                        type="text"
                      />
                    </div>
                  </div>
  
                  <div class="grid grid-cols-2 gap-4">
                    <div>
                      <label
                        class="block text-sm font-medium mb-2"
                        for="fw-version"
                        >Primary</label
                      >
                      <input
                        :value="getInputDict().ptInfo.pt1"
                        @input="
                          updateNestedField('ptInfo', 'pt1', $event.target.value)
                        "
                        id="fw-version"
                        class="form-input w-full"
                        type="text"
                      />
                    </div>
  
                    <div>
                      <label
                        class="block text-sm font-medium mb-2"
                        for="fw-version"
                        >Secondary</label
                      >
                      <input
                        :value="getInputDict().ptInfo.pt2"
                        @input="
                          updateNestedField('ptInfo', 'pt2', $event.target.value)
                        "
                        id="fw-version"
                        class="form-input w-full"
                        type="text"
                      />
                    </div>
                  </div>
                </div>
              </div>
              <!-- Sampling -->
              <div
                class="relative col-span-full xl:col-span-3 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg flex flex-col flex-1"
              >
                <div
                  class="absolute top-0 left-0 right-0 h-0.5 bg-blue-500"
                  aria-hidden="true"
                ></div>
                <div
                  class="px-5 pt-5 pb-6 border-b border-gray-200 dark:border-gray-700/60"
                >
                  <header class="flex items-center mb-2">
                    <div class="w-6 h-6 rounded-full shrink-0 bg-blue-500 mr-3">
                      <svg
                        class="w-6 h-6 fill-current text-white"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                      >
                        <path
                          d="M3 20h18M5 16l4-4 4 4 6-6"
                          stroke="currentColor"
                          stroke-width="2"
                          stroke-linecap="round"
                          stroke-linejoin="round"
                        />
                      </svg>
                    </div>
                    <h3
                      class="text-lg text-gray-800 dark:text-gray-100 font-semibold"
                    >
                      Sampling
                    </h3>
                  </header>
                </div>
                <div class="px-4 py-4 space-y-3">
                  <div>
                    <label class="block text-sm font-medium mb-1.5"
                      >Sampling Rate</label
                    >
                    <select
                      :value="getInputDict().sampling.rate"
                      @change="
                        updateNestedField(
                          'sampling',
                          'rate',
                          Number($event.target.value)
                        )
                      "
                      class="form-select w-full"
                    >
                      <option :value="8000">8 KHz</option>
                      <option :value="4000">4 KHz</option>
                      <option :value="2000">2 KHz</option>
                      <option :value="1000">1 KHz</option>
                    </select>
                  </div>
                  <div>
                    <label class="block text-sm font-medium mb-1.5"
                      >Sample Duration</label
                    >
                    <select
                      :value="getInputDict().sampling.duration"
                      @change="
                        updateNestedField(
                          'sampling',
                          'duration',
                          Number($event.target.value)
                        )
                      "
                      class="form-select w-full"
                    >
                      <option
                        v-for="(sec, index) in [5, 10, 15, 20, 25, 30]"
                        :key="index"
                        :value="sec"
                      >
                        {{ sec }} sec
                      </option>
                    </select>
                  </div>
                  <div>
                    <label class="block text-sm font-medium mb-1.5"
                      >Interval</label
                    >
                    <select
                      :value="getInputDict().sampling.period"
                      @change="
                        updateNestedField(
                          'sampling',
                          'period',
                          Number($event.target.value)
                        )
                      "
                      class="form-select w-full"
                    >
                      <option
                        v-for="(min, index) in [
                          1, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60,
                        ]"
                        :key="index"
                        :value="min"
                      >
                        {{ min }} {{ t("config.plansPanel.demand.minutes")
                        }}
                      </option>
                    </select>
                  </div>
                </div>
              </div>
              <!-- Demand -->
              <div
                class="relative col-span-full xl:col-span-3 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg flex flex-col flex-1"
              >
                <div
                  class="absolute top-0 left-0 right-0 h-0.5 bg-blue-500"
                  aria-hidden="true"
                ></div>
                <div
                  class="px-5 pt-5 pb-6 border-b border-gray-200 dark:border-gray-700/60"
                >
                  <header class="flex items-center mb-2">
                    <div class="w-6 h-6 rounded-full shrink-0 bg-blue-500 mr-3">
                      <svg
                        class="w-6 h-6 fill-current text-white"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                      >
                        <path
                          d="M3 20h18M5 16l4-4 4 4 6-6"
                          stroke="currentColor"
                          stroke-width="2"
                          stroke-linecap="round"
                          stroke-linejoin="round"
                        />
                      </svg>
                    </div>
                    <h3
                      class="text-lg text-gray-800 dark:text-gray-100 font-semibold"
                    >
                      Demand
                    </h3>
                  </header>
                </div>
                <div class="px-4 py-4 space-y-3">
                  <div>
                    <label class="block text-sm font-medium mb-1.5" for="rtu-id"
                      >Target (W)</label
                    >
                    <input
                      :value="getInputDict().demand.target"
                      @input="
                        updateNestedField(
                          'demand',
                          'target',
                          Number($event.target.value)
                        )
                      "
                      class="form-input w-full"
                      type="number"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium mb-1.5" for="model"
                      >{{ t("config.plansPanel.demand.interval")
                      }}<!--Interval--></label
                    >
                    <select
                      :value="getInputDict().demand.demand_interval"
                      @change="
                        updateNestedField(
                          'demand',
                          'demand_interval',
                          Number($event.target.value)
                        )
                      "
                      class="form-select w-full bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                    >
                      <option :value="1">
                        1 {{ t("config.plansPanel.demand.minutes")
                        }}<!--Minutes-->
                      </option>
                      <option :value="5">
                        5 {{ t("config.plansPanel.demand.minutes")
                        }}<!--Minutes-->
                      </option>
                      <option :value="15">
                        15 {{ t("config.plansPanel.demand.minutes")
                        }}<!--Minutes-->
                      </option>
                      <option :value="30">
                        30 {{ t("config.plansPanel.demand.minutes")
                        }}<!--Minutes-->
                      </option>
                      <option :value="60">
                        1 {{ t("config.plansPanel.demand.hours")
                        }}<!--Hours-->
                      </option>
                    </select>
                  </div>
                </div>
              </div>
              <TrendPanelCard v-if="devMode != 'device0'" />
              <AlarmCard
                v-if="devMode != 'device0'"
                :parameterOptions="parameterOptions"
              />
              <EventCard1 v-if="getInputDict().PowerQuality == 1 || getInputDict().PowerQuality == true" />
              <EventCard2 v-if="getInputDict().PowerQuality == 1 || getInputDict().PowerQuality == true"
                :title="'Event Setup 2'"
                :option1="'Sag'"
                :option2="'Swell'"
              />
              <EventCard2 v-if="getInputDict().PowerQuality == 1 || getInputDict().PowerQuality == true"
                :title="'Event Setup 3'"
                :option1="'Over Current'"
                :option2="'Interruption'"
              />
              <TrendPanelCard v-if="devMode != 'device0' && (getInputDict().PowerQuality == 1 || getInputDict().PowerQuality == true)" />
              <AssetCard
                v-if="
                  channel == 'Main' ? currentDiagnosis.Main : currentDiagnosis.Sub
                "
                :isAsset="isAssetExist"
                :channel="channel"
                :savefile="savefile"
              />
              <NameplateCard
                v-if="
                  (channel == 'Main'
                    ? currentDiagnosis.Main
                    : currentDiagnosis.Sub) && tableData.length > 0
                "
                :key="currentDiagnosis + '-' + channel"
              />
              <ParamCard
                v-if="
                  (channel == 'Main'
                    ? currentDiagnosis.Main
                    : currentDiagnosis.Sub) && paramData.length > 0
                "
                :key="currentDiagnosis + '-' + channel"
              />
            </div>
          </div>
        </section>
      </div>
  
      <!-- Panel footer -->
      <footer v-if="false">
        <div
          class="flex px-6 py-5 border-t border-gray-200 dark:border-gray-700/60"
        >
          <div class="flex gap-3">
            <button
              class="btn bg-gray-900 text-gray-100 hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white"
              @click.prevent="savefile"
            >
              {{ t("config.save")
              }}<!--Save Changes-->
            </button>
            <button
              v-if="isCurrentChannelDiagnosisActive"
              class="btn bg-emerald-900 text-emerald-100 hover:bg-emerald-800 dark:bg-emerald-100 dark:text-emerald-800 dark:hover:bg-white"
              @click.prevent="testDiagnosis"
            >
              {{ t("config.testSet")
              }}<!--Test Diagnosis setting-->
            </button>
            <button
              :disabled="!isRestartButtonEnabled"
              :class="[
                'btn transition-all duration-200 border',
                isRestartButtonEnabled
                  ? 'bg-emerald-900 text-emerald-100 hover:bg-emerald-800 dark:bg-emerald-100 dark:text-emerald-800 dark:hover:bg-white border-emerald-900 shadow-sm hover:shadow-md'
                  : 'bg-gray-300 text-gray-500 cursor-not-allowed dark:bg-gray-700 dark:text-gray-500 border-gray-300 dark:border-gray-600 opacity-60',
              ]"
              @click.prevent="restart"
            >
              Restart
            </button>
          </div>
        </div>
      </footer>
    </div>
  
    <ModalBasic
      id="test-modal"
      :modalOpen="testModalOpen"
      @close-modal="testModalOpen = false"
      title="Test Result of Asset Diagnosis setting"
    >
      <div
        class="max-h-[75vh] overflow-y-auto px-4 py-2 space-y-4 text-xs text-gray-800 dark:text-gray-100"
      >
        <div v-if="isLoading" class="text-center py-10 text-sm text-gray-500">
          Loading test result...
        </div>
        <div v-else>
          <div class="text-sm font-semibold uppercase">Summary Report</div>
          <div class="flex flex-wrap justify-between gap-4">
            <div>Asset Type: {{ testData.AssetName.AssemblyType }}</div>
            <div>Serial Number: {{ testData.SerialNumber }}</div>
            <div>Channel: {{ testData.Channel }}</div>
          </div>
          <div class="font-semibold mb-4 mt-2">
            Result: {{ testResult.err }} Errors, {{ testResult.warn }} Warnings
          </div>
          <div
            class="grid grid-cols-[10%_1fr_10%] gap-2 font-semibold text-gray-500 border-b border-gray-200 pb-1"
          >
            <div class="text-left">AssemblyId</div>
            <div class="text-center">Detail</div>
            <div class="text-center">Status</div>
          </div>
          <div class="space-y-1">
            <div
              v-for="(item, index) in testData.Commissions"
              :key="index"
              class="grid grid-cols-[10%_1fr_10%] gap-2 px-1 py-1 rounded hover:bg-gray-50"
            >
              <div class="text-left">{{ item.AssemblyID }}</div>
              <div class="text-left text-gray-600">
                <span v-for="(msg, i) in item.Messages" :key="i">
                  • {{ msg }}
                </span>
                {{ msg }}
              </div>
              <div
                class="text-center font-bold"
                :class="{
                  'text-green-600': item.Status < 2,
                  'text-red-600': item.Status === 3,
                  'text-yellow-500': item.Status === 2,
                }"
              >
                {{ stList[item.Status] }}
              </div>
            </div>
          </div>
        </div>
        <div v-if="showDiagChart" class="pt-4 space-y-4">
          <button
            v-for="option in ChartTypes"
            :key="option"
            :value="option"
            @click.prevent="set_chart(option)"
            :class="[
              'btn border px-4 py-2 transition-colors duration-200 rounded-none first:rounded-l-lg last:rounded-r-lg w-1/2',
              selectedChart === option
                ? 'bg-violet-500 text-white border-violet-500'
                : 'bg-white text-violet-500 border-gray-200 hover:bg-gray-100 dark:bg-gray-800 dark:text-white dark:hover:bg-gray-900',
            ]"
          >
            {{ option }}
          </button>
          <div class="text-sm font-semibold text-gray-600 dark:text-gray-300">
            {{ selectedChart }} Chart
          </div>
          <LineChart
            v-if="selectedChart == 'Time Domain(Voltage)'"
            :label="waveformLabelT"
            :data="waveformDataT"
            :title="'3-Phase' + selectedChart"
            :legend="['V1', 'V2', 'V3']"
            :mode="'3phase'"
          />
          <LineChart
            v-if="selectedChart == 'Time Domain(Current)'"
            :label="waveformLabelT"
            :data="waveformDataT"
            :title="'3-Phase' + selectedChart"
            :legend="['I1', 'I2', 'I3']"
            :mode="'3phase'"
          />
          <LineChart
            v-if="selectedChart == 'Frequency Domain(Voltage)'"
            :label="waveformLabelT"
            :data="waveformDataT"
            :title="selectedChart"
            :legend="['VFFT']"
            :mode="'1phase'"
          />
          <LineChart
            v-if="selectedChart == 'Frequency Domain(Current)'"
            :label="waveformLabelT"
            :data="waveformDataT"
            :title="selectedChart"
            :legend="['IFFT']"
            :mode="'1phase'"
          />
        </div>
      </div>
      <div class="px-5 py-4 border-t border-gray-200 dark:border-gray-700/60">
        <div class="flex justify-end space-x-2">
          <button
            class="btn-sm border-gray-200 dark:border-gray-700/60 hover:border-gray-300 dark:hover:border-gray-600 text-gray-800 dark:text-white"
            @click.stop="testModalOpen = false"
          >
            Close
          </button>
        </div>
      </div>
    </ModalBasic>
  </template>
  
  <script>
  import {
    ref,
    watch,
    onMounted,
    inject,
    computed,
    provide,
    nextTick,
  } from "vue";
  import axios from "axios";
  import { useRoute } from "vue-router";
  import { useAuthStore } from "@/store/auth";
  import { useSetupStore } from "@/store/setup";
  import ModalBasic from "../../../pages/common/ModalBasic.vue";
  import TrendPanelCard from "./TrendPanelCard.vue";
  import AssetCard from "./AssetCard.vue";
  import ParamCard from "./ParamCard.vue";
  import NameplateCard from "./NameplateCard.vue";
  import EventCard2 from "./EventCard2.vue";
  import EventCard1 from "./EventCard1.vue";
  import AlarmCard from "./AlarmCard.vue";
  import LineChart from "../../../charts/connect/LineChart01_Echart2.vue";
  import { useInputDict } from "@/composables/useInputDict";
  import { settingValidator } from "@/utils/validation.js";
  import { useI18n } from "vue-i18n"; //
  
  export default {
    name: "ChannelPanel",
    props: ["channel"],
    components: {
      ModalBasic,
      TrendPanelCard,
      AssetCard,
      NameplateCard,
      EventCard2,
      EventCard1,
      ParamCard,
      AlarmCard,
      LineChart,
    },
    setup(props) {
      const { t } = useI18n();
      const route = useRoute();
      const authStore = useAuthStore();
      const setupStore = useSetupStore();
      const General_inputDict = inject("inputDict");
      const inputDict_main = inject("channel_main");
      const inputDict_sub = inject("channel_sub");
      const useDiagnosis = inject("useDiagnosis");
      const currentDiagnosis = inject("currentDiagnosis");
      const devMode = inject("devMode");
      const saveAllSettings = inject("saveAllSettings");
      const changeDiagnosis = inject("changeDiagnosis");
      const checkNameplateflag = inject("checkNameplateflag");
      const diagnosis_detail = inject("diagnosis_detail");
      // ✅ 강제 리렌더링을 위한 키
      const componentKey = ref(0);
      const testModalOpen = ref(false);
      const feedbackModalOpen = ref(false);
      const channel = ref(props.channel);
      const selectedNodeName = ref("");
      const selectedNodeType = ref("");
      const isAssetExist = ref(false);
      const tableData = ref([]);
      const modalData = ref([]);
      const paramData = ref([]);
      const isEditNameplates = ref(false);
      const isEditParameters = ref(false);
      const testData = ref({});
      const isLoading = ref(false);
      const waveformDataT = ref([]);
      const SR = ref(0);
      const waveformLabelT = ref([]);
      const Duration = ref(0);
      const testResult = ref({});
      const chartResult = ref({});
      const ChartTypes = ref([
        "Time Domain(Voltage)",
        "Time Domain(Current)",
        "Frequency Domain(Voltage)",
        "Frequency Domain(Current)",
      ]);
      const showDiagChart = ref(false);
      const selectedChart = ref("Time Domain(Voltage)");
      const stList = ref(["Info", "Pass", "Warning", "Error"]);
  
      const { parameterOptions, selectedTrendSetup } = useInputDict();
  
      provide("selectedTrendSetup", selectedTrendSetup);
  
      // ✅ 현재 채널의 inputDict를 반환하는 함수
      const getInputDict = () => {
        return props.channel === "Main"
          ? inputDict_main.value
          : inputDict_sub.value;
      };
  
      const isRestartButtonEnabled = computed(() => {
        const currentDict = getInputDict();
        const isChannelEnabled = currentDict.Enable === 1;
        const isDiagnosisActive = isCurrentChannelDiagnosisActive.value;
  
        if (isDiagnosisActive) {
          // Diagnosis가 활성화된 경우: 채널 사용 + Asset 존재
          const result = isChannelEnabled && isAssetExist.value;
  
          return result;
        } else {
          // Diagnosis가 비활성화된 경우: 채널 사용만
          const result = isChannelEnabled;
          return result;
        }
      });
      // ✅ 필드 업데이트 함수들 (연동 로직 제거)
      const updateField = (field, value) => {
        const currentDict = getInputDict();
        currentDict[field] = value;
        componentKey.value++;
      };
  
      const updateNestedField = (parent, field, value) => {
        const currentDict = getInputDict();
        if (!currentDict[parent]) {
          currentDict[parent] = {};
        }
        currentDict[parent][field] = value;
        componentKey.value++;
      };
  
      const updateArrayField = (parent, field, index, value) => {
        const currentDict = getInputDict();
        if (!currentDict[parent]) {
          currentDict[parent] = {};
        }
        if (!currentDict[parent][field]) {
          currentDict[parent][field] = [];
        }
        currentDict[parent][field][index] = value;
        componentKey.value++;
      };
  
      // ✅ device0 모드 기본 설정 적용 함수
      const applyDevice0Defaults = () => {
        const device0Params = [
          "Temperature",
          "Frequency",
          "Line Voltage",
          "Phase Voltage",
          "Current",
          "Unbalance",
          "PF",
          "THD",
          "TDD",
          "Power",
        ];
  
        selectedTrendSetup.value.params = [...device0Params];
  
        const currentDict = getInputDict();
        if (!currentDict.trendInfo) {
          currentDict.trendInfo = {};
        }
        currentDict.trendInfo.params = [...device0Params];
      };
      // ✅ 안전한 devMode watch - 초기화 후에만 실행
      const initializeWatchers = () => {
        if (devMode && typeof devMode.value !== "undefined") {
          watch(
            () => devMode.value,
            (newDevMode) => {
              if (newDevMode === "device0") {
                applyDevice0Defaults();
              } else {
                const defaultParams = Array(8).fill("None");
                const currentDict = getInputDict();
                if (!currentDict.trendInfo) {
                  currentDict.trendInfo = {};
                  selectedTrendSetup.value.params = [...defaultParams];
                }
              }
            },
            { immediate: true }
          );
  
          watch(
            () => getInputDict(),
            (newInputDict) => {
              if (newInputDict && newInputDict.trendInfo) {
                const trendInfo = newInputDict.trendInfo;
                const periodValue =
                  typeof trendInfo.period === "number" ? trendInfo.period : 5;
  
                if (devMode.value === "device0") {
                  const device0Params = [
                    "Temperature",
                    "Frequency",
                    "Line Voltage",
                    "Phase Voltage",
                    "Current",
                    "Unbalance",
                    "PF",
                    "THD",
                    "TDD",
                    "Power",
                  ];
                  selectedTrendSetup.value = {
                    period: periodValue,
                    params: [...device0Params],
                  };
                } else {
                  const savedParams = Array.isArray(trendInfo.params)
                    ? trendInfo.params
                    : [];
                  const fullParams = [...savedParams];
                  while (fullParams.length < 8) {
                    fullParams.push("None");
                  }
                  selectedTrendSetup.value = {
                    period: periodValue,
                    params: fullParams.slice(0, 8),
                  };
                }
              }
            },
            { deep: true, immediate: false }
          );
        }
      };
      const validTrendItems = computed(() =>
        selectedTrendSetup.value.params.filter((param) => param !== "None")
      );
  
      // ✅ 현재 채널의 Diagnosis 활성 상태 계산
      const isCurrentChannelDiagnosisActive = computed(() => {
        if (channel.value === "Main") {
          return (
            General_inputDict.value.useFuction?.diagnosis_main === true ||
            General_inputDict.value.useFuction?.diagnosis_main === 1
          );
        } else if (channel.value === "Sub") {
          return (
            General_inputDict.value.useFuction?.diagnosis_sub === true ||
            General_inputDict.value.useFuction?.diagnosis_sub === 1
          );
        }
        return false;
      });
  
      // ✅ 채널 변경 감지
      watch(
        () => props.channel,
        async (newChannel, oldChannel) => {
          if (newChannel !== oldChannel) {
            channel.value = newChannel;
  
            const targetDict = getInputDict();
            if (targetDict.assetInfo && targetDict.assetInfo.name !== "") {
              isAssetExist.value = true;
              await fetchAssetData(targetDict.assetInfo.name, channel.value);
              await fetchParamData(targetDict.assetInfo.name, channel.value);
              //console.log('Channel - DIAGNOSIS', diagnosis_detail.value);
            } else {
              isAssetExist.value = false;
              tableData.value = [];
              paramData.value = [];
            }
  
            componentKey.value++;
          }
        }
      );
  
      // ✅ Asset 이름 변경 감지
      watch(
        () => getInputDict().assetInfo?.name,
        async (newValue) => {
          if (newValue) {
            isAssetExist.value = true;
            await fetchAssetData(newValue);
            await fetchParamData(newValue);
            //console.log('Asset- DIAGNOSIS', diagnosis_detail.value);
          } else {
            isAssetExist.value = false;
            tableData.value = [];
            paramData.value = [];
          }
        }
      );
  
      onMounted(async () => {
        await nextTick();
  
        initializeWatchers();
  
        const currentDict = getInputDict();
        if (currentDict?.assetInfo && currentDict.assetInfo.name !== "") {
          isAssetExist.value = true;
          await fetchAssetData(currentDict.assetInfo.name);
          await fetchParamData(currentDict.assetInfo.name);
          //console.log('Mount - DIAGNOSIS', diagnosis_detail.value);
        }
  
        if (currentDict?.trendInfo) {
          const trendInfo = currentDict.trendInfo;
          const periodValue =
            typeof trendInfo.period === "number" ? trendInfo.period : 5;
  
          selectedTrendSetup.value = {
            period: periodValue,
            params: Array.isArray(trendInfo.params)
              ? [...trendInfo.params]
              : Array(8).fill("None"),
          };
  
          while (selectedTrendSetup.value.params.length < 8) {
            selectedTrendSetup.value.params.push("None");
          }
          selectedTrendSetup.value.params = selectedTrendSetup.value.params.slice(
            0,
            8
          );
        } else {
          selectedTrendSetup.value = {
            period: 5,
            params: Array(8).fill("None"),
          };
        }
      });
  
      const fetchAssetData = async (asset) => {
        try {
          if (asset == "") {
            asset = "GlobalTemplate";
          }
          const response = await axios.get(`/setting/getAssetConfig/${asset}`);
          
          if (response.data.success === true) {
            const allData = response.data.data;
            if (Array.isArray(allData)) {
              tableData.value = allData.filter((item) => item.Type === 0);
              if(channel.value == 'Main'){
                diagnosis_detail.value.main.use = true;
                diagnosis_detail.value.main.assetName = asset;
                diagnosis_detail.value.main.tableData = tableData.value;
              }else{
                diagnosis_detail.value.sub.use = true;
                diagnosis_detail.value.sub.assetName = asset;
                diagnosis_detail.value.sub.tableData = tableData.value;
              }
              
              if (authStore.getUserRole == "1") {
                modalData.value = allData.filter((item) => item.Type === 1);
              } else if (authStore.getUserRole == "0") {
                modalData.value = allData.filter((item) =>
                  [1, 2].includes(item.Type)
                );
              }
            } else {
              tableData.value = [];
              modalData.value = [];
            }
          } else {
            console.log("Data Load Fail!", response.data.error);
          }
        } catch (error) {
          console.error("Data import failed:", error);
        }
      };
  
      const fetchParamData = async (asset) => {
        if (asset == "") {
          asset = "GlobalTemplate";
        }
        try {
          const response = await axios.get(`/setting/getAssetParams/${asset}`);
          if (response.data.success === true) {
            const allData = response.data.data;
            if (Array.isArray(allData)) {
              paramData.value = allData;
              if(channel.value == 'Main'){
                diagnosis_detail.value.main.paramData = paramData.value;
              }else{
                diagnosis_detail.value.sub.paramData = paramData.value;
              }
            } else {
              paramData.value = [];
            }
          } else {
            console.log("Data Load Fail!", response.data.error);
          }
        } catch (error) {
          console.error("Data import failed:", error);
        }
      };
  
      // ✅ 수정된 savefile 함수 - Enable과 Diagnosis 연동 제거
      const savefile = async (source = "default") => {
        const currentDict = getInputDict();
  
        const isFTPEnabled = General_inputDict.value.useFuction.ftp === 1;
        const isMainDiagnosisToggled =
          General_inputDict.value.useFuction.diagnosis_main === true ||
          General_inputDict.value.useFuction.diagnosis_main === 1;
        const isSubDiagnosisToggled =
          General_inputDict.value.useFuction.diagnosis_sub === true ||
          General_inputDict.value.useFuction.diagnosis_sub === 1;
        const isMainChannelEnabled = inputDict_main.value.Enable === 1;
        const isSubChannelEnabled = inputDict_sub.value.Enable === 1;
  
        const isMainDiagnosisActive =
          isMainDiagnosisToggled && isMainChannelEnabled;
        const isSubDiagnosisActive = isSubDiagnosisToggled && isSubChannelEnabled;
        const isAnyDiagnosisEnabled =
          isMainDiagnosisActive || isSubDiagnosisActive;
  
        if (isFTPEnabled && isAnyDiagnosisEnabled) {
          alert(
            "Cannot use Waveform FTP and Diagnosis simultaneously. Please disable one of them."
          );
          return;
        }
  
        currentDict.channel = channel.value;
  
        if (devMode && devMode.value === "device0") {
          selectedTrendSetup.value.params = [
            "Temperature",
            "Frequency",
            "Line Voltage",
            "Phase Voltage",
            "Current",
            "Unbalance",
            "PF",
            "THD",
            "TDD",
            "Power",
          ];
        }
  
        currentDict.trendInfo = {
          period: selectedTrendSetup.value.period,
          params: [...selectedTrendSetup.value.params],
        };
        //console.log('Nameplate TableData:', tableData.value.length);
        if (currentDict.assetInfo.name != "" && tableData.value.length > 0) {
          const success = await checkTableData();
          if (!success) return;
          else if (isEditNameplates.value) changeDiagnosis.value.nameplate = true;
        }
        if (currentDict.assetInfo.name != "" && paramData.value.length > 0) {
          await setAssetParams();
        }
  
        const validationResult = settingValidator.validateAllSettings(
          General_inputDict.value,
          inputDict_main.value,
          inputDict_sub.value
        );
  
        if (!validationResult.isValid) {
          alert(settingValidator.formatErrorOnlyMessage());
          return;
        }
  
        if (validationResult.hasWarnings) {
          const proceed = confirm(settingValidator.formatWarningMessage());
          if (!proceed) return;
        }
  
        await saveAllSettings();
      };
  
      provide(
        "channel_inputDict",
        computed(() => getInputDict())
      );
      provide("savefile", savefile);
      provide("tableData", tableData);
      provide("modalData", modalData);
      provide("paramData", paramData);
      provide("isEditNameplates", isEditNameplates);
      provide("isEditParameters", isEditParameters);
      provide(
        "channel_inputDict",
        computed(() => getInputDict())
      );
      provide("savefile", savefile);
      provide("tableData", tableData);
      provide("modalData", modalData);
      provide("paramData", paramData);
      provide("isEditNameplates", isEditNameplates);
      provide("isEditParameters", isEditParameters);
  
      // ✅ 추가: 필수 함수들을 provide
      provide("updateNestedField", updateNestedField);
      provide("updateField", updateField);
      provide("updateArrayField", updateArrayField);
  
      // ✅ 추가: selectedTrendSetup 명시적으로 provide
      provide("selectedTrendSetup", selectedTrendSetup);
  
      // ✅ 추가: parameterOptions도 provide
      provide("parameterOptions", parameterOptions);
  
      // ✅ 추가: 디버깅을 위한 getInputDict 함수도 provide
      provide("getInputDict", getInputDict);
  
      const getMinValue = (row) => {
        return row.Min === "none" ? null : row.Min;
      };
  
      const getMaxValue = (row) => {
        return row.Max === "none" ? null : row.Max;
      };
  
      provide("getMinValue", getMinValue);
      provide("getMaxValue", getMaxValue);
  
      const options = [
        { value: 60, label: "60Hz" },
        { value: 50, label: "50Hz" },
      ];
  
      const setbtnOption = (value) => {
        updateNestedField("ptInfo", "linefrequency", value);
      };
  
      const ct_direction = [
        [
          { value: 0, label: "Positive" },
          { value: 1, label: "Negative" },
        ],
        [
          { value: 0, label: "Positive" },
          { value: 1, label: "Negative" },
        ],
        [
          { value: 0, label: "Positive" },
          { value: 1, label: "Negative" },
        ],
      ];
  
      const setct_direction = (value, index) => {
        updateArrayField("ctInfo", "direction", index, value);
      };
  
      const setAssetParams = async () => {
        try {
          const plainTableData = paramData.value.map((item) => ({ ...item }));
          const currentDict = getInputDict();
  
          const response = await axios.post(
            `/setting/setAssetParams/${currentDict.assetInfo.name}`,
            plainTableData,
            {
              headers: { "Content-Type": "application/json" },
            }
          );
  
          if (!response.data?.success) {
            alert(
              "❌ Failed to save asset settings: " +
                (response.data.error || "unknown error")
            );
          }
          if (isEditParameters.value) changeDiagnosis.value.parameter = true;
        } catch (error) {
          console.error("Error occurred while saving asset:", error);
          alert(
            "❌ Error occurred while saving asset: " +
              (error.response?.data?.error || error.message)
          );
        }
      };
  
      const checkTableData = async () => {
        try {
          checkNameplateflag.value = await checkNameplateConfig();
          if (checkNameplateflag.value) await setNameplateConfig();
          return true;
        } catch (error) {
          console.error("Error occurred while saving:", error);
          return false;
        }
      };
  
      const setNameplateConfig = async () => {
        try {
          const plainTableData = tableData.value.map((item) => ({ ...item }));
          const currentDict = getInputDict();
  
          const response = await axios.post(
            `/setting/setAssetConfig/${currentDict.assetInfo.name}`,
            plainTableData,
            {
              headers: { "Content-Type": "application/json" },
            }
          );
  
          if (!response.data?.success) {
            alert(
              "❌ Failed to save asset settings: " +
                (response.data.error || "unknown error")
            );
          }
        } catch (error) {
          console.error("Error occurred while saving asset:", error);
          alert(
            "❌ Error occurred while saving asset: " +
              (error.response?.data?.error || error.message)
          );
        }
      };
  
      const checkNameplateConfig = async () => {
        try {
          const plainTableData = tableData.value.map((item) => ({ ...item }));
          const currentDict = getInputDict();
  
          const response = await axios.post(
            `/setting/checkAssetConfig/${currentDict.assetInfo.name}`,
            plainTableData,
            {
              headers: { "Content-Type": "application/json" },
            }
          );
  
          if (response.data.success) {
            return response.data.result;
          } else {
            console.log("Not Respond API");
          }
        } catch (error) {
          console.log(error);
          return false;
        }
      };
  
      const set_chart = (data) => {
        selectedChart.value = data;
        const deltaT = 1 / SR.value;
        const deltaF = 1 / Duration.value;
        if (data.includes("Time")) {
          waveformLabelT.value = Array.from(
            { length: chartResult.value[data][0].length },
            (_, i) => i * deltaT
          );
          waveformDataT.value = chartResult.value[data];
        } else {
          waveformLabelT.value = Array.from(
            { length: chartResult.value[data].length },
            (_, i) => i * deltaF
          );
          waveformDataT.value = [chartResult.value[data]];
        }
      };
  
      const testDiagnosis = async () => {
        testModalOpen.value = true;
        isLoading.value = true;
        try {
          const currentDict = getInputDict();
          const response = await axios.get(
            `/setting/test/${currentDict.assetInfo.name}`
          );
          if (response.data.success === true) {
            testData.value = response.data.data;
            let errCount = 0,
              warnCount = 0;
            for (let i = 0; i < testData.value["Commissions"].length; i++) {
              if (testData.value["Commissions"][i]["Status"] == 3) errCount += 1;
              else if (testData.value["Commissions"][i]["Status"] == 3)
                warnCount += 1;
            }
            testResult.value = { err: errCount, warn: warnCount };
            if (response.data.waveT.length > 0) {
              showDiagChart.value = true;
            }
            if (showDiagChart.value) {
              SR.value = response.data.waveT["SR"];
              chartResult.value = {
                "Time Domain(Voltage)": response.data.waveT["Vwave"],
                "Time Domain(Current)": response.data.waveT["Iwave"],
                "Frequency Domain(Voltage)": response.data.waveT["Vfft"],
                "Frequency Domain(Current)": response.data.waveT["Ifft"],
              };
              Duration.value = response.data.waveT["Duration"];
              set_chart(selectedChart.value);
            }
          }
        } catch (e) {
          console.error("Error occurred:", e);
        } finally {
          isLoading.value = false;
        }
      };
      const restart = async () => {
        try {
          const response = await axios.get(`/setting/restartasset`);
          //console.log("Restart Response:", response);
  
          if (response.data.success === true) {
            alert("✅ Asset restarted successfully");
          } else {
            alert(
              "❌ Restart failed: " + (response.data.error || "Unknown error")
            );
          }
        } catch (e) {
          console.error("Error occurred:", e);
          alert("❌ Error occurred while restarting: " + e.message);
        }
      };
  
      return {
        channel,
        currentDiagnosis,
        savefile,
        authStore,
        selectedNodeName,
        selectedNodeType,
        feedbackModalOpen,
        options,
        setbtnOption,
        ct_direction,
        setct_direction,
        tableData,
        fetchAssetData,
        isAssetExist,
        modalData,
        selectedTrendSetup,
        validTrendItems,
        paramData,
        isEditNameplates,
        isEditParameters,
        parameterOptions,
        testModalOpen,
        testDiagnosis,
        restart,
        testData,
        isLoading,
        waveformDataT,
        waveformLabelT,
        SR,
        Duration,
        testResult,
        ChartTypes,
        chartResult,
        set_chart,
        selectedChart,
        stList,
        useDiagnosis,
        devMode,
        General_inputDict,
        getInputDict,
        updateField,
        updateNestedField,
        updateArrayField,
        componentKey,
        isCurrentChannelDiagnosisActive,
        isRestartButtonEnabled,
        t,
        showDiagChart,
        checkNameplateflag,
        diagnosis_detail,
      };
    },
  };
  </script>
  
  <style scoped>
  .btn-direction:first-child {
    @apply rounded-l-lg;
    border-top-right-radius: 0;
    border-bottom-right-radius: 0;
    @apply h-9 flex items-center justify-center;
  }
  
  .btn-direction:last-child {
    @apply rounded-r-lg;
    border-top-left-radius: 0;
    border-bottom-left-radius: 0;
    @apply h-9 flex items-center justify-center;
  }
  </style>
  