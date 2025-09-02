<template>
  <div class="relative inline-flex">
    <button
      ref="trigger"
      class="inline-flex justify-center items-center group"
      aria-haspopup="true"
      @click.prevent="dropdownOpen = !dropdownOpen"
      :aria-expanded="dropdownOpen"
    >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      width="32"
      height="32"
      fill="currentColor"
    >
      <circle cx="12" cy="12" r="10" fill="#3498db" />
      <circle cx="12" cy="9" r="4" fill="white" />
      <path d="M6 19c0-3 2.5-5 6-5s6 2 6 5" fill="white" />
    </svg>
      <div class="flex items-center truncate">
        <span class="truncate ml-2 text-sm font-medium text-gray-600 dark:text-gray-100 group-hover:text-gray-800 dark:group-hover:text-white">{{user}}</span>
        <svg class="w-3 h-3 shrink-0 ml-1 fill-current text-gray-400 dark:text-gray-500" viewBox="0 0 12 12">
          <path d="M5.9 11.4L.5 6l1.4-1.4 4 4 4-4L11.3 6z" />
        </svg>
      </div>
    </button>
    <transition
      enter-active-class="transition ease-out duration-200 transform"
      enter-from-class="opacity-0 -translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition ease-out duration-200"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div v-show="dropdownOpen" class="origin-top-right z-10 absolute top-full min-w-44 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 py-1.5 rounded-lg shadow-lg overflow-hidden mt-1" :class="align === 'right' ? 'right-0' : 'left-0'">
        <div class="pt-0.5 pb-2 px-3 mb-1 border-b border-gray-200 dark:border-gray-700/60">
          <div class="font-medium text-gray-800 dark:text-gray-100">{{user}}</div>
        </div>
        <ul
          ref="dropdown"
          @focusin="dropdownOpen = true"
          @focusout="dropdownOpen = false"
        >
          <li v-if="!isAdmin">
            <router-link class="font-medium text-sm text-violet-500 hover:text-violet-600 dark:hover:text-violet-400 flex items-center py-1 px-3" to="/settings/Account/Profile" @click="dropdownOpen = false">
              <span>{{ t('header.profile') }}</span>
            </router-link>
          </li>
          <li v-else>
            <router-link class="font-medium text-sm text-violet-500 hover:text-violet-600 dark:hover:text-violet-400 flex items-center py-1 px-3" to="/settings/Account/User" @click="dropdownOpen = false">
              <span>{{ t('header.user') }}</span>
            </router-link>
          </li>
          <li v-if="isAdmin">
            <router-link class="font-medium text-sm text-violet-500 hover:text-violet-600 dark:hover:text-violet-400 flex items-center py-1 px-3" to="/config/Maintenance" @click="dropdownOpen = false">
              <span>{{ t('header.maintenance') }}</span>
            </router-link>
          </li>
          <li>
            <a
              href="#"
              class="font-medium text-sm text-violet-500 hover:text-violet-600 dark:hover:text-violet-400 flex items-center py-1 px-3"
              @click.prevent="handleLogout"
            >
            {{ t('header.signout') }}
            </a>
          </li>
        </ul>
      </div> 
    </transition>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, watch, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/store/auth'; // ✅ Pinia Store 사용
import { useI18n } from 'vue-i18n'  // ✅ 추가

export default {
  name: 'DropdownProfile',
  props: ['align'], // ✅ props 명칭 변경 ('lang' → 'langset' 으로 통일)
  setup(props) {
    const { t } = useI18n();      
    const authStore = useAuthStore(); // ✅ Pinia store 사용
    const router = useRouter();
    const user = computed(() => authStore?.getUser || ''); // ✅ undefined 
    const isAdmin = computed(() =>{
      const role = parseInt(authStore?.getUserRole);
      if (role > 1)
        return true;
      else
        return false;
    });
    // ✅ 로그인한 사용자 정보 가져오기

    // ✅ 드롭다운 상태 관리
    const dropdownOpen = ref(false);
    const trigger = ref(null);
    const dropdown = ref(null);

    // ✅ 언어 데이터 반응형으로 유지
    // const langs = ref(props.langset || {}); 

    // // ✅ props.langset이 변경될 때 자동 업데이트
    // watch(() => props.langset, (newLangset) => {
    //   langs.value = newLangset || {};
    // });

    // ✅ 로그아웃 함수 수정
    const handleLogout = async () => {
      try {
        if (authStore.logout) { // ✅ logout 함수가 정의되어 있는지 확인
          await authStore.logout();
          dropdownOpen.value = false; // ✅ 드롭다운 닫기
          router.push("/signin"); // ✅ 로그인 페이지로 이동
        } else {
          console.error("❌ Error: logout 함수가 정의되지 않음!");
        }
      } catch (err) {
        console.error("Logout Error:", err);
      }
    };

    // ✅ 드롭다운 외부 클릭 감지
    const clickHandler = ({ target }) => {
      if (
        dropdownOpen.value &&
        dropdown.value &&
        !dropdown.value.contains(target) &&
        trigger.value &&
        !trigger.value.contains(target)
      ) {
        dropdownOpen.value = false;
      }
    };

    // ✅ ESC 키를 눌렀을 때 드롭다운 닫기
    const keyHandler = ({ keyCode }) => {
      if (!dropdownOpen.value || keyCode !== 27) return;
      dropdownOpen.value = false;
    };

    // ✅ 컴포넌트 마운트 시 이벤트 추가
    onMounted(() => {
      if (!authStore.getLogin) { // ✅ 로그인 정보가 없으면 checkSession 실행
        authStore.checkSession();
      }
      document.addEventListener('click', clickHandler);
      document.addEventListener('keydown', keyHandler);
    });

    // ✅ 컴포넌트 언마운트 시 이벤트 제거
    onUnmounted(() => {
      document.removeEventListener('click', clickHandler);
      document.removeEventListener('keydown', keyHandler);
    });

    return {
      dropdownOpen,
      trigger,
      dropdown,
      user,
      handleLogout,
      //langs,
      t,
      isAdmin,
    };
  }
};
</script>

