package http

import (
	"fmt"
	"net/http"

	"gopkg.in/h2non/gentleman-mock.v2"

	ojson "github.com/fwtpe/owl-backend/common/json"
	"github.com/fwtpe/owl-backend/common/testing/http/gock"

	"github.com/fwtpe/owl-backend/modules/query/g"
	"github.com/fwtpe/owl-backend/modules/query/http/boss"
	bmodel "github.com/fwtpe/owl-backend/modules/query/model/boss"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("[queryPlatformJSON]", func() {
	gockConfig := gock.GockConfigBuilder.NewConfigByRandom()

	BeforeEach(func() {
		/**
		 * Set-up environment
		 */
		apiConfig := &g.ApiConfig{
			Name:     "mock-101",
			Token:    "mock-token-101",
			BossBase: gockConfig.NewHttpConfig().Url,
			Map:      gockConfig.NewHttpConfig().Url + g.BOSS_URI_BASE_MAP,
		}
		g.SetConfig(&g.GlobalConfig{
			Api: apiConfig,
		})
		// :~)

		boss.SetPlugins(mock.Plugin)
		boss.SetupServerUrl(apiConfig)

		gockConfig.New().Get(fmt.Sprintf(
			g.BOSS_URI_BASE_MAP+g.BOSS_PLATFORM_PATH_TMPL,
			apiConfig.Name, boss.SecureFctoken(apiConfig.Token),
		)).
			Reply(http.StatusOK).
			JSON(map[string]interface{}{})
	})

	AfterEach(func() {
		gockConfig.Off()
	})

	It("The data of platform should be as expected", func() {
	})
})
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

var _ = Describe("[loadIpDataOfPlatforms]", func() {
	gockConfig := gock.GockConfigBuilder.NewConfigByRandom()

	Context("Normal case", func() {
		BeforeEach(func() {
			/**
			 * Set-up environment
			 */
			apiConfig := &g.ApiConfig{
				Name:     "fmock-1",
				Token:    "fmock-token-1",
				BossBase: gockConfig.NewHttpConfig().Url,
			}
			g.SetConfig(&g.GlobalConfig{
				Api: apiConfig,
			})
			// :~)

			boss.SetPlugins(mock.Plugin)
			boss.SetupServerUrl(apiConfig)

			gockConfig.New().Get(
				fmt.Sprintf(
					g.BOSS_URI_BASE_MAP+g.BOSS_PLATFORM_PATH_TMPL,
					apiConfig.Name, boss.SecureFctoken(apiConfig.Token),
				),
			).
				Reply(http.StatusOK).
				JSON(&bmodel.PlatformIpResult{
					Status: 1,
					Info:   "当前操作成功了！",
					Result: []*bmodel.PlatformIps{
						{
							Name: "plt-981",
							IpList: []*bmodel.PlatformIp{
								{Hostname: "ga01.z.net", Ip: "97.6.1.41", PopId: "33", Status: "1"},
								{Hostname: "ga02.z.net", Ip: "97.6.1.42", PopId: "34", Status: "1"},
								{Hostname: "ga03.z.net", Ip: "97.6.1.43", PopId: "45", Status: "1"},
							},
						},
						{
							Name: "plt-982",
							IpList: []*bmodel.PlatformIp{
								{Hostname: "gb01.z.net", Ip: "97.6.12.41", PopId: "43", Status: "1"},
								{Hostname: "gb02.z.net", Ip: "97.6.12.42", PopId: "44", Status: "1"},
								{Hostname: "gb03.z.net", Ip: "97.6.12.43", PopId: "47", Status: "1"},
							},
						},
					},
				})
		})
		AfterEach(func() {
			gockConfig.Off()
		})

		It("JSON content should be as expected", func() {
			testedResult := make(map[string]interface{})

			loadIpDataOfPlatforms(testedResult, make(map[string]interface{}))

			GinkgoT().Logf("[Platform] JSON Content: %s", ojson.MarshalJSON(testedResult))
			Expect(testedResult["status"]).To(BeEquivalentTo(1))
			Expect(testedResult["result"]).To(HaveLen(2))
			Expect(testedResult["result"]).To(ContainElement(HaveLen(2)))
		})
	})
})

var _ = Describe("[getIpFromHostname]", func() {
	DescribeTable("Checks the converted ip string as expected",
		func(sourceHostName string, expectedResult string) {
			testedResult := getIpFromHostname(sourceHostName, make(map[string]interface{}))

			Expect(testedResult).To(Equal(expectedResult))
		},
		Entry("Normal", "bj-cnc-019-061-123-201", "19.61.123.201"),
		Entry("Cannot be parsed", "nothing", ""),
		Entry("Cannot be parsed(one of ip value)", "kz-abk-019-8c-123-201", ""),
	)
})
var _ = Describe("[getIspFromHostname]", func() {
	DescribeTable("Checks the converted ISP string as expected",
		func(sourceHostName string, expectedResult string) {
			testedResult := getIspFromHostname(sourceHostName)
			Expect(testedResult).To(Equal(expectedResult))
		},
		Entry("Normal", "bjb-ck-091-111-041-35", "bjb"),
		Entry("Short name", "ack_zs_091_111", "ack"),
		Entry("Short name(with partial IP)", "cjc-zs-091-111.ball.com", "cjc"),
		Entry("Cannot be parsed", "nothing", ""),
		Entry("Cannot be parsed", "ll.091", ""),
	)
})
