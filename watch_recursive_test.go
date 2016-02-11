package runit

import "testing"

func TestShouldIgnoreFile(t *testing.T) {
	ignore := []string{".git.*", ".*.pyc", "a/b/c$", "d/e.xx", "*badregex"}
	w, _ := NewRecursiveWatcher("./", ignore)

	var ignoretests = []struct {
		in  string
		out bool
	}{
		{"a", false},
		{"a/b/c", true},
		{"a/b/c.py", false},
		{"a/b/c.pyc", true},
		{"e.xx", false},
		{".git", true},
		{".git/d", true},
	}
	for _, tt := range ignoretests {
		s := w.ShouldIgnoreFile(tt.in)
		if s != tt.out {
			t.Errorf("ShouldIgnoreFile(%s) => %v, want %v", tt.in, s, tt.out)
		}
	}
}

func TestShouldIgnoreFileDefaults(t *testing.T) {
	return
	ignore := []string{}
	w, _ := NewRecursiveWatcher("./", ignore)

	var ignoretests = []struct {
		in  string
		out bool
	}{
		{"a", false},
		{"a/b/c", false},
		{"a/b/c.py", false},
		{"a/b/c.pyc", false},
		{"e.xx", false},
		{".git", true},
		{".git/d", true},
	}
	for _, tt := range ignoretests {
		s := w.ShouldIgnoreFile(tt.in)
		if s != tt.out {
			t.Errorf("ShouldIgnoreFile(%s) => %v, want %v", tt.in, s, tt.out)
		}
	}
}
