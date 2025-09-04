<template>
    <div class="flex h-[100dvh] overflow-hidden">
  
      <!-- Sidebar -->
      <Sidebar :sidebarOpen="sidebarOpen" @close-sidebar="sidebarOpen = true" :channel="channel"/>
  
      <!-- Content area -->
      <div class="relative flex flex-col flex-1 overflow-y-auto overflow-x-hidden">
        
        <!-- Site header -->
        <Header :sidebarOpen="sidebarOpen" @toggle-sidebar="sidebarOpen = !sidebarOpen" />
  
        <main class="grow">
          <div class="px-2 sm:px-4 lg:px-6 py-4 w-full max-w-full">
            
            <!-- Dashboard actions -->
            <div class="sm:flex sm:justify-between sm:items-center mb-4">
  
              <!-- Left: Title -->
              <div class="mb-4 sm:mb-0">
                <h2 class="text-xl md:text-2xl text-gray-800 dark:text-gray-100 font-bold"> {{ t('report.sitemap.title') }} > {{ channelComputed == 'Main'?t('report.sitemap.main'):t('report.sitemap.sub') }} </h2>
              </div>
              <button
                  class="btn h-6 px-5 bg-sky-900 text-sky-100 hover:bg-sky-800 dark:bg-sky-100 dark:text-sky-800 dark:hover:bg-white"
                  @click="openPdfPreview"
                >
                  PDF Download
                </button>
  
  
            </div>
<ModalBasic
    id="feedback-modal"
    :modalOpen="feedbackModalOpen"
    @close-modal="feedbackModalOpen = false"
    title="Preview for pdf"
  >
    <!-- Modal content -->
    <PdfView ref="pdfComponent" :channel="channel" :data = "tbdata" :pdf-data="pdfData" @render-complete="onPdfRenderComplete"/>
    <!-- Modal footer -->
    <div
      class="px-5 py-4 border-t border-gray-200 dark:border-gray-700/60"
    >
      <div class="flex flex-wrap justify-end space-x-2">
        <!-- 로딩 중일 때 표시할 메시지 -->
        <div v-if="!isPdfRendered" class="flex items-center text-gray-600 dark:text-gray-400">
          <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-gray-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          Loading PDF preview...
        </div>
        <button
          v-show="isPdfRendered"
          class="btn-sm bg-gray-900 text-gray-100 hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white"
          @click.prevent="savePDF"
        >
          Download
        </button>
        <button
          class="btn-sm border-gray-200 dark:border-gray-700/60 hover:border-gray-300 dark:hover:border-gray-600 text-gray-800 dark:text-white"
          @click.stop="feedbackModalOpen = false"
        >
          Close
        </button>
      </div>
    </div>
  </ModalBasic>
            <!-- Cards -->
            
            <div class="grid grid-cols-12 gap-6">
              <!--Report_Info v-if="mode" :channel="channelComputed" :mode="mode" :key="`info-${channelComputed}`" /-->
                          <div class="md:col-span-12 bg-white dark:bg-gray-800 shadow-md rounded-lg p-4 w-full">

              <!-- Tab Navigation -->
              <div class="px-4">
                <ul class="text-sm font-medium flex flex-nowrap overflow-x-auto no-scrollbar w-full">
                  <li v-for="(tab, index) in tabs" :key="index" class="mr-4 last:mr-0 relative">
                    <button
                      @click="changeTab(tab.name)"
                      class="px-4 py-2 whitespace-nowrap transition duration-200 ease-in-out"
                      :class="activeTab === tab.name
                        ? 'text-violet-500 border-b-2 border-violet-500'
                        : 'text-gray-500 dark:text-gray-400 hover:text-gray-600 dark:hover:text-gray-300'">
                      {{ t(`report.cardTitle.${tab.label}`)}}
                    </button>
                  </li>
                </ul>
              </div>

              <!-- Tab Content -->
              <div v-for="(tab, index) in tabs" :key="index">
                <div v-if="activeTab === tab.name" class="text-gray-700 dark:text-white text-left pt-3 px-4">
                
                  <!-- 차트 컨테이너 -->
                  <div class="flex flex-col space-y-2">
                      <Report_Diagnosis v-if="activeTab === 'Equipment' && mode" :channel="channelComputed" :key="`diag-${channelComputed}`" />
                      <ReportComponent v-if="activeTab === 'PowerQuality' && tbdata.length > 0" :data="tbdata" :channel="channelComputed" :mode="mode" :key="`component-${channelComputed}`" />
                      <Report_WattHour v-if="activeTab === 'Energy'" :mode="mode" :channel="channelComputed" :key="`wh-${channelComputed}`" />
                  </div>

                </div>
              </div>

            </div>    

            </div>
  
          </div>
        </main>
        <Footer />
      </div> 
  
    </div>
  </template>
  
  <script>
  import { ref, watch, computed, onMounted,nextTick } from 'vue'
  import { useRoute } from 'vue-router'
  import { useSetupStore } from '@/store/setup'; // Pinia Store 
  import axios from 'axios'
  import Sidebar from '../common/SideBarGems.vue'
  import Header from '../common/Header.vue'
  import Footer from "../common/Footer.vue";
  import Report_table from '../../partials/inners/report/Report_table.vue';
  import ReportComponent from '../../partials/inners/report/ReportComponent.vue';
  import Report_WattHour from '../../partials/inners/report/Report_WattHour.vue';
  //import Report_Energy from '../../partials/Report_Energy.vue';
  import Report_ITIC from '../../partials/inners/report/Report_ITIC.vue';
  import Report_Diagnosis from '../../partials/inners/report/Report_Diagnosis.vue';
  //import Report_Diagnosis from '../../partials/dashboard/Report_Diagnosis_trans.vue';
  //import Report_Info from "../../partials/Report_Info.vue";
  import Report_Info from "../../partials/inners/report/Report_Info.vue";
  import PdfView from "../../partials/inners/report/PDFPreview.vue";
  import ModalBasic from "../common/ModalBasic.vue";
  import { useI18n } from 'vue-i18n'  

  export default {
    name: 'Report',
    props:['channel'],
    components: {
      Sidebar,
      Header,
      Footer,
      Report_table,
      Report_ITIC,
      ReportComponent,
      Report_WattHour,
      Report_Diagnosis,
      Report_Info,
      ModalBasic,
      PdfView,
    },
    setup(props) {
      const { t } = useI18n();
      const route = useRoute()
      const sidebarOpen = ref(true)
      const channel = ref(props.channel)
      const setupStore = useSetupStore();
      const feedbackModalOpen = ref(false);
      const isPdfRendered = ref(false);
      const channelComputed = computed(() => props.channel || route.params.channel || 'Default')
      // const tabs = ref([
      //   {name:"Equipment", label:"Diagnosis"}, {name:"PowerQuality", label:"PowerQuality"}, {name:"Energy", label:"Energy"}
      // ]);
      const tbdata = ref([]);
      const activeTab = ref('Equipment');
  
      const channelStatus = computed(() => setupStore.getChannelSetting);
      const setupMenu = ref({});
      const mode = computed(() => {
        //console.log(channelComputed.value + ' : ' + setupMenu.value.MainDiagnosis);
        if (channelComputed.value === 'Main')
            return channelComputed.value === 'Main' && setupMenu.value.MainDiagnosis;
        else
            return channelComputed.value === 'Sub' && setupMenu.value.SubDiagnosis;
    });

    const tabs = computed(()=>{
          activeTab.value = 'PowerQuality';
          return [
            {name:"PowerQuality", label:"PowerQuality"}, {name:"Energy", label:"Energy"}
          ]
        // if (mode.value){
        //   activeTab.value = 'Equipment';
        //   return [
        //     {name:"Equipment", label:"Diagnosis"}, {name:"PowerQuality", label:"PowerQuality"}, {name:"Energy", label:"Energy"}
        //   ]
        // }else{
        //   activeTab.value = 'PowerQuality';
        //   return [
        //     {name:"PowerQuality", label:"PowerQuality"}, {name:"Energy", label:"Energy"}
        //   ]
        // }
    }
  );
  
    watch(() => route.params.channel, async (newChannel) => {
        channel.value = newChannel
        //console.log('Updated Channel:', channel.value)
        await fetchData();
      })

      watch(channelStatus, (newVal) => {
        setupMenu.value = newVal; // ✅ setupMenu 초기화
        }, { immediate: true });

      onMounted(async () => {
          await fetchData();  // 데이터 가져오기 후 렌더링
      });

      const changeTab = (tabName) => {
      activeTab.value = tabName;
      nextTick(() => {
      });
    };

      const fetchData = async () => {
      try {
        //console.log("Data Get");
        // const response = await axios.get('/api/getFormat');
        // if (response.data.success) {
        //     tbdata.value = [...response.data.data.tbdata]; 
        // } else {
        //     console.warn("서버 응답이 success: false 입니다.");
        //     tbdata.value = []; // 기본값 설정
        // }
          const res = await fetch('/en50160_info.json')
          const data = await res.json()
          tbdata.value = [...data.tbdata]
      } catch (error) {
        console.log("데이터 가져오기 실패:", error);
      }
    };

const onPdfRenderComplete = () => {
  isPdfRendered.value = true; // PDF 렌더링 완료
}
const openPdfPreview = () => {
  // 각 하위 컴포넌트에서 데이터 수집
  // pdfData.value = {
  //   channel: 'Main',
  //   datalist: infoRef.value?.getData(),
  // }
  isPdfRendered.value = false; // PDF 렌더링 상태 초기화
  feedbackModalOpen.value = true
}
  
      return {
        tabs,
        sidebarOpen,
        channel,
        //langset,
        tbdata,
        channelComputed,
        mode,
        feedbackModalOpen,
        isPdfRendered,
        t,
        activeTab,
        changeTab,
        openPdfPreview,
        onPdfRenderComplete,
      }  
    },
    methods: {
      savePDF() {
        this.$refs.pdfComponent.generatePDF();
      }
    }
  }
  </script>