<template>
  <div class="grow dark:text-white">
    <!-- Panel body -->
    <div class="p-6 space-y-6">
      <!-- Plans -->
      <section>
        <div class="mb-4">
          <!-- 기존 mb-8 → mb-4 로 축소 -->
          <div class="flex items-center space-x-4 text-sm relative">
            <!-- FTP Checkbox -->
            <label v-if="false" class="flex items-center space-x-2">
              <input
                type="checkbox"
                class="form-checkbox text-violet-500 focus:ring-violet-500"
                :checked="
                  inputDict.useFuction.ftp === 1 ||
                  inputDict.useFuction.ftp === true
                "
                @change="handleFTPToggle($event)"
              />
              <span>{{ t("config.plansPanel.useWaveFormFTP") }}</span>
            </label>

            <!-- SNTP Checkbox -->
            <label v-if="false" class="flex items-center space-x-2">
              <input
                type="checkbox"
                class="form-checkbox text-violet-500 focus:ring-violet-500"
                :checked="
                  inputDict.useFuction.sntp === 1 ||
                  inputDict.useFuction.sntp === true
                "
                @change="
                  inputDict.useFuction.sntp = $event.target.checked ? 1 : 0
                "
              />
              <span>{{ t("config.plansPanel.useSNTP") }}</span>
            </label>

            <!-- FTP Toast Notification (간결하게 수정됨) -->
            <div
              v-if="showFTPToast"
              class="ml-4 flex items-center px-2 py-1 bg-blue-50 border border-blue-200 rounded-md shadow-sm transition-all duration-300 ease-in-out whitespace-nowrap"
            >
              <div class="flex items-center">
                <div
                  class="w-4 h-4 rounded-full bg-blue-500 flex items-center justify-center mr-1"
                >
                  <svg
                    class="w-2.5 h-2.5 text-white"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                </div>
                <span class="text-xs font-medium text-blue-800">
                  {{ t("config.plansPanel.info") }}
                </span>
              </div>
            </div>

            <!-- File Upload Modal -->
            <ModalBasic
              id="feedback-modal"
              :modalOpen="feedbackModalOpen"
              @close-modal="feedbackModalOpen = false"
              title="File Upload"
            >
              <!-- Modal content -->
              <div class="px-4 py-3">
                <!-- 패딩 소폭 축소 -->
                <div class="text-sm">
                  <div
                    class="font-medium text-gray-800 dark:text-gray-100 mb-2"
                  >
                    Setting file upload
                  </div>
                </div>
                <div class="space-y-2">
                  <div>
                    <label class="block text-sm font-medium mb-1" for="name">
                      file path <span class="text-red-500">*</span>
                    </label>
                    <input
                      id="filename"
                      class="form-input w-full px-2 py-1"
                      @change="handleFileUpload"
                      type="file"
                      required
                    />
                  </div>
                </div>
                <div class="text-sm mt-2">
                  <div class="font-medium text-gray-800 dark:text-gray-100">
                    {{ message }} 🙌
                  </div>
                </div>
              </div>

              <!-- Modal footer -->
              <div
                class="px-4 py-3 border-t border-gray-200 dark:border-gray-700/60"
              >
                <div class="flex flex-wrap justify-end space-x-2">
                  <button
                    class="btn-sm bg-gray-900 text-gray-100 hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white"
                    @click.prevent="upload"
                  >
                    Import
                  </button>
                  <button
                    class="btn-sm border-gray-200 dark:border-gray-700/60 hover:border-gray-300 dark:hover:border-gray-600 text-gray-800 dark:text-white"
                    @click.stop="feedbackModalOpen = false"
                  >
                    Cancel
                  </button>
                </div>
              </div>
            </ModalBasic>
          </div>
        </div>

        <div class="grid grid-cols-12 gap-6">
          <!--Device Information,TCP/IP,Serial,Demand -->

          <!-- Device Information-->
          <div
            class="relative col-span-full xl:col-span-3 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg"
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
                  >
                    <path
                      d="M12 17a.833.833 0 01-.833-.833 3.333 3.333 0 00-3.334-3.334.833.833 0 110-1.666 3.333 3.333 0 003.334-3.334.833.833 0 111.666 0 3.333 3.333 0 003.334 3.334.833.833 0 110 1.666 3.333 3.333 0 00-3.334 3.334c0 .46-.373.833-.833.833z"
                    />
                  </svg>
                </div>
                <h3
                  class="text-lg text-gray-800 dark:text-gray-100 font-semibold"
                >
                  {{ t("config.plansPanel.devInformation.title")
                  }}<!--Device Information-->
                </h3>
              </header>
            </div>
            <div class="px-4 py-4 space-y-3">
              <div>
                <label
                  class="block text-sm font-medium mb-2"
                  for="device-model"
                >
                  {{ t("config.plansPanel.devInformation.name")
                  }}<!--Name--></label
                >
                <input
                  v-model="inputDict.deviceInfo.name"
                  class="form-input w-full"
                  type="text"
                  placeholder="Enter name"
                />
              </div>
              <div>
                <label class="block text-sm font-medium mb-2" for="location">
                  {{ t("config.plansPanel.devInformation.loc")
                  }}<!--Location--></label
                >
                <input
                  v-model="inputDict.deviceInfo.location"
                  class="form-input w-full"
                  type="text"
                  placeholder="Enter location"
                />
              </div>

              <div>
                <label
                  class="block text-sm font-medium mb-2"
                  for="serial-number"
                >
                  {{ t("config.plansPanel.devInformation.sn")
                  }}<!--Serial Number-->
                </label>
                <input
                  readonly
                  v-model="inputDict.deviceInfo.serial_number"
                  class="form-input w-full"
                  type="text"
                />
              </div>

              <div>
                <label class="block text-sm font-medium mb-2" for="mac-address">
                  {{ t("config.plansPanel.devInformation.mac")
                  }}<!--Mac Address-->
                </label>
                <input
                  readonly
                  v-model="inputDict.deviceInfo.mac_address"
                  class="form-input w-full"
                  type="text"
                />
              </div>
            </div>
          </div>

          <!--  TCP/IP -->
          <div
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
                    class="w-6 h-6 shrink-0 fill-current text-white"
                    viewBox="0 0 24 24"
                  >
                    <path
                      d="M12 20c.83 0 1.5.67 1.5 1.5S12.83 23 12 23s-1.5-.67-1.5-1.5S11.17 20 12 20Zm-5.5-4.5a1 1 0 0 0 1.41 0 5.99 5.99 0 0 1 8.18 0 1 1 0 1 0 1.41-1.41 7.99 7.99 0 0 0-11 0 1 1 0 0 0 0 1.41Zm-3.54-3.54a1 1 0 0 0 1.41 0 10.93 10.93 0 0 1 15.44 0 1 1 0 0 0 1.42-1.42 12.93 12.93 0 0 0-18.28 0 1 1 0 0 0 0 1.42Z"
                    />
                  </svg>
                </div>
                <h3
                  class="text-lg text-gray-800 dark:text-gray-100 font-semibold"
                >
                  {{ t("config.plansPanel.communication.title")
                  }}<!--Communication-->
                </h3>
              </header>
            </div>
            <div class="px-4 py-4 space-y-3">
              <div>
                <label class="block text-sm font-medium mb-2" for="ip-address">
                  {{ t("config.plansPanel.communication.ip")
                  }}<!--IP Address-->
                </label>
                <input
                  v-model="inputDict.tcpip.ip_address"
                  class="form-input w-full"
                  type="text"
                />
              </div>

              <div>
                <label class="block text-sm font-medium mb-2" for="subnet-mask">
                  {{ t("config.plansPanel.communication.sm")
                  }}<!--Subnet Mask-->
                </label>
                <input
                  v-model="inputDict.tcpip.subnet_mask"
                  class="form-input w-full"
                  type="text"
                />
              </div>

              <div>
                <label class="block text-sm font-medium mb-2" for="gateway">
                  {{ t("config.plansPanel.communication.gw")
                  }}<!--Gateway--></label
                >
                <input
                  v-model="inputDict.tcpip.gateway"
                  class="form-input w-full"
                  type="text"
                />
              </div>

              <div>
                <label class="block text-sm font-medium mb-2" for="dns-server">
                  {{ t("config.plansPanel.communication.dns")
                  }}<!--DNS Server--></label
                >
                <input
                  v-model="inputDict.tcpip.dnsserver"
                  class="form-input w-full"
                  type="text"
                />
              </div>
            </div>
          </div>

          <!-- Modbus -->
          <div
            class="relative col-span-full xl:col-span-3 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg"
          >
            <div
              class="absolute top-0 left-0 right-0 h-0.5 bg-violet-500"
              aria-hidden="true"
            ></div>
            <div
              class="px-5 pt-5 pb-6 border-b border-gray-200 dark:border-gray-700/60"
            >
              <header class="flex items-center mb-2">
                <div class="w-6 h-6 rounded-full shrink-0 bg-violet-500 mr-3">
                  <svg
                    class="w-6 h-6 fill-current text-white"
                    viewBox="0 0 24 24"
                  >
                    <path
                      d="M4 12h16m-4-4 4 4-4 4M8 8l-4 4 4 4"
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
                  {{ t("config.plansPanel.modbus.title") }}
                  <!--Modbus-->
                </h3>
              </header>
            </div>
            <div class="px-4 py-3 space-y-4">
              <!-- TCP Port Section -->
              <div
                class="text-xs text-gray-800 dark:text-gray-100 font-semibold uppercase mb-1 mt-2"
              >
                {{ t("config.plansPanel.modbus.tcp")
                }}<!--TCP Port-->
              </div>

              <div>
                <label class="block text-sm font-medium mb-2" for="tcp-port"
                  >{{ t("config.plansPanel.modbus.tcp")
                  }}<!--TCP Port--></label
                >
                <input
                  v-model.number="inputDict.modbus.tcp_port"
                  class="form-input w-full"
                  type="number"
                />
              </div>

              <div v-if="false">
                <label class="block text-sm font-medium mb-2" for="tcp-timeout">
                  {{ t("config.plansPanel.modbus.tcpTO")
                  }}<!--Timeout (ms)-->
                </label>
                <input
                  v-model.number="inputDict.modbus.tcp_timeout"
                  class="form-input w-full"
                  type="text"
                  placeholder="5000"
                />
              </div>

              <!-- Serial Information Section -->
              <div
                class="text-xs text-gray-800 dark:text-gray-100 font-semibold uppercase mb-1 mt-2"
              >
                {{ t("config.plansPanel.modbus.snInfor")
                }}<!--Serial Information-->
              </div>
              <div>
                <label class="block text-sm font-medium mb-2" for="rtu-id">
                  {{ t("config.plansPanel.modbus.rtuid")
                  }}<!--RTU/ID--></label
                >
                <input
                  v-model.number="inputDict.modbus.modbus_id"
                  class="form-input w-full"
                  type="number"
                />
              </div>

              <!-- Baud Rate and Parity in one row -->
              <div class="flex space-x-3">
                <div class="flex-1">
                  <label class="block text-sm font-medium mb-2" for="baud-rate">
                    {{ t("config.plansPanel.modbus.baudrate")
                    }}<!--Baud Rate--></label
                  >
                  <select
                    v-model.number="inputDict.modbus.baud_rate"
                    class="form-select w-full bg-white border-gray-300 rounded-md shadow-sm focus:ring-violet-500 focus:border-violet-500"
                  >
                    <option :value="0">9600</option>
                    <option :value="1">19200</option>
                    <option :value="2">38400</option>
                    <option :value="3">57600</option>
                    <option :value="4">115200</option>
                  </select>
                </div>

                <div class="flex-1">
                  <label class="block text-sm font-medium mb-2" for="parity"
                    >{{ t("config.plansPanel.modbus.parity")
                    }}<!--Parity--></label
                  >
                  <select
                    v-model.number="inputDict.modbus.parity"
                    class="form-select w-full bg-white border-gray-300 rounded-md shadow-sm focus:ring-violet-500 focus:border-violet-500"
                  >
                    <option :value="0">None</option>
                    <option :value="1">Odd</option>
                    <option :value="2">Even</option>
                  </select>
                </div>
              </div>

              <!-- Data Bits and Stop Bits -->
              <div class="flex space-x-3">
                <div class="flex-1">
                  <label class="block text-sm font-medium mb-2" for="data-bits">
                    {{ t("config.plansPanel.modbus.dbits")
                    }}<!--Data Bits-->
                  </label>
                  <select
                    v-model.number="inputDict.modbus.data_bits"
                    class="form-select w-full bg-white border-gray-300 rounded-md shadow-sm focus:ring-violet-500 focus:border-violet-500"
                  >
                    <option :value="7">7</option>
                    <option :value="8">8</option>
                  </select>
                </div>

                <div class="flex-1">
                  <label class="block text-sm font-medium mb-2" for="stop-bits">
                    {{ t("config.plansPanel.modbus.sbits")
                    }}<!--Stop Bits-->
                  </label>
                  <select
                    v-model.number="inputDict.modbus.stop_bits"
                    class="form-select w-full bg-white border-gray-300 rounded-md shadow-sm focus:ring-violet-500 focus:border-violet-500"
                  >
                    <option :value="1">1</option>
                    <option :value="2">2</option>
                  </select>
                </div>
              </div>
            </div>
          </div>
          <div class="relative col-span-full xl:col-span-3 flex flex-col gap-6">
            <!-- Demand-->
            <!-- <div
              class="relative xl:col-span-3 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg flex flex-col flex-1"
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
                    {{ t("config.plansPanel.demand.title") }}
                    
                  </h3>
                </header>
              </div>
              <div class="px-4 py-3 space-y-4">
               
                <div>
                  <label class="block text-sm font-medium mb-1.5" for="model"
                    >{{ t("config.plansPanel.demand.interval")
                    }}</label
                  >
                  <select
                    v-model="inputDict.demand_interval"
                    class="form-select w-full bg-white border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option value="1">
                      1 {{ t("config.plansPanel.demand.minutes")
                      }}
                    </option>
                    <option value="5">
                      5 {{ t("config.plansPanel.demand.minutes")
                      }}
                    </option>
                    <option value="15">
                      15 {{ t("config.plansPanel.demand.minutes")
                      }}
                    </option>
                    <option value="30">
                      30 {{ t("config.plansPanel.demand.minutes")
                      }}
                    </option>
                    <option value="60">
                      1 {{ t("config.plansPanel.demand.hours")
                      }}
                    </option>
                  </select>
                </div>
              </div>
            </div> -->
            <div
              class="relative xl:col-span-3 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg flex flex-col flex-1"
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
                    ETC
                    <!--Demand-->
                  </h3>
                </header>
              </div>
              <div class="px-4 py-3 space-y-4">
                <!-- ✅ 전체 여백 조정 -->

                <div>
                  <label class="block text-sm font-medium mb-1.5" for="model"
                    >PF Sign</label
                  >
                  <select
                    v-model.number="inputDict.pf_sign"
                    class="form-select w-full bg-white border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option :value="0">IEC</option>
                    <option :value="1">IEEE</option>
                  </select>
                </div>

                <div>
                  <label class="block text-sm font-medium mb-1.5" for="model"
                    >VA type</label
                  >
                  <select
                    v-model.number="inputDict.va_type"
                    class="form-select w-full bg-white border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option :value="0">RMS</option>
                    <option :value="1">vector</option>
                  </select>
                </div>
                <div>
                  <label class="block text-sm font-medium mb-1.5" for="model"
                    >Unbalance</label
                  >
                  <select
                    v-model.number="inputDict.unbalance"
                    class="form-select w-full bg-white border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option :value="0">Nema</option>
                    <option :value="1">Sequence Components</option>
                  </select>
                </div>
              </div>
            </div>
          </div>

          <WaveformFTPCard
            :showFTP="
              inputDict.useFuction.ftp === 1 ||
              inputDict.useFuction.ftp === true
            "
          />
          <SNTPCard
            :showSNTP="
              inputDict.useFuction.sntp === 1 ||
              inputDict.useFuction.sntp === true
            "
          />
          <div
            class="relative xl:col-span-6 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg flex flex-col flex-1"
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
                  Sub Channel Setting
                 
                </h3>
              </header>
            </div>
            <div class="px-4 py-3 space-y-4">
              <!-- ✅ 전체 여백 조정 -->

              <div>
                <label class="block text-sm font-medium mb-1.5" for="model"
                  >Channel Count</label
                >
                <select
                  
                  class="form-select w-full bg-white border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                >
                  <option :value="0">12</option>
                  <option :value="1">24</option>
                  <option :value="2">36</option>
                </select>
              </div>
            </div>
          </div>
          <DiagnosisSetupCard1 :showDiagnosis="useDiagnosis || showDiagnosis" />
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
        </div>
      </div>
    </footer>
  </div>
</template>

<script>
import { ref, onMounted, watch, inject, computed, provide } from "vue";
import axios from "axios";
import ModalBasic from "../../../components/ModalBasic.vue";
import flatPickr from "vue-flatpickr-component";
import "flatpickr/dist/flatpickr.css";
import { useSetupStore } from "@/store/setup";
import { useAuthStore } from "@/store/auth";
import { settingValidator } from "@/utils/validation.js";
import WaveformFTPCard from "./WaveformFTPCard.vue";
import SNTPCard from "./SNTPCard.vue";
import DiagnosisSetupCard1 from "./General_diagnosisCard.vue";
import { useI18n } from "vue-i18n"; //
export default {
  name: "PlansPanel",
  props: ["channel"],
  components: {
    ModalBasic,
    flatPickr,
    WaveformFTPCard,
    SNTPCard,
    DiagnosisSetupCard1,
  },
  setup(props) {
    const { t } = useI18n();
    const setupStore = useSetupStore();
    const channel = ref(props.channel);
    const errorMessage = ref("");
    const feedbackModalOpen = ref(false);
    const selectedFile = ref(null);
    const toggle1 = ref("YES");
    const message = ref("Select upload setting file");
    const inputDict = inject("inputDict");
    const diagnosisData = inject("diagnosisData");
    const advancedData = inject("advancedData");
    const useDiagnosis = inject("useDiagnosis");
    const devMode = inject("devMode");
    const saveAllSettings = inject("saveAllSettings");

    // FTP Toast 상태 추가
    const showFTPToast = ref(false);

    // 채널 데이터를 setup 시점에 inject
    const channel_main = inject("channel_main");
    const channel_sub = inject("channel_sub");

    // FTP 토글 핸들러 추가
    const handleFTPToggle = (event) => {
      const isChecked = event.target.checked;
      inputDict.value.useFuction.ftp = isChecked ? 1 : 0;

      // 토스트 표시/숨김 (반대로 변경)
      if (isChecked) {
        showFTPToast.value = false; // 체크되면 토스트 숨김
      } else {
        showFTPToast.value = true; // 체크 해제되면 토스트 표시
      }
    };

    // 초기 토스트 상태 설정을 위한 watch 추가
    watch(
      () => inputDict.value.useFuction?.ftp,
      (newVal) => {
        // FTP가 비활성화(0 또는 false)되어 있으면 토스트 표시
        showFTPToast.value = newVal !== 1;
      },
      { immediate: true }
    );
    // inputDict 초기화 시 기본값 설정 (필요한 경우)
    if (!inputDict.value.modbus) {
      inputDict.value.modbus = {
        mode: "serial",
        tcp_port: 502,
        tcp_timeout: 5000,
        modbus_id: 1,
        baud_rate: 0,
        parity: 0,
        data_bits: 8,
        stop_bits: 1,
      };
    }

    const setDiagnosisSetting = async () => {
      try {
        //console.log("diagnosisData", diagnosisData.value);
        const response = await axios.post(
          "/setting/setDiagnosisSetting",
          diagnosisData.value,
          {
            headers: { "Content-Type": "application/json" },
            withCredentials: true,
          }
        );
        console.log("🔍 Diagnosis 저장 응답:", response.data);
        if (response.data && response.data.success) {
          console.log("✅ Diagnosis 설정 저장 성공!");
        } else {
          console.log("❌ Diagnosis 설정 저장 실패!");
        }
      } catch (err) {
        console.log("❌ 저장 중 오류 발생: " + err.message);
      }
    };

    const setDiagnosisProfile = async () => {
      try {
        const response = await axios.post(
          "/setting/setDiagnosisProfile",
          advancedData.value,
          {
            headers: { "Content-Type": "application/json" },
            withCredentials: true,
          }
        );
        console.log("🔍 Diagnosis 저장 응답:", response.data);
        if (response.data && response.data.success) {
          console.log("✅ Diagnosis 설정 저장 성공!");
        } else {
          console.log("❌ Diagnosis 설정 저장 실패!");
        }
      } catch (err) {
        console.log("❌ 저장 중 오류 발생: " + err.message);
      }
    };

    const handleFileUpload = (event) => {
      selectedFile.value = event.target.files[0];
    };

    const upload = async () => {
      if (!selectedFile.value) {
        message.value = "파일을 선택하세요!";
        return;
      }

      const formData = new FormData();
      formData.append("file", selectedFile.value);
      //console.log(formData);
      try {
        const response = await axios.post("/setting/upload", formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });
        if (response.data.passOK == "1") {
          message.value = "Upload Success : " + response.data.file_path;
          feedbackModalOpen.value = false;
          await GetSetting();
        } else {
          message.value = response.data.error;
        }
      } catch (error) {
        message.value = "업로드 실패: " + error.response.data.error;
      }
    };

    const savefile = async () => {
      // 현재 inputDict 값을 직접 참조하여 최신 상태 확인
      const isFTPEnabled = inputDict.value.useFuction.ftp === 1;

      // inputDict의 diagnosis 설정값을 사용 (Setting.vue에서 watch로 연동됨)
      const isMainDiagnosisEnabled =
        inputDict.value.useFuction.diagnosis_main === true ||
        inputDict.value.useFuction.diagnosis_main === 1;
      const isSubDiagnosisEnabled =
        inputDict.value.useFuction.diagnosis_sub === true ||
        inputDict.value.useFuction.diagnosis_sub === 1;

      // channel Enable 값도 함께 확인 (실제 Enable 상태)
      const isMainChannelEnabled =
        channel_main.value && channel_main.value.Enable === 1;
      const isSubChannelEnabled =
        channel_sub.value && channel_sub.value.Enable === 1;

      // 둘 중 하나라도 활성화되어 있으면 true
      const isAnyDiagnosisEnabled =
        (isMainDiagnosisEnabled && isMainChannelEnabled) ||
        (isSubDiagnosisEnabled && isSubChannelEnabled);

      // FTP가 활성화되어 있고, 어느 채널이든 Diagnosis가 활성화되어 있으면 경고
      if (isFTPEnabled && isAnyDiagnosisEnabled) {
        alert(
          "Cannot use Waveform FTP and Diagnosis simultaneously. Please disable one of them."
        );
        return;
      }
      const validationResult = settingValidator.validateAllSettings(
        inputDict.value, // General 설정
        channel_main.value, // Main Channel 설정
        channel_sub.value // Sub Channel 설정
      );

      // 에러가 있으면 저장 차단
      if (!validationResult.isValid) {
        alert(settingValidator.formatErrorOnlyMessage());
        return;
      }

      // 경고만 있으면 사용자 확인 후 진행
      if (validationResult.hasWarnings) {
        const proceed = confirm(settingValidator.formatWarningMessage());
        if (!proceed) return;
      }

      // 저장 진행
      await saveAllSettings();
    };

    const GetSetting = async () => {
      try {
        const response = await axios.get("/setting/getSettingData/General");

        if (response.data.passOK == 1) {
          Object.assign(inputDict.value, response.data.data);
          setupStore.setsetupFromFile(true);
        }
      } catch (error) {
        console.error("데이터 가져오기 실패:", error);
        errorMessage.value = "데이터를 불러오는 중 오류가 발생했습니다.";
      }
    };

    const dateConfig = {
      dateFormat: "Y-m-d",
      defaultDate: new Date(),
      onChange: (selectedDates, dateStr) => {
        console.log("Selected Date:", dateStr);
      },
    };

    const triggerTimePicker = (e) => {
      e.target.showPicker();
    };

    return {
      inputDict,
      channel,
      errorMessage,
      feedbackModalOpen,
      message,
      selectedFile,
      upload,
      handleFileUpload,
      savefile,
      toggle1,
      GetSetting,
      useDiagnosis,
      setDiagnosisSetting,

      dateConfig,
      triggerTimePicker,
      devMode,
      diagnosisData,
      advancedData,
      t,
      showFTPToast,
      handleFTPToggle,
    };
  },
};
</script>
