<template>
  <router-view v-if="ready" />
</template>

<script setup>
import { onMounted, computed, watch, onUnmounted, ref } from "vue";
import { useAuthStore } from "@/store/auth"; // ✅ Pinia store 가져오기
import { useSetupStore } from "@/store/setup"; // ✅ Pinia store 가져오기
import "@/charts/ChartjsConfig";
import router from '@/router2';
import i18n from './i18n'  // 너의 i18n 인스턴스

const authStore = useAuthStore(); // ✅ Pinia store 사용
const setupStore = useSetupStore(); // ✅ Pinia store 사용
const ready = ref(false)
const isAuthenticated = computed(() => authStore.getLogin);
let sessionCheckInterval = null;

onMounted(() => {
  sessionCheckInterval = setInterval(() => {
    if (!authStore.sessionChecked) {  // ✅ 중복 실행 방지
      authStore.checkSession();  // ✅ Pinia action 호출
      authStore.setSessionChecked(true); // ✅ 실행 완료 표시
      //console.log("✅ checkSession 실행됨!");
    } else {
      //console.log("⚠️ checkSession 이미 실행됨 - 중복 실행 방지");
    }
  }, 300000); // 5분 (300,000ms)
});

onMounted(async () => {
  if (!setupStore.applysetup || !setupStore.opMode) {
    await setupStore.checkSetting();
  }
  ready.value = true;
  // authStore.fetchLangset(authStore.language);
});

onUnmounted(() => {
  clearInterval(sessionCheckInterval);
});

watch(
  () => i18n.global.locale.value,
  (newLang) => {
    document.documentElement.setAttribute('lang', newLang)
  },
  { immediate: true }
)

watch(isAuthenticated, async (newVal) => {
  if (newVal) {
    //console.log("🚀 로그인 후 Sidebar 즉시 업데이트!");
    await setupStore.checkSetting(true); // ✅ 강제 업데이트
  }
});

watch(() => router.path, async () => {
  if (isAuthenticated.value && !setupStore.applysetup) {
    await setupStore.checkSetting(true); // ✅ 강제 업데이트 적용
  }
  ready.value = true;
});

</script>