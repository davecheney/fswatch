// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin freebsd netbsd openbsd

package fswatch

import (
	"syscall"
	"testing"
)

func checkev(t *testing.T, expected, actual *syscall.Kevent_t) {
	if expected.Ident != actual.Ident {
		t.Errorf("ev.Ident: expected %d, received: %d", expected.Ident, actual.Ident)
	}
	if expected.Filter != actual.Filter {
		t.Errorf("ev.Filter: expected %d, received: %d", expected.Filter, actual.Filter)
	}
	if expected.Flags != actual.Flags {
		t.Errorf("ev.Flags: expected %d, received: %d", expected.Flags, actual.Flags)
	}
	if expected.Fflags != actual.Fflags {
		t.Errorf("ev.Fflags: expected %d, received: %d", expected.Fflags, actual.Fflags)
	}
}

func TestKqueuenewWatcher(t *testing.T) {
	w, err := newWatcher()
	assert(t, err)
	assert(t, w.close())
}

func TestKqueuewatcheradd(t *testing.T) {
	d, done := setup(t)
	defer done()

	w, err := newWatcher()
	assert(t, err)
	defer w.close()

	f := newfile(t, d)
	defer f.Close()

	w.add(f.Name())
}

func TestKqueuewait(t *testing.T) {
	d, done := setup(t)
	defer done()

	w, err := newWatcher()
	assert(t, err)
	defer w.close()

	f := newfile(t, d)
	defer f.Close()

	w.add(f.Name())

	n, err := f.Write([]byte("hey hey"))
	assert(t, err)
	var buf [1]syscall.Kevent_t
	n, err = w.kqueue.wait(buf[:])
	assert(t, err)
	if n != len(buf) {
		t.Fatalf("expecting 1 notification, received %d", n)
	}
}
