package http

import (
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/juju/errors"

	"github.com/fwtpe/owl-backend/common/logruslog"

	"github.com/fwtpe/owl-backend/modules/query/g"
)

var log = logruslog.NewDefaultLogger("WARN")

func InitDatabase() error {
	config := g.Config()
	// set default database
	//
	if err := orm.RegisterDataBase("default", "mysql", config.Db.Addr, config.Db.Idle, config.Db.Max); err != nil {
		return errors.Annotate(err, "init BeegoOrm for default database has error")
	}
	// register model
	orm.RegisterModel(new(Host), new(Grp), new(Grp_host), new(Grp_tpl), new(Plugin_dir), new(Tpl))

	// set grafana database
	strConn := strings.Replace(config.Db.Addr, "falcon_portal", "grafana", 1)

	if err := orm.RegisterDataBase("grafana", "mysql", strConn, config.Db.Idle, config.Db.Max); err != nil {
		return errors.Annotate(err, "init BeegoOrm for grafana database has error")
	}
	orm.RegisterModel(new(Province), new(City), new(Idc))

	if err := orm.RegisterDataBase("apollo", "mysql", config.ApolloDB.Addr, config.ApolloDB.Idle, config.ApolloDB.Max); err != nil {
		return errors.Annotate(err, "init BeegoOrm for apollo database has error")
	}

	if err := orm.RegisterDataBase("boss", "mysql", config.BossDB.Addr, config.BossDB.Idle, config.BossDB.Max); err != nil {
		return errors.Annotate(err, "init BeegoOrm for boss database has error")
	}

	orm.RegisterModel(new(Contacts), new(Hosts), new(Idcs), new(Ips), new(Platforms))
	if err := orm.RegisterDataBase("gz_nqm", "mysql", config.Nqm.Addr, config.Nqm.Idle, config.Nqm.Max); err != nil {
		return errors.Annotate(err, "init BeegoOrm for gz_nqm database has error")
	}

	orm.RegisterModel(new(Nqm_node))

	if config.Debug == true {
		orm.Debug = true
	}

	return nil
}
