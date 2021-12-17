package proc

type pEvent int

const (
	PEvProcBeforeStart pEvent = iota
	PEvProcWhenGetSigToExit
	PEvProcBeforeExit
	PEvWorkerBeforeStart
	PEvWorkerBeforeExit
)

func (p *process) OnEv(ev pEvent, fn func(proc *process, work *Worker, others ...interface{})) *process {
	if p.events == nil {
		p.events = make(map[pEvent]func(proc *process, work *Worker, others ...interface{}))
	}
	p.events[ev] = fn
	return p
}

func (p *process) OffEv(ev pEvent) *process {
	if _, ok := p.events[ev]; ok {
		delete(p.events, ev)
	}
	return p
}

func (p *process) emitEv(ev pEvent, proc *process, work *Worker, others ...interface{}) {
	if fn, ok := p.events[ev]; ok {
		fn(proc, work, others...)
	}
}
