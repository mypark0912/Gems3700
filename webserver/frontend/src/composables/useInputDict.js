// composables/useInputDict.js
import { ref, computed, reactive } from 'vue';

export function useInputDict() {
  const setupDict = ref({
    "mode":'',
    "lang":'',
    "General":{},
    "main":{},
    "ibsm": { channel: "ibsm", Enable: 1, tapboxs: [
      { CANid: "310941", index: 0, mtype: 10, cbtype: 5, cbcount: 2, stype: [1,1,1,1,1,1], main: 0, name: "" }
    ] },
    "ipsm72":{
      channel: "ipsm72",
      Enable: 1,
      enable1: 0,
      enable2: 0,
      di60_1: 0,
      di60_2: 0,
      di60Target1: 0,
      di60Target2: 0,
      ipsm72_1: Array.from({ length: 72 }, () => ({ mod: null, pt: null })),
      ipsm72_2: Array.from({ length: 72 }, () => ({ mod: null, pt: null })),
    },
    "mcs": {
      channel: "mcs",
      Enable: 1,
      feeders: Array.from({ length: 32 }, () => ({
        serialNumber: null,
        touName: '',
        systemType: 0,
        modbusId: 0,
        ct1: 1,
        ct2: 1,
        ct3: 1,
      })),
    }
  });

  const inputDict = ref({
    useFunction: {
      sntp: 0,
      modbus_serial: 0,
      mqtt:0
    },
    deviceInfo: {
      name: "",
      location:"",
      serial_number: "",
      mac_address: "",
      timezone: "",
    },
    tcpip: {
      ip_address: "",
      subnet_mask: "",
      gateway: "",
      dnsserver: "",
      dhcp: 0,
    },
    modbus: {
      tcp_port: 502,
      baud_rate: 0,
      parity: 0,
      data_bits:7,
      stop_bits:1,
    },
    etc:{
      backlight_off_time: 0,
      autorotation: 0,
      brightness: 50,
      testmode: 0,
      update_interval: 1,
    },
    sntpInfo: {
      host: "",
    },
    MQTT :{
      Type:"",
      host:"",
      port:"",
      device_id:"",
      username:"",
      password:"",
      externalport:"",
      url:""
    }
  });

  const channel_main = ref({
    channel: "",
    Enable: true,
    ctInfo: {
      direction: [0, 0, 0],
      startingcurrent: 1,
      inorminal: 100,
      ct1: 100,
      ct2: 1,
      zctscale: 1,
      zcttpye: 1,
    },
    opt: {
      pf_sign: 0,
      unbalance: 0,
      va_type: 0,
      iload: 0.0,
    },
    demand: {
      target: "",
      demand_interval: 15,
    },

    ptInfo: {
      wiringmode: 0,
      linefrequency: 60,
      vnorminal: 100,
      pt1: 100,
      pt2: 100,
      dash :0
    },
    eventInfo: {
      sag_action: 0,
      sag_holdofftime: 1,
      sag_level: 1,
      sag_hysteresis: 0,
      swell_action: 0,
      swell_holdofftime: 1,
      swell_level: 1,
      swell_hysteresis: 0,
      inter_action: 0,
      inter_holdofftime: 1,
      inter_level: 1,
      inter_hysteresis: 0,
      rvc_action: 0,
      rvc_holdofftime: 0,
      rvc_level: 0,
      rvc_hysteresis: 0,
      msv_level: 0,
      msv_carrierFrequency: 0,
      msv_recordLength: 0,
      flicker_model: 0,
    },
    trendInfo: {
      period: 5,
      params: [1, 2, 3, 4, 5, 6, 7, 8, 9, 0],
    },
    alarm: Object.fromEntries(
      Array.from({ length: 32 }, (_, i) => [i + 1, { chan: 0, cond: 0, action: 0, delay: 0, hysteresis: 0, threshold: 1 }])
    ),
  });

  const channel_sub = ref({
    channel: "",
    Enable: false,
    ctInfo: {
      direction: [0, 0, 0],
      startingcurrent: 1,
      inorminal: 100,
      ct1: 100,
      ct2: 1,
      zctscale: 1,
      zcttpye: 1,
    },
    opt: {
      pf_sign: 0,
      unbalance: 0,
      va_type: 0,
      iload: 0.0,
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
    eventInfo: {
      sag_action: 0,
      sag_holdofftime: 1,
      sag_level: 1,
      sag_hysteresis: 0,
      swell_action: 0,
      swell_holdofftime: 1,
      swell_level: 1,
      swell_hysteresis: 0,
      inter_action: 0,
      inter_holdofftime: 1,
      inter_level: 1,
      inter_hysteresis: 0,
      rvc_action: 0,
      rvc_holdofftime: 0,
      rvc_level: 0,
      rvc_hysteresis: 0,
      msv_level: 0,
      msv_carrierFrequency: 0,
      msv_recordLength: 0,
      flicker_model: 0,
    },
    trendInfo: {
      period: 5,
      params: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    },
    alarm: Object.fromEntries(
      Array.from({ length: 32 }, (_, i) => [i + 1, { chan: 0, cond: 0, action: 0, delay: 0, hysteresis: 0, threshold: 1 }])
    ),
  });


  const selectedTrendSetup = ref({
    period: 5,
    params: [1, 2, 3, 4, 5, 6, 7, 8, 9, 0],
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
    setupDict,
    formatToISOString,
    selectedTrendSetup,
  };
}
