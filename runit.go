package runit

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

// Runner ...
type Runner struct {
	cmdIn        string
	cmd          *exec.Cmd
	WatchPath    string
	Alive        bool
	Interrupt    chan os.Signal
	eventChan    chan bool
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
		cmdIn:        cmdIn,
		WatchPath:    watchPath,
		Alive:        alive,
		mu:           &sync.Mutex{},
		shutdownChan: make(chan struct{}),
		Interrupt:    make(chan os.Signal, 1),
	}
	if watchPath != "" {
		var err error
		runner.eventChan, err = runner.Watch(runner.shutdownChan)
		if err != nil {
			return nil, err
		}
	}

	return runner, nil
}

// Do determines the type of process to run based on Alive and WatchPath.
// It returns the exit status and any error.
func (r *Runner) Do() (int, error) {
	// just run the command and return.
	if !r.Alive && r.WatchPath == "" {
		return r.Run()
	}

	// begin the command
	err := r.Start()
	if err != nil {
		return 1, err
	}
	status := WaitFunc(r.Kill, r.Kill, r.Interrupt)
	return status, nil
}

// Run runs a subprocess and waits for it status
// to return.
func (r *Runner) Run() (int, error) {
	err := r.startCmd()
	if err != nil {
		return 1, err
	}
	exitStatus := 0
	if err = r.cmd.Wait(); err != nil {
		exitStatus = GetExitStatus(err)
	}
	return exitStatus, err
}

// Start begins the subprocess. If Alive and WatchPath
// are set it will restart on file changes.
// This should be used with Runner.WaitFunc or something
// like it.
func (r *Runner) Start() error {
	if r.WatchPath != "" || r.Alive {
		go r.RestartListen()
	}
	go func() {
		for {
			err := r.startCmd()
			if err != nil {
				perror(err)
				time.Sleep(time.Second)
			}

			select {
			case <-r.shutdownChan:
				return
			default:
				// run the command and restart after finished.
				err := r.cmd.Wait()
				if err != nil {
					time.Sleep(time.Second)
				}
			}
		}
	}()
	return nil
}

// RestartListen waits for restart events.
func (r *Runner) RestartListen() {
	for {
		select {
		case <-r.eventChan:
			pinfof("restart event")
			err := r.Kill()
			if err != nil {
				perror(err)
			}
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

	pdebugf("running %s", r.cmdIn)

	r.cmd = exec.Command("bash", "-c", "-e", r.cmdIn)
	r.cmd.Stdin = os.Stdin
	r.cmd.Stdout = os.Stdout
	r.cmd.Stderr = os.Stderr

	err := r.cmd.Start()
	if err != nil {
		perror(err)
	}
	return err
}

// Kill stops the runners subprocess
func (r *Runner) Kill() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.cmd == nil || r.cmd.Process == nil {
		return nil
	}
	pinfof("killing subprocess")
	return r.cmd.Process.Kill()
}

// Shutdown signals closing of the application.
func (r *Runner) Shutdown() {
	pinfof("shutting down")
	close(r.shutdownChan)
	r.Kill()
}
