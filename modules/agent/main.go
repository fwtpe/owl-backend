package main

import (
	"fmt"
	"github.com/fwtpe/owl-backend/common/logruslog"
	"github.com/fwtpe/owl-backend/common/vipercfg"
	"github.com/fwtpe/owl-backend/modules/agent/cron"
	"github.com/fwtpe/owl-backend/modules/agent/funcs"
	"github.com/fwtpe/owl-backend/modules/agent/g"
	"github.com/fwtpe/owl-backend/modules/agent/http"
	"os"
)

func main() {
	vipercfg.Parse()
	vipercfg.Bind()

	if vipercfg.Config().GetBool("version") {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if vipercfg.Config().GetBool("check") {
		funcs.CheckCollector()
		os.Exit(0)
	}

	vipercfg.Load()
	g.ParseConfig(vipercfg.Config().GetString("config"))
	logruslog.Init()

	g.InitRootDir()
	g.InitPublicIps()
	g.InitRpcClients()

	funcs.BuildMappers()

	go cron.InitDataHistory()

	cron.ReportAgentStatus()
	cron.SyncMinePlugins()
	cron.SyncBuiltinMetrics()
	cron.SyncTrustableIps()
	cron.Collect()

	go http.Start()

	select {}

}
