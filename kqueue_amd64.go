// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin,amd64 freebsd,amd64 netbsd,amd64 openbsd,amd64

package fswatch

import (
	"syscall"
)

func evset(ev *syscall.Kevent_t, fd uintptr, filter, flags int, fflags uint32) {
	ev.Ident = uint64(fd)
	ev.Filter = int16(filter)
	ev.Flags = uint16(flags)
	ev.Fflags = fflags
}
