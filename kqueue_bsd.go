// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build freebsd netbsd openbsd

package fswatch

import (
	"os"
	"syscall"
)

func (k *kqueue) add(fd uintptr, filter, flags int, fflags uint32) error {
	var buf [1]syscall.Kevent_t
	ev := &buf[0]
	flags = flags | syscall.EV_ADD
	evset(ev, fd, filter, flags, fflags)
	_, err := syscall.Kevent(k.fd, buf[:], nil, nil)
	if err != nil {
		return os.NewSyscallError("kevent", err)
	}
	return nil
}
