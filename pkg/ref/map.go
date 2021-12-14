package ref

func (e *Entry) initVMap() (ok bool) {
	if !e.IsMap() {
		return false
	}
	if e.vMap != nil {
		return true
	}
	for _, k := range e.GetValue().MapKeys() {
		e.vMap[k.Interface()] = e.GetValue().MapIndex(k)
	}
	return true
}

func (e *Entry) MapKeys() (keys []interface{}, ok bool) {
	if ok := e.initVMap(); !ok {
		return nil, false
	}
	keys = make([]interface{}, 0, len(e.vMap))
	for k, _ := range e.vMap {
		keys = append(keys, k)
	}
	return keys, true
}

func (e *Entry) MapHas(key interface{}) bool {
	if ok := e.initVMap(); !ok {
		return false
	}
	_, ok := e.vMap[key]
	return ok
}

func (e *Entry) MapGet(key interface{}) (value interface{}, ok bool) {
	if e.initVMap() {
		value, ok = e.vMap[key]
		return
	}
	return
}

func (e *Entry) MapLen() (l int, ok bool) {
	if e.initVMap() {
		return len(e.vMap), true
	}
	return
}
