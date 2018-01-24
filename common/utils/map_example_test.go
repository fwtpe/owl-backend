package utils

import (
	"fmt"
)

func ExampleAbstractMap_ToTypeOfTarget() {
	sampleMap := map[int16]string{
		20: "v1", 30: "v2",
	}

	targetMap := MakeAbstractMap(sampleMap).ToTypeOfTarget(
		int64(0), "string-type",
	).(map[int64]string)

	fmt.Printf("20[%s] 30[%s]", targetMap[20], targetMap[30])

	// Output:
	// 20[v1] 30[v2]
}

func ExampleAbstractMap_BatchProcess() {
	sampleMap := map[int16]string{
		1: "v1", 2: "v2", 3: "v3", 4: "v4", 5: "v5", 6: "v6",
		7: "v7", 8: "v8", 9: "v9",
	}

	batchTimes := 0
	restSize := 0

	MakeAbstractMap(sampleMap).BatchProcess(
		2,
		func(batch interface{}) {
			batchTimes++
		},
		func(rest interface{}) {
			restSize = len(rest.(map[int16]string))
		},
	)

	fmt.Printf("Batch times: %d. Rest size: %d", batchTimes, restSize)

	// Output:
	// Batch times: 4. Rest size: 1
}
