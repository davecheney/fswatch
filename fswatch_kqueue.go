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
}

func newWatcher() (*watcher, error) {
	kq, err := newkqueue()
	return &watcher{
		kqueue: kq,
	}, err
}

func (w *watcher) add(path string) error {
	fd, err := syscall.Open(path, syscall.O_NONBLOCK|syscall.O_RDONLY, 0700)
	if err != nil {
		return os.NewSyscallError("open", err)
	}
	return w.kqueue.add(uintptr(fd), syscall.EVFILT_VNODE, 0, syscall.NOTE_DELETE|syscall.NOTE_EXTEND|syscall.NOTE_WRITE)
}
