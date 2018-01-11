package http

import (
	"encoding/json"
	"net/http"
	_ "net/http/pprof"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/juju/errors"

	"github.com/fwtpe/owl-backend/modules/query/g"
)

type Dto struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewBossOrm() orm.Ormer {
	o := orm.NewOrm()
	o.Using("boss")

	return o
}

func Start() {
	if !g.Config().Http.Enabled {
		log.Warn("http.enabled is disabled in configuration")
		return
	}

	// config http routes
	configCommonRoutes()
	configProcHttpRoutes()
	configGraphRoutes()
	configAPIRoutes()
	configAlertRoutes()
	configGrafanaRoutes()
	configZabbixRoutes()
	configNqmRoutes()
	configNQMRoutes()

	// start mysql database
	if err := InitDatabase(); err != nil {
		log.Errorf("%s", errors.ErrorStack(err))
	}

	go SyncHostsAndContactsTable()

	// start http server
	addr := g.Config().Http.Listen
	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}

	log.Println("http.Start ok, listening on", addr)
	log.Fatalln(s.ListenAndServe())
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

func StdRender(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		w.WriteHeader(400)
		RenderMsgJson(w, err.Error())
		return
	}
	RenderJson(w, data)
}
