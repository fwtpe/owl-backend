package runtime

import (
	"fmt"
	"runtime"
	"time"

	rd "github.com/Pallinder/go-randomdata"

	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Gets information of caller", func() {
	Context("GetCallerInfo()", func() {
		It("The function should be \"sampleFunction1\"", func() {
			callerInfo := sampleFunction1()

			GinkgoT().Logf("CallerInfo: %s", callerInfo)
			Expect(callerInfo).To(PointTo(MatchFields(
				IgnoreExtras,
				Fields{
					"FileName": Equal("stack_test.go"),
					"PackageName": Equal("github.com/fwtpe/owl-backend/common/runtime"),
					"Line": BeNumerically(">=", 0),
					"FunctionName": Equal("sampleFunction1"),
				},
			)))
		})
	})

	Context("GetCurrentFuncInfo()", func() {
		It("The function should be \"sampleFunction2\"", func() {
			callerInfo := sampleFunction2()

			GinkgoT().Logf("CallerInfo: %s", callerInfo)
			Expect(callerInfo).To(PointTo(MatchFields(
				IgnoreExtras,
				Fields{
					"FileName": Equal("stack_test.go"),
					"PackageName": Equal("github.com/fwtpe/owl-backend/common/runtime"),
					"Line": BeNumerically(">=", 0),
					"FunctionName": Equal("sampleFunction2"),
				},
			)))
		})
	})

	Context("Gets stack of callers", func() {
		It("The stack should be \"stack3, stack2, and stack1\"", func() {
			testedStack := stack1()

			for i, stack := range testedStack {
				GinkgoT().Logf("Caller Stack[%d]: %+v", i, stack)
			}

			Expect(testedStack).To(HaveLen(3))
			Expect(testedStack[0].FunctionName).To(Equal("stack3"))
			Expect(testedStack[1].FunctionName).To(Equal("stack2"))
			Expect(testedStack[2].FunctionName).To(Equal("stack1"))
		})
	})
})

var _ = Describe("toCallerInfo(*runtime.Frame)", func() {
	Context("Without vendor", func() {
		It("\"sampleF1\" should be matched", func() {
			testedInfo := toCallerInfo(&runtime.Frame{
				Function: "github.com/some1/cool/utils.sampleF1",
				File: "/home/some1/go/src/github.com/some1/cool/utils/staff.go",
				Line: 30,
			})

			Expect(testedInfo).To(PointTo(MatchFields(
				IgnoreExtras,
				Fields {
					"PackageName": Equal("github.com/some1/cool/utils"),
					"FunctionName": Equal( "sampleF1"),
					"FileName": Equal( "staff.go"),
					"Line": Equal(30),
				},
			)))
		})
	})

	Context("With vendor", func() {
		It("\"doLeap\" should be matched", func() {
			testedInfo := toCallerInfo(&runtime.Frame{
				Function: "github.com/some2/vendor/github.com/some1/cleaning/editor.doLeap",
				File: "/home/some1/go/src/github.com/some2/vendor/github.com/some1/cleaning/editor/foot.go",
				Line: 30,
			})

			Expect(testedInfo).To(PointTo(MatchFields(
				IgnoreExtras,
				Fields {
					"PackageName": Equal("github.com/some1/cleaning/editor"),
					"FunctionName": Equal( "doLeap"),
					"FileName": Equal( "foot.go"),
					"Line": Equal(30),
				},
			)))
		})
	})
})

func stack1() CallerStack {
	return stack2()
}
func stack2() CallerStack {
	return stack3()
}
func stack3() CallerStack {
	return GetCallerInfoStack(0, 2)
}

func sampleFunction1() *CallerInfo {
	return infoRetriever()
}
func infoRetriever() *CallerInfo {
	return GetCallerInfo()
}
func sampleFunction2() *CallerInfo {
	return GetCurrentFuncInfo()
}

var _ = Describe("Benchmark \"toCallerInfo()\"", func() {
	countOfSamples := 0
	times := 6

	BeforeEach(func() {
		if countOfSamples == 0 {
			Skip("Skip because of \"countOfSamples == 0\"")
		}
	})

	Measure("The time used for convertion of \"*runtime.Frame\" to \"*CallerInfo\"", func(b Benchmarker) {

		b.Time("runtime", func() {
			zeroTimeCounter := 0

			for i := 0; i < countOfSamples; i++ {
				part1, part2, part3 := rd.Currency(), rd.LastName(), rd.Month()

				currentFrame := &runtime.Frame {
					Function: fmt.Sprintf("github.com/sample/car/%s/%s.%s", part1, part2, part3),
					File: fmt.Sprintf("C:/Code/go/src/github.com/sample/car/%s/%s.go", part1, part2),
					Line: rd.Number(5, 500),
				}

				beforeExecute := time.Now()
				testedInfo := toCallerInfo(currentFrame)
				timeDiff := time.Now().Sub(beforeExecute)
				if timeDiff > 0 {
					b.RecordValue("nanos/per toCallerInfo()", float64(timeDiff))
				} else {
					zeroTimeCounter++
				}

				Expect(testedInfo).To(PointTo(MatchFields(
					IgnoreExtras,
					Fields {
						"PackageName": HaveSuffix(fmt.Sprintf("%s/%s", part1, part2)),
						"FunctionName": Equal(part3),
						"FileName": Equal(part2 + ".go"),
						"Line": Equal(currentFrame.Line),
					},
				)))
			}

			b.RecordValue("zero-time counter", float64(zeroTimeCounter))
		})
	}, times)
})
