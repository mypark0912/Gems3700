<template>
    <div class="flex h-[100dvh] overflow-hidden">
  
      <!-- Sidebar -->
      <!--Sidebar :sidebarOpen="sidebarOpen" @close-sidebar="sidebarOpen = false" /-->
  
      <!-- Content area -->
      <div class="relative flex flex-col flex-1 overflow-y-auto overflow-x-hidden">
        
        <!-- Site header -->
        <!--Header :sidebarOpen="sidebarOpen" @toggle-sidebar="sidebarOpen = !sidebarOpen" /-->
  
        <main class="grow">
          <div class="px-2 sm:px-4 lg:px-6 py-4 w-full max-w-9xl mx-auto">
  
            <!-- Page header -->
            <div class="sm:flex sm:justify-between sm:items-center mb-6">
  
              <!-- Left: Title -->
              <div class="mb-4 sm:mb-0 flex items-center space-x-4">
  <img src="../../images/logo_transparent.png" alt="LOGO" width="48" height="48" />
  <h1 class="text-2xl md:text-3xl text-gray-800 dark:text-gray-100 font-bold">Device List</h1>
</div>

              
              <div class="grid grid-flow-col sm:auto-cols-max justify-start sm:justify-end gap-2">
              <!-- Search form -->
              <!-- Add member button -->
              <button class="btn bg-gray-900 text-gray-100 hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white" @click="feedbackModalOpen=true">
                <svg class="fill-current shrink-0 xs:hidden" width="16" height="16" viewBox="0 0 16 16">
                  <path d="M15 7H9V1c0-.6-.4-1-1-1S7 .4 7 1v6H1c-.6 0-1 .4-1 1s.4 1 1 1h6v6c0 .6.4 1 1 1s1-.4 1-1V9h6c.6 0 1-.4 1-1s-.4-1-1-1z" />
                </svg>
                <span class="max-xs:sr-only">Device Management</span>
              </button>
              <button
                class="btn-sm border-gray-200 dark:border-gray-700/60 hover:border-gray-300 dark:hover:border-gray-600 text-gray-800 dark:text-white"
                @click.stop="handleReset"
              >
                Reset
              </button>
              <button
                class="btn-sm border-gray-200 dark:border-gray-700/60 hover:border-gray-300 dark:hover:border-gray-600 text-gray-800 dark:text-white"
                @click.stop="handleLogout"
              >
                LOGOUT
              </button>
            </div>
            <ModalBasic
                id="feedback-modal"
                :modalOpen="feedbackModalOpen"
                @close-modal="closeDevice"
                title="Add Device"
              >
                <!-- Modal content -->
                <div class="px-5 py-4">
                  <div class="text-sm">
                    <div
                      class="font-medium text-gray-800 dark:text-gray-100 mb-3"
                    >
                     Device Information
                    </div>
                  </div>
                  <div class="space-y-3">
                    <!-- 콤보박스 추가 -->
                    <div>
                      <label class="block text-sm font-medium mb-1" for="comboBox">Select Item</label>
                      <select id="comboBox" v-model="selectedItem" class="form-select w-full px-2 py-1" @change="handleChange">
                        <option v-for="item in items" :key="item.id" :value="item.id">
                          {{ item.ip }}
                        </option>
                      </select>
                    </div>

                    <div>
                      <label class="block text-sm font-medium mb-1" for="ipaddr">
                        IP Address<span class="text-red-500">*</span>
                      </label>
                      <input
                        v-model="ipaddr"
                        id="ipaddr"
                        class="form-input w-full px-2 py-1"
                        type="text"
                        required
                      />
                    </div>
                    <div>
                      <label class="block text-sm font-medium mb-1" for="location">
                        Location<span class="text-red-500">*</span>
                      </label>
                      <input
                        v-model="location"
                        id="location"
                        class="form-input w-full px-2 py-1"
                        type="text"
                        required
                      />
                    </div>
                    <div class="flex items-center space-x-3 h-10">
                          <!-- ✅ 높이 조정 -->
                          <label
                            class="relative inline-flex items-center cursor-pointer"
                          >
                            <input
                              type="checkbox"
                              v-model="enable"
                              class="sr-only peer"
                            />
                            <div
                              class="w-11 h-6 bg-gray-300 rounded-full relative transition-all peer-checked:bg-violet-500"
                              :class="{
                                'bg-violet-500': enable,
                              }"
                            >
                              <div
                                class="absolute left-1 top-1 bg-white w-4 h-4 rounded-full transition-all transform"
                                :class="
                                  enable
                                    ? 'translate-x-5'
                                    : ''
                                "
                              ></div>
                            </div>
                          </label>

                          <span class="text-sm font-medium">
                            {{
                              enable
                                ? "Enable"
                                : "Disable"
                            }}
                          </span>
                        </div>
                    <div class="text-sm">
                      <div class="font-medium text-gray-800 dark:text-gray-100 mb-3">
                        {{ message }}🙌
                      </div>
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
                      @click.prevent="saveDevice"
                    >
                      Save
                    </button>
                    <button v-if="selectedItem != -1"
                      class="btn-sm border-gray-200 dark:border-gray-700/60 hover:border-gray-300 dark:hover:border-gray-600 text-gray-800 dark:text-white"
                      @click.prevent="deleteDevice"
                    >
                      Delete
                    </button>
                    <button
                      class="btn-sm border-gray-200 dark:border-gray-700/60 hover:border-gray-300 dark:hover:border-gray-600 text-gray-800 dark:text-white"
                      @click.stop="closeDevice"
                    >
                      Close
                    </button>
                  </div>
                </div>
              </ModalBasic>       
  
            </div>
  
            <!-- Cards -->
            <div class="grid grid-cols-12 gap-6">
              <DashCard
                v-for="item in datalist"
                :key="item.id"
                :item="item"
                :user="user"
              />
            </div>
            
            <!-- Pagination -->
            <div class="mt-8">
              <PaginationNumeric />
            </div>           
  
          </div>
        </main>
  
      </div> 
  
    </div>
  </template>
  
  <script>
  import { ref, onMounted, onUnmounted, computed } from 'vue'
  //import Sidebar from '../common/Sidebar.vue'
  import Header from '../common/Header.vue'
  import DashCard from '../../partials/inners/dashboard/DashCard.vue'
  import PaginationNumeric from '../common/PaginationNumeric.vue'
  import axios from 'axios'
  import ModalBasic from "../common/ModalBasic.vue";
  import { useAuthStore } from '@/store/auth'; // ✅ Pinia Store 사용
  import { useRouter } from 'vue-router';
  export default {
    name: 'MasterDashboard',
    components: {
      //Sidebar,
      Header,
      DashCard,
      PaginationNumeric,
      ModalBasic,
    },
    setup() {
      const authStore = useAuthStore();
      const router = useRouter();
      //const sidebarOpen = ref(false)
      const datalist = ref([]);
      const feedbackModalOpen = ref(false);
      const ipaddr = ref('');
      const location = ref('');
      const message = ref('');
      const selectedItem = ref(-1);
      const enable = ref(false);
      const items = ref([{
        "id":-1, "ip": "Add", "location" : ""
      }]);
      let updateInterval = null;
      const user = computed(() => authStore.getUser);
      const fetchData = async () => {
      try {
        //datalist.value = [];
        const response = await axios.get(`/master/getDeviceInfo`);
        if (response.data.result == 1) {
          const newData = response.data.data.map(device => ({
              id: device.id,
              name: device.ip,
              location: device.location,
              enable: device.enable,
              main: device.main,
              sub: device.sub
            }));
            newData.forEach(newItem => {
              const existingItem = datalist.value.find(item => item.id === newItem.id);

              if (existingItem) {
                // 🚀 기존 객체를 유지하면서 값만 업데이트 (전체 객체 교체 X)
                existingItem.name = newItem.name;
                existingItem.location = newItem.location;
                existingItem.enable = newItem.enable;
                existingItem.main = newItem.main;
                existingItem.sub = newItem.sub;
              } else {
                // 새로운 장치 추가
                datalist.value.push(newItem);
              }
            });

            // 기존 데이터에서 삭제된 항목 제거 (newData에 없는 항목 제거)
            datalist.value = datalist.value.filter(item => 
              newData.some(newItem => newItem.id === item.id)
            );

            updateItems();

            // `items` 리스트 업데이트 (렌더링 최적화)
            // items.value = datalist.value.map(({ id, name, location, enable }) => ({
            //   id, ip: name, location, enable
            // }));
        } else {
            console.warn(response.data.error);
        }
      } catch (error) {
        console.log("데이터 가져오기 실패:", error);
      }
    };

    const updateItems = () => {
      // 기존 `items` 배열 유지하면서 `datalist`의 데이터를 추가
      items.value = [
        { id: -1, ip: "Add", location: "" },  // 기본 데이터 유지
        ...datalist.value.map(({ id, name, location, enable }) => ({
          id, ip: name, location, enable
        }))
      ];
    };

    const handleLogout = async () => {
      try {
        if (authStore.logout) { // ✅ logout 함수가 정의되어 있는지 확인
          await authStore.logout();
          router.push("/signin"); // ✅ 로그인 페이지로 이동
        } else {
          console.error("❌ Error: logout 함수가 정의되지 않음!");
        }
      } catch (err) {
        console.error("Logout Error:", err);
      }
    };

    const handleChange = () =>{

      if(selectedItem.value == -1){
        ipaddr.value = "";
        location.value = "";
        enable.value = false;
      }else{
        for(let i = 0 ; i < items.value.length;i++){
          if(Number(selectedItem.value) == Number(items.value[i]["id"])){
            ipaddr.value = items.value[i]["ip"];
            location.value = items.value[i]["location"];
            enable.value = items.value[i]["enable"] == "true"?true:false;
            break;
          }
        }
      }
    };

    const handleReset = async () => {
      try {
        const response = await axios.get("/setting/Reset");
        if (response.data.success){
          message.value = "Setup initiated"
          authStore.setInstall(0);

          if (authStore.logout) {
            await authStore.logout();
            router.push("/signin"); 
          } else {
            console.error("❌ Error: logout 함수가 정의되지 않음!");
          }
          //dangerModalOpen.value = false;
        }else{
          console.error("데이터 가져오기 실패:", error);
        }
      } catch (error) {
        console.error("데이터 가져오기 실패:", error);
      }

    };

    const saveDevice = async () => {
      console.log('CLick Save');
      if (!ipaddr.value) {
        message.value = "IP 주소를 입력하세요.";
        return;
      }

      const formData = new FormData();
      if(selectedItem.value > 0){
        formData.append("id", selectedItem.value);
      }else{
        formData.append("id", "0");
      }
        formData.append("ipaddr", ipaddr.value);
        formData.append("location", location.value);
        formData.append("enable", enable.value);
        try {
            const response = await axios.post("/master/saveDevice", formData, {
                headers: { "Content-Type": "application/json" },
            });
            if (response.data.status == "1") {
              if(selectedItem.value >= 0)
                  message.value = "장치 정보가 수정되었습니다.";
              else
                  message.value = "장치가 추가되었습니다.";
              fetchData();
            } else {
            message.value = response.data.error;
            }
        } catch (error) {
            message.value = "업로드 실패: " + error.response.data.error;
        }
    };

    const deleteDevice = async () => {
      try {
        console.log("Data Get");
        const response = await axios.get(`/master/deleteDevice/${selectedItem.value}`);
        if (response.data.status == "1") {
          message.value = "장치가 삭제되었습니다.";
          fetchData();
        } else {
          message.value = "장치 삭제 실패";
        }
      } catch (error) {
        message.value = "장치 삭제 실패";
      }
    };

    const closeDevice = () =>{
        selectedItem.value = -1;
        ipaddr.value = "";
        location.value = "";
        enable.value = false;
        feedbackModalOpen.value = false;
        message.value="";
    };

    onMounted(()=>{
      fetchData();
      if (!updateInterval) {  // ✅ 중복 실행 방지
        updateInterval = setInterval(() => {
          fetchData();
        }, 60000);
      }
    });

    onUnmounted(() => {
      if (updateInterval) {
        clearInterval(updateInterval);
        updateInterval = null;
      }
    });


  
      return {
        //sidebarOpen,
        datalist,
        fetchData,
        saveDevice,
        ipaddr,
        location,
        feedbackModalOpen,
        message,
        items,
        selectedItem,
        handleChange,
        deleteDevice,
        enable,
        closeDevice,
        handleLogout,
        user,
        handleReset
      }  
    }
  }
  </script>