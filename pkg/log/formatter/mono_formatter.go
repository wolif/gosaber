package formatter

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

type MonoFormatter struct {
	ServiceName string
}

func (m *MonoFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	//formatStr := "{$time} {$name} {$level} {$traceID} {$message} {$data}"
	traceID := ""
	if tID, ok := entry.Data["TraceID"]; ok {
		if sTID, ok := tID.(string); ok {
			traceID = sTID
		}
	}
	time := entry.Time.Format("2006-01-02 15:04:05")
	lv, err := entry.Level.MarshalText()
	if err != nil {
		return nil, err
	}
	msg := entry.Message
	data, err := json.Marshal(entry.Data)
	if err != nil {
		return nil, err
	}

	logMsg := fmt.Sprintf("%s %s %s %s \"%s\" %s\n", time, m.ServiceName, strings.ToUpper(string(lv)), traceID, msg, data)
	return []byte(logMsg), nil
}
