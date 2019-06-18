
<template>
  <div class="alarmHistory_container">
    <el-row></el-row>
    <el-row class="title2" type="flex" justify="start" align="middle">
      <div></div>
    </el-row>
    <el-row class="content">
      <el-table :data="tableData" style="width: 100%" size="small">
        <el-table-column prop="groupName" label="分组" align="center"></el-table-column>
        <el-table-column prop="hostIp.String" label="节点" align="center"></el-table-column>
        <el-table-column prop="type" label="告警类型" align="center">
          <template slot-scope="scope">
            <span>{{scope.row.type | formatType}}</span>
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="告警时间" align="center">
          <template slot-scope="scope">
            <span>{{scope.row.createTime*1000 | formatTime}}</span>
          </template>
        </el-table-column>
        <el-table-column prop="warning" label="告警内容" align="center"></el-table-column>
      </el-table>
      <div class="corner left-top"></div>
      <div class="corner left-bottom"></div>
      <div class="corner right-top"></div>
      <div class="corner right-bottom"></div>
    </el-row>
    <el-row type="flex" justify="center" align="middle">
      <el-pagination
        @current-change="handleCurrentChange"
        :current-page.sync="requestContent.pageNum"
        :page-size="requestContent.pageSize"
        layout="total, prev, pager, next"
        :total="total"
      ></el-pagination>
    </el-row>
  </div>
</template>
<script>
import { alarmHistory } from "@/api/requestMethods";
import { nodeList } from "@/api/requestMethods";
import { formatTime } from "@/api/dateUtil";
export default {
  data() {
    return {
      requestContent: {
        pageNum: 1,
        pageSize: 10
      },
      total: 0,
      tableData: [],
      selection: [],
      nodeMap: {},
      groupMap: {}

      // tableData1: [
      //   {
      //     group: "分组1",
      //     node: "节点A",
      //     type: "业务监控告警",
      //     time: "2019年6月1日 19：30：00",
      //     content: "代扣地址余额不足"
      //   },
      //   {
      //     group: "分组1",
      //     node: "节点A",
      //     type: "业务监控告警",
      //     time: "2019年6月1日 19：30：00",
      //     content: "代扣地址余额不足"
      //   },
      //   {
      //     group: "分组1",
      //     node: "节点A",
      //     type: "业务监控告警",
      //     time: "2019年6月1日 19：30：00",
      //     content: "代扣地址余额不足"
      //   }
      // ],
      // currentPage: 1,
      // pageSize: 10,
      // total: 100
    };
  },
  methods: {
    handleDialogClose() {
      this.dialogVisible = false;
    },
    requestData() {
      alarmHistory(this.requestContent)
        .then(res => {
          this.tableData = res.data.values;
          this.total = res.data.total;
        })
        .catch(err => {
          this.$message({
            message: "网络错误，请稍后再试！",
            type: "error",
            offset: 125
          });
        });
    },
    handleCurrentChange() {
      this.requestData();
    }
  },
  mounted() {
    this.$store.commit("updateDefaultActive", "6");

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
    formatType(type) {
      if (type < 4) {
        return "资源监控告警";
      } else {
        return "业务监控告警";
      }
    }
  }
};
</script>

<style>
.alarmHistory_container {
}
</style>