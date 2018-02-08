package log

import (
	"encoding/json"
	cgin "github.com/fwtpe/owl-backend/common/gin"
	mvc "github.com/fwtpe/owl-backend/common/gin/mvc"
	"github.com/fwtpe/owl-backend/common/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Restful API of log framework", func() {
	var (
		mvcBuilder = mvc.NewMvcBuilder(mvc.NewDefaultMvcConfig())
		engine     = cgin.NewDefaultJsonEngine(&cgin.GinConfig{Mode: gin.ReleaseMode})
		group      = engine.Group("/_/test")
	)

	RestLogger(group, mvcBuilder.BuildHandler)

	Context("/v1/list-all", func() {
		It("The loggers is empty when no logger exists", func() {
			req := httptest.NewRequest(http.MethodGet, "/_/test/v1/list-all", nil)
			resp := httptest.NewRecorder()
			engine.ServeHTTP(resp, req)
			GinkgoT().Logf(resp.Body.String())

			expBody, _ := json.Marshal(model.NamedLoggerList{make([]*model.NamedLogger, 0)})
			Expect(resp.Body).To(MatchJSON(expBody))
		})

		It("The loggers contains the element when logger exists", func() {
			expName := "test/listall"
			GetLogger(expName)
			req := httptest.NewRequest(http.MethodGet, "/_/test/v1/list-all", nil)
			resp := httptest.NewRecorder()
			engine.ServeHTTP(resp, req)
			GinkgoT().Logf(resp.Body.String())

			expLevel := DEFAULT_LEVEL.String()
			l := model.NamedLoggerList{make([]*model.NamedLogger, 0, 1)}
			l.Loggers = append(l.Loggers, &model.NamedLogger{expName, expLevel})
			expBody, _ := json.Marshal(l)
			Expect(resp.Body).To(MatchJSON(expBody))
		})
	})
})
