package types

import (
	"fmt"
	DB "github.com/harrylee2015/monitor/common/db"
)

const (
	DBNAME       = "SQLITE3"
	WARNING_TYPE = int64(iota)
	MEM_wARNING
	CPU_WARING
	DISK_WARING
	BALANCE_WARING
	HASH_WARING
)

var (
	register = make(map[string]interface{})
)

func RegisterDB() {
	if db, ok := register[DBNAME]; !ok {
		db = DB.NewMonitorDB()
		register[DBNAME] = db
		return
	}
	panic(fmt.Errorf("can't duplicate registration db!"))
}

func GetDB() *DB.MonitorDB {
	if value, ok := register[DBNAME]; ok {
		if db, ok := value.(*DB.MonitorDB); ok {
			return db
		}
	}
	panic(fmt.Errorf("can't get db!"))
	return nil
}

func CloseDB() {
	GetDB().Close()
}