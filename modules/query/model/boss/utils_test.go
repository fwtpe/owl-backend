package boss

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("[getIpFromHostnameWithDefault]", func() {
	DescribeTable("Checks the converted ip string as expected",
		func(sourceHostName string, defaultIp string, expectedResult string) {
			testedResult := GetIpFromHostnameWithDefault(sourceHostName, defaultIp)

			Expect(testedResult).To(Equal(expectedResult))
		},
		Entry("Normal", "bj-cnc-019-061-123-201", "", "19.61.123.201"),
		Entry("Cannot be parsed", "nothing", "98.77.6.1", "98.77.6.1"),
		Entry("Cannot be parsed(one of ip value)", "kz-abk-019-8c-123-201", "109.177.6.1", "109.177.6.1"),
	)
})

var _ = Describe("[getIspFromHostname]", func() {
	DescribeTable("Checks the converted ISP string as expected",
		func(sourceHostName string, expectedResult string) {
			testedResult := GetIspFromHostname(sourceHostName)
			Expect(testedResult).To(Equal(expectedResult))
		},
		Entry("Normal", "bjb-ck-091-111-041-35", "bjb"),
		Entry("Short name", "ack_zs_091_111", "ack"),
		Entry("Short name(with partial IP)", "cjc-zs-091-111.ball.com", "cjc"),
		Entry("Cannot be parsed", "nothing", ""),
		Entry("Cannot be parsed", "ll.091", ""),
	)
})
