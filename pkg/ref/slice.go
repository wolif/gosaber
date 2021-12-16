package ref

func (e *Entity) SliceGet(index int) (interface{}, bool) {
	if e.IsSlice() {
		if index >= e.GetValue().Len() {
			return nil, false
		}
		return e.GetValue().Index(index).Interface(), true
	}
	return nil, false
}

func (e *Entity) SliceLen() (int, bool) {
	if e.IsSlice() {
		return e.GetValue().Len(), true
	}
	return 0, false
}
