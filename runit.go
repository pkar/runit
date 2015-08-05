package runit

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Runner ...
type Runner struct {
	cmdIn        string
	cmd          *exec.Cmd
	watchPath    string
	restartChan  chan bool
	shutdownChan chan bool
	mu           *sync.Mutex
}

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}

// New initializes a command runner and watches for changes
// in path if watch is given.
func New(cmdIn string, watchPath string) (*Runner, error) {
	if cmdIn == "" {
		return nil, fmt.Errorf("no command defined")
	}
	runner := &Runner{
		cmdIn:        cmdIn,
		watchPath:    watchPath,
		mu:           &sync.Mutex{},
		shutdownChan: make(chan bool),
	}
	if watchPath != "" {
		var err error
		runner.restartChan, err = runner.Watch()
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	return runner, nil
}

// Wait for the running command to complete.
func (r *Runner) Wait() error {
	return r.cmd.Wait()
}

// Run runs the subprocess with optional keep alive
// if it fails.
func (r *Runner) Run(alive bool) error {
	if r.watchPath != "" {
		go r.RestartListen()
	}
	if !alive {
		return r.runCmd()
	}
	go func() {
		for {
			select {
			case <-r.shutdownChan:
				if r.watchPath != "" {
					r.shutdownChan <- true
				}
				return
			default:
				log.Printf("running %s", r.cmdIn)
				err := r.runCmd()
				if err != nil {
					log.Println(err)
				}
				r.cmd.Wait()
				time.Sleep(time.Second)
			}
		}
	}()
	return nil
}

// RestartListen waits for restart events.
func (r *Runner) RestartListen() {
	for {
		select {
		case <-r.restartChan:
			log.Println("restart event")
			err := r.Restart()
			if err != nil {
				log.Println(err)
			}
			r.cmd.Wait()
		case <-r.shutdownChan:
			return
		}
	}
}

// runCmd starts the command and doesn't wait
// for it to complete.
func (r *Runner) runCmd() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	log.Printf("running %s", r.cmdIn)

	tokens := strings.Split(r.cmdIn, " ")
	r.cmd = exec.Command(tokens[0], tokens[1:]...)
	r.cmd.Stdin = os.Stdin
	r.cmd.Stdout = os.Stdout
	r.cmd.Stderr = os.Stderr

	err := r.cmd.Start()
	if err != nil {
		log.Println(err)
	}
	return err
}

// Kill stops the runners subprocess
func (r *Runner) Kill() error {
	if r.cmd.Process == nil {
		return nil
	}
	log.Println("killing subprocess")
	err := r.cmd.Process.Kill()
	if err != nil {
		log.Println(err)
	}
	return err
}

// Shutdown signals closing of the application.
func (r *Runner) Shutdown() {
	log.Println("shutting down")
	r.Kill()
	r.shutdownChan <- true
}

// Restart kills the runners subprocess and starts up a
// new one
func (r *Runner) Restart() error {
	log.Println("restarting")
	err := r.Kill()
	if err != nil {
		log.Println(err)
	}
	return r.runCmd()
}
