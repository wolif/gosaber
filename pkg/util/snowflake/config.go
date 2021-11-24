package snowflake

import (
	"fmt"
	"github.com/wolif/gosaber/pkg/log"
)

var snowFlakeWorker *Worker

type Config struct {
	WorkerID int64 `json:"WorkerIDs"`
}

func Init(conf *Config) error {
	var workerID int64
	if conf != nil {
		workerID = conf.WorkerID
	}

	log.Infof("new snowflake worker id = %d", workerID)
	worker, err := NewWorker(
		workerID,
	)

	if err != nil {
		return err
	}

	snowFlakeWorker = worker
	return nil
}

func GenerateInt64Id() int64 {
	if snowFlakeWorker != nil {
		return snowFlakeWorker.GetId()
	}

	return 0
}

func GenerateHex() string {
	return fmt.Sprintf("%016x", GenerateInt64Id())
}
