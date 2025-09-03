<template>
  <div class="py-3">
    <!-- 데이터 로딩 상태 표시 -->
    <div v-if="!data || Object.keys(data).length === 0" class="text-center py-8">
      <div class="text-gray-500 dark:text-gray-400">
        <svg class="animate-spin h-8 w-8 mx-auto mb-2" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        Loading data for Channel {{ channel }}...
      </div>
    </div>
    
    <!-- CB 카드 그리드 -->
    <div v-else class="grid grid-cols-12 gap-4">
      <div
        v-for="group in groupedData"
        :key="group.key"
        class="col-span-12 sm:col-span-6 xl:col-span-4"
      >
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-4">
          <!-- CB 헤더 정보 -->
          <div class="border-b border-gray-200 dark:border-gray-700/60 pb-3 mb-3">
            <div class="flex items-center justify-between">
              <div class="flex items-center">
                <svg v-if="group.notiType === 'alarm'" class="shrink-0 fill-current text-yellow-500 mr-2" width="16" height="16" viewBox="0 0 16 16">
                  <path d="M8 0C3.6 0 0 3.6 0 8s3.6 8 8 8 8-3.6 8-8-3.6-8-8-8zm0 12c-.6 0-1-.4-1-1s.4-1 1-1 1 .4 1 1-.4 1-1 1zm1-3H7V4h2v5z" />
                </svg>
                <svg v-else-if="group.notiType === 'error'" class="shrink-0 fill-current text-red-500 mr-2" width="16" height="16">
                  <path d="M8 0C3.6 0 0 3.6 0 8s3.6 8 8 8 8-3.6 8-8-3.6-8-8-8zm3.5 10.1l-1.4 1.4L8 9.4l-2.1 2.1-1.4-1.4L6.6 8 4.5 5.9l1.4-1.4L8 6.6l2.1-2.1 1.4 1.4L9.4 8l2.1 2.1z" />
                </svg>
                <svg v-else class="shrink-0 fill-current text-green-500 mr-2" width="16" height="16">
                  <path d="M8 0C3.6 0 0 3.6 0 8s3.6 8 8 8 8-3.6 8-8-3.6-8-8-8zM7 11.4L3.6 8 5 6.6l2 2 4-4L12.4 6 7 11.4z" />
                </svg>
                <div class="font-semibold text-gray-800 dark:text-gray-100">
                  {{ group.title }}
                </div>
              </div>
              <div class="text-xs font-semibold text-gray-400 dark:text-gray-500 uppercase">
                {{ group.statusText }}
              </div>
            </div>
          </div>

          <!-- 3상인 경우 -->
          <div v-if="group.is3Phase" class="space-y-3">
            <!-- Current 섹션 -->
            <div>
              <div class="text-xs font-bold text-black-400 dark:text-black-500 uppercase mb-1">Current</div>
              <div class="flex flex-wrap gap-2">
                <div class="flex items-center">
                  <div class="text-xs font-medium text-blue-700 px-1.5 bg-blue-500/20 rounded-full">L1</div>
                  <div class="text-sm font-bold text-gray-800 dark:text-gray-100 ml-1">{{ (group.items[0].data.irms && group.items[0].data.irms[0] ? group.items[0].data.irms[0] : 0).toFixed(2) }} A</div>
                </div>
                <div class="flex items-center">
                  <div class="text-xs font-medium text-blue-700 px-1.5 bg-blue-500/20 rounded-full">L2</div>
                  <div class="text-sm font-bold text-gray-800 dark:text-gray-100 ml-1">{{ (group.items[0].data.irms && group.items[0].data.irms[1] ? group.items[0].data.irms[1] : 0).toFixed(2) }} A</div>
                </div>
                <div class="flex items-center">
                  <div class="text-xs font-medium text-blue-700 px-1.5 bg-blue-500/20 rounded-full">L3</div>
                  <div class="text-sm font-bold text-gray-800 dark:text-gray-100 ml-1">{{ (group.items[0].data.irms && group.items[0].data.irms[2] ? group.items[0].data.irms[2] : 0).toFixed(2) }} A</div>
                </div>
              </div>
            </div>

            <!-- Power 섹션 -->
            <div>
              <div class="text-xs font-bold text-black-400 dark:text-black-500 uppercase mb-1">Power</div>
              <div class="flex flex-wrap gap-2">
                <div class="flex items-center">
                  <div class="text-xs font-medium text-green-700 px-1.5 bg-green-500/20 rounded-full">P1</div>
                  <div class="text-sm font-bold text-gray-800 dark:text-gray-100 ml-1">{{ ((group.items[0].data.power && group.items[0].data.power[0] ? group.items[0].data.power[0] : 0) / 1000).toFixed(2) }} kW</div>
                </div>
                <div class="flex items-center">
                  <div class="text-xs font-medium text-green-700 px-1.5 bg-green-500/20 rounded-full">P2</div>
                  <div class="text-sm font-bold text-gray-800 dark:text-gray-100 ml-1">{{ ((group.items[0].data.power && group.items[0].data.power[1] ? group.items[0].data.power[1] : 0) / 1000).toFixed(2) }} kW</div>
                </div>
                <div class="flex items-center">
                  <div class="text-xs font-medium text-green-700 px-1.5 bg-green-500/20 rounded-full">P3</div>
                  <div class="text-sm font-bold text-gray-800 dark:text-gray-100 ml-1">{{ ((group.items[0].data.power && group.items[0].data.power[2] ? group.items[0].data.power[2] : 0) / 1000).toFixed(2) }} kW</div>
                </div>
              </div>
            </div>

            <!-- 기타 측정값들 -->
            <div class="grid grid-cols-2 gap-2">
              <div>
                <div class="text-xs font-bold text-black-400 dark:text-black-500 uppercase mb-1">THD</div>
                <div class="text-sm font-bold text-gray-800 dark:text-gray-100">{{ group.items[0].data.pthd.toFixed(2) }} %</div>
              </div>
              <div>
                <div class="text-xs font-bold text-black-400 dark:text-black-500 uppercase mb-1">Power Factor</div>
                <div class="text-sm font-bold text-gray-800 dark:text-gray-100">{{ group.items[0].data.pf.toFixed(1) }} %</div>
              </div>
              <div>
                <div class="text-xs font-bold text-black-400 dark:text-black-500 uppercase mb-1">Unbalance</div>
                <div class="text-sm font-bold text-gray-800 dark:text-gray-100">{{ ((group.items[0].data.iunbal || 0) / 100).toFixed(1) }} %</div>
              </div>
              <div>
                <div class="text-xs font-bold text-black-400 dark:text-black-500 uppercase mb-1">Temperature</div>
                <div class="text-sm font-bold text-gray-800 dark:text-gray-100">{{ group.items[0].data.temp.toFixed(1) }} ℃</div>
              </div>
              <div>
                <div class="text-xs font-bold text-black-400 dark:text-black-500 uppercase mb-1">Total Energy</div>
                <div class="text-sm font-bold text-gray-800 dark:text-gray-100">{{ (group.items[0].data.kwh / 1000).toFixed(2) }} kWh</div>
              </div>
              <div v-if="group.items[0].data.cbtype === 5">
                <div class="text-xs font-bold text-black-400 dark:text-black-500 uppercase mb-1">Ground</div>
                <div class="text-sm font-bold text-gray-800 dark:text-gray-100">{{ group.items[0].data.ig.toFixed(2) }} mA</div>
              </div>
            </div>
          </div>

          <!-- 단상인 경우 -->
          <div v-else class="space-y-3">
            <!-- Current 정보 -->
            <div>
              <div class="text-xs font-bold text-black-400 dark:text-black-500 uppercase mb-1">Current</div>
              <div class="flex items-center gap-2">
                <span class="text-xs font-medium text-blue-700 px-1.5 bg-blue-500/20 rounded-full">
                  {{ group.phaseLabel }}
                </span>
                <span class="text-sm font-bold text-gray-800 dark:text-gray-100">
                  {{ group.items[0].data.irms[0].toFixed(2) }} A
                </span>
                <span v-if="group.items[0].data.cbtype > 5" class="flex items-center ml-2">
                  <span class="text-xs font-medium text-blue-700 px-1.5 bg-blue-500/20 rounded-full">Ig</span>
                  <span class="text-sm font-bold text-gray-800 dark:text-gray-100 ml-1">{{ group.items[0].data.ig.toFixed(2) }} mA</span>
                </span>
              </div>
            </div>

            <!-- Power 정보 -->
            <div>
              <div class="text-xs font-bold text-black-400 dark:text-black-500 uppercase mb-1">Power</div>
              <div class="flex items-center gap-2">
                <span class="text-xs font-medium text-green-700 px-1.5 bg-green-500/20 rounded-full">
                  P{{ group.phaseIndex + 1 }}
                </span>
                <span class="text-sm font-bold text-gray-800 dark:text-gray-100">
                  {{ (group.items[0].data.power[0] / 1000).toFixed(2) }} kW
                </span>
              </div>
            </div>

            <!-- 측정값들 -->
            <div class="grid grid-cols-2 gap-2">
              <div>
                <div class="text-xs font-bold text-black-400 dark:text-black-500 uppercase mb-1">THD</div>
                <div class="text-sm font-bold text-gray-800 dark:text-gray-100">{{ group.items[0].data.pthd.toFixed(2) }} %</div>
              </div>
              <div>
                <div class="text-xs font-bold text-black-400 dark:text-black-500 uppercase mb-1">Power Factor</div>
                <div class="text-sm font-bold text-gray-800 dark:text-gray-100">{{ group.items[0].data.pf.toFixed(1) }} %</div>
              </div>
              <div>
                <div class="text-xs font-bold text-black-400 dark:text-black-500 uppercase mb-1">Unbalance</div>
                <div class="text-sm font-bold text-gray-800 dark:text-gray-100">{{ ((group.items[0].data.iunbal || 0) / 100).toFixed(1) }} %</div>
              </div>
              <div>
                <div class="text-xs font-bold text-black-400 dark:text-black-500 uppercase mb-1">Temperature</div>
                <div class="text-sm font-bold text-gray-800 dark:text-gray-100">{{ group.items[0].data.temp.toFixed(1) }} ℃</div>
              </div>
              <div>
                <div class="text-xs font-bold text-black-400 dark:text-black-500 uppercase mb-1">Energy</div>
                <div class="text-sm font-bold text-gray-800 dark:text-gray-100">{{ (group.items[0].data.kwh / 1000).toFixed(2) }} kWh</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, watch, computed } from 'vue'

export default {
  name: 'GemsDashCard',
  props: {
    channel: String,
    data: Object,
  },
  setup(props) {
    const channel = computed(() => props.channel);
    const data = computed(() => props.data);
    
    // 데이터를 그룹핑하는 computed
    const groupedData = computed(() => {
      if (!data.value || Object.keys(data.value).length === 0) return [];
      
      const groups = [];
      let ctNumber = 1; // CT 번호 카운터
      
      // 데이터를 키 순서대로 정렬
      const sortedKeys = Object.keys(data.value).sort((a, b) => parseInt(a) - parseInt(b));
      
      for (const key of sortedKeys) {
        const cbData = data.value[key];
        const cbtype = cbData.cbtype;
        
        // 3상인 경우 (3P3W, 3P4W) - cblist가 3개
        if (cbtype === 4 || cbtype === 5) {
          const ctStart = ctNumber;
          const ctEnd = ctNumber + 2;
          
          groups.push({
            key: `3p-${key}`,
            title: `CT${ctStart}-${ctEnd} (${cbData.cbtype_text})`,
            is3Phase: true,
            items: [{ id: key, data: cbData }],
            notiType: getNotiType(cbData.status),
            statusText: getStatusText(cbData.status)
          });
          
          // CT 번호 3개 증가
          ctNumber += 3;
        } 
        // 단상인 경우 - cblist가 1개
        else {
          const phaseInfo = getPhaseInfo(cbData.cbtype);
          
          groups.push({
            key: `single-${key}`,
            title: `CT${ctNumber} (${cbData.cbtype_text})`,
            redisKey: key,
            is3Phase: false,
            items: [{
              id: key,
              data: cbData
            }],
            phaseLabel: phaseInfo.label,
            phaseIndex: phaseInfo.index,
            notiType: getNotiType(cbData.status),
            statusText: getStatusText(cbData.status)
          });
          
          // CT 번호 1개 증가
          ctNumber++;
        }
      }
      
      return groups;
    });
    
    // 상 정보 가져오기
    const getPhaseInfo = (cbtype) => {
      switch(cbtype) {
        case 1: return { label: 'L1', index: 0 };
        case 2: return { label: 'L2', index: 1 };
        case 3: return { label: 'L3', index: 2 };
        case 6: return { label: 'L1', index: 0 }; // L1+Z
        case 7: return { label: 'L2', index: 1 }; // L2+Z
        case 8: return { label: 'L3', index: 2 }; // L3+Z
        default: return { label: 'L1', index: 0 };
      }
    };

    // 상태 확인
    const getNotiType = (status) => {
      const devSt = status & 0xfffffc0e;
      const comSt = status & 0x01;
      
      if (comSt === 0) {
        return 'error';
      } else {
        if(devSt === 0) {
          return 'success';
        } else {
          return 'alarm';
        }
      }
    };

    // 상태 텍스트 가져오기
    const getStatusText = (status) => {
      const devSt = status & 0xfffffc0e;
      const comSt = status & 0x01;
      
      if (comSt === 0) return 'COMM ERROR';
      if (devSt === 0) return 'ONLINE';
      return 'ALARM';
    };

    return {
      channel,
      data,
      groupedData,
    }
  }
}
</script>