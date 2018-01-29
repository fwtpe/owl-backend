package http

import (
	"fmt"
	"net/http"

	ojson "github.com/fwtpe/owl-backend/common/json"
	"github.com/fwtpe/owl-backend/common/testing/http/gock"

	"github.com/fwtpe/owl-backend/modules/query/g"
	"github.com/fwtpe/owl-backend/modules/query/http/boss"
	bmodel "github.com/fwtpe/owl-backend/modules/query/model/boss"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("[queryIDCsBandwidths]", func() {
	var (
		gockConfig    *gock.GockConfig
		apiConfig     *g.ApiConfig
		idcName       string
		bandwidthData []*bmodel.IdcBandwidthRow
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
	var (
		gockConfig *gock.GockConfig
		apiConfig  *g.ApiConfig
	)

	Context("Normal case", func() {
		BeforeEach(func() {
			apiConfig, gockConfig = randomMockBoss()
			/**
			 * Set-up environment
			 */
			g.SetConfig(&g.GlobalConfig{
				Api: apiConfig,
			})
			// :~)

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

var _ = Describe("[getPlatformContact]", func() {
	var (
		gockConfig *gock.GockConfig
		apiConfig  *g.ApiConfig
	)

	Context("Normal Case", func() {
		BeforeEach(func() {
			apiConfig, gockConfig = randomMockBoss()

			/**
			 * Set-up environment
			 */
			g.SetConfig(&g.GlobalConfig{
				Api: apiConfig,
			})
			// :~)

			gockConfig.New().Post(g.BOSS_URI_BASE_CONTACT).
				JSON(map[string]interface{}{
					"fcname":       apiConfig.Name,
					"fctoken":      boss.SecureFctoken(apiConfig.Token),
					"platform_key": "ck01.pf93,ck02.gj01",
				}).
				Reply(http.StatusOK).
				JSON(&bmodel.PlatformContactResult{
					Status: 1,
					Info:   "当前操作成功了！",
					Result: map[string]*bmodel.ContactUsers{
						"ck01.pf93": {
							Principals: []*bmodel.ContactUser{
								{
									Id: "4901", RealName: "john-1",
									CellPhoneNumber: "19081141", TelephoneNumber: "0988117651", Email: "john-1@fw.com.going",
								},
								{
									Id: "4902", RealName: "bob-1",
									CellPhoneNumber: "21081141", TelephoneNumber: "0978187651", Email: "bob-1@fw.com.going",
								},
							},
							Backupers: []*bmodel.ContactUser{
								{
									Id: "5021", RealName: "john-2",
									CellPhoneNumber: "19081142", TelephoneNumber: "0988117652", Email: "john-2@fw.com.going",
								},
								{
									Id: "5022", RealName: "bob-2",
									CellPhoneNumber: "21081142", TelephoneNumber: "0978187652", Email: "bob-2@fw.com.going",
								},
							},
							Upgraders: []*bmodel.ContactUser{
								{
									Id: "7348", RealName: "john-3",
									CellPhoneNumber: "19081143", TelephoneNumber: "0988117653", Email: "john-3@fw.com.going",
								},
								{
									Id: "7349", RealName: "bob-3",
									CellPhoneNumber: "21081143", TelephoneNumber: "0978187653", Email: "bob-3@fw.com.going",
								},
							},
						},
						"ck02.gj01": {
							Principals: []*bmodel.ContactUser{
								{
									Id: "14051", RealName: "kgg-1",
									CellPhoneNumber: "22081141", TelephoneNumber: "0913117651", Email: "lukg33@fw.com.going",
								},
							},
							Backupers: []*bmodel.ContactUser{
								{
									Id: "14052", RealName: "kgg-2",
									CellPhoneNumber: "32081141", TelephoneNumber: "0983117441", Email: "lukg71@fw.com.going",
								},
							},
						},
					},
				})
		})
		AfterEach(func() {
			gockConfig.Off()
		})

		It("The data of \"result\" should be 2 backuper(deputy), 2 principals, and 2 upgraders", func() {
			testedNodes := make(map[string]interface{})

			getPlatformContact("ck01.pf93,ck02.gj01", testedNodes)

			GinkgoT().Logf("Platform contacts: %#v", testedNodes)

			Expect(testedNodes["count"]).To(Equal(2))
			Expect(testedNodes["platform"]).To(Equal("ck01.pf93,ck02.gj01"))

			resultPlatforms := testedNodes["result"].(map[string]interface{})["items"].(map[string]interface{})

			testedPlatform := resultPlatforms["ck01.pf93"].(map[string]map[string]string)
			Expect(testedPlatform["principal"]).To(And(
				HaveKeyWithValue(Equal("name"), Equal("john-1")),
				HaveKeyWithValue(Equal("phone"), Equal("19081141")),
				HaveKeyWithValue(Equal("email"), Equal("john-1@fw.com.going")),
			))
			Expect(testedPlatform["deputy"]).To(And(
				HaveKeyWithValue(Equal("name"), Equal("john-2")),
				HaveKeyWithValue(Equal("phone"), Equal("19081142")),
				HaveKeyWithValue(Equal("email"), Equal("john-2@fw.com.going")),
			))
			Expect(testedPlatform["upgrader"]).To(And(
				HaveKeyWithValue(Equal("name"), Equal("john-3")),
				HaveKeyWithValue(Equal("phone"), Equal("19081143")),
				HaveKeyWithValue(Equal("email"), Equal("john-3@fw.com.going")),
			))
		})
	})
})
