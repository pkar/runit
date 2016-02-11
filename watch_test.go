package runit

import (
	"io/ioutil"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestWatch(t *testing.T) {
	r, err := New("true", "test", []string{}, true, false)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll("test/test")
	}()

	go func() {
		time.Sleep(3 * time.Second)
		r.Interrupt <- syscall.SIGINT
	}()

	_, err = r.Do()
	if err != nil {
		t.Error(err)
	}

	// create file
	err = ioutil.WriteFile("test/test.txt", []byte("hello"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("testing creating test file")

	// write file
	err = ioutil.WriteFile("test/test.txt", []byte("goodbye"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("testing writing to test file")

	// rename file
	err = os.Rename("test/test.txt", "test/test1.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("testing renaming test file")

	// remove
	err = os.Remove("test/test1.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("testing removing test file")

	t.Log("cleanup test dir")
	err = os.RemoveAll("test/test")
	if err != nil {
		t.Fatal(err)
	}

	// create dir
	err = os.MkdirAll("test/test", 0777)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("testing creating test dir")
}
