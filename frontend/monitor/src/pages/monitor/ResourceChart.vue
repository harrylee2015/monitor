<template>
  <div ref="dom"></div>
</template>

<script>
import { resourceList } from "@/api/requestMethods";
import { formatTime } from "@/api/dateUtil";
import "../../../static/themes/chartTheme";
import echarts from "echarts";
import { on, off } from "@/libs/tools";
export default {
  name: "serviceRequests",
  data() {
    return {
      dom: null,
      option: {
        // tooltip: {
        //   trigger: "axis"
        // },
        tooltip: {
          trigger: "axis",
          position: function(pt) {
            return [pt[0], "10%"];
          }
        },
        // title: {
        //   // left: "center",
        //   text: "节点资源和时间得关系折线图"
        // },
        legend: {
          data: ["CPU使用率(百分比)", "内存使用率(百分比)", "磁盘空间使用率(百分比)"]
        },
        grid: {
          left: "3%",
          right: "4%",
          bottom: "3%",
          containLabel: true
        },
        // toolbox: {
        //   feature: {
        //     saveAsImage: {}
        //   }
        // },
        xAxis: {
          type: "category",
          boundaryGap: true,
          data: []
        },
        yAxis: {
          type: "value",
          boundaryGap: [0, "100%"]
        },
        dataZoom: [
          {
            type: "inside",
            start: 30,
            end: 100
          },
          {
            start: 0,
            end: 10,
            handleIcon:
              "M10.7,11.9v-1.3H9.3v1.3c-4.9,0.3-8.8,4.4-8.8,9.4c0,5,3.9,9.1,8.8,9.4v1.3h1.3v-1.3c4.9-0.3,8.8-4.4,8.8-9.4C19.5,16.3,15.6,12.2,10.7,11.9z M13.3,24.4H6.7V23h6.6V24.4z M13.3,19.6H6.7v-1.4h6.6V19.6z",
            handleSize: "80%",
            handleStyle: {
              color: "#fff",
              shadowBlur: 3,
              shadowColor: "rgba(0, 0, 0, 0.6)",
              shadowOffsetX: 2,
              shadowOffsetY: 2
            }
          }
        ],
        // series: [
        //   {
        //     name: "代扣地址余额",
        //     type: "line",
        //     stack: "总量",
        //     areaStyle: {
        //       normal: {
        //         color: "#2d8cf0"
        //       }
        //     },
        //     data: []
        //   }
        // ],
        series: [
          {
            name: "CPU使用率(百分比)",
            type: "line",
            stack: "总量",
            data: []
          },
          {
            name: "内存使用率(百分比)",
            type: "line",
            stack: "总量",
            data: []
          },
          {
            name: "磁盘空间使用率(百分比)",
            type: "line",
            stack: "总量",
            data: []
          }
        ]
      }
    };
  },
  methods: {
    resize() {
      this.dom.resize();
    },
    updateResourceChart(hostId) {
      resourceList(hostId).then(res => {
        let chartData = {
          time: [],
          cpu: [],
          mem: [],
          disk: []
        };

        for (let resource of res.data) {
          let date = new Date(resource.createTime * 1000);
          chartData.time.push(formatTime(date, "MM-dd hh:mm"));
          chartData.cpu.push(resource.cpuUsedPercent);
          chartData.mem.push(resource.memUsedPercent);
          chartData.disk.push(resource.diskUsedPercent);
        }

        this.option.xAxis.data = chartData.time;
        this.option.series[0].data = chartData.cpu;
        this.option.series[1].data = chartData.mem;
        this.option.series[2].data = chartData.disk;
        this.dom.setOption(this.option);
      });
    }
  },
  mounted() {
    this.$nextTick(() => {
      this.dom = echarts.init(this.$refs.dom, "macarons");
      on(window, "resize", this.resize);
    });
  },
  beforeDestroy() {
    off(window, "resize", this.resize);
  }
};
</script>
