package http

import (
	"encoding/json"
	"net/http"
	_ "net/http/pprof"

	"github.com/Cepave/open-falcon-backend/modules/hbs/rpc"

	commonGin "github.com/Cepave/open-falcon-backend/common/gin"
	log "github.com/Cepave/open-falcon-backend/common/logruslog"
	"github.com/Cepave/open-falcon-backend/modules/hbs/g"
	"github.com/gin-gonic/gin"
)

type Dto struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var logger = log.NewDefaultLogger("INFO")
var ginRouter *gin.Engine = nil
var GinConfig *commonGin.GinConfig = &commonGin.GinConfig{}

func init() {
	GinConfig.Mode = gin.ReleaseMode
	ginRouter = commonGin.NewDefaultJsonEngine(GinConfig)
	v1 := ginRouter.Group("/api/v1")
	v1.GET("/health", getHealth)

	configCommonRoutes(ginRouter)
	configProcRoutes(ginRouter)
}

func getHealth(c *gin.Context) {
	rpcInfo := g.Config().Listen
	httpInfo := g.Config().Http
	healthInfo := struct {
		Http        *g.HttpConfig      `json:"http"`
		Rpc         *g.RpcView         `json:"rpc"`
		FalconAgent *g.FalconAgentView `json:"falcon_agent"`
	}{
		httpInfo,
		&g.RpcView{rpcInfo},
		&g.FalconAgentView{
			&g.HeartbeatView{
				CurrentSize:         rpc.AgentHeartbeatService.CurrentSize(),
				CumulativeReceived:  rpc.AgentHeartbeatService.CumulativeAgentsPut(),
				CumulativeDropped:   rpc.AgentHeartbeatService.CumulativeAgentsDropped(),
				CumulativeProcessed: rpc.AgentHeartbeatService.CumulativeRowsAffected(),
			},
		},
	}
	c.JSON(http.StatusOK, healthInfo)
}

func RenderJson(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}

func RenderDataJson(w http.ResponseWriter, data interface{}) {
	RenderJson(w, Dto{Msg: "success", Data: data})
}

func RenderMsgJson(w http.ResponseWriter, msg string) {
	RenderJson(w, map[string]string{"msg": msg})
}

func AutoRender(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		RenderMsgJson(w, err.Error())
		return
	}
	RenderDataJson(w, data)
}

func Start() {
	if !g.Config().Http.Enabled {
		return
	}

	addr := g.Config().Http.Listen
	if addr == "" {
		return
	}
	s := &http.Server{
		Addr:           addr,
		Handler:        ginRouter,
		MaxHeaderBytes: 1 << 30,
	}

	logger.Println("http listening", addr)
	logger.Fatalln(s.ListenAndServe())
}
