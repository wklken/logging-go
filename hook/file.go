package hook

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type FileLogHook struct {
}

// file hook : https://github.com/rifflock/lfshook
func (f FileLogHook) New(name string, settings map[string]string, formatter logrus.Formatter) (logrus.Hook, error) {
	// 1. validate settings
	if err := validateRequiredHookSettings(name, settings, []string{"name"}); err != nil {
		return nil, err
	}

	// path is not required
	path, ok := settings["path"]
	if ok {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, errors.New("file path not exists")
		}
		// TODO: create the path or error
	} else {
		path = ""
	}

	filename := settings["name"]

	keep := 7
	keepStr, ok := settings["keep"]
	if ok {
		keepInt, err := strconv.Atoi(keepStr)
		if err != nil {
			return nil, errors.New("keep should be integer")
		}
		keep = keepInt
	}

	logPath := filename
	if path != "" {
		rawPath := strings.TrimSuffix(path, "/")
		logPath = fmt.Sprintf("%s/%s", rawPath, filename)
	}

	return newFileHook(logPath, keep, formatter)
}

func newFileHook(path string, keepDays int, formatter logrus.Formatter) (logrus.Hook, error) {
	// path+".%Y%m%d%H%M",
	rotateTime := time.Duration(keepDays*24*60*60) * time.Second
	maxAge := time.Duration(1*24*60*60) * time.Second
	writer, err := rotatelogs.New(
		path+".%Y%m%d",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(maxAge),
		rotatelogs.WithRotationTime(rotateTime),
	)
	if err != nil {
		return nil, err
	}

	hook := lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel:  writer,
			logrus.ErrorLevel: writer,
		},
		formatter,
	)

	return hook, nil

}
