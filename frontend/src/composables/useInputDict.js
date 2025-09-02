// composables/useInputDict.js
import { ref, computed } from 'vue';

export function useInputDict() {
  const setupDict = ref({
    "mode":'',
    "General":{},
    "main":{},
    "sub":{}
  });
  const inputDict = ref({
    useFuction: {
      ftp: false,
      sntp: false,
      diagnosis_main: false,
      diagnosis_sub: false,
    },
    deviceInfo: {
      name: "",
      // device_model: "SV-500",
      location:"",
      serial_number: "",
      mac_address: "",
    },
    tcpip: {
      ip_address: "",
      subnet_mask: "",
      gateway: "",
      dnsserver: "",
    },
    modbus: {
      isserial: false,
      tcp_port: "",
      modbus_id: "",
      baud_rate: 0,
      parity: 0,
      data_bits:7,
      stop_bits:1
    },
    
    pf_sign: 0,
    va_type: 0,
    unbalance:0,
    ftpInfo: {
      isManual: false,
      host: "",
      port: "",
      id: "",
      pass: "",
      upload_main: "",
      upload_sub: "",
    },
    sntpInfo: {
      // date: "",
      // time: "",
      host: "",
      // period: "",
      timezone: "",
    },
  });

  const diagnosis_detail = ref({
    main:{
      use:false,
      assetName:'',
      tableData:[],
      paramData:[]
    },
    sub:{
      use:false,
      assetName:'',
      tableData:[],
      paramData:[]
    }
  })

  const channel_main = ref({
    channel: "",
    Enable: false,
    PowerQuality: false,
    ctInfo: {
      direction: [0, 0, 0],
      startingcurrent: 1,
      inorminal: 100,
      ct1: 100,
      ct2: 1,
      zctscale: 1,
      zcttpye: 1,
    },
    demand: {
      target: "",
      demand_interval: "15",
    },
    
    ptInfo: {
      wiringmode: 0,
      linefrequency: 60,
      vnorminal: 100,
      pt1: 100,
      pt2: 100,
    },
    assetInfo: {
      name: "",
      type: "",
      nickname:"",
    },
    eventInfo: {
      oc_action: 0,
      oc_level: 1,
      oc_holdofftime: 1,
      tv_action: 0,
      tv_level: 1,
      tv_holdofftime: 1,
      tv_fastchange: 1,
      tc_action: 0,
      tc_level: 1,
      tc_holdofftime: 1,
      inter_action: 0,
      inter_level: 1,
      inter_holdofftime: 1,
      sag_action: 0,
      sag_level: 1,
      sag_holdofftime: 1,
      swell_action: 0,
      swell_level: 1,
      swell_holdofftime: 1,
      tc_fastchange: 1,
    },
    sampling: {
      rate: 8000,
      duration: 10,
      period: 5,
    },
    n_kva:0,
    trendInfo: [],
    alarm: {
      CompareTimeDelay: 0,
      ...Object.fromEntries(
        Array.from({ length: 32 }, (_, i) => [i + 1, [0, 0, 0, 0]])
      ),
    },
  });

  const channel_sub = ref({
    channel: "",
    Enable: false,
    PowerQuality: false,
    ctInfo: {
      direction: [0, 0, 0],
      startingcurrent: 1,
      inorminal: 100,
      ct1: 100,
      ct2: 1,
      zctscale: 1,
      zcttpye: 1,
    },
    demand: {
      target: 0,
      demand_interval: 15,
    },
    
    ptInfo: {
      wiringmode: 0,
      linefrequency: 60,
      vnorminal: 100,
      pt1: 100,
      pt2: 100,
    },
    assetInfo: {
      name: "",
      type: "",
      nickname:"",
    },
    eventInfo: {
      oc_action: 0,
      oc_level: 1,
      oc_holdofftime: 1,
      tv_action: 0,
      tv_level: 1,
      tv_holdofftime: 1,
      tv_fastchange: 1,
      tc_action: 0,
      tc_level: 1,
      tc_holdofftime: 1,
      inter_action: 0,
      inter_level: 1,
      inter_holdofftime: 1,
      sag_action: 0,
      sag_level: 1,
      sag_holdofftime: 1,
      swell_action: 0,
      swell_level: 1,
      swell_holdofftime: 1,
      tc_fastchange: 1,
    },
    sampling: {
      rate: 8000,
      duration: 5,
      period: 5,
    },
    n_kva:0,
    trendInfo: [],
    alarm: {
      CompareTimeDelay: 1,
      ...Object.fromEntries(
        Array.from({ length: 32 }, (_, i) => [i + 1, [0, 0, 0, 0]])
      ),
    },
  });

  const useDiagnosis = computed(() => {
    return (
      inputDict.value.useFuction.diagnosis_main ||
      inputDict.value.useFuction.diagnosis_sub
    );
  });

  
 const currentDiagnosis = computed(() => ({
    Main: inputDict.value.useFuction.diagnosis_main,
    Sub: inputDict.value.useFuction.diagnosis_sub,
  }));

  const diagnosisData = ref({});
  const advancedData = ref({});

  //const Diagnosis_inputDict = ref({});

  const selectedTrendSetup = ref({
    period: 5, // 기본값을 5분으로 설정
    params: ["None","None","None","None","None","None","None","None"], // 첫 번째는 Energy로 고정
  });

  const changeDiagnosis = ref({
    asset : false,
    nameplate : false,
    parameter : false
  });


  const parameterOptions = [
      "None",
      "Temperature",
      "Frequency",
      "Phase Voltage L1",
      "Phase Voltage L2",
      "Phase Voltage L3",
      "Phase Voltage Average",
      "Line Voltage L12",
      "Line Voltage L23",
      "Line Voltage L31",
      "Line Voltage Average",
      "Voltage Unbalance(Uo)",
      "Voltage Unbalance(Uu)",
      "Phase Current L1",
      "Phase Current L2",
      "Phase Current L3",
      "Phase Current Average",
      "Phase Current Total",
      "Phase Current Neutral",
      "Active Power L1",
      "Active Power L2",
      "Active Power L3",
      "Active Power Total",
      "Reactive Power L1",
      "Reactive Power L2",
      "Reactive Power L3",
      "Reactive Power Total",
      "Distortion Power L1",
      "Distortion Power L2",
      "Distortion Power L3",
      "Distortion Power Total",
      "Apparent Power L1",
      "Apparent Power L2",
      "Apparent Power L3",
      "Apparent Power Total",
      "Power Factor L1",
      "Power Factor L2",
      "Power Factor L3",
      "Power Factor Total",
      "THD Voltage L1",
      "THD Voltage L2",
      "THD Voltage L3",
      "THD Voltage L12",
      "THD Voltage L23",
      "THD Voltage L31",
      "THD Current L1",
      "THD Current L2",
      "THD Current L3"
    ];

    const formatToISOString = (date, soe) => {
      if (typeof date === "string") {
        date = new Date(date);
      }
      if (!(date instanceof Date) || isNaN(date)) {
        throw new Error("Invalid date");
      }

      const pad = (num, size = 2) => String(num).padStart(size, "0");

      const year = date.getFullYear();
      const month = pad(date.getMonth() + 1); // 월은 0부터 시작
      const day = pad(date.getDate());
      let hours, minutes, seconds, milliseconds;

      if (soe === 0) {
        hours = pad(0);
        minutes = pad(0);
        seconds = pad(0);
        milliseconds = pad(1, 7);
      } else if (soe === 1) {
        hours = pad(23);
        minutes = pad(59);
        seconds = pad(59);
        milliseconds = pad(999, 7);
      } else if (soe === 2) {
        hours = pad(0);
        minutes = pad(0);
        seconds = pad(0);
        milliseconds = pad(1, 2);
      } else {
        hours = pad(23);
        minutes = pad(59);
        seconds = pad(59);
        milliseconds = pad(99, 2);
      }
      // 타임존 오프셋 계산
      const timezoneOffset = -date.getTimezoneOffset();
      const offsetSign = timezoneOffset >= 0 ? "+" : "-";
      const offsetHours = pad(Math.abs(Math.floor(timezoneOffset / 60)));
      const offsetMinutes = pad(Math.abs(timezoneOffset % 60));
      if (soe === 2 || soe === 3) {
        return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}.${milliseconds}Z`;
      } else {
        return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}.${milliseconds}${offsetSign}${offsetHours}:${offsetMinutes}`;
      }
    };

  return {
    inputDict,
    channel_main,
    channel_sub,
    parameterOptions,
    useDiagnosis,
    diagnosisData,
    advancedData,
    currentDiagnosis,
    setupDict,
    formatToISOString,
    selectedTrendSetup,
    changeDiagnosis,
    diagnosis_detail,
  };
}
