package hook

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type LogHook interface {
	New(name string, settings map[string]string, formatter logrus.Formatter) (logrus.Hook, error)
}

var (
	ErrMissingLogHookSetting = errors.New("failed to init log hooks: missing required hook setting")
)

func validateRequiredHookSettings(name string, settings map[string]string, required []string) error {
	for i := range required {
		if _, ok := settings[required[i]]; !ok {
			logrus.WithFields(logrus.Fields{"logger_hook": name, "setting": required[i]}).Error("Missing required hook setting")
			return ErrMissingLogHookSetting
		}
	}
	return nil
}
