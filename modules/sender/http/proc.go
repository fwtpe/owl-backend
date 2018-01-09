package http

import (
	"fmt"
	"github.com/fwtpe/owl-backend/modules/sender/proc"
	"net/http"
)

func configProcRoutes() {

	http.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("sms:%v, mail:%v, qq:%v", proc.GetSmsCount(), proc.GetMailCount(), proc.GetQQCount())))
	})

}
