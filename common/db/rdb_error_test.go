package db

import (
	"errors"

	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("The caller information for database", func() {
	Context("Builds error by caller", func() {
		It("The error message should contain \"sampleDbWorker1\"", func() {
			testedError := sampleDbWorker1()

			Expect(testedError.Error()).To(And(
				ContainSubstring("SQL error: 87"),
				ContainSubstring("sampleDbWorker1"),
				ContainSubstring("rdb_error_test.go"),
			))
		})
	})

	Context("Error is panic!!", func() {
		It("The content of panic should contain sampleDbWorker1()", func() {
			testedError := samplePanic()

			Expect(testedError.Error()).To(And(
				ContainSubstring("SQL error: 71"),
				ContainSubstring("samplePanic"),
				ContainSubstring("rdb_error_test.go"),
			))
		})
	})
})

func sampleDbWorker1() error {
	return NewDatabaseError(errors.New("SQL error: 87"))
}

func samplePanic() (err *DbError) {
	defer func() {
		err = recover().(*DbError)
	}()

	PanicIfError(errors.New("SQL error: 71"))
	return
}
