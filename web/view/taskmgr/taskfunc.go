package taskmgr

import (
	"fmt"
	"github.com/harrylee2015/monitor/common"
	"github.com/harrylee2015/monitor/model"
	"time"
)

func getJrpc(ip string, port int64) string {
	return fmt.Sprintf("http://%s:%v", ip, port)
}

//type IsTrue = func(flag  bool) int64;

func collectMonitorData(db *common.MonitorDB) {
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
			lastHeader, err := common.QueryLastHeader(jrpc)
			if err != nil {
				//TODO  err need handler
				continue
			}
			isSync, err := common.QueryIsSync(jrpc)
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

func collectBalanceData(db *common.MonitorDB) {
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
		acounts, err := common.QueryBalance(model.MainJrpc, addrs)
		if err != nil {
			fmt.Println("QueryBalance have err:", err)
			return
		}
		for _, balance := range balances {
			for _, account := range acounts {
				if balance.Address == account.Addr {
					balance.Balance = account.Balance
					db.InsertData(balance)
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

func collectResourceData(db *common.MonitorDB) {
	//TODO
}

func clearResourceData(db *common.MonitorDB) {
	now := time.Now().Unix()
	lastTime := now - model.ResourceDataHoldTime
	db.DelResourceInfoByTime(lastTime)
}
func clearBalanceData(db *common.MonitorDB) {
	now := time.Now().Unix()
	lastTime := now - model.BalanceDataHoldTime
	db.DeleteBalanceByTime(lastTime)
}
