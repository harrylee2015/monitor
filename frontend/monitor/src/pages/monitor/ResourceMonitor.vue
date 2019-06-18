<template>
  <div class="resourcem_container">
    <!-- <el-collapse v-model="activeNames">
      <el-collapse-item name="1">
        <template slot="title">
          <font class="group">分组1</font>
    </template>-->
    <el-row type="flex" align="start">
      <el-select class="groups" v-model="group" @change="handleGroupChange">
        <el-option
          v-for="item in options"
          :key="item.value"
          :label="item.label"
          :value="item.value"
        ></el-option>
      </el-select>
    </el-row>
    <el-row class="title" type="flex" align="middle">
      <div></div>
      <div>节点列表</div>
    </el-row>
    <el-row class="content">
      <el-collapse v-model="activeNodes" @change="handleNodeChange" accordion>
        <el-collapse-item v-for="item in nodes" :key="item.hostId" :name="item.hostId">
          <template slot="title">
            <font class="group">{{item.hostName}}</font>
          </template>
          <div></div>
          <div class="content chart">
            <!-- <div class="title">代扣地址余额和时间的关系曲线</div> -->
            <balance-chart style="height: 240px;width: 500px;" ref="balanceChart"></balance-chart>
          </div>
          <div>
            <el-table
              v-loading="loading"
              :data="alarmTableData"
              style="width: 100%"
              size="small"
              @selection-change="handleSelectionChange"
            >
              <el-table-column type="selection" width="55"></el-table-column>
              <el-table-column prop="warning" align="left">
                <template slot="header">
                  <font color="red">异常警告</font>
                </template>
              </el-table-column>
              <el-table-column align="center" width="60">
                <template slot="header">
                  <span class="opt-head-clear" @click="clearAlarm('M', null, null)">清除</span>
                </template>
                <template slot-scope="scope">
                  <span class="opt-clear" @click="clearAlarm('S', scope.$index, scope.row)">清除</span>
                </template>
              </el-table-column>
            </el-table>
            <div class="corner left-top"></div>
            <div class="corner left-bottom"></div>
            <div class="corner right-top"></div>
            <div class="corner right-bottom"></div>
          </div>
        </el-collapse-item>
      </el-collapse>
    </el-row>
  </div>
</template>

<script>
import {
  groupList,
  nodeMonitor,
  addrBalance,
  balancelist,
  reswarning,
  warningRemove,
  warningBatchRemove,
  queryNodeListByGroupId
} from "@/api/requestMethods";
import { formatTime } from "@/api/dateUtil";
import BalanceChart from "./BalanceChart.vue";
import CpuChart from "./CpuChart.vue";
import DiskChart from "./DiskChart.vue";
import RamChart from "./RamChart.vue";
export default {
  components: {
    BalanceChart,
    CpuChart,
    DiskChart,
    RamChart
  },
  data() {
    return {
      group: "",
      nodes: [],
      activeNodes: [],
      options: [],
      nodeTableData: [],
      alarmTableData: [],
      selection: [],
      checkResult: {},
      loading: false
    };
  },
  methods: {
    requestData() {
      let data = {
        pageNum: 1,
        pageSize: 1000
      };
      groupList(data)
        .then(res => {
          if (res.data.values) {
            for (let item of res.data.values) {
              this.options.push({
                label: item.groupName,
                value: item.groupId * 1
              });
            }
            this.group = this.options[0].value;
            this.requestNodeList(this.group);
          }
        })
        .catch(err => {
          this.errMsg(err);
        });
    },
    requestNode(groupId) {
      nodeMonitor(groupId).then(res => {
        this.nodeTableData = res.data;
        // this.nodes = res.data;
      });
    },
    requestNodeList(groupId) {
      queryNodeListByGroupId(groupId)
        .then(res => {
          this.nodes = res.data;
        })
        .catch(err => {
          this.errMsg(err);
        });
    },
    requestBalance(groupId) {
      addrBalance(groupId).then(res => {
        this.balance = res.data[0].balance / 100000000;
        this.balanceUpdateTime = res.data[0].createTime * 1000;
      });

      this.$refs.balanceChart.updateChart(groupId);
    },
    requestAlarm(hostId) {
      reswarning(hostId).then(res => {
        this.alarmTableData = res.data;
      });
    },

    handleGroupChange(val) {
      this.requestData();
    },
    handleNodeChange(val) {
      this.requestAlarm(val);
    },
    handleSelectionChange(val) {
      this.selection = val;
    },
    clearAlarm(type, index, row) {
      if (type == "M") {
        warningBatchRemove(this.alarmTableData).then(res => {
          if (res.status == 200) {
            this.alarmTableData = [];
          }
        });
      } else {
        warningRemove(row).then(res => {
          if (res.status == 200) {
            this.alarmTableData.splice(index, 1);
          }
        });
      }
    },

    resMsg(res, opt) {
      if (res.status == 200) {
        this.dialogVisible = false;
        this.requestData();
        this.$message({
          message: opt + "成功！",
          type: "success",
          offset: 125
        });
      } else {
        this.$message({
          message: opt + "失败！",
          type: "error",
          offset: 125
        });
      }
    },
    errMsg(err) {
      this.$message({
        message: "网络错误，请稍后再试！",
        type: "error",
        offset: 125
      });
    }
  },
  mounted() {
    this.$store.commit("updateDefaultActive", "4");
    this.requestData();
  },
  filters: {
    formatHash(hash) {
      return hash.substring(0, 25) + "...";
    },
    formatTime(time) {
      let date = new Date(time);
      return formatTime(date, "yyyy-MM-dd hh:mm:ss");
    },
    formatCheck(result) {
      if (result.isConsistent == true) {
        return "正常";
      } else {
        return "异常";
      }
    }
  }
};
</script>

<style>
.resourcem_container {
  /* background: transparent !important; */
}
.resourcem_container .el-row:nth-of-type(1) {
  margin-bottom: 20px;
}
.resourcem_container .groups .el-input .el-input__inner {
  border: none;
  border-radius: 3px;
  background: #0862bc;
  caret-color: #ffffff;
  color: #ffffff;
}
.resourcem_container .groups .el-input .el-input__suffix {
  background: #0689e5;
  right: 0px;
  border-left: 2px solid #020f34;
  border-top-right-radius: 3px;
  border-bottom-right-radius: 3px;
}
.resourcem_container .groups .el-input .el-input__suffix-inner > i {
  color: #00d9dc;
}
</style>
