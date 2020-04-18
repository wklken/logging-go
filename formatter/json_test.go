package formatter

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// ! copy from logrus, change to testify

func TestErrorNotLost(t *testing.T) {
	t.Parallel()

	formatter := &JSONFormatter{}

	b, err := formatter.Format(logrus.WithField("error", errors.New("wild walrus")))
	assert.NoError(t, err)

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	assert.NoError(t, err)

	assert.Equal(t, entry["error"], "wild walrus")
}

func TestErrorNotLostOnFieldNotNamedError(t *testing.T) {
	t.Parallel()

	formatter := &JSONFormatter{}

	b, err := formatter.Format(logrus.WithField("omg", errors.New("wild walrus")))
	assert.NoError(t, err)

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	assert.NoError(t, err)

	assert.Equal(t, entry["omg"], "wild walrus")
}

func TestFieldClashWithTime(t *testing.T) {
	t.Parallel()

	formatter := &JSONFormatter{}

	b, err := formatter.Format(logrus.WithField("time", "right now!"))
	assert.NoError(t, err)

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	assert.NoError(t, err)

	assert.Equal(t, entry["fields.time"], "right now!")

	assert.Equal(t, entry["time"], "0001-01-01T00:00:00Z")
}

func TestFieldClashWithMsg(t *testing.T) {
	t.Parallel()

	formatter := &JSONFormatter{}

	b, err := formatter.Format(logrus.WithField("msg", "something"))
	assert.NoError(t, err)

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	assert.NoError(t, err)

	assert.Equal(t, entry["fields.msg"], "something")
}

func TestFieldClashWithLevel(t *testing.T) {
	t.Parallel()

	formatter := &JSONFormatter{}

	b, err := formatter.Format(logrus.WithField("level", "something"))
	assert.NoError(t, err)

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	assert.NoError(t, err)

	assert.Equal(t, entry["fields.level"], "something")
}

func TestFieldClashWithRemappedFields(t *testing.T) {
	t.Parallel()

	formatter := &JSONFormatter{
		FieldMap: FieldMap{
			FieldKeyTime:  "@timestamp",
			FieldKeyLevel: "@level",
			FieldKeyMsg:   "@message",
		},
	}

	b, err := formatter.Format(logrus.WithFields(logrus.Fields{
		"@timestamp": "@timestamp",
		"@level":     "@level",
		"@message":   "@message",
		"timestamp":  "timestamp",
		"level":      "level",
		"msg":        "msg",
	}))
	assert.NoError(t, err)

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	assert.NoError(t, err)

	for _, field := range []string{"timestamp", "level", "msg"} {
		assert.Equal(t, field, entry[field])

		remappedKey := fmt.Sprintf("fields.%s", field)
		assert.NotContains(t, entry, remappedKey)
	}

	for _, field := range []string{"@timestamp", "@level", "@message"} {
		assert.NotEqual(t, field, entry[field])

		remappedKey := fmt.Sprintf("fields.%s", field)

		assert.Contains(t, entry, remappedKey)
		assert.Equal(t, field, entry[remappedKey])
	}
}

func TestFieldsInNestedDictionary(t *testing.T) {
	t.Parallel()

	formatter := &JSONFormatter{
		DataKey: "args",
	}

	logEntry := logrus.WithFields(logrus.Fields{
		"level": "level",
		"test":  "test",
	})
	logEntry.Level = logrus.InfoLevel

	b, err := formatter.Format(logEntry)
	assert.NoError(t, err)

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	assert.NoError(t, err)

	args := entry["args"].(map[string]interface{})

	for _, field := range []string{"test", "level"} {
		assert.Contains(t, args, field)
		assert.Equal(t, args[field], field)
	}

	for _, field := range []string{"test", "fields.level"} {
		assert.NotContains(t, entry, field)
	}

	// with nested object, "level" shouldn't clash
	assert.Equal(t, "info", entry["level"])
}

func TestJSONEntryEndsWithNewline(t *testing.T) {
	t.Parallel()

	formatter := &JSONFormatter{}

	b, err := formatter.Format(logrus.WithField("level", "something"))
	assert.NoError(t, err)

	assert.Equal(t, "\n", string(b[len(b)-1]))
}

func TestJSONMessageKey(t *testing.T) {
	t.Parallel()

	formatter := &JSONFormatter{
		FieldMap: FieldMap{
			FieldKeyMsg: "message",
		},
	}

	b, err := formatter.Format(&logrus.Entry{Message: "oh hai"})
	assert.NoError(t, err)

	s := string(b)
	assert.Contains(t, s, "message")
	assert.Contains(t, s, "oh hai")
}

func TestJSONLevelKey(t *testing.T) {
	t.Parallel()

	formatter := &JSONFormatter{
		FieldMap: FieldMap{
			FieldKeyLevel: "somelevel",
		},
	}

	b, err := formatter.Format(logrus.WithField("level", "something"))
	assert.NoError(t, err)

	s := string(b)
	assert.Contains(t, s, "somelevel")
}

func TestJSONTimeKey(t *testing.T) {
	t.Parallel()

	formatter := &JSONFormatter{
		FieldMap: FieldMap{
			FieldKeyTime: "timeywimey",
		},
	}

	b, err := formatter.Format(logrus.WithField("level", "something"))
	assert.NoError(t, err)

	s := string(b)
	assert.Contains(t, s, "timeywimey")
}

func TestFieldDoesNotClashWithCaller(t *testing.T) {
	t.Parallel()

	logrus.SetReportCaller(false)
	formatter := &JSONFormatter{}

	b, err := formatter.Format(logrus.WithField("func", "howdy pardner"))
	assert.NoError(t, err)

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	assert.NoError(t, err)

	assert.Equal(t, entry["func"], "howdy pardner")
}

func TestFieldClashWithCaller(t *testing.T) {
	// t.Parallel()

	logrus.SetReportCaller(true)
	formatter := &JSONFormatter{}
	e := logrus.WithField("func", "howdy pardner")
	e.Caller = &runtime.Frame{Function: "somefunc"}
	b, err := formatter.Format(e)
	assert.NoError(t, err)

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	assert.NoError(t, err)

	assert.Equal(t, entry["fields.func"], "howdy pardner")

	assert.Equal(t, entry["func"], "somefunc")

	logrus.SetReportCaller(false) // return to default value
}

func TestJSONDisableTimestamp(t *testing.T) {
	t.Parallel()

	formatter := &JSONFormatter{
		DisableTimestamp: true,
	}

	b, err := formatter.Format(logrus.WithField("level", "something"))
	assert.NoError(t, err)

	s := string(b)
	assert.NotContains(t, s, FieldKeyTime)
}

func TestJSONEnableTimestamp(t *testing.T) {
	t.Parallel()

	formatter := &JSONFormatter{}

	b, err := formatter.Format(logrus.WithField("level", "something"))
	assert.NoError(t, err)

	s := string(b)
	assert.Contains(t, s, FieldKeyTime)
}
