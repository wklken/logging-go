package hook

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewSentryHook(t *testing.T) {

	s := SentryLogHook{}

	name := "test"
	emptySettings := map[string]string{}
	formatter := &logrus.JSONFormatter{}

	_, err := s.New(name, emptySettings, formatter)
	assert.Error(t, err)

	// wrong dsn
	settings := map[string]string{
		"dsn": "xxxx",
	}
	_, err = s.New(name, settings, formatter)
	assert.Error(t, err)
	// t.Log(err)

	// TODO: normal sentry
}
