// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"github.com/issue9/assert"
)

func TestDetectLangType(t *testing.T) {
	a := assert.New(t)

	l, err := detectLangType([]string{".abc1", ".abc1", ".abc1"})
	a.Error(err).Equal(0, len(l))

	l, err = detectLangType([]string{".js", ".php", ".abc1"})
	a.NotError(err).Equal("js", l)
}

func TestDetectDirLangType(t *testing.T) {
	a := assert.New(t)

	l, err := detectDirLangType("./")
	a.NotError(err).Equal(l, "go")

	l, err = detectDirLangType("./testdir")
	a.Error(err).Equal(0, len(l))
}

func TestRecursivePath(t *testing.T) {
	a := assert.New(t)

	paths, err := recursivePath("./testdir", false, ".1", ".2")
	a.NotError(err)
	a.Equal(paths, []string{
		"testdir/testfile.1",
		"testdir/testfile.2",
	})

	paths, err = recursivePath("./testdir", true, ".1", ".2")
	a.NotError(err)
	a.Contains(paths, []string{
		"testdir/testdir1/testfile.1",
		"testdir/testdir1/testfile.2",
		"testdir/testdir2/testfile.1",
		"testdir/testfile.1",
		"testdir/testfile.2",
	})

	paths, err = recursivePath("./testdir/testdir1", true, ".1", ".2")
	a.NotError(err)
	a.Equal(paths, []string{
		"testdir/testdir1/testfile.1",
		"testdir/testdir1/testfile.2",
	})

	paths, err = recursivePath("./testdir", true, ".1")
	a.NotError(err)
	a.Equal(paths, []string{
		"testdir/testdir1/testfile.1",
		"testdir/testdir2/testfile.1",
		"testdir/testfile.1",
	})
}
