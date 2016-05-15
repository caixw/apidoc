// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"encoding/json"
	"io/ioutil"
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

	cfg := &config{Input: &input{Dir: "./testdir", Recursive: false, Exts: []string{".1", ".2"}}}
	paths, err := recursivePath(cfg)
	a.NotError(err)
	a.Equal(paths, []string{
		"testdir/testfile.1",
		"testdir/testfile.2",
	})

	cfg.Input.Dir = "./testdir"
	cfg.Input.Recursive = true
	cfg.Input.Exts = []string{".1", ".2"}
	paths, err = recursivePath(cfg)
	a.NotError(err)
	a.Contains(paths, []string{
		"testdir/testdir1/testfile.1",
		"testdir/testdir1/testfile.2",
		"testdir/testdir2/testfile.1",
		"testdir/testfile.1",
		"testdir/testfile.2",
	})

	cfg.Input.Dir = "./testdir/testdir1"
	cfg.Input.Recursive = true
	cfg.Input.Exts = []string{".1", ".2"}
	paths, err = recursivePath(cfg)
	a.NotError(err)
	a.Equal(paths, []string{
		"testdir/testdir1/testfile.1",
		"testdir/testdir1/testfile.2",
	})

	cfg.Input.Dir = "./testdir"
	cfg.Input.Recursive = true
	cfg.Input.Exts = []string{".1"}
	paths, err = recursivePath(cfg)
	a.NotError(err)
	a.Equal(paths, []string{
		"testdir/testdir1/testfile.1",
		"testdir/testdir2/testfile.1",
		"testdir/testfile.1",
	})
}

func TestGenConfigFile(t *testing.T) {
	a := assert.New(t)
	a.NotError(genConfigFile())

	data, err := ioutil.ReadFile("./" + configFilename)
	a.NotError(err).NotNil(data)
	cfg := &config{}
	a.NotError(json.Unmarshal(data, cfg))

	a.Equal(cfg.Input.Dir, "./").Equal(cfg.Input.Recursive, true)
}
