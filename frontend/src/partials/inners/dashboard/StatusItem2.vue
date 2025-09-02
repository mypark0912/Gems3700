<template>
    <div class="bg-white shadow rounded-xl px-5 py-4 dark:bg-gray-700">
  <h3 class="text-sm font-semibold text-gray-800 mb-2 dark:text-gray-200"  @click="goToEquipmentDetail"> {{ stData.devNickname }}</h3>
  <!---h3
  @click="goToEquipmentDetail"
  class="text-sm font-semibold text-blue-600 dark:text-blue-400 cursor-pointer hover:underline hover:text-blue-800 transition-colors"
  >
  {{ t('dashboard.diagnosis.equipment') }} : {{ stData.devNickname }}
</h3-->
 <div class="flex flex-col">
    <div class="flex items-center space-x-4">
  </div>
    <div v-if="stData.devType.includes('Transformer')" class="flex items-center py-2">
        <!-- 이미지 -->
        <div class="w-12 h-12 shrink-0">
        <img class="w-12 h-12 rounded-full" :src="motorImageSrc" />
        </div>
        <div class="mr-5 ml-2">
            <div class="flex items-center">
                <div class="text-xl font-bold text-gray-800 mr-2 dark:text-gray-200">{{ transData.Temp.toFixed(2) }} ℃ </div>
            </div>
            <div class="text-sm text-gray-500 dark:text-gray-200">Temperature</div>
        </div>
        <div class="hidden md:block w-px h-8 bg-gray-200 mr-5" aria-hidden="true"></div>
        <div class="mr-5">
            <div class="flex items-center">
                <div class="text-xl font-bold text-gray-800 mr-2 dark:text-gray-200">{{ LoadRate }} %</div>
            </div>
            <div class="text-sm text-gray-500 dark:text-gray-200">Load Factor</div>
        </div>
        <div class="hidden md:block w-px h-8 bg-gray-200 mr-5" aria-hidden="true"></div>
        <div class="mr-5">
            <div class="flex items-center">
                <div class="text-xl font-bold text-gray-800 mr-2 dark:text-gray-200"> {{ transData.Ig.toFixed(2) }} A </div>
            </div>
            <div class="text-sm text-gray-500 dark:text-gray-200">Ig</div>
        </div>
    </div>
    <div v-else-if="stData.devType=='VFD'" class="flex items-center py-2">
        <!-- 이미지 -->
        <div class="w-12 h-12 shrink-0">
        <img class="w-12 h-12 rounded-full" :src="motorImageSrc" />
        </div>
        <div class="mr-5 ml-2">
            <div class="flex items-center">
                <div class="text-xl font-bold text-gray-800 mr-2">{{ transData.Temp.toFixed(2) }} ℃ </div>
            </div>
            <div class="text-sm text-gray-500">Temperature</div>
        </div>
        <div class="hidden md:block w-px h-8 bg-gray-200 mr-5" aria-hidden="true"></div>
        <div class="mr-5">
            <div class="flex items-center">
                <div class="text-xl font-bold text-gray-800 mr-2">56.9 %</div>
            </div>
            <div class="text-sm text-gray-500">Load Factor</div>
        </div>
        <div class="hidden md:block w-px h-8 bg-gray-200 mr-5" aria-hidden="true"></div>
        <div class="mr-5">
            <div class="flex items-center">
                <div class="text-xl font-bold text-gray-800 mr-2"> {{ transData.Ig.toFixed(2) }} A </div>
            </div>
            <div class="text-sm text-gray-500">Ig</div>
        </div>
    </div>
    <div v-else class="flex items-center py-2">
        <!-- 이미지 -->
        <div class="w-12 h-12 shrink-0">
        <img class="w-12 h-12 rounded-full" :src="motorImageSrc" />
        </div>
        <div class="mr-5 ml-2">
            <div class="flex items-center">
                <div class="text-xl font-bold text-gray-800 mr-2 dark:text-gray-200">132 </div>
            </div>
            <div class="text-sm text-gray-500 dark:text-gray-200">Running Hours</div>
        </div>
        <div class="hidden md:block w-px h-8 bg-gray-200 mr-5" aria-hidden="true"></div>
        <div class="mr-5 ml-2">
            <div class="flex items-center">
                <div class="text-xl font-bold text-gray-800 mr-2 dark:text-gray-200">{{ transData[0]["Value"].toFixed(2) }} {{ transData[0]["Unit"] }} </div>
            </div>
            <div class="text-sm text-gray-500 dark:text-gray-200">{{ transData[0]["Assembly"] }} : {{ transData[0]["Title"] }}</div>
        </div>
        <div class="hidden md:block w-px h-8 bg-gray-200 mr-5" aria-hidden="true"></div>
            <div v-if="transData.length > 2" class="mr-5 ml-2">
                <div class="flex items-center">
                    <div class="text-xl font-bold text-gray-800 mr-2 dark:text-gray-200">{{ transData[1]["Value"].toFixed(2) }} {{ transData[1]["Unit"] }} </div>
                </div>
                <div class="text-sm text-gray-500 dark:text-gray-200">{{ transData[1]["Assembly"] }} : {{ transData[1]["Title"] }}</div>
            </div>
            <div v-if="transData.length > 2" class="hidden md:block w-px h-8 bg-gray-200 mr-5" aria-hidden="true"></div>
            <div v-if="transData.length > 2" class="mr-5">
                <div class="flex items-center">
                    <div class="text-xl font-bold text-gray-800 mr-2 dark:text-gray-200">{{ transData[2]["Value"].toFixed(2) }} %</div>
                </div>
                <div class="text-sm text-gray-500 dark:text-gray-200">{{ transData[2]["Assembly"] }} : {{ transData[2]["Title"] }}</div>
            </div>
            <div v-else class="mr-5 ml-2">
                <div class="flex items-center">
                    <div class="text-xl font-bold text-gray-800 mr-2 dark:text-gray-200">{{ transData[1].Value.toFixed(2) }} {{ transData[1].Unit }} </div>
                </div>
                <div class="text-sm text-gray-500 dark:text-gray-200">{{ transData[1].Assembly }} : {{ transData[1].Title }}</div>
            </div>
    </div>
    <!--div class="flex items-center py-2">
      <div class="mr-5">
            <div class="flex items-center space-x-1">
                <span
                class="text-sm rounded-full px-3 py-1 min-w-[100px] w-fit text-center whitespace-nowrap"
                :class="statusClass"
                >
                {{ statusStr }}
                </span>
                  <span
                    v-for="(item, idx) in strList"
                    :key="idx"
                    class="text-sm rounded-full px-3 py-1 min-w-[100px] w-fit text-center whitespace-nowrap transition-all"
                    :class="stData.devStatus === idx ? item.css : 'bg-gray-200 text-gray-400'"
                  >
                    {{ item.text }}
                  </span>
            </div>
        </div>
    </div-->
 </div>
</div>

  </template>
  
  <script>
import motorImg from '@/images/motor_m.png'
import fanImg from '@/images/fan_m.png'
import pumpImg from '@/images/pump_m.png'
import compImg from '@/images/comp_m.png'
import powerImg from '@/images/power_m.png'
import defaultImg from '@/images/cleaned_logo.png'
import transImg from '@/images/trans.png'
import { ref, watch, computed, onMounted } from 'vue';
import { useI18n } from 'vue-i18n'  // ✅ 추가
import { useSetupStore } from '@/store/setup'
import { useRouter } from 'vue-router'
  export default {
    name: 'StatusItem',
    props: {
        data: Object,
        channel: String,
        transData:Object,
    },
    setup(props){
      const { t } = useI18n();
      const setupStroe = useSetupStore();
      const route = useRouter();
        const channel = ref(props.channel);
        const computedChannel = computed(()=>{
          if (channel.value == 'Main'|| channel.value == 'main')
            return 'Main';
          else
            return 'Sub';
        })
        const stData = ref(props.data);
        const transData = ref(props.transData);
        const LoadRate = ref(0);
        const strList = computed(()=>{
          return [{"text":t('dashboard.diagnosis.st0'), "css":'bg-gray-500/20 text-gray-700 font-semibold'},
          {"text":t('dashboard.diagnosis.st1'), "css":'bg-green-500/20 text-green-700 font-semibold'},
          {"text":t('dashboard.diagnosis.st2'), "css":'bg-yellow-500/20 text-yellow-700 font-semibold'},
          {"text":t('dashboard.diagnosis.st3'), "css":'bg-orange-500/20 text-orange-700 font-semibold'},
          {"text":t('dashboard.diagnosis.st4'), "css":'bg-red-500/20 text-red-700 font-semibold'}
          ]

        })
        const motorImageSrc = computed(() => {
            switch (stData.value.devType) {
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

        const goToEquipmentDetail = () => {
            route.push(`/diagnosis/${computedChannel.value}`)  // devId나 다른 식별자 기반 경로
          }

        const statusClass = computed(() => {
      switch (stData.value.devStatus) {
        case 1:
          return 'bg-green-500/20 text-green-700 font-semibold';
        case 2:
        return 'bg-yellow-500/20 text-yellow-700 font-semibold';
        case 3:
          return 'bg-orange-500/20 text-orange-700 font-semibold';
        case 4:
          return 'bg-red-500/20 text-red-700 font-semibold';
        default:
          return 'bg-gray-500/20 text-gray-700 font-semibold';
      }
    });

    const LoadFactor = computed(()=>{
      let kva = -1;
      if(computedChannel.value == 'Main' && stData.value.devType=='Transformer'){
        kva = setupStroe.getMkva;
      }
      if(computedChannel.value == 'sub' && stData.value.devType=='Transformer'){
        kva = setupStroe.getSkva;
      }
      return kva;
    })

    watch(
      () => LoadFactor.value, // ✅ 이렇게 해야 실제 값이 감지됨
      (newVal) => {
        if (newVal > 0 && transData.value?.Stotal) {
          LoadRate.value = ((transData.value.Stotal / newVal) * 100).toFixed(2);
        }
      },
      { immediate: true }
    );


    const statusStr =  computed(() => {
      switch (stData.value.devStatus) {
          case 1:
            return t('dashboard.diagnosis.st1'); //'OK';
          case 2:
          return t('dashboard.diagnosis.st2');// 'Warning';
          case 3:
            return t('dashboard.diagnosis.st3'); //'Inspect';
          case 4:
            return t('dashboard.diagnosis.st4');// 'Repair';
          default:
            return t('dashboard.diagnosis.st0'); //'No Data';
        }
    });

        return {
            stData,
            motorImageSrc,
            statusStr,
            statusClass,
            transData,
            t,
            strList,
            LoadFactor,
            channel,
            LoadRate,
            goToEquipmentDetail,
            computedChannel,
        }
    }
  }
  </script>