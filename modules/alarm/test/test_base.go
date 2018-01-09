package test

import (
	"github.com/fwtpe/owl/modules/alarm/g"
	"github.com/fwtpe/owl/modules/alarm/model"
)

func initTest() {
	g.ParseConfig("../test_cfg.json")
	model.InitDatabase()
}
