package runit

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/pkar/log"
)

// WaitFunc listens for signals and restarts the given restart runner
// on SIGHUP or exits for where appropriate.
func WaitFunc(do func() error, die func() error, interrupt chan os.Signal) int {
	signal.Notify(interrupt)
	for {
		select {
		case sig := <-interrupt:
			switch sig {
			case syscall.SIGHUP:
				err := do()
				log.Info.Printf("captured %v restarting...err: %s", sig, err)
				continue
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL:
				die()
				log.Info.Printf("captured %v exiting", sig)
				return 0
			default:
				log.Info.Printf("captured %v continue...", sig)
				continue
			}
		}
	}
}

// GetExitStatus determines the exit status code of an err
// from a command that was run.
func GetExitStatus(waitError error) int {
	if exitError, ok := waitError.(*exec.ExitError); ok {
		if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
			return waitStatus.ExitStatus()
		}
	}
	return 1
}
