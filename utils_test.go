package runit

import (
	"syscall"
	"testing"
	"time"
)

func TestSighup(t *testing.T) {
	r, err := New("sleep 1", "", []string{}, true, false)
	if err != nil {
		t.Error(err)
	}
	go func() {
		time.Sleep(500 * time.Millisecond)
		r.Interrupt <- syscall.SIGHUP
		r.Interrupt <- syscall.SIGINT
	}()
	r.Do()
}
