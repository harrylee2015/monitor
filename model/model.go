package model

import (
	"database/sql"
)

//group
type GroupInfo struct {
	GroupID   int64  `json:"groupId"`
	GroupName string `json:"groupName"`
	Describe  string `json:"describe"`
	Title     string `json:"title"`
}

//paymentAddress
type PaymentAddress struct {
	ID        int64  `json:"id"`
	GroupID   int64  `json:"groupId"`
	GroupName string `json:"groupName"`
	Address   string `json:"address"`
}

//hostInfo
type HostInfo struct {
	HostID    int64  `json:"hostId"`
	HostName  string `json:"hostName"`
	GroupID   int64  `json:"groupId"`
	GroupName string `json:"groupName"`
	HostIp    string `json:"hostIp"`
	SSHPort   int64  `json:"SSHPort"`
	UserName  string `json:"userName"`
	PassWd    string `json:"passWd"`
	// 0表示不检查节点资源状况 1表示检查
	IsCheckResource int64  `json:"isCheckResource"`
	ProcessName     string `json:"processName"`
	ServerPort      int64  `json:"serverPort"`
	CreateTime      int64  `json:"createTime"`
	UpdateTime      int64  `json:"updateTime"`
}

//专门用于返回list数据类型
type List struct {
	Total  int64       `json:"total"`
	Values interface{} `json:"values"`
}

//resourceInfo
type ResourceInfo struct {
	ID      int64 `json:"id"`
	HostID  int64 `json:"hostId"`
	GroupID int64 `json:"groupId"`
	//单位MB
	MemTotal       uint64  `json:"memTotal"`
	MemUsedPercent float64 `json:"memUsedPercent"`
	//单位核
	CpuTotal       uint64  `json:"cpuTotal"`
	CpuUsedPercent float64 `json:"cpuUsedPercent"`
	//空间大小单位GB
	DiskTotal       uint64  `json:"diskTotal"`
	DiskUsedPercent float64 `json:"diskUsedPercent"`
	CreateTime      int64   `json:"createTime"`
}

//ServerMonitor
type Monitor struct {
	ID         int64  `json:"id"`
	HostID     int64  `json:"hostId"`
	GroupID    int64  `json:"groupId"`
	HostIp     string `json:"hostIp"`
	ServerPort int64  `json:"serverPort"`
	//0表示服务正常，1表示服务异常
	ServerStatus    int64 `json:"serverStatus"`
	LastBlockHeight int64 `json:"lastBlockHeight"`
	//0表示同步，1表示未同步
	IsSync        int64  `json:"isSync"`
	LastBlockHash string `json:"lastBlockHash"`
	UpdateTime    int64  `json:"updateTime"`
}

//balance
type Balance struct {
	ID         int64  `json:"id"`
	GroupID    int64  `json:"groupId"`
	Address    string `json:"address"`
	Balance    int64  `json:"balance"`
	CreateTime int64  `json:"createTime"`
}

//warning
type Warning struct {
	ID      int64 `json:"id"`
	HostID  int64 `json:"hostId"`
	GroupID int64 `json:"groupId"`
	//告警类型 1是内存告警，2是cpu超负荷，3是磁盘不足， 4表示代扣地址余额告警，5表示区块不一致告警
	Type        int64  `json:"type"`
	Warning     string `json:"warning"`
	BlockHeight int64  `json:"blockHeight"`
	CreateTime  int64  `json:"createTime"`
	//是否处理过了 0表示没有处理，1表示处理过了
	IsClosed   int64 `json:"isClosed"`
	UpdateTime int64 `json:"updateTime"`
}

type WarningDetail struct {
	ID        int64          `json:"id"`
	HostID    sql.NullInt64  `json:"hostId"`
	HostIp    sql.NullString `json:"hostIp"`
	HostName  sql.NullString `json:"hostName"`
	GroupID   int64          `json:"groupId"`
	GroupName string         `json:"groupName"`
	//告警类型 1是内存告警，2是cpu超负荷，3是磁盘不足， 4表示代扣地址余额告警，5表示区块不一致告警
	Type        int64  `json:"type"`
	Warning     string `json:"warning"`
	BlockHeight int64  `json:"blockHeight"`
	CreateTime  int64  `json:"createTime"`
	//是否处理过了 0表示没有处理，1表示处理过了
	IsClosed   int64 `json:"isClosed"`
	UpdateTime int64 `json:"updateTime"`
}

//HASH 一致性告警
type Hash struct {
	//是否一致，如果是false的话，会把hash不一致得告警信息附上
	IsConsistent bool        `json:"isConsistent"`
	Values       interface{} `json:"values"`
}
type Page struct {
	PageNum  int64 `json:"pageNum"`
	PageSize int64 `json:"pageSize"`
}
