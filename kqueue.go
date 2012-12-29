// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin freebsd netbsd openbsd

package fswatch

import (
	"os"
	"syscall"
)

type kqueue struct {
	fd int
}

func newkqueue() (*kqueue, error) {
	fd, err := syscall.Kqueue()
	if err != nil {
		return nil, os.NewSyscallError("kqueue", err)
	}
	return &kqueue{
		fd: fd,
	}, nil
}

func (k *kqueue) wait(buf []syscall.Kevent_t) (int, error) {
	n, err := syscall.Kevent(k.fd, nil, buf, nil)
	if err != nil {
		return 0, os.NewSyscallError("kevent", err)
	}
	return n, nil
}

func (k *kqueue) close() error {
	return syscall.Close(k.fd)
}
