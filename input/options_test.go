// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package input

import (
	"path/filepath"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/internal/config"
)

var _ config.Sanitizer = &Options{}

func TestOptions_Sanitize(t *testing.T) {
	a := assert.New(t)

	o := &Options{Dir: "not exists"}
	a.Error(o.Sanitize())

	o.Dir = "./"
	o.Lang = "not exists"
	a.Error(o.Sanitize())

	// 未指定扩展名，则使用系统默认的
	o.Lang = "go"
	a.NotError(o.Sanitize())
	a.Equal(o.Exts, langExts["go"])

	// 指定了 Exts，自动调整扩展名样式。
	o.Lang = "go"
	o.Exts = []string{"go", ".g2"}
	a.NotError(o.Sanitize())
	a.Equal(o.Exts, []string{".go", ".g2"})
}

func TestDetectExts(t *testing.T) {
	a := assert.New(t)

	files, err := detectExts("./testdir", false)
	a.NotError(err)
	a.Equal(len(files), 4)
	a.Equal(files[".php"], 1).Equal(files[".c"], 1)

	files, err = detectExts("./testdir", true)
	a.NotError(err)
	a.Equal(len(files), 5)
	a.Equal(files[".php"], 1).Equal(files[".1"], 3)
}

func TestDetect(t *testing.T) {
	a := assert.New(t)

	o, err := Detect("./testdir", true)
	a.NotError(err).NotEmpty(o)
	a.NotContains(o.Exts, ".1") // .1 不存在于已定义的语言中
}

func TestRecursivePath(t *testing.T) {
	a := assert.New(t)

	opt := &Options{Dir: "./testdir", Recursive: false, Exts: []string{".c", ".h"}}
	paths, err := recursivePath(opt)
	a.NotError(err)
	a.Contains(paths, []string{
		filepath.Join("testdir", "testfile.c"),
		filepath.Join("testdir", "testfile.h"),
	})

	opt.Dir = "./testdir"
	opt.Recursive = true
	opt.Exts = []string{".1", ".2"}
	paths, err = recursivePath(opt)
	a.NotError(err)
	a.Contains(paths, []string{
		filepath.Join("testdir", "testdir1", "testfile.1"),
		filepath.Join("testdir", "testdir1", "testfile.2"),
		filepath.Join("testdir", "testdir2", "testfile.1"),
		filepath.Join("testdir", "testfile.1"),
	})

	opt.Dir = "./testdir/testdir1"
	opt.Recursive = true
	opt.Exts = []string{".1", ".2"}
	paths, err = recursivePath(opt)
	a.NotError(err)
	a.Contains(paths, []string{
		filepath.Join("testdir", "testdir1", "testfile.1"),
		filepath.Join("testdir", "testdir1", "testfile.2"),
	})

	opt.Dir = "./testdir"
	opt.Recursive = true
	opt.Exts = []string{".1"}
	paths, err = recursivePath(opt)
	a.NotError(err)
	a.Contains(paths, []string{
		filepath.Join("testdir", "testdir1", "testfile.1"),
		filepath.Join("testdir", "testdir2", "testfile.1"),
		filepath.Join("testdir", "testfile.1"),
	})
}
