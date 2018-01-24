package utils

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Abstract Map", func() {
	Context("ToType() of different types for key and value", func() {
		type s2 string

		It("result map should be same as expected", func() {
			sampleAMap := MakeAbstractMap(map[int16]s2{
				1: "Nice",
				2: "Good",
			})

			testedMap := sampleAMap.ToTypeOfTarget(int32(0), "").(map[int32]string)
			Expect(testedMap).To(Equal(
				map[int32]string{
					1: "Nice",
					2: "Good",
				},
			))
		})
	})

	Context("Batch Processor", func() {
		sampleData := map[string]int{
			"v1": 10, "v2": 20, "v3": 30,
			"v4": 40, "v5": 50, "v6": 60,
			"v7": 70, "v8": 80, "v9": 90,
			"v10": 100, "v11": 110, "v12": 120,
		}

		DescribeTable("Batch should get expected number of maps",
			func(batchSize int) {
				testedBatchTimes := 0
				expectedRest := 0

				testedMap := MakeAbstractMap(sampleData)
				testedMap.BatchProcess(
					batchSize,
					func(data interface{}) {
						Expect(data).To(HaveLen(batchSize))
						testedBatchTimes++
					},
					func(data interface{}) {
						expectedRest = len(data.(map[string]int))
					},
				)

				Expect(testedBatchTimes).To(Equal(len(sampleData) / batchSize))
				Expect(expectedRest).To(Equal(len(sampleData) % batchSize))
			},
			Entry("Just fetched", 3),
			Entry("Have rest", 5),
			Entry("Not batch", 13),
		)

		Context("Empty map", func() {
			DescribeTable("Nothing happened",
				func(mapData interface{}) {
					testedMap := MakeAbstractMap(mapData)
					testedMap.BatchProcess(
						10,
						func(data interface{}) {
							Fail("Should not get called(batch)")
						},
						func(data interface{}) {
							Fail("Should not get called(rest of batch)")
						},
					)
				},
				Entry("Empty map", map[int]int{}),
				Entry("Nil map", map[int]int(nil)),
			)
		})
	})
})
