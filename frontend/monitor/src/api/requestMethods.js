import http from './fetch'
import axios from 'axios'

export const getRequestData = () => {
    let data = {
        head: {
            version: '1.0',
            uid: $cookies.get('id'),
            token: $cookies.get('f'),
            timestamp: Date.parse(new Date())
        },
        body: {}
    };
    return data;
}

// axios.interceptors.response.use(function (res) {
//     // let redirectUrl = "http://114.55.11.139:1101/#/";
//     let redirectUrl = "http://47.244.139.179/#/";

//     if (res.data.errCode == 10100) {
//         $cookies.remove('id'),
//             $cookies.remove('f'),
//             $cookies.remove('account'),
//             alert("token已失效，即将跳转至登录页~");
//         window.location.href = redirectUrl;
//     } else if (res.data.errCode == 10001) {
//         $cookies.remove('id'),
//             $cookies.remove('f'),
//             $cookies.remove('account'),
//             alert("服务错误，即将跳转至登录页~");
//         window.location.href = redirectUrl;
//     }
//     return res;
// }, function (err) {
//     return Promise.reject(err);
// })

// 管理员管理
// -登录
export const adminLogin = (data) => {
    return http.postFetch("admin/adminLogin", data);
}


// 分组管理
// -分组列表
export const groupList = (data) => {
    return http.postFetch("hostgroup/list", data);
}
// -新增分组
export const groupAdd = (data) => {
    return http.postFetch("hostgroup/add", data);
}
// -修改分组
export const groupUpdate = (data) => {
    return http.postFetch("hostgroup/update", data);
}
// -删除分组
export const groupDel = (groupId) => {
    return http.deleteFetch("hostgroup/delete/" + groupId);
}

// 节点管理
// -节点列表
export const nodeList = (data) => {
    return http.postFetch("hostmgr/list", data);
}
// -通过groupId查询所有节点
export const queryNodeListByGroupId = (groupId) => {
    return http.getFetch("hostmgr/" + groupId);
}
// -新增节点
export const nodeAdd = (data) => {
    return http.postFetch("hostmgr/add", data);
}
// -修改节点
export const nodeUpdate = (data) => {
    return http.postFetch("hostmgr/update", data);
}
// -删除节点
export const nodeDel = (nodeId) => {
    return http.deleteFetch("hostmgr/delete/" + nodeId);
}

// 代扣地址管理
// -代扣地址列表
export const addrList = (data) => {
    return http.postFetch("addressmgr/list", data);
}
// -新增代扣地址
export const addrAdd = (data) => {
    return http.postFetch("addressmgr/add", data);
}
// -修改代扣地址
export const addrUpdate = (data) => {
    return http.postFetch("addressmgr/update", data);
}
// -删除代扣地址
export const addrDel = (addrId) => {
    return http.deleteFetch("addressmgr/delete/" + addrId);
}

// 节点业务监控
// -统计业务告警总数
export const servAlarmCount = (groupId) => {
    return http.getFetch("monitormgr/buswarningcount/" + groupId);
}
// -

// 节点资源监控
// -统计资源告警总数
export const resAlarmCount = (groupId) => {
    return http.getFetch("monitormgr/reswarningcount/" + groupId);
}
// -查看monitor信息
export const nodeMonitor = (groupId) => {
    return http.getFetch("monitormgr/group/" + groupId);
}

// -查看主节点monitor信息
export const mainNetMonitor = (groupId) => {
    return http.getFetch("monitormgr/group/mainnet/" + groupId);
}
// -查看hash一致性
export const hashCheck = (groupId) => {
    return http.getFetch("monitormgr/hashcheck/" + groupId);
}
// -余额
export const addrBalance = (groupId) => {
    return http.getFetch("monitormgr/balance/" + groupId);
}
// -余额变化
export const balancelist = (groupId) => {
    return http.getFetch("monitormgr/balancelist/" + groupId);
}
// -业务告警信息
export const buswarning = (groupId) => {
    return http.getFetch("monitormgr/buswarning/" + groupId);
}
// -资源告警信息
export const reswarning = (hostId) => {
    return http.getFetch("monitormgr/reswarning/" + hostId);
}
// -移除警告
export const warningRemove = (data) => {
    return http.postFetch("monitormgr/warning/remove", data);
}
// -批量移除警告
export const warningBatchRemove = (data) => {
    return http.postFetch("monitormgr/warning/batchremove", data);
}
// -查看节点资源信息
export const resourceList = (hostId) => {
    return http.getFetch("monitormgr/resourcelist/" + hostId);
}
// 历史告警信息
export const alarmHistory = (data) => {
    return http.postFetch("monitormgr/warning/history", data);
}
// 告警信息数
export const alarmCount = () => {
    return http.getFetch("monitormgr/warning");
}
