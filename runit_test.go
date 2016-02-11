package runit

// TODO better synchronization to exit rather than time.Sleep

import (
	"syscall"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	_, err := New("ls", "test", []string{}, false, false)
	t.Log(err)
	if err != nil {
		t.Error(err)
	}
}

func TestNewNoCommand(t *testing.T) {
	_, err := New("", "test", []string{}, false, false)
	if err == nil {
		t.Error("cmd empty should be err")
	}
}

func TestNewWatchInvalidPath(t *testing.T) {
	_, err := New("true", "nothere", []string{}, false, false)
	if err == nil {
		t.Fatal("should get no folders to watch here")
	}
}

func TestDoRun(t *testing.T) {
	runner, err := New("true", "", []string{}, false, false)
	if err != nil {
		t.Error(err)
	}
	status, err := runner.Do()
	if err != nil {
		t.Error(err)
	}
	if status != 0 {
		t.Error("status not 0 got", status)
	}
}

func TestDoRepeat(t *testing.T) {
	r, err := New("true", ".", []string{}, true, true)
	if err != nil {
		t.Error(err)
	}
	go func() {
		time.Sleep(500 * time.Millisecond)
		r.Shutdown()
		r.Interrupt <- syscall.SIGINT
	}()
	status, err := r.Do()
	if err != nil {
		t.Error(err)
	}
	if status != 0 {
		t.Error("status not 0 got", status)
	}
}

func TestDoStart(t *testing.T) {
	r, err := New("true", "", []string{}, true, false)
	if err != nil {
		t.Error(err)
	}
	go func() {
		time.Sleep(500 * time.Millisecond)
		r.Shutdown()
		r.Interrupt <- syscall.SIGINT
	}()
	status, err := r.Do()
	if err != nil {
		t.Error(err)
	}
	if status != 0 {
		t.Error("status not 0 got", status)
	}
}

func TestRun(t *testing.T) {
	runner, err := New("true", "", []string{}, false, false)
	if err != nil {
		t.Fatal(err)
	}
	status, err := runner.Run()
	if err != nil {
		t.Fatal(err)
	}
	if status != 0 {
		t.Error("status not 0 got", status)
	}
}

func TestKill(t *testing.T) {
	r, err := New("sleep 1", "", []string{}, true, false)
	if err != nil {
		t.Error(err)
	}
	go func() {
		time.Sleep(500 * time.Millisecond)
		r.Shutdown()
		r.Interrupt <- syscall.SIGINT
	}()
	status, err := r.Do()
	if err != nil {
		t.Error(err)
	}
	if status != 0 {
		t.Error("status not 0 got", status)
	}
	r.Kill()
}

func TestShutdown(t *testing.T) {
	r, err := New("sleep 1", "", []string{}, true, false)
	if err != nil {
		t.Error(err)
	}
	go func() {
		time.Sleep(500 * time.Millisecond)
		r.Interrupt <- syscall.SIGINT
	}()
	status, err := r.Do()
	if err != nil {
		t.Error(err)
	}
	if status != 0 {
		t.Error("status not 0 got", status)
	}
	r.Shutdown()
}
