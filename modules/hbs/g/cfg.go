package g

import (
	"encoding/json"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/toolkits/file"
)

type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type MysqlApiConfig struct {
	Host     string `json:"host"`
	Resource string `json:"resource"`
}

type GlobalConfig struct {
	Debug     bool            `json:"debug"`
	Hosts     string          `json:"hosts"`
	MaxIdle   int             `json:"maxIdle"`
	Listen    string          `json:"listen"`
	Trustable []string        `json:"trustable"`
	Http      *HttpConfig     `json:"http"`
	MysqlApi  *MysqlApiConfig `json:"mysql_api"`
}

type RpcView struct {
	Listen string `json:"listen"`
}

type FalconAgentView struct {
	Heartbeat *HeartbeatView `json:"heartbeat"`
}

type HeartbeatView struct {
	CurrentSize         int   `json:"current_size"`
	CumulativeDropped   int64 `json:"cumulative_dropped"`
	CumulativeReceived  int64 `json:"cumulative_received"`
	CumulativeProcessed int64 `json:"cumulative_processed"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	SetConfig(&c)

	log.Println("read config file:", cfg, "successfully")
}

func SetConfig(newConfig *GlobalConfig) {
	configLock.Lock()
	defer configLock.Unlock()

	config = newConfig
}
