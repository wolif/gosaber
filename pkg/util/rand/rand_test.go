package rand

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRand(t *testing.T) {
	int31 := Int31Range(10, 20)
	assert.True(t, int31 >= 10)
	assert.True(t, int31 < 20)
	t.Log(int31)

	int63 := Int63Range(10, 20)
	assert.True(t, int63 >= 10)
	assert.True(t, int63 < 20)
	t.Log(int63)

	bytes, err := Bytes(3)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(bytes))
}
