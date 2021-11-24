package rand

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"github.com/martinlindhe/base36"
	"strings"
)

func BytesAsBase64String(bytes int) (string, error) {
	b, err := Bytes(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), err
}

func BytesAsBase36String(bytes int) (string, error) {
	b, err := Bytes(bytes)
	if err != nil {
		return "", err
	}
	return strings.ToLower(base36.EncodeBytes(b)), err
}

func BytesAsBase16String(bytes int) (string, error) {
	b, err := Bytes(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func fillStringWithZero(str string, width int) string {
	l := len(str)
	if l >= width {
		start := l - width
		return str[start:]
	}

	var buf bytes.Buffer
	buf.WriteString(str)

	for i := 0; i < width-l; i++ {
		buf.WriteString("0")
	}

	return buf.String()
}
