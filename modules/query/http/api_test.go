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

	BeforeEach(func() {
		/**
		 * Set-up environment
		 */
		apiConfig := &g.ApiConfig{
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

		gockConfig.New().Post(g.BOSS_URI_BASE_UPLINK).
			JSON(map[string]interface{}{
				"fcname":   apiConfig.Name,
				"fctoken":  boss.SecureFctoken(apiConfig.Token),
				"pop_name": "sample-idc1",
			}).
			Reply(http.StatusOK).
			JSON(&bmodel.IdcBandwidthResult{
				Status: 1,
				Info:   "当前操作成功了！",
				Result: []*bmodel.IdcBandwidthRow{
					{UplinkTop: 200},
					{UplinkTop: 300},
					{UplinkTop: 400},
				},
			})
	})
	AfterEach(func() {
		gockConfig.Off()
	})

	It("The sum of bandwidth should be 900", func() {
		testedResult := make(map[string]interface{})

		queryIDCsBandwidths("sample-idc1", testedResult)

		testedItems := testedResult["items"].(map[string]interface{})

		Expect(testedItems["upperLimitMB"]).To(Equal(float64(900)))
		Expect(testedItems["IDCName"]).To(Equal("sample-idc1"))
	})
})
