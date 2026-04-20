<template>
    <div class="col-span-full sm:col-span-6 xl:col-span-3 bg-white dark:bg-gray-800 shadow-sm rounded-xl">
      <div class="flex flex-col h-full">
        <!-- Card top -->
        <div class="grow p-5">
          <!-- Menu button -->
          <!-- Image + name -->
          <header>
            <!--div class="flex justify-center mb-2">
                <svg xmlns="http://www.w3.org/2000/svg" width="128" height="128" viewBox="0 0 128 128">

                  <rect x="16" y="64" width="96" height="24" rx="8" ry="8" fill="#d0d7e1"/>
                  <circle cx="88" cy="76" r="4" fill="#007aff"/>
                  <circle cx="100" cy="76" r="4" fill="#007aff"/>
                  

                  <rect x="28" y="88" width="12" height="8" rx="4" ry="4" fill="#007aff"/>
                  <rect x="88" y="88" width="12" height="8" rx="4" ry="4" fill="#007aff"/>
                  

                  <rect x="58" y="40" width="12" height="28" rx="4" ry="4" fill="#aab2bd"/>

                  <path d="M64 16 q20 0 32 12" fill="none" stroke="#007aff" stroke-width="6" stroke-linecap="round"/>
                  <path d="M64 24 q14 0 24 10" fill="none" stroke="#007aff" stroke-width="6" stroke-linecap="round"/>
                  <path d="M64 32 q8 0 16 8" fill="none" stroke="#007aff" stroke-width="6" stroke-linecap="round"/>
              </svg>
            </div-->
            <div class="flex items-center justify-center gap-4">
                <!-- SVG 아이콘 -->
                 <div v-if="!main_data.st_ico && !sub_data.st_ico">
                <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 128 128">
                  <rect x="16" y="64" width="96" height="24" rx="8" ry="8" fill="#d0d7e1"/>
                  <circle cx="88" cy="76" r="4" fill="#007aff"/>
                  <circle cx="100" cy="76" r="4" fill="#007aff"/>
                  <rect x="28" y="88" width="12" height="8" rx="4" ry="4" fill="#007aff"/>
                  <rect x="88" y="88" width="12" height="8" rx="4" ry="4" fill="#007aff"/>
                  <rect x="58" y="40" width="12" height="28" rx="4" ry="4" fill="#aab2bd"/>
                  <path d="M64 16 q20 0 32 12" fill="none" stroke="#007aff" stroke-width="6" stroke-linecap="round"/>
                  <path d="M64 24 q14 0 24 10" fill="none" stroke="#007aff" stroke-width="6" stroke-linecap="round"/>
                  <path d="M64 32 q8 0 16 8" fill="none" stroke="#007aff" stroke-width="6" stroke-linecap="round"/>
                </svg>
                </div>
                <div v-else-if="!main_data.st_ico && sub_data.st_ico" class="flex items-center gap-4">
                <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 128 128">
                  <rect x="16" y="64" width="96" height="24" rx="8" ry="8" fill="#d0d7e1"/>
                  <circle cx="88" cy="76" r="4" fill="#007aff"/>
                  <circle cx="100" cy="76" r="4" fill="#007aff"/>
                  <rect x="28" y="88" width="12" height="8" rx="4" ry="4" fill="#007aff"/>
                  <rect x="88" y="88" width="12" height="8" rx="4" ry="4" fill="#007aff"/>
                  <rect x="58" y="40" width="12" height="28" rx="4" ry="4" fill="#aab2bd"/>
                  <path d="M64 16 q20 0 32 12" fill="none" stroke="#007aff" stroke-width="6" stroke-linecap="round"/>
                  <path d="M64 24 q14 0 24 10" fill="none" stroke="#007aff" stroke-width="6" stroke-linecap="round"/>
                  <path d="M64 32 q8 0 16 8" fill="none" stroke="#007aff" stroke-width="6" stroke-linecap="round"/>
                </svg>
                <img :src="equipImageSrc" alt="장비 이미지" class="w-14 h-14 object-cover" />
                </div>
                <div v-else-if="main_data.st_ico && !sub_data.st_ico" class="flex items-center gap-4">
                <img :src="equipImageSrc" alt="장비 이미지" class="w-14 h-14 object-cover" />
                <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 128 128">
                  <rect x="16" y="64" width="96" height="24" rx="8" ry="8" fill="#d0d7e1"/>
                  <circle cx="88" cy="76" r="4" fill="#007aff"/>
                  <circle cx="100" cy="76" r="4" fill="#007aff"/>
                  <rect x="28" y="88" width="12" height="8" rx="4" ry="4" fill="#007aff"/>
                  <rect x="88" y="88" width="12" height="8" rx="4" ry="4" fill="#007aff"/>
                  <rect x="58" y="40" width="12" height="28" rx="4" ry="4" fill="#aab2bd"/>
                  <path d="M64 16 q20 0 32 12" fill="none" stroke="#007aff" stroke-width="6" stroke-linecap="round"/>
                  <path d="M64 24 q14 0 24 10" fill="none" stroke="#007aff" stroke-width="6" stroke-linecap="round"/>
                  <path d="M64 32 q8 0 16 8" fill="none" stroke="#007aff" stroke-width="6" stroke-linecap="round"/>
                </svg>                
                </div>
                <div v-else class="flex items-center gap-4">
                <img :src="equipImageSrc" alt="장비 이미지" class="w-14 h-14 object-cover" />
                <img :src="equipImageSrc2" alt="장비 이미지" class="w-14 h-14 object-cover" />
                </div>

                <!-- 텍스트 영역 (수직 정렬) -->
                <div class="flex items-center justify-center gap-4">
                  <div class="flex flex-col">
                    <h3 class="text-xl leading-snug justify-center font-semibold">{{item.name}}</h3>
                    <div class="flex justify-center items-center">
                      <span class="text-xs font-medium text-gray-400 dark:text-gray-500 -mt-0.5 mr-1">-&gt;</span> 
                      <span class="text-xs">{{item.location}}</span>
                    </div>
                  </div>
                </div>
              </div>
            <!--div class="text-center">
                <h2 class="text-xl leading-snug justify-center font-semibold">{{item.name}}</h2>
            </div-->
          </header>
          <!-- Bio -->
          <div class="flex justify-center">
            <div class="grid grid-cols-2 gap-4 mt-2">
            <!-- Main Channel (왼쪽) -->
              <div class="text-left">
                <div class="text-sm font-bold text-black-400 dark:text-black-500 uppercase mb-1">Main Channel</div>
                <div class="flex flex-col gap-1">
                  <div class="flex items-center gap-2">
                    <span class="text-xs font-medium text-violet-700 px-1.5 bg-violet-500/20 rounded-full">V</span>
                    <span class="text-basic font-bold text-gray-800 dark:text-gray-100">{{ Number(main_data.volt).toFixed(2) }} V</span>
                  </div>
                  <div class="flex items-center gap-2">
                    <span class="text-xs font-medium text-sky-700 px-1.5 bg-sky-500/20 rounded-full">C</span>
                    <span class="text-basic font-bold text-gray-800 dark:text-gray-100">{{ Number(main_data.cur).toFixed(2) }} A</span>
                  </div>
                  <div class="flex items-center gap-2">
                    <span class="text-xs font-medium text-lime-700 px-1.5 bg-lime-500/20 rounded-full">F</span>
                    <span class="text-basic font-bold text-gray-800 dark:text-gray-100">{{ Number(main_data.freq) .toFixed(2) }} Hz</span>
                  </div>
                  <div class="flex items-center gap-2">
                    <span class="text-xs font-medium text-pink-700 px-1.5 bg-pink-500/20 rounded-full">D</span>
                    <span :class="getStatusColor(main_data.st)">{{main_data.st}}</span>
                  </div>
                </div>
              </div>
            <!-- Sub Channel (오른쪽) -->
              <div class="text-left">
                <div class="text-sm font-bold text-black-400 dark:text-black-500 uppercase mb-1">Sub Channel</div>
                <div class="flex flex-col gap-1">
                  <div class="flex items-center gap-2">
                    <span class="text-xs font-medium text-violet-700 px-1.5 bg-violet-500/20 rounded-full">V</span>
                    <span class="text-basic font-bold text-gray-800 dark:text-gray-100">{{ Number(sub_data.volt).toFixed(2) }} V</span>
                  </div>
                  <div class="flex items-center gap-2">
                    <span class="text-xs font-medium text-sky-700 px-1.5 bg-sky-500/20 rounded-full">C</span>
                    <span class="text-basic font-bold text-gray-800 dark:text-gray-100">{{ Number(sub_data.cur).toFixed(2) }} A</span>
                  </div>
                  <div class="flex items-center gap-2">
                    <span class="text-xs font-medium text-lime-700 px-1.5 bg-lime-500/20 rounded-full">F</span>
                    <span class="text-basic font-bold text-gray-800 dark:text-gray-100">{{ Number(sub_data.freq).toFixed(2) }} Hz</span>
                  </div>
                  <div class="flex items-center gap-2">
                    <span class="text-xs font-medium text-pink-700 px-1.5 bg-pink-500/20 rounded-full">D</span>
                    <span :class="getStatusColor(sub_data.st)">{{sub_data.st}}</span>
                  </div>
                </div>
              </div>
          </div>
          </div>



        </div>
        <!-- Card footer -->
        <div class="border-t border-gray-100 dark:border-gray-700/60">
          <a v-if="enabled" 
            :href="`http://${item.name}:4000/dashboard/${user}`"
            target="_blank"
            rel="noopener noreferrer"
            class="block text-center text-sm font-medium px-3 py-4 text-violet-500 hover:text-violet-600 dark:hover:text-violet-400">
            <div class="flex items-center justify-center">
              <svg class="w-6 h-6 mr-2" viewBox="0 0 256 256" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
                <path d="M 21,120 L 20,121 L 19,121 L 17,123 L 17,124 L 16,125 L 16,130 L 17,131 L 17,132 L 19,134 L 20,134 L 21,135 L 133,135 L 134,134 L 135,134 L 138,131 L 138,124 L 135,121 L 134,121 L 133,120 Z" />
                <path d="M 143,101 L 142,102 L 141,102 L 139,104 L 139,105 L 138,106 L 138,111 L 139,112 L 139,113 L 149,123 L 149,124 L 150,125 L 150,130 L 149,131 L 149,132 L 139,142 L 139,143 L 138,144 L 138,149 L 139,150 L 139,151 L 141,153 L 142,153 L 143,154 L 148,154 L 149,153 L 150,153 L 171,132 L 171,131 L 172,130 L 172,125 L 171,124 L 171,123 L 150,102 L 149,102 L 148,101 Z" />
                <path d="M 79,28 L 78,29 L 77,29 L 74,32 L 74,34 L 73,35 L 73,88 L 74,89 L 74,91 L 77,94 L 78,94 L 79,95 L 83,95 L 84,94 L 85,94 L 88,91 L 88,44 L 89,43 L 223,43 L 224,44 L 224,211 L 223,212 L 89,212 L 88,211 L 88,164 L 85,161 L 84,161 L 83,160 L 79,160 L 78,161 L 77,161 L 74,164 L 74,166 L 73,167 L 73,220 L 74,221 L 74,223 L 77,226 L 78,226 L 79,227 L 234,227 L 235,226 L 236,226 L 238,224 L 238,223 L 239,222 L 239,33 L 238,32 L 238,31 L 236,29 L 235,29 L 234,28 Z" />
              </svg>
              <span class="text-basic font-semibold">Access</span>
            </div>
          </a>

          <a v-else 
            href="#" 
            class="block text-center text-sm font-medium px-3 py-4 text-gray-400 cursor-not-allowed"
            @click.prevent>
            <div class="flex items-center justify-center">
              <span class="text-basic font-semibold">Unavailable</span>
            </div>
          </a>
        </div>



      </div>
    </div>
  </template>
  
  <script>
  import { ref, computed, watchEffect } from 'vue'
import motorImg from '@/images/motor_m.png'
import fanImg from '@/images/fan_m.png'
import pumpImg from '@/images/pump_m.png'
import compImg from '@/images/comp_m.png'
import powerImg from '@/images/power_m.png'
import defaultImg from '@/images/cleaned_logo.png'
import transImg from '@/images/trans.png'

  export default {
    name: 'DashCard',
    props: ['item','user'],
    setup(props){
      const item = ref(props.item);
      const user = ref(props.user);
      const main_data = computed(()=>{
        if (item.value.main && Object.keys(item.value.main).length > 0) {
            return{
              "volt":parseFloat( item.value.main["v"][0]) + parseFloat( item.value.main["v"][1]) + parseFloat( item.value.main["v"][2])/3,
              "cur": parseFloat( item.value.main["i"][0]) + parseFloat( item.value.main["i"][1]) + parseFloat( item.value.main["i"][2]),
              "freq":item.value.main["freq"],
              "st":getStatusText(item.value.main["status"]),
              "st_ico":item.value.main["status"] == -1? false:true
            }
        }
        return {
          "volt":0,
          "cur":0,
          "freq":0,
          "st":getStatusText(-1),
          "st_ico":false
        }
      })

      const sub_data = computed(()=>{
        if (item.value.sub && Object.keys(item.value.sub).length > 0) {
            return{
              "volt":parseFloat( item.value.sub["v"][0]) + parseFloat( item.value.sub["v"][1]) + parseFloat( item.value.sub["v"][2])/3,
              "cur": parseFloat( item.value.sub["i"][0]) + parseFloat( item.value.sub["i"][1]) + parseFloat( item.value.sub["i"][2]),
              "freq":item.value.sub["freq"],
              "st":getStatusText(item.value.sub["status"]),
              "st_ico":item.value.sub["status"] == -1? false:true
            }
        }
        return {
          "volt":0,
          "cur":0,
          "freq":0,
          "st":getStatusText(-1),
          "st_ico":false
        }
      })
      const enabled = computed(() => item.value["enable"] == 'true'?true:false);
      const getStatusColor = (status) => {
      switch (status) {
        case 'OK': return 'text-basic text-green-700 font-semibold';
        case 'Warning': return 'text-basic text-yellow-700 font-semibold';
        case 'Inspect': return 'text-basic text-orange-700 font-semibold';
        case 'Repair': return 'text-basic text-red-700 font-semibold';
        default: return 'text-basic text-gray-700 font-semibold';
      }
    };

    const equipImageSrc = computed(() => {
      //console.log(item.value.main["st_type"]);
      switch (item.value.main["st_type"]) {
        case 'Motor':
          return motorImg;
        case 'MotorFeed':
          return motorImg;          
        case 'Pump':
          return pumpImg;//'/images/motor_pump.png';
        case 'Fan':
          return fanImg;
        case 'Compressor':
          return compImg;//'/images/fan_m.png';
        case 'PSupply':
          return powerImg;//'/images/power.png';
        case 'PowerSupply':
          return powerImg;//'/images/power.png';  
        case 'PrimaryTransformer':
          return transImg;//'/images/trans.png';       
        case 'Transformer':
          return transImg;//'/images/trans.png';                  
        default:
          return defaultImg;//'@/images/cleaned_logo.png'
      }
    });
    const equipImageSrc2 = computed(() => {
      switch (item.value.sub["st_type"]) {
        case 'Motor':
          return motorImg;
        case 'MotorFeed':
          return motorImg;          
        case 'Pump':
          return pumpImg;//'/images/motor_pump.png';
        case 'Fan':
          return fanImg;
        case 'Compressor':
          return compImg;//'/images/fan_m.png';
        case 'PSupply':
          return powerImg;//'/images/power.png';
        case 'PowerSupply':
          return powerImg;//'/images/power.png';  
        case 'PrimaryTransformer':
          return transImg;//'/images/trans.png';       
        case 'Transformer':
          return transImg;//'/images/trans.png';                  
        default:
          return defaultImg;//'@/images/cleaned_logo.png'
      }
    });

// ✅ 상태 변환 함수
    function getStatusText(status) {
      switch (status) {
        case -1: return 'None';
        case 0: return 'No Data';
        case 1: return 'OK';
        case 2: return 'Warning';
        case 3: return 'Inspect';
        case 4: return 'Repair';
      }
    }

      return {
        item,
        main_data,
        sub_data,
        getStatusColor,
        equipImageSrc,
        equipImageSrc2,
        enabled,
        user,
      }
    }
  }
  </script>