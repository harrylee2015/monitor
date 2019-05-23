package taskmgr

import (
	"github.com/harrylee2015/monitor/common"
	"github.com/harrylee2015/monitor/model"
	"sync"
	"time"
)

func CronTask() {
	db := model.GetDB()
	var group sync.WaitGroup
	group.Add(5)
	go collectMonitor(db, &group)

	go collectBalance(db, &group)

	go collectResource(db, &group)

	go clearBalanceTable(db, &group)

	go clearResourceTable(db, &group)

	group.Wait()

}

func collectMonitor(db *common.MonitorDB, group *sync.WaitGroup) {
	tick := time.Tick(time.Duration(model.CollectMonitorCycle) * time.Second)
	for {
		<-tick
		collectMonitorData(db)
	}
	defer group.Done()
}
func collectResource(db *common.MonitorDB, group *sync.WaitGroup) {
	tick := time.Tick(time.Duration(model.CollectResourceCycle) * time.Second)
	for {
		<-tick
		//TODO
	}
	defer group.Done()
}
func collectBalance(db *common.MonitorDB, group *sync.WaitGroup) {
	tick := time.Tick(time.Duration(model.CollectBalanceCycle) * time.Second)
	for {
		<-tick
		collectBalanceData(db)
	}
	defer group.Done()
}
func clearResourceTable(db *common.MonitorDB, group *sync.WaitGroup) {
	tick := time.Tick(time.Duration(model.ClearDataCycle) * time.Second)
	for {
		<-tick
		clearResourceData(db)
	}
	defer group.Done()
}
func clearBalanceTable(db *common.MonitorDB, group *sync.WaitGroup) {
	tick := time.Tick(time.Duration(model.ClearDataCycle) * time.Second)
	for {
		<-tick
		clearBalanceData(db)
	}
	defer group.Done()
}
