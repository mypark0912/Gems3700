<template>
    <!-- Modal backdrop -->
    <transition
      enter-active-class="transition ease-out duration-200"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition ease-out duration-100"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div v-show="modalOpen" class="fixed inset-0 bg-gray-900 bg-opacity-30 z-50 transition-opacity" aria-hidden="true"></div>
    </transition>
    <!-- Modal dialog -->
    <transition
      enter-active-class="transition ease-in-out duration-200"
      enter-from-class="opacity-0 translate-y-4"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition ease-in-out duration-200"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 translate-y-4"
    >
      <div v-if="modalOpen" :id="id" class="fixed inset-0 z-50 overflow-hidden flex items-center justify-center px-4 sm:px-6" role="dialog" aria-modal="true">
        <div ref="modalContent" class="bg-white dark:bg-gray-800 rounded-lg shadow-lg w-auto min-w-[300px] max-w-screen-lg max-h-full">
          <!-- Modal header -->
          <div class="px-5 py-3 border-b border-gray-200 dark:border-gray-700/60">
            <div class="flex justify-between items-center">
              <div class="flex items-start">
                <svg v-if="notiType === 'alarm'" class="shrink-0 fill-current text-yellow-500 mt-[3px] mr-3" width="16" height="16" viewBox="0 0 16 16">
                  <path d="M8 0C3.6 0 0 3.6 0 8s3.6 8 8 8 8-3.6 8-8-3.6-8-8-8zm0 12c-.6 0-1-.4-1-1s.4-1 1-1 1 .4 1 1-.4 1-1 1zm1-3H7V4h2v5z" />
                </svg>
                <svg v-else-if="notiType === 'error'" class="shrink-0 fill-current text-red-500 mt-[3px] mr-3" width="16" height="16">
                  <path d="M8 0C3.6 0 0 3.6 0 8s3.6 8 8 8 8-3.6 8-8-3.6-8-8-8zm3.5 10.1l-1.4 1.4L8 9.4l-2.1 2.1-1.4-1.4L6.6 8 4.5 5.9l1.4-1.4L8 6.6l2.1-2.1 1.4 1.4L9.4 8l2.1 2.1z" />
                </svg>
                <svg v-else-if="notiType === 'success'" class="shrink-0 fill-current text-green-500 mt-[3px] mr-3" width="16" height="16">
                  <path d="M8 0C3.6 0 0 3.6 0 8s3.6 8 8 8 8-3.6 8-8-3.6-8-8-8zM7 11.4L3.6 8 5 6.6l2 2 4-4L12.4 6 7 11.4z" />
                </svg> 
                <div class="font-semibold text-gray-800 dark:text-gray-100">
                  {{ tabTitle }}
                </div> 
              </div>
              <button class="text-gray-400 dark:text-gray-500 hover:text-gray-500 dark:hover:text-gray-400" @click.stop="$emit('close-modal')">
                <div class="sr-only">Close</div>
                <svg class="fill-current" width="16" height="16" viewBox="0 0 16 16">
                  <path d="M7.95 6.536l4.242-4.243a1 1 0 111.415 1.414L9.364 7.95l4.243 4.242a1 1 0 11-1.415 1.415L7.95 9.364l-4.243 4.243a1 1 0 01-1.414-1.415L6.536 7.95 2.293 3.707a1 1 0 011.414-1.414L7.95 6.536z" />
                </svg>
              </button>
            </div>
          </div>
          <div class="px-4">
            <div class="flex justify-between items-center py-4">
              <div class="text-sm font-bold text-black-400 dark:text-black-500 uppercase">
                {{ subTitle }}
              </div>
              <div class="text-sm font-semibold text-gray-400 dark:text-gray-500 uppercase">
                {{ updated }}
              </div>
            </div>
            <ul class="text-sm font-medium flex flex-nowrap overflow-x-auto no-scrollbar w-full">
                <li v-for="(tab, index) in tabs" :key="index" class="mr-4 last:mr-0 relative">
                  <button
                    @click="changeTab(tab.name)"
                    class="px-4 py-2 whitespace-nowrap transition duration-200 ease-in-out"
                    :class="activeTab === tab.name
                      ? 'text-violet-500 border-b-2 border-violet-500'
                      : 'text-gray-500 dark:text-gray-400 hover:text-gray-600 dark:hover:text-gray-300'">
                    {{ tab.label }}
                  </button>
                </li>
              </ul>
            </div>

                <!-- Tab Content -->
                <div v-for="(tab, index) in tabs" :key="index">
                  <div v-if="activeTab === tab.name" class="text-gray-700 dark:text-gray-300 text-left pt-3 p-4">                 
                    <SubMeters v-if="data.cblist.length > 0" :data="data.cblist[index]" :temp="data.temp" />
                  </div>
                </div>
        </div>
      </div>
    </transition>
  </template>
  
  <script>
  import { ref, onMounted, onUnmounted, watch } from 'vue'
  //import SubDash from '../partials/job/SubDash.vue';
  import SubMeters from './SubMeters.vue'
  
  export default {
    name: 'ModalBasic',
    props: {
        modalOpen:Boolean,
        title:String,
        data:Object,
        id:String,
        notiType:String,
    },
    components:{
      SubMeters,
    },
    emits: ['close-modal'],
    setup(props, { emit }) {
      const modalContent = ref(null)
        const activeTab = ref('');
        const tabs = ref([]);
        const data = ref(props.data);
        const tabTitle = ref(props.title);
        const subTitle = ref('');
        const updated = ref('');
        const notiType = ref(props.notiType);
        const updateNotiType = (status) => {
        const devSt = status & 0xfffffc0e;
        const comSt = status & 0x01;
        let notiType = ''
        if (comSt == 0)
          notiType = 'Comm Error';
        else{
          if(devSt == 0)
            notiType = 'Online';
          else
            notiType = 'Alarm';
        }
        return notiType;
      };

      function convertTimestampToKST(timestamp) {
        // timestamp는 초 단위이므로 밀리초로 변환
        const date = new Date(timestamp * 1000); // 주의: 이미 소수점 포함된 초 단위

        // 한국 시간으로 변환
        const options = {
          timeZone: 'Asia/Seoul',
          year: 'numeric',
          month: '2-digit',
          day: '2-digit',
          hour: '2-digit',
          minute: '2-digit',
          second: '2-digit',
          hour12: false,
        };

        const formatter = new Intl.DateTimeFormat('en-CA', options);
        const parts = formatter.formatToParts(date);

        let result = {};
        parts.forEach(({ type, value }) => {
          result[type] = value;
        });

        // 밀리초 추출
        const milliseconds = Math.floor(date.getMilliseconds());

        return `${result.year}-${result.month}-${result.day} ${result.hour}:${result.minute}:${result.second}.${milliseconds.toString().padStart(3, '0')}`;
      }



        watch(
            () => props.data,
            (newVal) => {
                data.value = newVal;
                if(data.value.cblist.length >0){
                  tabs.value = [];
                  for (let i = 0; i < data.value.cblist.length; i++) {
                    tabs.value.push({
                        name: i.toString(),
                        label: 'CB #' + (i + 1),
                    });
                  }
                  updateNotiType(data.value.status);
                  subTitle.value = updateNotiType(data.value.status);
                  activeTab.value = tabs.value[0].name;
                  updated.value = convertTimestampToKST(data.value.updateTime);
                }
            },
            { immediate: true, deep: true }
            );
  
        // const setupTabs = () => {
           
        //     tabs.value = [];
        //     // if (data.value?.cblist?.length > 0) {
        //     //     for (let i = 0; i < data.value.cblist.length; i++) {
        //     //     tabs.value.push({
        //     //         name: i.toString(),
        //     //         label: 'CB #' + (i + 1),
        //     //     });
        //     //     }
        //     //     activeTab.value = tabs.value[0]?.name || '';
        //     // }
        //     };

      const changeTab = (tabName) => {
        activeTab.value = tabName;
        //tabTitle.value = tabName;
        };
  
      const clickHandler = (event) => {
        if (!props.modalOpen) return;
  
        // 클릭한 대상이 모달 내부인지 확인
        if (modalContent.value && modalContent.value.contains(event.target)) {
          event.stopPropagation(); // 이벤트 전파 방지
          return;
        }
  
        emit('close-modal');
      };
  
  
      // close if the esc key is pressed
      const keyHandler = ({ keyCode }) => {
        if (!props.modalOpen || keyCode !== 27) return
        emit('close-modal')
      }
  
      onMounted(() => {
        //document.addEventListener('click', clickHandler)
        document.addEventListener('keydown', keyHandler)
        // if(data.value.cblist.length > 0){
        //     for(let i = 0 ; i < data.value.cblist.length;i++){
        //         tabs.value.push({"name" : i.toString(), "label":"CB #"+(i+1).toString()});
        //     }
        // }
      })
  
      onUnmounted(() => {
        //document.removeEventListener('click', clickHandler)
        document.removeEventListener('keydown', keyHandler)
      })
  
      return {
        modalContent,
        changeTab,
        activeTab,
        data,
        tabs,
        tabTitle,
        clickHandler,
        subTitle,
        updated,
      }
    }  
  }
  </script>
  
  <style scoped>
  table {
    table-layout: fixed;
    width: 100%;
    min-height: calc(40px * 20); /* 40px * 20줄 높이 고정 */
    border-collapse: collapse;
  }
  
  th,
  td {
    text-align: center;
    padding: 8px;
    border: 1px solid #ddd;
    height: 40px; /* 각 행 높이 고정 */
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  th:nth-child(1) {
    width: 300px;
  } /* 첫 번째 컬럼 */
  th:nth-child(2) {
    width: 150px;
  }
  th:nth-child(3) {
    width: 150px;
  }
  th:nth-child(4) {
    width: 200px;
  }
  th:nth-child(5) {
    width: 100px;
  }
  .empty-row {
    visibility: hidden; /* 빈 행 안보이게 처리 */
  }
  
  
  </style>
  