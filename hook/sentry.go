package hook

import (
	"github.com/evalphobia/logrus_sentry"
	"github.com/sirupsen/logrus"
)

type SentryLogHook struct {
}

// sentry: https://github.com/evalphobia/logrus_sentry
func (s SentryLogHook) New(name string, settings map[string]string, formatter logrus.Formatter) (logrus.Hook, error) {
	// 1. validate settings
	if err := validateRequiredHookSettings(name, settings, []string{"dsn"}); err != nil {
		return nil, err
	}

	return newSentryHook(settings["dsn"])
}

func newSentryHook(dsn string) (logrus.Hook, error) {
	hook, err := logrus_sentry.NewSentryHook(dsn, []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	})
	return hook, err
}
