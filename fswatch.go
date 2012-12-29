// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fswatch

import (
	"errors"
	"sync"
)

var (
	ErrWatchExists = errors.New("existing watch for path")
)

type Watcher struct {
	// C receives the stream of Events from watched paths.
	C   <-chan *Event
	err struct {
		sync.Mutex
		val error
	}
	*watcher
}

type Event struct {
	Source, Target string
	event
}

// NewWatcher creates a new Watcher.
func NewWatcher() (*Watcher, error) {
	w, err := newWatcher()
	if err != nil {
		return nil, err
	}
	if w == nil {
		panic(w)
	}
	go w.run()
	return &Watcher{
		C:       w.C,
		watcher: w,
	}, nil
}

// Add adds path to the list of paths monitored by this Watcher.
// If path is a directory, the watcher will observe all changes to
// the directory and its direct decentants.
func (w *Watcher) Add(path string) error {
	return w.setError(w.add(path))
}

// Watch creates a new watcher for the supplied path.
func Watch(path string) (*Watcher, error) {
	w, err := NewWatcher()
	if err != nil {
		return nil, err
	}
	if err := w.Add(path); err != nil {
		w.Close()
		return nil, err
	}
	return w, nil
}

func (w *Watcher) Close() error {
	return w.close()
}

// setError sets the internal error value, if not already set.
func (w *Watcher) setError(err error) error {
	w.err.Lock()
	if w.err.val != nil {
		w.err.val = err
	}
	w.err.Unlock()
	return err
}
