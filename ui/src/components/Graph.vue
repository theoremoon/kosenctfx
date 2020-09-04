<template>
  <canvas ref="chart"></canvas>
</template>

<script>
import Chart from "chart.js/dist/Chart";
export default {
  props: {
    chartdata: {
      type: Array,
      default: null
    }
  },
  mounted() {
    let ctx = this.$refs.chart.getContext("2d");
    this.chart = new Chart(ctx, {
      type: "line",
      data: [],
      options: {
        animation: {
          duration: 0
        },
        tooltips: {
          callbacks: {
            label(item, data) {
              let c = data.datasets[item.datasetIndex].data[item.index];
              return c.team + " " + c.name + ":" + c.score;
            }
          }
        },
        scales: {
          xAxes: [
            {
              type: "time",
              display: true,
              ticks: {
                padding: 10
              }
            }
          ],
          yAxes: [
            {
              ticks: {
                padding: 10
              }
            }
          ]
        }
      }
    });
    this.chart.data.datasets = this.chartdata;
    this.chart.update();
  },
  watch: {
    chartdata() {
      this.chart.data.datasets = this.chartdata;
      this.chart.update();
    }
  }
};
</script>
