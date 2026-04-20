<template>
  <div class="grow dark:text-white">
    <div class="p-6 space-y-6">

      <!-- Row 1: Device Info / Device Function / Communication / Measurement Type -->
      <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-5">

        <!-- Device Information -->
        <div class="relative bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg">
          <div class="absolute top-0 left-0 right-0 h-0.5 bg-green-500" aria-hidden="true"></div>
          <div class="px-5 pt-5 pb-4 border-b border-gray-200 dark:border-gray-700/60">
            <header class="flex items-center">
              <div class="w-6 h-6 rounded-full shrink-0 bg-green-500 mr-3">
                <svg class="w-6 h-6 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M6 2h12a2 2 0 0 1 2 2v16a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2z"/><path d="M8 6h8"/><path d="M8 10h8"/><path d="M8 14h4"/>
                </svg>
              </div>
              <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold">Device Information</h3>
            </header>
          </div>
          <div class="px-4 py-4 space-y-3">
            <div>
              <label class="block text-sm font-medium mb-2">Name</label>
              <input v-model="inputDict.deviceInfo.name" class="form-input w-full" type="text" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-2">Location</label>
              <input v-model="inputDict.deviceInfo.location" class="form-input w-full" type="text" />
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium mb-2">Serial Number</label>
                <input v-model="inputDict.deviceInfo.serial_number" class="form-input w-full" type="text" />
              </div>
              <div>
                <label class="block text-sm font-medium mb-2">MAC Address</label>
                <input v-model="inputDict.deviceInfo.mac_address" class="form-input w-full font-mono" type="text" readonly />
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium mb-2">Timezone</label>
              <div class="flex items-center gap-3">
                <select v-model="inputDict.deviceInfo.timezone" class="form-select flex-1">
                  <option v-for="tz in timezones" :key="tz" :value="tz">{{ tz }}</option>
                </select>
                <span class="text-sm font-mono font-medium text-green-600 dark:text-green-400 whitespace-nowrap">{{ utcOffsetLabel }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Device Function + SNTP (수직 배치) -->
        <div class="flex flex-col gap-5">
          <!-- Device Function -->
          <div class="relative flex-1 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg">
            <div class="absolute top-0 left-0 right-0 h-0.5 bg-violet-500" aria-hidden="true"></div>
            <div class="px-5 pt-5 pb-4 border-b border-gray-200 dark:border-gray-700/60">
              <header class="flex items-center">
                <div class="w-6 h-6 rounded-full shrink-0 bg-violet-500 mr-3">
                  <svg class="w-6 h-6 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                  </svg>
                </div>
                <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold">Device Function</h3>
              </header>
            </div>
            <div class="px-4 py-4 space-y-4">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium">SNTP</span>
                <label class="relative inline-flex items-center cursor-pointer">
                  <input type="checkbox" class="sr-only peer" :checked="inputDict.useFunction.sntp === 1" @change="inputDict.useFunction.sntp = $event.target.checked ? 1 : 0" />
                  <div class="w-9 h-5 bg-gray-200 peer-focus:outline-none rounded-full peer dark:bg-gray-600 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-violet-500"></div>
                </label>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium">MQTT</span>
                <label class="relative inline-flex items-center cursor-pointer">
                  <input type="checkbox" class="sr-only peer" :checked="inputDict.useFunction.mqtt === 1" @change="inputDict.useFunction.mqtt = $event.target.checked ? 1 : 0" />
                  <div class="w-9 h-5 bg-gray-200 peer-focus:outline-none rounded-full peer dark:bg-gray-600 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-violet-500"></div>
                </label>
              </div>
              <div class="flex items-center justify-between gap-3">
                <span class="text-sm font-medium">Modbus Serial</span>
                <div class="inline-flex rounded-md border border-gray-200 dark:border-gray-700 overflow-hidden text-xs">
                  <button type="button"
                    class="px-3 py-1 transition-colors"
                    :class="inputDict.useFunction.modbus_serial === 0 ? 'bg-violet-500 text-white' : 'bg-white dark:bg-gray-800 text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700'"
                    @click="inputDict.useFunction.modbus_serial = 0">Disabled</button>
                  <button type="button"
                    class="px-3 py-1 transition-colors border-l border-gray-200 dark:border-gray-700"
                    :class="inputDict.useFunction.modbus_serial === 1 ? 'bg-violet-500 text-white' : 'bg-white dark:bg-gray-800 text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700'"
                    @click="inputDict.useFunction.modbus_serial = 1">Master</button>
                  <button type="button"
                    class="px-3 py-1 transition-colors border-l border-gray-200 dark:border-gray-700"
                    :class="inputDict.useFunction.modbus_serial === 2 ? 'bg-violet-500 text-white' : 'bg-white dark:bg-gray-800 text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700'"
                    @click="inputDict.useFunction.modbus_serial = 2">Slave</button>
                </div>
              </div>
            </div>
          </div>

          <!-- SNTP (SNTP 활성화 시) -->
          <div v-if="inputDict.useFunction.sntp === 1" class="relative flex-1 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg">
            <div class="absolute top-0 left-0 right-0 h-0.5 bg-teal-500" aria-hidden="true"></div>
            <div class="px-5 pt-5 pb-4 border-b border-gray-200 dark:border-gray-700/60">
              <header class="flex items-center">
                <div class="w-6 h-6 rounded-full shrink-0 bg-teal-500 mr-3">
                  <svg class="w-6 h-6 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>
                  </svg>
                </div>
                <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold">SNTP</h3>
              </header>
            </div>
            <div class="px-4 py-4 space-y-3">
              <div>
                <label class="block text-sm font-medium mb-2">Host</label>
                <input v-model="inputDict.sntpInfo.host" class="form-input w-full font-mono" type="text" />
              </div>
            </div>
          </div>
        </div>

        <!-- Communication (TCP/IP + Modbus TCP) -->
        <div class="relative bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg">
          <div class="absolute top-0 left-0 right-0 h-0.5 bg-sky-500" aria-hidden="true"></div>
          <div class="px-5 pt-5 pb-4 border-b border-gray-200 dark:border-gray-700/60">
            <header class="flex items-center">
              <div class="w-6 h-6 rounded-full shrink-0 bg-sky-500 mr-3">
                <svg class="w-6 h-6 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M5 12.55a11 11 0 0 1 14.08 0"/><path d="M1.42 9a16 16 0 0 1 21.16 0"/><path d="M8.53 16.11a6 6 0 0 1 6.95 0"/><line x1="12" y1="20" x2="12.01" y2="20"/>
                </svg>
              </div>
              <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold">Communication</h3>
            </header>
          </div>
          <div class="px-4 py-4 space-y-3">
            <div class="flex items-center gap-2 mb-2">
              <div class="w-1 h-4 bg-sky-500 rounded-full"></div>
              <span class="text-xs text-gray-800 dark:text-gray-100 font-semibold uppercase tracking-wider">TCP/IP</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-sm font-medium">DHCP</span>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" class="sr-only peer" :checked="inputDict.tcpip.dhcp === 1" @change="inputDict.tcpip.dhcp = $event.target.checked ? 1 : 0" />
                <div class="w-9 h-5 bg-gray-200 peer-focus:outline-none rounded-full peer dark:bg-gray-600 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-violet-500"></div>
              </label>
            </div>
            <template v-if="inputDict.tcpip.dhcp !== 1">
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium mb-2">IP Address</label>
                  <input v-model="inputDict.tcpip.ip_address" class="form-input w-full font-mono" type="text" />
                </div>
                <div>
                  <label class="block text-sm font-medium mb-2">Subnet Mask</label>
                  <input v-model="inputDict.tcpip.subnet_mask" class="form-input w-full font-mono" type="text" />
                </div>
              </div>
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium mb-2">Gateway</label>
                  <input v-model="inputDict.tcpip.gateway" class="form-input w-full font-mono" type="text" />
                </div>
                <div>
                  <label class="block text-sm font-medium mb-2">DNS Server</label>
                  <input v-model="inputDict.tcpip.dnsserver" class="form-input w-full font-mono" type="text" />
                </div>
              </div>
            </template>

            <div class="flex items-center gap-2 mt-4 mb-2">
              <div class="w-1 h-4 bg-sky-500 rounded-full"></div>
              <span class="text-xs text-gray-800 dark:text-gray-100 font-semibold uppercase tracking-wider">Modbus TCP</span>
            </div>
            <div>
              <label class="block text-sm font-medium mb-2">TCP Port</label>
              <input v-model="inputDict.modbus.tcp_port" class="form-input w-full font-mono" type="text" />
            </div>
          </div>
        </div>

        <!-- ETC (Display / System) -->
        <div class="relative bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg">
          <div class="absolute top-0 left-0 right-0 h-0.5 bg-amber-500" aria-hidden="true"></div>
          <div class="px-5 pt-5 pb-4 border-b border-gray-200 dark:border-gray-700/60">
            <header class="flex items-center">
              <div class="w-6 h-6 rounded-full shrink-0 bg-amber-500 mr-3">
                <svg class="w-6 h-6 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <rect x="2" y="3" width="20" height="14" rx="2" ry="2" /><line x1="8" y1="21" x2="16" y2="21" /><line x1="12" y1="17" x2="12" y2="21" />
                </svg>
              </div>
              <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold">ETC</h3>
            </header>
          </div>
          <div class="px-4 py-4 space-y-3">
            <div>
              <label class="block text-sm font-medium mb-2">Backlight Off Time (sec)</label>
              <input v-model.number="inputDict.etc.backlight_off_time" class="form-input w-full" type="number" min="0" />
            </div>
            <div class="flex items-center justify-between">
              <span class="text-sm font-medium">Auto Rotation</span>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" class="sr-only peer" :checked="inputDict.etc.autorotation === 1" @change="inputDict.etc.autorotation = $event.target.checked ? 1 : 0" />
                <div class="w-9 h-5 bg-gray-200 peer-focus:outline-none rounded-full peer dark:bg-gray-600 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-amber-500"></div>
              </label>
            </div>
            <div>
              <label class="block text-sm font-medium mb-2">Brightness</label>
              <input v-model.number="inputDict.etc.brightness" class="form-input w-full" type="number" min="0" max="100" />
            </div>
            <div class="flex items-center justify-between">
              <span class="text-sm font-medium">Test Mode</span>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" class="sr-only peer" :checked="inputDict.etc.testmode === 1" @change="inputDict.etc.testmode = $event.target.checked ? 1 : 0" />
                <div class="w-9 h-5 bg-gray-200 peer-focus:outline-none rounded-full peer dark:bg-gray-600 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-amber-500"></div>
              </label>
            </div>
            <div>
              <label class="block text-sm font-medium mb-2">Update Interval (sec)</label>
              <input v-model.number="inputDict.etc.update_interval" class="form-input w-full" type="number" min="1" />
            </div>
          </div>
        </div>

      </div>

      <!-- Row 2: MQTT -->
      <div v-if="inputDict.useFunction.mqtt === 1"
        class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-5">

        <!-- MQTT -->
        <div v-if="inputDict.useFunction.mqtt === 1" class="relative bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg">
          <div class="absolute top-0 left-0 right-0 h-0.5 bg-indigo-500" aria-hidden="true"></div>
          <div class="px-5 pt-5 pb-4 border-b border-gray-200 dark:border-gray-700/60">
            <header class="flex items-center">
              <div class="w-6 h-6 rounded-full shrink-0 bg-indigo-500 mr-3">
                <svg class="w-6 h-6 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M7.5 21L3 16.5m0 0L7.5 12M3 16.5h13.5m0-13.5L21 7.5m0 0L16.5 12M21 7.5H7.5"/>
                </svg>
              </div>
              <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold">MQTT</h3>
            </header>
          </div>
          <div class="px-4 py-4 space-y-3">
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium mb-2">Type</label>
                <select v-model="inputDict.MQTT.Type" class="form-select w-full">
                  <option value="public">Public</option>
                  <option value="private">Private</option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-medium mb-2">Host</label>
                <input v-model="inputDict.MQTT.host" class="form-input w-full font-mono" type="text" />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium mb-2">Port</label>
                <input v-model="inputDict.MQTT.port" class="form-input w-full font-mono" type="text" />
              </div>
              <div>
                <label class="block text-sm font-medium mb-2">Device ID</label>
                <input v-model="inputDict.MQTT.device_id" class="form-input w-full font-mono" type="text" />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium mb-2">Username</label>
                <input v-model="inputDict.MQTT.username" class="form-input w-full" type="text" />
              </div>
              <div>
                <label class="block text-sm font-medium mb-2">Password</label>
                <input v-model="inputDict.MQTT.password" class="form-input w-full" type="password" />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium mb-2">External Port</label>
                <input v-model="inputDict.MQTT.externalport" class="form-input w-full font-mono" type="text" />
              </div>
              <div>
                <label class="block text-sm font-medium mb-2">External URL</label>
                <input v-model="inputDict.MQTT.url" class="form-input w-full font-mono" type="text" />
              </div>
            </div>
          </div>
        </div>

      </div>

    </div>
  </div>
</template>

<script setup>
import { inject, computed } from 'vue'
import { getUtcOffsetString } from '@/utils/Utils'

const inputDict = inject('inputDict')

const utcOffsetLabel = computed(() => {
  const tz = inputDict.value?.deviceInfo?.timezone
  if (!tz) return ''
  return getUtcOffsetString(tz)
})

const timezones = [
  "Africa/Abidjan","Africa/Accra","Africa/Addis_Ababa","Africa/Algiers","Africa/Asmara","Africa/Bamako","Africa/Bangui","Africa/Banjul","Africa/Bissau","Africa/Blantyre","Africa/Brazzaville","Africa/Bujumbura","Africa/Cairo","Africa/Casablanca","Africa/Ceuta","Africa/Conakry","Africa/Dakar","Africa/Dar_es_Salaam","Africa/Djibouti","Africa/Douala","Africa/El_Aaiun","Africa/Freetown","Africa/Gaborone","Africa/Harare","Africa/Johannesburg","Africa/Juba","Africa/Kampala","Africa/Khartoum","Africa/Kigali","Africa/Kinshasa","Africa/Lagos","Africa/Libreville","Africa/Lome","Africa/Luanda","Africa/Lubumbashi","Africa/Lusaka","Africa/Malabo","Africa/Maputo","Africa/Maseru","Africa/Mbabane","Africa/Mogadishu","Africa/Monrovia","Africa/Nairobi","Africa/Ndjamena","Africa/Niamey","Africa/Nouakchott","Africa/Ouagadougou","Africa/Porto-Novo","Africa/Sao_Tome","Africa/Tripoli","Africa/Tunis","Africa/Windhoek",
  "America/Adak","America/Anchorage","America/Anguilla","America/Antigua","America/Araguaina","America/Argentina/Buenos_Aires","America/Argentina/Catamarca","America/Argentina/Cordoba","America/Argentina/Jujuy","America/Argentina/La_Rioja","America/Argentina/Mendoza","America/Argentina/Rio_Gallegos","America/Argentina/Salta","America/Argentina/San_Juan","America/Argentina/San_Luis","America/Argentina/Tucuman","America/Argentina/Ushuaia","America/Aruba","America/Asuncion","America/Atikokan","America/Bahia","America/Bahia_Banderas","America/Barbados","America/Belem","America/Belize","America/Blanc-Sablon","America/Boa_Vista","America/Bogota","America/Boise","America/Buenos_Aires","America/Cambridge_Bay","America/Campo_Grande","America/Cancun","America/Caracas","America/Cayenne","America/Cayman","America/Chicago","America/Chihuahua","America/Costa_Rica","America/Cuiaba","America/Curacao","America/Danmarkshavn","America/Dawson","America/Dawson_Creek","America/Denver","America/Detroit","America/Dominica","America/Edmonton","America/Eirunepe","America/El_Salvador","America/Fort_Nelson","America/Fortaleza","America/Glace_Bay","America/Goose_Bay","America/Grand_Turk","America/Grenada","America/Guadeloupe","America/Guatemala","America/Guayaquil","America/Guyana","America/Halifax","America/Havana","America/Hermosillo","America/Indiana/Indianapolis","America/Indiana/Knox","America/Indiana/Marengo","America/Indiana/Petersburg","America/Indiana/Tell_City","America/Indiana/Vevay","America/Indiana/Vincennes","America/Indiana/Winamac","America/Inuvik","America/Iqaluit","America/Jamaica","America/Juneau","America/Kentucky/Louisville","America/Kentucky/Monticello","America/La_Paz","America/Lima","America/Los_Angeles","America/Maceio","America/Managua","America/Manaus","America/Martinique","America/Matamoros","America/Mazatlan","America/Menominee","America/Merida","America/Metlakatla","America/Mexico_City","America/Miquelon","America/Moncton","America/Monterrey","America/Montevideo","America/Montreal","America/Montserrat","America/Nassau","America/New_York","America/Nipigon","America/Nome","America/Noronha","America/Nuuk","America/Ojinaga","America/Panama","America/Pangnirtung","America/Paramaribo","America/Phoenix","America/Port-au-Prince","America/Port_of_Spain","America/Porto_Velho","America/Puerto_Rico","America/Punta_Arenas","America/Rainy_River","America/Rankin_Inlet","America/Recife","America/Regina","America/Resolute","America/Rio_Branco","America/Santiago","America/Santo_Domingo","America/Sao_Paulo","America/Scoresbysund","America/Sitka","America/St_Johns","America/Swift_Current","America/Tegucigalpa","America/Thule","America/Thunder_Bay","America/Tijuana","America/Toronto","America/Tortola","America/Vancouver","America/Whitehorse","America/Winnipeg","America/Yakutat","America/Yellowknife",
  "Antarctica/Casey","Antarctica/Davis","Antarctica/DumontDUrville","Antarctica/Macquarie","Antarctica/Mawson","Antarctica/McMurdo","Antarctica/Palmer","Antarctica/Rothera","Antarctica/South_Pole","Antarctica/Syowa","Antarctica/Troll","Antarctica/Vostok",
  "Arctic/Longyearbyen",
  "Asia/Aden","Asia/Almaty","Asia/Amman","Asia/Anadyr","Asia/Aqtau","Asia/Aqtobe","Asia/Ashgabat","Asia/Atyrau","Asia/Baghdad","Asia/Bahrain","Asia/Baku","Asia/Bangkok","Asia/Barnaul","Asia/Beirut","Asia/Bishkek","Asia/Brunei","Asia/Chita","Asia/Choibalsan","Asia/Colombo","Asia/Damascus","Asia/Dhaka","Asia/Dili","Asia/Dubai","Asia/Dushanbe","Asia/Famagusta","Asia/Gaza","Asia/Hebron","Asia/Ho_Chi_Minh","Asia/Hong_Kong","Asia/Hovd","Asia/Irkutsk","Asia/Istanbul","Asia/Jakarta","Asia/Jayapura","Asia/Jerusalem","Asia/Kabul","Asia/Kamchatka","Asia/Karachi","Asia/Kathmandu","Asia/Khandyga","Asia/Kolkata","Asia/Krasnoyarsk","Asia/Kuala_Lumpur","Asia/Kuching","Asia/Kuwait","Asia/Macau","Asia/Magadan","Asia/Makassar","Asia/Manila","Asia/Muscat","Asia/Nicosia","Asia/Novokuznetsk","Asia/Novosibirsk","Asia/Omsk","Asia/Oral","Asia/Phnom_Penh","Asia/Pontianak","Asia/Pyongyang","Asia/Qatar","Asia/Qostanay","Asia/Qyzylorda","Asia/Riyadh","Asia/Sakhalin","Asia/Samarkand","Asia/Seoul","Asia/Shanghai","Asia/Singapore","Asia/Srednekolymsk","Asia/Taipei","Asia/Tashkent","Asia/Tbilisi","Asia/Tehran","Asia/Thimphu","Asia/Tokyo","Asia/Tomsk","Asia/Ulaanbaatar","Asia/Urumqi","Asia/Ust-Nera","Asia/Vientiane","Asia/Vladivostok","Asia/Yakutsk","Asia/Yangon","Asia/Yekaterinburg","Asia/Yerevan",
  "Atlantic/Azores","Atlantic/Bermuda","Atlantic/Canary","Atlantic/Cape_Verde","Atlantic/Faroe","Atlantic/Madeira","Atlantic/Reykjavik","Atlantic/South_Georgia","Atlantic/St_Helena","Atlantic/Stanley",
  "Australia/Adelaide","Australia/Brisbane","Australia/Broken_Hill","Australia/Darwin","Australia/Eucla","Australia/Hobart","Australia/Lindeman","Australia/Lord_Howe","Australia/Melbourne","Australia/Perth","Australia/Sydney",
  "CET","CST6CDT","EET","EST","EST5EDT","Etc/GMT","Etc/GMT+1","Etc/GMT+10","Etc/GMT+11","Etc/GMT+12","Etc/GMT+2","Etc/GMT+3","Etc/GMT+4","Etc/GMT+5","Etc/GMT+6","Etc/GMT+7","Etc/GMT+8","Etc/GMT+9","Etc/GMT-1","Etc/GMT-10","Etc/GMT-11","Etc/GMT-12","Etc/GMT-13","Etc/GMT-14","Etc/GMT-2","Etc/GMT-3","Etc/GMT-4","Etc/GMT-5","Etc/GMT-6","Etc/GMT-7","Etc/GMT-8","Etc/GMT-9","Etc/UTC",
  "Europe/Amsterdam","Europe/Andorra","Europe/Astrakhan","Europe/Athens","Europe/Belgrade","Europe/Berlin","Europe/Bratislava","Europe/Brussels","Europe/Bucharest","Europe/Budapest","Europe/Busingen","Europe/Chisinau","Europe/Copenhagen","Europe/Dublin","Europe/Gibraltar","Europe/Guernsey","Europe/Helsinki","Europe/Isle_of_Man","Europe/Istanbul","Europe/Jersey","Europe/Kaliningrad","Europe/Kyiv","Europe/Kirov","Europe/Lisbon","Europe/Ljubljana","Europe/London","Europe/Luxembourg","Europe/Madrid","Europe/Malta","Europe/Mariehamn","Europe/Minsk","Europe/Monaco","Europe/Moscow","Europe/Nicosia","Europe/Oslo","Europe/Paris","Europe/Podgorica","Europe/Prague","Europe/Riga","Europe/Rome","Europe/Samara","Europe/San_Marino","Europe/Sarajevo","Europe/Saratov","Europe/Simferopol","Europe/Skopje","Europe/Sofia","Europe/Stockholm","Europe/Tallinn","Europe/Tirane","Europe/Ulyanovsk","Europe/Vaduz","Europe/Vatican","Europe/Vienna","Europe/Vilnius","Europe/Volgograd","Europe/Warsaw","Europe/Zagreb","Europe/Zurich",
  "GMT","HST","MET","MST","MST7MDT","PST8PDT","WET","UTC",
  "Indian/Antananarivo","Indian/Chagos","Indian/Christmas","Indian/Cocos","Indian/Comoro","Indian/Kerguelen","Indian/Mahe","Indian/Maldives","Indian/Mauritius","Indian/Mayotte","Indian/Reunion",
  "Pacific/Apia","Pacific/Auckland","Pacific/Bougainville","Pacific/Chatham","Pacific/Chuuk","Pacific/Easter","Pacific/Efate","Pacific/Fakaofo","Pacific/Fiji","Pacific/Funafuti","Pacific/Galapagos","Pacific/Gambier","Pacific/Guadalcanal","Pacific/Guam","Pacific/Honolulu","Pacific/Kanton","Pacific/Kiritimati","Pacific/Kosrae","Pacific/Kwajalein","Pacific/Majuro","Pacific/Marquesas","Pacific/Midway","Pacific/Nauru","Pacific/Niue","Pacific/Norfolk","Pacific/Noumea","Pacific/Pago_Pago","Pacific/Palau","Pacific/Pitcairn","Pacific/Pohnpei","Pacific/Port_Moresby","Pacific/Rarotonga","Pacific/Tahiti","Pacific/Tarawa","Pacific/Tongatapu","Pacific/Wake","Pacific/Wallis",
  "US/Alaska","US/Aleutian","US/Arizona","US/Central","US/East-Indiana","US/Eastern","US/Hawaii","US/Indiana-Starke","US/Michigan","US/Mountain","US/Pacific","US/Samoa",
]
</script>
