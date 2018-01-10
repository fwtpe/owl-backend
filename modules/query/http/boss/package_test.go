package boss

import (
	"testing"

	qtest "github.com/Cepave/open-falcon-backend/modules/query/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}

func SetupBossEnvOrSkip() {
	BeforeEach(func() {
		qtest.SkipIfNoBossConfig()
		SetupServerUrl(qtest.GetApiConfigByTestFlag())
	})
}
