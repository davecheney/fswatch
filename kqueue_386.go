// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin,386 freebsd,386 netbsd,386 openbsd,386

package fswatch

import (
	"syscall"
)

func evset(ev *syscall.Kevent_t, fd uintptr, filter, flags int, fflags uint32) {
	ev.Ident = uint32(fd)
	ev.Filter = int16(filter)
	ev.Flags = uint16(flags)
	ev.Fflags = fflags
}
