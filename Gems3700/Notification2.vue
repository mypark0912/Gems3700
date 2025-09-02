<template>
    <div v-show="open" role="alert">
      <div class="inline-flex flex-col w-full max-w-lg px-4 py-2 rounded-lg text-sm bg-white dark:bg-gray-800 shadow-sm border border-gray-200 dark:border-gray-700/60 text-gray-600 dark:text-gray-400">
        <div class="flex w-full justify-between items-start">
          <div class="flex cursor-pointer" @click="handleClick">
            <!-- ✅ `notiType.value` 대신 `notiType` 사용 (Vue가 감지 가능하도록) -->
            <svg v-if="notiType === 'alarm'" class="shrink-0 fill-current text-yellow-500 mt-[3px] mr-3" width="16" height="16" viewBox="0 0 16 16">
              <path d="M8 0C3.6 0 0 3.6 0 8s3.6 8 8 8 8-3.6 8-8-3.6-8-8-8zm0 12c-.6 0-1-.4-1-1s.4-1 1-1 1 .4 1 1-.4 1-1 1zm1-3H7V4h2v5z" />
            </svg>
            <svg v-else-if="notiType === 'error'" class="shrink-0 fill-current text-red-500 mt-[3px] mr-3" width="16" height="16">
              <path d="M8 0C3.6 0 0 3.6 0 8s3.6 8 8 8 8-3.6 8-8-3.6-8-8-8zm3.5 10.1l-1.4 1.4L8 9.4l-2.1 2.1-1.4-1.4L6.6 8 4.5 5.9l1.4-1.4L8 6.6l2.1-2.1 1.4 1.4L9.4 8l2.1 2.1z" />
            </svg>
            <svg v-else-if="notiType === 'success'" class="shrink-0 fill-current text-green-500 mt-[3px] mr-3" width="16" height="16">
              <path d="M8 0C3.6 0 0 3.6 0 8s3.6 8 8 8 8-3.6 8-8-3.6-8-8-8zM7 11.4L3.6 8 5 6.6l2 2 4-4L12.4 6 7 11.4z" />
            </svg>
            <div>
              <div class="font-medium text-gray-800 dark:text-gray-100 mb-1">{{ chNum }}</div>
              <div class="flex items-center mt-2">
                <div class="flex flex-col">
                   <div class="flex items-center" >
                    <div v-if="devType.p == 4" class="text-xs font-medium text-sky-700 px-1.5 bg-sky-500/20 rounded-full">Total</div>
                    <div v-else class="text-xs font-medium text-sky-700 px-1.5 bg-sky-500/20 rounded-full">L{{ devType.p }}</div>
                    <div class="text-sm font-bold text-gray-800 dark:text-gray-100 mr-2">{{ printData.cur }} A</div>
                   </div>
                   <div class="flex items-center" >
                    <!--div class="text-sm font-bold text-gray-800 dark:text-gray-100 mr-2">{{ printData.power }} kW</div-->
                    <div class="text-sm font-bold text-gray-800 dark:text-gray-100 mr-2">{{ printData.kwh }} kWh</div>
                   </div>
                  
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <ModalBasic v-if="feedbackModalOpen"
                id="feedback-modal"
                :modalOpen="feedbackModalOpen"
                @close-modal="feedbackModalOpen = false"
                :data="subData" :numid="id" 
                :title="chNum" :notiType="notiType"
              />
  </template>
  
  <script>
  import { ref, watch, computed, onMounted } from 'vue';
  import ModalBasic from "./ModalBasic_modify.vue";
  //import SubDash from '../partials/job/SubDash.vue'
  export default {
    name: 'Notification',
    components:{
      ModalBasic,
    },
    props: {
      open: Boolean,
      data:Object,
      id:String,
      channel:String,
    },
    setup(props) {
      const subData = ref(props.data);
      const id = ref(props.id);
      const channel = ref(props.channel);
      const chNum = ref('');
      const printData = ref({});
      const feedbackModalOpen = ref(false);
      const devType = ref({});
    //   const stData = ref(props.data);
  
      const notiType = ref('none');
  
      const updateNotiType = (status) => {
        const devSt = status & 0xfffffc0e;
        const comSt = status & 0x01;
        if (comSt == 0)
          notiType.value = 'error';
        else{
          if(devSt == 0)
            notiType.value = 'success';
          else
            notiType.value = 'alarm';
        }
      };

      const defType = (cbtype) =>{
        let re='L1';
        switch(cbtype){
          case 1:
            re = "L1";
            break;
          case 2:
            re = "L2";
            break;
          case 3:
            re = "L3";
            break;
          case 4:
            re = "3P3W";
            break;
          case 5:
            re = "3P4W";
            break;
          case 6:
            re = "L1+Z";
            break;
          case 7:
            re = "L2+Z";
            break;
          case 8:
            re = "L3+Z";
            break;
        }
        return re;
        
      }
  
      watch(
        () => props.data,
        (newData) => {
          if (newData && Object.keys(newData).length > 0) {
           // console.log("🔄 새 데이터 업데이트:", newData);
            subData.value = { ...newData }; 
            updateNotiType(subData.value.status);
            let irms = 0.0;
            const idx = subData.value.cblist[0].cbtype%3;
            if(subData.value.cblist[0].cbtype > 5){
              devType.value = { p: idx + 1, g:1}
              irms = parseFloat(subData.value.cblist[0].irms[0])
            }else if(subData.value.cblist[0].cbtype == 4 || subData.value.cblist[0].cbtype == 5){
              devType.value = { p: 4, g:1}
              irms = parseFloat(subData.value.cblist[0].irms[0]) + parseFloat(subData.value.cblist[0].irms[1]) + parseFloat(subData.value.cblist[0].irms[2]);
            }else{
              devType.value = { p: subData.value.cblist[0].cbtype, g:0}
              irms = parseFloat(subData.value.cblist[0].irms[0])
            }

            printData.value = {
              "cur" : irms.toFixed(2),
              "power" : parseFloat(subData.value.cblist[0].power[0]).toFixed(2),
              "kwh" : parseFloat(subData.value.cblist[0].kwh).toFixed(2),
              "type" : defType(subData.value.cblist[0].cbtype),
            }
            chNum.value = channel.value +"-" + id.value +" : " + defType(subData.value.cblist[0].cbtype);              
            //console.log("✅ 업데이트된 notiType:", notiType.value);
          }
        },
        { immediate: true, deep: true }
      );

    const handleClick = () => {
      feedbackModalOpen.value = true;
    };
  
      return {
        subData,
        notiType,
        chNum,
        id,
        channel,
        printData,
        feedbackModalOpen,
        handleClick,
        devType,
      };
    },
  };
  </script>
  