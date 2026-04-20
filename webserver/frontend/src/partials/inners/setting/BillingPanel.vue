<template>
  <div class="grow">
    <!-- Panel body -->
    <div class="p-6 space-y-6">
      <!-- Billing Information -->
      <section class="flex flex-col col-span-full xl:col-span-6">
        <h3 class="text-xl leading-snug text-gray-800 dark:text-gray-100 font-bold mb-3">{{ t('User.user.title') }}</h3>
        <div class="text-sm text-violet-500 mb-0.5">{{ t('User.user.msg') }}</div>
        <div class="flex py-3 gap-6">
          <!-- 왼쪽: 사용자 리스트 + Add User 버튼 -->
          <div class="w-1/4 flex flex-col">
            <ul class="border rounded overflow-y-auto max-h-80">
              <li
                v-for="item in userlist"
                :key="item.username"
                @click="selectAccount(item)"
                :class="[
                  'px-4 py-2 cursor-pointer border-b',
                  selectedUser?.username === item.username
                    ? 'bg-gray-200 dark:bg-gray-700 text-black dark:text-white font-semibold'
                    : 'hover:bg-gray-100 dark:hover:bg-gray-800'
                ]"
              >
                {{ item.account }}
              </li>
            </ul>
            <div class="mt-4 flex justify-center">
              <button
                @click="add"
                class="btn w-full dark:bg-gray-800 border-gray-200 dark:border-gray-700/60 
                       hover:border-gray-300 dark:hover:border-gray-600 text-gray-800 dark:text-white"
              >
              {{ t('User.user.btnAdd') }}
              </button>
            </div>
          </div>

          <!-- 오른쪽: 사용자 정보 입력 or 수정 -->
          <div v-if="addUser" class="flex-1 border rounded p-4 space-y-4">
            <form>
              <div class="space-y-4">
                <div>
                  <label class="block text-sm font-medium mb-1" for="name">
                    {{ t('User.user.username') }} <span class="text-red-500">*</span>
                  </label>
                  <input id="name" class="form-input w-full" type="text" v-model="username" />
                </div>
                <div>
                  <label class="block text-sm font-medium mb-1" for="email">
                    {{ t('User.user.email') }} <span class="text-red-500">*</span>
                  </label>
                  <input id="email" class="form-input w-full" type="email" v-model="email" />
                </div>
                <div>
                  <label class="block text-sm font-medium mb-1" for="account">
                    {{ t('User.user.account') }} <span class="text-red-500">*</span>
                  </label>
                  <input id="account" class="form-input w-full" type="text" v-model="account" />
                </div>
                <div>
                  <label class="block text-sm font-medium mb-1" for="role">
                    {{ t('User.user.role') }}<span class="text-red-500">*</span>
                  </label>
                  <select id="role" class="form-select w-full" v-model="role">
                    <option value="1">User</option>
                    <option value="0">Operator</option>
                    <option value="2">Admin</option>
                  </select>
                </div>
                <div>
                  <label class="block text-sm font-medium mb-1" for="password">
                    {{ t('User.user.Pass') }}
                  </label>
                  <input
                    id="password"
                    class="form-input w-full"
                    type="password"
                    autoComplete="on"
                    v-model="password"
                  />
                </div>
              </div>
            </form>
          </div>

          <div
            v-if="!addUser && selectedUser.account !== ''"
            class="flex-1 border rounded p-4 space-y-4"
          >
            <!-- 사용자 이름 -->
            <h4 class="text-lg font-semibold">{{ t('User.user.account') }} : {{ selectedUser.account }}</h4>
            <!-- Role 수평 정렬 -->
            <div class="flex items-center gap-4 mb-4">
              <label class="text-sm font-medium">{{ t('User.user.username') }}</label>
              <input
                type="text"
                class="w-48 h-9 form-input text-sm dark:bg-gray-800 dark:text-white" v-model="selectedUser.username" disabled
              />
            </div>
            <div class="flex items-center gap-4 mb-4">
              <label class="text-sm font-medium">{{ t('User.user.role') }}</label>
              <select
                v-model="selectedUser.role"
                class="w-48 px-2 py-1 border rounded text-sm dark:bg-gray-800 dark:text-white"
              >
                <option value="1">User</option>
                <option value="0">Operator</option>
                <option value="2">Admin</option>
              </select>
            </div>
            <div v-if="selectedUser.role !== '2'">
              <label class="flex items-center space-x-2">
              <input
                type="checkbox"
                class="form-checkbox text-violet-500 focus:ring-violet-500"
                v-model="removeUser"
              />
              <span class="text-sm">{{ t('User.user.remove') }}</span>
            </label> 
            </div>
          </div>
        </div>
      </section>
    </div>

    <footer v-if="selectedUser.account !== '' || addUser">
      <div class="flex px-6 py-5 border-t border-gray-200 dark:border-gray-700/60">
        <!-- ✅ 왼쪽 정렬 -->
        <div class="flex gap-3">
          <button
            class="btn bg-gray-900 text-gray-100 hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white"
            @click.prevent="save"
          >
          {{ t('User.btnSave') }}
          </button>
        </div>
      </div>
    </footer>
  </div>
</template>

<script>
import { ref , onMounted, computed } from 'vue'
import axios from 'axios';
import { useAuthStore } from '@/store/auth'; 
import { useI18n } from 'vue-i18n'  // ✅ 추가

export default {
  name: 'BillingPanel',
  props:['mode'],
  setup(props){
    const { t } = useI18n();
    const authStore = useAuthStore();
    const loginUser = computed(() => authStore.getUser);
    const userlist = ref([]);
    const addUser = ref(false);
    const email = ref("");
    const account = ref("");
    const username = ref("");
    const password = ref("");
    const role = ref("2"); // 기본값 Admin
    const selectedUser = ref({
      username:'',
      account: '',
      role:'',
      api :'',
    });
    const mode = ref(props.mode);
    const removeUser = ref(false);
    const fetchData = async () => {
      try {
        const response = await axios.get("/auth/getUserList");
        if (response.data.passOK == 1) {
          //userlist.value = response.data.data;
          if(loginUser.value == 'ntek'){
            userlist.value = response.data.data.filter(user =>
                              user.account != "ntek" &&
                              !["operator", "guest", "admin"].includes(user.account)
                            );
          }else{
            userlist.value = response.data.data.filter(user =>                                                
                              parseInt(user.role) < 2 &&
                              !["operator", "guest", "admin"].includes(user.account)
                            );
          }
        }
      }catch (error) {
        console.error("데이터 가져오기 실패:", error);
      }
    };

    const selectAccount = (user) => {
      selectedUser.value = user;
      addUser.value = false;
    };

    const add = () => {
      addUser.value = true;
      removeUser.value = false;
      selectedUser.value = { account:'', username: '', role: '', api: '' };
    };

    const save = async () => {
      if(addUser.value){
        if (!username.value || !account.value || !password.value) {
          alert("Please fill in all the fields.");
          return;
        }
        if (account.value == 'operator' || account.value == 'guest'|| account.value == 'admin') {
          alert("This account cannot be used.");
          return;
        }
        const data = {
          account: account.value,
          username: username.value,
          password: password.value,
          email: email.value,
          role: role.value,
        };
        try {
          const response = await axios.post("/auth/join", data, {
            headers: { "Content-Type": "application/json" },
          });

          if (response.data.passOK === "1") {
            alert("✅ The account has been created successfully!");
            fetchData();
          } else {
            alert("❌ Account creation failed: " + (response.data.msg || "Unknown error"));
          }
        } catch (error) {
          console.error("Signup error:", error);
          alert("❌ Error occurred: " + error.message);
        }
        selectedUser.value.account ='';
        addUser.value = false;
        removeUser.value = false;
      }else{
        if(selectedUser.value.username == ''){
          alert('Please select a user.');
          return;
        }
        if(removeUser.value){
            try {
              const response = await axios.get(`/auth/removeUser/${selectedUser.value.account}`);
              if (response.data.passOK == "1"){
                alert('The user has been deleted.');
                fetchData();
              }
              // else if (response.data.passOK == "2")
              //   alert('No subscription information exists for the diagnostic service.');
              // else if (response.data.passOK == "3")
              //   alert('The diagnostic service is not running.');
              else
                alert('Failed to delete the user.');
              // 저장 성공 처리
            } catch (error) {
              console.error("Deletion failed.:", error);
            }
        }else{
          try {
              const response = await axios.get(`/auth/saveUser/${selectedUser.value.account}/${selectedUser.value.role}`);
              if (response.data.passOK == "1")
                alert('The information has been changed.');
              else if (response.data.passOK == "2")
                alert('No subscription information is available for the diagnostic service.');
              else if (response.data.passOK == "3")
                alert('The diagnostic service is not running.');
              else
                alert('Failed to change permissions');
              // 저장 성공 처리
            } catch (error) {
              console.error("Failed to save:", error);
            }
        }
      }
    };

    onMounted(()=>{
      fetchData();
    });

    return {
      userlist,
      save,
      selectedUser,
      addUser,
      add,
      username,
      password,
      email,
      role,
      removeUser,
      selectAccount,
      mode,
      account,
      loginUser,
      t,
    }
  }
}
</script>