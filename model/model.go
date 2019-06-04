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
	HostName        string `json:"hostName"`
	GroupID         int64  `json:"groupId"`
	GroupName       string `json:"groupName"`
	HostIp          string `json:"hostIp"`
	SSHPort         int64  `json:"SSHPort"`
	UserName        string `json:"userName"`
	PassWd          string `json:"passWd"`
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
	ID              int64  `json:"id"`
	HostID          int64  `json:"hostId"`
	GroupID         int64  `json:"groupId"`
	HostIp          string `json:"hostIp"`
	ServerPort      int64  `json:"serverPort"`
	LastBlockHeight int64  `json:"lastBlockHeight"`
	IsSync          int64  `json:"isSync"`
	LastBlockHash   string `json:"lastBlockHash"`
	UpdateTime      int64  `json:"updateTime"`
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
type Page struct {
	PageNum  int64 `json:"pageNum"`
	PageSize int64 `json:"pageSize"`
}

// 创建到此SSH主机的连接
//func (node *HostInfo) Connect() (*ssh.Client, error) {
//	// 使用密码认证
//	config := &ssh.ClientConfig{
//		User: node.UserName,
//		Auth: []ssh.AuthMethod{
//			ssh.Password(node.PassWd),
//		},
//		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
//	}
//	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", node.HostIp, node.SSHPort), config)
//	if err != nil {
//		log.Error("connecting error", err)
//		return nil, err
//	}
//	return client, nil
//}
