package http

import (
	"net/http"

	"gopkg.in/h2non/gock.v1"
	"gopkg.in/h2non/gentleman-mock.v2"

	"github.com/Cepave/open-falcon-backend/modules/query/http/boss"
	"github.com/Cepave/open-falcon-backend/modules/query/g"
	bmodel "github.com/Cepave/open-falcon-backend/modules/query/model/boss"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("[queryIDCsBandwidths]", func() {
	var fakeUrl = "http://fake-idcbd.net"

	BeforeEach(func() {
		/**
		 * Set-up environment
		 */
		apiConfig := &g.ApiConfig {
			Name: "mock-2",
			Token: "mock-token-2",
			BossBase: fakeUrl,
		}
		g.SetConfig(&g.GlobalConfig {
			Api: apiConfig,
		})
		// :~)

		boss.SetPlugins(mock.Plugin)
		boss.SetupServerUrl(apiConfig)

		gock.New(fakeUrl).Post(g.BOSS_URI_BASE_UPLINK).
			JSON(map[string]interface{} {
				"fcname":   apiConfig.Name,
				"fctoken":  boss.SecureFctoken(apiConfig.Token),
				"pop_name": "sample-idc1",
			}).
			Reply(http.StatusOK).
			JSON(&bmodel.IdcBandwidthResult{
				Status: 1,
				Info: "当前操作成功了！",
				Result: []*bmodel.IdcBandwidthRow {
					{ UplinkTop: 200 },
					{ UplinkTop: 300 },
					{ UplinkTop: 400 },
				},
			})
	})
	AfterEach(func() {
		gock.Off()
	})

	It("The sum of bandwidth should be 900", func() {
		testedResult := make(map[string]interface{})

		queryIDCsBandwidths("sample-idc1", testedResult)

		testedItems := testedResult["items"].(map[string]interface{})

		Expect(testedItems["upperLimitMB"]).To(Equal(float64(900)))
		Expect(testedItems["IDCName"]).To(Equal("sample-idc1"))
	})
})
