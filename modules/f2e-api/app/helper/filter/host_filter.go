package filter

import (
	"fmt"

	gp "github.com/fwtpe/owl/modules/f2e-api/app/model/graph"
	con "github.com/fwtpe/owl/modules/f2e-api/config"
)

func HostFilter(filterTxt string, limit int) []gp.Endpoint {
	db := con.Con()
	res := []gp.Endpoint{}
	db.Graph.Model(&res).Where(fmt.Sprintf("endpoint regexp '%s'", filterTxt)).Limit(limit).Scan(&res)
	return res
}
