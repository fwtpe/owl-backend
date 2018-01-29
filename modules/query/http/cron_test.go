package http

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	osqlx "github.com/fwtpe/owl-backend/common/db/sqlx"
	ojson "github.com/fwtpe/owl-backend/common/json"
	"github.com/fwtpe/owl-backend/common/testing"
	"github.com/fwtpe/owl-backend/common/testing/http/gock"

	db "github.com/fwtpe/owl-backend/modules/query/database"
	"github.com/fwtpe/owl-backend/modules/query/g"
	"github.com/fwtpe/owl-backend/modules/query/http/boss"
	bmodel "github.com/fwtpe/owl-backend/modules/query/model/boss"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("[syncHostData] Mocked BOSS", skipBossDb.PrependBeforeEach(func() {
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
			Api:   apiConfig,
			Hosts: &g.HostsConfig{Interval: 8},
		})
		// :~)

		gockConfig.New().Get(
			fmt.Sprintf(
				g.BOSS_URI_BASE_MAP+g.BOSS_PLATFORM_PATH_TMPL,
				apiConfig.Name, boss.SecureFctoken(apiConfig.Token),
			),
		).
			Times(3).
			Reply(http.StatusOK).
			JSON(&bmodel.PlatformIpResult{
				Status: 1,
				Info:   "当前操作成功了！",
				Result: []*bmodel.PlatformIps{
					{
						Name: "c01.g01",
						IpList: []*bmodel.PlatformIp{
							{Hostname: "bj-cnc-097-006-001-041", Ip: "97.6.1.41", PopId: "141", Status: "1"},
							{Hostname: "bj-cnc-097-006-001-042", Ip: "97.6.1.42", PopId: "142", Status: "1"},
							{Hostname: "bj-cnc-097-006-001-043", Ip: "97.6.1.43", PopId: "143", Status: "1"},
						},
					},
					{
						Name: "c01.g02",
						IpList: []*bmodel.PlatformIp{
							{Hostname: "ac-fww-097-006-012-041", Ip: "97.6.12.41", PopId: "151", Status: "1"},
							{Hostname: "ac-fww-097-006-012-042", Ip: "97.6.12.42", PopId: "152", Status: "1"},
							{Hostname: "ac-fww-097-006-012-043", Ip: "97.6.12.43", PopId: "153", Status: "1"},
						},
					},
					{
						Name: "p01.k02",
						IpList: []*bmodel.PlatformIp{
							{Hostname: "bj-cnc-097-006-001-041", Ip: "97.6.1.41", PopId: "141", Status: "1"},
							{Hostname: "ac-fww-097-006-012-042", Ip: "97.6.12.42", PopId: "152", Status: "1"},
						},
					},
				},
			})

		gockConfig.New().Post(g.BOSS_URI_BASE_PLATFORM).
			Times(3).
			JSON(map[string]string{
				"fcname":  apiConfig.Name,
				"fctoken": boss.SecureFctoken(apiConfig.Token),
			}).
			Reply(http.StatusOK).
			JSON(&bmodel.PlatformDetailResult{
				Status: 1,
				Info:   "当前操作成功了！",
				Result: []*bmodel.PlatformDetail{
					{
						Name: "c01.g01", Type: "intel",
						Department: "dev-hs", Team: "main",
						Visible: "1", Description: "For Intel Evaluation",
					},
					{
						Name: "c01.g02", Type: "amd",
						Department: "dev-kc", Team: "minor-1",
						Visible: "1", Description: "For AMD Evaluation",
					},
				},
			})
	})

	AfterEach(func() {
		gockConfig.Off()

		bossInTx(
			`
			DELETE FROM ips
			WHERE ip LIKE '97.6.%'
			`,
			`
			DELETE FROM hosts
			WHERE hostname LIKE 'bj-%'
				OR hostname LIKE 'ac-%'
			`,
			`
			DELETE FROM platforms
			WHERE platform LIKE 'c01.g%'
			`,
		)
	})

	syncAndAssertData := func(checkedTime time.Time) {
		checkedTime = checkedTime.Add(-time.Second)

		syncHostData()

		testedResult := &struct {
			CountIps                    int `db:"count_ips"`
			CountHosts                  int `db:"count_hosts"`
			CountHostsForMultiPlatforms int `db:"count_hosts_multi_platforms"`
			CountPlatforms              int `db:"count_platforms"`
		}{}

		db.BossDbFacade.SqlxDbCtrl.Get(
			testedResult,
			`
			SELECT
				(
					SELECT COUNT(*) FROM ips
					WHERE ip LIKE '97.6.%'
						AND updated >= ?
				) AS count_ips,
				(
					SELECT COUNT(*) FROM hosts
					WHERE (hostname LIKE 'bj-%' OR hostname LIKE 'ac-%')
						AND updated >= ?
				) AS count_hosts,
				(
					SELECT COUNT(*) FROM hosts
					WHERE platform = 'p01.k02'
						AND platforms LIKE '%,p01.k02'
						AND updated >= ?
				) AS count_hosts_multi_platforms,
				(
					SELECT COUNT(*) FROM platforms
					WHERE platform LIKE 'c01.g%'
						AND count = 3
						AND updated >= ?
				) AS count_platforms
			`,
			checkedTime, checkedTime, checkedTime, checkedTime,
		)

		Expect(testedResult).To(PointTo(MatchAllFields(Fields{
			"CountIps":                    Equal(8),
			"CountHosts":                  Equal(6),
			"CountHostsForMultiPlatforms": Equal(2),
			"CountPlatforms":              Equal(2),
		})))
	}

	It(`The number of rows in "ips", "hosts", and "platforms" should have be (8, 6, 2)`, func() {
		By("New data")
		syncAndAssertData(time.Now())

		By("Update data(under interval)")
		syncAndAssertData(time.Now())

		By("Set time to past(before interval)")
		bossInTx(
			`
			UPDATE ips
			SET updated = NOW() - INTERVAL 24 HOUR
			WHERE ip LIKE '97.6.%'
			`,
			`
			UPDATE hosts
			SET updated = NOW() - INTERVAL 24 HOUR
			WHERE (hostname LIKE 'bj-%' OR hostname LIKE 'ac-%')
			`,
			`
			UPDATE platforms
			SET updated = NOW() - INTERVAL 24 HOUR
			WHERE platform LIKE 'c01.g%'
			`,
		)

		By("Update data again(exceeding interval)")
		syncAndAssertData(time.Now())
	})
}))

var _ = Describe("[syncIdcData] Mocked BOSS", skipBossDb.PrependBeforeEach(func() {
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
			Api:      apiConfig,
			Contacts: &g.ContactsConfig{Interval: 8},
		})
		// :~)

		gockConfig.New().Get(
			fmt.Sprintf(
				g.BOSS_URI_BASE_MAP+g.BOSS_IDC_PATH_TMPL,
				apiConfig.Name, boss.SecureFctoken(apiConfig.Token),
			),
		).
			Times(3).
			Reply(http.StatusOK).
			JSON(&bmodel.IdcResult{
				Status: 1,
				Info:   "当前操作成功了！",
				Result: []*bmodel.IdcIps{
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

		gockConfig.New().Times(18).Post(g.BOSS_URI_BASE_UPLINK).
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

		gockConfig.New().Times(18).Post(g.BOSS_URI_BASE_GEO).
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

		bossInTx(
			`DELETE FROM idcs WHERE area = 'area-v1'`,
		)
	})

	syncAndAssert := func(checkedTime time.Time) {
		checkedTime = checkedTime.Add(-time.Second)

		syncIdcData()

		testedData := &struct {
			CountLocation  int `db:"count_location"`
			CountIdcName   int `db:"count_idc_name"`
			CountBandwidth int `db:"count_bandwidth"`
		}{}

		db.BossDbFacade.SqlxDbCtrl.Get(
			testedData,
			`
			SELECT
				(
					SELECT COUNT(*) FROM idcs
					WHERE area = 'area-v1' AND province = '廣西' AND city = 'city-v1'
						AND updated >= ?
				) AS count_location,
				(
					SELECT COUNT(*) FROM idcs
					WHERE idc LIKE 'tt%'
						AND updated >= ?
				) AS count_idc_name,
				(
					SELECT COUNT(*) FROM idcs
					WHERE bandwidth = 350 AND count = 0
						AND updated >= ?
				) AS count_bandwidth
			`,
			checkedTime, checkedTime, checkedTime,
		)

		Expect(testedData).To(PointTo(MatchAllFields(Fields{
			"CountLocation":  Equal(6),
			"CountIdcName":   Equal(6),
			"CountBandwidth": Equal(6),
		})))
	}

	It("The number of rows(\"idcs\") should be 6", func() {
		By("1st sync(insert)")
		syncAndAssert(time.Now())

		By("2nd sync(do nothing)")
		syncAndAssert(time.Now())

		By("Set time to past(before interval)")
		db.BossDbFacade.SqlxDb.MustExec(
			`UPDATE idcs SET updated = NOW() - INTERVAL 24 HOUR`,
		)

		By("3rd sync(updated)")
		syncAndAssert(time.Now())
	})
}))

var _ = Describe("[syncContactData] Mocked BOSS", skipBossDb.PrependBeforeEach(func() {
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
			Api:      apiConfig,
			Contacts: &g.ContactsConfig{Interval: 8},
		})
		// :~)

		gockConfig.New().Post(g.BOSS_URI_BASE_CONTACT).
			JSON(map[string]interface{}{
				"fcname":       apiConfig.Name,
				"fctoken":      boss.SecureFctoken(apiConfig.Token),
				"platform_key": "mc01.tp01,mc01.tp02",
			}).
			Reply(http.StatusOK).
			JSON(&bmodel.PlatformContactResult{
				Status: 1,
				Info:   "当前操作成功了！",
				Result: map[string]*bmodel.ContactUsers{
					"mc01.tp01": {
						Principals: []*bmodel.ContactUser{
							{
								Id: "9981", RealName: "bok-1",
								CellPhoneNumber: "875-19081141", TelephoneNumber: "0988117651", Email: "bok-1@fw.com.going",
							},
						},
						Backupers: []*bmodel.ContactUser{
							{
								Id: "9982", RealName: "bok-2",
								CellPhoneNumber: "875-19081142", TelephoneNumber: "0988117652", Email: "bok-2@fw.com.going",
							},
						},
					},
					"mc01.tp02": {
						Principals: []*bmodel.ContactUser{
							{
								Id: "8041", RealName: "zs3-1",
								CellPhoneNumber: "875-22081141", TelephoneNumber: "0913117651", Email: "lukg33@fw.com.going",
							},
						},
						Backupers: []*bmodel.ContactUser{
							{
								Id: "8042", RealName: "zs3-2",
								CellPhoneNumber: "875-32081141", TelephoneNumber: "0983117441", Email: "lukg71@fw.com.going",
							},
						},
					},
				},
			})

		bossInTx(
			`
			INSERT INTO platforms(platform)
			VALUES ('mc01.tp01'), ('mc01.tp02')
			`,
		)
	})
	AfterEach(func() {
		gockConfig.Off()

		bossInTx(
			`
			DELETE FROM platforms
			WHERE platform LIKE 'mc01.%'
			`,
			`
			DELETE FROM contacts
			WHERE name LIKE 'bok%'
				OR name LIKE 'zs3%'
			`,
		)
	})

	It("The users of platform[mc01.tp01, mc01.tp02] should have [bok, zs3] users", func() {
		syncContactData()

		testedResult := &struct {
			CountTp01     int `db:"count_tp01"`
			CountContacts int `db:"count_contacts"`
		}{}

		db.BossDbFacade.SqlxDbCtrl.Get(
			testedResult,
			`
			SELECT
				(
					SELECT COUNT(*) FROM platforms
					WHERE platform = 'mc01.tp01'
						AND principal = 'bok-1' AND deputy = 'bok-2'
				) AS count_tp01,
				(
					SELECT COUNT(*) FROM contacts
					WHERE phone LIKE '875-%'
				) AS count_contacts
			`,
		)

		Expect(testedResult).To(PointTo(MatchAllFields(
			Fields{
				"CountTp01":     Equal(1),
				"CountContacts": Equal(4),
			},
		)))
	})
}))

var _ = Describe("[updateHostsTable]", skipBossDb.PrependBeforeEach(func() {
	BeforeEach(func() {
		g.SetConfig(&g.GlobalConfig{
			Hosts: &g.HostsConfig{Interval: 8},
		})

		bossInTx(
			`
			INSERT INTO hosts(
				ip, hostname, platform, platforms, idc,
				exist, activate,
				isp, province, city, updated
			)
			VALUES
				('61.103.22.56', 'gg-zhc-061-103-022-056', 'pc99', 'pc99,pc33', 'fwcc', 0, 0,
					'fwtt', 'p01', 'c01', NOW() - INTERVAL 1 MINUTE
				),
				('67.103.23.12', 'gg-zhc-067-103-023-012', 'pc99', 'pc99,pc33', 'fwcc', 1, 0,
					'fwtt', 'p01', 'c01', NOW() - INTERVAL 1 HOUR
				),
				('67.103.23.13', 'gg-zhc-067-103-023-013', 'pc99', 'pc99,pc33', 'fwcc', 1, 0,
					'fwtt', 'p01', 'c01', NOW() - INTERVAL 1 HOUR
				),
				('33.57.123.112', 'gg-zhc-033-057-123-112', 'za-3', 'za-3,az-45', 'idk1', 1, 0,
					'fwtt', 'p08', 'c08', NOW() - INTERVAL 1 MINUTE
				)
			`,
			`
			INSERT INTO idcs(popid, idc, province, city, area)
			VALUES
				(103, 'idc-3', 'pkz', '山東', 'a1'),
				(104, 'idc-21', 'p02', 'c02', 'a2'),
				(105, 'idc-gc', 'p03', 'c03', 'a3'),
				(106, 'fwcc', 'p04', 'c04', 'a4')
			`,
		)
	})

	AfterEach(func() {
		bossInTx(
			`
			DELETE FROM hosts
			WHERE hostname LIKE 'gg-zhc-%'
			`,
			`
			DELETE FROM idcs
			WHERE popid >= 103 AND popid <= 106
			`,
		)
	})

	It("The number of inserted/updated rows should be as expected", func() {
		updateHostsTable([]*bmodel.Host{
			{ // Updated data
				Ip: "61.103.22.57", Hostname: "gg-zhc-061-103-022-056",
				Isp: "gg", Platform: "pc33", Platforms: []string{"pc34", "pc35"},
				Activate: "1", IdcId: "104",
			},
			{ // Updated data(IDC cannot be found)
				Ip: "33.57.123.112", Hostname: "gg-zhc-033-057-123-112",
				Isp: "gg", Platform: "pc33", Platforms: []string{"pc34", "pc35"},
				Activate: "1", IdcId: "501",
			},
			/**
			 * New data
			 */
			{
				Ip: "61.93.22.33", Hostname: "gg-zhc-061-093-022-033",
				Isp: "gg", Platform: "np1", Platforms: []string{"np1", "g2"},
				Activate: "1", IdcId: "103",
			},
			{
				Ip: "61.93.22.34", Hostname: "gg-zhc-061-093-022-034",
				Isp: "gg", Platform: "np1", Platforms: []string{"np1", "g2"},
				Activate: "1", IdcId: "103",
			},
			// :~)
			/**
			 * New data(IDC cannot be found)
			 */
			{
				Ip: "33.57.123.113", Hostname: "gg-zhc-033-057-123-113",
				Isp: "gg", Platform: "np1", Platforms: []string{"np1", "g2"},
				Activate: "1", IdcId: "501",
			},
			// :~)
		})

		testedResult := &struct {
			CountInserted     int `db:"count_inserted"`
			CountUpdated      int `db:"count_updated"`
			CountTurnOffExist int `db:"count_turn_off_exist"`
			NullIdcInserted   int `db:"count_inserted_null_idc"`
			NullIdcUpdated    int `db:"count_updated_null_idc"`
		}{}

		db.BossDbFacade.SqlxDbCtrl.Get(
			testedResult,
			`
			SELECT
				(
					SELECT COUNT(*)
					FROM hosts
					WHERE ip LIKE '61.93.22%'
						AND platform = 'np1' AND platforms = 'np1,g2'
						AND isp = 'gg' AND exist = 1 AND activate = 1
						AND idc = 'idc-3' AND province = 'pkz' AND city = '山東'
				) AS count_inserted,
				(
					SELECT COUNT(*)
					FROM hosts
					WHERE ip = '61.103.22.57'
						AND platform = 'pc33' AND platforms = 'pc34,pc35'
						AND isp = 'gg' AND exist = 1 AND activate = 1
						AND idc = 'idc-21' AND province = 'p02' AND city = 'c02'
				) AS count_updated,
				(
					SELECT COUNT(*)
					FROM hosts
					WHERE exist = 0 AND ip LIKE '67.103.23.%'
				) AS count_turn_off_exist,
				(
					SELECT COUNT(*)
					FROM hosts
					WHERE ip = '33.57.123.113'
						AND idc IS NULL AND province IS NULL AND city IS NULL
				) AS count_inserted_null_idc,
				(
					SELECT COUNT(*)
					FROM hosts
					WHERE ip = '33.57.123.112'
						AND idc IS NULL AND province IS NULL AND city IS NULL
				) AS count_updated_null_idc
			`,
		)

		Expect(testedResult).To(PointTo(MatchAllFields(Fields{
			"CountInserted":     Equal(2),
			"CountUpdated":      Equal(1),
			"CountTurnOffExist": Equal(2),
			"NullIdcInserted":   Equal(1),
			"NullIdcUpdated":    Equal(1),
		})))

		idcCount := &struct {
			Count103 int `db:"count_103"`
			Count104 int `db:"count_104"`
			Count105 int `db:"count_105"`
			Count106 int `db:"count_106"`
		}{}

		db.BossDbFacade.SqlxDbCtrl.Get(
			idcCount,
			`
			SELECT
				(
					SELECT count FROM idcs WHERE popid = 103
				) AS count_103,
				(
					SELECT count FROM idcs WHERE popid = 104
				) AS count_104,
				(
					SELECT count FROM idcs WHERE popid = 105
				) AS count_105,
				(
					SELECT count FROM idcs WHERE popid = 106
				) AS count_106
			`,
		)

		Expect(idcCount).To(PointTo(MatchAllFields(Fields{
			"Count103": Equal(2),
			"Count104": Equal(1),
			"Count105": Equal(0),
			"Count106": Equal(2),
		})))
	})
}))

var _ = Describe("[updateIpsTable]", skipBossDb.PrependBeforeEach(func() {
	BeforeEach(func() {
		g.SetConfig(&g.GlobalConfig{
			Hosts: &g.HostsConfig{
				Interval: 8,
			},
		})
	})

	Context("Insert or Update existing data", func() {
		BeforeEach(func() {
			bossInTx(
				`
				INSERT INTO ips(id, ip, exist, status, type, hostname, platform, updated)
				VALUES
					(3301, '10.11.87.191', 0, 0, 'OK', 'bp-cnc-010-011-087-191', 'CNC-01', NOW() - INTERVAL 30 SECOND),
					(3303, '10.12.87.191', 1, 0, 'OK', 'bp-cnc-010-012-087-191', 'g01.y55', NOW() - INTERVAL 11 MINUTE),
					(3304, '10.12.87.192', 1, 0, 'OK', 'bp-cnc-010-012-087-192', 'g01.y55', NOW() - INTERVAL 11 MINUTE)
				`,
			)
		})
		AfterEach(func() {
			bossInTx(
				`
				DELETE FROM ips
				WHERE ip LIKE '10.1%'
				`,
			)
		})

		It("The inserted/updated content should be as expected", func() {
			updateIpsTable(
				[]*bmodel.PlatformIps{
					{
						Name: "CNC-01",
						IpList: []*bmodel.PlatformIp{
							{Ip: "10.11.87.191", Hostname: "ubp-cnc-010-011-087-193", Status: "0", Type: "UDD"},
						},
					},
					{
						Name: "CNC-02",
						IpList: []*bmodel.PlatformIp{
							{Ip: "10.11.87.51", Hostname: "bp-cnc-010-011-087-051", Status: "1", Type: "NDD"},
							{Ip: "10.11.87.52", Hostname: "bp-cnc-010-011-087-052", Status: "1", Type: "NDD"},
						},
					},
					{
						// Duplicated ip
						Name: "CNC-01",
						IpList: []*bmodel.PlatformIp{
							{Ip: "10.11.87.191", Hostname: "ubp-cnc-010-011-087-193", Status: "1", Type: "UDD"},
						},
					},
				},
			)

			testedResult := &struct {
				CountInserted          int `db:"count_inserted"`
				CountUpdated           int `db:"count_updated"`
				CountUpdatedToNotExist int `db:"count_updated_not_exist"`
			}{}

			db.BossDbFacade.SqlxDbCtrl.Get(
				testedResult,
				`
				SELECT
					(
						SELECT COUNT(*)
						FROM ips
						WHERE ip LIKE '10.11.87.%' AND hostname LIKE 'bp-cnc-010-011-087-05%'
							AND platform = 'CNC-02' AND status = 1 AND type = 'NDD'
					) AS count_inserted,
					(
						SELECT COUNT(*)
						FROM ips
						WHERE hostname = 'ubp-cnc-010-011-087-193' AND type = 'UDD'
							AND exist = 1 AND status = 1
					) AS count_updated,
					(
						SELECT COUNT(*)
						FROM ips
						WHERE ip LIKE '10.12.87.%'
							AND exist = 0
					) AS count_updated_not_exist
				`,
			)

			Expect(testedResult).To(PointTo(MatchAllFields(Fields{
				"CountInserted":          Equal(2),
				"CountUpdated":           Equal(1),
				"CountUpdatedToNotExist": Equal(2),
			})))
		})
	})
}))

var _ = Describe("[updatePlatformsTable]", skipBossDb.PrependBeforeEach(func() {
	BeforeEach(func() {
		bossInTx(
			`
			INSERT INTO ips(id, ip, exist, status, type, hostname, platform)
			VALUES
				(4501, '115.61.3.11', 1, 1, 'ok', 'gc-zhj-115-061-003-011', 'k01.u01'),
				(4502, '115.61.3.12', 0, 1, 'ok', 'gc-zhj-115-061-003-012', 'k01.u01'),
				(4503, '115.61.3.13', 0, 1, 'ok', 'gc-zhj-115-061-003-013', 'k01.u01'),
				(4504, '115.61.3.21', 1, 1, 'ok', 'gc-zhj-115-061-003-021', 'k01.r01')
			`,
			`
			INSERT INTO platforms(platform)
			VALUES ('k01.r01')
			`,
		)
	})
	AfterEach(func() {
		bossInTx(
			`
			DELETE FROM ips
			WHERE id >= 4501 AND id <= 4504
			`,
			`
			DELETE FROM platforms
			WHERE platform LIKE 'k01.%'
			`,
		)
	})

	Context("inserted/updated data should be as expected", func() {
		It("The number of inserted/updated rows should be as expected", func() {
			updatePlatformsTable([]*bmodel.PlatformDetail{
				{
					Name: "k01.u01", Type: "np-95", Visible: "1",
					Department: "dep-88", Team: "adg-tpe",
					Description: "desc-66",
				},
				{
					Name: "k01.r01", Type: "gc-95", Visible: "1",
					Department: "dep-70", Team: "it-spp",
					Description: "desc-33",
				},
			})

			testedResult := &struct {
				CountInserted int `db:"count_inserted"`
				CountUpdated  int `db:"count_updated"`
			}{}

			db.BossDbFacade.SqlxDbCtrl.Get(
				testedResult,
				`
				SELECT
					(
						SELECT COUNT(*)
						FROM platforms
						WHERE platform = 'k01.u01' AND type = 'np-95' AND visible = '1'
							AND department = 'dep-88' AND team = 'adg-tpe'
							AND description = 'desc-66'
							AND count = 1
					) AS count_inserted,
					(
						SELECT COUNT(*)
						FROM platforms
						WHERE platform = 'k01.r01' AND type = 'gc-95' AND visible = '1'
							AND department = 'dep-70' AND team = 'it-spp'
							AND description = 'desc-33'
							AND count = 1
					) AS count_updated
				`,
			)

			Expect(testedResult).To(PointTo(MatchAllFields(Fields{
				"CountInserted": Equal(1),
				"CountUpdated":  Equal(1),
			})))
		})
	})
}))

var _ = Describe("[updateIdcData]", skipBossDb.PrependBeforeEach(func() {
	BeforeEach(func() {
		bossInTx(
			`
			INSERT INTO idcs(popid, idc, area, province, city)
			VALUES
				(2231, 'iuidc-ck1', 'ar1', 'pv-1', 'ct-1')
			`,
		)
	})
	AfterEach(func() {
		bossInTx(
			`
			DELETE FROM idcs
			WHERE idc LIKE 'iuidc-%'
			`,
		)
	})

	sampleData := map[string]*sourceIdcRow{
		"iuidc-ck1": {
			id: 22315, name: "iuidc-ck1",
			location: &bmodel.Location{
				Area: "au3", Province: "pu3", City: "cu3",
			},
			bandwidth: 100,
		},
		"iuidc-jk1": {
			id: 2232, name: "iuidc-jk1",
			location: &bmodel.Location{
				Area: "a3", Province: "p3", City: "c3",
			},
			bandwidth: 200,
		},
		"iuidc-jk2": {
			id: 2233, name: "iuidc-jk2",
			location: &bmodel.Location{
				Area: "a3", Province: "p3", City: "c3",
			},
			bandwidth: 300,
		},
	}

	It("The inserted/updated data should be matched", func() {
		updateIdcData(sampleData)

		result := &struct {
			InsertedCount int `db:"count_inserted"`
			UpdatedCount  int `db:"count_updated"`
		}{}

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
				) AS count_updated
			`,
		)

		Expect(result).To(PointTo(MatchAllFields(Fields{
			"InsertedCount": Equal(3),
			"UpdatedCount":  Equal(1),
		})))
	})
}))

var _ = Describe("[loadDetailOfMatchedPlatforms]", func() {
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

		gockConfig.New().Post(g.BOSS_URI_BASE_PLATFORM).
			JSON(map[string]string{
				"fcname":  apiConfig.Name,
				"fctoken": boss.SecureFctoken(apiConfig.Token),
			}).
			Reply(http.StatusOK).
			JSON(&bmodel.PlatformDetailResult{
				Status: 1,
				Info:   "当前操作成功了！",
				Result: []*bmodel.PlatformDetail{
					{
						Name: "p01.y31", Type: "high-end",
						Department: "wuc", Team: "bjdev",
						Visible: "1", Description: "wuc is fine",
					},
					{
						Name: "p01.y32", Type: "high-end",
						Department: "wuc", Team: "bjdev-s1",
						Visible: "1", Description: "wuc is fine",
					},
					{
						Name: "p01.y33", Type: "low-end",
						Department: "gks", Team: "admin",
						Visible: "1", Description: "since 1981",
					},
					/**
					 * Non-imported platform
					 */
					{
						Name: "p01.y34", Type: "low-end",
						Department: "gks", Team: "admin",
						Visible: "1", Description: "since 1981",
					},
					{
						Name: "p01.y35", Type: "low-end",
						Department: "gks", Team: "admin",
						Visible: "1", Description: "since 1981",
					},
					// :~)
				},
			})
	})
	AfterEach(func() {
		gockConfig.Off()
	})

	Context("Loads detail of platforms from BOSS API", func() {
		It("The loaded data should be as expected one", func() {
			samplePlatforms := []*bmodel.PlatformIps{
				{Name: "p01.y31"},
				{Name: "p01.y32"},
				{Name: "p01.y33"},
			}
			testedDetailOfPlatforms := loadDetailOfMatchedPlatforms(samplePlatforms)

			GinkgoT().Logf("Platforms: %s", ojson.MarshalJSON(testedDetailOfPlatforms))

			Expect(testedDetailOfPlatforms).To(HaveLen(3))
			Expect(testedDetailOfPlatforms[1]).To(PointTo(MatchAllFields(Fields{
				"Name":        Equal("p01.y32"),
				"Type":        Equal("high-end"),
				"Department":  Equal("wuc"),
				"Team":        Equal("bjdev-s1"),
				"Visible":     Equal("1"),
				"Description": Equal("wuc is fine"),
			})))
		})
	})
})

var _ = Describe("[updateContactsTable]", skipBossDb.PrependBeforeEach(func() {
	BeforeEach(func() {
		bossInTx(
			`
			INSERT INTO contacts(name, phone, email)
			VALUES
				('adu-05', '98070101', 'adu05@gmail.com')
			`,
		)
	})
	AfterEach(func() {
		bossInTx(
			`
			DELETE FROM contacts
			WHERE name LIKE 'adu-%'
			`,
		)
	})

	It("The updated(1)/inserted(2) data of \"contacts\" should be as expected", func() {
		updateContactsTable(
			map[string]*bmodel.ContactUsers{
				"a01.k01": {
					Principals: []*bmodel.ContactUser{
						{
							RealName:        "adu-05",
							CellPhoneNumber: "0932-601601", Email: "adu31@nb.ak44.com",
						},
					},
					Backupers: []*bmodel.ContactUser{
						{
							RealName:        "adu-31",
							CellPhoneNumber: "0944-987001", Email: "adu77-1@dev.ak22.com",
						},
					},
				},
				"a01.k02": {
					Principals: []*bmodel.ContactUser{
						{
							RealName:        "adu-32",
							CellPhoneNumber: "0944-987002", Email: "adu77-2@dev.ak22.com",
						},
					},
				},
			},
		)

		testedResult := &struct {
			CountInserted int `db:"count_inserted"`
			CountUpdated  int `db:"count_updated"`
		}{}

		db.BossDbFacade.SqlxDbCtrl.Get(
			testedResult,
			`
			SELECT
				(
					SELECT COUNT(*) FROM contacts
					WHERE name LIKE 'adu-3%'
						AND phone LIKE '0944-%'
						AND email LIKE 'adu77-%'
				) AS count_inserted,
				(
					SELECT COUNT(*) FROM contacts
					WHERE name = 'adu-05'
						AND phone = '0932-601601'
						AND email = 'adu31@nb.ak44.com'
				) AS count_updated
			`,
		)

		Expect(testedResult).To(PointTo(MatchAllFields(
			Fields{
				"CountInserted": Equal(2),
				"CountUpdated":  Equal(1),
			},
		)))
	})
}))

var _ = Describe("[addContactsToPlatformsTable]", skipBossDb.PrependBeforeEach(func() {
	BeforeEach(func() {
		bossInTx(
			`
			INSERT INTO platforms(platform)
			VALUES ('pac-01'), ('pac-02'), ('pac-03')
			`,
		)
	})
	AfterEach(func() {
		bossInTx(
			`
			DELETE FROM platforms
			WHERE platform LIKE 'pac-%'
			`,
		)
	})

	It("The contact[3], principal[1], deputy[1] and upgrader[1] should be as expected", func() {
		addContactsToPlatformsTable(
			map[string]*bmodel.ContactUsers{
				"pac-01": {
					Principals: []*bmodel.ContactUser{{RealName: "pr01"}},
					Backupers:  []*bmodel.ContactUser{{RealName: "dp01"}},
					Upgraders:  []*bmodel.ContactUser{{RealName: "up01"}},
				},
				"pac-02": {},
				"pac-03": {
					Upgraders: []*bmodel.ContactUser{{RealName: "up03"}},
				},
			},
		)

		testedResult := &struct {
			Count01 int `db:"count_01"`
			Count02 int `db:"count_02"`
			Count03 int `db:"count_03"`
		}{}

		db.BossDbFacade.SqlxDbCtrl.Get(
			testedResult,
			`
			SELECT
				(
					SELECT COUNT(*) FROM platforms
					WHERE platform = 'pac-01'
						AND contacts = 'pr01,dp01,up01'
						AND principal = 'pr01'
						AND deputy = 'dp01'
						AND upgrader = 'up01'
				) AS count_01,
				(
					SELECT COUNT(*) FROM platforms
					WHERE platform = 'pac-02'
						AND contacts = ''
						AND principal = ''
						AND deputy = ''
						AND upgrader = ''
				) AS count_02,
				(
					SELECT COUNT(*) FROM platforms
					WHERE platform = 'pac-03'
						AND contacts = 'up03'
						AND principal = ''
						AND deputy = ''
						AND upgrader = 'up03'
				) AS count_03
			`,
		)

		Expect(testedResult).To(PointTo(MatchAllFields(
			Fields{
				"Count01": Equal(1),
				"Count02": Equal(1),
				"Count03": Equal(1),
			},
		)))
	})
}))

var _ = Describe("Checking on passing of elapsed time for customized table and name(idcs)", skipBossDb.PrependBeforeEach(func() {
	Context("Empty table", func() {
		It("Should be passed", func() {
			testedResult := isElapsedTimePassed("idcs", "updated", time.Now(), 0)
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
				testedResult := isElapsedTimePassed("idcs", "updated", sampleTime, 30)

				Expect(testedResult).To(Equal(expected))
			},
			Entry("Has not passed", "2015-10-02T10:21:00+08:00", false),
			Entry("Has passed", "2015-10-02T10:21:01+08:00", true),
			Entry("Has passed", "2015-10-03T07:35:24+08:00", true),
		)
	})
}))

var _ = XDescribe("Use online BOSS", func() {
	SetupBossEnv()

	var sqlxDb *osqlx.DbController
	BeforeEach(func() {
		sqlxDb = db.BossDbFacade.SqlxDbCtrl

		g.SetConfig(&g.GlobalConfig{
			Contacts: &g.ContactsConfig{Interval: 300},
			Hosts:    &g.HostsConfig{Interval: 300},
		})

		log.Level = logrus.WarnLevel
	})

	AfterEach(func() {
		bossInTx(
			`DELETE FROM idcs`,
			`DELETE FROM ips`,
			`DELETE FROM hosts`,
			`DELETE FROM platforms`,
			`DELETE FROM contacts`,
		)

		log.Level = logrus.DebugLevel
	})

	showRows := func(sql string, container interface{}) {
		rows := sqlxDb.QueryxExt(sql)
		defer rows.Close()

		for rows.Next() {
			rows.StructScan(container)
			GinkgoT().Logf("Row: %+v", container)
		}
	}

	Context("[syncHostData, syncHostData, syncContactData] online BOSS", func() {
		syncAndAssert := func() {
			syncIdcData()
			syncHostData()
			syncContactData()

			testedResult := &struct {
				CountIdcs      int `db:"count_idcs"`
				CountIps       int `db:"count_ips"`
				CountHosts     int `db:"count_hosts"`
				CountPlatforms int `db:"count_platforms"`
				CountContacts  int `db:"count_contacts"`
			}{}

			sqlxDb.Get(
				testedResult,
				`
				SELECT
					(
						SELECT COUNT(*) FROM idcs
					) AS count_idcs,
					(
						SELECT COUNT(*) FROM ips
					) AS count_ips,
					(
						SELECT COUNT(*) FROM hosts
					) AS count_hosts,
					(
						SELECT COUNT(*) FROM platforms
					) AS count_platforms,
					(
						SELECT COUNT(*) FROM contacts
					) AS count_contacts
				`,
			)

			GinkgoT().Logf(
				`Number of "idcs[%d]", "ips"[%d], "hosts"[%d], "platforms"[%d] and "contacts[%d]".`,
				testedResult.CountIdcs, testedResult.CountIps,
				testedResult.CountHosts, testedResult.CountPlatforms,
				testedResult.CountContacts,
			)

			Expect(testedResult).To(PointTo(MatchAllFields(Fields{
				"CountIdcs":      BeNumerically(">=", 1),
				"CountIps":       BeNumerically(">=", 1),
				"CountHosts":     BeNumerically(">=", 1),
				"CountPlatforms": BeNumerically(">=", 1),
				"CountContacts":  BeNumerically(">=", 1),
			})))

			GinkgoT().Logf("5 rows of \"idcs\"")
			showRows(
				`
				SELECT popid, idc, count, area, province, city, updated
				FROM idcs ORDER BY count DESC LIMIT 5
				`,
				&struct {
					PopId       int       `db:"popid"`
					Idc         string    `db:"idc"`
					Count       int       `db:"count"`
					Area        string    `db:"area"`
					Province    string    `db:"province"`
					City        string    `db:"city"`
					UpdatedTime time.Time `db:"updated"`
				}{},
			)

			GinkgoT().Logf("5 rows of \"platforms\" table")
			showRows(
				`
				SELECT platform, type, count, visible,
					department, team, description, updated
				FROM platforms ORDER BY count DESC LIMIT 5
				`,
				&struct {
					Platform    string    `db:"platform"`
					Type        string    `db:"type"`
					Count       int       `db:"count"`
					Visible     int       `db:"visible"`
					Department  string    `db:"department"`
					Team        string    `db:"team"`
					Description string    `db:"description"`
					UpdatedTime time.Time `db:"updated"`
				}{},
			)

			GinkgoT().Logf("5 rows of \"ips\" table")
			showRows(
				`
				SELECT hostname, ip, exist, status, type, updated, platform
				FROM ips
				WHERE exist = 1 AND status = 1
				ORDER BY ip ASC LIMIT 5
				`,
				&struct {
					Ip         string    `db:"ip"`
					Exist      int       `db:"exist"`
					Status     int       `db:"status"`
					Type       string    `db:"type"`
					Hostname   string    `db:"hostname"`
					Platform   string    `db:"platform"`
					UpdateTime time.Time `db:"updated"`
				}{},
			)

			GinkgoT().Logf("5 rows of \"hosts\" table")
			showRows(
				`
				SELECT ip, hostname, exist, activate,
					platform, platforms, isp, province, city, updated
				FROM hosts
				WHERE activate = 1 AND exist = 1
				ORDER BY idc DESC LIMIT 5
				`,
				&struct {
					Hostname    string         `db:"hostname"`
					Ip          string         `db:"ip"`
					Exist       int            `db:"exist"`
					Activate    int            `db:"activate"`
					Platform    string         `db:"platform"`
					Platforms   string         `db:"platforms"`
					Isp         sql.NullString `db:"isp"`
					Province    sql.NullString `db:"province"`
					City        sql.NullString `db:"city"`
					UpdatedTime time.Time      `db:"updated"`
				}{},
			)

			GinkgoT().Logf("5 rows of \"contacts\" table")
			showRows(
				`
				SELECT name, phone, email
				FROM contacts
				ORDER BY updated DESC LIMIT 5
				`,
				&struct {
					Name  string `db:"name"`
					Phone string `db:"phone"`
					Email string `db:"email"`
				}{},
			)
		}

		It("The number of idcs, platforms, ips, hosts, and contacts should >= 1", func() {
			By("1st sync(insert)")
			syncAndAssert()

			By("Set time to past(before interval)")
			bossInTx(
				`UPDATE idcs SET updated = NOW() - INTERVAL 24 HOUR`,
				`UPDATE ips SET updated = NOW() - INTERVAL 24 HOUR`,
				`UPDATE hosts SET updated = NOW() - INTERVAL 24 HOUR`,
				`UPDATE platforms SET updated = NOW() - INTERVAL 24 HOUR`,
				`UPDATE contacts SET updated = NOW() - INTERVAL 24 HOUR`,
			)

			By("2nd sync(updated)")
			syncAndAssert()
		})
	})
})
