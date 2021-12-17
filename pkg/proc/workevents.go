package proc

type wEvent int

const (
	WEventBeforeStart wEvent = iota
	WEventBeforeExit
)


func (w *Worker) OnEv(ev wEvent, fn func(w *Worker, others ...interface{})) *Worker {
	if w.events == nil {
		w.events = make(map[wEvent]func(w *Worker, others ...interface{}))
	}
	w.events[ev] = fn
	return w
}

func (w *Worker) OffEv(ev wEvent)  *Worker {
	if _, ok := w.events[ev]; ok {
		delete(w.events, ev)
	}
	return w
}

func (w *Worker) emitEv(ev wEvent, work *Worker, others ...interface{}) {
	if fn, ok := w.events[ev]; ok {
		fn(work, others...)
	}
}