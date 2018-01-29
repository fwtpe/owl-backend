package log

import (
	lf "github.com/sirupsen/logrus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Factory functions for named loggers", func() {
	testLoggerFactory := newLoggerFactory()
	testLoggerFactory.GetLogger("m/mysqlapi")
	testLoggerFactory.GetLogger("m/mysqlapi/db")
	testLoggerFactory.GetLogger("m/mysqlapi/restful")
	testLoggerFactory.GetLogger("m/query")
	testLoggerFactory.GetLogger("m/query/db")

	Context("Get loggers", func() {
		It("The loggers with same name should be identical", func() {
			firstLogger := testLoggerFactory.GetLogger("m/mysqlapi")

			Expect(testLoggerFactory.GetLogger("m/mysqlapi")).To(BeIdenticalTo(firstLogger))
			Expect(testLoggerFactory.GetLogger("m/mysqlapi")).To(BeIdenticalTo(firstLogger))
		})
	})

	Context("Listing for all named loggers", func() {
		It("The map of logger should contains \"m/mysqlapi/...\"", func() {
			loggers := testLoggerFactory.ListAll()

			Expect(loggers).To(Equal([]*LoggerEntry{
				{testLoggerFactory.GetLogger("m/mysqlapi"), "m/mysqlapi"},
				{testLoggerFactory.GetLogger("m/mysqlapi/db"), "m/mysqlapi/db"},
				{testLoggerFactory.GetLogger("m/mysqlapi/restful"), "m/mysqlapi/restful"},
				{testLoggerFactory.GetLogger("m/query"), "m/query"},
				{testLoggerFactory.GetLogger("m/query/db"), "m/query/db"},
			}))
		})
	})
})

var _ = Describe("Output message of logger", func() {
	It("The message should contain module name[t/m/assert/content]", func() {
		testLoggerFactory := newLoggerFactory()
		logger := testLoggerFactory.GetLogger("test/message/assert/content")

		messageCatcher := new(catchMessageHook)

		testLoggerFactory.AddHook("test/message/assert/content", messageCatcher)

		logger.Warnf("[GBC11] Testing[%d] on logging(WARN level)", 195)

		Expect(messageCatcher.message).To(And(
			ContainSubstring("GBC11"),
			ContainSubstring("195"),
			ContainSubstring("t/m/assert/content"),
		))
	})

	It("The message should contain caller infomation: \"sampleCaller1\"", func() {
		messageCatcher := new(catchMessageHook)

		sampleCaller1(messageCatcher)

		Expect(messageCatcher.message).To(And(
			ContainSubstring("8891"),
			ContainSubstring("t/message/sample/cp1"),
			ContainSubstring("sampleCaller1"),
		))
	})
})

type catchMessageHook struct {
	message string
}

func (h *catchMessageHook) Levels() []lf.Level {
	return AllLevels
}
func (h *catchMessageHook) Fire(entry *lf.Entry) (err error) {
	h.message, err = entry.String()
	return
}

func sampleCaller1(messageCatcher *catchMessageHook) {
	logger := WithCurrentFrame(GetLogger("test/message/sample/cp1"))

	AddHook("test/message/sample/cp1", messageCatcher)

	logger.Warnln("This is 8891 testing")
}
