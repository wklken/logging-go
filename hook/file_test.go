package hook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFileHook(t *testing.T) {
	// formatter := logrus.JSONFormatter{}
	f := FileLogHookBuilder{}

	name := "test"

	var data = []struct {
		settings  map[string]string
		willError bool
	}{
		{map[string]string{}, true},
		// normal
		{map[string]string{"name": "test"}, false},
		// normal with path
		{map[string]string{"name": "test", "path": "/tmp"}, false},
		// error, path not exists
		{map[string]string{"name": "test", "path": "/xxxx/tmp"}, true},
		// normal with keep
		{map[string]string{"name": "test", "path": "/tmp", "keep": "3"}, false},
		// normal with wrong keep
		{map[string]string{"name": "test", "path": "/tmp", "keep": "aa"}, true},
	}
	for _, d := range data {
		_, err := f.New(name, d.settings)
		if d.willError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
