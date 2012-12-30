// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin freebsd netbsd openbsd

package fswatch

import (
	"os"
	"sync"
	"syscall"
)

type watcher struct {
	*kqueue
	C     chan *Event
	waker *os.File
	sync.RWMutex
	fds   map[string]int
	paths map[int]string
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
		fds:    make(map[string]int),
		paths:  make(map[int]string),
	}, nil
}

func (w *watcher) add(path string) error {
	fd, err := syscall.Open(path, syscall.O_NONBLOCK|syscall.O_RDONLY, 0700)
	if err != nil {
		return os.NewSyscallError("open", err)
	}
	w.Lock()
	defer w.Unlock()
	if _, exists := w.fds[path]; exists {
		defer syscall.Close(fd)
		return ErrWatchExists
	}
	w.fds[path] = fd
	w.paths[fd] = path
	return w.kqueue.add(uintptr(fd), syscall.EVFILT_VNODE, 0, syscall.NOTE_DELETE|syscall.NOTE_EXTEND|syscall.NOTE_WRITE)
}

func (w *watcher) run() {
	defer w.kqueue.close()
	var buf [18]syscall.Kevent_t
	for done := false; !done; {
		n, err := w.kqueue.wait(buf[:])
		if err != nil {
			println(err.Error())
			return
		}
		for _, e := range buf[:n] {
			if uintptr(e.Ident) == w.waker.Fd() {
				done = true
				continue
			}
			w.C <- &Event{
				Target: w.paths[int(e.Ident)],
				event:  event(e),
			}
		}
	}
	println("done")
}

func (w *watcher) close() error {
	return w.waker.Close()
}

type event syscall.Kevent_t

func (e *event) IsCreate() bool { return e.Flags&syscall.NOTE_EXTEND > 0 }
func (e *event) IsRemove() bool { return e.Flags&syscall.NOTE_DELETE > 0 }
