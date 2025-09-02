<template>
    <div class="col-span-full xl:col-span-6 bg-white dark:bg-gray-800 shadow-sm rounded-xl">
      <header class="px-5 py-4 border-b border-gray-100 dark:border-gray-700/60">
        <h2 class="font-semibold text-gray-800 dark:text-gray-100">Voltage/Current Phase</h2>
      </header>
      <div class="p-3 flex justify-center items-center">
        <canvas ref="canvasRef" width="400" height="400" />
      </div>
    </div>
  </template>
  
  <script>
  import { onMounted, ref, watch } from 'vue';
  
  export default {
    name: 'CanvasAngle',
    props: {
      degree: Array,     // 예: [0, 120, 240]
      colors: Array,     // 예: ['#f00', '#0f0', '#00f']
      texts: Array,      // 예: ['L1', 'L2', 'L3']
      dcount: Number     // 예: 3
    },
    setup(props) {
      const canvasRef = ref(null);
      let ctx = null;
      let radius = 0;
  
      onMounted(() => {
        const canvas = canvasRef.value;
        ctx = canvas.getContext('2d');
        radius = canvas.height / 2;
        ctx.translate(radius, radius);
        radius *= 0.9;
  
        drawAll();
      });
  
      watch(() => [props.degree, props.colors, props.texts, props.dcount], drawAll, { deep: true });
  
      function drawAll() {
        if (!ctx) return;
        ctx.clearRect(-radius, -radius, radius * 2, radius * 2);
        drawFace(ctx, radius, props.degree, props.colors, props.texts, props.dcount);
      }
  
      function drawFace(ctx, radius, degree, colors, texts, dcount) {
        ctx.beginPath();
        ctx.arc(0, 0, radius, 0, 2 * Math.PI);
        ctx.fillStyle = '#ccc';
        ctx.fill();
        ctx.strokeStyle = "#808080";
        ctx.lineWidth = radius * 0.01;
        ctx.stroke();
  
        for (let i = 0; i < 3; i++) {
          ctx.beginPath();
          ctx.arc(0, 0, radius * i / 3, 0, 2 * Math.PI);
          ctx.stroke();
        }
  
        for (let i = 0; i < 12; i++) {
          const second = (i * Math.PI / 6);
          drawHand(ctx, second, radius * 0.99, radius * 0.01);
        }
  
        for (let i = 0; i < dcount; i++) {
          ctx.strokeStyle = colors[i];
          const second = -(((degree[i] / 10) - 9) * Math.PI / 18);
          if (i < 3) {
            drawHand_arrow(ctx, second, radius * 0.99, radius * 0.03, colors[i], texts[i], 0, (i * 23) - 120);
          } else {
            drawHand_arrow(ctx, second, radius * 0.55, radius * 0.03, colors[i], texts[i], 33, ((i - 3) * 23) - 120);
          }
        }
  
        ctx.strokeStyle = "#808080";
        drawNumbers(ctx, radius + 45);
        ctx.beginPath();
        ctx.arc(0, 0, radius * 0.01, 0, 2 * Math.PI);
        ctx.fillStyle = '#333';
        ctx.fill();
      }
  
      function drawHand(ctx, pos, length, width) {
        ctx.beginPath();
        ctx.lineWidth = width;
        ctx.lineCap = "square";
        ctx.moveTo(0, 0);
        ctx.rotate(pos);
        ctx.lineTo(0, -length);
        ctx.stroke();
        ctx.rotate(-pos);
      }
  
      function drawHand_arrow(ctx, pos, length, width, c_order, texts) {
        ctx.beginPath();
        ctx.lineWidth = width;
        ctx.lineCap = "square";
        ctx.moveTo(0, 0);
        ctx.rotate(pos);
        ctx.lineTo(0, -length);
        ctx.lineTo(-1, -length - 2);
        ctx.lineTo(-4, -length + 2);
        ctx.lineTo(4, -length + 2);
        ctx.lineTo(0, -length - 2);
        ctx.stroke();
  
        ctx.fillStyle = c_order;
        ctx.font = "bold 16px arial";
        ctx.fillText(texts, 5, -length - 2);
        ctx.rotate(-pos);
      }
  
      function drawNumbers(ctx, radius) {
        ctx.font = radius * 0.05 + "px arial";
        ctx.textBaseline = "middle";
        ctx.textAlign = "center";
        ctx.fillStyle = '#333';
  
        for (let num = 0; num < 13; num++) {
          let ang = (num + 3) * Math.PI / 6;
          ctx.rotate(ang);
          ctx.translate(0, -radius * 0.85);
          ctx.rotate(-ang);
          let tnum = 360 - (num * 30);
          if (tnum === 360) tnum = 0;
          if (tnum % 60 === 0) {
            ctx.fillText(tnum + '°', 0, 0);
          }
          ctx.rotate(ang);
          ctx.translate(0, radius * 0.85);
          ctx.rotate(-ang);
        }
      }
  
      return {
        canvasRef
      };
    }
  };
  </script>
  