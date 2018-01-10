package http

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego/orm"
	"gopkg.in/h2non/gock.v1"
	"gopkg.in/h2non/gentleman-mock.v2"

	"github.com/Cepave/open-falcon-backend/modules/query/http/boss"
	"github.com/Cepave/open-falcon-backend/modules/query/g"
	bmodel "github.com/Cepave/open-falcon-backend/modules/query/model/boss"

	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("[syncIDCsTable]", func() {
	var fakeUrl = "http://fake-sidct.com"
	var bossOrm orm.Ormer

	BeforeEach(func() {
		RegisterBossOrmOrSkip()

		/**
		 * Set-up environment
		 */
		mapUrl := fakeUrl + g.BOSS_URI_BASE_MAP
		apiConfig := &g.ApiConfig {
			Name: "mock-1",
			Token: "mock-token-1",
			BossBase: fakeUrl,
			Map: mapUrl,
		}
		g.SetConfig(&g.GlobalConfig {
			Api: apiConfig,
			Contacts: &g.ContactsConfig {
				Interval: 8,
			},
		})

		boss.SetPlugins(mock.Plugin)
		boss.SetupServerUrl(apiConfig)
		// :~)

		bossOrm = NewBossOrm()

		uri := fmt.Sprintf(
			"%s/fcname/mock-1/fctoken/%s/pop/yes/pop_id/yes.json",
			g.BOSS_URI_BASE_MAP, boss.SecureFctoken(apiConfig.Token),
		)

		gock.New(fakeUrl).Get(uri).
			Reply(http.StatusOK).
			JSON(&bmodel.IdcResult {
				Status: 1,
				Info: "当前操作成功了！",
				Result: []*bmodel.IdcRow {
					{
						Platform: "idc-001",
						IpList: []*bmodel.IdcIp {
							{ Ip: "19.20.6.1", Pop: "浙江-一", PopId: "601" },
							{ Ip: "19.20.6.2", Pop: "浙江-二", PopId: "602" },
							{ Ip: "19.20.6.3", Pop: "浙江-三", PopId: "603" },
						},
					},
					{
						Platform: "idc-002",
						IpList: []*bmodel.IdcIp {
							{ Ip: "32.120.77.1", Pop: "北京-一", PopId: "321" },
							{ Ip: "32.120.77.2", Pop: "北京-二", PopId: "322" },
							{ Ip: "32.120.77.3", Pop: "北京-三", PopId: "323" },
						},
					},
				},
			})

		gock.New(fakeUrl).Times(6).Post(g.BOSS_URI_BASE_UPLINK).
			Reply(http.StatusOK).
			JSON(&bmodel.IdcBandwidthResult{
				Status: 1,
				Info: "当前操作成功了！",
				Result: []*bmodel.IdcBandwidthRow {
					{ UplinkTop: 120 },
					{ UplinkTop: 130 },
					{ UplinkTop: 100 },
				},
			})

		gock.New(fakeUrl).Times(6).Post(g.BOSS_URI_BASE_GEO).
			Reply(http.StatusOK).
			JSON(&bmodel.LocationResult {
				Status: 1,
				Info: "当前操作成功了！",
				Count: 3,
				Result: &bmodel.Location {
					Area: "area-v1",
					Province: "廣西",
					City: "city-v1",
				},
			})
	})
	AfterEach(func() {
		gock.Off()

		bossOrm.Raw(
			`DELETE FROM idcs WHERE area = 'area-v1'`,
		).Exec()
	})

	It("The number of tables should be as expected", func() {
		syncIDCsTable()

		testedNumberOfRows := 0
		bossOrm.Raw(
			`SELECT COUNT(*) FROM idcs WHERE area = 'area-v1'`,
		).QueryRow(&testedNumberOfRows)

		Expect(testedNumberOfRows).To(Equal(6))
	})
})

var _ = Describe("Update host table of BOSS", skipBossDb.PrependBeforeEach(func() {
	var ormDb orm.Ormer

	BeforeEach(func() {
		RegisterBossOrmOrSkip()

		ormDb = NewBossOrm()

		g.SetConfig(&g.GlobalConfig {
			Hosts: &g.HostsConfig {
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
