// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux

package fswatch

import (
	"os"
	"syscall"
)

type inotify struct {
	fd int
}

func newinotify() (*inotify, error) {
	fd, err := syscall.InotifyInit()
	if fd == -1 {
		return nil, os.NewSyscallError("inotify_init", err)
	}
	return &inotify{
		fd: fd,
	}, nil
}

func (i *inotify) add(path string, flags uint32) (int, error) {
	return syscall.InotifyAddWatch(i.fd, path, flags)
}

func (i *inotify) close() error {
	return syscall.Close(i.fd)
}
