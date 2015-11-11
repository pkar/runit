package runit

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

const (
	// MAXERRORS is a hardcoded maximum number of attempts
	// to run a process before quitting.
	MAXERRORS = 10
)

// Runner ...
type Runner struct {
	cmdIn        string
	Cmd          *exec.Cmd
	WatchPath    string
	MaxErrors    int
	Alive        bool
	restartChan  chan bool
	shutdownChan chan struct{}
	mu           *sync.Mutex
}

// New initializes a command runner and watches for changes
// in path if watch is given.
func New(cmdIn string, watchPath string, alive bool) (*Runner, error) {
	if cmdIn == "" {
		return nil, fmt.Errorf("no command defined")
	}
	runner := &Runner{
		MaxErrors:    MAXERRORS,
		cmdIn:        cmdIn,
		WatchPath:    watchPath,
		Alive:        alive,
		mu:           &sync.Mutex{},
		shutdownChan: make(chan struct{}),
	}
	if watchPath != "" {
		var err error
		runner.restartChan, err = runner.Watch(runner.shutdownChan)
		if err != nil {
			return nil, err
		}
	}

	return runner, nil
}

// Run runs the subprocess with optional keep alive
// if it fails. If Alive is not set it will just
// finish the command and return. Optionally if
// Alive and WatchPath are set it will restart on
// file changes.
func (r *Runner) Run() error {
	if !r.Alive {
		err := r.startCmd()
		if err != nil {
			return err
		}
		return r.Cmd.Wait()
	}

	if r.WatchPath != "" {
		go r.RestartListen()
	}
	nErrs := 0
	go func(nerrs int) {
		for {
			if nerrs > r.MaxErrors {
				return
			}
			err := r.startCmd()
			if err != nil {
				perror(err)
				nErrs++
				time.Sleep(time.Second)
			}

			select {
			case <-r.shutdownChan:
				return
			default:
				err := r.Cmd.Wait()
				if err != nil {
					nErrs++
				}
			}
		}
	}(nErrs)
	if nErrs >= r.MaxErrors {
		return fmt.Errorf("maximum error retries attempted")
	}
	return nil
}

// RestartListen waits for restart events.
func (r *Runner) RestartListen() {
	for {
		select {
		case <-r.restartChan:
			pinfof("restart event")
			err := r.Restart()
			if err != nil {
				perror(err)
				return
			}
			//r.Cmd.Wait()
		case <-r.shutdownChan:
			return
		}
	}
}

// startCmd starts the command and doesn't wait
// for it to complete.
func (r *Runner) startCmd() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	pinfof("running %s", r.cmdIn)

	r.Cmd = exec.Command("bash", "-c", "-e", r.cmdIn)
	r.Cmd.Stdin = os.Stdin
	r.Cmd.Stdout = os.Stdout
	r.Cmd.Stderr = os.Stderr

	err := r.Cmd.Start()
	if err != nil {
		perror(err)
	}
	return err
}

// Kill stops the runners subprocess
func (r *Runner) Kill() error {
	if r.Cmd == nil || r.Cmd.Process == nil {
		return nil
	}
	pinfof("killing subprocess")
	err := r.Cmd.Process.Kill()
	if err != nil {
		perror(err)
	}
	return err
}

// Shutdown signals closing of the application.
func (r *Runner) Shutdown() {
	pinfof("shutting down")
	r.Kill()
	close(r.shutdownChan)
}

// Restart kills the runners subprocess and starts up a
// new one
func (r *Runner) Restart() error {
	pdebug("restarting")
	err := r.Kill()
	return err
}
