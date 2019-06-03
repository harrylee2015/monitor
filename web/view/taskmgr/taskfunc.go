package taskmgr

import (
	"fmt"
	DB "github.com/harrylee2015/monitor/common/db"
	"github.com/harrylee2015/monitor/common/exec"
	"github.com/harrylee2015/monitor/common/rpc"
	"github.com/harrylee2015/monitor/conf"
	"github.com/harrylee2015/monitor/model"
	"github.com/inconshreveable/log15"
	"time"
)

func getJrpc(ip string, port int64) string {
	return fmt.Sprintf("http://%s:%v", ip, port)
}

//type IsTrue = func(flag  bool) int64;

func collectMonitorData(db *DB.MonitorDB) {
	page := &model.Page{
		PageSize: 10,
		PageNum:  1,
	}
	for {
		items := db.QueryHostInfoByPageNum(page)
		if len(items) == 0 {
			return
		}

		for _, item := range items {
			jrpc := getJrpc(item.HostIp, item.ServerPort)
			lastHeader, err := rpc.QueryLastHeader(jrpc)
			if err != nil {
				//TODO  发生错误时，需要更新服务状态
				continue
			}
			isSync, err := rpc.QueryIsSync(jrpc)
			monitor := &model.Monitor{
				HostIp:          item.HostIp,
				HostID:          item.HostID,
				GroupID:         item.GroupID,
				ServerPort:      item.ServerPort,
				LastBlockHeight: lastHeader.Height,
				LastBlockHash:   lastHeader.Hash,
				IsSync: func(flag bool) int64 {
					if flag {
						return 0
					} else {
						return 1
					}
				}(isSync),
			}
			values := db.QueryMonitorById(item.GroupID, item.HostID)
			if len(values) == 0 {
				db.InsertData(monitor)

			} else {
				db.UpdateData(monitor)
			}

		}

		if len(items) < 10 {
			return
		}
		page.PageNum++
	}
}

func collectBalanceData(db *DB.MonitorDB) {
	page := &model.Page{
		PageSize: 10,
		PageNum:  1,
	}
	for {
		items := db.QueryPaymentAddressByPageNum(page)
		if len(items) == 0 {
			return
		}
		var addrs []string
		var balances []*model.Balance
		for _, item := range items {
			addrs = append(addrs, item.Address)
			balance := &model.Balance{
				Address: item.Address,
				GroupID: item.GroupID,
			}
			balances = append(balances, balance)
		}
		acounts, err := rpc.QueryBalance(conf.MainJrpc, addrs)
		if err != nil {
			log15.Error("QueryBalance", "err:", err.Error())
			return
		}
		for _, balance := range balances {
			for _, account := range acounts {
				if balance.Address == account.Addr {
					balance.Balance = account.Balance
					db.InsertData(balance)
					if balance.Balance <= int64(conf.BalanceWarning*1e8) {
						warnings := db.QueryWarningByGroupId(balance.GroupID)
						var flag bool
						for _, w := range warnings {
							if w.Type == DB.BALANCE_WARING && w.IsClosed == 0 {
								flag = true
								break
							}
						}
						if !flag {
							warning := &model.Warning{
								GroupID:  balance.GroupID,
								Type:     DB.BALANCE_WARING,
								Warning:  fmt.Sprintf("groupId %d balance is %v,please handle it as soon as possible!", balance.GroupID, balance.Balance),
								IsClosed: 0,
							}
							db.InsertData(warning)
						}
					}
					break
				}
			}
		}
		if len(items) < 10 {
			return
		}
		page.PageNum++
	}
}

func collectResourceData(db *DB.MonitorDB) {
	page := &model.Page{
		PageSize: 10,
		PageNum:  1,
	}
	for {
		items := db.QueryHostInfoByPageNum(page)
		if len(items) == 0 {
			return
		}

		for _, item := range items {
			if item.UserName == "" || item.PassWd == "" || item.SSHPort == 0 {
				continue
			}
			resource, err := exec.Exec_CollectResource(item)
			if err != nil {
				log15.Error("collectResourceData", "err", err.Error())
				continue
			}
			resource.HostID = item.HostID
			resource.GroupID = item.GroupID
			db.InsertData(resource)
			//TODO 对比告警指标，生成告警信息
			if resource.DiskUsedPercent >= conf.CpuUsedPercentWarning {
				warnings := db.QueryWarningByHostId(item.HostID)
				var flag bool
				for _, w := range warnings {
					if w.Type == DB.DISK_WARING && w.IsClosed == 0 {
						flag = true
						break
					}
				}
				if !flag {
					warning := &model.Warning{
						HostID:   item.HostID,
						GroupID:  item.GroupID,
						Type:     DB.DISK_WARING,
						Warning:  fmt.Sprintf("%s disk usedPercent is %v,please handle it as soon as possible!", item.HostIp, resource.DiskUsedPercent),
						IsClosed: 0,
					}
					db.InsertData(warning)
				}

			}
			if resource.CpuUsedPercent >= conf.CpuUsedPercentWarning {
				warnings := db.QueryWarningByHostId(item.HostID)
				var flag bool
				for _, w := range warnings {
					if w.Type == DB.CPU_WARING && w.IsClosed == 0 {
						flag = true
						break
					}
				}
				if !flag {
					warning := &model.Warning{
						HostID:   item.HostID,
						GroupID:  item.GroupID,
						Type:     DB.CPU_WARING,
						Warning:  fmt.Sprintf("%s cpu usedPercent is %v,please handle it as soon as possible!", item.HostIp, resource.CpuUsedPercent),
						IsClosed: 0,
					}
					db.InsertData(warning)
				}

			}
			if resource.MemUsedPercent >= conf.MemUsedPercentWarning {
				warnings := db.QueryWarningByHostId(item.HostID)
				var flag bool
				for _, w := range warnings {
					if w.Type == DB.MEM_WARNING && w.IsClosed == 0 {
						flag = true
						break
					}
				}
				if !flag {
					warning := &model.Warning{
						HostID:   item.HostID,
						GroupID:  item.GroupID,
						Type:     DB.MEM_WARNING,
						Warning:  fmt.Sprintf("%s mem usedPercent is %v,please handle it as soon as possible!", item.HostIp, resource.MemUsedPercent),
						IsClosed: 0,
					}
					db.InsertData(warning)
				}

			}
		}

		if len(items) < 10 {
			return
		}
		page.PageNum++
	}
}

func clearResourceData(db *DB.MonitorDB) {
	now := time.Now().Unix()
	lastTime := now - conf.ResourceDataHoldTime
	db.DelResourceInfoByTime(lastTime)
}
func clearBalanceData(db *DB.MonitorDB) {
	now := time.Now().Unix()
	lastTime := now - conf.BalanceDataHoldTime
	db.DeleteBalanceByTime(lastTime)
}
