package proc

import (
	"fmt"
	"log"
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
	proc.OnEv(PEvProcBeforeStart, func(proc *process, work *Worker, others ...interface{}) {
		log.Println("process before start here")
	})
	proc.OnEv(PEvProcWhenGetSigToExit, func(proc *process, work *Worker, others ...interface{}) {
		log.Println(fmt.Sprintf("process get system signal: %v", others[0]))
	})
	proc.OnEv(PEvProcBeforeExit, func(proc *process, work *Worker, others ...interface{}) {
		log.Println("process before exit here")
	})
	proc.OnEv(PEvWorkerBeforeStart, func(proc *process, work *Worker, others ...interface{}) {
		log.Println(fmt.Sprintf("worker [%s] before start here", work.Name))
	})

	proc.NewWorker("worker1", worker1)
	proc.NewWorkers("worker2", worker2, 2)
	proc.RunInterval().Wait()
}
