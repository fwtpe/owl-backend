package http

import (
	ch "gopkg.in/check.v1"
	"testing"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"

	odb "github.com/fwtpe/owl-backend/common/db"
	"github.com/fwtpe/owl-backend/common/db/facade"
	tFlag "github.com/fwtpe/owl-backend/common/testing/flag"
	"github.com/fwtpe/owl-backend/common/testing/http/gock"

	db "github.com/fwtpe/owl-backend/modules/query/database"
	"github.com/fwtpe/owl-backend/modules/query/http/boss"
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

var skipBossDb = tFlag.BuildSkipFactoryOfOwlDb(
	tFlag.OWL_DB_BOSS, tFlag.OwlDbHelpString(tFlag.OWL_DB_BOSS),
)

func SetupBossEnv() {
	BeforeEach(func() {
		boss.SetupServerUrl(qtest.GetApiConfigByTestFlag())
	})
}

var bossInTx func(sql ...string)
var _ = BeforeSuite(func() {
	log.Level = logrus.DebugLevel

	if !testFlags.HasMySqlOfOwlDb(tFlag.OWL_DB_BOSS) {
		return
	}

	bossFacade := &facade.DbFacade{}
	bossFacade.Open(&odb.DbConfig{
		Dsn: testFlags.GetMysqlOfOwlDb(tFlag.OWL_DB_BOSS), MaxIdle: 2,
	})

	db.BossDbFacade = bossFacade
	bossInTx = db.BossDbFacade.SqlDbCtrl.ExecQueriesInTx

	orm.RegisterModel(new(Contacts), new(Hosts), new(Idcs), new(Ips), new(Platforms))

	orm.RegisterDataBase("default", "mysql", testFlags.GetMysqlOfOwlDb(tFlag.OWL_DB_BOSS), 30)
	orm.RegisterDataBase("boss", "mysql", testFlags.GetMysqlOfOwlDb(tFlag.OWL_DB_BOSS), 30)
})
var _ = AfterSuite(func() {
	log.Level = logrus.WarnLevel

	db.BossDbFacade = nil
})
