<template>
  <div class="col-span-full xl:col-span-5 bg-white dark:bg-gray-800 shadow-sm rounded-xl">
  <div class="card">
    <div class="premium-card-header">
      <div class="header-content">
        <div class="header-left">
          <h2 class="card-title">
            {{t('meter.cardTitle.title_phase')  }}
          </h2>    
        </div>  
      </div>
    </div> 
  </div>
    <!-- <header class="px-5 py-4 border-b border-gray-100 dark:border-gray-700/60">
      <h2 class="font-semibold text-gray-800 dark:text-gray-100">{{ t('meter.cardTitle.title_phase') }}</h2>
    </header> -->
    <div class="p-2 flex justify-center items-center">
      <canvas v-if="channel === 'Main'" ref="mainCanvasRef" width="420" height="420"/>
      <canvas v-else ref="subCanvasRef" width="420" height="420" />
    </div>
    <!--div class="p-2 space-y-2 flex flex-col justify-center items-center">
      <div class="grid grid-cols-4 gap-x-4 gap-y-2 p-3">
        <template v-for="i in 3" :key="i">
          <div class="flex items-start py-2">
            <div class="text-xs font-medium text-violet-700 px-1.5 bg-violet-500/20 
                        dark:bg-violet-600/40 dark:text-violet-300 rounded-full">
              {{ labelDict[i - 1]?.label }}
            </div>
            <div class="text-lg font-bold text-gray-800 dark:text-gray-100 mr-2">
              {{ labelDict[i - 1]?.value?.toFixed(2) }} {{ labelDict[i - 1]?.unit }}
            </div>
          </div>
          <div class="flex items-start py-2">
            <div class="text-xs font-medium text-violet-700 px-1.5 bg-violet-500/20 
                        dark:bg-violet-600/40 dark:text-violet-300 rounded-full">
              {{ angleDict[i - 1]?.label }}
            </div>
            <div class="text-lg font-bold text-gray-800 dark:text-gray-100 mr-2">
              {{ angleDict[i - 1]?.value?.toFixed(2) }} {{ angleDict[i - 1]?.unit }}
            </div>
          </div>
          <div class="flex items-start py-2">
            <div class="text-xs font-medium text-violet-700 px-1.5 bg-violet-500/20 
                        dark:bg-violet-600/40 dark:text-violet-300 rounded-full">
              {{ labelDict[i + 2]?.label }}
            </div>
            <div class="text-lg font-bold text-gray-800 dark:text-gray-100 mr-2">
              {{ labelDict[i + 2]?.value?.toFixed(2) }} {{ labelDict[i + 2]?.unit }}
            </div>
          </div>
          <div class="flex items-start py-2">
            <div class="text-xs font-medium text-violet-700 px-1.5 bg-violet-500/20 
                        dark:bg-violet-600/40 dark:text-violet-300 rounded-full">
              {{ angleDict[i + 2]?.label }}
            </div>
            <div class="text-lg font-bold text-gray-800 dark:text-gray-100 mr-2">
              {{ angleDict[i + 2]?.value?.toFixed(2) }} {{ angleDict[i + 2]?.unit }}
            </div>
          </div>
        </template>
      </div>
    </div-->
  </div>
</template>

<script>
import { onMounted, ref, watch, nextTick } from 'vue';
import { useI18n } from 'vue-i18n'
import { useDark } from '@vueuse/core'
import { tailwindConfig } from '../../../utils/Utils'

export default {
  name: 'CanvasAngle',
  props: {
    degree: Array,      // 각도 (도 단위) [0 ~ 360]
    magnitude: Array,   // 크기 (비례해서 화살표 길이에 반영)
    maxlist: Array,     // 크기 최대값 리스트 (크기 비율용)
    texts: Array,        // 라벨 텍스트 배열
    channel: String,
  },
  setup(props) {
    const mainCanvasRef = ref(null);
    const subCanvasRef = ref(null);
    const channel = ref(props.channel);
    const darkMode = useDark();
    let ctx = null;
    let radius = 0;
    const { t } = useI18n();
    const dcount = 6;
    const baseColors = [
      tailwindConfig().theme.colors.gray[500], 
      tailwindConfig().theme.colors.orange[500], 
      tailwindConfig().theme.colors.sky[500],
      tailwindConfig().theme.colors.slate[800], 
      tailwindConfig().theme.colors.red[600], 
      tailwindConfig().theme.colors.blue[600]
    ];

    const labelDict = ref([]);
    const angleDict = ref([]);
    let isInitialized = false;

    onMounted(async () => {
      await nextTick(); // DOM 렌더링 완료 대기
      await initCanvas();
    });

    // 채널 변경 감지
    watch(() => props.channel, async (newChannel) => {
      channel.value = newChannel;
      await nextTick();
      await initCanvas();
    }, { immediate: false });

    async function initCanvas() {
      // 약간의 지연을 추가하여 DOM이 완전히 준비될 때까지 대기
      await new Promise(resolve => setTimeout(resolve, 10));
      
      const canvas = channel.value === 'Main' ? mainCanvasRef.value : subCanvasRef.value;
      
      if (!canvas) {
        console.warn('Canvas element not found:', channel.value);
        return;
      }

      try {
        ctx = canvas.getContext('2d');
        if (!ctx) {
          console.error('Failed to get canvas context');
          return;
        }

        radius = canvas.height / 2;

        // transform 초기화 및 재적용
        ctx.setTransform(1, 0, 0, 1, 0, 0);
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        ctx.translate(radius, radius);

        radius *= 0.85;
        isInitialized = true;

        //console.log('Canvas initialized successfully:', channel.value);

        // 초기 데이터가 있으면 그리기
        if (props.degree && props.magnitude && props.maxlist && props.texts) {
          updateDataAndDraw();
        }
      } catch (error) {
        console.error('Error initializing canvas:', error);
        isInitialized = false;
      }
    }

    // props 변경 감지
    watch(
      [() => props.degree, () => props.magnitude, () => props.maxlist, () => props.texts],
      () => {
        if (isInitialized) {
          updateDataAndDraw();
        } else {
          // 초기화가 안 된 경우 재시도
          setTimeout(() => {
            if (!isInitialized) {
              initCanvas();
            }
          }, 100);
        }
      },
      { immediate: true, deep: true }
    );

    // 다크모드 변경 감지
    watch(darkMode, () => {
      if (isInitialized) {
        drawCanvas();
      }
    });

    function updateDataAndDraw() {
      const { degree, magnitude, maxlist, texts } = props;
      
      // 데이터 유효성 검사
      if (!Array.isArray(degree) || degree.length < 6 ||
          !Array.isArray(magnitude) || magnitude.length < 6 ||
          !Array.isArray(maxlist) || maxlist.length < 6 ||
          !Array.isArray(texts) || texts.length < 6) {
        console.warn("⚠️ CanvasAngle: props 데이터가 부족합니다", {
          degree: degree?.length,
          magnitude: magnitude?.length,
          maxlist: maxlist?.length,
          texts: texts?.length
        });
        return;
      }

      // 데이터 업데이트
      labelDict.value = [
        { label: "U1", value: magnitude[0], unit: "V" },
        { label: "U2", value: magnitude[1], unit: "V" },
        { label: "U3", value: magnitude[2], unit: "V" },
        { label: "I1", value: magnitude[3], unit: "A" },
        { label: "I2", value: magnitude[4], unit: "A" },
        { label: "I3", value: magnitude[5], unit: "A" }
      ];
      
      angleDict.value = [
        { label: "U1 Angle", value: degree[0], unit: "°" },
        { label: "U2 Angle", value: degree[1], unit: "°" },
        { label: "U3 Angle", value: degree[2], unit: "°" },
        { label: "I1 Angle", value: degree[3], unit: "°" },
        { label: "I2 Angle", value: degree[4], unit: "°" },
        { label: "I3 Angle", value: degree[5], unit: "°" }
      ];

      // console.log('Drawing canvas with data:', {
      //   channel: channel.value,
      //   degree: degree.slice(0, 3),
      //   magnitude: magnitude.slice(0, 3)
      // });

      // 캔버스 그리기
      drawCanvas();
    }

    function drawCanvas() {
      if (!ctx || !isInitialized) {
        console.warn('Canvas not ready for drawing');
        return;
      }
      
      try {
        // 캔버스 클리어
        ctx.clearRect(-radius / 0.9, -radius / 0.9, (radius / 0.9) * 2, (radius / 0.9) * 2);
        
        // 배경 원 그리기
        drawBackground();
        
        // 격자 그리기
        drawGrid();
        
        // 숫자 그리기
        drawNumbers();
        
        // 화살표 그리기
        drawArrows();
        
        // 중심점 그리기
        drawCenter();
        
        //console.log('Canvas drawing completed');
      } catch (error) {
        console.error('Error drawing canvas:', error);
      }
    }

    function drawBackground() {
      ctx.beginPath();
      ctx.arc(0, 0, radius, 0, 2 * Math.PI);
      ctx.fillStyle = darkMode.value ? 
        tailwindConfig().theme.colors.gray[100] : 
        tailwindConfig().theme.colors.gray[200];
      ctx.fill();
      
      ctx.setLineDash([5, 3]);
      ctx.strokeStyle = darkMode.value ? 
        tailwindConfig().theme.colors.gray[200] : 
        tailwindConfig().theme.colors.gray[400];
      ctx.lineWidth = radius * 0.01;
      ctx.stroke();
      ctx.setLineDash([]); // 점선 해제
    }

    function drawGrid() {
      ctx.setLineDash([5, 3]);
      ctx.strokeStyle = darkMode.value ? 
        tailwindConfig().theme.colors.gray[200] : 
        tailwindConfig().theme.colors.gray[400];
      ctx.lineWidth = radius * 0.01;

      // 동심원 그리기
      for (let i = 1; i <= 3; i++) {
        ctx.beginPath();
        ctx.arc(0, 0, radius * i / 3, 0, 2 * Math.PI);
        ctx.stroke();
      }

      // 방사선 그리기
      for (let i = 0; i < 12; i++) {
        const angle = (i * Math.PI / 6);
        drawRadialLine(angle, radius * 0.99, radius * 0.01);
      }
      
      ctx.setLineDash([]); // 점선 해제
    }

    function drawRadialLine(angle, length, width) {
      ctx.save();
      ctx.beginPath();
      ctx.lineWidth = width;
      ctx.lineCap = "square";
      ctx.moveTo(0, 0);
      ctx.rotate(angle);
      ctx.lineTo(0, -length);
      ctx.stroke();
      ctx.restore();
    }

    function drawArrows() {
      const { degree, magnitude, maxlist, texts } = props;
      
      // console.log('Drawing arrows with data:', {
      //   degree: degree?.slice(0, 6),
      //   magnitude: magnitude?.slice(0, 6),
      //   maxlist: maxlist?.slice(0, 6)
      // });
      
      for (let i = 0; i < dcount; i++) {
        if (degree[i] == null || magnitude[i] == null || 
            maxlist[i] == null || texts[i] == null) {
          console.warn(`Missing data for arrow ${i}:`, {
            degree: degree[i],
            magnitude: magnitude[i],
            maxlist: maxlist[i],
            texts: texts[i]
          });
          continue;
        }

        const angleRad = (-degree[i] + 90) * Math.PI / 180;
        const ratio = maxlist[i] === 0 ? 0 : magnitude[i] / maxlist[i];
        const len = i < 3 ? radius * ratio : radius * 0.55 * ratio;

        //console.log(`Arrow ${i}: angle=${degree[i]}, magnitude=${magnitude[i]}, ratio=${ratio}, length=${len}`);

        drawArrow(angleRad, len, radius * 0.03, baseColors[i], texts[i], magnitude[i]);
      }
    }

    function drawArrow(angle, length, width, color, text, value) {
      if (length <= 0) return;
      
      ctx.save();
      ctx.beginPath();
      ctx.lineWidth = width;
      ctx.lineCap = "square";
      ctx.strokeStyle = color;
      ctx.moveTo(0, 0);
      ctx.rotate(angle);
      ctx.lineTo(0, -length);
      
      // 화살표 머리 그리기
      ctx.lineTo(-4, -length + 8);
      ctx.moveTo(0, -length);
      ctx.lineTo(4, -length + 8);
      ctx.stroke();

      // 텍스트 그리기
      ctx.fillStyle = color;
      ctx.font = "bold 12px arial";
      ctx.fillText(text.substring(0, 2), 5, -length - 6);
      
      ctx.restore();
    }

    function drawNumbers() {
      ctx.save();
      ctx.font = radius * 0.06 + "px arial";
      ctx.textBaseline = "middle";
      ctx.textAlign = "center";
      ctx.fillStyle = tailwindConfig().theme.colors.violet[500];

      for (let num = 0; num < 12; num++) {
        let ang = (num + 3) * Math.PI / 6;
        ctx.rotate(ang);
        ctx.translate(0, -radius * 1.08);
        ctx.rotate(-ang);
        
        let tnum = 360 - (num * 30);
        if (tnum === 360) tnum = 0;
        if (tnum % 60 === 0) {
          ctx.fillText(tnum + '°', 0, 0);
        }
        
        ctx.rotate(ang);
        ctx.translate(0, radius * 1.08);
        ctx.rotate(-ang);
      }
      ctx.restore();
    }

    function drawCenter() {
      ctx.beginPath();
      ctx.arc(0, 0, radius * 0.02, 0, 2 * Math.PI);
      ctx.fillStyle = '#333';
      ctx.fill();
      ctx.stroke();
    }

    return {
      mainCanvasRef,
      subCanvasRef,
      t,
      labelDict,
      angleDict,
      channel,
    };
  }
};
</script>

<style>
/* CSS 파일 import */
@import '../../../css/card-styles.css';


</style>