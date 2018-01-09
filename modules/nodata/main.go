package main

import (
	"fmt"
	"github.com/fwtpe/owl/common/logruslog"
	"github.com/fwtpe/owl/common/vipercfg"
	"os"

	"github.com/fwtpe/owl/modules/nodata/collector"
	"github.com/fwtpe/owl/modules/nodata/config"
	"github.com/fwtpe/owl/modules/nodata/g"
	"github.com/fwtpe/owl/modules/nodata/http"
	"github.com/fwtpe/owl/modules/nodata/judge"
	"github.com/fwtpe/owl/modules/nodata/sender"
)

func main() {
	vipercfg.Parse()
	vipercfg.Bind()

	if vipercfg.Config().GetBool("version") {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}
	if vipercfg.Config().GetBool("vg") {
		fmt.Println(g.VERSION, g.COMMIT)
		os.Exit(0)
	}

	// global config
	vipercfg.Load()
	g.ParseConfig(vipercfg.Config().GetString("config"))
	logruslog.Init()
	// proc
	g.StartProc()

	// config
	config.Start()
	// collector
	collector.Start()
	// judge
	judge.Start()
	// sender
	sender.Start()

	// http
	http.Start()

	select {}
}
