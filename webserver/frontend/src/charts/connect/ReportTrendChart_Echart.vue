<template>
  <div ref="chartRef" class="w-full h-[400px]"></div>
</template>

<script>
import { onMounted, onBeforeUnmount, ref, watch, inject, computed } from "vue";
import * as echarts from "echarts";
import { useDark } from "@vueuse/core";
import { chartColors } from "../ChartjsConfig";
import { useI18n } from "vue-i18n";

export default {
  name: "LineChart",
  props: {
    chartData: {
      type: Array,
      required: true,
    },
    chartLabels: {
      type: Array,
      required: true,
    },
    title: String,
    mode: {
      type: String,
      default: 'DiagnosisDetail'
    }
  },
  setup(props) {
    const isPdfMode = inject('isPdfMode', false)
    
    const { t, locale } = useI18n();
    const chartRef = ref(null);
    let chartInstance = null;

    const darkMode = useDark();
    const { textColor, gridColor } = chartColors;
    
    const getTextColor = () => {
      if (isPdfMode) return '#000000'
      return darkMode.value ? textColor.dark : textColor.light
    }
    
    // 상태 라벨 computed로 생성 (언어 변경 시 자동 업데이트)
    const statusLabels = computed(() => {
      if (props.mode === 'PowerQuality') {
        return [
          t('diagnosis.tabContext.pqfe0'),
          t('diagnosis.tabContext.pqfe1'),
          t('diagnosis.tabContext.pqfe2'),
          t('diagnosis.tabContext.pqfe3'),
          t('diagnosis.tabContext.pqfe4'),
        ];
      } else {
        return [
          t('diagnosis.tabContext.st0'),
          t('diagnosis.tabContext.st1'),
          t('diagnosis.tabContext.st2'),
          t('diagnosis.tabContext.st3'),
          t('diagnosis.tabContext.st4'),
        ];
      }
    });
    
    // 상태 텍스트 반환
    const getStatusText = (value) => {
      const statusValue = Math.floor(value);
      const labels = statusLabels.value;
      return labels[statusValue] || labels[0];
    };
    
    const thresholdColorMap = {
      warning: "#ffff00",
      repair: "#ff0000",
      inspect: "#ff7f00",
    };
    
    const getLineColor = (index) => {
      const colors = [
        "#1f77b4",
        "#ff7f0e",
        "#2ca02c",
        "#d62728",
        "#9467bd",
        "#8c564b",
        "#e377c2",
        "#7f7f7f",
        "#bcbd22",
        "#17becf",
      ];
      return colors[index % colors.length];
    };

    const buildSeriesAndLegends = () => {
      let thresholdCounter = 0;
      const normalSeries = [];
      const thresholdSeries = [];
      const normalLegendNames = [];
      const thresholdLegendNames = [];

      props.chartData.forEach((dataset, index) => {
        if (dataset.isThreshold) {
          const color =
            thresholdColorMap[dataset.name.toLowerCase()] ??
          thresholdCounter++;
          thresholdSeries.push({
            name: dataset.name,
            type: "line",
            data: dataset.data,
            smooth: true,
            symbol: "none",
            lineStyle: { width: 1, type: "dashed" },
            itemStyle: { color },
            label: {
              show: true,
              position: "end",
              formatter: function (params) {
                const rawDate = new Date(params.name);
                const yyyy = rawDate.getFullYear();
                const MM = String(rawDate.getMonth() + 1).padStart(2, "0");
                const dd = String(rawDate.getDate()).padStart(2, "0");
                const HH = String(rawDate.getHours()).padStart(2, "0");
                const mm = String(rawDate.getMinutes()).padStart(2, "0");
                const ss = String(rawDate.getSeconds()).padStart(2, "0");
                return `${yyyy}-${MM}-${dd} ${HH}:${mm}:${ss}`;
              },
              color,
            },
          });
          thresholdLegendNames.push(dataset.name);
        } else {
          normalSeries.push({
            name: dataset.name,
            type: "line",
            data: dataset.data,
            smooth: true,
            lineStyle: { width: 2 },
            itemStyle: { color: getLineColor(index) },
          });
          normalLegendNames.push(dataset.name);
        }
      });

      const allSeries = [...normalSeries, ...thresholdSeries];

      const legends = [
        {
          data: normalLegendNames,
          bottom: 15,
          textStyle: {
            color: getTextColor(),
          },
        },
      ];

      return { allSeries, legends };
    };

    const initChart = () => {
      if (chartRef.value) {
        chartInstance = echarts.init(chartRef.value);

        const { allSeries, legends } = buildSeriesAndLegends();

        const isEmpty =
          props.chartData.length === 0 ||
          allSeries.every((s) => s.data.length === 0);

        // 현재 상태 라벨 캡처 (클로저용)
        const currentLabels = statusLabels.value;

        chartInstance.setOption({
          title: {
            text: props.title,
            left: "center",
            top: 0,
            textStyle: {
              fontSize: 16,
              fontWeight: "bold",
              color: getTextColor(),
            },
          },
          tooltip: {
            trigger: "axis",
            backgroundColor: isPdfMode ? 'rgba(255, 255, 255, 0.95)' : darkMode.value 
                  ? '#000000'
                  : 'rgba(255, 255, 255, 0.95)',
            textStyle: {
              color: isPdfMode 
                ? '#000000' 
                : darkMode.value 
                  ? undefined
                  : '#000000',
            },
            formatter: function (params) {
              const rawDate = new Date(params[0].axisValue);
              const yyyy = rawDate.getFullYear();
              const MM = String(rawDate.getMonth() + 1).padStart(2, "0");
              const dd = String(rawDate.getDate()).padStart(2, "0");
              const HH = String(rawDate.getHours()).padStart(2, "0");
              const mm = String(rawDate.getMinutes()).padStart(2, "0");
              const ss = String(rawDate.getSeconds()).padStart(2, "0");
              const formattedTime = `${yyyy}-${MM}-${dd} ${HH}:${mm}:${ss}`;
              let result = formattedTime + "<br/>";

              const sorted = [...params].sort((a, b) => {
                const va = typeof a.data === "number" ? a.data : -Infinity;
                const vb = typeof b.data === "number" ? b.data : -Infinity;
                return vb - va;
              });

              sorted.forEach((item) => {
                const rawValue =
                  typeof item.data === "number" ? item.data : item.value;
                // Y축 값을 상태 텍스트로 변환 (캡처된 라벨 사용)
                const statusValue = Math.floor(rawValue);
                const formatted = typeof rawValue === "number" 
                  ? (currentLabels[statusValue] || currentLabels[0])
                  : rawValue;
                result += `${item.marker} ${item.seriesName}: ${formatted}<br/>`;
              });

              return result;
            },
          },
          legend: {
             show: false,
          },
          grid: {
            top: 60,
            left: "3%",
            right: "4%",
            bottom: 30,
            containLabel: true,
          },
          xAxis: {
            type: "category",
            data: props.chartLabels,
            axisLine: { 
              lineStyle: { 
                color: isPdfMode ? '#cccccc' : '#ccc'
              } 
            },
            axisLabel: {
              color: getTextColor(),
              formatter: function (value) {
                const date = new Date(value);
                const yyyy = date.getFullYear();
                const MM = String(date.getMonth() + 1).padStart(2, "0");
                const dd = String(date.getDate()).padStart(2, "0");
                const HH = String(date.getHours()).padStart(2, "0");
                const mm = String(date.getMinutes()).padStart(2, "0");
                const ss = String(date.getSeconds()).padStart(2, "0");
                return `${yyyy}-${MM}-${dd} ${HH}:${mm}:${ss}`;
              },
            },
          },
          yAxis: {
            type: "value",
            min: 0,
            max: 4,
            interval: 1,
            axisLine: { 
              lineStyle: { 
                color: isPdfMode ? '#cccccc' : '#ccc'
              } 
            },
            splitLine: { 
              lineStyle: { 
                color: isPdfMode ? '#e5e5e5' : '#eee'
              } 
            },
            axisLabel: {
              color: getTextColor(),
              formatter: function (value) {
                // Y축 라벨을 상태 텍스트로 변환 (캡처된 라벨 사용)
                const statusValue = Math.floor(value);
                return currentLabels[statusValue] || currentLabels[0];
              },
            },
          },
          series: allSeries,
          graphic: isEmpty
            ? [
                {
                  type: "text",
                  left: "center",
                  top: "middle",
                  style: {
                    text: t("trend.Linechart.nodata"),
                    fontSize: 20,
                    fill: isPdfMode ? '#000000' : '#999',
                  },
                },
              ]
            : [],
        });
      }
    };

    const disposeChart = () => {
      if (chartInstance) {
        chartInstance.dispose();
        chartInstance = null;
      }
    };

    const resizeChart = () => {
      if (chartInstance) {
        chartInstance.resize();
      }
    };

    onMounted(() => {
      initChart();
      window.addEventListener("resize", resizeChart);
    });

    onBeforeUnmount(() => {
      disposeChart();
      window.removeEventListener("resize", resizeChart);
    });

    watch(
      () => [props.chartData, props.chartLabels],
      () => {
        disposeChart();
        setTimeout(() => {
          initChart();
        }, 0);
      },
      { deep: true }
    );

    // mode 변경 시에도 차트 다시 그리기
    watch(
      () => props.mode,
      () => {
        disposeChart();
        setTimeout(() => {
          initChart();
        }, 0);
      }
    );

    watch(locale, () => {
      disposeChart();
      setTimeout(() => {
        initChart();
      }, 0);
    });

    watch(
      () => [darkMode.value, isPdfMode],
      () => {
        if (!isPdfMode) {
          disposeChart();
          setTimeout(() => {
            initChart();
          }, 0);
        }
      }
    );

    return {
      chartRef,
      t,
    };
  },
};
</script>

<style scoped></style>