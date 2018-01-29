package http

import (
	"net/http"

	"github.com/fwtpe/owl-backend/common/testing/http/gock"

	"github.com/fwtpe/owl-backend/modules/query/g"
	"github.com/fwtpe/owl-backend/modules/query/http/boss"
	bmodel "github.com/fwtpe/owl-backend/modules/query/model/boss"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("[getLocation(int)]", func() {
	var (
		gockConfig *gock.GockConfig
		apiConfig  *g.ApiConfig
	)

	BeforeEach(func() {
		apiConfig, gockConfig = randomMockBoss()

		/**
		 * Set-up environment
		 */
		g.SetConfig(&g.GlobalConfig{
			Api: apiConfig,
		})
		// :~)

		gockConfig.New().Post(g.BOSS_URI_BASE_GEO).
			JSON(map[string]interface{}{
				"fcname":  apiConfig.Name,
				"fctoken": boss.SecureFctoken(apiConfig.Token),
				"pop_id":  381,
			}).
			Reply(http.StatusOK).
			JSON(&bmodel.LocationResult{
				Status: 1,
				Info:   "当前操作成功了！",
				Count:  3,
				Result: &bmodel.Location{
					Area:     "area-v1",
					Province: "province-v1",
					City:     "city-v1",
				},
			})
	})
	AfterEach(func() {
		gockConfig.Off()
	})

	It("The location data should be as expected", func() {
		testedResult := getLocation(381)

		Expect(testedResult).To(HaveKeyWithValue(Equal("area"), Equal("area-v1")))
		Expect(testedResult).To(HaveKeyWithValue(Equal("province"), Equal("province-v1")))
		Expect(testedResult).To(HaveKeyWithValue(Equal("city"), Equal("city-v1")))
	})
})
