<template>
  <div class="servicem_container">
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
      <div>节点监控</div>
    </el-row>
    <el-row class="content">
      <el-table :data="nodeTableData" style="width: 100%" size="small">
        <el-table-column prop="hostIp" label="节点" align="center" width="130"></el-table-column>
        <el-table-column prop="serverStatus" label="进程" align="center" width="100">
          <template slot-scope="scope">
            <span v-if="scope.row.isSync == 0">运行</span>
            <span v-else>停止</span>
          </template>
        </el-table-column>
        <el-table-column prop="serverPort" label="RPC端口" align="center" width="100"></el-table-column>
        <el-table-column prop="lastBlockHeight" label="最新区块高度" align="center" width="120"></el-table-column>
        <el-table-column prop="isSync" label="主链同步状态" align="center" width="120">
          <template slot-scope="scope">
            <span v-if="scope.row.isSync == 0">已同步</span>
            <span v-else>未同步</span>
          </template>
        </el-table-column>
        <el-table-column prop="lastBlockHash" label="最新区块Hash" align="center">
          <template slot-scope="scope">{{scope.row.lastBlockHash | formatHash}}</template>
        </el-table-column>
      </el-table>
      <div class="corner left-top"></div>
      <div class="corner left-bottom"></div>
      <div class="corner right-top"></div>
      <div class="corner right-bottom"></div>
    </el-row>
    <el-row class="title" type="flex" align="middle">
      <div></div>
      <div>代扣地址余额： {{balance}} BTY &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; 更新时间： {{balanceUpdateTime | formatTime}}</div>
    </el-row>
    <el-row class="content chart">
      <!-- <div class="title">代扣地址余额和时间的关系曲线</div> -->
      <balance-chart style="height: 240px;width: 500px;" ref="balanceChart"></balance-chart>
    </el-row>
    <el-row class="title" type="flex" align="middle">
      <div></div>
      <div>区块一致性状态：{{checkResult |formatCheck}}</div>
    </el-row>
    <el-row class="content">
      <el-table
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
    </el-row>
    <!-- </el-collapse-item>
      <el-collapse-item name="2">
        <template slot="title">
          <font class="group">分组2</font>
        </template>
        <div>控制反馈：通过界面样式和交互动效让用户可以清晰的感知自己的操作；</div>
        <div>页面反馈：操作后，通过页面元素的变化清晰地展现当前状态。</div>
      </el-collapse-item>
    </el-collapse>-->
  </div>
</template>

<script>
import {
  groupList,
  nodeMonitor,
  addrBalance,
  balancelist,
  buswarning,
  hashCheck,
  warningRemove,
  warningBatchRemove
} from "@/api/requestMethods";
import { formatTime } from "@/api/dateUtil";
import BalanceChart from "./BalanceChart.vue";
export default {
  components: {
    BalanceChart
  },
  data() {
    return {
      group: "",
      options: [],
      nodeTableData: [],
      balance: 0,
      balanceUpdateTime: 0,
      alarmTableData: [],
      selection: [],
      checkResult: {}
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
            this.requestNode(this.group);
            this.requestBalance(this.group);
            this.requestAlarm(this.group);
            this.requestHashCheck(this.group);
          }
        })
        .catch(err => {
          this.errMsg(err);
        });
    },
    requestNode(groupId) {
      nodeMonitor(groupId).then(res => {
        this.nodeTableData = res.data;
      });
    },
    requestHashCheck(groupId) {
      hashCheck(groupId).then(res => {
        this.checkResult = res.data;
      });
    },
    requestBalance(groupId) {
      addrBalance(groupId).then(res => {
        this.balance = res.data[0].balance / 100000000;
        this.balanceUpdateTime = res.data[0].createTime * 1000;
      });

      this.$refs.balanceChart.updateChart(groupId);
    },
    requestAlarm(groupId) {
      buswarning(groupId).then(res => {
        this.alarmTableData = res.data;
      });
    },

    handleGroupChange(val) {
      this.group = val;
      this.requestNode(this.group);
      this.requestBalance(this.group);
      this.requestAlarm(this.group);
      this.requestHashCheck(this.group);
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
.servicem_container {
  /* background: transparent !important; */
}
.servicem_container .el-row:nth-of-type(1) {
  margin-bottom: 20px;
}
.servicem_container .groups .el-input .el-input__inner {
  border: none;
  border-radius: 3px;
  background: #0862bc;
  caret-color: #ffffff;
  color: #ffffff;
}
.servicem_container .groups .el-input .el-input__suffix {
  background: #0689e5;
  right: 0px;
  border-left: 2px solid #020f34;
  border-top-right-radius: 3px;
  border-bottom-right-radius: 3px;
}
.servicem_container .groups .el-input .el-input__suffix-inner > i {
  color: #00d9dc;
}
</style>
