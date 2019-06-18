import Vue from 'vue'
import Router from 'vue-router'
import store from '../store'
import { getRequestData, getAdminPermission } from "@/api/requestMethods";

Vue.use(Router)

const MONITOR = {
  path: '/',
  name: "index",
  component: () => import('../pages/layout/Layout.vue'),
  children: [
    {
      path: '/',
      name: 'service-monitor',
      component: () => import('../pages/monitor/ServiceMonitor.vue')
    },
    {
      path: '/resourceMonitor',
      name: 'resource-monitor',
      component: () => import('../pages/monitor/ResourceMonitor.vue')
    },
    {
      path: '/alarmHistory',
      name: 'alarm-history',
      component: () => import('../pages/monitor/AlarmHistory.vue')
    },
    {
      path: '/groupMgr',
      name: 'group-mgr',
      component: () => import('../pages/settings/GroupMgr.vue')
    },
    {
      path: '/nodeMgr',
      name: 'node-mgr',
      component: () => import('../pages/settings/NodeMgr.vue')
    },
    {
      path: '/addressMgr',
      name: 'address-mgr',
      component: () => import('../pages/settings/AddressMgr.vue')
    }
  ]
}

const router = new Router({
  routes: [
    MONITOR,
    {
      name: '404',
      path: '/404',
      component: () => import('../components/NotFound.vue')
    },
    { path: '*', redirect: '/404' }
  ]
});

export default router;