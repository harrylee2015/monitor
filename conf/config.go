package conf

var (
	DbPath              = "datadir"
	CollectMonitorCycle = int64(30)

	CollectResourceCycle = int64(3600 * 6)

	CollectBalanceCycle = int64(60 * 30)

	ResourceDataHoldTime = int64(3600 * 24 * 30)

	BalanceDataHoldTime = int64(3600 * 24 * 7)

	ClearDataCycle = int64(3600)

	CheckBlockHashCycle = int64(300)

	MemUsedPercentWarning = float64(80)

	CpuUsedPercentWarning = float64(80)

	DiskUsedPercentWarning = float64(80)

	BalanceWarning = float64(20)

	MainJrpc = "http://192.168.0.194:8801"

	CurrDir = ""

	FromEmail = "harrylee2015@qq.com"
	PassWd    = "vzbvfmsipfkrhfce"
	Host      = "smtp.qq.com"
	Port      = 465

	RemoteGrpcClient = "101.37.227.226:8802,39.97.20.242:8802,47.107.15.126:8802,jiedian2.bityuan.com,cloud.bityuan.com"
)

type Config struct {
	Version  string
	Database DbConfig
	Monitor  Monitor
	Log      LogConfig
}

type DbConfig struct {
	DataSource string
	DbPath     string
}

type Monitor struct {
	MainChain string

	RemoteGrpcClient string

	CollectMonitorCycle int64

	CollectResourceCycle int64

	CollectBalanceCycle int64

	CheckBlockHashCycle int64

	ResourceDataHoldTime int64

	BalanceDataHoldTime int64

	ClearDataCycle int64
	//单位百分比
	MemUsedPercentWarning float64

	CpuUsedPercentWarning float64

	DiskUsedPercentWarning float64

	BalanceWarning float64

	FromEmail string

	PassWd string
}

type LogConfig struct {
	// 日志级别，支持debug(dbug)/info/warn/error(eror)/crit
	Loglevel        string
	LogConsoleLevel string
	// 日志文件名，可带目录，所有生成的日志文件都放到此目录下
	LogFile string
	// 单个日志文件的最大值（单位：兆）
	MaxFileSize uint32
	// 最多保存的历史日志文件个数
	MaxBackups uint32
	// 最多保存的历史日志消息（单位：天）
	MaxAge uint32
	// 日志文件名是否使用本地事件（否则使用UTC时间）
	LocalTime bool
	// 历史日志文件是否压缩（压缩格式为gz）
	Compress bool
	// 是否打印调用源文件和行号
	CallerFile bool
	// 是否打印调用方法
	CallerFunction bool
}

func SetConf(conf *Config) {
	MainJrpc = conf.Monitor.MainChain

	if conf.Monitor.RemoteGrpcClient != "" {
		RemoteGrpcClient = conf.Monitor.RemoteGrpcClient
	}
	CollectMonitorCycle = conf.Monitor.CollectMonitorCycle
	CollectResourceCycle = conf.Monitor.CollectResourceCycle

	CollectBalanceCycle = conf.Monitor.CollectBalanceCycle

	CheckBlockHashCycle = conf.Monitor.CheckBlockHashCycle

	ResourceDataHoldTime = conf.Monitor.ResourceDataHoldTime

	BalanceDataHoldTime = conf.Monitor.BalanceDataHoldTime

	ClearDataCycle = conf.Monitor.ClearDataCycle

	MemUsedPercentWarning = conf.Monitor.MemUsedPercentWarning

	CpuUsedPercentWarning = conf.Monitor.CpuUsedPercentWarning

	DiskUsedPercentWarning = conf.Monitor.DiskUsedPercentWarning

	BalanceWarning = conf.Monitor.BalanceWarning

	if conf.Monitor.FromEmail != "" {
		FromEmail = conf.Monitor.FromEmail
	}

	if conf.Monitor.PassWd != "" {
		PassWd = conf.Monitor.PassWd
	}

}
