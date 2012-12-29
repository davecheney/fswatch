// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin freebsd netbsd openbsd

package fswatch

import (
	"testing"
)

func TestKqueuenewWatcher(t *testing.T) {
	w, err := newWatcher()
	assert(t, err)
	assert(t, w.close())
}

func TestKqueuewatcheradd(t *testing.T) {
	d, done := setup(t)
	defer done()

	w, err := newWatcher()
	assert(t, err)
	defer w.close()

	f := newfile(t, d)
	defer f.Close()

	w.add(f.Name())
}
