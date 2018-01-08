package boss

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	gt "gopkg.in/h2non/gentleman.v2"

	"github.com/Cepave/open-falcon-backend/common/http/client"

	"github.com/Cepave/open-falcon-backend/modules/query/g"
)

var bossFcname string
var bossPlainFctoken string

var bossContactBase *gt.Request
var bossEventBase *gt.Request
var bossMapBase *gt.Request
var bossGeoBase *gt.Request
var bossPlatformBase *gt.Request
var bossUplinkBase *gt.Request

func SetupServerUrl(apiConfig *g.ApiConfig) {
	bossFcname = apiConfig.Name
	bossPlainFctoken = apiConfig.Token

	url := apiConfig.BossBase

	defaultClient := client.CommonGentleman.NewDefaultClient()

	bossContactBase = defaultClient.BaseURL(url).Post().AddPath(
		"/Base/platform/get_platform_linkman",
	)
	bossEventBase = defaultClient.BaseURL(url).Post().AddPath(
		"/Monitor/add_zabbix_event",
	)
	bossMapBase = defaultClient.BaseURL(url).Get().AddPath(
		"/Base/platform/get_ip",
	)
	bossGeoBase = defaultClient.BaseURL(url).Post().AddPath(
		"/pop/get_area",
	)
	bossPlatformBase = defaultClient.BaseURL(url).Post().AddPath(
		"/base/platform/get_all_platform_pbc",
	)
	bossUplinkBase = defaultClient.BaseURL(url).Post().AddPath(
		"/base/get_uplink_list",
	)
}

func SecureFctokenByConfig() string {
	return secureFctoken(bossPlainFctoken)
}
func secureFctoken(plainToken string) string {
	firstPhase := md5.Sum([]byte(plainToken))

	timedValue := hex.EncodeToString(firstPhase[:md5.Size])
	timedValue = time.Now().Format("20060102") + timedValue

	secondPhase := md5.Sum([]byte(timedValue))

	return hex.EncodeToString(secondPhase[:md5.Size])
}
