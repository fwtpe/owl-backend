// This package provides management of named loggers
//
// Named logger
//
// The name of logger is used to
//	1. Show the name of logger in output logging message
//	2. Sets the tree of logger(by name convention) with certain level
//
// Field of stack
//
// You could use "WithCurrentFrame(*Logger)" to add "frame=<Stack information>"
//
// 	WithCurrentFrame(logger).Errorf(...)
//
// Stack information
//
// Since it is unstable to depend on calling depth of 3-party libraries,
// the user of this package should use "common/runtime"."GetCurrentFunctionInfo()"
// to retrieve frame information of current function.
//
//	import rt "github.com/fwtpe/owl-backend/common/log"
//
//	logger.Infof("Something I want to tell you: %s", rt.GetCurrentFuncInfo())
package log

import (
	"io"
	"os"
	"sort"
	"strings"
	"sync"

	lf "github.com/sirupsen/logrus"

	or "github.com/fwtpe/owl-backend/common/runtime"
)

// Alias to "github.com/sirupsen/logrus"
var (
	AllLevels = lf.AllLevels
)

// Alias to "github.com/sirupsen/logrus"
const (
	PanicLevel = lf.PanicLevel
	FatalLevel = lf.FatalLevel
	ErrorLevel = lf.ErrorLevel
	WarnLevel  = lf.WarnLevel
	InfoLevel  = lf.InfoLevel
	DebugLevel = lf.DebugLevel

	ROOT_LOGGER = "_root_"
)

var defaultFactory = newLoggerFactory()

// Get a logger with name
func GetLogger(name string) *lf.Logger {
	return defaultFactory.GetLogger(name)
}

// List all of the named loggers
func ListAll() []*LoggerEntry {
	return defaultFactory.ListAll()
}

// Sets the level to the named logger
func SetLevel(name string, level lf.Level) {
	defaultFactory.SetLevel(name, level)
}

// Sets the level to tree of named loggers
//
// returns the number of matched loggers
func SetLevelToTree(name string, level lf.Level) int {
	return defaultFactory.SetLevelToTree(name, level)
}

func AddHook(name string, hook lf.Hook) {
	defaultFactory.AddHook(name, hook)
}
func AddHookToTree(name string, hook lf.Hook) int {
	return defaultFactory.AddHookToTree(name, hook)
}

func SetFormatter(name string, formatter lf.Formatter) {
	defaultFactory.SetFormatter(name, formatter)
}
func SetFormatterToTree(name string, formatter lf.Formatter) int {
	return defaultFactory.SetFormatterToTree(name, formatter)
}

func SetOut(name string, out io.Writer) {
	defaultFactory.SetOut(name, out)
}
func SetOutToTree(matchName string, out io.Writer) int {
	return defaultFactory.SetOutToTree(matchName, out)
}

// Add information of frame="<filename>:<line>:<package>"
func WithCurrentFrame(logger *lf.Logger) lf.FieldLogger {
	return logger.WithField("frame", or.GetCallerInfoWithDepth(1))
}

func newLoggerFactory() *loggerFactory {
	return &loggerFactory{
		namedLoggers: make(map[string]*lf.Logger),
		lock:         &sync.Mutex{},
	}
}

// Implementation for public functions
type loggerFactory struct {
	namedLoggers map[string]*lf.Logger
	lock         *sync.Mutex
}

func (l *loggerFactory) GetLogger(name string) *lf.Logger {
	l.lock.Lock()
	defer l.lock.Unlock()

	if _, ok := l.namedLoggers[name]; !ok {
		l.namedLoggers[name] = newLogger(name)
	}

	return l.namedLoggers[name]
}

type LoggerEntry struct {
	Logger *lf.Logger
	Name   string
}

func (l *loggerFactory) ListAll() []*LoggerEntry {
	loggers := make([]*LoggerEntry, 0)
	//

	for name, logger := range l.namedLoggers {
		loggers = append(loggers, &LoggerEntry{
			Logger: logger,
			Name:   name,
		})
	}
	by(name).Sort(loggers)
	return loggers
}

type by func(a, b *LoggerEntry) bool

func (b by) Sort(entries []*LoggerEntry) {
	s := &sorter{
		entries: entries,
		less:    b,
	}
	sort.Sort(s)
}

var name = func(a, b *LoggerEntry) bool {
	return a.Name < b.Name
}

type sorter struct {
	entries []*LoggerEntry
	less    by
}

func (s *sorter) Len() int {
	return len(s.entries)
}

func (s *sorter) Swap(i, j int) {
	s.entries[i], s.entries[j] = s.entries[j], s.entries[i]
}

func (s *sorter) Less(i, j int) bool {
	return s.less(s.entries[i], s.entries[j])
}

func (l *loggerFactory) SetLevel(name string, level lf.Level) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if logger, ok := l.namedLoggers[name]; ok {
		logger.SetLevel(level)
	}
}
func (l *loggerFactory) SetLevelToTree(matchName string, level lf.Level) int {
	l.lock.Lock()
	defer l.lock.Unlock()

	return l.iterateMatchLoggers(matchName, func(logger *lf.Logger) {
		logger.SetLevel(level)
	})
}

func (l *loggerFactory) AddHook(name string, hook lf.Hook) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if logger, ok := l.namedLoggers[name]; ok {
		logger.AddHook(hook)
	}
}
func (l *loggerFactory) AddHookToTree(matchName string, hook lf.Hook) int {
	l.lock.Lock()
	defer l.lock.Unlock()

	return l.iterateMatchLoggers(matchName, func(logger *lf.Logger) {
		logger.AddHook(hook)
	})
}

func (l *loggerFactory) SetFormatter(name string, formatter lf.Formatter) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if logger, ok := l.namedLoggers[name]; ok {
		logger.Formatter = formatter
	}
}
func (l *loggerFactory) SetFormatterToTree(matchName string, formatter lf.Formatter) int {
	l.lock.Lock()
	defer l.lock.Unlock()

	return l.iterateMatchLoggers(matchName, func(logger *lf.Logger) {
		logger.Formatter = formatter
	})
}

func (l *loggerFactory) SetOut(name string, out io.Writer) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if logger, ok := l.namedLoggers[name]; ok {
		logger.Out = out
	}
}
func (l *loggerFactory) SetOutToTree(matchName string, out io.Writer) int {
	l.lock.Lock()
	defer l.lock.Unlock()

	return l.iterateMatchLoggers(matchName, func(logger *lf.Logger) {
		logger.Out = out
	})
}

func (l *loggerFactory) iterateMatchLoggers(matchName string, callback func(*lf.Logger)) int {
	if matchName == ROOT_LOGGER {
		for _, logger := range l.namedLoggers {
			callback(logger)
		}

		return len(l.namedLoggers)
	}

	/**
	 * Checks for prefix
	 */
	count := 0
	for name, logger := range l.namedLoggers {
		if strings.HasPrefix(name, matchName) {
			callback(logger)
			count++
		}
	}
	// :~)

	return count
}

func toLogger(fieldLogger *lf.Logger) *lf.Logger {
	return interface{}(fieldLogger).(*lf.Logger)
}

func newLogger(name string) *lf.Logger {
	newLogger := lf.New()
	newLogger.Level = lf.WarnLevel
	newLogger.Formatter = NewTextFormatter(name)
	newLogger.Out = os.Stdout

	return newLogger
}
