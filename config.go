package logging

import (
	"io"
	"io/ioutil"
	"os"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/wklken/logging-go/formatter"
	"github.com/wklken/logging-go/hook"
)

// reference: https://github.com/hellofresh/logging-go/blob/master/hooks.go

var (
	// ErrUnknownLogHookFormat is the error returned when trying to initialize hook of unknown format
	ErrUnknownLogHookFormat = errors.New("failed to init log hooks: unknown hook found")

	// ErrMissingLogHookSetting is the error returned when trying to initialize hook with required settings missing
	// ErrMissingLogHookSetting = errors.New("Failed to init log hooks: missing required hook setting")
	// ErrFailedToConfigureLogHook is the error returned when hook configuring failed for some reasons
	// ErrFailedToConfigureLogHook = errors.New("Failed to init log hooks: failed to configure hook")
)

// LogFormat type for enumerating available log formats
type LogFormat string

// LogWriter for enumerating available log writers
type LogWriter string

const (
	// StdErr is os stderr log writer
	StdErr LogWriter = "stderr"
	// StdOut is os stdout log writer
	StdOut LogWriter = "stdout"
	// Discard is the quite mode for log writer aka /dev/null
	Discard LogWriter = "discard"

	// Text is plain text log format
	Text LogFormat = "text"
	// JSON is json log format
	JSON LogFormat = "json"
	// NULL is null log format
	Null LogFormat = "null"

	HookFile   = "file"
	HookSentry = "sentry"
	HookRedis  = "redis"
)

// LogHook is a struct holding settings for each enabled hook
type LogHook struct {
	// Format   string
	Type     string
	Settings map[string]string
}

// LogHooks is collection of enabled hooks
type LogHooks []LogHook

// UnmarshalText is an implementation of encoding.TextUnmarshaler for LogHooks type
func (lh *LogHooks) UnmarshalText(text []byte) error {
	var hooks []LogHook
	err := jsoniter.Unmarshal(text, &hooks)
	if nil != err {
		return err
	}

	*lh = hooks

	return nil
}

// LogConfig is the struct that stores all the logging configuration and routines for applying configurations
// to logger
type LogConfig struct {
	Level          string
	Format         LogFormat
	FormatSettings map[string]string
	Writer         LogWriter
	Hooks          LogHooks
}

func (c LogConfig) NewLogger() (*log.Logger, error) {
	var logger = log.New()

	level, err := log.ParseLevel(strings.ToLower(c.Level))
	if nil != err {
		return logger, err
	}
	logger.SetLevel(level)

	// logger.SetOutput(c.getWriter())
	logger.SetOutput(ioutil.Discard)
	// TODO: move the setFormat another place, default without any format
	// TODO: default NullFormatter
	logger.SetFormatter(c.getDefaultFormatter())
	// logger.SetFormatter(c.getFormatter())

	hooks, err := c.initHooks()
	if err != nil {
		log.WithError(err).Error("initHooks fail")
	}
	for _, hook := range hooks {
		logger.AddHook(hook)
	}
	return logger, nil
}

// Apply configures logger and all enabled hooks
func (c LogConfig) ApplyAsStdLogger() error {
	level, err := log.ParseLevel(strings.ToLower(c.Level))
	if nil != err {
		return err
	}
	log.SetLevel(level)

	log.SetOutput(c.getWriter())
	log.SetFormatter(c.getFormatter())

	hooks, err := c.initHooks()
	if err != nil {
		return err
	}
	for _, hook := range hooks {
		log.AddHook(hook)
	}
	return nil
}

func (c LogConfig) getWriter() io.Writer {
	switch c.Writer {
	case StdOut:
		return os.Stdout
	case Discard:
		return ioutil.Discard
	case StdErr:
		fallthrough
	default:
		return os.Stderr
	}
}

func (c LogConfig) getFormatter() log.Formatter {
	switch c.Format {
	case JSON:
		// return &log.JSONFormatter{}
		return &formatter.JSONFormatter{}
	case Null:
		return &formatter.NullFormatter{}
	case Text:
		fallthrough
	default:
		return &log.TextFormatter{}
	}
}

func (c LogConfig) getDefaultFormatter() log.Formatter {
	return &formatter.NullFormatter{}
}

type Errors []error

func (e Errors) Error() string {
	if e != nil && len(e) > 0 {
		messages := make([]string, 0, len(e))
		for _, err := range e {
			messages = append(messages, err.Error())
		}
		return strings.Join(messages, " ")
	}
	return ""
}

func (c LogConfig) initHooks() ([]log.Hook, error) {
	hooks := []log.Hook{}

	errs := Errors{}
	formatter := c.getFormatter()

	for _, h := range c.Hooks {
		var loghook hook.LogHookBuilder
		switch h.Type {
		case HookFile:
			loghook = hook.FileLogHookBuilder{Formatter: formatter}
		case HookSentry:
			loghook = hook.SentryLogHookBuilder{}
		case HookRedis:
			loghook = hook.RedisLogHookBuilder{}
		default:
			loghook = nil
		}

		// should match one of the types
		if loghook == nil {
			return nil, ErrUnknownLogHookFormat
		}

		lh, err := loghook.New(h.Type, h.Settings)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "init log hook %s fail", h.Type))
		} else {
			hooks = append(hooks, lh)
		}
	}

	if len(errs) != 0 {
		return hooks, errors.New(errs.Error())
	}

	return hooks, nil
}
