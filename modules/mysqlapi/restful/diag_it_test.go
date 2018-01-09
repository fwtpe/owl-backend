package restful

import (
	"net/http"

	json "github.com/fwtpe/owl/common/json"
	ogko "github.com/fwtpe/owl/common/testing/ginkgo"
	testingHttp "github.com/fwtpe/owl/common/testing/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test /health", itSkipOnPortal.PrependBeforeEach(func() {
	It("returns the JSON data", func() {
		resp := testingHttp.NewResponseResultBySling(
			httpClientConfig.NewClient().
				Get("health"),
		)
		jsonBody := resp.GetBodyAsJson()
		GinkgoT().Logf("[Mysql API Module Response] JSON Result:\n%s", json.MarshalPrettyJSON(jsonBody))
		Expect(resp).To(ogko.MatchHttpStatus(http.StatusOK))
	})
}))
