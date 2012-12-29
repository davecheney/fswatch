// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fswatch

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// watcher is platform dependent. This check ensures that all
// watcher implementations conform to the specification.
var _ interface {
	add(string) error
} = &watcher{}

// event is platform dependant. Events are passed as pointers,
// not interfaces to allow inlining to operate effectively. This
// check ensures all event implementations confirm to the
// specificatoin
var _ interface {
	IsCreate() bool
	IsRemove() bool
} = &event{}

// setup creates a temporary working directory, and a cleanup function.
func setup(t *testing.T) (string, func()) {
	d, err := ioutil.TempDir("", "fswatch-")
	if err != nil {
		t.Fatal(err)
	}
	return d, func() {
		os.RemoveAll(d)
	}
}

func assert(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func expect(t *testing.T, expected, actual error) {
	if expected != actual {
		t.Fatalf("expected: %v, actual: %v", expected, actual)
	}
}

func mkdir(t *testing.T, dir string) {
	assert(t, os.Mkdir(dir, 0777))
}

func newdir(t *testing.T, dir string) string {
	p, err := ioutil.TempDir(dir, "dir")
	assert(t, err)
	return p
}

func newfile(t *testing.T, dir string) *os.File {
	f, err := ioutil.TempFile(dir, "file")
	assert(t, err)
	return f
}

func TestNewWatcher(t *testing.T) {
	w, err := NewWatcher()
	assert(t, err)
	assert(t, w.Close())
}

func TestWatcherAddPath(t *testing.T) {
	d, done := setup(t)
	defer done()
	w, err := NewWatcher()
	assert(t, err)
	assert(t, w.Add(d))
	assert(t, w.Close())
}

func TestWatch(t *testing.T) {
	d, done := setup(t)
	defer done()
	w, err := Watch(d)
	assert(t, err)
	assert(t, w.Close())
}

func TestWatchMissingPath(t *testing.T) {
	d, done := setup(t)
	defer done()
	w, err := Watch(filepath.Join(d, "missing"))
	if err == nil {
		w.Close()
		t.Fatalf("expected: error, got: %v", err)
	}
}

func TestWatchCreateFile(t *testing.T) {
	d, done := setup(t)
	defer done()
	w, err := Watch(d)
	assert(t, err)
	defer w.Close()
	f := newfile(t, d)
	defer f.Close()
	select {
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("timeout")
	case e := <-w.C:
		if f.Name() != e.Target || !e.IsCreate() {
			t.Fatalf("expected create, got %v", e)
		}
	}
}

func TestWatchDeleteFile(t *testing.T) {
	d, done := setup(t)
	defer done()
	f := newfile(t, d)
	defer f.Close()
	w, err := Watch(d)
	assert(t, err)
	defer w.Close()
	if err := os.Remove(f.Name()); err != nil {
		t.Fatal(err)
	}
	select {
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("timeout")
	case e := <-w.C:
		if f.Name() != e.Target || !e.IsRemove() {
			t.Fatalf("expected remove, got %v", e)
		}
	}
}

func TestWatchCreateDir(t *testing.T) {
	d, done := setup(t)
	defer done()
	w, err := Watch(d)
	assert(t, err)
	defer w.Close()
	p := newdir(t, d)
	select {
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("timeout")
	case e := <-w.C:
		if p != e.Target || !e.IsCreate() {
			t.Fatalf("expected create, got %v", e)
		}
	}
}

func TestWatchDeleteDir(t *testing.T) {
	d, done := setup(t)
	defer done()
	p := newdir(t, d)
	w, err := Watch(d)
	assert(t, err)
	defer w.Close()
	if err := os.Remove(p); err != nil {
		t.Fatal(err)
	}
	select {
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("timeout")
	case e := <-w.C:
		if p != e.Target || !e.IsRemove() {
			t.Fatalf("expected remove, got %v", e)
		}
	}
}

func TestWatchSelfFile(t *testing.T) {
	d, done := setup(t)
	defer done()
	w, err := Watch(d)
	assert(t, err)
	defer w.Close()
	f := newfile(t, d)
	assert(t, err)
	defer f.Close()
	select {
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("timeout")
	case e := <-w.C:
		if f.Name() != e.Target || !e.IsCreate() {
			t.Fatalf("expected create, got %v", e)
		}
	}
}

func TestCannotWatchSameFileTwice(t *testing.T) {
	d, done := setup(t)
	defer done()
	w, err := NewWatcher()
	assert(t, err)
	defer w.Close()
	f := newfile(t, d)
	defer f.Close()
	assert(t, w.Add(f.Name()))
	expect(t, ErrWatchExists, w.Add(f.Name()))
}

func TestCannotWatchSameDirTwice(t *testing.T) {
	d, done := setup(t)
	defer done()
	w, err := NewWatcher()
	assert(t, err)
	defer w.Close()
	d1 := newdir(t, d)
	assert(t, w.Add(d1))
	expect(t, ErrWatchExists, w.Add(d1))
}

// disabled, hangs on linux
func testCanWatchFileInsideWatchedDir(t *testing.T) {
	d, done := setup(t)
	defer done()
	w, err := NewWatcher()
	assert(t, err)
	defer w.Close()
	d1 := newdir(t, d)
	assert(t, w.Add(d1))
	f := newfile(t, d1)
	defer f.Close()
	assert(t, w.Add(f.Name()))
}

func TestWatchTwoPaths(t *testing.T) {
	d, done := setup(t)
	defer done()
	w, err := NewWatcher()
	assert(t, err)
	defer w.Close()
	d1 := newdir(t, d)
	d2 := newdir(t, d)
	assert(t, w.Add(d1))
	assert(t, w.Add(d2))

	f := newfile(t, d1)
	defer f.Close()
	select {
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("timeout")
	case e := <-w.C:
		if f.Name() != e.Target || !e.IsCreate() {
			t.Errorf("d1: expected %q, got %q", f.Name(), e.Target)
		}
	}
	f = newfile(t, d2)
	defer f.Close()
	select {
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("timeout")
	case e := <-w.C:
		if f.Name() != e.Target || !e.IsCreate() {
			t.Errorf("d2: expected %q, got %q", f.Name(), e.Target)
		}
	}
}
