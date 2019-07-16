<template>
  <div class="groupMgr_container">
    <el-row></el-row>
    <el-row class="title2" type="flex" justify="space-between" align="middle">
      <div></div>
      <div>
        <!-- <el-button-group> -->
        <el-button type="info" @click="handleAddGroup">新 增</el-button>
        <!-- <el-button type="danger" @click="handleDelGroup('M')">删 除</el-button> -->
        <!-- </el-button-group> -->
      </div>
    </el-row>
    <el-row class="content">
      <el-table
        :data="tableData"
        style="width: 100%"
        size="small"
        @selection-change="handleSelectionChange"
      >
        <!-- <el-table-column type="selection" width="40"></el-table-column> -->
        <el-table-column prop="groupId" label="序号" align="center" width="60"></el-table-column>
        <el-table-column prop="groupName" label="组名" align="center" width="100"></el-table-column>
        <el-table-column prop="title" label="前缀" align="center" width="180">
          <template slot-scope="scope">
            <span v-if="scope.row.title">{{scope.row.title}}</span>
            <span v-else>--</span>
          </template>
        </el-table-column>
        <el-table-column prop="describe" label="介绍" align="center"></el-table-column>
        <el-table-column prop="email" label="运维邮箱" align="center"></el-table-column>
        <el-table-column label="操作" align="center" width="90">
          <template slot-scope="scope">
            <span class="opt-mod" @click="handleUpdateGroup(scope.$index, scope.row)">修改</span>
            <span class="opt-del" @click="handleDelGroup('S', scope.$index, scope.row)">删除</span>
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
    <el-dialog
      :title="dialogTitle"
      :visible="dialogVisible"
      width="650px"
      @close="handleDialogClose"
    >
      <el-form label-width="110px" style="width:550px;" :model="addForm" ref="addForm">
        <el-form-item
          label="分组名称："
          prop="groupName"
          :rules="[{ required: true, message: '请输入分组名称', trigger: 'blur' }]"
        >
          <el-input v-model="addForm.groupName"></el-input>
        </el-form-item>
        <el-form-item
          label="前  缀："
          prop="title"
          :rules="[{ required: true, message: '请输入前缀', trigger: 'blur' }]"
        >
          <el-input v-model="addForm.title"></el-input>
        </el-form-item>
        <el-form-item
          label="运维邮箱："
          prop="email"
          :rules="[{ required: true, message: '请输入邮箱地址', trigger: 'blur' }]"
        >
          <el-input v-model="addForm.email"></el-input>
        </el-form-item>
        <el-form-item label="介  绍：" prop="describe">
          <el-input type="textarea" v-model="addForm.describe" rows="6" resize="none"></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer">
        <el-button type="primary" @click="dialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="AddGroup">保 存</el-button>
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
  groupAdd,
  groupUpdate,
  groupDel
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
      selection: [],

      addForm: {
        groupId: 0,
        groupName: "",
        describe: "",
        title: "",
        email:""
      },

      dialogVisible: false,
      dialogTitle: "新增分组"
    };
  },
  methods: {
    requestData() {
      groupList(this.requestContent)
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
    },
    handleSelectionChange(selection) {
      this.selection = selection;
    },

    handleDelGroup(type, index, row) {
      if (type == "M") {
        // if (this.selection.length == 0) {
        //   this.$message({
        //     message: "请选择需要删除的分组！",
        //     type: "warning",
        //     offset: 125
        //   });
        // } else {
        //   for(let item of this.selection){
        //     // groupDel()
        //   }
        // }
      } else if (type == "S") {
        this.$confirm("此操作将永久删除该分组, 是否继续?", "提示", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning"
        })
          .then(() => {
            groupDel(row.groupId)
              .then(res => {
                this.resMsg(res, "删除");
              })
              .catch(err => {
                this.errMsg(err);
              });
          })
          .catch(() => {});
      }
    },

    handleAddGroup() {
      this.dialogTitle = "新增分组";
      this.dialogVisible = true;
    },
    handleUpdateGroup(index, row) {
      this.dialogTitle = "修改分组";
      this.addForm.groupId = row.groupId;
      this.addForm.groupName = row.groupName;
      this.addForm.describe = row.describe;
      this.addForm.title = row.title;
      this.addForm.email = row.email;
      this.dialogVisible = true;
    },
    AddGroup() {
      this.$refs.addForm.validate(valid => {
        if (valid) {
          if (this.addForm.groupId == 0) {
            groupAdd(this.addForm)
              .then(res => {
                this.resMsg(res, "添加");
              })
              .catch(err => {
                this.errMsg(err);
              });
          } else {
            groupUpdate(this.addForm)
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
    handleDialogClose() {
      this.dialogVisible = false;
      this.$refs.addForm.resetFields();
      this.addForm.groupId = 0;
      this.addForm.groupName = "";
      this.addForm.describe = "";
      this.addForm.title = "";
      this.addForm.email = "";
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
    this.$store.commit("updateDefaultActive", "1");
    this.requestData();
  }
};
</script>
<style>
.groupMgr_container {
}
</style>
