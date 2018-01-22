package log

import (
	"strings"
	"regexp"

	lf "github.com/sirupsen/logrus"
)

var defaultTextFormatter = new(lf.TextFormatter)

type NameFormatter struct {
	name string
}

func (f *NameFormatter) Format(entry *lf.Entry) ([]byte, error) {
	entry.Data["module"] = f.name
	return defaultTextFormatter.Format(entry)
}

const maximumModuleName = 16

var componentSeperator, _ = regexp.Compile("[/.]")
func truncateModuleName(name string, maximumSize int) string {
	if len(name) <= maximumSize {
		return strings.Replace(name, ".", "/", -1)
	}

	nameComponents := componentSeperator.Split(name, -1)

	currentLength := len(nameComponents[len(nameComponents) - 1])
	// Keeps the full path of last component
	for i := len(nameComponents) - 2; i >= 0; i-- {
		if currentLength > maximumSize ||
			currentLength + len(nameComponents[i]) > maximumSize {
			nameComponents[i] = nameComponents[i][0:1]
		}

		currentLength += len(nameComponents[i])
	}

	return strings.Join(nameComponents, "/")
}
