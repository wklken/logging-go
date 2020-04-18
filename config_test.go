package logging

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/wklken/logging-go/formatter"
)

func TestNewLogger(t *testing.T) {
	// text logger
	l := LogConfig{
		Level:          "debug",
		Format:         Text,
		FormatSettings: map[string]string{},
		Writer:         Discard,
		Hooks:          []LogHook{},
	}

	logger, err := l.NewLogger()
	assert.NoError(t, err)
	assert.NotNil(t, logger)

	// json logger
	l = LogConfig{
		Level:          "debug",
		Format:         JSON,
		FormatSettings: map[string]string{},
		Writer:         StdOut,
		Hooks:          []LogHook{},
	}

	logger, err = l.NewLogger()
	assert.NoError(t, err)
	assert.NotNil(t, logger)

	// apply as system default
	err = l.ApplyAsStdLogger()
	assert.NoError(t, err)
}

func TestLogConfigApplyGetWriter(t *testing.T) {
	var data = []struct {
		writer   LogWriter
		expected io.Writer
	}{
		{StdErr, os.Stderr},
		{StdOut, os.Stdout},
		{Discard, ioutil.Discard},
		{LogWriter("unknown"), os.Stderr},
	}
	for _, d := range data {
		c := LogConfig{Writer: d.writer}
		assert.Equal(t, d.expected, c.getWriter())
	}
}

func TestLogConfigApplyGetFormatter(t *testing.T) {
	var data = []struct {
		format   LogFormat
		expected log.Formatter
	}{

		{Text, &log.TextFormatter{}},
		{JSON, &formatter.JSONFormatter{}},
		{LogFormat("unknown"), &log.TextFormatter{}},
	}
	for _, d := range data {
		c := LogConfig{Format: d.format}
		assert.IsType(t, d.expected, c.getFormatter())
	}
}

func TestLogConfigApply(t *testing.T) {
	c := LogConfig{Level: "warning"}

	err := c.ApplyAsStdLogger()
	assert.NoError(t, err)

	assert.Equal(t, log.WarnLevel, log.GetLevel())

	c.Level = "unknown"
	err = c.ApplyAsStdLogger()
	assert.Error(t, err)
}

func TestLogHooksUnmarshalText(t *testing.T) {
	// setGlobalConfigEnv()
	// os.Setenv("LOG_HOOKS", `[{broken:json"]}`)

	// v := viper.New()
	// InitDefaults(v, "")

	// _, err := LoadConfigFromEnv(v)
	// assert.Error(t, err)

	lh := LogHooks{}
	err := lh.UnmarshalText([]byte(""))
	assert.Error(t, err)

	hooks := `[{"format":"logstash", "settings":{"type":"MyService","ts":"RFC3339Nano", "network": "udp",
	"host":"logstash.mycompany.io","port": "8911"}},{"format":"syslog","settings":{"network": "udp",
	"host":"localhost", "port": "514", "tag": "MyService", "facility": "LOG_LOCAL0", "severity": "LOG_INFO"}}]`

	err = lh.UnmarshalText([]byte(hooks))
	assert.NoError(t, err)
	// assert.Len(t, 2, len(lh))
}

func TestInitHooks(t *testing.T) {
	// init ErrUnknownLogHookFormat
	l := LogConfig{
		Level:          "debug",
		Format:         Text,
		FormatSettings: map[string]string{},
		Writer:         Discard,
		Hooks: []LogHook{
			{Type: "unknow", Settings: map[string]string{}},
		},
	}

	_, err := l.initHooks()
	assert.Error(t, err)

	// file, will init fail
	l = LogConfig{
		Level:          "debug",
		Format:         Text,
		FormatSettings: map[string]string{},
		Writer:         Discard,
		Hooks: []LogHook{
			{Type: "file", Settings: map[string]string{}},
		},
	}

	_, err = l.initHooks()
	assert.Error(t, err)

	// file, will init success
	l = LogConfig{
		Level:          "debug",
		Format:         Text,
		FormatSettings: map[string]string{},
		Writer:         Discard,
		Hooks: []LogHook{
			{Type: "file", Settings: map[string]string{"name": "test.log"}},
		},
	}

	hooks, err := l.initHooks()
	assert.NoError(t, err)
	assert.Len(t, hooks, 1)
}

func TestErrorArray(t *testing.T) {
	errs := Errors{
		errors.New("a"),
		errors.New("b"),
	}

	assert.Equal(t, "a b", errs.Error())
}
