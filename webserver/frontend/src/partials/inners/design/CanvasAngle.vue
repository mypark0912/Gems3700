<template>
  <div class="col-span-full xl:col-span-5 meter-card">
    <div class="meter-card-header">
      <h3 class="meter-card-title meter-accent-indigo">{{ t('meter.cardTitle.title_phase') }}</h3>
    </div>
    <div class="meter-card-body flex justify-center items-center">
      <canvas v-if="channel === 'Main'" ref="mainCanvasRef" width="420" height="420"/>
      <canvas v-else ref="subCanvasRef" width="420" height="420" />
    </div>
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
    degree: Array,
    magnitude: Array,
    maxlist: Array,
    texts: Array,
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
      await nextTick();
      await initCanvas();
    });

    watch(() => props.channel, async (newChannel) => {
      channel.value = newChannel;
      await nextTick();
      await initCanvas();
    }, { immediate: false });

    async function initCanvas() {
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
        ctx.setTransform(1, 0, 0, 1, 0, 0);
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        ctx.translate(radius, radius);
        radius *= 0.85;
        isInitialized = true;

        if (props.degree && props.magnitude && props.maxlist && props.texts) {
          updateDataAndDraw();
        }
      } catch (error) {
        console.error('Error initializing canvas:', error);
        isInitialized = false;
      }
    }

    watch(
      [() => props.degree, () => props.magnitude, () => props.maxlist, () => props.texts],
      () => {
        if (isInitialized) {
          updateDataAndDraw();
        } else {
          setTimeout(() => {
            if (!isInitialized) {
              initCanvas();
            }
          }, 100);
        }
      },
      { immediate: true, deep: true }
    );

    watch(darkMode, () => {
      if (isInitialized) {
        drawCanvas();
      }
    });

    function updateDataAndDraw() {
      const { degree, magnitude, maxlist, texts } = props;

      if (!Array.isArray(degree) || degree.length < 6 ||
          !Array.isArray(magnitude) || magnitude.length < 6 ||
          !Array.isArray(maxlist) || maxlist.length < 6 ||
          !Array.isArray(texts) || texts.length < 6) {
        return;
      }

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

      drawCanvas();
    }

    function drawCanvas() {
      if (!ctx || !isInitialized) return;

      try {
        ctx.clearRect(-radius / 0.9, -radius / 0.9, (radius / 0.9) * 2, (radius / 0.9) * 2);
        drawBackground();
        drawGrid();
        drawNumbers();
        drawArrows();
        drawCenter();
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
      ctx.setLineDash([]);
    }

    function drawGrid() {
      ctx.setLineDash([5, 3]);
      ctx.strokeStyle = darkMode.value ?
        tailwindConfig().theme.colors.gray[200] :
        tailwindConfig().theme.colors.gray[400];
      ctx.lineWidth = radius * 0.01;

      for (let i = 1; i <= 3; i++) {
        ctx.beginPath();
        ctx.arc(0, 0, radius * i / 3, 0, 2 * Math.PI);
        ctx.stroke();
      }

      for (let i = 0; i < 12; i++) {
        const angle = (i * Math.PI / 6);
        drawRadialLine(angle, radius * 0.99, radius * 0.01);
      }

      ctx.setLineDash([]);
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

      for (let i = 0; i < dcount; i++) {
        if (degree[i] == null || magnitude[i] == null ||
            maxlist[i] == null || texts[i] == null) {
          continue;
        }

        const angleRad = (-degree[i] + 90) * Math.PI / 180;
        const ratio = maxlist[i] === 0 ? 0 : magnitude[i] / maxlist[i];
        const len = i < 3 ? radius * ratio : radius * 0.55 * ratio;

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

      ctx.lineTo(-4, -length + 8);
      ctx.moveTo(0, -length);
      ctx.lineTo(4, -length + 8);
      ctx.stroke();

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

<style scoped>
@import '../../../css/meter-card.css';
</style>
