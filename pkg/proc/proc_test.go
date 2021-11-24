package proc

import (
	"fmt"
	"testing"
	"time"
)

func worker1(w *Worker) {
	fmt.Println("worker111")
	w.SetNextInterval(1 * time.Second)
}

func worker2(w *Worker) {
	fmt.Println("worker222")
	w.SetNextInterval(1 * time.Second)
}

func TestProcess_Run(t *testing.T) {
	proc := New()
	proc.NewWorker("worker1", worker1)
	proc.NewWorkers("worker2", worker2, 2)

	proc.RunInterval().Wait()
}
