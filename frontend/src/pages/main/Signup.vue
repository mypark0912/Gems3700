<template>
  <main class="bg-white dark:bg-gray-900">
    <div class="relative flex">
      <!-- Content -->
      <div class="w-full md:w-1/2">
        <div class="min-h-[100dvh] h-full flex flex-col after:flex-1">
          <div class="flex-1">
            <div
              class="flex items-center justify-between h-16 px-4 sm:px-6 lg:px-8"
            >
              <!-- Logo -->
              <router-link class="block" to="/">
                <svg
                  class="fill-violet-500"
                  xmlns="http://www.w3.org/2000/svg"
                  width="32"
                  height="32"
                >
                  <path
                    d="M31.956 14.8C31.372 6.92 25.08.628 17.2.044V5.76a9.04 9.04 0 0 0 9.04 9.04h5.716ZM14.8 26.24v5.716C6.92 31.372.63 25.08.044 17.2H5.76a9.04 9.04 0 0 1 9.04 9.04Zm11.44-9.04h5.716c-.584 7.88-6.876 14.172-14.756 14.756V26.24a9.04 9.04 0 0 1 9.04-9.04ZM.044 14.8C.63 6.92 6.92.628 14.8.044V5.76a9.04 9.04 0 0 1-9.04 9.04H.044Z"
                  />
                </svg>
              </router-link>
            </div>
          </div>
            <div class="max-w-sm mx-auto w-full px-4 py-4">
            <h1
              class="text-3xl text-gray-800 dark:text-gray-100 font-bold mb-6"
            >
             {{  t('signup.title1') }}
            </h1>
            <!-- Form -->
            <form>
              <div class="flex space-x-4">
                <div class="w-1/3">
                  <label class="block text-sm font-medium mb-1" for="language">LANGUAGE</label>
                  <select
                    id="language"
                    class="form-input w-full px-2 py-1 border rounded bg-white"
                    v-model="lang"
                    @change="onLangChange"
                  >
                    <option value="ko">KOREAN</option>
                    <option value="en">ENGLISH</option>
                    <option value="ja">JAPANESE</option>
                  </select>
                </div>
                <div class="w-2/3">
                  <label class="block text-sm font-medium mb-1" for="devtype">
                    {{  t('signup.deviceType') }} <span class="text-red-500">*</span>
                  </label>
                  <select id="devtype" class="form-input w-full px-2 py-1 border rounded bg-white" v-model="devType">
                    <option value="0">{{  t('signup.device0') }}</option>
                    <option value="1">{{  t('signup.device1') }}</option>
                    <option value="2">{{  t('signup.device2') }}</option>
                    <!--option value="3">{{  t('signup.device3') }}</option-->
                  </select>
                </div>
              </div>

            </form>
          </div>
          <div class="max-w-sm mx-auto w-full px-4 py-4">
            <h1
              class="text-3xl text-gray-800 dark:text-gray-100 font-bold mb-6"
            >
              {{  t('signup.title2') }}
            </h1>
            <!-- Form -->
            <form>
              <div class="space-y-4">
                <div>
                  <label class="block text-sm font-medium mb-1" for="name">
                    {{  t('signup.name') }} <span class="text-red-500">*</span>
                  </label>
                  <input id="name" class="form-input w-full" type="text"  v-model="username"/>
                </div>
                <div>
                  <label class="block text-sm font-medium mb-1" for="email">
                    {{  t('signup.email') }} <span class="text-red-500">*</span>
                  </label>
                  <input id="email" class="form-input w-full" type="email" v-model="email"/>
                </div>
                <div>
                  <label class="block text-sm font-medium mb-1" for="account">
                    {{  t('signup.account') }} <span class="text-red-500">*</span>
                  </label>
                  <input id="account" class="form-input w-full" type="text"  v-model="account" disabled/>
                </div>
                <div>
                  <label class="block text-sm font-medium mb-1" for="role">
                    {{  t('signup.accountType') }} <span class="text-red-500">*</span>
                  </label>
                  <select id="role" class="form-select w-full" v-model="role" disabled>
                    <option value="0">Operator</option>
                    <option value="1">User</option>
                    <option value="2">Admin</option>
                    <option value="3">Ntek Admin</option>
                  </select>
                </div>
                <!--div>
                  <label class="block text-sm font-medium mb-1" for="password">
                    {{  t('signup.password') }}
                  </label>
                  <input
                    id="password"
                    class="form-input w-full"
                    type="password"
                    autoComplete="on"
                    v-model="password"
                  />
                </div-->
                <div>
                  <label class="block text-sm font-medium mb-1" for="adminpassword">
                    {{  t('signup.Adminpassword') }}
                  </label>
                  <input
                    id="adminpassword"
                    class="form-input w-full"
                    type="password"
                    autoComplete="on"
                    v-model="adminpassword"
                  />
                </div>
                <!-- Form button wrapped with a div for alignment -->
                <div class="flex justify-end">
                  <button
                   @click.prevent="signup"
                    type="button"
                    class="btn bg-gray-900 text-gray-100 hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white whitespace-nowrap"
                  >
                    {{  t('signup.register') }}
                  </button>
                </div>
              </div>
            </form>

            <!-- Footer -->
            <div
              class="pt-5 mt-6 border-t border-gray-100 dark:border-gray-700/60"
            >
              <div class="text-sm">
                {{  t('signup.msg') }}
                <router-link
                  class="font-medium text-violet-500 hover:text-violet-600 dark:hover:text-violet-400"
                  to="/signin"
                  >{{  t('signup.link') }}</router-link
                >
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Image -->
      <div
        class="hidden md:block absolute top-0 bottom-0 right-0 md:w-1/2"
        aria-hidden="true"
      >
        <img
          class="object-cover object-center w-full h-full"
          src="../../images/background2.png"
          width="760"
          height="1024"
          alt="Authentication"
        />
      </div>
    </div>
  </main>
</template>

<script setup>


import { ref, onMounted } from "vue";
import axios from "axios";
import { useRouter } from "vue-router";
import { useAuthStore } from '@/store/auth'
import i18n from '@/i18n'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
const email = ref("");
const username = ref("");
const account = ref("ntek");
const password = ref("");
const lang = ref("ko");
const devType = ref(0);
const role = ref("3"); // 기본값 Admin
const adminpassword = ref("");
const router = useRouter();
const authStore = useAuthStore(); // ✅ Pinia store 사용

function onLangChange(event) {
  const lang = event.target.value
  i18n.global.locale.value = lang
  localStorage.setItem('lang', lang)

  // ✅ html 태그에 lang 속성 적용 (폰트 적용을 위해 꼭 필요)
  document.documentElement.setAttribute('lang', lang)
}

onMounted(async () => {
    try {
      i18n.global.locale.value = lang.value
      localStorage.setItem('lang', lang.value)
      document.documentElement.setAttribute('lang', lang.value)
      //console.log('Installed:',installed.value)
    } catch (error) {
      //console.error("CheckInstalled Error:", error);
    }
});

const signup = async () => {
  if (!email.value || !username.value || !adminpassword.value) {
    alert("Please fill in all the fields.");
    return;
  }

  const data = {
    devType : devType.value,
    username: username.value,
    account: account.value,
    password: adminpassword.value,
    email: email.value,
    role: role.value,
    adminPass: adminpassword.value,
    lang: lang.value
  };

  try {
    const response = await axios.post("/auth/joinAdmin", data, {
      headers: { "Content-Type": "application/json" },
    });

    if (response.data.passOK === "1") {
      alert("✅ Account created successfully!");
      authStore.setInstall(1);
      authStore.setLangDefault(data.lang);
      router.push("/signin"); // 회원가입 후 로그인 페이지로 이동
    } else {
      alert("❌ Account creation failed.: " + (response.data.msg || "Unknown error"));
    }
  } catch (error) {
    console.error("Signup error:", error);
    alert("❌ Error occurred: " + error.message);
  }
};
</script>
