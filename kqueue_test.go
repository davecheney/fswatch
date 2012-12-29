// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin freebsd netbsd openbsd

package fswatch

import (
	"os"
	"syscall"
	"testing"
)

func TestNewKqueue(t *testing.T) {
	kq, err := newkqueue()
	assert(t, err)
	assert(t, kq.close())
}

func TestKqueueAddSocket(t *testing.T) {
	kq, err := newkqueue()
	assert(t, err)
	s, err := syscall.Socket(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	assert(t, err)
	defer syscall.Close(s)
	assert(t, kq.add(uintptr(s), syscall.EVFILT_WRITE, 0, 0))
	assert(t, kq.close())
}

func TestKqueueAddFile(t *testing.T) {
	d, done := setup(t)
	defer done()
	kq, err := newkqueue()
	assert(t, err)
	f := newfile(t, d)
	defer f.Close()
	assert(t, kq.add(f.Fd(), syscall.EVFILT_WRITE, 0, 0))
	assert(t, kq.close())
}

func TestKqueueWaitPipeClose(t *testing.T) {
	kq, err := newkqueue()
	assert(t, err)
	pr, pw, err := os.Pipe()
	assert(t, err)
	defer pw.Close()
	assert(t, kq.add(uintptr(pr.Fd()), syscall.EVFILT_READ, syscall.EV_ONESHOT, 0))
	result := make(chan error)
	go func() {
		_, err := kq.wait(make([]syscall.Kevent_t, 1))
		result <- err
	}()
	assert(t, pw.Close())
	assert(t, <-result)
	assert(t, kq.close())
}
