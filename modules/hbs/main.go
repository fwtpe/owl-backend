package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/fwtpe/owl-backend/common/logruslog"
	oos "github.com/fwtpe/owl-backend/common/os"
	"github.com/fwtpe/owl-backend/common/vipercfg"
	"github.com/fwtpe/owl-backend/modules/hbs/g"
	"github.com/fwtpe/owl-backend/modules/hbs/http"
	"github.com/fwtpe/owl-backend/modules/hbs/rpc"
	"github.com/fwtpe/owl-backend/modules/hbs/service"
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
