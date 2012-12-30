// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux

package fswatch

import (
	"os"
	"strings"
	"sync"
	"syscall"
	"unsafe"
)

type watcher struct {
	*inotify
	C     chan *Event
	err   chan error
	buf   [syscall.SizeofInotifyEvent * 4096]byte
	paths struct {
		sync.RWMutex
		m map[int]string
	}
}

func newWatcher() (*watcher, error) {
	i, err := newinotify()
	w := &watcher{
		inotify: i,
		C:       make(chan *Event, 1),
		err:     make(chan error, 1),
	}
	w.paths.m = make(map[int]string)
	return w, err
}

func (w *watcher) run() {
	defer w.close()
	for {
		n, err := syscall.Read(w.fd, w.buf[:])
		if n == 0 || err != nil {
			w.err <- err
			return
		}
		// We don't know how many events we just read into the buffer
		// While the offset points to at least one whole event...
		w.paths.RLock()
		for offset := uint32(0); offset <= uint32(n-syscall.SizeofInotifyEvent); {
			// Point "raw" to the event in the buffer
			raw := (*syscall.InotifyEvent)(unsafe.Pointer(&w.buf[offset]))
			nameLen := uint32(raw.Len)
			target := w.paths.m[int(raw.Wd)]
			if nameLen > 0 {
				// Point "bytes" at the first byte of the filename
				bytes := (*[syscall.PathMax]byte)(unsafe.Pointer(&w.buf[offset+syscall.SizeofInotifyEvent]))
				// The filename is padded with NUL bytes. TrimRight() gets rid of those.
				target += "/" + strings.TrimRight(string(bytes[0:nameLen]), "\000")
			}
			w.C <- &Event{
				Target: target,
				event: event{
					mask: raw.Mask,
				},
			}

			// Move to the next event in the buffer
			offset += syscall.SizeofInotifyEvent + nameLen
		}
		w.paths.RUnlock()
	}

}

func (w *watcher) add(path string) error {
	w.paths.RLock()
	for _, v := range w.paths.m {
		if v == path {
			w.paths.RUnlock()
			return ErrWatchExists
		}
	}
	w.paths.RUnlock()

	const flags = syscall.IN_ALL_EVENTS
	wd, err := w.inotify.add(path, flags)
	if err != nil {
		return &os.PathError{
			Op:   "inotify_add_watch",
			Path: path,
			Err:  err,
		}
	}
	w.paths.Lock()
	w.paths.m[wd] = path
	w.paths.Unlock()
	return nil
}

func (w *watcher) close() error {
	return syscall.Close(w.fd)
}

type event struct {
	mask uint32
}

func (e *event) IsCreate() bool { return e.mask&syscall.IN_CREATE > 0 }
func (e *event) IsRemove() bool { return e.mask&syscall.IN_DELETE > 0 }
