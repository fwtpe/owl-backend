package boss

import (
	"testing"

	qtest "github.com/fwtpe/owl-backend/modules/query/test"
	"github.com/sirupsen/logrus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}

func SetupBossEnv() {
	BeforeEach(func() {
		SetupServerUrl(qtest.GetApiConfigByTestFlag())
	})
}

var _ = BeforeSuite(func() {
	logger.Level = logrus.DebugLevel
})

var _ = AfterSuite(func() {
	logger.Level = logrus.WarnLevel
})
