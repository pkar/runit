package runit

import (
	"os"
	"path/filepath"

	"github.com/go-fsnotify/fsnotify"
)

// Watch watches the runner watch path for changes and
// notifies the runner of change events.
func (r *Runner) Watch(shutdown <-chan struct{}) (chan bool, error) {
	restartChan := make(chan bool, 2)

	watcher, err := NewRecursiveWatcher(r.WatchPath)
	if err != nil {
		perror(err)
		return restartChan, err
	}

	go func(restart chan<- bool) {
		for {
			select {
			case <-shutdown:
				watcher.Close()
			case event := <-watcher.Events:
				pinfo("event:", event)
				switch {
				case event.Op&fsnotify.Create == fsnotify.Create:
					// create a file or directory
					fi, err := os.Stat(event.Name)
					if err != nil {
						// eg. stat .subl513.tmp : no such file or directory
						perror(err)
						continue
					}

					if fi.IsDir() {
						pdebugf("detected new directory %s", event.Name)
						if !shouldIgnoreFile(filepath.Base(event.Name)) {
							watcher.AddFolder(event.Name)
							restart <- true
							pinfof("added new folder: %s", event.Name)
						}
						continue
					}
					// created a file
					watcher.Files <- event.Name
					restart <- true
					pdebugf("added new file: %s", event.Name)
				case event.Op&fsnotify.Write == fsnotify.Write:
					watcher.Files <- event.Name
					restart <- true
					pdebugf("added modified file: %s", event.Name)
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					restart <- true
					watcher.Remove(event.Name)
					pdebugf("removed file: %s", event.Name)
				case event.Op&fsnotify.Chmod == fsnotify.Chmod:
					restart <- true
					pdebugf("chmod file: %s", event.Name)
				case event.Op&fsnotify.Rename == fsnotify.Rename:
					// renaming a file triggers a create event
					pdebugf("renamed file: %s", event.Name)
				}
			case err := <-watcher.Errors:
				perror(err)
			}
		}
	}(restartChan)

	return restartChan, nil
}
