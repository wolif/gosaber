package ref

func (e *Entry) initStructValues() bool {
	if e.vStructValues != nil {
		return true
	}
	if !e.IsStruct() {
		return false
	}
	e.vStructValues = make(map[string]interface{}, e.GetValue().NumField())
	for i := 0; i < e.GetValue().NumField(); i++ {
		e.vStructValues[e.GetType().Field(i).Name] = e.GetValue().Field(i).Interface()
	}
	return true
}

func (e *Entry) StructValues() (map[string]interface{}, bool) {
	if !e.initStructValues() {
		return nil, false
	}
	return e.vStructValues, true
}

func (e *Entry) StructValueGet(name string) (value interface{}, ok bool) {
	if e.initStructValues() {
		value, ok = e.vStructValues[name]
		return
	}
	return nil, false
}