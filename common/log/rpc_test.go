package log

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Struct method of log framework", func() {
	var testLogger = new(RpcLogger)

	BeforeEach(func() {
		defaultFactory = newLoggerFactory()
	})

	Context("ListAll", func() {
		expName := "struct/listall"
		var reply []*LoggerEntry
		id := func(element interface{}) string {
			return element.(*LoggerEntry).Name
		}

		It("The entry of loggers should be listed", func() {
			GetLogger(expName)
			err := testLogger.ListAll(nil, &reply)

			Expect(err).To(Succeed())
			Expect(len(reply)).To(Equal(1))
			Expect(reply).To(MatchAllElements(id, Elements{
				expName: Not(BeZero()),
			}))
		})
	})

	Context("SetLevel", func() {
		var (
			name        = "struct/setlevel"
			validArgs   = NamedLogger{name, "panic"}
			invalidArgs = NamedLogger{name, "invalid-level"}
		)

		BeforeEach(func() {
			defaultFactory = newLoggerFactory()
		})

		Context("When valid args are given", func() {
			It("A positive result should be returned if there is an exact match", func() {
				var reply bool
				GetLogger(name)
				err := testLogger.SetLevel(validArgs, &reply)

				Expect(err).To(Succeed())
				Expect(reply).To(BeTrue())
			})

			It("A negative result should be returned if there is no match", func() {
				var reply bool
				err := testLogger.SetLevel(validArgs, &reply)

				Expect(err).To(Succeed())
				Expect(reply).To(BeFalse())
			})
		})

		Context("When invalid args are given", func() {
			It("An error should have occurred", func() {
				var reply bool
				err := testLogger.SetLevel(invalidArgs, &reply)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("SetLevelToTree", func() {
		var (
			name        = "struct/totree"
			validArgs   = NamedLogger{name, "error"}
			invalidArgs = NamedLogger{name, "invalid-level"}
		)

		BeforeEach(func() {
			defaultFactory = newLoggerFactory()
		})

		Context("When valid args are given", func() {
			It("exact match count should be returned", func() {
				var reply int
				GetLogger(name)
				GetLogger(name + "1")
				GetLogger("struct/notmatch")
				err := testLogger.SetLevelToTree(validArgs, &reply)

				Expect(err).To(Succeed())
				Expect(reply).To(Equal(2))
			})
		})

		Context("When invalid args are given", func() {
			It("An error should have occurred", func() {
				var reply int
				err := testLogger.SetLevelToTree(invalidArgs, &reply)

				Expect(err).To(HaveOccurred())
			})
		})
	})

})
