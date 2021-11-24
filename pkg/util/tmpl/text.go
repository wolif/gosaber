package tmpl

import (
	"bytes"
	"text/template"
)

func Text(name, tmpl string, params interface{}) (string, error) {
	t, err := template.New(name).Parse(tmpl)
	b := new(bytes.Buffer)
	err = t.Execute(b, params)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
