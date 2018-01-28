package log

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Truncate name of module", func() {
	DescribeTable("The truncated name should be as expected one",
		func(sampleName string, expectedName string) {
			testedName := truncateModuleName(sampleName, maximumModuleName)

			GinkgoT().Logf("Named[%s] to [%s]", sampleName, testedName)
			Expect(testedName).To(Equal(expectedName))
		},
		Entry("Nothing truncated", "common/db", "common/db"),
		Entry("Nothing truncated", "github.com/someone/vim-go/easy/ginkgo/utils", "g/c/s/v/easy/ginkgo/utils"),
		Entry("Nothing truncated", "common.db", "common/db"),
		Entry("Nothing truncated", "github.com.someone.vim-go.easy.ginkgo.utils", "g/c/s/v/easy/ginkgo/utils"),
	)
})
