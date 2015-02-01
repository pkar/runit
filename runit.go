package runit

import (
	"os"
	"os/exec"
	"strings"

	"github.com/pkar/runit/vendor/log"
)

type Runner struct {
	cmdIn       string
	cmd         *exec.Cmd
	watchPath   string
	restartChan chan bool
}

// Run initializes a command runner and watches for changes
// in path if watch is given.
func Run(cmdIn string, watchPath string) (*Runner, error) {
	runner := &Runner{
		cmdIn:     cmdIn,
		watchPath: watchPath,
	}
	if watchPath != "" {
		var err error
		runner.restartChan, err = runner.Watch()
		if err != nil {
			log.Error(err)
		}
	}

	go runner.runCmd()
	go func() {
		for {
			select {
			case <-runner.restartChan:
				err := runner.Restart()
				if err != nil {
					log.Error(err)
				}
			}
		}
	}()

	return runner, nil
}

// runCmd ...
func (r *Runner) runCmd() {
	tokens := strings.Split(r.cmdIn, " ")
	r.cmd = exec.Command(tokens[0], tokens[1:]...)
	r.cmd.Stdin = os.Stdin
	r.cmd.Stdout = os.Stdout
	err := r.cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = r.cmd.Wait()
	if err != nil {
		log.Error(err)
	}
}

// Restart kills the runners subprocess and starts up a
// new one
func (r *Runner) Restart() error {
	log.Info("restart event ", r.cmdIn)
	err := r.cmd.Process.Kill()
	if err != nil {
		log.Error(err)
	}
	r.runCmd()

	return nil
}
