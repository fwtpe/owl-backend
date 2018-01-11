package test

import (
	"fmt"

	gk "github.com/onsi/ginkgo"

	tflag "github.com/fwtpe/owl-backend/common/testing/flag"

	"github.com/fwtpe/owl-backend/modules/query/g"
)

var props = tflag.NewTestFlags().GetViper()
var sampleFcname = props.GetString("boss.fcname")
var sampleFctoken = props.GetString("boss.fctoken")
var sampleBossUrl = props.GetString("boss.url")

func GetApiConfigByTestFlag() *g.ApiConfig {
	if sampleBossUrl == "" || sampleFcname == "" || sampleFctoken == "" {
		return nil
	}

	return &g.ApiConfig{
		BossBase: sampleBossUrl,
		Name:     sampleFcname,
		Token:    sampleFctoken,
		Contact:  fmt.Sprintf("%s%s", sampleBossUrl, g.BOSS_URI_BASE_CONTACT),
		Event:    fmt.Sprintf("%s%s", sampleBossUrl, g.BOSS_URI_BASE_EVENT),
		Geo:      fmt.Sprintf("%s%s", sampleBossUrl, g.BOSS_URI_BASE_GEO),
		Map:      fmt.Sprintf("%s%s", sampleBossUrl, g.BOSS_URI_BASE_MAP),
		Platform: fmt.Sprintf("%s%s", sampleBossUrl, g.BOSS_URI_BASE_PLATFORM),
		Uplink:   fmt.Sprintf("%s%s", sampleBossUrl, g.BOSS_URI_BASE_UPLINK),
	}
}

var BossSkipper = tflag.BuildSkipFactoryByBool(
	sampleBossUrl == "" || sampleFcname == "" || sampleFctoken == "",
	`Need "boss.fcname", "boss.fctoken", and "boss.url" in test properties`,
)

func SkipIfNoBossConfig() {
	if sampleBossUrl == "" || sampleFcname == "" || sampleFctoken == "" {
		gk.Skip(`Need "boss.fcname", "boss.fctoken", and "boss.url" in test properties`)
	}
}
