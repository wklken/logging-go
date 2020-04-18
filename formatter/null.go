package formatter

import (
	"github.com/sirupsen/logrus"
)

// NullFormatter formats logs into text
type NullFormatter struct {
}

// Format renders a single log entry
func (f *NullFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(""), nil
}
