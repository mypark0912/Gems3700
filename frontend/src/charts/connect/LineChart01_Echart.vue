<template>
  <div ref="chartRef" class="w-full h-[700px]"></div>
</template>

<script>
import { onMounted, onBeforeUnmount, ref, watch } from "vue";
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
  },
  setup(props) {
    const { t, locale } = useI18n();
    const chartRef = ref(null);
    let chartInstance = null;

    const darkMode = useDark();
    const { textColor, gridColor } = chartColors;
    const thresholdColorMap = {
      warning: "#ffff00", // 노랑
      repair: "#ff0000", // 빨강
      inspect: "#ff7f00", // 주황
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

    //const thresholdColors = ["#ffff00", "#ff7f00", "#ff0000"];

    const buildSeriesAndLegends = () => {
      let thresholdCounter = 0;
      const normalSeries = [];
      const thresholdSeries = [];
      const normalLegendNames = [];
      const thresholdLegendNames = [];

      console.log(props.chartData);

      props.chartData.forEach((dataset, index) => {
        if (dataset.isThreshold) {
          const color =
            thresholdColorMap[dataset.name.toLowerCase()] ??
            //thresholdColors[thresholdCounter % thresholdColors.length];
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
            color: darkMode.value ? textColor.dark : textColor.light,
          },
        },
        // {
        //   data: thresholdLegendNames,
        //   top: 30,
        //   right: 10,
        //   textStyle: {
        //     color: darkMode.value ? textColor.dark : textColor.light,
        //   },
        // },
      ];

      return { allSeries, legends };
    };

    const getYAxisRange = (seriesList) => {
      const allValues = seriesList.flatMap(s =>
        s.data.filter(v => typeof v === "number" && !isNaN(v))
      );

      if (allValues.length === 0) return [0, 1]; // fallback

      const min = Math.min(...allValues);
      const max = Math.max(...allValues);

      return [Math.floor(Math.min(min, 0)), Math.ceil(max)];
    };


    // const getYAxisRange = (data) => {
    //   let min = Math.min(...data);
    //   if (min < 0) min = 0;
    //   const max = Math.max(...data);
    //   console.log([Math.floor(min), Math.ceil(max)]);
    //   return [Math.floor(min), Math.ceil(max)];
    // };

    const initChart = () => {
      if (chartRef.value) {
        chartInstance = echarts.init(chartRef.value);

        const { allSeries, legends } = buildSeriesAndLegends();

        const [minY, maxY] = getYAxisRange(allSeries);

        // const allValues = allSeries.flatMap((s) =>
        //   s.data.filter((v) => typeof v === "number")
        // );
        // const [minY, maxY] = getYAxisRange(allValues);

        const isEmpty =
          props.chartData.length === 0 ||
          allSeries.every((s) => s.data.length === 0);

        chartInstance.setOption({
          title: {
            text: t("trend.Linechart.TrendChart"),
            left: "center",
            top: 0,
            textStyle: {
              fontSize: 16,
              fontWeight: "bold",
              color: darkMode.value ? textColor.dark : textColor.light,
            },
          },
          tooltip: {
            trigger: "axis",
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

              // ✅ 값 기준으로 내림차순 정렬 (단, number만 비교)
              const sorted = [...params].sort((a, b) => {
                const va = typeof a.data === "number" ? a.data : -Infinity;
                const vb = typeof b.data === "number" ? b.data : -Infinity;
                return vb - va;
              });

              sorted.forEach((item) => {
                const rawValue =
                  typeof item.data === "number" ? item.data : item.value;
                const formatted =
                  typeof rawValue === "number" ? rawValue.toFixed(2) : rawValue;
                result += `${item.marker} ${item.seriesName}: ${formatted}<br/>`;
              });

              return result;
            },
          },
          legend: legends,
          grid: {
            top: 60,
            left: "3%",
            right: "4%",
            bottom: "20%", // 기존보다 여유있게 줘야 slider 표시됨
            containLabel: true,
          },
          dataZoom: [
            {
              type: "inside",
              xAxisIndex: 0,
              filterMode: "filter",
            },
            {
              type: "slider",
              xAxisIndex: 0,
              height: 20, // ✅ 슬라이더 바 자체를 넓게
              bottom: 120, // ✅ 여백도 충분히
              zoomLock: true, // ✅ 줌 잠금
              handleSize: 0, // ✅ 핸들 완전 제거 (줌 못하게)
              showDetail: false, // ✅ 마우스오버 시 값 안 뜨게
              brushSelect: false, // ✅ 클릭/드래그로 범위 선택 못하게
              showDataShadow: false,
              moveOnMouseMove: true, // ✅ 바 본체를 드래그 가능하게
              moveOnMouseWheel: false,
              backgroundColor: "rgba(200,200,255,0.1)",
              fillerColor: "rgba(120, 144, 200, 0.1)", // 선택 영역 색
              borderColor: "transparent",
              handleStyle: {
                color: "#888",
              },
              textStyle: {
                color: darkMode.value ? textColor.dark : textColor.light,
              },
            },
          ],
          xAxis: {
            type: "category",
            data: props.chartLabels,
            axisLine: { lineStyle: { color: "#ccc" } },
            axisLabel: {
              color: darkMode.value ? textColor.dark : textColor.light,
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
            min: minY,
            max: maxY,
            axisLine: { lineStyle: { color: "#ccc" } },
            splitLine: { lineStyle: { color: "#eee" } },
            axisLabel: {
              color: darkMode.value ? textColor.dark : textColor.light,
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
                    fill: "#999",
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

    // ✅ 언어(locale) 변경 시 차트 다시 그리기
    watch(locale, () => {
      disposeChart();
      setTimeout(() => {
        initChart();
      }, 0);
    });

    return {
      chartRef,
      t,
    };
  },
};
</script>

<style scoped></style>
