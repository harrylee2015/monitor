package common

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/harrylee2015/monitor/model"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

var (
	InsertGroupInfoSql      = "INSERT INTO GroupInfo (groupName, describle, title) VALUES (?,?,?)"
	InsertPaymentAddressSql = "INSERT INTO PaymentAddress (groupId, address) VALUES (?,?)"
	InsertHostInfoSql       = "INSERT INTO HostInfo (groupId, hostIp, sshPort,userName,passWd,isCheckResource,processName,serverPort,createTime,updateTime) VALUES (?,?,?,?,?,?,?,?,?,?)"
	InsertResourceInfoSql   = "INSERT INTO ResourceInfo (groupId, hostId, mem,cpu,disk,createTime) VALUES (?,?,?,?,?,?)"
	InsertMonitorSql        = "INSERT INTO Monitor (groupId, hostId, hostIp,serverPort,lastBlockHeight,isSync,lastBlockHash,updateTime) VALUES (?,?,?,?,?,?,?,?)"
	InsertBalanceSql        = "INSERT INTO Balance (groupId, address, balance,createTime) VALUES (?,?,?,?)"
	InsertWarningSql        = "INSERT INTO Warning (groupId, hostId, type,warning,blockHeight,createTime,isClosed,updateTime) VALUES (?,?,?,?,?,?,?,?)"

	QueryGroupInfoByPageNum      = "SELECT * FROM GroupInfo LIMIT ? OFFSET ?"
	QueryHostInfoByPageNum       = "SELECT * FROM HostInfo LIMIT ? OFFSET ?"
	QueryHostInfoByGroupId       = "SELECT * FROM HostInfo WHERE groupId=? ORDER BY hostId ASC"
	QueryPaymentAddressByPageNum = "SELECT * FROM PaymentAddress LIMIT ? OFFSET ?"
	QueryPaymentAddress          = "SELECT * FROM PaymentAddress WHERE groupId=?"
	QueryResourceInfo            = "SELECT * FROM ResourceInfo WHERE hostId=? ORDER BY createTime DESC limit ?"
	QueryMonitorSql              = "SELECT * FROM Monitor WHERE groupId=? ORDER BY hostId ASC"
	QueryBalanceSql              = "SELECT * FROM Balance WHERE groupId=? ORDER BY createTime ASC"
	QueryBusWarningCount         = "SELECT COUNT(*) FROM Warning WHERE isClosed=0  AND groupId=? AND type IN ( 4, 5 )"
	QueryResWarningCount         = "SELECT COUNT(*) FROM Warning WHERE isClosed=0  AND groupId=? AND type IN ( 1, 2, 3 )"
	QueryBusWarningByGroupId     = "SELECT * FROM Warning WHERE isClosed=0 AND groupId=? AND type IN ( 4, 5 ) ORDER BY createTime ASC"
	QueryResWarningByHostId      = "SELECT * FROM Warning WHERE isClosed=0 AND hostId=? AND type IN ( 1, 2, 3 ) ORDER BY createTime ASC"
	QueryHistoryWarning          = "SELECT * FROM Warning WHERE isClosed=1 ORDER BY createTime DESC LIMIT ? OFFSET ? "

	UpdateGroupInfoSql      = "UPDATE GroupInfo SET groupName=?,describle=?,title=? WHERE groupId=?"
	UpdatePaymentAddressSql = "UPDATE PaymentAddress SET groupId=?,address=? WHERE id=?"
	UpdateHostInfoSql       = "UPDATE HostInfo SET hostIp=?,sshPort=?,userName=?,passWd=?,isCheckResource=?,processName=?,serverPort=?,updateTime=? WHERE hostId=?"
	UpdateMonitorSql        = "UPDATE Monitor SET hostIp=?,serverPort=?,lastBlockHeight=?,isSync=?,lastBlockHash=?,updateTime=? WHERE hostId=?"
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
)

type MonitorDB struct {
	sync.Mutex
	db *sql.DB
}

func NewMonitorDB() *MonitorDB {
	monitorDb := &MonitorDB{}
	db, err := sql.Open("sqlite3", fmt.Sprintf("%v/%v", model.DbPath, "monitor.db"))
	checkErr(err)
	monitorDb.db = db
	return monitorDb
}

func (mdb *MonitorDB) Close() {
	mdb.db.Close()
}
func checkErr(err error) {
	if err != nil {
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
		mdb.insertData(InsertResourceInfoSql, g.GroupID, g.HostID, g.Mem, g.CPU, g.Disk, time.Now().Unix())
		return
	}
	if g, ok := data.(*model.Monitor); ok {
		mdb.insertData(InsertMonitorSql, g.GroupID, g.HostID, g.HostIp, g.ServerPort, g.LastBlockHeight, g.IsSync, g.LastBlockHash, time.Now().Unix())
		return
	}
	if g, ok := data.(*model.Balance); ok {
		mdb.insertData(InsertBalanceSql, g.GroupID, g.Address, g.Balance, time.Now().Unix())
		return
	}
	if g, ok := data.(*model.Warning); ok {
		mdb.insertData(InsertWarningSql, g.GroupID, g.HostID, g.Type, g.Warning, g.BlockHeight, time.Now().Unix(), g.IsClosed, time.Now().Unix())
		return
	}
	if g, ok := data.(*model.GroupInfo); ok {
		mdb.insertData(InsertGroupInfoSql, g.GroupName, g.Describe, g.Title)
		return
	}
	if p, ok := data.(*model.PaymentAddress); ok {
		mdb.insertData(InsertPaymentAddressSql, p.GroupID, p.Address)
		return
	}
	if g, ok := data.(*model.HostInfo); ok {
		mdb.insertData(InsertHostInfoSql, g.GroupID, g.HostIp, g.SSHPort, g.UserName, g.PassWd, g.IsCheckResource, g.ProcessName, g.ServerPort, time.Now().Unix(), time.Now().Unix())
		return
	}
	panic(fmt.Errorf("unknow type data!,data=%v", data))
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
	defer stmt.Close()
	defer rows.Close()
	for rows.Next() {
		value := model.GroupInfo{}
		//TODO:这里应该是字段一一对应关系
		err := rows.Scan(&value.GroupID, &value.GroupName, &value.Describe, &value.Title)
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
	defer stmt.Close()
	defer rows.Close()
	for rows.Next() {
		value := model.HostInfo{}
		//TODO:这里应该是字段一一对应关系
		err := rows.Scan(&value.HostID, &value.GroupID, &value.HostIp, &value.SSHPort, &value.UserName, &value.PassWd, &value.IsCheckResource, &value.ProcessName, &value.ServerPort, &value.CreateTime, &value.UpdateTime)
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
	defer stmt.Close()
	defer rows.Close()
	for rows.Next() {
		value := model.HostInfo{}
		//TODO:这里应该是字段一一对应关系
		err := rows.Scan(&value.HostID, &value.GroupID, &value.HostIp, &value.SSHPort, &value.UserName, &value.PassWd, &value.IsCheckResource, &value.ProcessName, &value.ServerPort, &value.CreateTime, &value.UpdateTime)
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
	defer stmt.Close()
	defer rows.Close()
	for rows.Next() {
		value := model.PaymentAddress{}
		//TODO:这里应该是字段一一对应关系
		err := rows.Scan(&value.ID, &value.GroupID, &value.Address)
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
	defer stmt.Close()
	defer rows.Close()
	for rows.Next() {
		value := model.PaymentAddress{}
		//TODO:这里应该是字段一一对应关系
		err := rows.Scan(&value.ID, &value.GroupID, &value.Address)
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
	defer stmt.Close()
	defer rows.Close()
	for rows.Next() {
		value := model.ResourceInfo{}
		//TODO:这里应该是字段一一对应关系
		err := rows.Scan(&value.ID, &value.GroupID, &value.HostID, &value.Mem, &value.CPU, &value.Disk, &value.CreateTime)
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
	defer stmt.Close()
	defer rows.Close()
	for rows.Next() {
		value := model.Monitor{}
		//TODO:这里应该是字段一一对应关系
		err := rows.Scan(&value.ID, &value.GroupID, &value.HostID, &value.HostIp, &value.ServerPort, &value.LastBlockHeight, &value.IsSync, &value.LastBlockHash, &value.UpdateTime)
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
	defer stmt.Close()
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
	defer stmt.Close()
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
	defer stmt.Close()
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
func (mdb *MonitorDB) QueryHistoryWarning(page *model.Page) (items []*model.Warning) {
	mdb.Lock()
	defer mdb.Unlock()
	stmt, err := mdb.db.Prepare(QueryHistoryWarning)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(page.PageSize, (page.PageNum-1)*page.PageSize)
	checkErr(err)
	defer stmt.Close()
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

func (mdb *MonitorDB) QueryBusWarningCount(groupId int64) (count int64) {
	mdb.Lock()
	defer mdb.Unlock()
	stmt, err := mdb.db.Prepare(QueryBusWarningCount)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(groupId)
	checkErr(err)
	defer stmt.Close()
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
	defer stmt.Close()
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
	affect, err := res.RowsAffected()
	checkErr(err)
	if affect != 1 {
		panic(fmt.Errorf("update data error! sql=%v", sql))
	}
}

// UpdateData
func (mdb *MonitorDB) UpdateData(data interface{}) {
	mdb.Lock()
	defer mdb.Unlock()
	if g, ok := data.(*model.GroupInfo); ok {
		mdb.updateData(UpdateGroupInfoSql, g.GroupName, g.Describe, g.Title, g.GroupID)
		return
	}
	if g, ok := data.(*model.PaymentAddress); ok {
		mdb.updateData(UpdatePaymentAddressSql, g.GroupID, g.Address, g.ID)
		return
	}
	if g, ok := data.(*model.HostInfo); ok {
		mdb.updateData(UpdateHostInfoSql, g.HostID, g.SSHPort, g.UserName, g.PassWd, g.IsCheckResource, g.ProcessName, g.ServerPort, time.Now().Unix(), g.HostID)
		return
	}
	//UpdateMonitorSql
	if g, ok := data.(*model.Monitor); ok {
		mdb.updateData(UpdateMonitorSql, g.HostIp, g.ServerPort, g.LastBlockHeight, g.IsSync, g.LastBlockHash, time.Now().Unix(), g.HostID)
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
	affect, err := res.RowsAffected()
	checkErr(err)
	if affect != 1 {
		panic("Delete data error!")
	}
	stmt, err = tx.Prepare(DelPaymentAddressByGroupId)
	checkErr(err)
	res, err = stmt.Exec(groupId)
	checkErr(err)
	affect, err = res.RowsAffected()
	checkErr(err)
	if affect != 1 {
		panic("Delete data error!")
	}
	stmt, err = tx.Prepare(DelHostInfoByGroupId)
	checkErr(err)
	res, err = stmt.Exec(groupId)
	checkErr(err)
	affect, err = res.RowsAffected()
	checkErr(err)
	if affect != 1 {
		panic("Delete data error!")
	}
	stmt, err = tx.Prepare(DelResourceInfoByGroupId)
	checkErr(err)
	res, err = stmt.Exec(groupId)
	checkErr(err)
	affect, err = res.RowsAffected()
	checkErr(err)
	if affect != 1 {
		panic("Delete data error!")
	}
	stmt, err = tx.Prepare(DelMonitorByGroupId)
	checkErr(err)
	res, err = stmt.Exec(groupId)
	checkErr(err)
	affect, err = res.RowsAffected()
	checkErr(err)
	if affect != 1 {
		panic("Delete data error!")
	}
	//DelBalanceByGroupId
	stmt, err = tx.Prepare(DelBalanceByGroupId)
	checkErr(err)
	res, err = stmt.Exec(groupId)
	checkErr(err)
	affect, err = res.RowsAffected()
	checkErr(err)
	if affect != 1 {
		panic("Delete data error!")
	}
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
	affect, err := res.RowsAffected()
	checkErr(err)
	if affect != 1 {
		panic("Delete data error!")
	}
	stmt, err = tx.Prepare(DelResourceInfoByHostId)
	checkErr(err)
	res, err = stmt.Exec(hostId)
	checkErr(err)
	affect, err = res.RowsAffected()
	checkErr(err)
	if affect != 1 {
		panic("Delete data error!")
	}
	stmt, err = tx.Prepare(DelMonitorByHostId)
	checkErr(err)
	res, err = stmt.Exec(hostId)
	checkErr(err)
	affect, err = res.RowsAffected()
	checkErr(err)
	if affect != 1 {
		panic("Delete data error!")
	}
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
	affect, err := res.RowsAffected()
	checkErr(err)
	if affect != 1 {
		panic("Delete data error!")
	}
	//DelBalanceByGroupId
	stmt, err = tx.Prepare(DelBalanceByGroupId)
	checkErr(err)
	res, err = stmt.Exec(groupId)
	checkErr(err)
	affect, err = res.RowsAffected()
	checkErr(err)
	if affect != 1 {
		panic("Delete data error!")
	}
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
	affect, err := res.RowsAffected()
	checkErr(err)
	if affect != 1 {
		panic("Delete data error!")
	}
}

func (mdb *MonitorDB) DelResourceInfoByTime(time int64) {
	mdb.Lock()
	defer mdb.Unlock()
	stmt, err := mdb.db.Prepare(DelResourceInfoByTime)
	checkErr(err)
	res, err := stmt.Exec(time)
	checkErr(err)
	affect, err := res.RowsAffected()
	checkErr(err)
	if affect != 1 {
		panic("Delete data error!")
	}
}
