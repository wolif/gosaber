package frequency

import (
	"testing"
	"time"
)

func TestFrequency_Ctrl_Sync(t *testing.T) {
	c := make(chan string)
	fc := New(10, time.Second*1)

	n := 1
	go func() {
		for {
			fc.Ctrl()
			t.Logf("N: %d", n)
			n++
		}
	}()

	m := 1
	go func() {
		for {
			fc.Ctrl()
			t.Logf("M: %d", m)
			m++
		}
	}()

	<-c
}
