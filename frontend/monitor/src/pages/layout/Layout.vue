<template>
  <el-container class="manager_wrap" ref="layout">
    <el-header class="manager_header">
      <manager-header></manager-header>
    </el-header>
    <el-container class="manager_container">
      <el-aside class="layout_aside">
        <side-bar class="manager-side-bar"></side-bar>
      </el-aside>
      <el-main class="main_container">
        <router-view></router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
import managerHeader from "../../components/ManagerHeader";
import sideBar from "../../components/SideBar";
export default {
  inject: ["reload"],
  components: {
    "manager-header": managerHeader,
    "side-bar": sideBar,
  }
};
</script>

<style scoped>
.manager_wrap {
  width: 1150px;
  /* height: 2000px; */
  background: url("../../../static/img/background.png") top center no-repeat;
  flex: 0 0 1150px;
}
.manager_header {
  margin-top: 5px;
  height: 95px !important;
  background: transparent;
  width: 1150px;
  padding: 0 0;
}
.manager_container {
  margin-top: 20px;
  height: calc(100vh - 120px);
  border-top: 1px solid #022d4d;
}
.layout_aside {
  width: 170px !important;
  background: #020f34;
  border-right: 1px solid #022d4d;
}
.manager-side-bar {
  width: 100%;
}
.nav_container {
  padding: 0;
  height: 56px !important;
}
.main_container {
  padding: 30px 20px;
  overflow: auto;
}
.main_container::-webkit-scrollbar
{
    width: 5px;
}
.main_container::-webkit-scrollbar-thumb
{
    border-radius: 10px;
    background-color: #022d4d;
}
</style>
