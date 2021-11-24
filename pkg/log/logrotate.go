package log

import (
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/wolif/gosaber/pkg/log/formatter"
	"strings"
	"time"
)

func newLfsHook(config *Config) (logrus.Hook, error) {
	// rotation time
	var rotationTime time.Duration
	switch config.RotationTime {
	case RotationTimeHour:
		rotationTime = time.Hour
	case RotationTimeDay:
		rotationTime = 24 * time.Hour
	default:
		d, err := time.ParseDuration(config.RotationTime)
		if err != nil {
			return nil, err
		}
		rotationTime = d
	}

	// log suffix
	var logSuffix string
	switch {
	case rotationTime >= 24*time.Hour:
		logSuffix = "%Y-%m-%d"
	case rotationTime >= time.Hour:
		logSuffix = "%Y-%m-%d--%H"
	case rotationTime >= time.Minute:
		logSuffix = "%Y-%m-%d--%H:%M"
	default:
		panic("RotationTime too small")
	}

	writer, err := rotatelogs.New(
		fmt.Sprintf("%s.%s.log", strings.TrimSuffix(strings.TrimSpace(config.Output), ".log"), logSuffix),

		// WithLinkName为最新的日志建立软连接,以方便随着找到当前日志文件
		rotatelogs.WithLinkName(config.Output),

		// WithRotationTime设置日志分割的时间,这里设置为一小时分割一次
		rotatelogs.WithRotationTime(rotationTime),

		// WithMaxAge和WithRotationCount二者只能设置一个,
		// WithMaxAge设置文件清理前的最长保存时间,
		// WithRotationCount设置文件清理前最多保存的个数.
		//rotatelogs.WithMaxAge(time.Hour*24),
		rotatelogs.WithRotationCount(config.RotationCount),
	)

	if err != nil {
		logrus.Errorf("config local file system for logger error: %v", err)
	}

	// format
	var formater logrus.Formatter
	switch config.Format {
	case FormatText:
		formater = &logrus.TextFormatter{DisableColors: true}
	case FormatJson:
		formater = &logrus.JSONFormatter{}
	case FormatMono:
		formater = &formatter.MonoFormatter{ServiceName: config.ServiceName}
	}

	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, formater)

	return lfsHook, nil
}
