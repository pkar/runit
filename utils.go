package runit

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// Wait listens for signals and restarts the given restart runner
// on SIGHUP or exits for where appropriate.
func Wait(restart func() error, interrupt chan os.Signal) int {
	signal.Notify(interrupt)
	for {
		select {
		case sig := <-interrupt:
			switch sig {
			case syscall.SIGHUP:
				pdebugf("captured %v restarting...", sig)
				restart()
				continue
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL:
				pdebugf("captured %v exiting", sig)
				return 0
			default:
				pdebugf("captured %v continue...", sig)
				continue
			}
		}
	}
}

// GetExitStatus determines the exit status code of an err
// from a command that was run.
func GetExitStatus(waitError error) (int, error) {
	exitError, ok := waitError.(*exec.ExitError)
	if ok {
		waitStatus, ok := exitError.Sys().(syscall.WaitStatus)
		if ok {
			return waitStatus.ExitStatus(), nil
		}
	}
	return 1, waitError
}
