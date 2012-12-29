// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin freebsd netbsd openbsd

package fswatch

import (
	"os"
	"syscall"
)

type watcher struct {
	*kqueue
	C     chan *Event
	waker *os.File
}

func newWatcher() (*watcher, error) {
	kq, err := newkqueue()
	if err != nil {
		kq.close()
		return nil, err
	}
	pr, pw, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	if err := kq.add(uintptr(pr.Fd()), syscall.EVFILT_READ, 0, 0); err != nil {
		kq.close()
		pw.Close()
		return nil, err
	}
	return &watcher{
		kqueue: kq,
		C:      make(chan *Event, 1),
		waker:  pw,
	}, nil
}

func (w *watcher) add(path string) error {
	fd, err := syscall.Open(path, syscall.O_NONBLOCK|syscall.O_RDONLY, 0700)
	if err != nil {
		return os.NewSyscallError("open", err)
	}
	return w.kqueue.add(uintptr(fd), syscall.EVFILT_VNODE, 0, syscall.NOTE_DELETE|syscall.NOTE_EXTEND|syscall.NOTE_WRITE)
}

func (w *watcher) run() {
	// nothing
}

func (w *watcher) close() error {
	return w.waker.Close()
}

type event struct {
	mask uint32
}

func (e *event) IsCreate() bool { return e.mask&syscall.NOTE_EXTEND > 0 }
func (e *event) IsRemove() bool { return e.mask&syscall.NOTE_DELETE > 0 }
