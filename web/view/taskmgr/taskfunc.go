package taskmgr

import (
	"fmt"
	"time"

	DB "github.com/harrylee2015/monitor/common/db"
	"github.com/harrylee2015/monitor/common/email"
	"github.com/harrylee2015/monitor/common/exec"
	"github.com/harrylee2015/monitor/common/rpc"
	"github.com/harrylee2015/monitor/common/sftp"
	"github.com/harrylee2015/monitor/conf"
	"github.com/harrylee2015/monitor/model"
	"github.com/inconshreveable/log15"
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
			//检查主链是否同步
			jrpc := getJrpc(item.MainNet, item.NetPort)
			lastHeader, err := rpc.QueryLastHeader(jrpc)
			if err != nil {
				//TODO  发生错误时，需要更新服务状态,这里需要定义一个通用函数更新即可
				updateAbnormalServerStatus(db, item)
				continue
			}
			isSync, err := rpc.QueryIsSync(jrpc)
			monitor := &model.MainNetMonitor{
				HostIp:          item.MainNet,
				HostID:          item.HostID,
				GroupID:         item.GroupID,
				ServerPort:      item.NetPort,
				ServerStatus:    0,
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
			values := db.QueryMainNetMonitorById(item.GroupID, item.HostID)
			if len(values) == 0 {
				db.InsertData(monitor)

			} else {
				db.UpdateData(monitor)
			}
			//检查平行链
			jrpc = getJrpc(item.HostIp, item.ServerPort)
			lastHeader, err = rpc.QueryLastHeader(jrpc)
			if err != nil {
				//TODO  发生错误时，需要更新服务状态,这里需要定义一个通用函数更新即可
				updateAbnormalServerStatus(db, item)
				continue
			}
			isSync, err = rpc.QueryIsSync(jrpc)
			m := &model.Monitor{
				HostIp:          item.HostIp,
				HostID:          item.HostID,
				GroupID:         item.GroupID,
				ServerPort:      item.ServerPort,
				ServerStatus:    0,
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
			monitors := db.QueryMonitorById(item.GroupID, item.HostID)
			if len(monitors) == 0 {
				db.InsertData(m)

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

func updateAbnormalServerStatus(db *DB.MonitorDB, hostInfo *model.HostInfo) {
	monitor := &model.Monitor{
		HostIp:       hostInfo.HostIp,
		HostID:       hostInfo.HostID,
		GroupID:      hostInfo.GroupID,
		ServerPort:   hostInfo.ServerPort,
		ServerStatus: 1,
		IsSync:       1,
	}
	values := db.QueryMonitorById(hostInfo.GroupID, hostInfo.HostID)
	if len(values) == 0 {
		db.InsertData(monitor)

	} else {
		if values[0].ServerStatus != 1 && values[0].IsSync != 1 {
			db.UpdateData(monitor)
		}

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
				Address:   item.Address,
				GroupID:   item.GroupID,
				Email:     item.Email,
				GroupName: item.GroupName,
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
						warnings := db.QueryWarningByGroupIdAndType(balance.GroupID, DB.BALANCE_WARING)
						if len(warnings) == 0 {
							warning := &model.Warning{
								GroupID:  balance.GroupID,
								Type:     DB.BALANCE_WARING,
								Warning:  fmt.Sprintf("分组%v 当前代扣地址余额 %v BTY,不足 %v BTY,请尽快充值!", balance.GroupID, balance.Balance/1e8, conf.BalanceWarning),
								IsClosed: 0,
							}
							db.InsertData(warning)
							e := &model.Email{
								FromMail: conf.FromEmail,
								PassWd:   conf.PassWd,
								Host:     conf.Host,
								Port:     conf.Port,
								ToMail:   balance.Email,
								Subject:  fmt.Sprintf("%v平行公链运维告警", balance.GroupName),
								Body:     warning.Warning,
							}
							//发送告警邮件
							go email.SendMail(e)
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

//采集机器信息
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
			sftClient, err := sftp.NewSftpClient(item)
			if err != nil {
				log15.Error("collectResourceData", "NewSftpClient err", err.Error())
				continue
			}

			if !sftClient.Exists(exec.GetRemoteScriptsPath()) {
				err = sftClient.UploadFile(exec.GetScriptsPath(), exec.GetRemoteScriptsPath())
				sftClient.Close()
				if err != nil {
					log15.Error("collectResourceData", "UploadFile err", err.Error())
					continue
				}
			}
			//手动关闭连接，不然会造成资源泄漏
			sftClient.Close()
			resource, err := exec.Exec_CollectResource(item)
			if err != nil {
				log15.Error("collectResourceData", "err", err.Error())
				continue
			}
			resource.HostID = item.HostID
			resource.GroupID = item.GroupID
			db.InsertData(resource)
			//TODO 对比告警指标，生成告警信息
			if resource.DiskUsedPercent >= conf.DiskUsedPercentWarning {
				warnings := db.QueryWarningByHostIdAndType(item.HostID, DB.DISK_WARING)
				if len(warnings) == 0 {
					warning := &model.Warning{
						HostID:   item.HostID,
						GroupID:  item.GroupID,
						Type:     DB.DISK_WARING,
						Warning:  fmt.Sprintf("分组%v,%s 节点，disk usedPercent %f%%,达到预警值%f%%!!", item.GroupID, item.HostIp, resource.DiskUsedPercent, conf.DiskUsedPercentWarning),
						IsClosed: 0,
					}

					db.InsertData(warning)
					e := &model.Email{
						FromMail: conf.FromEmail,
						PassWd:   conf.PassWd,
						Host:     conf.Host,
						Port:     conf.Port,
						ToMail:   item.Email,
						Subject:  fmt.Sprintf("%v平行公链运维告警", item.GroupName),
						Body:     warning.Warning,
					}
					//发送告警邮件
					go email.SendMail(e)
				}

			}
			if resource.CpuUsedPercent >= conf.CpuUsedPercentWarning {
				warnings := db.QueryWarningByHostIdAndType(item.HostID, DB.CPU_WARING)
				if len(warnings) == 0 {
					warning := &model.Warning{
						HostID:   item.HostID,
						GroupID:  item.GroupID,
						Type:     DB.CPU_WARING,
						Warning:  fmt.Sprintf("分组%v,%s 节点，cpu usedPercent是%f%%,达到预警值%f%%!", item.GroupID, item.HostIp, resource.CpuUsedPercent, conf.CpuUsedPercentWarning),
						IsClosed: 0,
					}
					db.InsertData(warning)
					e := &model.Email{
						FromMail: conf.FromEmail,
						PassWd:   conf.PassWd,
						Host:     conf.Host,
						Port:     conf.Port,
						ToMail:   item.Email,
						Subject:  fmt.Sprintf("%v平行公链运维告警", item.GroupName),
						Body:     warning.Warning,
					}
					//发送告警邮件
					go email.SendMail(e)
				}

			}
			if resource.MemUsedPercent >= conf.MemUsedPercentWarning {
				warnings := db.QueryWarningByHostIdAndType(item.HostID, DB.MEM_WARNING)
				if len(warnings) == 0 {
					warning := &model.Warning{
						HostID:   item.HostID,
						GroupID:  item.GroupID,
						Type:     DB.MEM_WARNING,
						Warning:  fmt.Sprintf("分组%v,%s 节点，mem usedPercent是%f%%,达到预警值%f%%!", item.GroupID, item.HostIp, resource.MemUsedPercent, conf.MemUsedPercentWarning),
						IsClosed: 0,
					}
					db.InsertData(warning)
					e := &model.Email{
						FromMail: conf.FromEmail,
						PassWd:   conf.PassWd,
						Host:     conf.Host,
						Port:     conf.Port,
						ToMail:   item.Email,
						Subject:  fmt.Sprintf("%v平行公链运维告警", item.GroupName),
						Body:     warning.Warning,
					}
					//发送告警邮件
					go email.SendMail(e)
				}

			}
		}

		if len(items) < 10 {
			return
		}
		page.PageNum++
	}
}

//区块hash一致性检查,检查主链,平行链
func checkBlockHash(db *DB.MonitorDB) {
	//TODO :检查方法如下
	// 1：根据groupId遍历节点moniitor，过滤出有效节点
	// 2：根据有效节点信息获得最新blockHeight,以最小得区块高度为准，进行rpc请求查询。对比blockhash
	page := &model.Page{
		PageSize: 10,
		PageNum:  1,
	}
	for {
		items := db.QueryGroupInfoByPageNum(page)
		if len(items) == 0 {
			return
		}
		for _, item := range items {
			//TODO这里需要做优化
			monitors := db.QueryMonitor(item.GroupID)
			if len(monitors) == 0 {
				continue
			}
			var normals []*model.Monitor
			for _, monitor := range monitors {
				if monitor.ServerStatus == 0 && monitor.IsSync == 0 {
					normals = append(normals, monitor)
				}
			}
			if len(normals) == 0 {
				continue
			}

			// if len(normals) == 1 {
			// 	//TODO: 如果只有一个节点，默认区块hash是一致的
			// 	continue
			// }
			lowestBlock := normals[0]
			for i := 0; i < len(normals); i++ {
				if normals[i].LastBlockHeight < lowestBlock.LastBlockHeight {
					lowestBlock = normals[i]
				}
			}
			//调用rpc接口查询,kv map 用于接收查询结果
			hostInfos := db.QueryHostInfoByGroupId(lowestBlock.GroupID)
			height := lowestBlock.LastBlockHeight
			KV1 := make(map[*model.HostInfo]string)
			COUNT1 := make(map[string]int)
		HERE:
			KV := make(map[*model.HostInfo]string)
			COUNT := make(map[string]int)
			for _, hostInfo := range hostInfos {
				//如果告警已经触发了，就不重复检查
				warnings := db.QueryWarningByHostIdAndType(hostInfo.HostID, DB.HASH_WARING)
				if len(warnings) >= 1 {
					continue
				}
				reply, err := rpc.QueryBlockHash(getJrpc(hostInfo.HostIp, hostInfo.ServerPort), height)
				if err != nil {
					updateAbnormalServerStatus(db, hostInfo)
					continue
				}
				KV[hostInfo] = reply.Hash
				if count, ok := COUNT[reply.Hash]; ok {
					COUNT[reply.Hash] = count + 1
					continue
				}
				COUNT[reply.Hash] = 1
			}
			if len(COUNT) <= 1 {
				//如果前一个区块hash一样得话，停止搜索，回退到下一个块，生成告警信息
				for hash, count := range COUNT1 {
					if len(KV) > count*2 { //如果小于一半，认为这一小半的hash有问题，生成hash告警
						for k, v := range KV1 {
							if v == hash {
								warning := &model.Warning{
									HostID:   k.HostID,
									GroupID:  k.GroupID,
									Type:     DB.HASH_WARING,
									Warning:  fmt.Sprintf("分组 %v,%s节点,on %v 高度 ,blockHash: %s,blockhash 与其他节点blockhash不一致!", k.GroupID, k.HostIp, height+1, v),
									IsClosed: 0,
								}
								db.InsertData(warning)
								e := &model.Email{
									FromMail: conf.FromEmail,
									PassWd:   conf.PassWd,
									Host:     conf.Host,
									Port:     conf.Port,
									ToMail:   item.Email,
									Subject:  fmt.Sprintf("%v平行公链运维告警", item.GroupName),
									Body:     warning.Warning,
								}
								//发送告警邮件
								go email.SendMail(e)
							}
						}
					}
				}
			}
			//对查询结果进行归类比对，我们认为多数节点得结果时正确的
			if len(COUNT) >= 2 { //D当COUNT >=2时，说明发生分叉了，这时需要具体定位了
				if height-1 > 0 {
					height = height - 1
					KV1 = KV
					COUNT1 = COUNT
					goto HERE
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
