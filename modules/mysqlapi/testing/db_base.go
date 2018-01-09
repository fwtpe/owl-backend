package testing

import (
	dbTest "github.com/fwtpe/owl/common/testing/db"
	"github.com/fwtpe/owl/modules/mysqlapi/rdb"
	check "gopkg.in/check.v1"
)

// The base environment for RDB testing
func InitRdb(c *check.C) {
	dbTest.SetupByViableDbConfig(c, rdb.InitPortalRdb)
}
func ReleaseRdb(c *check.C) {
	rdb.ReleaseAllRdb()
}
