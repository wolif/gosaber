package log

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestError(t *testing.T) {
	orgOp := "/tmp.log/logs/a.go.log       "

	t.Log(fmt.Sprintf("%s.%s.log", strings.TrimSuffix(strings.TrimSpace(orgOp), ".log"), time.Now().Format("2006-01-02")))
	t.Log(orgOp)
}
