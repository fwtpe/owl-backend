package boss

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
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
	if err := bindJson(req.SendAndStatusMustMatch(http.StatusOK), idcBandwidthResult);
		err != nil {
		return nil
	}

	if idcBandwidthResult.Status != 1 {
		message := fmt.Sprintf(
			"Cannot load bandwidth data of IDC. Error Code[%d]. Message[%s]",
			idcBandwidthResult.Status, idcBandwidthResult.Info,
		)
		logger.Warn(message)
		panic(message)
	}

	logger.Debugf("Bandwidth: %v", idcBandwidthResult.Result)

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
	if err := bindJson(req.SendAndStatusMustMatch(http.StatusOK), locationResult);
		err != nil {
		return &bmodel.Location{}
	}

	if locationResult.Status != 1 {
		message := fmt.Sprintf(
			"Cannot load location data of Pop Id[%d]. Error Code[%d]. Message[%s]",
			idcId, locationResult.Status, locationResult.Info,
		)
		logger.Warn(message)
		panic(message)
	}

	logger.Debugf("Location data: %s", locationResult.Result)

	return locationResult.Result
}

func LoadPlatformData() []*bmodel.Platform {
	platformData := new(bmodel.PlatformResult)
	if err := bindJson(loadPlatformResp(), platformData);
		err != nil {
		return nil
	}

	return platformData.Result
}

func LoadPlatformDataAsMap(container *map[string]interface{}) {
	err := bindJson(loadPlatformResp(), container)

	if err != nil {
		panic(err)
	}
}

func loadPlatformResp() *gt.Response {
	uri := fmt.Sprintf(
		g.BOSS_PLATFORM_PATH_TMPL,
		bossFcname, SecureFctokenByConfig(),
	)

	req := client.ToGentlemanReq(
		bossMapBase.Clone().AddPath(uri),
	)

	return req.SendAndStatusMustMatch(http.StatusOK)
}

func LoadIdcData() []*bmodel.IdcRow {
	uri := fmt.Sprintf(
		g.BOSS_IDC_PATH_TMPL,
		bossFcname, SecureFctokenByConfig(),
	)
	req := client.ToGentlemanReq(
		bossMapBase.Clone().AddPath(uri),
	)

	idcResult := new(bmodel.IdcResult)
	if err := bindJson(req.SendAndStatusMustMatch(http.StatusOK), idcResult);
		err != nil {
		return nil
	}

	if idcResult.Status != 1 {
		message := fmt.Sprintf(
			"Cannot load IDC data. Error Code[%d]. Message[%s]",
			idcResult.Status, idcResult.Info,
		)
		logger.Warn(message)
		panic(message)
	}

	logger.Debugf("Number of IDC data: %d", len(idcResult.Result))

	return idcResult.Result
}

func bindJson(resp *gt.Response, container interface{}) error {
	jsonContent := resp.Bytes()

	if err := json.Unmarshal(jsonContent, container); err != nil {
		logger.Errorf("JSON cannot be unmarshalled to \"%T\". Error %v", container, err)
		logger.Errorf("JSON: %s", string(jsonContent))
		return err
	}

	return nil
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
