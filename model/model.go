package model

//group
type GroupInfo struct {
	GroupID   int64  `json:"groupId"`
	GroupName string `json:"groupName"`
	Describe  string `json:"describe"`
	Title     string `json:"title"`
}

//paymentAddress
type PaymentAddress struct {
	ID      int64  `json:"id"`
	GroupID int64  `json:"groupId"`
	Address string `json:"address"`
}

//hostInfo
type HostInfo struct {
	HostID          int64  `json:"hostId"`
	GroupID         int64  `json:"groupId"`
	HostIp          string `json:"hostIp"`
	SSHPort         int64  `json:"SSHPort"`
	UserName        string `json:"userName"`
	PassWd          string `json:"passWd"`
	IsCheckResource int64  `json:"isCheckResource"`
	ProcessName     string `json:"processName"`
	ServerPort      int64  `json:"serverPort"`
	CreateTime      int64  `json:"createTime"`
	UpdateTime      int64  `json:"updateTime"`
}

//resourceInfo
type ResourceInfo struct {
	ID         int64
	HostID     int64
	GroupID    int64
	Mem        int64
	CPU        int64
	Disk       int64
	CreateTime int64
}

//ServerMonitor
type Monitor struct {
	ID              int64
	HostID          int64
	GroupID         int64
	HostIp          string
	ServerPort      int64
	LastBlockHeight int64
	IsSync          int64
	LastBlockHash   string
	UpdateTime      int64
}

//balance
type Balance struct {
	ID         int64
	GroupID    int64
	Address    string
	Balance    int64
	CreateTime int64
}

//warning
type Warning struct {
	ID          int64
	HostID      int64
	GroupID     int64
	Type        int64
	Warning     string
	BlockHeight int64
	CreateTime  int64
	IsClosed    int64
	UpdateTime  int64
}
