package http

import (
	ch "gopkg.in/check.v1"
	"testing"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"

	tFlag "github.com/fwtpe/owl-backend/common/testing/flag"
	"github.com/fwtpe/owl-backend/common/testing/http/gock"

	"github.com/fwtpe/owl-backend/modules/query/g"
	qtest "github.com/fwtpe/owl-backend/modules/query/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}

func TestByCheck(t *testing.T) {
	ch.TestingT(t)
}

var testFlags = tFlag.NewTestFlags()

var mockMySqlApi = gock.GockConfigBuilder.NewConfig(
	"ack.com.cc", 22060,
)

var skipItOnMySqlApi = tFlag.BuildSkipFactory(tFlag.F_ItWeb, tFlag.FeatureHelpString(tFlag.F_ItWeb))
var skipBossDb = tFlag.BuildSkipFactoryOfOwlDb(tFlag.OWL_DB_BOSS, tFlag.OwlDbHelpString(tFlag.OWL_DB_BOSS))

func SetupBossEnvOrSkip() {
	qtest.SkipIfNoBossConfig()
	g.SetConfig(&g.GlobalConfig{
		Api: qtest.GetApiConfigByTestFlag(),
	})
}

var initBoss = false

func RegisterBossOrmOrSkip() {
	skipBossDb.Skip()

	if initBoss {
		return
	}

	orm.RegisterModel(new(Contacts), new(Hosts), new(Idcs), new(Ips), new(Platforms))

	orm.RegisterDataBase("default", "mysql", testFlags.GetMysqlOfOwlDb(tFlag.OWL_DB_BOSS), 30)
	orm.RegisterDataBase("boss", "mysql", testFlags.GetMysqlOfOwlDb(tFlag.OWL_DB_BOSS), 30)

	initBoss = true
}
