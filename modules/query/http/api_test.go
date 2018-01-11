package http

import (
	"net/http"

	"gopkg.in/h2non/gentleman-mock.v2"

	"github.com/fwtpe/owl-backend/common/testing/http/gock"

	"github.com/fwtpe/owl-backend/modules/query/g"
	"github.com/fwtpe/owl-backend/modules/query/http/boss"
	bmodel "github.com/fwtpe/owl-backend/modules/query/model/boss"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("[queryIDCsBandwidths]", func() {
	gockConfig := gock.GockConfigBuilder.NewConfigByRandom()

	var idcName string
	var bandwidthData []*bmodel.IdcBandwidthRow
	var apiConfig *g.ApiConfig

	BeforeEach(func() {
		/**
		 * Set-up environment
		 */
		apiConfig = &g.ApiConfig{
			Name:     "mock-2",
			Token:    "mock-token-2",
			BossBase: gockConfig.NewHttpConfig().Url,
		}
		g.SetConfig(&g.GlobalConfig{
			Api: apiConfig,
		})
		// :~)

		boss.SetPlugins(mock.Plugin)
		boss.SetupServerUrl(apiConfig)
	})
	AfterEach(func() {
		gockConfig.Off()
	})

	JustBeforeEach(func() {
		gockConfig.New().Post(g.BOSS_URI_BASE_UPLINK).
			JSON(map[string]interface{}{
				"fcname":   apiConfig.Name,
				"fctoken":  boss.SecureFctoken(apiConfig.Token),
				"pop_name": idcName,
			}).
			Reply(http.StatusOK).
			JSON(&bmodel.IdcBandwidthResult{
				Status: 1,
				Info:   "当前操作成功了！",
				Result: bandwidthData,
			})
	})

	Context("Viable bandwidth", func() {
		BeforeEach(func() {
			idcName = "sample-idc1"
			bandwidthData = []*bmodel.IdcBandwidthRow{
				{UplinkTop: 200},
				{UplinkTop: 300},
				{UplinkTop: 400},
			}
		})
		It("The sum of bandwidth should be 900", func() {
			testedResult := make(map[string]interface{})

			queryIDCsBandwidths(idcName, testedResult)

			testedItems := testedResult["items"].(map[string]interface{})

			Expect(testedItems["upperLimitMB"]).To(Equal(float64(900)))
			Expect(testedItems["IDCName"]).To(Equal(idcName))
		})
	})

	Context("Bandwidth is empty", func() {
		BeforeEach(func() {
			idcName = "sample-idc2"
			bandwidthData = make([]*bmodel.IdcBandwidthRow, 0)
		})

		It("The sum of bandwidth should be 0", func() {
			testedResult := make(map[string]interface{})

			queryIDCsBandwidths(idcName, testedResult)

			testedItems := testedResult["items"].(map[string]interface{})

			Expect(testedItems["upperLimitMB"]).To(Equal(float64(0)))
			Expect(testedItems["IDCName"]).To(Equal(idcName))
		})
	})
})
