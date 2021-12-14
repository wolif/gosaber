package ref

func (e *Entry) StructTagGet(fieldName, tagName string) (tag string, ok bool) {
	field, ok := e.StructFieldGet(fieldName)
	if !ok {
		return
	}
	return field.Tag.Lookup(tagName)
}
