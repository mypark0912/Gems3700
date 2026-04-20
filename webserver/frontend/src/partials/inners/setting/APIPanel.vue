<template>
    <div class="grow">
      <!-- Panel body -->
      <div class="p-6 space-y-6">
        <!-- Billing Information -->
        <section class="flex flex-col col-span-full xl:col-span-6">
          <h3 class="text-xl leading-snug text-gray-800 dark:text-gray-100 font-bold mb-3">{{ t('User.api.title')}}</h3>
          <div v-if="userlist.length > 0" class="text-sm text-violet-500 mb-0.5">{{ t('User.api.msg')}}</div>
          <div class="flex py-3 gap-6">
            <!-- 왼쪽: 사용자 리스트 + Add User 버튼 -->
            <div class="w-1/4 flex flex-col">
              <div v-if="!addUser && userlist.length == 0">{{ t('User.api.noexist')}}</div>
              <ul v-else class="border rounded overflow-y-auto max-h-80">
                <li
                  v-for="item in userlist"
                  :key="item.account"
                  @click="selectAccount(item)"
                  :class="[
                    'px-4 py-2 cursor-pointer border-b',
                    selectedUser?.account === item.Username
                      ? 'bg-gray-200 dark:bg-gray-700 text-black dark:text-white font-semibold'
                      : 'hover:bg-gray-100 dark:hover:bg-gray-800'
                  ]"
                >
                  {{ item.Username }}
                </li>
              </ul>
              <div class="mt-4 flex justify-center">
                <button
                  @click="add"
                  class="btn w-full dark:bg-gray-800 border-gray-200 dark:border-gray-700/60 
                         hover:border-gray-300 dark:hover:border-gray-600 text-gray-800 dark:text-white"
                >
                {{ t('User.api.btnAdd')}}
                </button>
              </div>
            </div>
  
            <!-- 오른쪽: 사용자 정보 입력 or 수정 -->
            <div v-if="addUser" class="flex-1 border rounded p-4 space-y-4">
              <form>
                <div class="space-y-4">
                  <div>
                    <label class="block text-sm font-medium mb-1" for="role">
                      {{ t('User.api.role')}} <span class="text-red-500">*</span>
                    </label>
                    <select id="role" class="form-select w-1/2" v-model="role">
                      <option value="2">User</option>
                      <option value="1">Operator</option>
                    </select>
                  </div>
                </div>
              </form>
            </div>
  
            <div
              v-if="!addUser && selectedUser.account !== ''"
              class="flex-1 border rounded p-4 space-y-4"
            >
              <!-- 사용자 이름 -->
              <h4 class="text-lg font-semibold">{{ t('User.api.account')}} : {{ selectedUser.account }}</h4>
  
              <!-- Role 수평 정렬 -->
              <div class="flex items-center gap-4 mb-4">
                <label class="text-sm font-medium">{{ t('User.api.role')}}</label>
                <select
                  v-model="selectedUser.role"
                  class="w-24 px-2 py-1 border rounded text-sm dark:bg-gray-800 dark:text-white" disabled
                >
                  <option value="2">User</option>
                  <option value="1">Operator</option>
                </select>
              </div>
              <div>
                <label class="flex items-center space-x-2">
                <input
                  type="checkbox"
                  class="form-checkbox text-violet-500 focus:ring-violet-500"
                  v-model="removeUser"
                />
                <span class="text-sm">{{ t('User.api.remove')}}</span>
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
            {{ t('User.btnSave')}}
            </button>
          </div>
        </div>
      </footer>
    </div>
  </template>
  
  <script>
  import { ref , onMounted } from 'vue'
  import axios from 'axios';
  import { useI18n } from 'vue-i18n'  // ✅ 추가
  export default {
    name: 'APIPanel',
    setup(){
      const { t } = useI18n();
      const userlist = ref([]);
      const addUser = ref(false);
      const username = ref("");
      const password = ref("");
      const role = ref("2"); // 기본값 Admin
      const selectedUser = ref({
        account:'',
        role:'',
        api :'',
      });
      const removeUser = ref(false);
      const fetchData = async () => {
        try {
          const response = await axios.get("/auth/getAPIUsers");
          if (response.data.success) {
            //userlist.value = response.data.data.filter(user => user.Username !== "admin");
            userlist.value = response.data.data.filter(user =>
                              user.Level !== "0" &&
                              ["operator", "guest"].includes(user.Username)
                            );
          }
        } catch (error) {
          console.error("데이터 가져오기 실패:", error);
        }
      };
  
      const selectAccount = (user) => {
        // selectedUser.value = user;
        selectedUser.value.account = user.Username;
        selectedUser.value.role = user.Level;
        selectedUser.value.api ='Yes';
        addUser.value = false;
      };
  
      const add = () => {
        addUser.value = true;
        removeUser.value = false;
        selectedUser.value = { account: '', role: '', api: '' };
      };
  
      const save = async () => {
        if(addUser.value){
          try {
            const response = await axios.get(`/auth/joinAPI/${role.value}`);
  
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
          if(selectedUser.value.account == ''){
            alert('Please select a user.');
            return;
          }
          if(removeUser.value){
              
              try {
                const response = await axios.get(`/auth/removeAPI/${selectedUser.value.role}`);
                if (response.data.passOK == "1"){
                  alert('The user has been deleted.');
                  fetchData();
                }
                else if (response.data.passOK == "2")
                  alert('No subscription information exists for the diagnostic service.');
                else if (response.data.passOK == "3")
                  alert('The diagnostic service is not currently running.');
                else
                  alert('Failed to delete the user');
                // 저장 성공 처리
              } catch (error) {
                console.error("Failed to delete:", error);
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
        role,
        removeUser,
        selectAccount,
        t,
      }
    }
  }
  </script>