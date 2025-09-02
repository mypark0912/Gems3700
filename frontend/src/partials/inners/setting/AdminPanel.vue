<template>
  <div class="w-full">
    <!-- 탭 메뉴를 맨 위로 이동 -->
    <div class="mb-6">
      <SettingsSidebar :api="APIUse" />
    </div>
    
    <!-- Content Panels -->
    <div class="bg-white dark:bg-gray-800 shadow-sm rounded-xl">
      <AccountPanel v-if="channel =='Profile'" />
      <BillingPanel v-else-if="channel == 'User'" />
      <APIPanel v-else-if="channel == 'APIUser' && APIUse" />
    </div>
  </div>
</template>

<script>
import { computed, watch, ref } from "vue";
import SettingsSidebar from "./SettingsSidebar.vue";
import { useSetupStore } from "@/store/setup"; // ✅ Pinia Store 사용
import BillingPanel from "./BillingPanel.vue";
import AccountPanel from "./AccountPanel.vue";
import APIPanel from "./APIPanel.vue";
import { useRoute } from "vue-router";

export default {
  name: "Setting",
  props:["channel"],
  components: {
    AccountPanel,
    SettingsSidebar,
    BillingPanel,
    APIPanel,
  },
  setup(props) {
    const route = useRoute();
    const channel = computed(() => props.channel || route.params.channel);

    const setupStore = useSetupStore();
    const APIUse = computed(()=>{
      const result = setupStore.getChannelSetting;
      return (result.MainDiagnosis || result.SubDiagnosis);
    })

  watch(
    () => route.params.channel,
    (newChannel) => {
      channel.value = newChannel;
      //console.log("Updated Channel:", channel.value);
    }
  );

    return{
      channel,
      APIUse,
      //isReady,
    }
  },
};
</script>