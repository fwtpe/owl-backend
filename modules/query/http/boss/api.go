package boss

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	gt "gopkg.in/h2non/gentleman.v2"
	gtp "gopkg.in/h2non/gentleman.v2/plugin"

	"github.com/fwtpe/owl-backend/common/http/client"

	"github.com/fwtpe/owl-backend/modules/query/g"
	bmodel "github.com/fwtpe/owl-backend/modules/query/model/boss"
)

var bossFcname string
var bossPlainFctoken string

var bossContactBase *gt.Request
var bossEventBase *gt.Request
var bossMapBase *gt.Request
var bossGeoBase *gt.Request
var bossPlatformBase *gt.Request
var bossUplinkBase *gt.Request

var httpClientPlugins []gtp.Plugin

func SetupServerUrl(apiConfig *g.ApiConfig) {
	bossFcname = apiConfig.Name
	bossPlainFctoken = apiConfig.Token

	url := apiConfig.BossBase

	defaultClient := client.CommonGentleman.NewDefaultClient()
	for _, plugin := range httpClientPlugins {
		defaultClient.Use(plugin)
	}

	bossContactBase = defaultClient.BaseURL(url).Post().Path(
		g.BOSS_URI_BASE_CONTACT,
	)
	bossEventBase = defaultClient.BaseURL(url).Post().Path(
		g.BOSS_URI_BASE_EVENT,
	)
	bossMapBase = defaultClient.BaseURL(url).Get().Path(
		g.BOSS_URI_BASE_MAP,
	)
	bossGeoBase = defaultClient.BaseURL(url).Post().Path(
		g.BOSS_URI_BASE_GEO,
	)
	bossPlatformBase = defaultClient.BaseURL(url).Post().Path(
		g.BOSS_URI_BASE_PLATFORM,
	)
	bossUplinkBase = defaultClient.BaseURL(url).Post().Path(
		g.BOSS_URI_BASE_UPLINK,
	)
}

func SetPlugins(newPlugins ...gtp.Plugin) {
	httpClientPlugins = newPlugins
}

func LoadIdcBandwidth(idcName string) []*bmodel.IdcBandwidthRow {
	req := client.ToGentlemanReq(
		bossUplinkBase.Clone().JSON(
			map[string]interface{}{
				"fcname":   bossFcname,
				"fctoken":  SecureFctokenByConfig(),
				"pop_name": idcName,
			},
		),
	)

	idcBandwidthResult := &bmodel.IdcBandwidthResult{}
	client.ToGentlemanResp(req.SendAndStatusMustMatch(http.StatusOK)).
		MustBindJson(idcBandwidthResult)

	if idcBandwidthResult.Status != 1 {
		message := fmt.Sprintf(
			"Cannot load bandwidth data of IDC[POST /base/get_uplink_list]. Error Code[%d]. Message[%s]",
			idcBandwidthResult.Status, idcBandwidthResult.Info,
		)
		logger.Warn(message)
		panic(message)
	}

	return idcBandwidthResult.Result
}

func LoadLocationData(idcId int) *bmodel.Location {
	req := client.ToGentlemanReq(
		bossGeoBase.Clone().JSON(
			map[string]interface{}{
				"fcname":  bossFcname,
				"fctoken": SecureFctokenByConfig(),
				"pop_id":  idcId,
			},
		),
	)

	locationResult := new(bmodel.LocationResult)

	client.ToGentlemanResp(req.SendAndStatusMustMatch(http.StatusOK)).
		MustBindJson(locationResult)

	if locationResult.Status != 1 {
		message := fmt.Sprintf(
			"Cannot load location data of Pop Id[%d]. Error Code[%d]. Message[%s]",
			idcId, locationResult.Status, locationResult.Info,
		)
		logger.Warn(message)
		panic(message)
	}

	return locationResult.Result
}

func LoadIdcData() []*bmodel.IdcRow {
	uri := fmt.Sprintf(
		"/fcname/%s/fctoken/%s/pop/yes/pop_id/yes.json",
		bossFcname, SecureFctokenByConfig(),
	)
	req := client.ToGentlemanReq(
		bossMapBase.Clone().AddPath(uri),
	)

	idcResult := new(bmodel.IdcResult)
	client.ToGentlemanResp(req.SendAndStatusMustMatch(http.StatusOK)).
		MustBindJson(idcResult)

	if idcResult.Status != 1 {
		message := fmt.Sprintf(
			"Cannot load IDC data[GET /Base/platform/get_ip ... /pop/yes/pop_id/yes.json]. Error Code[%d]. Message[%s]",
			idcResult.Status, idcResult.Info,
		)
		logger.Warn(message)
		panic(message)
	}

	return idcResult.Result
}

func SecureFctokenByConfig() string {
	return SecureFctoken(bossPlainFctoken)
}
func SecureFctoken(plainToken string) string {
	firstPhase := md5.Sum([]byte(plainToken))

	timedValue := hex.EncodeToString(firstPhase[:md5.Size])
	timedValue = time.Now().Format("20060102") + timedValue

	secondPhase := md5.Sum([]byte(timedValue))

	return hex.EncodeToString(secondPhase[:md5.Size])
}
