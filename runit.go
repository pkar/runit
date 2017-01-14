package runit

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

// Runner ...
type Runner struct {
	Alive        bool
	Interrupt    chan os.Signal
	Wait         bool
	WatchPath    string
	Ignore       []string
	cmd          *exec.Cmd
	cmdIn        string
	eventChan    chan bool
	shutdownChan chan struct{}
	mu           *sync.Mutex
}

// New initializes a command runner and watches for changes
// in path if watch is given.
func New(cmdIn string, watchPath string, ignore []string, alive bool, wait bool) (*Runner, error) {
	if cmdIn == "" {
		return nil, fmt.Errorf("no command defined")
	}
	runner := &Runner{
		Alive:        alive,
		Wait:         wait,
		Ignore:       ignore,
		Interrupt:    make(chan os.Signal, 1),
		WatchPath:    watchPath,
		cmdIn:        cmdIn,
		mu:           &sync.Mutex{},
		shutdownChan: make(chan struct{}),
	}
	if watchPath != "" {
		var err error
		runner.eventChan, err = runner.Watch(runner.shutdownChan, ignore)
		if err != nil {
			return nil, err
		}
	}

	return runner, nil
}

// Do determines the type of process to run based on Alive and WatchPath/Wait.
// It returns the exit status and any error.
func (r *Runner) Do() (status int, err error) {
	switch {
	// just run the command and return.
	case !r.Alive && r.WatchPath == "":
		return r.Run()
	// wait for file changes before running the command.
	case r.WatchPath != "" && r.Wait:
		go r.Repeat()
		status = WaitFunc(r.Kill, r.Kill, r.Interrupt)
	// begin the command and restart on finish or if with watch
	// on file changes as well.
	default:
		err := r.Start()
		if err != nil {
			return 1, err
		}
		if r.WatchPath != "" {
			go r.RestartListen()
		}
		status = WaitFunc(r.Kill, r.Kill, r.Interrupt)
	}
	return status, nil
}

// Repeat runs the command on event signals.
func (r *Runner) Repeat() {
	for {
		select {
		case <-r.shutdownChan:
			return
		case <-r.eventChan:
			status, err := r.Run()
			if status != 0 || err != nil {
				log.Println("ERRO:", status, err)
			}
		}
	}
}

// Run runs a subprocess and waits for its status
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
	go func() {
		for {
			err := r.startCmd()
			if err != nil {
				log.Println("ERRO:", err)
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
			log.Println("INFO: restart event")
			err := r.Kill()
			if err != nil {
				log.Println("ERRO:", err)
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

	log.Printf("INFO: running %s\n", r.cmdIn)

	r.cmd = exec.Command("bash", "-c", "-e", r.cmdIn)
	r.cmd.Stdin = os.Stdin
	r.cmd.Stdout = os.Stdout
	r.cmd.Stderr = os.Stderr

	err := r.cmd.Start()
	if err != nil {
		log.Println("ERRO:", err)
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
	log.Println("INFO: killing subprocess")
	return r.cmd.Process.Kill()
}

// Shutdown signals closing of the application.
func (r *Runner) Shutdown() {
	log.Println("INFO: shutting down")
	close(r.shutdownChan)
	r.Kill()
}
