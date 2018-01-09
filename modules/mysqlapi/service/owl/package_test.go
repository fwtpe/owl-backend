package owl

import (
	"testing"

	f "github.com/fwtpe/owl-backend/common/db/facade"
	tDb "github.com/fwtpe/owl-backend/common/testing/db"
	tFlag "github.com/fwtpe/owl-backend/common/testing/flag"
	owlDb "github.com/fwtpe/owl-backend/modules/mysqlapi/rdb/owl"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}

var ginkgoDb = &tDb.GinkgoDb{}
var dbFacade = &f.DbFacade{}

var _ = BeforeSuite(func() {
	dbFacade = ginkgoDb.InitDbFacade()
	owlDb.DbFacade = dbFacade
})

var _ = AfterSuite(func() {
	ginkgoDb.ReleaseDbFacade(dbFacade)
	dbFacade = nil
	owlDb.DbFacade = dbFacade
})

func inTx(sql ...string) {
	dbFacade.SqlDbCtrl.ExecQueriesInTx(sql...)
}

var (
	itFeatures    = tFlag.F_MySql
	itSkipMessage = tFlag.FeatureHelpString(itFeatures)
	itSkip        = tFlag.BuildSkipFactory(tFlag.F_MySql, itSkipMessage)
)
