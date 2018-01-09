package testing

import (
	dbTest "github.com/fwtpe/owl-backend/common/testing/db"
	"github.com/fwtpe/owl-backend/modules/mysqlapi/rdb"
	check "gopkg.in/check.v1"
)

// The base environment for RDB testing
func InitRdb(c *check.C) {
	dbTest.SetupByViableDbConfig(c, rdb.InitPortalRdb)
}
func ReleaseRdb(c *check.C) {
	rdb.ReleaseAllRdb()
}
