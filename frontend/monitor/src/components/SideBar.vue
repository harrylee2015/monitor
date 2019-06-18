<template>
  <div class="sidebar_container">
    <el-menu
      class="el-menu-vertical-demo"
      text-color="#00FFFF"
      active-text-color="#fff"
      :unique-opened="true"
      :default-active="$store.state.defaultActive"
    >
      <el-menu-item
        index="1"
        @click="toLink('group-mgr')"
      >分组管理</el-menu-item>
      <el-menu-item
        index="2"
        @click="toLink('node-mgr')"
      >节点管理</el-menu-item>
      <el-menu-item
        index="3"
        @click="toLink('address-mgr')"
      >代扣地址管理</el-menu-item>
      <el-menu-item
        index="4"
        @click="toLink('service-monitor')"
      >节点业务监控</el-menu-item>
      <el-menu-item
        index="5"
        @click="toLink('resource-monitor')"
      >节点资源监控</el-menu-item>
      <el-menu-item
        index="6"
        @click="toLink('alarm-history')"
      >历史告警信息</el-menu-item>
    </el-menu>
  </div>
</template>

<script>
export default {
  data() {
    return {};
  },
  methods: {
    toLink(url) {
      this.$router.push({ name: url });
    }
  },
  mounted() {}
};
</script>

<style>
.sidebar_container {
  width: 100%;
}
.sidebar_container .el-menu {
  border-right: solid 0px #e6e6e6 !important;
}
.sidebar_container .el-submenu .title-icon {
  font-size: 18px;
  width: 18px;
  margin-right: 10px;
  text-align: center;
}
.sidebar_container .el-menu-item .item-icon {
  font-size: 18px;
  width: 18px;
  margin-right: 15px;
  text-align: center;
}
.sidebar_container .el-submenu .item-icon {
  font-size: 10px;
  width: 10px;
  margin-right: 10px;
  text-align: center;
}
.sidebar_container > ul{
  margin-top: 5px;
  background-color: #020F34 !important;
}
.sidebar_container > ul > li {
  background-color: #020F34 !important;
  font-size: 18px;
  border-left: 3px solid #020F34;
  margin-bottom: 10px;
}
.sidebar_container > ul .is-active{
  background: linear-gradient(to right, #01195D, #021039);;
  border-left-color: #00FFFF;
}
</style>
