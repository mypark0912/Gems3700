<template>
  <div class="grow">
    <!-- Panel body -->
    <div class="p-6 space-y-6">
      <!-- Picture -->
      <!-- Business Profile -->
      <section>
        <h3 class="text-xl leading-snug text-gray-800 dark:text-gray-100 font-bold mb-1">{{ t('User.profile.title') }}</h3>
        <div class="py-3">
          <div class="flex items-center w-full max-w-md">
          <label for="account" class="text-sm font-medium text-gray-700 dark:text-white mr-1 w-32">
            {{ t('User.profile.account') }}
          </label>
          <input id="account" class="form-input h-10 px-3" type="text" v-model="user" disabled/>
        </div>
        <div class="flex items-center w-full max-w-md py-2">
            <label for="username" class="text-sm font-medium text-gray-700 dark:text-white w-32 mr-1">
              {{ t('User.profile.username') }}
            </label>
            <input id="username" class="form-input h-10 px-3 flex-1" type="text" v-model="userDict.username" />
          </div>

          <div class="flex items-center w-full max-w-md py-2">
            <label for="email" class="text-sm font-medium text-gray-700 dark:text-white w-32 mr-1">
              {{ t('User.profile.email') }}
            </label>
            <input id="email" class="form-input h-10 px-3 flex-1" type="email" v-model="userDict.email" />
          </div>
        </div>
        </section>
        <section v-if="!devMode && (userRole =='2' || userRole =='3')">
          <div class="flex flex-col w-full max-w-md">
            <label for="email" class="text-sm font-medium text-gray-700 dark:text-white mb-2">
              {{ t('User.profile.useAPI') }}
            </label>
            <select v-model="userDict.api" class="form-select w-1/2 bg-white border-gray-300 rounded-md shadow-sm focus:ring-violet-500 focus:border-violet-500">
              <option value="Yes">Yes</option>
              <option value="No">No</option>
            </select>
        </div>
        </section>
        <section>
        <!-- New Password & Confirm Password (가로 정렬) -->
        <div v-if="!isNtek" class="flex items-start w-full gap-6">
          <!-- 현재 비밀번호 -->
          <div class="flex flex-col w-1/6">
            <label for="Pass" class="text-sm font-medium text-gray-700 dark:text-white mb-2">
              {{ t('User.profile.nowPass') }}
            </label>
            <input id="Pass" class="form-input h-10 px-3" type="password" v-model="userDict.password" placeholder="Enter current password" @focus="hidePassword"/>
          </div>
          <!-- 새 비밀번호 -->
          <div class="flex flex-col w-1/6">
            <label for="newPass" class="text-sm font-medium text-gray-700 dark:text-white mb-2">
              {{ t('User.profile.newPass') }}
            </label>
            <input id="newPass" class="form-input h-10 px-3" type="password" v-model="newPass" placeholder="Enter New password"/>
          </div>
          <!-- 새 비밀번호 확인 -->
          <div class="flex flex-col w-1/6">
            <label for="confirmPass" class="text-sm font-medium text-gray-700 dark:text-white mb-2">
              {{ t('User.profile.cnfPass') }}
            </label>
            <input id="confirmPass" class="form-input h-10 px-3" type="password" v-model="confirmPass" placeholder="Enter New password"/>
          </div>
        </div>

    </section>
      <!-- Smart Sync -->
    </div>
    <!-- Panel footer -->
    <footer>
      <div class="flex px-6 py-5 border-t border-gray-200 dark:border-gray-700/60">
        <!-- ✅ 왼쪽 정렬 -->
        <div class="flex gap-3">
          <button class="btn bg-gray-900 text-gray-100 hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white" 
            @click.prevent="save">
            {{ t('User.btnSave') }}
          </button>
          <!-- <button class="btn dark:bg-gray-800 border-gray-200 dark:border-gray-700/60 hover:border-gray-300 dark:hover:border-gray-600 text-gray-800 dark:text-white">
            Cancel
          </button> -->
        </div>
      </div>
    </footer>
  </div>  
</template>

<script>
import { ref , computed, onMounted, watch, watchEffect} from 'vue'
import { useAuthStore } from '@/store/auth'; 
import axios from 'axios';
import { useI18n } from 'vue-i18n'  // ✅ 추가
export default {
  name: 'AccountPanel',
  setup() {
  const { t } = useI18n();
  const authStore = useAuthStore();
  const user = computed(() => authStore.getUser);
  const userRole = computed(() => authStore.getUserRole);
  const sync = ref('Off')
  const userDict = ref({
    username:'',
    email:'',
    password:'',
    api:'',
  })
  const newPass=ref('');
  const confirmPass = ref('');
  const errorMessage = ref("");
  const isNtek = ref(false);
  const devMode = computed(()=> authStore.getOpMode == 'device0'?true:false);
  
  watch([newPass, confirmPass], () => {
    if (!newPass.value || !confirmPass.value) {
      errorMessage.value = "";
      return;
    }

    errorMessage.value =
      newPass.value !== confirmPass.value
        ? "❌ 입력한 새 비밀번호와 확인 비밀번호가 다릅니다!"
        : "✅ 비밀번호가 일치합니다!";
  });

  const hidePassword = () => {
      userDict.value.password = ''; // 입력 필드를 클릭하면 비워짐
  };

  watchEffect(() => {
        console.log("✅ 사용자 데이터 로드 완료:", userDict.value);
      });

  const fetchData = async () => {
        try {
          const response = await axios.get("/auth/getUser");
          if (response.data.passOK == 1) {
            //Object.assign(userDict.value, response.data[0]); // ✅ 기존 데이터 유지하면서 새로운 데이터 병합
            //userDict.value = { ...response.data.data };
            userDict.value.email = response.data.data["email"];
            userDict.value.username = response.data.data["username"];
            if(parseInt(userRole.value) > 1){
              userDict.value.api = response.data.data["api"];
              if(userRole.value == '3')
                isNtek.value = true;
            }
            else
              userDict.value.api = '';
          }
        } catch (error) {
          console.error("데이터 가져오기 실패:", error);
        }
      };

  const save = async () => {
    if (!userDict.value) {
      console.error("❌ userDict.value가 정의되지 않았습니다!");
      return;
    }

    // ✅ ntek 사용자가 아닌 경우에만 비밀번호 검증
    if (!isNtek.value) {
      // 새 비밀번호가 입력되었는지 확인
      const isPasswordChange = newPass.value.trim() !== '' || confirmPass.value.trim() !== '';
      
      if (isPasswordChange) {
        // 새 비밀번호 변경 시 검증
        if (!newPass.value.trim()) {
          alert('Please enter a new password.');
          return;
        }
        if (!confirmPass.value.trim()) {
          alert('Please confirm your new password.');
          return;
        }
        if (newPass.value !== confirmPass.value) {
          alert('The new password and the confirmation password do not match.');
          return;
        }
        if (!userDict.value.password.trim()) {
          alert('Please enter your current password.');
          return;
        }
      }
    }

    // ✅ 데이터 구성 - 새 비밀번호가 비어있으면 현재 비밀번호 유지
    const data = {
      email: userDict.value.email,
      username: userDict.value.username,
      password: userDict.value.password,
      newPass: newPass.value.trim() || userDict.value.password, // ✅ 비어있으면 현재 비밀번호 사용
      role: userRole.value,
      api: userDict.value.api
    };

    try {
      const response = await axios.post('/auth/updateProfile', data, {
        headers: { 'Content-Type': 'application/json' },
        withCredentials: true,
      });
      
      if (response.data?.passOK === "1") {
        alert('The information has been updated.');
        // ✅ 저장 후 비밀번호 필드 초기화
        userDict.value.password = '';
        newPass.value = '';
        confirmPass.value = '';
      } else if (response.data?.passOK === "2") {
        alert('Failed to update the diagnostic service information.');
      } else if (response.data?.passOK === "3") {
        alert('The diagnostic service is not running.');
      } else {
        alert('Incorrect password!');
      }
    } catch (err) {
      console.error("Update Error:", err);
      alert('An error occurred while updating information.');
    }
  };

  onMounted(fetchData);

  return {
    sync,
    user,
    userDict,
    save,
    newPass,
    confirmPass,
    errorMessage,
    hidePassword,
    userRole,
    t,
    isNtek,
    devMode,
  }
}
}
</script>