package hook

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	AsyncEnable            = true
	DefaultAsyncBufferSize = 100000
	// DefaultAsyncBlock if block when the chan buffer is full(waitUntilBufferFrees), default will drop all the events
	DefaultAsyncBlock = false
)

type LogHookBuilder interface {
	New(name string, settings map[string]string) (logrus.Hook, error)
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

func getAsyncSettings(settings map[string]string) (bool, int, bool) {
	// async
	asyncEnable := AsyncEnable
	asyncBufferSize := DefaultAsyncBufferSize
	asyncBlock := DefaultAsyncBlock

	asyncEnableStr, ok := settings["async_enable"]
	if ok {
		if asyncEnableStr == "true" || asyncEnableStr == "1" {
			asyncEnable = true
		} else {
			asyncEnable = false
		}
	}

	if asyncBufferSizeStr, exists := settings["async_buffer_size"]; exists {
		size, cErr := strconv.Atoi(asyncBufferSizeStr)
		if cErr != nil {
			fmt.Println("async_buffer_size is not a valid integer, will use the default value")
		} else {
			asyncBufferSize = size
		}
	}

	asyncBlockStr, ok := settings["async_block"]
	if ok {
		if asyncBlockStr == "true" || asyncBlockStr == "1" {
			asyncBlock = true
		} else {
			asyncBlock = false
		}
	}

	return asyncEnable, asyncBufferSize, asyncBlock
}
