// composables/useTrendTree.js
// 트렌드 탭에서 사용하는 트리 데이터 및 파라미터 맵 정의

// Meter 탭용 트리 데이터
export const trendTreeData = [
  {
    ID: 1,
    Name: "Temp",
    Title: "Temperature",
    Titles: { en: "Temperature", ko: "온도", ja: "温度" },
    isParent: true,
  },
  {
    ID: 2,
    Name: "Freq",
    Title: "Frequency",
    Titles: { en: "Frequency", ko: "주파수", ja: "周波数" },
    isParent: true
  },
  {
    ID: 3,
    Name: "PF",
    Title: "Power Factor",
    Titles: { en: "Power Factor", ko: "역률", ja: "力率" },
    isParent: true,
    children: [
      { ID: 4, Name: "PF1", Title: "Power Factor L1", Titles: { en: "Power Factor L1", ko: "역률 L1", ja: "力率 L1" } },
      { ID: 5, Name: "PF2", Title: "Power Factor L2", Titles: { en: "Power Factor L2", ko: "역률 L2", ja: "力率 L2" } },
      { ID: 6, Name: "PF3", Title: "Power Factor L3", Titles: { en: "Power Factor L3", ko: "역률 L3", ja: "力率 L3" } },
    ],
  },
  {
    ID: 7,
    Name: "linevoltage",
    Title: "Line Voltage",
    Titles: { en: "Line Voltage", ko: "선간전압", ja: "線間電圧" },
    isParent: true,
    children: [
      { ID: 8, Name: "Upp1", Title: "Line Voltage L12", Titles: { en: "Line Voltage L12", ko: "선간전압 L12", ja: "線間電圧 L12" } },
      { ID: 9, Name: "Upp2", Title: "Line Voltage L23", Titles: { en: "Line Voltage L23", ko: "선간전압 L23", ja: "線間電圧 L23" } },
      { ID: 10, Name: "Upp3", Title: "Line Voltage L31", Titles: { en: "Line Voltage L31", ko: "선간전압 L31", ja: "線間電圧 L31" } },
      { ID: 11, Name: "Upp4", Title: "Line Voltage Average", Titles: { en: "Line Voltage Average", ko: "선간전압 평균", ja: "線間電圧 平均" } },
    ],
  },
  {
    ID: 12,
    Name: "voltage",
    Title: "Phase Voltage",
    Titles: { en: "Phase Voltage", ko: "상전압", ja: "相電圧" },
    isParent: true,
    children: [
      { ID: 13, Name: "U1", Title: "Voltage L1", Titles: { en: "Voltage L1", ko: "전압 L1", ja: "電圧 L1" } },
      { ID: 14, Name: "U2", Title: "Voltage L2", Titles: { en: "Voltage L2", ko: "전압 L2", ja: "電圧 L2" } },
      { ID: 15, Name: "U3", Title: "Voltage L3", Titles: { en: "Voltage L3", ko: "전압 L3", ja: "電圧 L3" } },
    ],
  },
  {
    ID: 16,
    Name: "current",
    Title: "Current",
    Titles: { en: "Current", ko: "전류", ja: "電流" },
    isParent: true,
    children: [
      { ID: 17, Name: "I1", Title: "Current L1", Titles: { en: "Current L1", ko: "전류 L1", ja: "電流 L1" } },
      { ID: 18, Name: "I2", Title: "Current L2", Titles: { en: "Current L2", ko: "전류 L2", ja: "電流 L2" } },
      { ID: 19, Name: "I3", Title: "Current L3", Titles: { en: "Current L3", ko: "전류 L3", ja: "電流 L3" } },
    ],
  },
  {
    ID: 20,
    Name: "unbalance",
    Title: "Unbalance",
    Titles: { en: "Unbalance", ko: "불평형", ja: "不平衡" },
    isParent: true,
    children: [
      { ID: 21, Name: "Ibal1", Title: "Current Unbalance", Titles: { en: "Current Unbalance", ko: "전류 불평형", ja: "電流不平衡" } },
      { ID: 22, Name: "Ubal1", Title: "Voltage Unbalance", Titles: { en: "Voltage Unbalance", ko: "전압 불평형", ja: "電圧不平衡" } },
    ],
  },
  {
    ID: 23,
    Name: "Power",
    Title: "Power",
    Titles: { en: "Power", ko: "전력", ja: "電力" },
    isParent: true,
    children: [
      { ID: 24, Name: "P4", Title: "Active Power", Titles: { en: "Active Power", ko: "유효전력", ja: "有効電力" } },
      { ID: 25, Name: "Q4", Title: "Reactive Power", Titles: { en: "Reactive Power", ko: "무효전력", ja: "無効電力" } },
      { ID: 26, Name: "S4", Title: "Apparent Power", Titles: { en: "Apparent Power", ko: "피상전력", ja: "皮相電力" } },
    ],
  },
  {
    ID: 31,
    Name: "THD",
    Title: "THD",
    Titles: { en: "THD", ko: "THD", ja: "THD" },
    isParent: true,
    children: [
      { ID: 32, Name: "THD_U1", Title: "THD Voltage L1", Titles: { en: "THD Voltage L1", ko: "THD 전압 L1", ja: "THD 電圧 L1" } },
      { ID: 33, Name: "THD_U2", Title: "THD Voltage L2", Titles: { en: "THD Voltage L2", ko: "THD 전압 L2", ja: "THD 電圧 L2" } },
      { ID: 34, Name: "THD_U3", Title: "THD Voltage L3", Titles: { en: "THD Voltage L3", ko: "THD 전압 L3", ja: "THD 電圧 L3" } },
      { ID: 35, Name: "THD_Upp1", Title: "THD Voltage L12", Titles: { en: "THD Voltage L12", ko: "THD 전압 L12", ja: "THD 電圧 L12" } },
      { ID: 36, Name: "THD_Upp2", Title: "THD Voltage L23", Titles: { en: "THD Voltage L23", ko: "THD 전압 L23", ja: "THD 電圧 L23" } },
      { ID: 37, Name: "THD_Upp3", Title: "THD Voltage L31", Titles: { en: "THD Voltage L31", ko: "THD 전압 L31", ja: "THD 電圧 L31" } },
      { ID: 38, Name: "THD_I1", Title: "THD Current L1", Titles: { en: "THD Current L1", ko: "THD 전류 L1", ja: "THD 電流 L1" } },
      { ID: 39, Name: "THD_I2", Title: "THD Current L2", Titles: { en: "THD Current L2", ko: "THD 전류 L2", ja: "THD 電流 L2" } },
      { ID: 40, Name: "THD_I3", Title: "THD Current L3", Titles: { en: "THD Current L3", ko: "THD 전류 L3", ja: "THD 電流 L3" } },
    ],
  },
  {
    ID: 41,
    Name: "TDD",
    Title: "TDD",
    Titles: { en: "TDD", ko: "TDD", ja: "TDD" },
    isParent: true,
    children: [
      { ID: 42, Name: "TDD_I1", Title: "TDD Current L1", Titles: { en: "TDD Current L1", ko: "TDD 전류 L1", ja: "TDD 電流 L1" } },
      { ID: 43, Name: "TDD_I2", Title: "TDD Current L2", Titles: { en: "TDD Current L2", ko: "TDD 전류 L2", ja: "TDD 電流 L2" } },
      { ID: 44, Name: "TDD_I3", Title: "TDD Current L3", Titles: { en: "TDD Current L3", ko: "TDD 전류 L3", ja: "TDD 電流 L3" } },
    ],
  },
];

// Energy 탭 전용 트리 데이터
export const energyTreeData = [
  {
    ID: 27,
    Name: "Energy",
    Title: "Energy",
    Titles: { en: "Energy", ko: "에너지", ja: "エネルギー" },
    isParent: true,
    children: [
      { ID: 28, Name: "kwh_import_consumption", Title: "Active Energy", Titles: { en: "Active Energy", ko: "유효 에너지", ja: "有効エネルギー" } },
      { ID: 29, Name: "kvarh_import_consumption", Title: "Reactive Energy", Titles: { en: "Reactive Energy", ko: "무효 에너지", ja: "無効エネルギー" } },
      { ID: 30, Name: "kvah_import_consumption", Title: "Apparent Energy", Titles: { en: "Apparent Energy", ko: "피상 에너지", ja: "皮相エネルギー" } },
    ],
  },
];

// Demand 탭 전용 트리 데이터
export const demandTreeData = [
  {
    ID: 50,
    Name: "Demand",
    Title: "Demand",
    Titles: { en: "Demand", ko: "수요전력", ja: "デマンド" },
    isParent: true,
    children: [
      { ID: 51, Name: "CD_P_import", Title: "Active Power Demand", Titles: { en: "Active Power Demand", ko: "유효전력 디맨드", ja: "有効電力デマンド" } },
      { ID: 52, Name: "CD_Q_import", Title: "Reactive Power Demand", Titles: { en: "Reactive Power Demand", ko: "무효전력 디맨드", ja: "無効電力デマンド" } },
      { ID: 53, Name: "CD_S", Title: "Apparent Power Demand", Titles: { en: "Apparent Power Demand", ko: "피상전력 디맨드", ja: "皮相電力デマンド" } },
    ],
  },
];

// Meter 탭: 설정에서 활성화된 파라미터를 트리 Name으로 매핑
export const meterParamMap = {
  Temperature: ["Temp"],
  Frequency: ["Freq"],
  "Phase Voltage": ["voltage"],
  "Line Voltage": ["linevoltage"],
  "Line Voltage L12": ["Upp1"],
  "Line Voltage L23": ["Upp2"],
  "Line Voltage L31": ["Upp3"],
  "Line Voltage Average": ["Upp4"],
  Current: ["current"],
  "Current L1": ["I1"],
  "Current L2": ["I2"],
  "Current L3": ["I3"],
  Power: ["Power"],
  "Active Power": ["P4"],
  "Reactive Power": ["Q4"],
  "Apparent Power": ["S4"],
  PF: ["PF"],
  Unbalance: ["unbalance"],
  "Current Unbalance": ["Ibal1"],
  "Voltage Unbalance": ["Ubal1"],
  THD: ["THD"],
  "THD Voltage L1": ["THD_U1"],
  "THD Voltage L2": ["THD_U2"],
  "THD Voltage L3": ["THD_U3"],
  "THD Voltage L12": ["THD_Upp1"],
  "THD Voltage L23": ["THD_Upp2"],
  "THD Voltage L31": ["THD_Upp3"],
  "THD Current L1": ["THD_I1"],
  "THD Current L2": ["THD_I2"],
  "THD Current L3": ["THD_I3"],
  TDD: ["TDD"],
  "TDD I1": ["TDD_I1"],
  "TDD I2": ["TDD_I2"],
  "TDD I3": ["TDD_I3"],
  "TDD Iavg": ["tdd_i_avg"],
};

// Meter 차트: 선택된 Name을 API 필드로 매핑
export const meterChartParamMap = {
  Temp: ["Temp"],
  U1: ["U1"],
  U2: ["U2"],
  U3: ["U3"],
  U4: ["U4"],
  Upp1: ["Upp1"],
  Upp2: ["Upp2"],
  Upp3: ["Upp3"],
  Upp4: ["Upp4"],
  I1: ["I1"],
  I2: ["I2"],
  I3: ["I3"],
  P4: ["P4"],
  Q4: ["Q4"],
  S4: ["S4"],
  Freq: ["Freq"],
  PF1: ["PF1"],
  PF2: ["PF2"],
  PF3: ["PF3"],
  vunbal: ["Ubal1"],
  curunbal: ["Ibal1"],
  THD_U1: ["THD_U1"],
  THD_U2: ["THD_U2"],
  THD_U3: ["THD_U3"],
  THD_Upp1: ["THD_Upp1"],
  THD_Upp2: ["THD_Upp2"],
  THD_Upp3: ["THD_Upp3"],
  THD_I1: ["THD_I1"],
  THD_I2: ["THD_I2"],
  THD_I3: ["THD_I3"],
  TDD_I1: ["TDD_I1"],
  TDD_I2: ["TDD_I2"],
  TDD_I3: ["TDD_I3"],
  Ubal1: ["Ubal1"],
  Ibal1: ["Ibal1"],
};

// Demand 파라미터 맵
export const demandParamMap = {
  CD_P_import: ["CD_P_import"],
  CD_Q_import: ["CD_Q_import"],
  CD_S: ["CD_S"],
};

// 트리 데이터로부터 Name → Title 맵 생성
export function buildNameToTitleMap(...trees) {
  const map = {};
  const build = (nodes) => {
    nodes.forEach((node) => {
      map[node.Name] = node.Title;
      if (node.children) {
        build(node.children);
      }
    });
  };
  trees.forEach((tree) => build(tree));
  return map;
}
