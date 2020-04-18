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

type FileLogHookBuilder struct {
	Formatter logrus.Formatter
}

// file hook : https://github.com/rifflock/lfshook
func (b FileLogHookBuilder) New(name string, settings map[string]string) (logrus.Hook, error) {
	// 1. validate settings
	if err := validateRequiredHookSettings(name, settings, []string{"name"}); err != nil {
		return nil, err
	}

	// path is not required
	path, ok := settings["path"]
	if ok {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, errors.New(fmt.Sprintf("file path %s not exists", path))
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

	asyncEnable, asyncBufferSize, asyncBlock := getAsyncSettings(settings)

	logPath := filename
	if path != "" {
		rawPath := strings.TrimSuffix(path, "/")
		logPath = fmt.Sprintf("%s/%s", rawPath, filename)
	}

	return newFileHook(logPath, keep, b.Formatter, asyncEnable, asyncBufferSize, asyncBlock)
}

type FileLogHook struct {
	fireChannel     chan *logrus.Entry
	asyncEnable     bool
	asyncBufferSize int
	asyncBlock      bool

	// loghook *logrus.Hook
	loghook *lfshook.LfsHook
}

func newFileHook(path string, keepDays int, formatter logrus.Formatter,
	asyncEnable bool, asyncBufferSize int, asyncBlock bool) (*FileLogHook, error) {
	// create rotate file hook
	rotateTime := 24 * time.Hour
	writer, err := rotatelogs.New(
		path+".%Y%m%d",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithRotationCount(uint(keepDays)),
		rotatelogs.WithRotationTime(rotateTime),
	)
	if err != nil {
		return nil, err
	}
	loghook := lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel:  writer,
			logrus.ErrorLevel: writer,
		},
		formatter,
	)

	// new fileloghook
	hook := &FileLogHook{
		loghook: loghook,
	}
	if asyncEnable {
		hook.asyncEnable = asyncEnable
		hook.asyncBufferSize = asyncBufferSize
		hook.asyncBlock = asyncBlock
		hook.makeAsync()
		fmt.Printf("init a async logger enable=%t buffer_size=%d, block=%t\n", asyncEnable, asyncBufferSize, asyncBlock)
	}

	return hook, nil
}

func (f *FileLogHook) makeAsync() {
	f.fireChannel = make(chan *logrus.Entry, f.asyncBufferSize)
	fmt.Printf("file hook will use a async buffer with size %d\n", f.asyncBufferSize)
	go func() {
		for entry := range f.fireChannel {
			if err := f.send(entry); err != nil {
				fmt.Println("Error during sending message to file:", err)
			}
		}
	}()
}

// Fire is called when a log event is fired.
func (f *FileLogHook) Fire(entry *logrus.Entry) error {
	if f.fireChannel != nil { // Async mode.
		select {
		case f.fireChannel <- entry: // try and put into chan, if fail will to default
		default:
			if f.asyncBlock {
				fmt.Println("the log buffered chan is full! will block")
				f.fireChannel <- entry // Blocks the goroutine because buffer is full.
				return nil
			}
			fmt.Println("the log buffered chan is full! will drop")
			// Drop message by default.
		}
		return nil
	}

	return f.send(entry)
}

func (f *FileLogHook) send(entry *logrus.Entry) error {
	return f.loghook.Fire(entry)
}

func (f *FileLogHook) Levels() []logrus.Level {
	return f.loghook.Levels()
}
