<template>
  <div>
    <div v-for="(chartData, index) in chartDataList" :key="index" class="chart-container">
      <canvas :ref="(el) => chartRefs[index] = el"></canvas>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, nextTick } from 'vue';
import Chart from 'chart.js/auto';
import axios from 'axios';

const chartDataList = ref([]);
const chartRefs = ref([]);

onMounted(async () => {
  try {
    const response = await axios.get('http://localhost:8082/chartData');
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
      new Chart(ctx, {
        type: 'line',
        data: {
          labels: data.labels,
          datasets: data.datasets.map((dataset) => ({
            ...dataset,
            borderColor: dataset.borderColor || getRandomColor(),
            fill: true,
          })),
        },
        options: {
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
  const randomColor = () => {
    const letters = '0123456789ABCDEF';
    return '#' + letters[Math.floor(Math.random() * 16)];
  };
  return randomColor();
}
</script>

<style scoped>
.chart-container {
  margin-bottom: 20px;
}
</style>