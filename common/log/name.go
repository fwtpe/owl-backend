package log

import (
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
