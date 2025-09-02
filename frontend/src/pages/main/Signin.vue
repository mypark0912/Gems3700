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
                <svg width="32" height="32" viewBox="0 0 32 32">
                  <defs>
                    <linearGradient
                      x1="28.538%"
                      y1="20.229%"
                      x2="100%"
                      y2="108.156%"
                      id="logo-a"
                    >
                      <stop stop-color="#B7ACFF" stop-opacity="0" offset="0%" />
                      <stop stop-color="#B7ACFF" offset="100%" />
                    </linearGradient>
                    <linearGradient
                      x1="88.638%"
                      y1="29.267%"
                      x2="22.42%"
                      y2="100%"
                      id="logo-b"
                    >
                      <stop stop-color="#7BC8FF" stop-opacity="0" offset="0%" />
                      <stop stop-color="#7BC8FF" offset="100%" />
                    </linearGradient>
                  </defs>
                  <rect fill="#8470FF" width="32" height="32" rx="16" />
                  <path
                    d="M18.277.16C26.035 1.267 32 7.938 32 16c0 8.837-7.163 16-16 16a15.937 15.937 0 01-10.426-3.863L18.277.161z"
                    fill="#755FF8"
                  />
                  <path
                    d="M7.404 2.503l18.339 26.19A15.93 15.93 0 0116 32C7.163 32 0 24.837 0 16 0 10.327 2.952 5.344 7.404 2.503z"
                    fill="url(#logo-a)"
                  />
                  <path
                    d="M2.223 24.14L29.777 7.86A15.926 15.926 0 0132 16c0 8.837-7.163 16-16 16-5.864 0-10.991-3.154-13.777-7.86z"
                    fill="url(#logo-b)"
                  />
                </svg>
              </router-link>
            </div>
          </div>

          <div class="max-w-sm mx-auto w-full px-4 py-8">
            <h1
              class="text-3xl text-gray-800 dark:text-gray-100 font-bold mb-6"
            >
              {{ t('Title')}}
            </h1>
            
            <!-- Error Message Alert -->
            <div v-if="errorMessage" class="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
              {{ errorMessage }}
            </div>
            
            <!-- Form -->
            <form @submit.prevent="login">
              <div class="space-y-4">
                <div>
                  <label class="block text-sm font-medium mb-1" for="username"
                    >{{ t('signin.account') }}</label
                  >
                  <input
                    id="username"
                    ref="usernameInput"
                    class="form-input w-full px-2 py-1 border rounded"
                    v-model="username"
                    type="text"
                    placeholder="Enter username"
                    :disabled="isLoading"
                  />
                </div>
                <div>
                  <label class="block text-sm font-medium mb-1" for="language"
                    >{{ t('signin.language') }}</label
                  >
                  <select
                    id="language"
                    class="form-input w-full px-2 py-1 border rounded bg-white"
                    v-model="lang"  
                    @change="onLangChange"
                    :disabled="isLoading"
                  >
                    <option value="ko">KOREAN</option>
                    <option value="en">ENGLISH</option>
                    <option value="ja">JAPANESE</option>
                  </select>
                </div>
                <div>
                  <label class="block text-sm font-medium mb-1" for="password"
                    >{{ t('signin.password') }}</label
                  >
                  <input
                    v-model="password"
                    ref="passwordInput"
                    class="form-input w-full px-2 py-1 border rounded"
                    type="password"
                    autoComplete="on"
                    id="password"
                    :disabled="isLoading"
                  />
                </div>
              </div>
              <div class="flex items-center justify-between mt-6">
                <div class="mr-1">
                  <div class="text-sm">
                    <a
                      class="font-medium text-violet-500 hover:text-violet-600 dark:hover:text-violet-400"
                      href="#"
                      @click="feedbackModalOpen = true"
                      >{{ t('signin.context') }}
                    </a>
                  </div>
                </div>
                <button
                  type="submit"
                  class="btn bg-gray-900 text-gray-100 hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white whitespace-nowrap ml-3 disabled:opacity-50 disabled:cursor-not-allowed"
                  :disabled="isLoading"
                >
                  <span v-if="isLoading">{{ t('signin.logging_in') || 'Logging in...' }}</span>
                  <span v-else>{{ t('signin.login') }}</span>
                </button>
              </div>
            </form>
            
            <!-- Modal and footer sections remain the same -->
            <ModalBasic
              id="feedback-modal"
              :modalOpen="feedbackModalOpen"
              @close-modal="feedbackModalOpen = false"
              title="Password Reset"
            >
              <!-- Modal content -->
              <div class="px-5 py-4">
                <div class="space-y-3">
                  <div>
                    <label class="block text-sm font-medium mb-1" for="account"
                      >{{ t('resetPass.account') }} <span class="text-red-500">*</span></label
                    >
                    <input
                      id="account"
                      class="form-input w-full px-2 py-1 border rounded"
                      type="text"
                      v-model="reAccount"
                      required
                    />
                  </div>

                  <div>
                    <label class="block text-sm font-medium mb-1" for="email"
                      >{{ t('resetPass.email') }} <span class="text-red-500">*</span></label
                    >
                    <input
                      id="email"
                      class="form-input w-full px-2 py-1 border rounded"
                      type="email"
                      v-model="email"
                      required
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium mb-1" for="password"
                      >{{ t('resetPass.password') }} <span class="text-red-500">*</span></label
                    >
                    <input
                      id="password"
                      class="form-input w-full px-2 py-1 border rounded"
                      type="password"
                      v-model="rePass"
                      required
                    />
                  </div>
                </div>
                <div class="text-sm">
                  <div
                    class="font-medium text-gray-800 dark:text-gray-100 mb-3"
                  >
                    {{ message }}🙌
                  </div>
                </div>
              </div>

              <!-- Modal footer -->
              <div
                class="px-5 py-4 border-t border-gray-200 dark:border-gray-700/60"
              >
                <div class="flex flex-wrap justify-end space-x-2">
                  <button
                    class="btn-sm bg-gray-900 text-gray-100 hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white"
                    @click.prevent="reset"
                  >
                  {{ t('resetPass.save') }}
                  </button>
                  <button
                    class="btn-sm border-gray-200 dark:border-gray-700/60 hover:border-gray-300 dark:hover:border-gray-600 text-gray-800 dark:text-white"
                    @click.stop="feedbackModalOpen = false"
                  >
                  {{ t('resetPass.close') }}
                  </button>
                </div>
              </div>
            </ModalBasic>
            <!-- Footer -->
            <div
              v-if="installed == 0"
              class="pt-5 mt-6 border-t border-gray-100 dark:border-gray-700/60"
            >
              <div class="text-sm">
                {{ t('signin.msg')}}
                <a
                  class="font-medium text-violet-500 hover:text-violet-600 dark:hover:text-violet-400 cursor-pointer"
                  @click="Signup"
                >
                {{ t('signin.link')}}
                </a>
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
import { ref, computed, onMounted, watchEffect, watch, nextTick } from "vue";
import { useAuthStore } from '@/store/auth'
import { useSetupStore } from "@/store/setup";
import { useRouter } from "vue-router";
import ModalBasic from "../common/ModalBasic.vue";
import axios from "axios";
import i18n from '@/i18n'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const username = ref("");
const password = ref("");
const lang = ref("ko");
const router = useRouter();
const authStore = useAuthStore();
const setupStore = useSetupStore();
const feedbackModalOpen = ref(false);
const reAccount = ref("");
const rePass = ref("");
const email = ref("");
const message = ref("Enter the Password");

// Added variables
const errorMessage = ref("");
const isLoading = ref(false);

// ref elements
const usernameInput = ref(null);
const passwordInput = ref(null);

const installed = computed(() => authStore.getInstalled);
const devMode = computed(() => authStore.getOpMode);
const defaultLang = computed(()=> authStore.getLangDefault);

function onLangChange(event) {
  lang.value = event.target.value
  i18n.global.locale.value = lang.value
  localStorage.setItem('lang', lang.value)
  document.documentElement.setAttribute('lang', lang.value)
}

onMounted(async () => {
    try {
      await authStore.checkInstalled();
      await setupStore.checkSetting();
      
      // Focus on account input field when page loads
      nextTick(() => {
        usernameInput.value?.focus();
      });
    } catch (error) {
      console.error(error);
    }
});

watch(defaultLang, (newLang) => {
  if (newLang && newLang !== lang.value) {
    lang.value = newLang;
    i18n.global.locale.value = newLang;
    localStorage.setItem('lang', newLang);
    document.documentElement.setAttribute('lang', newLang);
  }
}, { immediate: true });

async function waitForDevMode() {
  return new Promise((resolve) => {
    const stop = watchEffect(() => {
      const val = devMode.value;
      if (val && val !== 'undefined') {
        resolve(val);
        queueMicrotask(() => stop());
      }
    });
  });
}

// Modified login function
const login = async () => {
  // Initialize error message
  errorMessage.value = "";
  
  // Input validation
  if (!username.value) {
    errorMessage.value = t('signin.error_account_required') || "Please enter your account.";
    // Focus on account input field
    nextTick(() => {
      usernameInput.value?.focus();
    });
    return;
  }
  
  if (!password.value) {
    errorMessage.value = t('signin.error_password_required') || "Please enter your password.";
    // Focus on password input field
    nextTick(() => {
      passwordInput.value?.focus();
    });
    return;
  }
  
  isLoading.value = true;
  
  const data = {
    account: username.value,
    password: password.value,
    lang: lang.value,
  };

  try {
    const result = await authStore.login(data);
    
    if(result && result.passOK){
      devMode.value = result.mode;
      //console.log('Login successful, devMode:', devMode.value);
      
      // Navigate on successful login
      if (devMode.value.includes("device")) {
        router.push({ name: "dashboard" });
      } else {
        router.push({ name: "master" });
      }
    } else {
      // Show error message on login failure
      const errorCode = result?.errorCode || "4";
      
      switch(errorCode) {
        case "4":
          errorMessage.value = t('signin.error_wrong_password') || "The password is incorrect.";
          password.value = ""; // Clear password field
          // Focus on password field
          nextTick(() => {
            passwordInput.value?.focus();
          });
          break;
        case "5":
          errorMessage.value = t('signin.error_account_not_found') || "Account does not exist.";
          username.value = ""; // Clear account field
          password.value = ""; // Clear password field too
          // Focus on account field
          nextTick(() => {
            usernameInput.value?.focus();
          });
          break;
        default:
          errorMessage.value = t('signin.error_invalid_credentials') || "Invalid account or password.";
          password.value = ""; // Clear password field
          // Focus on password field
          nextTick(() => {
            passwordInput.value?.focus();
          });
      }
    }
  } catch (error) {
    console.error("Login error:", error);
    errorMessage.value = t('signin.error_login_failed') || "An error occurred during login. Please try again.";
    // Focus on account field for network errors
    nextTick(() => {
      usernameInput.value?.focus();
    });
  } finally {
    isLoading.value = false;
  }
};

const Signup = () => {
  router.push({ name: "signup" });
};

const reset = async () => {
  if (email.value == "") {
    message.value = "Please enter your email address.";
  } else if (rePass.value == "") {
    message.value = "Please enter your password.";
  } else {
    const data = {
      username: reAccount.value,
      password: rePass.value,
      email: email.value,
    };
    //console.log('🚀 Sending reset password request:', data); 
    try {
      const response = await axios.post("/auth/resetPassword", data, {
        headers: { "Content-Type": "application/json" },
        withCredentials: true,
      });
      if (response.data?.passOK === "1") {
        message.value = t('resetPass.msg_success');
        // Show success message briefly then close modal
        setTimeout(() => {
          feedbackModalOpen.value = false;
          // Clear input fields after modal closes
          reAccount.value = "";
          rePass.value = "";
          email.value = "";
          message.value = "Enter the Password";
        }, 1500); // Close after 1.5 seconds
      } else if (response.data?.passOK === "0") {
        message.value = t('resetPass.msg_fail');          
      } else if (response.data?.passOK === "2") {
        message.value = t('resetPass.msg_undefined');  
      }
    } catch (err) {
      console.error("Reset Password Error:", err);
      message.value = "An error occurred while resetting the password.";
    }
  }
};
</script>