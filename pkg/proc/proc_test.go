package proc

import (
	"fmt"
	"testing"
	"time"
)

func worker1(w *Worker) {
	fmt.Println("worker111")
	w.SetInterval(1 * time.Second)
}

func worker2(w *Worker) {
	fmt.Println("worker222")
	w.SetInterval(1 * time.Second)
}

func TestProcess_Run(t *testing.T) {
	proc := New().DefaultEvent()
	proc.AddWorker("worker1", worker1)
	proc.AddWorker("worker2", worker2, 2)
	proc.RunInterval().Wait()
}
