<template>
    <div class="flex gap-4">
      <!-- Channel 1 -->
      <div class="flex-1 flex flex-col gap-4">
        <FormBlock :forms="Forms" title="Main Channel" />
        <GroupTable :channels="channels.main" tableTitle="Main Channel" />
      </div>
  
      <!-- Channel 2 -->
      <div class="flex-1 flex flex-col gap-4">
        <FormBlock :forms="Forms" title="Sub Channel" />
        <GroupTable :channels="channels.sub" tableTitle="Sub Channel" />
      </div>
    </div>
  </template>
  
  <script setup>

  import { reactive, ref } from 'vue'
  import GroupTable from './GroupTable.vue'
  import FormBlock from './FormBlock.vue'
    
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
    
  </script>
  
  <style scoped>
  table, th, td {
    border: 1px solid gray;
  }
  </style>