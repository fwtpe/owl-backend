package service

import (
	"testing"

	f "github.com/fwtpe/owl/common/db/facade"
	tDb "github.com/fwtpe/owl/common/testing/db"
	tFlag "github.com/fwtpe/owl/common/testing/flag"
	"github.com/fwtpe/owl/modules/mysqlapi/rdb"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}

var (
	ginkgoDb = &tDb.GinkgoDb{}
	dbFacade = &f.DbFacade{}

	itDbs  = tFlag.OWL_DB_PORTAL
	itSkip = tFlag.BuildSkipFactoryOfOwlDb(itDbs, tFlag.OwlDbHelpString(itDbs))
)

var _ = BeforeSuite(func() {
	dbFacade = ginkgoDb.InitDbFacadeByFlag(tFlag.OWL_DB_PORTAL)
	rdb.DbFacade = dbFacade
})

var _ = AfterSuite(func() {
	ginkgoDb.ReleaseDbFacade(dbFacade)
	dbFacade = nil
	rdb.DbFacade = dbFacade
})
