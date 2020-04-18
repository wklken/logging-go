package hook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSentryHook(t *testing.T) {
	s := SentryLogHookBuilder{}

	name := "test"

	var data = []struct {
		settings  map[string]string
		willError bool
	}{
		{map[string]string{}, true},
		// wrong dsn
		{map[string]string{"dsn": "xxxx"}, true},
	}
	for _, d := range data {
		_, err := s.New(name, d.settings)
		if d.willError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
