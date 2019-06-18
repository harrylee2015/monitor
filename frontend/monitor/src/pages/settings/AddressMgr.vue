<template>
  <div class="addressMgr_container">
    <el-row></el-row>
    <el-row class="title2" type="flex" justify="space-between" align="middle">
      <div></div>
      <div>
        <!-- <el-button-group> -->
        <el-button type="info" @click="handleAddAddr">新 增</el-button>
        <!-- <el-button type="danger">删 除</el-button> -->
        <!-- </el-button-group> -->
      </div>
    </el-row>
    <el-row class="content">
      <el-table :data="tableData" style="width: 100%" size="small">
        <!-- <el-table-column type="selection" width="40"></el-table-column> -->
        <el-table-column prop="id" label="分组ID" align="center" width="100"></el-table-column>
        <el-table-column prop="groupName" label="分组名称" align="center" width="100"></el-table-column>
        <el-table-column prop="address" label="代扣地址" align="center"></el-table-column>
        <el-table-column label="操作" align="center" width="90">
          <template slot-scope="scope">
            <span class="opt-mod" @click="handleUpdateAddr(scope.$index, scope.row)">修改</span>
            <span class="opt-del" @click="handleDeleteAddr(scope.$index, scope.row)">删除</span>
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
    <el-dialog title="新增代扣地址" :visible="dialogVisible" width="605px" @close="handleDialogClose">
      <el-form
        label-width="110px"
        style="width:510px;"
        ref="addForm"
        :model="addForm"
        :rules="addRules"
      >
        <el-form-item label="所属分组：" prop="option" >
          <el-select
            v-model="addForm.option"
            placeholder="请选择"
            style="width: 100%;"
            :disabled="groupDisable"
          >
            <el-option
              v-for="item in options"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="代扣地址：" prop="address">
          <el-input v-model="addForm.address"></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer">
        <el-button type="primary" @click="dialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="addAddr">保 存</el-button>
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
  addrList,
  addrAdd,
  addrUpdate,
  addrDel
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
      addForm: {
        id: 0,
        option: "",
        address: ""
      },
      addRules: {
        option: [{ required: true, message: "请选择分组", trigger: "change" }],
        address: [{ required: true, message: "请输入地址", trigger: "blur" }]
      },
      options: [],
      groupDisable: false
    };
  },
  methods: {
    requestData() {
      addrList(this.requestContent)
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

    handleAddAddr() {
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
    handleUpdateAddr(index, row) {
      let groupVal = row.groupId + ":" + row.groupName;
      this.options.push({ label: row.groupName, value: groupVal });

      this.addForm.id = row.id;
      this.addForm.option = groupVal;
      this.addForm.address = row.address;

      this.groupDisable = true;
      this.dialogVisible = true;
    },
    getAddRequestContent() {
      let data = {};
      if (this.addForm.id) {
        data.id = this.addForm.id;
      }
      let groupVal = this.addForm.option.split(":");
      data.groupId = groupVal[0] * 1;
      data.groupName = groupVal[1];
      data.address = this.addForm.address;

      return data;
    },
    addAddr() {
      this.$refs.addForm.validate(valid => {
        if (valid) {
          let data = this.getAddRequestContent();
          
          if(this.addForm.id == 0){
            addrAdd(data)
              .then(res => {
                this.resMsg(res, "添加");
              })
              .catch(err => {
                this.errMsg(err);
              });
          }else{
            addrUpdate(data)
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
    handleDeleteAddr(index, row) {
      this.$confirm("此操作将永久删除该代扣地址, 是否继续?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      })
        .then(() => {
          addrDel(row.id)
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
      this.addForm.id = 0;
      this.addForm.address = "";
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
    this.$store.commit("updateDefaultActive", "3");
    this.requestData();
  }
};
</script>
<style>
.addressMgr_container {
}
</style>
