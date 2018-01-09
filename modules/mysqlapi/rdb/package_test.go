package rdb

import (
	"testing"

	ch "gopkg.in/check.v1"

	tDb "github.com/fwtpe/owl-backend/common/testing/db"
	tFlag "github.com/fwtpe/owl-backend/common/testing/flag"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	itSkipMessage = tFlag.OwlDbHelpString(tFlag.OWL_DB_PORTAL)
	itSkip        = tFlag.BuildSkipFactoryOfOwlDb(tFlag.OWL_DB_PORTAL, itSkipMessage)
)

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}

func TestByCheck(t *testing.T) {
	ch.TestingT(t)
}

func inTx(sql ...string) {
	DbFacade.SqlDbCtrl.ExecQueriesInTx(sql...)
}

var ginkgoDb = &tDb.GinkgoDb{}
var _ = BeforeSuite(func() {
	DbFacade = ginkgoDb.InitDbFacadeByFlag(tFlag.OWL_DB_PORTAL)
})

var _ = AfterSuite(func() {
	ginkgoDb.ReleaseDbFacade(DbFacade)
	DbFacade = nil
})
