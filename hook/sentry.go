package hook

import (
	"github.com/evalphobia/logrus_sentry"
	"github.com/sirupsen/logrus"
)

type SentryLogHookBuilder struct {
}

// sentry: https://github.com/evalphobia/logrus_sentry
func (b SentryLogHookBuilder) New(name string, settings map[string]string) (logrus.Hook, error) {
	// 1. validate settings
	if err := validateRequiredHookSettings(name, settings, []string{"dsn"}); err != nil {
		return nil, err
	}

	asyncEnable := AsyncEnable
	asyncEnableStr, ok := settings["async_enable"]
	if ok {
		if asyncEnableStr == "true" || asyncEnableStr == "1" {
			asyncEnable = true
		} else {
			asyncEnable = false
		}
	}

	return newSentryHook(settings["dsn"], asyncEnable)
}

func newSentryHook(dsn string, asyncEnable bool) (logrus.Hook, error) {
	newSentryHookFunc := logrus_sentry.NewSentryHook
	if asyncEnable {
		newSentryHookFunc = logrus_sentry.NewAsyncSentryHook
	}

	hook, err := newSentryHookFunc(dsn, []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	})
	return hook, err
}
