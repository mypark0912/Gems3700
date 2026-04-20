<template>
  <div class="flex h-[100dvh] overflow-hidden">

    <!-- Sidebar -->
    <Sidebar :sidebarOpen="sidebarOpen" @close-sidebar="sidebarOpen = true" />

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
              <h2 class="text-xl md:text-2xl text-gray-800 dark:text-gray-100 font-bold ml-1">{{ t('meter.sitemap.title') }}</h2>
            </div>

          </div>

          <!-- Cards -->
          <div class="grid grid-cols-12 gap-4">

            <MeterDetail2 v-if="Object.keys(meterDatas).length > 0" :data="meterDatas.meterData" :channel="channel" :title="'Meter'" />
            <div class="col-span-5 md:col-span-5 flex flex-col gap-4">
              <div class="col-span-5">
                <CanvasAngle2 v-if="Object.keys(phaseDict).length > 0"
                  :degree="phaseDict.degree"
                  :magnitude="phaseDict.magnitude"
                  :texts="phaseDict.texts"
                  :maxlist="phaseDict.max"
                  :channel="channel"
                />
              </div>
              <div class="flex gap-4">
              <div class="basis-1/2">
                <MeterKwh
                  v-if="Object.keys(meterDatas).length > 0"
                  :data="energyData"
                  :channel="channel"
                  :title="'Energy'"
                  :mode="'import'"
                />
              </div>
              <div class="basis-1/2">
                <MeterKwh
                  v-if="Object.keys(energyData).length > 0"
                  :data="energyData"
                  :channel="channel"
                  :title="'Energy'"
                  :mode="'export'"
                />
              </div>
            </div>
            </div>

            <MeterDetail2 v-if="Object.keys(powerThd).length > 0 && powerThd.thdData" :data="powerThd.thdData" :channel="channel" :title="'THD'" />
            <MeterDetail2 v-if="Object.keys(powerThd).length > 0 && powerThd.demandDataP" :data="powerThd.demandDataP" :channel="channel" :title="'Demand'" />

            <MeterDetail2 v-if="Object.keys(powerThd).length > 0 && powerThd.powerData" :data="powerThd.powerData" :channel="channel" :title="'Power'" />
            <MeterDetail2 v-if="Object.keys(powerThd).length > 0 && powerThd.demandDataI" :data="powerThd.demandDataI" :channel="channel" :title="'Demand I'" />

          </div>

        </div>
      </main>
      <Footer />
    </div>

  </div>
</template>

<script>
import { reactive,ref, watch, computed, onMounted, onUnmounted } from 'vue'

import axios from 'axios'
import Sidebar from '../common/SideBar.vue'
import Header from '../common/Header.vue'
import Footer from "../common/Footer.vue";

import CanvasAngle2 from '../../partials/inners/design/CanvasAngle.vue'

import MeterKwh from '../../partials/inners/design/MeterKwh.vue'
import MeterDetail2 from '../../partials/inners/design/MeterDetail.vue'
import { useSetupStore } from '@/store/setup'
import { useI18n } from 'vue-i18n'

export default {
  name: 'MeterDesign',
  props: {
    channel: {
      type: String,
      default: 'Main',
    },
  },
  components: {
    Sidebar,
    Header,
    Footer,
    CanvasAngle2,
    MeterKwh,
    MeterDetail2,
  },
  setup(props) {
    const { t } = useI18n();
    const sidebarOpen = ref(true)
    const channel = ref(props.channel)
    const meterDatas = ref({});
    const powerThd = ref({});
    const phaseDict = ref({});
    const energyData = ref({});

    const setupStore = useSetupStore();
    let updateInterval_ones = null;
    let updateInterval_onem = null;
    let updateInterval_fifthm = null;
    let updateInterval_oneh = null;
    let timeout_fifthm = 15*60*1000;

    const setup = computed(() => setupStore.getSetup);

    const unbal = computed(()=> setupStore.getUnbalance);

    const ReleaseInterval = () =>{
      if (updateInterval_ones) {
        clearInterval(updateInterval_ones);
        updateInterval_ones = null
      }
      if (updateInterval_onem) {
        clearInterval(updateInterval_onem);
        updateInterval_onem = null
      }
      if (updateInterval_fifthm) {
        clearInterval(updateInterval_fifthm);
        updateInterval_fifthm = null
      }
      if (updateInterval_oneh) {
        clearInterval(updateInterval_oneh);
        updateInterval_oneh = null
      }
    };

    const fetchInterval = async (ch) => {
      try {
        const response = await axios.get(`/api/getInteverval/demand/${ch}`);
        if (response.data.success) {
          timeout_fifthm = parseInt(response.data.data);
        }
      } catch (error) {
        console.log("데이터 가져오기 실패:", error);
      }
    };

    onUnmounted(() => {
      ReleaseInterval();
    });

    watch(() => props.channel, (newChannel) => {
      if (newChannel !== channel.value) {
        console.log("Channel changed to:", newChannel);
        channel.value = newChannel;
        meterDatas.value = {};
        powerThd.value = {};
        phaseDict.value = {};
        energyData.value ={};
        demandData.value = {
          demandData: [
            { subTitle: "Active", data: [] },
            { subTitle: "Reactive", data: [] },
            { subTitle: "Apparent", data: [] }
          ],
          demandIData: [
            { subTitle: "L1", data: [] },
            { subTitle: "L2", data: [] },
            { subTitle: "L3", data: [] }
          ]
        };
        ReleaseInterval();
        startFetching();
      }
    },{ immediate: true });

    const fetchRedisOnesData = async (ch) => {
      try {
        const response = await axios.get(`/api/getOnesfromRedis/${ch}/${unbal.value}`);
        if (response.data.success) {
          Object.assign(meterDatas.value, response.data.retData.retData);
        }
      } catch (error) {
        console.log("데이터 가져오기 실패:", error);
      }
    };

    const fetchRedisOnemData = async (ch) => {
      try {
        const response = await axios.get(`/api/getonemfromRedis/${ch}`);
        if (response.data.success) {
          phaseDict.value = response.data.retData.retData.angleData;
        }
      } catch (error) {
        console.log("데이터 가져오기 실패:", error);
      }
    };

    const fetchRedisfifthmData = async (ch) => {
      try {
        const response = await axios.get(`/api/getFifthMfromRedis/${ch}`);
        if (response.data.success) {
          powerThd.value = response.data.retData.retData;
          powerThd.value.powerData.forEach(section => {
            section.data.forEach(item => {
              item.value = (item.value / 1000).toFixed(2);
            });
          });
        }
      } catch (error) {
        console.log("데이터 가져오기 실패:", error);
      }
    };

    const fetchRedisOnehData = async (ch) => {
      try {
        const response = await axios.get(`/api/getOnehfromRedis/${ch}`);

        if (response.data.success) {
          energyData.value = response.data.retData.retData.energyData;
        }else{
          if ('retData' in response.data) {
            energyData.value = response.data.retData.retData.energyData;
          }
        }
      } catch (error) {
        console.log("데이터 가져오기 실패:", error);
      }
    };

    const startFetching = async() => {
      await fetchInterval(channel.value);

      if (updateInterval_ones) {
        clearInterval(updateInterval_ones);
        updateInterval_ones = null
      }
      fetchRedisOnesData(channel.value);
        updateInterval_ones = setInterval(() => {
          fetchRedisOnesData(channel.value);
        }, 1000);

      if (updateInterval_onem){
        clearInterval(updateInterval_onem);
        updateInterval_onem = null
      }
      fetchRedisOnemData(channel.value);
        updateInterval_onem = setInterval(() => {
          fetchRedisOnemData(channel.value);
        }, 60*1000);

      if (updateInterval_fifthm){
        clearInterval(updateInterval_fifthm);
        updateInterval_fifthm = null
      }
      fetchRedisfifthmData(channel.value);
        updateInterval_fifthm = setInterval(() => {
          fetchRedisfifthmData(channel.value);
        }, timeout_fifthm*1000);

      if (updateInterval_oneh) {
        clearInterval(updateInterval_oneh);
        updateInterval_oneh = null
      }
      fetchRedisOnehData(channel.value);
        updateInterval_oneh = setInterval(() => {
          fetchRedisOnehData(channel.value);
        }, 60*60*1000);
    };

    watch(() => setup.value, (newChannel) => {
      if (newChannel) {
        startFetching();
      }
    },{ immediate: true });

    return {
      sidebarOpen,
      channel,
      meterDatas,
      phaseDict,
      t,
      powerThd,
      energyData,
      setup,
      unbal,
    }
  }
}
</script>
