package rand

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFillStringWithZero(t *testing.T) {
	assert.Equal(t, "ABC", fillStringWithZero("ABC", 3))
	assert.Equal(t, "A00", fillStringWithZero("A", 3))
	assert.Equal(t, "BCD", fillStringWithZero("ABCD", 3))
}

func TestBytes(t *testing.T) {
	str, err := BytesAsBase16String(4)
	assert.Nil(t, err)
	t.Log(str)

	str, err = BytesAsBase36String(4)
	assert.Nil(t, err)
	t.Log(str)

	str, err = BytesAsBase64String(4)
	assert.Nil(t, err)
	t.Log(str)
}
