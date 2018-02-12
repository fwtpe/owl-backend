package log

import (
	"bytes"
	"encoding/json"
	cgin "github.com/fwtpe/owl-backend/common/gin"
	mvc "github.com/fwtpe/owl-backend/common/gin/mvc"
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

	RestfulLogger(group, mvcBuilder.BuildHandler)

	Context("/v1/list-all", func() {
		testApi := "/_/test/v1/list-all"
		BeforeEach(func() {
			defaultFactory = newLoggerFactory()
		})

		It("The loggers is empty when no logger exists", func() {
			req := httptest.NewRequest(http.MethodGet, testApi, nil)
			resp := httptest.NewRecorder()
			engine.ServeHTTP(resp, req)
			GinkgoT().Logf("Resp body: %s", resp.Body.String())

			expBody, _ := json.Marshal(NamedLoggerList{make([]*NamedLogger, 0)})
			Expect(resp.Body).To(MatchJSON(expBody))
		})

		It("The loggers contains the element when logger exists", func() {
			expName := "test/listall"
			GetLogger(expName)
			req := httptest.NewRequest(http.MethodGet, testApi, nil)
			resp := httptest.NewRecorder()
			engine.ServeHTTP(resp, req)
			GinkgoT().Logf("Resp body: %s", resp.Body.String())

			expLevel := DEFAULT_LEVEL.String()
			l := NamedLoggerList{make([]*NamedLogger, 0, 1)}
			l.Loggers = append(l.Loggers, &NamedLogger{expName, expLevel})
			expBody, _ := json.Marshal(l)

			Expect(resp.Body).To(MatchJSON(expBody))
		})
	})

	Context("/v1/set-level", func() {
		testApi := "/_/test/v1/set-level"
		testApiTree := testApi + "?tree=true"
		BeforeEach(func() {
			defaultFactory = newLoggerFactory()
		})

		It("Invalid logger name", func() {
			rawJson, _ := json.Marshal(NamedLogger{"test/setlevel-null", "debug"})
			req := httptest.NewRequest(http.MethodPost, testApi, bytes.NewReader(rawJson))
			resp := httptest.NewRecorder()
			engine.ServeHTTP(resp, req)
			GinkgoT().Logf("Resp body: %s", resp.Body.String())

			expBody, _ := json.Marshal(map[string]interface{}{
				"affected_loggers": 0,
			})
			Expect(resp.Body).To(MatchJSON(expBody))
		})

		It("Set 1 logger (default mode)", func() {
			name := "test/setlevel-default"
			GetLogger(name)
			rawJson, _ := json.Marshal(NamedLogger{name, "debug"})
			req := httptest.NewRequest(http.MethodPost, testApi, bytes.NewReader(rawJson))
			resp := httptest.NewRecorder()
			engine.ServeHTTP(resp, req)
			GinkgoT().Logf("Resp body: %s", resp.Body.String())

			expBody, _ := json.Marshal(map[string]interface{}{
				"affected_loggers": 1,
			})
			Expect(resp.Body).To(MatchJSON(expBody))
		})

		It("Set match loggers (tree mode)", func() {
			name := "test/setlevel-tree"
			GetLogger("not-match/setlevel-tree")
			GetLogger(name)
			GetLogger(name + "1")
			rawJson, _ := json.Marshal(NamedLogger{name, "debug"})
			req := httptest.NewRequest(http.MethodPost, testApiTree, bytes.NewReader(rawJson))
			resp := httptest.NewRecorder()
			engine.ServeHTTP(resp, req)
			GinkgoT().Logf("Resp body: %s", resp.Body.String())

			expBody, _ := json.Marshal(map[string]interface{}{
				"affected_loggers": 2,
			})
			Expect(resp.Body).To(MatchJSON(expBody))
		})
	})
})
