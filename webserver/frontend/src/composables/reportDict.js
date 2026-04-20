// composables/useReportData.js
import { reactive } from 'vue'
import axios from 'axios'

const reportData = reactive({
  infoData: {},
  tbdata:{},
  diagnosisData: {
    barchart:{},
    chStatus:{},
    chChart:{}
  },
  powerQualityData: {
    barchart:{},
    chStatus:{},
    chChart : {}
  },
  energyHourlyData: {},
  energyDailyData: {},
  energyMonthlyData: {},
  loadrateData:{},
  weeklyloadData:{},
  heatmaploadData:{},
  isLoaded: {
    info: false,
    tbdata:false,
    diagnosis: false,
    powerQuality: false,
    energy: false,
    loadrate: false,
    weekload:false,
    heatmapload:false,
  }
})

const baseChart = {
  labels: [],  // ✅ 빈 labels 추가 (없으면 에러날 수 있음)
  datasets: [
    {
      label: 'Swell Upper Limit',
      data: [
        { x: 0.0001, y: 500 },
        { x: 0.001, y: 200 },
        { x: 0.003, y: 140 },
        { x: 0.003, y: 120 },
        { x: 0.02, y: 120 },
        { x: 0.5, y: 120 },
        { x: 0.5, y: 110 },
        { x: 10, y: 110 },
        { x: 100, y: 110 },
      ],
      borderColor: 'blue',
      backgroundColor: 'transparent',
      borderWidth: 2,
      pointRadius: 0,
      tension: 0,
    },
    {
      label: 'Sag Lower Limit',
      data: [
        { x: 0.02, y: 0.5 },   // ✅ 0 → 0.5 (로그스케일에서 0 불가)
        { x: 0.02, y: 70 },
        { x: 0.5, y: 70 },
        { x: 0.5, y: 80 },
        { x: 10, y: 80 },
        { x: 10, y: 90 },
        { x: 100, y: 90 },
      ],
      borderColor: 'red',
      backgroundColor: 'transparent',
      borderWidth: 2,
      pointRadius: 0,
      tension: 0,
    }
  ]
}

export function useReportData() {
  // ✅ API 호출 함수들을 composable에 포함
    const loadInfoData = async (chName) => {
        try {
            //const ch = 'Fan';
            const response = await axios.get(`/api/getAsset/${chName}`);
            if (response.data.success) {
              reportData.infoData = response.data.data;
              reportData.info = true;
              reportData.mac = response.data.mac;
              reportData.driveType = response.data.driveType;
              return reportData;
            }else{
              console.log('No Data');
            }
        }catch (error) {
            console.log("데이터 가져오기 실패:", error);
        }
    };

    const saveReportData = async (chName, chType, location) => {
      try {
          const data={
            "assetName":chName,
            "assetType":chType,
            "location":location
          }
          const response = await axios.post(`/report/generate`,  data, {
            headers: { "Content-Type": "application/json" },
          });
          console.log(response.data);
      }catch (error) {
          console.log("데이터 가져오기 실패:", error);
      }
  };

 const formatToISOString = (dateStr, endOfDay) => {
    const date = new Date(dateStr)
    if (endOfDay) {
      date.setHours(23, 59, 59, 999)
    } else {
      date.setHours(0, 0, 0, 0)
    }
    return date.toISOString()
  }

  const setParamIds = async (asset, datalist) => {
    const idList = []
    try {
      const response = await axios.get(`/api/getParameters/${asset}/diagnostic`)
      if (response.data.success) {
        const paramData = response.data.data
        for (let i = 0; i < paramData.length; i++) {
          for (let j = 0; j < datalist.length; j++) {
            if (paramData[i]["Name"] == datalist[j].Name) {
              idList.push({ idx: paramData[i]["ID"], title: paramData[i]["Title"] })
            }
          }
        }
      }
    } catch (error) {
      console.error(error)
    }
    return idList
  }

  const setPQParamIds = async(asset, datalist) =>{
      const idList = [];
      try{
        const response = await axios.get(`/api/getParameters/${asset}/powerquality`);
        if (response.data.success) {
          const paramData = response.data.data;
          //console.log('datalist',datalist);
          for(let i = 0 ; i < paramData.length ; i++){
            for(let j = 0; j < datalist.length; j++){
              if(paramData[i]["Name"] == datalist[j].Name){
                idList.push({idx:paramData[i]["ID"], title:paramData[i]["Title"]});
              }
            }
          }
        }
      }catch (error){
        console.error(error);
      }
      return idList;
    }

  const setChartData = async (effectiveIds, Title) => {
    let option = {}

    const trendDataRequest = {
      ParametersIds: [effectiveIds],
      StartDate: formatToISOString('2025-06-01', 0),
      EndDate: formatToISOString('2025-06-07', 1),
    }

    try {
      const url = `/api/getTrendData`
      const response = await axios.post(url, trendDataRequest, {
        headers: { "Content-Type": "application/json" },
      })

      if (response.data.success) {
        const resData = response.data.data
        let datasets = []
        let labels = []

        // 1. 메인 데이터 처리 (Thresholds 제외)
        Object.keys(resData).forEach((key) => {
          if (key !== "Thresholds") {
            const dataPoints = resData[key].data
            if (dataPoints && dataPoints.length > 0) {
              if (labels.length === 0) {
                labels = dataPoints.map((point) => point.XAxis)
              }
              datasets.push({
                name: resData[key].Title,
                data: dataPoints.map((point) => point.YAxis),
                isThreshold: false,
              })
            }
          }
        })

        // 2. Thresholds 처리
        // if (resData.Thresholds && resData.Thresholds.length == 2 && labels.length > 0) {
        //   let timeList = []
        //   for (let i = 0; i < resData.Thresholds.length; i++) {
        //     timeList.push(new Date(resData.Thresholds[i].XAxis))
        //   }
        //   const t1 = new Date(resData.Thresholds[0].XAxis)
        //   const t2 = new Date(resData.Thresholds[1].XAxis)

        //   resData.Thresholds[0].Thresholds.forEach((value, idx) => {
        //     if (value !== "NaN" && value !== null && value !== undefined) {
        //       const secondValue = resData.Thresholds[1].Thresholds[idx]
        //       if (secondValue === "NaN" || secondValue === null || secondValue === undefined) {
        //         return
        //       }

        //       const thresholdData = labels.map((lbl) => {
        //         const dt = new Date(lbl)
        //         return dt < t1 ? value : secondValue
        //       })
              
        //       const ThresholdString = [
        //         "Out of Range(Down side)",
        //         "Repair",
        //         "Inspect", 
        //         "Warning",
        //         "Warning",
        //         "Inspect",
        //         "Repair",
        //         "Out of Range(Upper side)",
        //       ]
              
        //       datasets.push({
        //         name: ThresholdString[idx],
        //         data: thresholdData,
        //         isThreshold: true,
        //       })
        //     }
        //   })
        // }
        if (resData.Thresholds && resData.Thresholds.length > 0 && labels.length > 0) {
          // Threshold 문자열
            const ThresholdString = [
              "Out of Range(Down side)",
              "Repair",
              "Inspect",
              "Warning",
              "Warning",
              "Inspect",
              "Repair",
              "Out of Range(Upper side)",
            ];

            // Threshold 데이터가 있는 인덱스 찾기
            if(resData.Thresholds[0].Thresholds != null){
              const thresholdCount = resData.Thresholds[0].Thresholds.length;
            
            for (let idx = 0; idx < thresholdCount; idx++) {
              // 해당 인덱스에 유효한 값이 하나라도 있는지 확인
              const hasValidValue = resData.Thresholds.some(t => {
                const value = t.Thresholds[idx];
                return value !== "NaN" && value !== null && value !== undefined && typeof value === 'number';
              });

              if (!hasValidValue) continue;

              // Threshold 시간 리스트 생성 (유효한 값만)
              const timeList = resData.Thresholds
                .filter(t => {
                  const value = t.Thresholds[idx];
                  return value !== "NaN" && value !== null && value !== undefined && typeof value === 'number';
                })
                .map(t => ({
                  time: new Date(t.XAxis),
                  value: t.Thresholds[idx]
                }))
                .sort((a, b) => a.time - b.time);

              if (timeList.length === 0) continue;

              // 각 label에 대해 적절한 threshold 값 찾기
              const thresholdData = labels.map((lbl) => {
                const labelTime = new Date(lbl);
                
                // labelTime에 해당하는 threshold 값 찾기
                // labelTime보다 작거나 같은 가장 최근 threshold 사용
                let applicableThreshold = timeList[0].value;
                
                for (let i = 0; i < timeList.length; i++) {
                  if (labelTime >= timeList[i].time) {
                    applicableThreshold = timeList[i].value;
                  } else {
                    break;
                  }
                }
                
                return applicableThreshold;
              });

              datasets.push({
                name: ThresholdString[idx],
                data: thresholdData,
                isThreshold: true,
              });
            }
            }
            
          }

        // 3. 차트 업데이트
        option = {
          lineLabels: labels,
          lineData: datasets,
          lineTitle: Title,
        }
      } else {
        console.error("서버 오류:", response.data.error)
        option = {
          lineLabels: [],
          lineData: [],
          lineTitle: '',
        }
      }
    } catch (error) {
      console.error("요청 실패:", error)
      option = {
        lineLabels: [],
        lineData: [],
        lineTitle: '',
      }
    }
    return option
  }

  // ✅ PDF 전용 - 단순히 데이터만 반환
  const loadDiagnosisData = async (chName, locale) => {
    try {
      const response = await axios.get(`/api/getDiagnosisDetail/${chName}`)

      if (response.data.success) {
        let itemlist = [], valuelist = [], datalist = []
        
        for (let i = 0; i < response.data.data_status.length; i++) {
          itemlist.push(response.data.data_status[i]["Titles"][locale])
          valuelist.push(response.data.data_status[i]["Status"])
          datalist.push(response.data.data_status[i]["Titles"])
        }
        
        const chartdata = { "Names": itemlist, "Values": valuelist, "Titles": datalist }
        
        const treeItems = response.data.data_tree
        const chartList = []
        const items = []
        
        for (let i = 0; i < treeItems.length; i++) {
          if (treeItems[i]["Status"] > 1) {
            const childDict = []
            for (let j = 0; j < treeItems[i]["children"].length; j++) {
              if (treeItems[i]["children"][j]["Status"] > 1) {
                chartList.push(treeItems[i]["children"][j])
                childDict.push({
                  Title: treeItems[i]["children"][j]["Title"],
                  Value: treeItems[i]["children"][j]["Value"]
                })
              }
            }
            items.push({ Item: treeItems[i], Child: childDict })
          }
        }

        const chartOption = []
        if (chartList.length > 0) {
          let chartValue = null
          const effectiveIds = await setParamIds(chName, chartList)
          for (let i = 0; i < effectiveIds.length; i++) {
            chartValue = await setChartData(effectiveIds[i].idx, effectiveIds[i].title)
            chartOption.push(chartValue)
          }
        }

        // ✅ PDF에 필요한 모든 데이터 반환
        return {
          items,
          chartOption,
          chartdata
        }
      } else {
        //console.log('No Data')
        return {
          items: [],
          chartOption: [],
          chartdata: {}
        }
      }
    } catch (error) {
      //console.log("데이터 가져오기 실패:", error)
      return {
        items: [],
        chartOption: [],
        chartdata: {}
      }
    }
  }

  const loadPQData = async (chName, locale) => {
      try {

        const response = await axios.get(`/api/getDiagPQ/${chName}`);

        const items = []
        if (response.data.success) {
          let itemlist = [], valuelist=[], datalist=[];
          for(let i = 0; i < response.data.data_status.length;i++){
       
            itemlist.push(response.data.data_status[i]["Titles"][locale]);
            valuelist.push(response.data.data_status[i]["Status"]);
            datalist.push(response.data.data_status[i]["Titles"])
          }
          const chartdata = {"Names" : itemlist, "Values" : valuelist, "Titles": datalist};
          const treeItems = response.data.data_tree;
          const chartList = [];
          for(let i = 0 ; i< treeItems.length;i++){
            if(treeItems[i]["Status"] > 1 ){
              //items.value.push();
              const childDict = [];
              if(treeItems[i]["children"].length > 0){
                for(let j = 0 ; j < treeItems[i]["children"].length ; j++ ){
                  if(treeItems[i]["children"][j]["isSub"]){
                    for(let k = 0 ; k < treeItems[i]["children"][j]["children"].length ; k++ ){
                      if( treeItems[i]["children"][j]["children"][k]["Status"] > 1){
                        chartList.push(treeItems[i]["children"][j]["children"][k]);
                        childDict.push({Title:treeItems[i]["children"][j]["children"][k]["Title"], Value:treeItems[i]["children"][j]["children"][k]["Value"]})
                      }
                    }
                  }else{
                    if(treeItems[i]["children"][j]["Status"] > 1){
                      chartList.push(treeItems[i]["children"][j]);
                      childDict.push({Title:treeItems[i]["children"][j]["Title"], Value:treeItems[i]["children"][j]["Value"]})
                    }
                  }
                }
              }else{
                childDict = [{Title:treeItems[i]["Title"], Value:treeItems[i]["Value"]} ]
              }
              items.push({Item:treeItems[i], Child:childDict});            
            }
          }
          const chartOption = []
          if(chartList.length > 0){
            let chartValue=null;
            const effectiveIds = await setPQParamIds(chName, chartList)
            for (let i = 0 ; i < effectiveIds.length;i++){
                chartValue = await setChartData(effectiveIds[i].idx, effectiveIds[i].title);
                chartOption.push(chartValue);
            }
          }
           return {
            items,
            chartOption,
            chartdata
            }

        }else{
          console.log('No Data');
           return {
                items: [],
                chartOption: [],
                chartdata: {}
            }
        }
      }catch (error) {
        console.log("데이터 가져오기 실패:", error);
         return {
            items: [],
            chartOption: [],
            chartdata: {}
        }
      }
    };

  const loadPowerQualityData = async (channel) => {
    try {
      const response = await axios.get(`/api/getEn50160/${channel}`)
      reportData.powerQualityData = response.data.data
      reportData.isLoaded.powerQuality = true
      return reportData.powerQualityData
    } catch (error) {
      console.error('PowerQuality 데이터 로딩 실패:', error)
    }
  }

  const loadEnergyHourlyData = async (channel) => {
    try {
      //const response = await axios.get(`/api/getHourlyEnergyData/${channel}`)
      const response = await axios.get(`/api/getHourlyEnergyData/${channel}`, {
        params: {
          date: new Date().toISOString().split('T')[0],
        }
      });
      if (response.data.success && response.data.hourlyData) {
        reportData.energyHourlyData = response.data.hourlyData;
        reportData.isLoaded.energy = true
        //console.log("Received hourly data:", reportData.energyHourlyData);
      }
      // reportData.energyData = response.data.data
      return reportData.energyHourlyData
    } catch (error) {
      console.error('Energy 데이터 로딩 실패:', error)
    }
  }

  const loadEnergyDailyData = async (channel) => {
    try {
      //const response = await axios.get(`/api/getHourlyEnergyData/${channel}`)
      const response = await axios.get(`/api/getDailyEnergyData/${channel}`, {
        params: {
          date: new Date().toISOString().split('T')[0],
        }
      });
      if (response.data.success && response.data.dailyData) {
        reportData.energyDailyData = response.data.dailyData;
        reportData.isLoaded.energy = true
        //console.log("Received daily data:", reportData.energyDailyData);
      }
      // reportData.energyData = response.data.data
      return reportData.energyDailyData
    } catch (error) {
      console.error('Energy 데이터 로딩 실패:', error)
    }
  }

  const loadEnergyMonthlyData = async (channel) => {
    try {
      //const response = await axios.get(`/api/getHourlyEnergyData/${channel}`)
      const response = await axios.get(`/api/getMonthlyEnergyData/${channel}`, {
        params: {
          date: new Date().toISOString().split('T')[0],
        }
      });
      if (response.data.success && response.data.monthlyData) {
        reportData.energyMonthlyData = response.data.monthlyData;
        reportData.isLoaded.energy = true
        //console.log("Received monthly data:", reportData.energyMonthlyData);
      }
      // reportData.energyData = response.data.data
      return reportData.monthlyData
    } catch (error) {
      console.error('Energy 데이터 로딩 실패:', error)
    }
  }

  const getLoadFactorCalculated = async (channel) => {
    try {
      //const response = await axios.get(`/api/getHourlyEnergyData/${channel}`)
      const response = await axios.get(`/api/getLoadFactorCalculated/${channel}`, {
        params: {
          date: new Date().toISOString().split('T')[0],
        }
      });
      if (response.data.success && response.data.loadFactorData) {
        reportData.loadrateData = response.data.loadFactorData;
        reportData.isLoaded.loadrate = true
        //console.log("Received LoadFactor data:", reportData.loadrateData);
      }
      // reportData.energyData = response.data.data
      return reportData.loadrateData
    } catch (error) {
      console.error('Energy 데이터 로딩 실패:', error)
    }
  }

  const getWeeklyLoadFactorData = async (channel) => {
    try {
      //const response = await axios.get(`/api/getHourlyEnergyData/${channel}`)
      const response = await axios.get(`/api/getWeeklyLoadFactorData/${channel}`, {
        params: {
          date: new Date().toISOString().split('T')[0],
        }
      });
      if (response.data.success && response.data.weeklyData) {
        reportData.weeklyloadData = response.data.weeklyData;
        reportData.isLoaded.weekload = true
        //console.log("Received LoadFactor data:", reportData.weeklyloadData);
      }
      // reportData.energyData = response.data.data
      return reportData.weeklyloadData
    } catch (error) {
      console.error('Energy 데이터 로딩 실패:', error)
    }
  }

  const getHeatmapLoadFactorData = async (channel, weeks = 4) => {
    try {
      const response = await axios.get(`/api/getHeatmapLoadFactorData/${channel}`, {
        params: {
          weeks: weeks // weeks 파라미터 사용
        }
      });
      
      if (response.data.success && response.data.heatmapData) {
        // 전체 응답 데이터를 저장 (heatmapData + distribution + overallStats 등)
        reportData.heatmapLoadFactorData = response.data;
        reportData.isLoaded.heatmap = true;
        //console.log("Received heatmap LoadFactor data:", response.data);
        
        // 전체 응답 반환 (Vue 컴포넌트에서 success, distribution 등 사용 가능)
        return response.data;
      } else {
        //console.error('히트맵 데이터 로드 실패');
        return {
          success: false,
          heatmapData: [],
          distribution: { light: 0, medium: 0, high: 0, overload: 0 }
        };
      }
    } catch (error) {
      //console.error('히트맵 데이터 로딩 실패:', error);
      return {
        success: false,
        heatmapData: [],
        distribution: { light: 0, medium: 0, high: 0, overload: 0 }
      };
    }
  }

  // ✅ 모든 데이터를 한번에 로드
  const loadAllData = async (channel) => {
    await Promise.all([
      loadInfoData(channel),
      loadDiagnosisData(channel),
      loadPowerQualityData(channel),
      loadEnergyData(channel)
    ])
  }

  const getAllData = () => ({
    channel,
    datalist: reportData.infoData,
    diagnosisData: reportData.diagnosisData,
    powerQualityData: reportData.powerQualityData,
    energyData: reportData.energyData
  })

  const parseMask = (mask) => {
    const phases = [];
    if (mask & 0b001) phases.push("L1");
    if (mask & 0b010) phases.push("L2");
    if (mask & 0b100) phases.push("L3");
    return phases;
}

const getfinValue = (phaselist, levellist, mode) =>{
  let finalValue = 0;
  if (phaselist.length === 1) {
      // phaselist가 하나일 때 해당하는 Level 값 선택
      switch (phaselist[0]) {
        case "L1":
          finalValue = levellist[0];
          break;
        case "L2":
          finalValue = levellist[1];
          break;
        case "L3":
          finalValue = levellist[2];
          break;
        default:
          finalValue = null; // 예상 외 값일 경우 처리
      }
    } else if (phaselist.length > 1) {
      // 2개 이상일 때 가장 큰 값 선택
      if (mode == 'SWELL')
        finalValue = Math.max(...levellist);
      else
      finalValue = Math.min(...levellist);
    } else {
      finalValue = null; // phaselist가 비어있는 경우
    }
    return finalValue;
}
const makeKey = (param, phase) => {
  const suffixMap = {
    L1: "L1",
    L2: "L2",
    L3: "L3",
    "Multi Phase": "Multi Phase"
  }

  // 예외 처리
  if (param === "Frequency Variation 1") {
    return phase === 'L1' ? "Frequency Variation 1" : undefined
  }
  if (param === "Frequency Variation 2") {
    return phase === 'L1' ? "Frequency Variation 2" : undefined
  }
  if (param === "Voltage Unbalance") {
    return phase === 'L1' ? "Voltage Unbalance" : undefined
  }
  if (param === "Voltage Variation 1") return `Voltage Variation 1 ${suffixMap[phase]}(%)`
  if (param === "Voltage Variation 2") return `Voltage Variation 2 ${suffixMap[phase]}(%)`
  if (param === "THD") return `THDs Variation ${suffixMap[phase]}(%)`
  if (param === "Harmonics") return `Harmonics Variatiopn ${suffixMap[phase]}(%)`
  if (param === "Pst") return `Flickers Pst ${suffixMap[phase]}(%)`
  if (param === "Plt") return `Flickers Plt ${suffixMap[phase]}(%)`
  if (param === "Signal Vol.") return `Signaling Voltage ${suffixMap[phase]}(%)`
  if (param === "Voltage Sag") return `Voltage Dips ${suffixMap[phase]}`
  if (param === "Voltage Swell") return `Voltage Swells ${suffixMap[phase]}`
  if (param === "Short Interruption") return `Short Interruptions ${suffixMap[phase]}`
  if (param === "Long Interruption") return `Long Interruptions ${suffixMap[phase]}`
  if (param === "Signaling Volt.") return `Signaling Voltage ${suffixMap[phase]}(%)`

  // 기본값
  return `${param} ${suffixMap[phase]}`
}

  return {
    reportData,
    baseChart,
    loadInfoData,
    loadDiagnosisData,
    loadPowerQualityData,
    loadEnergyHourlyData,
    loadEnergyDailyData,
    loadEnergyMonthlyData,
    getLoadFactorCalculated,
    getWeeklyLoadFactorData,
    getHeatmapLoadFactorData,
    loadAllData,
    getAllData,
    loadPQData,
    parseMask,
    getfinValue,
    makeKey,
    saveReportData,
  }
}