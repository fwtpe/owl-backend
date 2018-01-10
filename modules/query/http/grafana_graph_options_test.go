
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

var _ = Describe("[getLocation(int)]", func() {
	var fakeUrl = "http://fake-glidc.net"

	BeforeEach(func() {
		/**
		 * Set-up environment
		 */
		apiConfig := &g.ApiConfig {
			Name: "mock-3",
			Token: "mock-token-3",
			BossBase: fakeUrl,
		}
		g.SetConfig(&g.GlobalConfig {
			Api: apiConfig,
		})
		// :~)

		boss.SetPlugins(mock.Plugin)
		boss.SetupServerUrl(apiConfig)

		gock.New(fakeUrl).Post(g.BOSS_URI_BASE_GEO).
			JSON(map[string]interface{} {
				"fcname":   apiConfig.Name,
				"fctoken":  boss.SecureFctoken(apiConfig.Token),
				"pop_id": 381,
			}).
			Reply(http.StatusOK).
			JSON(&bmodel.LocationResult {
				Status: 1,
				Info: "当前操作成功了！",
				Count: 3,
				Result: &bmodel.Location {
					Area: "area-v1",
					Province: "province-v1",
					City: "city-v1",
				},
			})
	})
	AfterEach(func() {
		gock.Off()
	})

	It("The location data should be as expected", func() {
		testedResult := getLocation(381)

		Expect(testedResult).To(HaveKeyWithValue(Equal("area"), Equal("area-v1")))
		Expect(testedResult).To(HaveKeyWithValue(Equal("province"), Equal("province-v1")))
		Expect(testedResult).To(HaveKeyWithValue(Equal("city"), Equal("city-v1")))
	})
})
