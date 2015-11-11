package runit

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-fsnotify/fsnotify"
)

// RecursiveWatcher https://github.com/nathany/looper/blob/master/watch.go#L13
type RecursiveWatcher struct {
	*fsnotify.Watcher
	Files   chan string
	Folders chan string
}

// NewRecursiveWatcher https://github.com/nathany/looper/blob/master/watch.go#L19
func NewRecursiveWatcher(path string) (*RecursiveWatcher, error) {
	folders := Subfolders(path)
	if len(folders) == 0 {
		return nil, fmt.Errorf("No folders to watch.")
	}
	pinfo("watching folders in", path)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	rw := &RecursiveWatcher{Watcher: watcher}

	rw.Files = make(chan string, 10000)
	rw.Folders = make(chan string, 10000)

	for _, folder := range folders {
		rw.AddFolder(folder)
	}
	return rw, nil
}

// AddFolder https://github.com/nathany/looper/blob/master/watch.go#L40
func (watcher *RecursiveWatcher) AddFolder(folder string) {
	err := watcher.Add(folder)
	if err != nil {
		perrorf("error watching: %s %v\n", folder, err)
		return
	}
	watcher.Folders <- folder
	pdebugf("added folder %s", folder)
}

// Subfolders returns a slice of subfolders (recursive), including the folder provided.
func Subfolders(path string) (paths []string) {
	filepath.Walk(path, func(newPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			name := info.Name()
			// skip folders that begin with a dot
			if shouldIgnoreFile(name) && name != "." && name != ".." {
				return filepath.SkipDir
			}
			paths = append(paths, newPath)
		}
		return nil
	})
	return paths
}

// shouldIgnoreFile determines if a file should be ignored.
// File names that begin with "." or "_" are ignored by the go tool.
func shouldIgnoreFile(name string) bool {
	return strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_")
}
