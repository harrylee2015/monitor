package taskmgr

import (
	"github.com/harrylee2015/monitor/common/db"
	"github.com/harrylee2015/monitor/conf"
	"github.com/harrylee2015/monitor/types"
	"sync"
	"time"
)

func CronTask() {
	db := types.GetDB()
	var group sync.WaitGroup
	group.Add(5)
	go collectMonitor(db, &group)

	go collectBalance(db, &group)

	go collectResource(db, &group)

	go clearBalanceTable(db, &group)

	go clearResourceTable(db, &group)

	group.Wait()

}

func collectMonitor(db *db.MonitorDB, group *sync.WaitGroup) {
	tick := time.Tick(time.Duration(conf.CollectMonitorCycle) * time.Second)
	for {
		<-tick
		collectMonitorData(db)
	}
	defer group.Done()
}
func collectResource(db *db.MonitorDB, group *sync.WaitGroup) {
	tick := time.Tick(time.Duration(conf.CollectResourceCycle) * time.Second)
	for {
		<-tick
		collectResourceData(db)
	}
	defer group.Done()
}
func collectBalance(db *db.MonitorDB, group *sync.WaitGroup) {
	tick := time.Tick(time.Duration(conf.CollectBalanceCycle) * time.Second)
	for {
		<-tick
		collectBalanceData(db)
	}
	defer group.Done()
}
func clearResourceTable(db *db.MonitorDB, group *sync.WaitGroup) {
	tick := time.Tick(time.Duration(conf.ClearDataCycle) * time.Second)
	for {
		<-tick
		clearResourceData(db)
	}
	defer group.Done()
}
func clearBalanceTable(db *db.MonitorDB, group *sync.WaitGroup) {
	tick := time.Tick(time.Duration(conf.ClearDataCycle) * time.Second)
	for {
		<-tick
		clearBalanceData(db)
	}
	defer group.Done()
}
