package runit

// TODO better synchronization to exit rather than time.Sleep

import (
	"io/ioutil"
	"os"
	"syscall"
	"testing"
	"time"
)

func init() {
	LogLevel = 0
}

func TestNew(t *testing.T) {
	_, err := New("ls", "test", false, false)
	t.Log(err)
	if err != nil {
		t.Error(err)
	}
}

func TestNewNoCommand(t *testing.T) {
	_, err := New("", "test", false, false)
	if err == nil {
		t.Error("cmd empty should be err")
	}
}

func TestNewWatchInvalidPath(t *testing.T) {
	_, err := New("true", "nothere", false, false)
	if err == nil {
		t.Fatal("should get no folders to watch here")
	}
}

func TestDoRun(t *testing.T) {
	runner, err := New("true", "", false, false)
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
	r, err := New("true", ".", true, true)
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
	r, err := New("true", "", true, false)
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
	runner, err := New("true", "", false, false)
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
	r, err := New("sleep 1", "", true, false)
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
	r, err := New("sleep 1", "", true, false)
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

func TestSighup(t *testing.T) {
	r, err := New("sleep 1", "", true, false)
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

func TestWatch(t *testing.T) {
	r, err := New("true", "test", true, false)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll("test/test")
	}()

	go func() {
		time.Sleep(3 * time.Second)
		r.Interrupt <- syscall.SIGINT
	}()

	_, err = r.Do()
	if err != nil {
		t.Error(err)
	}

	// create file
	err = ioutil.WriteFile("test/test.txt", []byte("hello"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("testing creating test file")

	// write file
	err = ioutil.WriteFile("test/test.txt", []byte("goodbye"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("testing writing to test file")

	// rename file
	err = os.Rename("test/test.txt", "test/test1.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("testing renaming test file")

	// remove
	err = os.Remove("test/test1.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("testing removing test file")

	t.Log("cleanup test dir")
	err = os.RemoveAll("test/test")
	if err != nil {
		t.Fatal(err)
	}

	// create dir
	err = os.MkdirAll("test/test", 0777)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("testing creating test dir")
}
