<template>
  <div class="chart-wrapper">
    <div class="chart-container" v-for="(chartData, index) in chartDataList" :key="index">
      <canvas :ref="(el) => chartRefs[index] = el"></canvas>
    </div>
  </div>
</template>

<style scoped>
.chart-wrapper {
  display: flex;
  margin-left: 10%;
  flex-wrap: wrap;
  /* 添加 flex-wrap 属性，使容器换行显示 */
  justify-content: center;
  align-items: center;
  width: 900px;
  height: 100vh;
}

.chart-container {
  width: 100%;
  /* 实现居中对齐 */
  margin-bottom: 20px;
}
</style>

<script setup>
import { ref, onMounted, watch, nextTick } from 'vue';
import Chart from 'chart.js/auto';
import axios from 'axios';

const chartDataList = ref([]);
const chartRefs = ref([]);

function getHostnameWithoutPort() {
  const host = window.location.host;
  return host.split(':').shift();
}

onMounted(async () => {
  try {
    let host = getHostnameWithoutPort();
    let url = `http://${host}:8082/chartData`;
    console.log(">>> host", url)

    const response = await axios.get(url);
    chartDataList.value = response.data;
    await nextTick();
    renderCharts();
  } catch (error) {
    console.error('数据获取失败：', error);
  }
});

watch(chartDataList, () => {
  nextTick(() => {
    renderCharts();
  });
});

function renderCharts() {
  chartDataList.value.forEach((data, index) => {
    const canvas = chartRefs.value[index];
    if (canvas && canvas.getContext) {
      const ctx = canvas.getContext('2d');
      // if (ctx._chart) {
      //   Chart.destroy(ctx._chart);
      // }
      if (ctx._chart) {
        ctx._chart.destroy();
      }
      new Chart(ctx, {
        type: 'line',
        data: {
          labels: data.labels,
          datasets: data.datasets.map((dataset) => ({
            ...dataset,
            // borderColor: dataset.borderColor||randomRGBA(),
             borderColor: dataset.borderColor||getRandomColor(),
            fill: true,
          })),
        },
        options: {
          plugins: {
            title: {
              display: true,
              text: data.title,
            },
          },
          scales: {
            y: {
              beginAtZero: true
            }
          }
        }
      });
    }
  });
}

function getRandomColor() {
  // 生成十六进制颜色值
  const randomColor = () => {
    const letters = '0123456789ABCDEF';
    let color = '#';
    for (let i = 0; i < 6; i++) {
      color += letters[Math.floor(Math.random() * 16)];
    }
    return color;
  };

  // 生成rgba颜色值
  const randomRGBA = () => {
    const r = Math.floor(Math.random() * 256);
    const g = Math.floor(Math.random() * 256);
    const b = Math.floor(Math.random() * 256);
    const a = Math.random().toFixed(2); // 随机生成透明度
    return `rgba(${r}, ${g}, ${b}, ${a})`;
  };

  // 从上面两个函数中随机选择一个
  const colorFuncs = [randomColor, randomRGBA];
  return colorFuncs[Math.floor(Math.random() * colorFuncs.length)]();
}
</script>

<style scoped>
.chart-container {
  margin-bottom: 20px;
}
</style>
