package test

import (
	"github.com/fwtpe/owl-backend/modules/alarm/g"
	"github.com/fwtpe/owl-backend/modules/alarm/model"
)

func initTest() {
	g.ParseConfig("../test_cfg.json")
	model.InitDatabase()
}
