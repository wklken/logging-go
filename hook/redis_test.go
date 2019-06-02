package hook

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewRedisHook(t *testing.T) {

	f := RedisLogHook{}

	name := "test"
	emptySettings := map[string]string{}
	formatter := &logrus.JSONFormatter{}

	_, err := f.New(name, emptySettings, formatter)
	assert.Error(t, err)

	// normal
	settings := map[string]string{"host": "127.1.1.1", "port": "6379", "db": "0", "key": "test", "poolsize": "3"}
	_, err = f.New(name, settings, formatter)
	// assert.NoError(t, err)
	// TODO: will error, no redis available
	assert.Error(t, err)

	// normal wrong port
	settings = map[string]string{"host": "127.1.1.1", "port": "a", "db": "0", "key": "test"}
	_, err = f.New(name, settings, formatter)
	assert.Error(t, err)

	// normal wrong db
	settings = map[string]string{"host": "127.1.1.1", "port": "6379", "db": "a", "key": "test"}
	_, err = f.New(name, settings, formatter)
	assert.Error(t, err)

}

func TestRedisLogHookLevels(t *testing.T) {
	f := RedisLogHook{}

	assert.Len(t, f.Levels(), 7)
}

func TestRedisLogHookFire(t *testing.T) {

	// f := RedisLogHook{}
	// name := "test"
	// formatter := &logrus.JSONFormatter{}
	// settings := map[string]string{"host": "1.1.1.1", "port": "1234", "db": "0", "key": "test", "poolsize": "3"}
	// h, err := f.New(name, settings, formatter)
	// // assert.NoError(t, err)
	// entry := &logrus.Entry{
	// 	Message: "hello",
	// 	Level:   logrus.DebugLevel,
	// 	Time:    time.Now(),
	// 	Data:    logrus.Fields{"a": 1},
	// }

	// err = h.Fire(entry)
	// assert.Error(t, err)
}

func TestCreateMessage(t *testing.T) {

	entry := &logrus.Entry{
		Message: "hello",
		Level:   logrus.DebugLevel,
		Time:    time.Now(),
		Data:    logrus.Fields{"a": 1},
	}

	m1 := createMessage(entry)
	assert.Equal(t, "hello", m1["message"].(string))

	m2 := createV0Message(entry, "app1", "localhost")
	assert.Equal(t, "hello", m2["@message"])
	assert.Equal(t, "localhost", m2["@source_host"])
	assert.Len(t, m2["@fields"], 3)

	m3 := createV1Message(entry, "app1", "localhost")
	assert.Equal(t, "hello", m3["message"])
	assert.Equal(t, "localhost", m3["host"])
	assert.Equal(t, "app1", m3["application"])
	assert.Equal(t, 1, m3["a"])

}
