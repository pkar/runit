package runit

import (
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	_, err := New("ls", "test")
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewNoCommand(t *testing.T) {
	_, err := New("", "abc")
	if err == nil {
		t.Fatal("cmd empty should be err")
	}
}

func TestRun(t *testing.T) {
	runner, err := New("echo", "")
	if err != nil {
		t.Fatal(err)
	}
	err = runner.Run(false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRunKeepAlive(t *testing.T) {
	runner, err := New("echo", "")
	if err != nil {
		t.Fatal(err)
	}
	err = runner.Run(true)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	runner.Shutdown()
}

func TestKill(t *testing.T) {
	runner, err := New("test/test.sh", "")
	if err != nil {
		t.Fatal(err)
	}
	err = runner.Run(true)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	err = runner.Kill()
	if err != nil {
		t.Fatal(err)
	}
	runner.Shutdown()
}

func TestRestart(t *testing.T) {
	runner, err := New("test/test.sh", "")
	if err != nil {
		t.Fatal(err)
	}
	err = runner.Run(false)
	if err != nil {
		t.Fatal(err)
	}
	err = runner.Restart()
	if err != nil {
		t.Fatal(err)
	}
	runner.Kill()
}

func TestRestartListen(t *testing.T) {
	runner, err := New("echo", "test")
	if err != nil {
		t.Fatal(err)
	}
	err = runner.Run(false)
	if err != nil {
		t.Fatal(err)
	}
	runner.restartChan <- true
	runner.Kill()
}

func TestWatch(t *testing.T) {
	runner, err := New("echo", "test")
	if err != nil {
		t.Fatal(err)
	}
	err = runner.Run(false)

	cmd := exec.Command("touch", "test/test")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()

	time.Sleep(1 * time.Second)

	cmd = exec.Command("rm", "test/test")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()

	runner.Shutdown()
}
