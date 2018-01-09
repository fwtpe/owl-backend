package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/fwtpe/owl/common/logruslog"
	oos "github.com/fwtpe/owl/common/os"
	"github.com/fwtpe/owl/common/vipercfg"
	"github.com/fwtpe/owl/modules/hbs/g"
	"github.com/fwtpe/owl/modules/hbs/http"
	"github.com/fwtpe/owl/modules/hbs/rpc"
	"github.com/fwtpe/owl/modules/hbs/service"
)

func main() {
	vipercfg.Parse()
	vipercfg.Bind()

	if vipercfg.Config().GetBool("version") {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	vipercfg.Load()
	g.ParseConfig(vipercfg.Config().GetString("config"))
	logruslog.Init()

	service.InitPackage(vipercfg.Config())
	rpc.InitPackage(vipercfg.Config())

	go http.Start()
	go rpc.Start()

	oos.HoldingAndWaitSignal(
		func(signal os.Signal) {
			rpc.Stop()
		},
		os.Interrupt, os.Kill,
		syscall.SIGTERM,
	)
}
