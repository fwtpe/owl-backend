package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/astaxie/beego/orm"
	"gopkg.in/h2non/gentleman-mock.v2"

	"github.com/fwtpe/owl-backend/common/testing/http/gock"
	"github.com/fwtpe/owl-backend/common/testing"

	"github.com/fwtpe/owl-backend/modules/query/g"
	"github.com/fwtpe/owl-backend/modules/query/http/boss"
	bmodel "github.com/fwtpe/owl-backend/modules/query/model/boss"
	db "github.com/fwtpe/owl-backend/modules/query/database"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("[syncHostData]", func() {
})

var _ = Describe("[syncIdcData]", skipBossDb.PrependBeforeEach(func() {
	var bossOrm orm.Ormer
	gockConfig := gock.GockConfigBuilder.NewConfigByRandom()

	Context("Basic testing on mocked API of BOSS", func() {
		BeforeEach(func() {
			/**
			 * Set-up environment
			 */
			apiConfig := &g.ApiConfig{
				Name:     "mock-77",
				Token:    "mock-token-77",
				BossBase: gockConfig.NewHttpConfig().Url,
				Map:      gockConfig.NewHttpConfig().Url + g.BOSS_URI_BASE_MAP,
			}
			g.SetConfig(&g.GlobalConfig{
				Api: apiConfig,
				Contacts: &g.ContactsConfig{
					Interval: 8,
				},
			})

			boss.SetPlugins(mock.Plugin)
			boss.SetupServerUrl(apiConfig)
			// :~)

			bossOrm = NewBossOrm()

			gockConfig.New().Get(
				fmt.Sprintf(
					g.BOSS_URI_BASE_MAP+ g.BOSS_IDC_PATH_TMPL,
					apiConfig.Name, boss.SecureFctoken(apiConfig.Token),
				),
			).
				Reply(http.StatusOK).
				JSON(&bmodel.IdcResult{
					Status: 1,
					Info:   "当前操作成功了！",
					Result: []*bmodel.IdcRow{
						{
							Platform: "idc-001",
							IpList: []*bmodel.IdcIp{
								{Ip: "19.20.6.1", Pop: "tt浙江-一", PopId: "601"},
								{Ip: "19.20.6.2", Pop: "tt浙江-二", PopId: "602"},
								{Ip: "19.20.6.3", Pop: "tt浙江-三", PopId: "603"},
							},
						},
						{
							Platform: "idc-002",
							IpList: []*bmodel.IdcIp{
								{Ip: "19.20.7.21", Pop: "tt北京-一", PopId: "611"},
								{Ip: "19.20.7.22", Pop: "tt北京-二", PopId: "612"},
								{Ip: "19.20.7.23", Pop: "tt北京-三", PopId: "613"},
							},
						},
					},
				})

			gockConfig.New().Times(6).Post(g.BOSS_URI_BASE_UPLINK).
				Reply(http.StatusOK).
				JSON(&bmodel.IdcBandwidthResult{
					Status: 1,
					Info:   "当前操作成功了！",
					Result: []*bmodel.IdcBandwidthRow{
						{UplinkTop: 120},
						{UplinkTop: 130},
						{UplinkTop: 100},
					},
				})

			gockConfig.New().Times(6).Post(g.BOSS_URI_BASE_GEO).
				Reply(http.StatusOK).
				JSON(&bmodel.LocationResult{
					Status: 1,
					Info:   "当前操作成功了！",
					Count:  3,
					Result: &bmodel.Location{
						Area:     "area-v1",
						Province: "廣西",
						City:     "city-v1",
					},
				})
		})
		AfterEach(func() {
			gockConfig.Off()

			bossOrm.Raw(
				`DELETE FROM idcs WHERE area = 'area-v1'`,
			).Exec()
		})

		It("The number of rows(\"idcs\") should be as expected", func() {
			syncIdcData()

			testedData := &struct {
				CountLocation int `db:"count_location"`
				CountIdcName int `db:"count_idc_name"`
				CountBandwidth int `db:"count_bandwidth"`
			} {}

			db.BossDbFacade.SqlxDbCtrl.Get(
				testedData,
				`
				SELECT
					(
						SELECT COUNT(*) FROM idcs
						WHERE area = 'area-v1' AND province = '廣西' AND city = 'city-v1'
					) AS count_location,
					(
						SELECT COUNT(*) FROM idcs
						WHERE idc LIKE 'tt%'
					) AS count_idc_name,
					(
						SELECT COUNT(*) FROM idcs
						WHERE bandwidth = 350 AND count = 0
					) AS count_bandwidth
				`,
			)

			Expect(testedData).To(PointTo(MatchAllFields(Fields{
				"CountLocation": Equal(6),
				"CountIdcName": Equal(6),
				"CountBandwidth": Equal(6),
			})))
		})
	})
	XContext("Testing on real BOSS", func() {
		SetupBossEnv()

		intervalSeconds := 8
		var timeBeforeStarted time.Time

		BeforeEach(func() {
			g.SetConfig(&g.GlobalConfig{
				Contacts: &g.ContactsConfig{
					Interval: intervalSeconds,
				},
			})
		})

		AfterEach(func() {
			bossInTx(`DELETE FROM idcs`)
		})

		syncAndAssert := func() {
			syncIdcData()

			testedResult := &struct {
				Count int `db:"count_of_all"`
				MaxUpdatedTime time.Time `db:"max_updated_time"`
			} {}

			db.BossDbFacade.SqlxDbCtrl.Get(
				testedResult,
				`
				SELECT COUNT(*) AS count_of_all, MAX(updated) AS max_updated_time
				FROM idcs
				`,
			)

			GinkgoT().Logf("IDCs on real BOSS: COUNT[%d]. Updated Time[%s]", testedResult.Count, testedResult.MaxUpdatedTime)

			Expect(testedResult).To(PointTo(MatchAllFields(Fields{
				"Count": BeNumerically(">", 0),
				"MaxUpdatedTime": BeTemporally(">=", timeBeforeStarted),
			})))
		}

		It("The number of IDCs should be viable", func() {
			timeBeforeStarted = time.Now()

			By("1st sync(insert)")
			syncAndAssert()

			By("2nd sync(do nothing)")
			syncAndAssert()

			timeBeforeStarted = timeBeforeStarted.Add(time.Duration(-intervalSeconds - 1) * time.Second)
			By("Set time to past(before interval)")
			db.BossDbFacade.SqlxDb.MustExec(
				`
				UPDATE idcs SET updated = FROM_UNIXTIME(?)
				`,
				timeBeforeStarted.Unix(),
			)

			By("3rd sync(updated)")
			syncAndAssert()
		})
	})
}))

var _ = Describe("Update host table of BOSS", skipBossDb.PrependBeforeEach(func() {
	var ormDb orm.Ormer

	BeforeEach(func() {
		ormDb = NewBossOrm()

		g.SetConfig(&g.GlobalConfig{
			Hosts: &g.HostsConfig{
				Interval: 8,
			},
		})
	})

	AfterEach(func() {
		_, err := ormDb.Raw(
			`
			DELETE FROM hosts
			WHERE hostname LIKE 'uh-host-%'
			`,
		).Exec()

		Expect(err).To(Succeed())

		g.SetConfig(nil)
		ormDb = nil
	})

	It("The number of inserted rows should be 3", func() {
		updateHostsTable(
			[]string{
				"uh-host-01", "uh-host-02", "uh-host-03",
			},
			map[string]map[string]string{
				"uh-host-01": {
					"hostname":  "uh-host-01",
					"isp":       "isp-1",
					"province":  "province-1",
					"city":      "city-1",
					"platforms": "p1,p2",
					"platform":  "p1",
					"activate":  "1",
				},
				"uh-host-02": {
					"hostname":  "uh-host-02",
					"isp":       "isp-2",
					"province":  "province-2",
					"city":      "city-2",
					"platforms": "p1,p2",
					"platform":  "p1",
					"activate":  "1",
				},
				"uh-host-03": {
					"hostname":  "uh-host-03",
					"isp":       "isp-3",
					"province":  "province-3",
					"city":      "city-3",
					"platforms": "p1,p2",
					"platform":  "p1",
					"activate":  "1",
				},
			},
		)

		var count int
		ormDb.Raw("SELECT COUNT(*) FROM hosts WHERE hostname LIKE 'uh-host-%'").
			QueryRow(&count)

		Expect(count).To(Equal(3))
	})
}))

var _ = Describe("Checking on passing of elapsed time for table \"idcs\"", skipBossDb.PrependBeforeEach(func() {
	Context("Empty table", func() {
		It("Should be passed", func() {
			testedResult := isElapsedTimePassedForIdcsTable(time.Now(), 0)
			Expect(testedResult).To(BeTrue())
		})
	})

	Context("Has data", func() {
		BeforeEach(func() {
			bossInTx(
				`
				INSERT INTO idcs(popid, idc, area, province, city, updated)
				VALUES
					(1301, 't01-cc1', 'ar1', 'pv-1', 'ct-1', '2015-10-02 10:19:30'),
					(1302, 't01-cc2', 'ar1', 'pv-1', 'ct-1', '2015-10-02 10:20:30')
				`,
			)
		})
		AfterEach(func() {
			bossInTx(
				`
				DELETE FROM idcs
				WHERE idc LIKE 't01-%'
				`,
			)
		})

		DescribeTable("Passed result should be as expected",
			func(time string, expected bool) {
				sampleTime := testing.ParseTimeByGinkgo(time)
				testedResult := isElapsedTimePassedForIdcsTable(sampleTime, 30)

				Expect(testedResult).To(Equal(expected))
			},
			Entry("Has not passed", "2015-10-02T10:21:00+08:00", false),
			Entry("Has passed", "2015-10-02T10:21:01+08:00", true),
			Entry("Has passed", "2015-10-03T07:35:24+08:00", true),
		)
	})
}))

var _ = Describe("Insert or Update table \"idcs\"", skipBossDb.PrependBeforeEach(func() {
	BeforeEach(func() {
		bossInTx(
			`
			INSERT INTO idcs(popid, idc, area, province, city)
			VALUES
				(2231, 'iuidc-ck1', 'ar1', 'pv-1', 'ct-1')
			`,
			`
			INSERT INTO hosts(hostname, ip, isp, idc)
			VALUES ('gp01.zc1.com', '87,6.55.1', 'isp-1', 'iuidc-jk1'),
				 ('gp02.zc1.com', '87,6.55.2', 'isp-1', 'iuidc-jk1'),
				 ('gp03.zc1.com', '87,6.55.3', 'isp-1', 'iuidc-jk1')
			`,
		)
	})
	AfterEach(func() {
		bossInTx(
			`
			DELETE FROM idcs
			WHERE idc LIKE 'iuidc-%'
			`,
			`
			DELETE FROM hosts
			WHERE idc LIKE 'iuidc-%'
			`,
		)
	})

	sampleData := map[string]*sourceIdcRow {
		"iuidc-ck1": {
			id: 22315, name: "iuidc-ck1",
			location: &bmodel.Location {
				Area: "au3", Province: "pu3", City: "cu3",
			},
			bandwidth: 100,
		},
		"iuidc-jk1": {
			id: 2232, name: "iuidc-jk1",
			location: &bmodel.Location {
				Area: "a3", Province: "p3", City: "c3",
			},
			bandwidth: 200,
		},
		"iuidc-jk2": {
			id: 2233, name: "iuidc-jk2",
			location: &bmodel.Location {
				Area: "a3", Province: "p3", City: "c3",
			},
			bandwidth: 300,
		},
	}

	It("The inserted/updated data should be matched", func() {
		updateIdcData(sampleData)

		result := &struct {
			InsertedCount int `db:"count_inserted"`
			UpdatedCount int `db:"count_updated"`
			CountForHosts int `db:"count_for_hosts"`
		} {}

		db.BossDbFacade.SqlxDbCtrl.Get(
			result,
			`
			SELECT
				(
					SELECT COUNT(*)
					FROM idcs
					WHERE idc LIKE 'iuidc-%'
				) AS count_inserted,
				(
					SELECT COUNT(*)
					FROM idcs
					WHERE idc = 'iuidc-ck1'
						AND popid = 22315 AND bandwidth = 100
						AND area = 'au3' AND province = 'pu3' AND city = 'cu3'
				) AS count_updated,
				(
					SELECT COUNT(*)
					FROM idcs
					WHERE idc = 'iuidc-jk1' AND count = 3
				) AS count_for_hosts
			`,
		)

		Expect(result).To(PointTo(MatchAllFields(Fields{
			"InsertedCount": Equal(3),
			"UpdatedCount": Equal(1),
			"CountForHosts": Equal(1),
		})))
	})
}))
