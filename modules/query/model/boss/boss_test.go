package boss

import (
	rd "github.com/Pallinder/go-randomdata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("PlatformDetail", func() {
	Context("ShortenDescription()", func() {
		longText := rd.RandStringRunes(201)

		DescribeTable("The result should be as expected one",
			func(source string, expectedResult string) {
				testedResult := (&PlatformDetail{Description: source}).ShortenDescription()

				Expect(testedResult).To(Equal(expectedResult))
			},
			Entry("Simple text", "This is a car", "This is a car"),
			Entry("Simple text(trimmed)", " True    color ", "True color"),
			Entry("Text contains linefeed and tab", "This\tis\ra\ncat", "This is a cat"),
			Entry("Text is more the 200 characters", longText, string([]rune(longText)[0:100])),
		)
	})
})
