<template>
  <div class="nodeMgr_container">
    <el-row></el-row>
    <el-row class="title2" type="flex" justify="space-between" align="middle">
      <div></div>
      <div>
        <!-- <el-button-group> -->
        <el-button type="info" @click="handleAddNode">新 增</el-button>
        <!-- <el-button type="danger">删 除</el-button> -->
        <!-- </el-button-group> -->
      </div>
    </el-row>
    <el-row class="content">
      <el-table :data="tableData" style="width: 100%">
        <!-- <el-table-column type="selection" width="20"></el-table-column> -->
        <el-table-column prop="hostId" label="节点ID" align="center" width="80"></el-table-column>
        <el-table-column prop="hostName" label="节点名称" align="center"></el-table-column>
        <el-table-column prop="groupName" label="所属分组" align="center"></el-table-column>
        <el-table-column prop="hostIp" label="节点IP" align="center" width="125"></el-table-column>
        <el-table-column prop="SSHPort" label="SSH端口" align="center">
          <template slot-scope="scope">
            <span v-if="scope.row.SSHPort && scope.row.SSHPort != 0">{{scope.row.SSHPort}}</span>
            <span v-else>--</span>
          </template>
        </el-table-column>
        <el-table-column prop="userName" label="用户名" align="center">
          <template slot-scope="scope">
            <span v-if="scope.row.userName">{{scope.row.userName}}</span>
            <span v-else>--</span>
          </template>
        </el-table-column>
        <el-table-column prop="passWd" label="密码" align="center" width="60">
          <template slot-scope="scope">
            <span v-if="scope.row.passWd">***</span>
            <span v-else>--</span>
          </template>
        </el-table-column>
        <el-table-column prop="processName" label="进程" align="center"></el-table-column>
        <el-table-column prop="serverPort" label="服务端口" align="center" width="90"></el-table-column>
        <el-table-column label="操作" align="center" width="100">
          <template slot-scope="scope">
            <span class="opt-mod" @click="handleUpdateNode(scope.$index, scope.row)">修改</span>
            <span class="opt-del" @click="handleDeleteNode(scope.$index, scope.row)">删除</span>
          </template>
        </el-table-column>
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
    <el-dialog title="新增节点" :visible="dialogVisible" width="705px" @close="handleDialogClose">
      <el-form label-width="110px" ref="addForm" :model="addForm" :rules="addRules">
        <el-row type="flex" justify="start">
          <el-col :span="11">
            <el-form-item label="节点名称：" prop="hostName">
              <el-input v-model="addForm.hostName"></el-input>
            </el-form-item>
            <el-form-item label="所属分组：" prop="option">
              <el-select v-model="addForm.option" placeholder="请选择" :disabled="groupDisable">
                <el-option
                  v-for="item in options"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                ></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="节点IP：" prop="hostIp">
              <el-input v-model="addForm.hostIp"></el-input>
            </el-form-item>
            <el-form-item label="SSH端口：" prop="SSHPort">
              <el-input v-model="addForm.SSHPort"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="1"></el-col>
          <el-col :span="11">
            <el-form-item label="用户名：" prop="userName">
              <el-input v-model="addForm.userName"></el-input>
            </el-form-item>
            <el-form-item label="密码：" prop="passWd">
              <el-input show-password v-model="addForm.passWd"></el-input>
            </el-form-item>
            <el-form-item label="进程名：" prop="processName">
              <el-input v-model="addForm.processName"></el-input>
            </el-form-item>
            <el-form-item label="服务端口：" prop="serverPort">
              <el-input v-model="addForm.serverPort"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>

      <span slot="footer">
        <el-button type="primary" @click="dialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="addNode">保 存</el-button>
      </span>
      <span class="corner left-top"></span>
      <span class="corner left-bottom"></span>
      <span class="corner right-top"></span>
      <span class="corner right-bottom"></span>
    </el-dialog>
  </div>
</template>
<script>
import {
  groupList,
  nodeList,
  nodeAdd,
  nodeUpdate,
  nodeDel
} from "@/api/requestMethods";
export default {
  data() {
    return {
      requestContent: {
        pageNum: 1,
        pageSize: 10
      },
      total: 0,
      tableData: [],

      dialogVisible: false,
      options: [],
      addRules: {
        hostName: [
          { required: true, message: "请输入节点名称", trigger: "blur" }
        ],
        option: [
          { required: true, message: "请选择所属分组", trigger: "change" }
        ],
        hostIp: [{ required: true, message: "请输入节点IP", trigger: "blur" }],
        // SSHPort: [{ required: true, message: "请输入SSH端口", trigger: "blur" }],
        // userName: [
        //   { required: true, message: "请输入用户名", trigger: "blur" }
        // ],
        // passWd: [{ required: true, message: "请输入密码", trigger: "blur" }],
        processName: [
          { required: true, message: "请输入进程名", trigger: "blur" }
        ],
        serverPort: [
          { required: true, message: "请输入服务端口", trigger: "blur" }
        ]
      },
      addForm: {
        hostId: 0,
        hostName: "",
        hostIp: "",
        SSHPort: "",
        userName: "",
        passWd: "",
        processName: "",
        serverPort: "",
        option: ""
      },
      groupDisable: false
    };
  },
  methods: {
    requestData() {
      nodeList(this.requestContent)
        .then(res => {
          this.tableData = res.data.values;
          this.total = res.data.total;
        })
        .catch(err => {
          this.errMsg(err);
        });
    },
    handleCurrentChange() {
      this.requestData();
    },

    handleAddNode() {
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
                value: item.groupId + ":" + item.groupName
              });
            }
          }
        })
        .catch(err => {
          this.errMsg(err);
        });
      this.dialogVisible = true;
    },
    handleUpdateNode(index, row) {
      let optionVal = row.groupId + ":" + row.groupName;
      this.options.push({ label: row.groupName, value: optionVal });
      this.groupDisable = true;

      this.addForm.hostId = row.hostId;
      this.addForm.hostName = row.hostName;
      this.addForm.hostIp = row.hostIp;
      this.addForm.SSHPort = row.SSHPort == 0 ? "" : row.SSHPort;
      this.addForm.userName = row.userName;
      this.addForm.passWd = row.passWd;
      this.addForm.processName = row.processName;
      this.addForm.serverPort = row.serverPort;
      this.addForm.option = optionVal;
      this.dialogVisible = true;
    },
    getAddRequestContent() {
      let data = {};
      if (this.addForm.hostId) {
        data.hostId = this.addForm.hostId;
      }
      data.hostName = this.addForm.hostName;
      let groupVal = this.addForm.option.split(":");
      data.groupId = groupVal[0] * 1;
      data.groupName = groupVal[1];
      data.hostIp = this.addForm.hostIp;
      data.isCheckResource = 0;
      if (this.addForm.SSHPort) {
        data.SSHPort = this.addForm.SSHPort * 1;
        data.isCheckResource = 1;
      }
      if (this.addForm.userName) {
        data.userName = this.addForm.userName;
      }
      if (this.addForm.passWd) {
        data.passWd = this.addForm.passWd;
      }
      data.processName = this.addForm.processName;
      data.serverPort = this.addForm.serverPort * 1;

      return data;
    },
    addNode() {
      this.$refs.addForm.validate(valid => {
        if (valid) {
          if (
            !(
              (this.addForm.SSHPort &&
                this.addForm.userName &&
                this.addForm.passWd) ||
              (!this.addForm.SSHPort &&
                !this.addForm.userName &&
                !this.addForm.passWd)
            )
          ) {
            this.$message({
              message: "参数错误！",
              type: "warning",
              offset: 125
            });
            return false;
          }

          let data = this.getAddRequestContent();

          if (this.addForm.hostId == 0) {
            nodeAdd(data)
              .then(res => {
                this.resMsg(res, "添加");
              })
              .catch(err => {
                this.errMsg(err);
              });
          } else {
            nodeUpdate(data)
              .then(res => {
                this.resMsg(res, "修改");
              })
              .catch(err => {
                this.errMsg(err);
              });
          }
        }
      });
    },
    handleDeleteNode(index, row) {
      this.$confirm("此操作将永久删除该节点, 是否继续?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      })
        .then(() => {
          nodeDel(row.hostId)
            .then(res => {
              this.resMsg(res, "删除");
            })
            .catch(err => {
              this.errMsg(err);
            });
        })
        .catch(() => {});
    },
    handleDialogClose() {
      this.dialogVisible = false;
      this.$refs.addForm.resetFields();
      this.groupDisable = false;
      this.options = [];
      this.addForm.hostId = 0;
      this.addForm.hostName = "";
      this.addForm.hostIp = "";
      this.addForm.SSHPort = "";
      this.addForm.userName = "";
      this.addForm.passWd = "";
      this.addForm.processName = "";
      this.addForm.serverPort = "";
      this.addForm.option = "";
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
    this.$store.commit("updateDefaultActive", "2");
    this.requestData();
  }
};
</script>
<style>
.nodeMgr_container {
}
</style>
