package hook

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewFileHook(t *testing.T) {

	f := FileLogHook{}

	name := "test"
	emptySettings := map[string]string{}
	formatter := &logrus.JSONFormatter{}

	_, err := f.New(name, emptySettings, formatter)
	assert.Error(t, err)

	// normal
	settings := map[string]string{"name": "test"}
	_, err = f.New(name, settings, formatter)
	assert.NoError(t, err)

	// normal with path
	settings = map[string]string{"name": "test", "path": "/tmp"}
	_, err = f.New(name, settings, formatter)
	assert.NoError(t, err)

	// error, path not exists
	settings = map[string]string{"name": "test", "path": "/xxxx/tmp"}
	_, err = f.New(name, settings, formatter)
	assert.Error(t, err)

	// normal with keep
	settings = map[string]string{"name": "test", "path": "/tmp", "keep": "3"}
	_, err = f.New(name, settings, formatter)
	assert.NoError(t, err)

	// normal with wrong keep
	settings = map[string]string{"name": "test", "path": "/tmp", "keep": "aa"}
	_, err = f.New(name, settings, formatter)
	assert.Error(t, err)
}
