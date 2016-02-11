package runit

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-fsnotify/fsnotify"
)

// RecursiveWatcher https://github.com/nathany/looper/blob/master/watch.go#L13
type RecursiveWatcher struct {
	*fsnotify.Watcher
	Ignore []string
}

// NewRecursiveWatcher https://github.com/nathany/looper/blob/master/watch.go#L19
func NewRecursiveWatcher(path string, ignore []string) (*RecursiveWatcher, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	rw := &RecursiveWatcher{Watcher: watcher, Ignore: ignore}

	log.Println("watching folders in", path)

	rw.AddFolder(path)
	return rw, nil
}

// AddFolder https://github.com/nathany/looper/blob/master/watch.go#L40
func (w *RecursiveWatcher) AddFolder(path string) {
	w.Add(path)

	subFolders := w.Subfolders(path)
	if len(subFolders) == 0 {
		return
	}
	for _, folder := range subFolders {
		err := w.Add(folder)
		if err != nil {
			log.Printf("error watching: %s %v\n", folder, err)
			return
		}
		log.Printf("added folder %s", folder)
	}
}

// Subfolders returns a slice of subfolders (recursive), including the folder provided.
func (w *RecursiveWatcher) Subfolders(path string) (paths []string) {
	filepath.Walk(path, func(newPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			name := info.Name()
			if w.ShouldIgnoreFile(name) && name != "." && name != ".." {
				return filepath.SkipDir
			}
			paths = append(paths, newPath)
		}
		return nil
	})
	return paths
}

// ShouldIgnoreFile tests whether or not to ignore the event name triggered.
func (w *RecursiveWatcher) ShouldIgnoreFile(name string) bool {
	if len(w.Ignore) > 0 {
		for _, ignoreInput := range w.Ignore {
			ignore, err := regexp.Compile(ignoreInput)
			if err != nil {
				log.Println("[ERR]", err, ignoreInput)
				continue
			}
			if ignore.MatchString(name) {
				log.Printf("ignoring %s", name)
				return true
			}
		}
		return false
	}
	// by default skip folders and files that begin with a dot
	return strings.HasPrefix(name, ".")
}
