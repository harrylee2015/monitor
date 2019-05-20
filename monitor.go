package main

import (
	"flag"
	"os"
	"path/filepath"
	"sync"

	"github.com/inconshreveable/log15"
	"gitlab.33.cn/chain33/chain33/common/limits"
	"github.com/harrylee2015/guesshand/charge"
	"github.com/harrylee2015/guesshand/common"
	"github.com/harrylee2015/guesshand/types"
)

var (
	configPath = flag.String("f", "guess.toml", "configuration file")
	datadir    = flag.String("datadir", "", "data dir of guesshand")
	group      sync.WaitGroup //定义一个同步等待的组
)

func main() {
	flag.Parse()
	os.Chdir(pwd())
	d, _ := os.Getwd()
	log15.Info("current dir:", "dir", d)
	err := limits.SetLimits()
	if err != nil {
		panic(err)
	}

	cfg := types.InitCfg(*configPath)
	if *datadir != "" {
		types.ResetDatadir(cfg, *datadir)
	}
	types.SetCfgValue(cfg)

	db := common.NewGuessDB()

	//比对高度,若数据库中存储高度小于配置文件的，则更新
	if db.GetGameNetHeight() < int(types.GameChainStartHeight) {
		db.SetGameNetHeight(int(types.GameChainStartHeight))
	}
	if db.GetMainNetHeight() < int(types.MainChainStartHeight) {
		db.SetMainNetHeight(int(types.MainChainStartHeight))
	}

	group.Add(4)
	// 定时从主链查询充值交易，并生成游戏链的冲币交易
	go charge.CheckCharge(db, &group)

	// 定时检查游戏链的冲币交易，确保成功或进行重做等
	go charge.ConfirmCharge(db, &group)

	// 定时从游戏链查询提币交易，并生成主链的转账交易
	go charge.CheckWithdraw(db, &group)

	// 定时检查主链的转账交易，确保成功或进行重做
	go charge.ConfirmWithdraw(db, &group)


	// 阻塞等待
	group.Wait()
	defer func() {
		db.Close()
	}()
}

func pwd() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return dir
}
