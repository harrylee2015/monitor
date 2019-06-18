<template>
  <div class="manager_header">
    <section class="empty">
      <!-- <div class="line left-top"></div>
      <div class="line left-bottom"></div> -->
    </section>
   <section class="title" >chain33平行链运行监控系统</section>
    <!-- <section class="title" >chain33平行链运行监控系统</section> -->
    <!-- <section class="alarm">
      <el-popover placement="bottom" width="200" trigger="click" id="alarm-popover">
        <ul class="alarm_items">
          <li>分组1节点业务异常警告<font> 3</font></li>
          <li>分组1节点业务异常警告<font> 3</font></li>
          <li>分组1节点业务异常警告<font> 3</font></li>
          <li>分组1节点业务异常警告<font> 3</font></li>
        </ul>
        <div slot="reference" class="btn">
          <div>异常告警</div>
          <div class="arrow"></div>
          <el-badge :value="4"></el-badge>
        </div>
      </el-popover>
    </section> -->
  </div>
</template>

<script>
import { getRequestData, logout } from "@/api/requestMethods";
export default {
  data() {
    return {};
  },
  methods: {}
};
</script>

<style>
.manager_header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.manager_header .empty {
  width: 358px;
  height: 95px;
  position: relative;
}
.manager_header .title {
  width: 380px;
  height: 95px;
  line-height: 95px;
  text-align: center;
  color: #00ffff;
  font-size: 24px;
}
.manager_header .alarm {
  width: 365px;
  height: 95px;
  padding-top: 30px;
  padding-left: 50px;
  position: relative;
}
.manager_header .alarm .btn {
  width: 200px;
  height: 50px;
  background: url("../../static/img/alarm-btn.png") center center;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #ffffff;
  font-size: 18px;
}
.manager_header .alarm .btn .arrow {
  width: 15px;
  height: 12px;
  margin-left: 6px;
  margin-right: 10px;
  background: url("../../static/img/alarm-arrow.png") center center;
}
.el-popover {
  margin-top: 0px !important;
  margin-left: -26px;
  width: 168px !important;
  border-radius: 0;
  background: linear-gradient(to bottom, #0277b8, #0195a7);
  color: #ffffff;
  border: none;
  padding: 20px 6px;
}
.alarm_items + .popper__arrow {
  visibility: hidden;
}
.alarm_items {
}
.alarm_items > li {
  color: #ffffff;
  font-size: 12px;
  margin-bottom: 10px;
  padding: 5px;
  border: 1px solid transparent;
}
.alarm_items > li:hover {
  border: 1px solid #02426D;
  border-radius: 8px;
  cursor: pointer;
  background: #02426D;
}
.alarm_items > li:hover > font{
  color: yellow;
}
/* .manager_header .line {
  height: 1px;
  position: absolute;
}
.manager_header .left-top{
  background: linear-gradient(to left, #0699EB, #082156);
  width: 322px;
  right: 36px;
  top: 18px;
}
.manager_header .left-bottom{
  background: linear-gradient(to left, #0699EB, #082156);
  width: 357px;
  right: 1px;
  bottom: 3px;
}
.manager_header .right-top{
  background: linear-gradient(to right, #0699EB, #082156);
  width: 329px;
  left: 36px;
  top: 18px;
}
.manager_header .right-bottom{
  background: linear-gradient(to right, #0699EB, #082156);
  width: 364px;
  left: 1px;
  bottom: 3px;
} */
</style>
