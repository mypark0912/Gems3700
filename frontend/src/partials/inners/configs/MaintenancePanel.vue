<template>
    <div class="grow">
  
      <!-- Panel body -->
      <div class="p-6">
        <h2 class="text-2xl text-gray-800 dark:text-gray-100 font-bold mb-5"> {{ t('config.maintenance.header') }} </h2>
  
        <!-- General -->
        <div class="mb-6">
          <!-- Filters -->
          <div class="mt-4 flex gap-2 border-b border-gray-200 dark:border-gray-700/60 p-2 pb-4">
            <div v-if="message"
            class="font-medium text-gray-800 dark:text-gray-100 mb-3"
          >
            {{ message }}🙌
          </div>
          </div>
        </div>
  
        <!-- Connected Apps cards -->
        <section class="pb-6 border-b border-gray-200 dark:border-gray-700/60">
            <div class="flex flex-col space-y-10 sm:flex-row sm:space-x-6 sm:space-y-0 md:flex-col md:space-x-0 md:space-y-10 xl:flex-row xl:space-x-6 xl:space-y-0 mt-9">
                <div class="w-full">
                  <div class="bg-white dark:bg-gray-800 shadow-lg rounded-xl p-5 mb-2">
                    <div class="space-y-6">
                      <!-- Title + Version -->
                      <div class="flex flex-col sm:flex-row sm:items-end gap-4">
                        <div class="flex flex-col w-full">
                          <label for="title-input" class="text-sm font-medium text-gray-700 dark:text-gray-300 w-20">{{ t('config.maintenance.Title') }}</label>
                          <input
                            id="title-input" v-model="title"
                            class="form-input flex-1 bg-gray-100 dark:bg-gray-700 border-transparent dark:border-transparent focus:bg-white dark:focus:bg-gray-800 placeholder-gray-500"
                            type="text"
                            placeholder=""
                          />
                        </div>
                      </div>
                      <div class="flex flex-col sm:flex-row sm:items-end gap-4">
                        <div class="flex space-x-5">
                        <label class="inline-flex items-center text-sm font-medium cursor-pointer text-gray-400 peer-checked:text-violet-500">
                          <input type="radio" class="sr-only peer" name="action" value="0" v-model="mtype" />
                          <svg class="shrink-0 mr-2 stroke-current peer-checked:text-violet-500 transition-colors" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
                            <polyline points="7 10 12 15 17 10" />
                            <line x1="12" y1="15" x2="12" y2="3" />
                          </svg>
                          <span class="peer-checked:text-violet-500">{{ t('config.maintenance.Installation') }}</span>
                        </label>

                        <label class="inline-flex items-center text-sm font-medium cursor-pointer text-gray-400 peer-checked:text-violet-500">
                          <input type="radio" class="sr-only peer" name="action" value="1" v-model="mtype" />
                          <svg class="shrink-0 mr-2 stroke-current peer-checked:text-violet-500 transition-colors" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <polyline points="23 4 23 10 17 10" />
                            <polyline points="1 20 1 14 7 14" />
                            <path d="M3.51 9a9 9 0 0 1 14.13-3.36L23 10" />
                            <path d="M20.49 15a9 9 0 0 1-14.13 3.36L1 14" />
                          </svg>
                          <span class="peer-checked:text-violet-500">{{ t('config.maintenance.Update') }}</span>
                        </label>

                        <label class="inline-flex items-center text-sm font-medium cursor-pointer text-gray-400 peer-checked:text-violet-500">
                          <input type="radio" class="sr-only peer" name="action" value="2" v-model="mtype" />
                          <svg class="shrink-0 mr-2 stroke-current peer-checked:text-violet-500 transition-colors" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <polyline points="16 3 21 3 21 8" />
                            <line x1="4" y1="20" x2="21" y2="3" />
                            <polyline points="8 21 3 21 3 16" />
                          </svg>
                          <span class="peer-checked:text-violet-500">{{ t('config.maintenance.Replacement') }}</span>
                        </label>
                      </div>
                      <div v-if="mtype == 0 || mtype == 1" class="flex space-x-5">
                          <!-- Installation -->
                          <label class="inline-flex items-center text-sm font-medium cursor-pointer text-gray-600 dark:text-gray-300">
                            <input
                              type="checkbox"
                              class="form-checkbox text-violet-500 border-gray-300 dark:border-gray-600 focus:ring-violet-400"
                              v-model="selectedTarget.fw"
                              value="install"
                            />
                            <span class="ml-2">{{ t('config.maintenance.FirmWare') }}</span>
                          </label>

                          <!-- Update -->
                          <label class="inline-flex items-center text-sm font-medium cursor-pointer text-gray-600 dark:text-gray-300">
                            <input
                              type="checkbox"
                              class="form-checkbox text-violet-500 border-gray-300 dark:border-gray-600 focus:ring-violet-400"
                              v-model="selectedTarget.app"
                              value="update"
                            />
                            <span class="ml-2">{{ t('config.maintenance.App') }}</span>
                          </label>

                          <!-- Replacement -->
                          <label class="inline-flex items-center text-sm font-medium cursor-pointer text-gray-600 dark:text-gray-300">
                            <input
                              type="checkbox"
                              class="form-checkbox text-violet-500 border-gray-300 dark:border-gray-600 focus:ring-violet-400"
                              v-model="selectedTarget.web"
                              value="replace"
                            />
                            <span class="ml-2">{{ t('config.maintenance.Web') }}</span>
                          </label>
                        </div>
                      </div>

                      <div v-if="!isNoneSelected" class="flex flex-col sm:flex-row sm:items-end gap-4">
                        <div v-if="selectedTarget.fw" class="flex items-center space-x-2">
                          <label for="version-input" class="text-sm font-medium text-gray-700 dark:text-gray-300 w-20">{{ t('config.maintenance.FirmWare') }}</label>
                          <input
                            id="version-input"
                            class="form-input w-48 bg-gray-100 dark:bg-gray-700 border-transparent dark:border-transparent focus:bg-white dark:focus:bg-gray-800 placeholder-gray-500"
                            type="text"
                            placeholder=""
                            v-model="versions.fw"
                          />
                        </div>
                        <div v-if="selectedTarget.app" class="flex items-center space-x-2">
                          <label for="version-input" class="text-sm font-medium text-gray-700 dark:text-gray-300 w-10">{{ t('config.maintenance.App') }}</label>
                          <input
                            id="version-input"
                            class="form-input w-48 bg-gray-100 dark:bg-gray-700 border-transparent dark:border-transparent focus:bg-white dark:focus:bg-gray-800 placeholder-gray-500"
                            type="text"
                            placeholder="" v-model="versions.app"
                          />
                        </div>
                        <div v-if="selectedTarget.web" class="flex items-center space-x-2">
                          <label for="version-input" class="text-sm font-medium text-gray-700 dark:text-gray-300 w-10">{{ t('config.maintenance.Web') }}</label>
                          <input
                            id="version-input"
                            class="form-input w-48 bg-gray-100 dark:bg-gray-700 border-transparent dark:border-transparent focus:bg-white dark:focus:bg-gray-800 placeholder-gray-500"
                            type="text"
                            placeholder="" v-model="versions.web"
                          />
                        </div>
                      </div>
                      <!-- Context + Send -->
                      <div class="flex flex-col sm:flex-row sm:items-end gap-4">
                        <div class="flex flex-col w-full">
                          <label for="context-input" class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">{{ t('config.maintenance.Context') }}</label>
                          <textarea
                            id="context-input" v-model="context"
                            rows="3"
                            class="form-textarea w-full bg-gray-100 dark:bg-gray-700 border-transparent dark:border-transparent focus:bg-white dark:focus:bg-gray-800 placeholder-gray-500"
                            placeholder=""
                          ></textarea>
                        </div>
                      </div>
                      <div class="flex gap-2">
                        <button @click="save" class="btn-sm w-32 bg-gray-900 text-gray-100 hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white whitespace-nowrap">{{ t('config.maintenance.Save') }}</button>
                        <button v-if="id !== 0" @click="remove" class="btn-sm w-32 bg-blue-900 text-blue-100 hover:bg-blue-800 dark:bg-blue-100 dark:text-blue-800 dark:hover:bg-white whitespace-nowrap">{{ t('config.maintenance.Delete') }}</button>
                      </div>
                    </div>
                  </div>

                    <!-- Job list -->
                    <div class="space-y-2">
                      <JobListItem
                          v-for="item in paginatedContents"
                          :key="item.id + '-' + contents.length"
                          :item="item"
                          @editTitle="handleEdit"
                      />             
                    </div>

                    <!-- Real Pagination -->
                    <div v-if="totalPages > 1" class="mt-6">
                        <nav class="flex justify-center" aria-label="Pagination">
                          <div class="flex space-x-1">
                            <!-- Previous button -->
                            <button
                              @click="goToPage(currentPage - 1)"
                              :disabled="currentPage === 1"
                              class="px-3 py-2 text-sm leading-tight text-gray-500 bg-white border border-gray-300 rounded-l-lg hover:bg-gray-100 hover:text-gray-700 disabled:opacity-50 disabled:cursor-not-allowed dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white"
                            >
                              <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                                <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd"></path>
                              </svg>
                            </button>

                            <!-- Page numbers -->
                            <button
                              v-for="page in visiblePages"
                              :key="page"
                              @click="goToPage(page)"
                              :class="[
                                'px-3 py-2 text-sm leading-tight border',
                                currentPage === page
                                  ? 'text-blue-600 border-blue-300 bg-blue-50 hover:bg-blue-100 hover:text-blue-700 dark:border-gray-700 dark:bg-gray-700 dark:text-white'
                                  : 'text-gray-500 bg-white border-gray-300 hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white'
                              ]"
                            >
                              {{ page }}
                            </button>

                            <!-- Next button -->
                            <button
                              @click="goToPage(currentPage + 1)"
                              :disabled="currentPage === totalPages"
                              class="px-3 py-2 text-sm leading-tight text-gray-500 bg-white border border-gray-300 rounded-r-lg hover:bg-gray-100 hover:text-gray-700 disabled:opacity-50 disabled:cursor-not-allowed dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white"
                            >
                              <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                                <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd"></path>
                              </svg>
                            </button>
                          </div>
                        </nav>

                        <!-- Pagination info -->
                        <div class="flex justify-center mt-2">
                          <span class="text-sm text-gray-700 dark:text-gray-300">
                            {{ startIndex + 1 }} - {{ Math.min(endIndex, contents.length) }} / {{ contents.length }} 
                          </span>
                        </div>
                    </div>
                </div>
          </div>
        </section>
        <section>
        </section>
      </div>
  
    </div>
  </template>
  
  <script>
  import { useI18n } from 'vue-i18n'
  import JobListItem from './maintenanceItem.vue'
  //import { useRoute } from 'vue-router'
  import { ref , computed, onMounted, reactive, watchEffect} from 'vue'
  import axios from 'axios'
  export default {
    name: 'MaintenancePanel',
    components:{
      JobListItem,
    },
    setup(){
      const { t } = useI18n();
      const sidebarOpen = ref(false)
      const id = ref(0);
      const message = ref('')
      const title = ref('')
      const context = ref('')
      const contents = ref([]);
      const selContext = ref(0);
      const versions = ref({
        fw:'',
        app:'',
        web:''
      });
      const mtype = ref(-1);
      const selectedTarget=reactive({
        fw:false,
        app:false,
        web:false
      })

      // Pagination state
      const currentPage = ref(1);
      const itemsPerPage = ref(5); // 페이지당 아이템 수

      const utype = computed(() => {
        const types = [];
        if (selectedTarget.fw) types.push("fw");
        if (selectedTarget.app) types.push("app");
        if (selectedTarget.web) types.push("web");
        return types.join(",");
      });

      const isNoneSelected = computed(() => {
        return !selectedTarget.fw &&
              !selectedTarget.app &&
              !selectedTarget.web
      });

      // Pagination computed properties
      const totalPages = computed(() => {
        return Math.ceil(contents.value.length / itemsPerPage.value);
      });

      const startIndex = computed(() => {
        return (currentPage.value - 1) * itemsPerPage.value;
      });

      const endIndex = computed(() => {
        return startIndex.value + itemsPerPage.value;
      });

      const paginatedContents = computed(() => {
        // 최신 순으로 정렬 (id가 클수록 최신이라고 가정)
        const sortedContents = [...contents.value].sort((a, b) => b.id - a.id);
        return sortedContents.slice(startIndex.value, endIndex.value);
      });

      const visiblePages = computed(() => {
        const pages = [];
        const maxVisiblePages = 5;
        let start = Math.max(1, currentPage.value - Math.floor(maxVisiblePages / 2));
        let end = Math.min(totalPages.value, start + maxVisiblePages - 1);
        
        // Adjust start if we're near the end
        if (end - start + 1 < maxVisiblePages) {
          start = Math.max(1, end - maxVisiblePages + 1);
        }
        
        for (let i = start; i <= end; i++) {
          pages.push(i);
        }
        return pages;
      });

      // Pagination methods
      const goToPage = (page) => {
        if (page >= 1 && page <= totalPages.value) {
          currentPage.value = page;
        }
      };

      // Reset to first page when data changes
      watchEffect(() => {
        if (contents.value.length > 0 && currentPage.value > totalPages.value) {
          currentPage.value = 1;
        }
      });
  
  onMounted(()=>{
    getContents();
  });

  function applyUtype(utypeStr) {
  const types = utypeStr.split(',').map(t => t.trim().toLowerCase())
  //console.log(types);
  selectedTarget.fw = types.includes('fw')
  selectedTarget.app = types.includes('app')
  selectedTarget.web = types.includes('web')
}

  const handleEdit = (item) =>{
    selContext.value = 1;
    id.value = item.id;
    title.value = item.title;
    mtype.value = item.mtype;
    applyUtype(item.utype);
    context.value = item.context;
    versions.value.fw = item.f_version;
    versions.value.app = item.a_version;
    versions.value.web = item.w_version;
  }

  const save = async () => {
      const data = {
        title: title.value,
        mtype: mtype.value,
        utype: utype.value,
        context: context.value,
        f_version: versions.value.fw,
        a_version:versions.value.app,
        w_version:versions.value.web,
      };
      try {
        const response = await axios.post(
          `/config/savePost/${selContext.value}/${id.value}`,
          data,
          {
            headers: { "Content-Type": "application/json" },
            withCredentials: true,
          }
        );

        if (response.data && response.data.passOK === "1") {
          alert("✅ Data saved successfully!");
          selContext.value = 0;
          id.value = 0;
          title.value = "";
          mtype.value = -1;
          context.value = "";
          versions.value.fw = "";
          versions.value.app = "";
          versions.value.web = "";
          applyUtype("");
          await getContents();
        } else {
          alert("❌ Data save failed!");
        }
      } catch (err) {
        console.log("Error occurred while saving:", err);
      }
    };

    const getContents = async () => {
      try {
         const response = await axios.get('/config/getPost');
        if (response.data.result === 1) {
          const newData = response.data.data || [];
          contents.value.splice(0, contents.value.length, ...newData);
          // Reset to first page when data is refreshed
          currentPage.value = 1;
        }else{
           contents.value.splice(0);
        }
      } catch (err) {
        contents.value.splice(0);
        //console.log("읽기 오류 발생:", err);
      }
    };

    const remove = async () => {
      try {
         const response = await axios.get(`/config/deletePost/${id.value}`);
        if (response.data.result === 1) {
          alert("✅ Data deleted successfully!");
          id.value = 0;
          selContext.value = 0;
          title.value = "";
          mtype.value = -1;
          context.value = "";
          versions.value.fw = "";
          versions.value.app = "";
          versions.value.web = "";
          applyUtype("");
          // 데이터 다시 불러오기
          await getContents();
          if(contents.value.length == 0){
            contents.value = [];
            console.log('Update List');
          }
        } 
      } catch (err) {
        console.log("읽기 오류 발생:", err);
      }
    };

    return {
      sidebarOpen,
      message,
      versions,
      mtype,
      selectedTarget,
      isNoneSelected,
      title,
      context,
      utype,
      save,
      contents,
      handleEdit,
      id,
      remove,
      // Pagination
      currentPage,
      itemsPerPage,
      totalPages,
      startIndex,
      endIndex,
      paginatedContents,
      visiblePages,
      goToPage,
      t,
    }  

    }
  }
  </script>