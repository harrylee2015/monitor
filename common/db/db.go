package db

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/harrylee2015/monitor/conf"

	"time"

	"github.com/harrylee2015/monitor/model"
	_ "github.com/mattn/go-sqlite3"
)

const (
	WARNING_TYPE = int64(iota)
	MEM_WARNING
	CPU_WARING
	DISK_WARING
	BALANCE_WARING
	HASH_WARING
	M_HASH_WARING
)
const (
	Type_Group = iota + 1
	Type_Host
	Type_Addr
)

var (
	InsertGroupInfoSql      = "INSERT INTO GroupInfo (groupName, describe, title,email) VALUES (?,?,?,?)"
	InsertPaymentAddressSql = "INSERT INTO PaymentAddress (groupId, groupName, address) VALUES (?,?,?)"
	InsertHostInfoSql       = "INSERT INTO HostInfo (hostName,groupId,groupName, hostIp, sshPort,userName,passWd,isCheckResource,processName,serverPort,mainNet,netPort,createTime,updateTime) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	InsertResourceInfoSql   = "INSERT INTO ResourceInfo (groupId, hostId, memTotal,memUsedPercent,cpuTotal,cpuUsedPercent,diskTotal,diskUsedPercent,createTime) VALUES (?,?,?,?,?,?,?,?,?)"
	InsertMonitorSql        = "INSERT INTO Monitor (groupId, hostId, hostIp,serverPort,serverStatus,lastBlockHeight,isSync,lastBlockHash,updateTime) VALUES (?,?,?,?,?,?,?,?,?)"
	InsertMainNetMonitorSql = "INSERT INTO MainNetMonitor (groupId, hostId, hostIp,serverPort,serverStatus,lastBlockHeight,isSync,lastBlockHash,updateTime) VALUES (?,?,?,?,?,?,?,?,?)"
	InsertBalanceSql        = "INSERT INTO Balance (groupId, address, balance,createTime) VALUES (?,?,?,?)"
	InsertWarningSql        = "INSERT INTO Warning (hostId,groupId,type,warning,blockHeight,createTime,isClosed,updateTime) VALUES (?,?,?,?,?,?,?,?)"

	QueryGroupInfoCount           = "SELECT COUNT(*) FROM GroupInfo"
	QueryHostInfoCount            = "SELECT COUNT(*) FROM HostInfo"
	QueryGroupInfoByPageNum       = "SELECT * FROM GroupInfo LIMIT ? OFFSET ?"
	QueryHostInfoByPageNum        = "SELECT HostInfo.hostId, HostInfo.hostName,HostInfo.groupId,GroupInfo.groupName, HostInfo.hostIp, HostInfo.sshPort,HostInfo.userName,HostInfo.passWd,HostInfo.isCheckResource,HostInfo.processName,HostInfo.serverPort,HostInfo.mainNet,HostInfo.netPort,HostInfo.createTime,HostInfo.updateTime,GroupInfo.email FROM HostInfo LEFT OUTER JOIN GroupInfo ON HostInfo.groupId = GroupInfo.groupId LIMIT ? OFFSET ?"
	QueryHostInfoByGroupId        = "SELECT * FROM HostInfo WHERE groupId=? ORDER BY hostId ASC"
	QueryHostInfoList             = "SELECT HostInfo.hostId, HostInfo.hostName,HostInfo.groupId,GroupInfo.groupName, HostInfo.hostIp, HostInfo.sshPort,HostInfo.userName,HostInfo.passWd,HostInfo.isCheckResource,HostInfo.processName,HostInfo.serverPort,HostInfo.mainNet,HostInfo.netPort,HostInfo.createTime,HostInfo.updateTime,GroupInfo.email FROM HostInfo LEFT OUTER JOIN GroupInfo ON HostInfo.groupId = GroupInfo.groupId LIMIT ? OFFSET ?"
	QueryPaymentAddressCount      = "SELECT COUNT(*) FROM PaymentAddress"
	QueryPaymentAddressByPageNum  = "SELECT PaymentAddress.id,PaymentAddress.groupId, GroupInfo.groupName, PaymentAddress.address,GroupInfo.email FROM PaymentAddress LEFT OUTER JOIN GroupInfo ON PaymentAddress.groupId = GroupInfo.groupId LIMIT ? OFFSET ?"
	QueryPaymentAddress           = "SELECT * FROM PaymentAddress WHERE groupId=?"
	QueryPaymentAddressList       = "SELECT PaymentAddress.id,PaymentAddress.groupId, GroupInfo.groupName, PaymentAddress.address,GroupInfo.email FROM PaymentAddress LEFT OUTER JOIN GroupInfo ON PaymentAddress.groupId = GroupInfo.groupId LIMIT ? OFFSET ?"
	QueryResourceInfo             = "SELECT * FROM ResourceInfo WHERE hostId=? ORDER BY createTime DESC limit ?"
	QueryMonitorSql               = "SELECT * FROM Monitor WHERE groupId=? ORDER BY hostId ASC"
	QueryMainNetMonitorSql        = "SELECT * FROM MainNetMonitor WHERE groupId=? ORDER BY hostId ASC"
	QueryMonitorById              = "SELECT * FROM Monitor WHERE groupId=? AND hostId=?"
	QueryMainNetMonitorById       = "SELECT * FROM MainNetMonitor WHERE groupId=? AND hostId=?"
	QueryLastBalance              = "SELECT * FROM Balance WHERE groupId=? ORDER BY createTime DESC LIMIT 1"
	QueryBalanceSql               = "SELECT * FROM Balance WHERE groupId=? ORDER BY createTime ASC"
	QueryBusWarningCount          = "SELECT COUNT(*) FROM Warning WHERE isClosed=0  AND groupId=? AND type IN ( 4, 5 )"
	QueryResWarningCount          = "SELECT COUNT(*) FROM Warning WHERE isClosed=0  AND groupId=? AND type IN ( 1, 2, 3 )"
	QueryBusWarningCountByGroupId = "SELECT COUNT(Warning.id),Warning.groupId,groupInfo.groupName FROM Warning LEFT OUTER JOIN GroupInfo ON Warning.groupId=GroupInfo.groupId WHERE isClosed=0 AND type IN ( 4, 5 ) GROUP BY GroupInfo.groupId"
	QueryResWarningCountByGroupId = "SELECT COUNT(Warning.id),Warning.groupId,groupInfo.groupName FROM Warning LEFT OUTER JOIN GroupInfo ON Warning.groupId=GroupInfo.groupId WHERE isClosed=0 AND type IN ( 1, 2, 3) GROUP BY GroupInfo.groupId"
	QueryBusWarningByGroupId      = "SELECT * FROM Warning WHERE isClosed=0 AND groupId=? AND type IN ( 4, 5 ) ORDER BY createTime ASC"
	QueryResWarningByHostId       = "SELECT * FROM Warning WHERE isClosed=0 AND hostId=? AND type IN ( 1, 2, 3 ) ORDER BY createTime ASC"
	QueryWarningByGroupIdAndType  = "SELECT * FROM Warning WHERE isClosed=0  AND groupId=? AND type=?"
	QueryWarningByHostIdAndType   = "SELECT * FROM Warning WHERE isClosed=0  AND hostId=? AND type=?"
	QueryHistoryWarningCount      = "SELECT COUNT(*) FROM Warning WHERE isClosed=1"
	QueryHistoryWarning           = "SELECT warning.id,Warning.hostId,HostInfo.hostIp,HostInfo.hostName,Warning.groupId,GroupInfo.groupName,Warning.type,Warning.warning,Warning.blockHeight,Warning.createTime,Warning.isClosed,Warning.updateTime FROM Warning LEFT OUTER JOIN HostInfo ON Warning.hostId = HostInfo.hostId LEFT OUTER JOIN GroupInfo ON Warning.groupId = GroupInfo.groupId WHERE Warning.isClosed=1 ORDER BY Warning.createTime DESC LIMIT ? OFFSET ?"

	UpdateGroupInfoSql      = "UPDATE GroupInfo SET groupName=?,describe=?,title=?,email=? WHERE groupId=?"
	UpdatePaymentAddressSql = "UPDATE PaymentAddress SET groupId=?,groupName=?,address=? WHERE id=?"
	UpdateHostInfoSql       = "UPDATE HostInfo SET hostName=?,hostIp=?,sshPort=?,userName=?,passWd=?,isCheckResource=?,processName=?,serverPort=?,mainNet=?,netPort=?,updateTime=? WHERE hostId=?"
	UpdateMonitorSql        = "UPDATE Monitor SET hostIp=?,serverPort=?,serverStatus=?,lastBlockHeight=?,isSync=?,lastBlockHash=?,updateTime=? WHERE groupId =? AND hostId=?"
	UpdateMainNetMonitorSql = "UPDATE MainNetMonitor SET hostIp=?,serverPort=?,serverStatus=?,lastBlockHeight=?,isSync=?,lastBlockHash=?,updateTime=? WHERE groupId =? AND hostId=?"
	UpdateWarningSql        = "UPDATE Warning SET isClosed=?,updateTime=? WHERE id=?"

	DelGroupInfoSql            = "DELETE FROM GroupInfo WHERE groupId=?"
	DelPaymentAddressSql       = "DELETE FROM PaymentAddress WHERE id=?"
	DelPaymentAddressByGroupId = "DELETE FROM PaymentAddress WHERE groupId=?"
	DelHostInfoSql             = "DELETE FROM HostInfo WHERE hostId=?"
	DelHostInfoByGroupId       = "DELETE FROM HostInfo WHERE groupId=?"
	DelResourceInfoByHostId    = "DELETE FROM ResourceInfo WHERE hostId=?"
	DelResourceInfoByGroupId   = "DELETE FROM ResourceInfo WHERE groupId=?"
	DelResourceInfoByTime      = "DELETE FROM ResourceInfo WHERE createTime<=?"
	DelMonitorByHostId         = "DELETE FROM Monitor WHERE hostId=?"
	DelMonitorByGroupId        = "DELETE FROM Monitor WHERE groupId=?"
	DelBalanceByTime           = "DELETE FROM Balance WHERE createTime<=?"
	DelBalanceByGroupId        = "DELETE FROM Balance WHERE groupId=?"

	//DbPath = "datadir"
)

type MonitorDB struct {
	sync.Mutex
	db *sql.DB
}

func NewMonitorDB() *MonitorDB {
	monitorDb := &MonitorDB{}
	db, err := sql.Open("sqlite3", fmt.Sprintf("%v/%v", conf.DbPath, "monitor.db"))
	checkErr(err)
	monitorDb.db = db
	return monitorDb
}

func (mdb *MonitorDB) Close() {
	mdb.db.Close()
}
func checkErr(err error) {
	if err != nil {
		fmt.Println("err:", err.Error())
		panic(err)
	}
}
func (mdb *MonitorDB) insertData(sql string, args ...interface{}) {
	stmt, err := mdb.db.Prepare(sql)
	checkErr(err)
	defer stmt.Close()
	res, err := stmt.Exec(args...)
	checkErr(err)
	id, err := res.LastInsertId()
	checkErr(err)
	if id < 0 {
		panic(fmt.Errorf("insertData error! sql=%v", sql))
	}
}

// insertData
func (mdb *MonitorDB) InsertData(data interface{}) {
	mdb.Lock()
	defer mdb.Unlock()
	if g, ok := data.(*model.ResourceInfo); ok {
		mdb.insertData(InsertResourceInfoSql, g.GroupID, g.HostID, g.MemTotal, g.MemUsedPercent, g.CpuTotal, g.CpuUsedPercent, g.DiskTotal, g.DiskUsedPercent, time.Now().Unix())
		return
	}
	if g, ok := data.(*model.Monitor); ok {
		mdb.insertData(InsertMonitorSql, g.GroupID, g.HostID, g.HostIp, g.ServerPort, g.ServerStatus, g.LastBlockHeight, g.IsSync, g.LastBlockHash, time.Now().Unix())
		return
	}
	if g, ok := data.(*model.MainNetMonitor); ok {
		mdb.insertData(InsertMainNetMonitorSql, g.GroupID, g.HostID, g.HostIp, g.ServerPort, g.ServerStatus, g.LastBlockHeight, g.IsSync, g.LastBlockHash, time.Now().Unix())
		return
	}
	if g, ok := data.(*model.Balance); ok {
		mdb.insertData(InsertBalanceSql, g.GroupID, g.Address, g.Balance, time.Now().Unix())
		return
	}
	if g, ok := data.(*model.Warning); ok {
		mdb.insertData(InsertWarningSql, g.HostID, g.GroupID, g.Type, g.Warning, g.BlockHeight, time.Now().Unix(), g.IsClosed, time.Now().Unix())
		return
	}
	if g, ok := data.(*model.GroupInfo); ok {
		mdb.insertData(InsertGroupInfoSql, g.GroupName, g.Describe, g.Title, g.Email)
		return
	}
	if p, ok := data.(*model.PaymentAddress); ok {
		mdb.insertData(InsertPaymentAddressSql, p.GroupID, p.GroupName, p.Address)
		return
	}
	if g, ok := data.(*model.HostInfo); ok {
		mdb.insertData(InsertHostInfoSql, g.HostName, g.GroupID, g.GroupName, g.HostIp, g.SSHPort, g.UserName, g.PassWd, g.IsCheckResource, g.ProcessName, g.ServerPort, g.MainNet, g.NetPort, time.Now().Unix(), time.Now().Unix())
		return
	}
	panic(fmt.Errorf("unknow type data!,data=%v", data))
}

//根据类型统计count数量
func (mdb *MonitorDB) QueryCount(types int) (count int64) {
	mdb.Lock()
	defer mdb.Unlock()
	query := ""
	switch types {
	case Type_Group:
		query = QueryGroupInfoCount
	case Type_Host:
		query = QueryHostInfoCount
	case Type_Addr:
		query = QueryPaymentAddressCount
	}
	stmt, err := mdb.db.Prepare(query)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	checkErr(err)
	defer stmt.Close()
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&count)
		checkErr(err)
	}
	return count
}

func (mdb *MonitorDB) QueryBusWarningCountByGroup() (items []*model.WarningCount) {
	mdb.Lock()
	defer mdb.Unlock()
	stmt, err := mdb.db.Prepare(QueryBusWarningCountByGroupId)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		value := model.WarningCount{}
		//TODO:这里应该是字段一一对应关系
		err := rows.Scan(&value.Count, &value.GroupId, &value.GroupName)
		checkErr(err)
		items = append(items, &value)
	}
	return items
}

func (mdb *MonitorDB) QueryResWarningCountByGroup() (items []*model.WarningCount) {
	mdb.Lock()
	defer mdb.Unlock()
	stmt, err := mdb.db.Prepare(QueryResWarningCountByGroupId)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		value := model.WarningCount{}
		//TODO:这里应该是字段一一对应关系
		err := rows.Scan(&value.Count, &value.GroupId, &value.GroupName)
		checkErr(err)
		items = append(items, &value)
	}
	return items
}

// QueryGroupInfoByPageNum
func (mdb *MonitorDB) QueryGroupInfoByPageNum(page *model.Page) (items []*model.GroupInfo) {
	mdb.Lock()
	defer mdb.Unlock()
	stmt, err := mdb.db.Prepare(QueryGroupInfoByPageNum)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(page.PageSize, (page.PageNum-1)*page.PageSize)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		value := model.GroupInfo{}
		//TODO:这里应该是字段一一对应关系
		err := rows.Scan(&value.GroupID, &value.GroupName, &value.Describe, &value.Title, &value.Email)
		checkErr(err)
		items = append(items, &value)
	}
	return items
}

func (mdb *MonitorDB) QueryHostInfoByPageNum(page *model.Page) (items []*model.HostInfo) {
	mdb.Lock()
	defer mdb.Unlock()
	stmt, err := mdb.db.Prepare(QueryHostInfoByPageNum)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(page.PageSize, (page.PageNum-1)*page.PageSize)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		value := model.HostInfo{}
		//TODO:这里应该是字段一一对应关系
		err := rows.Scan(&value.HostID, &value.HostName, &value.GroupID, &value.GroupName, &value.HostIp, &value.SSHPort, &value.UserName, &value.PassWd, &value.IsCheckResource, &value.ProcessName, &value.ServerPort, &value.MainNet, &value.NetPort, &value.CreateTime, &value.UpdateTime, &value.Email)
		checkErr(err)
		items = append(items, &value)
	}
	return items
}

func (mdb *MonitorDB) QueryHostInfoByGroupId(groupId int64) (items []*model.HostInfo) {
	mdb.Lock()
	defer mdb.Unlock()
	stmt, err := mdb.db.Prepare(QueryHostInfoByGroupId)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(groupId)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		value := model.HostInfo{}
		//TODO:这里应该是字段一一对应关系
		err := rows.Scan(&value.HostID, &value.HostName, &value.GroupID, &value.GroupName, &value.HostIp, &value.SSHPort, &value.UserName, &value.PassWd, &value.IsCheckResource, &value.ProcessName, &value.ServerPort, &value.MainNet, &value.NetPort, &value.CreateTime, &value.UpdateTime)
		checkErr(err)
		items = append(items, &value)
	}
	return items
}

func (mdb *MonitorDB) QueryPaymentAddress(groupId int64) (items []*model.PaymentAddress) {
	mdb.Lock()
	defer mdb.Unlock()
	stmt, err := mdb.db.Prepare(QueryPaymentAddress)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(groupId)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		value := model.PaymentAddress{}
		//TODO:这里应该是字段一一对应关系
		err := rows.Scan(&value.ID, &value.GroupID, &value.GroupName, &value.Address)
		checkErr(err)
		items = append(items, &value)
	}
	return items
}

func (mdb *MonitorDB) QueryPaymentAddressByPageNum(page *model.Page) (items []*model.PaymentAddress) {
	mdb.Lock()
	defer mdb.Unlock()
	stmt, err := mdb.db.Prepare(QueryPaymentAddressByPageNum)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(page.PageSize, (page.PageNum-1)*page.PageSize)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		value := model.PaymentAddress{}
		//TODO:这里应该是字段一一对应关系
		err := rows.Scan(&value.ID, &value.GroupID, &value.GroupName, &value.Address, &value.Email)
		checkErr(err)
		items = append(items, &value)
	}
	return items
}

func (mdb *MonitorDB) QueryResourceInfo(hostId, limit int64) (items []*model.ResourceInfo) {
	mdb.Lock()
	defer mdb.Unlock()
	stmt, err := mdb.db.Prepare(QueryResourceInfo)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(hostId, limit)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		value := model.ResourceInfo{}
		//TODO:这里应该是字段一一对应关系
		err := rows.Scan(&value.ID, &value.GroupID, &value.HostID, &value.MemTotal, &value.MemUsedPercent, &value.CpuTotal, &value.CpuUsedPercent, &value.DiskTotal, &value.DiskUsedPercent, &value.CreateTime)
		checkErr(err)
		items = append(items, &value)
	}
	return items
}

//QueryMonitorSql
func (mdb *MonitorDB) QueryMonitor(groupId int64) (items []*model.Monitor) {
	mdb.Lock()
	defer mdb.Unlock()
	stmt, err := mdb.db.Prepare(QueryMonitorSql)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(groupId)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		value := model.Monitor{}
		//TODO:这里应该是字段一一对应关系
		err := rows.Scan(&value.ID, &value.GroupID, &value.HostID, &value.HostIp, &value.ServerPort, &value.ServerStatus, &value.LastBlockHeight, &value.IsSync, &value.LastBlockHash, &value.UpdateTime)
		checkErr(err)
		items = append(items, &value)
	}
	return items
}

func (mdb *MonitorDB) QueryMainNetMonitor(groupId int64) (items []*model.MainNetMonitor) {
	mdb.Lock()
	defer mdb.Unlock()
	stmt, err := mdb.db.Prepare(QueryMainNetMonitorSql)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(groupId)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		value := model.MainNetMonitor{}
		//TODO:这里应该是字段一一对应关系
		err := rows.Scan(&value.ID, &value.GroupID, &value.HostID, &value.HostIp, &value.ServerPort, &value.ServerStatus, &value.LastBlockHeight, &value.IsSync, &value.LastBlockHash, &value.UpdateTime)
		checkErr(err)
		items = append(items, &value)
	}
	return items
}

func (mdb *MonitorDB) QueryMonitorById(groupId, hostId int64) (items []*model.Monitor) {
	mdb.Lock()
	defer mdb.Unlock()
	stmt, err := mdb.db.Prepare(QueryMonitorById)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(groupId, hostId)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		value := model.Monitor{}
		//TODO:这里应该是字段一一对应关系
		err := rows.Scan(&value.ID, &value.GroupID, &value.HostID, &value.HostIp, &value.ServerPort, &value.ServerStatus, &value.LastBlockHeight, &value.IsSync, &value.LastBlockHash, &value.UpdateTime)
		checkErr(err)
		items = append(items, &value)
	}
	return items
}

func (mdb *MonitorDB) QueryMainNetMonitorById(groupId, hostId int64) (items []*model.MainNetMonitor) {
	mdb.Lock()
	defer mdb.Unlock()
	stmt, err := mdb.db.Prepare(QueryMainNetMonitorById)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(groupId, hostId)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		value := model.MainNetMonitor{}
		//TODO:这里应该是字段一一对应关系
		err := rows.Scan(&value.ID, &value.GroupID, &value.HostID, &value.HostIp, &value.ServerPort, &value.ServerStatus, &value.LastBlockHeight, &value.IsSync, &value.LastBlockHash, &value.UpdateTime)
		checkErr(err)
		items = append(items, &value)
	}
	return items
}

//查询最新的余额信息
func (mdb *MonitorDB) QueryLastBalance(groupId int64) (items []*model.Balance) {
	mdb.Lock()
	defer mdb.Unlock()
	stmt, err := mdb.db.Prepare(QueryLastBalance)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(groupId)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		value := model.Balance{}
		//TODO:这里应该是字段一一对应关系
		err := rows.Scan(&value.ID, &value.GroupID, &value.Address, &value.Balance, &value.CreateTime)
		checkErr(err)
		items = append(items, &value)
	}
	return items
}

//QueryBalanceSql
func (mdb *MonitorDB) QueryBalance(groupId int64) (items []*model.Balance) {
	mdb.Lock()
	defer mdb.Unlock()
	stmt, err := mdb.db.Prepare(QueryBalanceSql)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(groupId)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		value := model.Balance{}
		//TODO:这里应该是字段一一对应关系
		err := rows.Scan(&value.ID, &value.GroupID, &value.Address, &value.Balance, &value.CreateTime)
		checkErr(err)
		items = append(items, &value)
	}
	return items
}

//QueryBusWarningByGroupId
func (mdb *MonitorDB) QueryWarningByGroupId(groupId int64) (items []*model.Warning) {
	mdb.Lock()
	defer mdb.Unlock()
	stmt, err := mdb.db.Prepare(QueryBusWarningByGroupId)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(groupId)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		value := model.Warning{}
		//TODO:这里应该是字段一一对应关系
		err := rows.Scan(&value.ID, &value.HostID, &value.GroupID, &value.Type, &value.Warning, &value.BlockHeight, &value.CreateTime, &value.IsClosed, &value.UpdateTime)
		checkErr(err)
		items = append(items, &value)
	}
	return items
}

//QueryResWarningByHostId
func (mdb *MonitorDB) QueryWarningByHostId(hostId int64) (items []*model.Warning) {
	mdb.Lock()
	defer mdb.Unlock()
	stmt, err := mdb.db.Prepare(QueryResWarningByHostId)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(hostId)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		value := model.Warning{}
		//TODO:这里应该是字段一一对应关系
		err := rows.Scan(&value.ID, &value.HostID, &value.GroupID, &value.Type, &value.Warning, &value.BlockHeight, &value.CreateTime, &value.IsClosed, &value.UpdateTime)
		checkErr(err)
		items = append(items, &value)
	}
	return items
}

//QueryWarningByGroupIdAndType
func (mdb *MonitorDB) QueryWarningByGroupIdAndType(groupId, ty int64) (items []*model.Warning) {
	mdb.Lock()
	defer mdb.Unlock()
	stmt, err := mdb.db.Prepare(QueryWarningByGroupIdAndType)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(groupId, ty)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		value := model.Warning{}
		//TODO:这里应该是字段一一对应关系
		err := rows.Scan(&value.ID, &value.HostID, &value.GroupID, &value.Type, &value.Warning, &value.BlockHeight, &value.CreateTime, &value.IsClosed, &value.UpdateTime)
		checkErr(err)
		items = append(items, &value)
	}
	return items
}

//QueryWarningByHostIdAndType
func (mdb *MonitorDB) QueryWarningByHostIdAndType(hostId, ty int64) (items []*model.Warning) {
	mdb.Lock()
	defer mdb.Unlock()
	stmt, err := mdb.db.Prepare(QueryWarningByHostIdAndType)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(hostId, ty)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		value := model.Warning{}
		//TODO:这里应该是字段一一对应关系
		err := rows.Scan(&value.ID, &value.HostID, &value.GroupID, &value.Type, &value.Warning, &value.BlockHeight, &value.CreateTime, &value.IsClosed, &value.UpdateTime)
		checkErr(err)
		items = append(items, &value)
	}
	return items
}

//QueryHistoryWarning
func (mdb *MonitorDB) QueryHistoryWarning(page *model.Page) (items []*model.WarningDetail) {
	mdb.Lock()
	defer mdb.Unlock()
	stmt, err := mdb.db.Prepare(QueryHistoryWarning)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(page.PageSize, (page.PageNum-1)*page.PageSize)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		value := model.WarningDetail{}
		//TODO:这里应该是字段一一对应关系
		err := rows.Scan(&value.ID, &value.HostID, &value.HostIp, &value.HostName, &value.GroupID, &value.GroupName, &value.Type, &value.Warning, &value.BlockHeight, &value.CreateTime, &value.IsClosed, &value.UpdateTime)
		checkErr(err)
		items = append(items, &value)
	}
	return items
}

func (mdb *MonitorDB) QueryHistoryWarningCount() (count int64) {
	mdb.Lock()
	defer mdb.Unlock()
	stmt, err := mdb.db.Prepare(QueryHistoryWarningCount)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&count)
		checkErr(err)
	}
	return count
}

func (mdb *MonitorDB) QueryBusWarningCount(groupId int64) (count int64) {
	mdb.Lock()
	defer mdb.Unlock()
	stmt, err := mdb.db.Prepare(QueryBusWarningCount)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(groupId)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&count)
		checkErr(err)
	}
	return count
}

func (mdb *MonitorDB) QueryResWarningCount(groupId int64) (count int64) {
	mdb.Lock()
	defer mdb.Unlock()
	stmt, err := mdb.db.Prepare(QueryResWarningCount)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(groupId)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&count)
		checkErr(err)
	}
	return count
}

func (mdb *MonitorDB) updateData(sql string, args ...interface{}) {
	stmt, err := mdb.db.Prepare(sql)
	checkErr(err)
	defer stmt.Close()
	res, err := stmt.Exec(args...)
	checkErr(err)
	_, err = res.RowsAffected()
	checkErr(err)
	// if affect != 1 {
	// 	panic(fmt.Errorf("update data error! sql=%v", sql))
	// }
}

// UpdateData
func (mdb *MonitorDB) UpdateData(data interface{}) {
	mdb.Lock()
	defer mdb.Unlock()
	if g, ok := data.(*model.GroupInfo); ok {
		mdb.updateData(UpdateGroupInfoSql, g.GroupName, g.Describe, g.Title, g.Email, g.GroupID)
		return
	}
	if g, ok := data.(*model.PaymentAddress); ok {
		mdb.updateData(UpdatePaymentAddressSql, g.GroupID, g.GroupName, g.Address, g.ID)
		return
	}
	if g, ok := data.(*model.HostInfo); ok {
		mdb.updateData(UpdateHostInfoSql, g.HostName, g.HostIp, g.SSHPort, g.UserName, g.PassWd, g.IsCheckResource, g.ProcessName, g.ServerPort, g.MainNet, g.NetPort, time.Now().Unix(), g.HostID)
		return
	}
	//UpdateMonitorSql
	if g, ok := data.(*model.Monitor); ok {
		mdb.updateData(UpdateMonitorSql, g.HostIp, g.ServerPort, g.ServerStatus, g.LastBlockHeight, g.IsSync, g.LastBlockHash, time.Now().Unix(), g.GroupID, g.HostID)
		return
	}
	//UpdateMonitorSql
	if g, ok := data.(*model.MainNetMonitor); ok {
		mdb.updateData(UpdateMainNetMonitorSql, g.HostIp, g.ServerPort, g.ServerStatus, g.LastBlockHeight, g.IsSync, g.LastBlockHash, time.Now().Unix(), g.GroupID, g.HostID)
		return
	}
	//UpdateWarningSql
	if g, ok := data.(*model.Warning); ok {
		mdb.updateData(UpdateWarningSql, 1, time.Now().Unix(), g.ID)
		return
	}
	panic(fmt.Errorf("Unkown data type! data=%v", data))

}

func (mdb *MonitorDB) DeleteDataByGroupId(groupId int64) {
	mdb.Lock()
	defer mdb.Unlock()
	// 因为操作两个表，这是一个事务操作
	tx, err := mdb.db.Begin()
	checkErr(err)
	stmt, err := tx.Prepare(DelGroupInfoSql)
	checkErr(err)
	res, err := stmt.Exec(groupId)
	checkErr(err)
	_, err = res.RowsAffected()
	checkErr(err)

	stmt, err = tx.Prepare(DelPaymentAddressByGroupId)
	checkErr(err)
	res, err = stmt.Exec(groupId)
	checkErr(err)
	_, err = res.RowsAffected()
	checkErr(err)

	stmt, err = tx.Prepare(DelHostInfoByGroupId)
	checkErr(err)
	res, err = stmt.Exec(groupId)
	checkErr(err)
	_, err = res.RowsAffected()
	checkErr(err)

	stmt, err = tx.Prepare(DelResourceInfoByGroupId)
	checkErr(err)
	res, err = stmt.Exec(groupId)
	checkErr(err)
	_, err = res.RowsAffected()
	checkErr(err)

	stmt, err = tx.Prepare(DelMonitorByGroupId)
	checkErr(err)
	res, err = stmt.Exec(groupId)
	checkErr(err)
	_, err = res.RowsAffected()
	checkErr(err)

	//DelBalanceByGroupId
	stmt, err = tx.Prepare(DelBalanceByGroupId)
	checkErr(err)
	res, err = stmt.Exec(groupId)
	checkErr(err)
	_, err = res.RowsAffected()
	checkErr(err)

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}
func (mdb *MonitorDB) DeleteDataByHostId(hostId int64) {
	mdb.Lock()
	defer mdb.Unlock()
	// 因为操作两个表，这是一个事务操作
	tx, err := mdb.db.Begin()
	checkErr(err)
	stmt, err := tx.Prepare(DelHostInfoSql)
	checkErr(err)
	res, err := stmt.Exec(hostId)
	checkErr(err)
	_, err = res.RowsAffected()
	checkErr(err)
	stmt, err = tx.Prepare(DelResourceInfoByHostId)
	checkErr(err)
	res, err = stmt.Exec(hostId)
	checkErr(err)
	_, err = res.RowsAffected()
	checkErr(err)

	stmt, err = tx.Prepare(DelMonitorByHostId)
	checkErr(err)
	res, err = stmt.Exec(hostId)
	checkErr(err)
	_, err = res.RowsAffected()
	checkErr(err)

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}
func (mdb *MonitorDB) DeleteAddressByGroupId(groupId int64) {
	mdb.Lock()
	defer mdb.Unlock()
	// 因为操作两个表，这是一个事务操作
	tx, err := mdb.db.Begin()
	checkErr(err)
	stmt, err := tx.Prepare(DelPaymentAddressByGroupId)
	checkErr(err)
	res, err := stmt.Exec(groupId)
	checkErr(err)
	_, err = res.RowsAffected()
	checkErr(err)

	//DelBalanceByGroupId
	stmt, err = tx.Prepare(DelBalanceByGroupId)
	checkErr(err)
	res, err = stmt.Exec(groupId)
	checkErr(err)
	_, err = res.RowsAffected()
	checkErr(err)

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}

func (mdb *MonitorDB) DeleteBalanceByTime(time int64) {
	mdb.Lock()
	defer mdb.Unlock()
	stmt, err := mdb.db.Prepare(DelBalanceByTime)
	checkErr(err)
	res, err := stmt.Exec(time)
	checkErr(err)
	_, err = res.RowsAffected()
	checkErr(err)
	//if affect != 1 {
	//	panic("Delete data error!")
	//}
}

func (mdb *MonitorDB) DelResourceInfoByTime(time int64) {
	mdb.Lock()
	defer mdb.Unlock()
	stmt, err := mdb.db.Prepare(DelResourceInfoByTime)
	checkErr(err)
	res, err := stmt.Exec(time)
	checkErr(err)
	_, err = res.RowsAffected()
	checkErr(err)
	//if affect != 1 {
	//	panic("Delete data error!")
	//}
}
