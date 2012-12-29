// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fswatch

import (
	"io/ioutil"
	"os"
	"testing"
)

// setup creates a temporary working directory, and a cleanup function.
func setup(t *testing.T) (string, func()) {
	d, err := ioutil.TempDir("", "fswatch-")
	if err != nil {
		t.Fatal(err)
	}
	return d, func() {
		os.RemoveAll(d)
	}
}

func assert(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func expect(t *testing.T, expected, actual error) {
	if expected != actual {
		t.Fatalf("expected: %v, actual: %v", expected, actual)
	}
}

func mkdir(t *testing.T, dir string) {
	assert(t, os.Mkdir(dir, 0777))
}

func newdir(t *testing.T, dir string) string {
	p, err := ioutil.TempDir(dir, "dir")
	assert(t, err)
	return p
}

func newfile(t *testing.T, dir string) *os.File {
	f, err := ioutil.TempFile(dir, "file")
	assert(t, err)
	return f
}
