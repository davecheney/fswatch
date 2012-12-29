// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin

package fswatch

import (
	"errors"
	"os"
	"syscall"
)

func (k *kqueue) add(fd uintptr, filter, flags int, fflags uint32) error {
	var buf [1]syscall.Kevent_t
	ev := &buf[0]
	flags = flags | syscall.EV_ADD | syscall.EV_RECEIPT
	evset(ev, fd, filter, flags, fflags)
	n, err := syscall.Kevent(k.fd, buf[:], buf[:], nil)
	if err != nil {
		return os.NewSyscallError("kevent", err)
	}
	if n != 1 || (ev.Flags&syscall.EV_ERROR) == 0 || uintptr(ev.Ident) != fd || int(ev.Filter) != filter {
		return errors.New("kqueue phase error")
	}
	if ev.Data != 0 {
		return syscall.Errno(ev.Data)
	}
	return nil
}
