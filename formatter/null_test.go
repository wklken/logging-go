package formatter

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNullFormatter(t *testing.T) {
	f := NullFormatter{}

	entry := logrus.WithField("hello", "world")

	b, err := f.Format(entry)
	assert.NoError(t, err)

	assert.Equal(t, []byte(""), b)
}
