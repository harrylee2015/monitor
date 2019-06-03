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
	ID      int64
	HostID  int64
	GroupID int64
	//单位MB
	MemTotal       uint64
	MemUsedPercent float64
	//单位核
	CpuTotal       uint64
	CpuUsedPercent float64
	//空间大小单位GB
	DiskTotal       uint64
	DiskUsedPercent float64
	CreateTime      int64
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
	ID      int64
	HostID  int64
	GroupID int64
	//告警类型 1是内存告警，2是cpu超负荷，3是磁盘不足， 4表示代扣地址余额告警，5表示区块不一致告警
	Type        int64
	Warning     string
	BlockHeight int64
	CreateTime  int64
	//是否处理过了 0表示没有处理，1表示处理过了
	IsClosed   int64
	UpdateTime int64
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
