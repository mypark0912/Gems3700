<template>
    <div class="grow">
  
      <!-- Panel body -->
      <div class="p-4">
        <div class="mt-2 grid grid-cols-2 gap-4 border-b border-gray-200 dark:border-gray-700/60 p-2 pb-4">
            <div>
              <h2 class="text-2xl text-gray-800 dark:text-gray-100 font-bold mb-5"> Main Channel</h2>

              <div class="flex items-center gap-2">
                <label for="selectItem" class="text-xs font-bold flex items-center text-gray-700 dark:text-gray-300">
                Select Setup
              </label>

              <select
                id="selectItem"
                class="h-8 w-32 p-2 text-xs border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400"
                v-model="selectSetup_main"
              >
              <option v-for="(setup, index) in setupList_main" :key="index" :value="index">
                {{ setup.item }}
              </option>
              </select>
              </div>
            </div>
            <div>
              <h2 class="text-2xl text-gray-800 dark:text-gray-100 font-bold mb-5"> Sub Channel</h2>
              <div class="flex items-center gap-2">
                <label for="selectItem" class="text-xs font-bold flex items-center text-gray-700 dark:text-gray-300">
                Select Setup
              </label>
              <select
                id="selectItem"
                class="h-8 w-32 p-2 text-xs border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400"
                v-model="selectSetup_sub"
              >
              <option v-for="(setup, index) in setupList_sub" :key="index" :value="index">
                {{ setup.item }}
              </option>
              </select>
              </div>
            </div>            
        </div>


  
        <!-- Connected Apps cards -->
        <section class="mt-2 pb-2 border-b border-gray-200 dark:border-gray-700/60">
            <div class="grid grid-cols-2 gap-4">
                <GroupTable :channels="channels.main" />
                <GroupTable :channels="channels.sub" />
            </div>
          
        </section>
        <!--section class="mt-2 pb-6 border-b border-gray-200 dark:border-gray-700/60">
          <CaliCard v-for="item in items" :item="item"/>
        </section-->
        <section>
            <div class="mt-4 mb-1 flex gap-2">
                <button
                class="h-9 w-32 px-4 text-sm bg-gray-900 text-white rounded-md hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white transition-colors duration-200"
            >Save</button>
            <button
                class="h-9 w-32 px-4 text-sm bg-blue-900 text-white rounded-md hover:bg-blue-800 dark:bg-blue-100 dark:text-blue-800 dark:hover:bg-white transition-colors duration-200"
            >Send Report</button>
            </div>
        </section>
      </div>
  
    </div>
  </template>
  
  <script>
  import CaliCard from './CaliCard2.vue';
  import CaliCard_New from './CaliCard4.vue';
  import FormBlock from './FormBlock.vue';
  import GroupTable from './GroupTable.vue';
  import { useRoute } from 'vue-router'
  import { ref , watch, reactive, onMounted} from 'vue'
  import { useSetupStore } from "@/store/setup"; // ✅ Pinia Store 사용
  import axios from 'axios';
  export default {
    name: 'CaliPanel',
    props:['items','channel'],
    components:{
        CaliCard,
        CaliCard_New,
        FormBlock,
        GroupTable,
    },
    setup(props){
        const setupStore = useSetupStore();
        const route = useRoute();
        const items = ref(props.items);

        const channel = ref(props.channel);
        const selectSetup_main = ref(0);
        const selectSetup_sub = ref(0);
        const tvalue = ref('');
        const setupList_main = ref([]);
        const setupList_sub = ref([]);
        const selectItem = ref('Time');
        const channels = reactive({
      "main":{       
        "Phase Voltage": {
            "view":[
              { subTitle: "U_A", value: 0, error: 0 },
              { subTitle: "U_B", value: 0, error: 0 },
              { subTitle: "U_C", value: 0, error: 0 },
              { subTitle: "Upp", value: 0, error: 0 }
            ],
          },
          "Phase Current":{
              "view": [
              { subTitle: "I_A", value: 0, error: 0 },
              { subTitle: "I_B", value: 0, error: 0 },
              { subTitle: "I_C", value: 0, error: 0 },
              { subTitle: "In", value: 0, error: 0 }
            ],
          },
        "Power Angle": {
            "view":[
              { subTitle: "Angle_A", value: 0, error: 0 },
              { subTitle: "Angle_B", value: 0, error: 0 },
              { subTitle: "Angle_C", value: 0, error: 0 }
            ],
          },
          "Active Power":{
            "view" :[
              { subTitle: "Watt_A", value: 0, error: 0 },
              { subTitle: "Watt_B", value: 0, error: 0 },
              { subTitle: "Watt_C", value: 0, error: 0 }
            ]
          },
          "Reactive Power":{
            "view":[
              { subTitle: "Var_A", value: 0, error: 0 },
              { subTitle: "Var_B", value: 0, error: 0 },
              { subTitle: "Var_C", value: 0, error: 0 }
            ]
          },
          "Apparent Power":{
            "view":[
              { subTitle: "VA_A", value: 0, error: 0 },
              { subTitle: "VA_B", value: 0, error: 0 },
              { subTitle: "VA_C", value: 0, error: 0 }
            ]
          },
      },
      "sub":{       
        "Phase Voltage": {
            "view":[
              { subTitle: "U_A", value: 0, error: 0 },
              { subTitle: "U_B", value: 0, error: 0 },
              { subTitle: "U_C", value: 0, error: 0 },
              { subTitle: "Upp", value: 0, error: 0 }
            ],
          },
          "Phase Current":{
              "view": [
              { subTitle: "I_A", value: 0, error: 0 },
              { subTitle: "I_B", value: 0, error: 0 },
              { subTitle: "I_C", value: 0, error: 0 },
              { subTitle: "In", value: 0, error: 0 }
            ],
          },
        "Power Angle": {
            "view":[
              { subTitle: "Angle_A", value: 0, error: 0 },
              { subTitle: "Angle_B", value: 0, error: 0 },
              { subTitle: "Angle_C", value: 0, error: 0 }
            ],
          },
          "Active Power":{
            "view" :[
              { subTitle: "Watt_A", value: 0, error: 0 },
              { subTitle: "Watt_B", value: 0, error: 0 },
              { subTitle: "Watt_C", value: 0, error: 0 }
            ]
          },
          "Reactive Power":{
            "view":[
              { subTitle: "Var_A", value: 0, error: 0 },
              { subTitle: "Var_B", value: 0, error: 0 },
              { subTitle: "Var_C", value: 0, error: 0 }
            ]
          },
          "Apparent Power":{
            "view":[
              { subTitle: "VA_A", value: 0, error: 0 },
              { subTitle: "VA_B", value: 0, error: 0 },
              { subTitle: "VA_C", value: 0, error: 0 }
            ]
          },
      },
    })
        const Forms = ref({
            "row1":{
                "lbl":"U Ref",
                "text":"",
                "btn0":"Send",
                "btn1":"Set U",
                "btn2":"Clear U",
            },
            "row2":{
                "lbl":"Upp Ref",
                "text":"",
                "btn0":"Send",
                "btn1":"Set Upp",
                "btn2":"Clear Upp"
            },
            "row3":{
                "lbl":"I Ref",
                "text":"",
                "btn0":"Send",
                "btn1":"Set I",
                "btn2":"Clear I",
            },
            "row4":{
                "lbl":"In Ref",
                "text":"",
                "btn0":"Send",
                "btn1":"Set In",
                "btn2":"Clear In"
            },
            "row5":{
              "lbl":"Phase Angle",
              "text":"",
              "btn0":"Send",
              "btn1":"Set Phase",
              "btn2":"Clear Phase",
            },
            "row6":{
              "btn1":"Set Power",
              "btn2":"Clear Power"
            }
        })
    
        const fetchsetupList = async() =>{
          try{
            const response = await axios.get("/config/checkSetup");
            if(response.data.passOK == "1"){
              setupList_main.value = response.data.data;
              setupList_sub.value = response.data.data;
              if(!setupStore.getCalib)
                setupStore.setCalib(true);
              //console.log(setupList_main.value);
            }
          } catch (error) {
            const msg = "업로드 실패: " + error.response.data.error;
            alert(msg);
          }
        }

    onMounted(()=>{
      fetchsetupList();
    })

    watch(() => props.channel, (newChannel) => {
      if (newChannel !== channel.value) {
        channel.value = newChannel;
      }
    }, { immediate: true });

  // route.params.channel 변경 감시
    watch(() => route.params.channel, (newChannel) => {
      if (newChannel !== channel.value) {
        channel.value = newChannel;
      }
    }, { immediate: true });

    watch(() => setupStore.calib, (newVal) => {
      if (newVal) fetchsetupList(newVal);
    }, { immediate: true });
    

      return {
          items,
          channel,
          selectSetup_main,
          selectSetup_sub,
          tvalue,
          selectItem,
          Forms,
          channels,
          setupList_main,
          setupList_sub,
      }
    }
  }
  </script>

<style scoped>
/* 다크모드에서 select option 스타일링 추가 */
@media (prefers-color-scheme: dark) {
  select option {
    background-color: rgb(31 41 55); /* gray-800 */
    color: rgb(243 244 246); /* gray-100 */
  }
  
  select option:hover,
  select option:focus,
  select option:checked {
    background-color: rgb(55 65 81); /* gray-700 */
  }
}

/* 라이트모드에서 option 스타일링 */
@media (prefers-color-scheme: light) {
  select option {
    background-color: white;
    color: rgb(17 24 39); /* gray-900 */
  }
  
  select option:hover,
  select option:focus,
  select option:checked {
    background-color: rgb(243 244 246); /* gray-100 */
  }
}
</style>