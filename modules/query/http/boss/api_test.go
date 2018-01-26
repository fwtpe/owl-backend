package boss

import (
	"strconv"

	ojson "github.com/fwtpe/owl-backend/common/json"

	"github.com/fwtpe/owl-backend/modules/query/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var bossSkipper = test.BossSkipper

var _ = Describe("Load ip data of platforms", bossSkipper.PrependBeforeEach(func() {
	SetupBossEnv()

	It("The data \"[]*PlatformIps\" [1] platform data", func() {
		testedData := LoadIpDataOfPlatforms()

		GinkgoT().Logf("[Platform data] Size[%d]. First row(JSON): %s", len(testedData), ojson.MarshalJSON(testedData[0]))
		Expect(len(testedData)).To(BeNumerically(">=", 1))
	})
}))

var _ = Describe("Load IDC data", bossSkipper.PrependBeforeEach(func() {
	SetupBossEnv()

	// See model/boss for JSON specification of data
	It("The data \"[]*bmodel.IdcIps\" should have at least [1] IDC data", func() {
		testedData := LoadIdcData()

		GinkgoT().Logf("[IDC data] Size[%d] First row(JSON): %s", len(testedData), ojson.MarshalJSON(testedData[0]))
		Expect(len(testedData)).To(BeNumerically(">=", 1))
	})
}))

var _ = Describe("Load IDC bandwidth", bossSkipper.PrependBeforeEach(func() {
	SetupBossEnv()

	// See model/boss for JSON specification of data
	It("The data \"[]*IdcBandwidthRow\" should have at least 1 bandwidth data", func() {
		sampleIdcName := LoadIdcData()[0].IpList[0].Pop

		testedData := LoadIdcBandwidth(sampleIdcName)

		GinkgoT().Logf(
			"[Bandwidth data of IDC(%s)] (JSON): %s",
			sampleIdcName, ojson.MarshalJSON(testedData),
		)
		Expect(len(testedData)).To(BeNumerically(">=", 1))
	})
}))

var _ = Describe("Load location data", bossSkipper.PrependBeforeEach(func() {
	SetupBossEnv()

	It("The data \"*Location\" should have corresponding data of location", func() {
		sampleIdcId, _ := strconv.Atoi(LoadIdcData()[0].IpList[0].PopId)

		testedLocation := LoadLocationData(sampleIdcId)

		GinkgoT().Logf(
			"[Location data of IDC(%d)] (JSON): %s",
			sampleIdcId, ojson.MarshalJSON(testedLocation),
		)
		Expect(testedLocation).To(
			PointTo(MatchAllFields(Fields{
				"Area":     Not(BeEmpty()),
				"City":     Not(BeEmpty()),
				"Province": Not(BeEmpty()),
			})),
		)
	})
}))

var _ = XDescribe("encrypt the \"fctoken\" of BOSS service", bossSkipper.PrependBeforeEach(func() {
	SetupBossEnv()

	// Any of "Load xxx" testing would test the encryption of token
	It("The encrypted value of \"SecureFctoken()\" should be as expected", func() {
		testedResult := SecureFctoken("hello")
		Expect(testedResult).To(Equal("ecc65534b21a39c5df1c554dec7662c2"))
	})
}))

var _ = Describe("Load data of platforms' detail", func() {
	SetupBossEnv()

	It("The data \"[]*PlatformDetail\" should have at least 1 row", func() {
		testedResult := LoadDetailOfPlatforms()

		GinkgoT().Logf("Total size: [%d]. First row: %s", len(testedResult), ojson.MarshalJSON(testedResult[0]))

		Expect(len(testedResult)).To(BeNumerically(">=", 1))
	})
})
