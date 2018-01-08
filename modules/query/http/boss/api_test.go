package boss

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("encrypt the \"fctoken\" of BOSS service", func() {
	It("The encrypted value of \"SecureFctoken()\" should be as expected", func() {
		testedResult := secureFctoken("hello")

		Expect(testedResult).To(Equal("ecc65534b21a39c5df1c554dec7662c2"))
	})
})
