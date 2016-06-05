// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package input

import (
	"testing"

	"github.com/caixw/apidoc/doc"
	"github.com/issue9/assert"
)

func TestParseFile(t *testing.T) {
	a := assert.New(t)

	testParseFile(a, "go", "./testdata/go/test1.go")
	testParseFile(a, "php", "./testdata/php/test1.php")
	testParseFile(a, "c", "./testdata/c/test1.c")
	testParseFile(a, "ruby", "./testdata/ruby/test1.rb")
	testParseFile(a, "java", "./testdata/java/test1.java")
	testParseFile(a, "javascript", "./testdata/javascript/test1.js")
	testParseFile(a, "python", "./testdata/python/test1.py")
}

func testParseFile(a *assert.Assertion, lang string, path string) {
	docs := doc.New()
	a.NotNil(docs)

	b, found := langs[lang]
	if !found {
		a.T().Error("不支持该语言")
	}

	parseFile(docs, path, b)
	a.Equal(2, len(docs.Apis))

	api0 := docs.Apis[0]
	api1 := docs.Apis[1]

	a.Equal(api0.URL, "/users/login").
		Equal(api1.URL, "/users/login").
		Equal(api0.Group, "users").
		Equal(api1.Group, "users")

	if api0.Method == "POST" {
		a.Equal(api1.Method, "DELETE").
			Equal(1, len(api1.Request.Headers))
	} else {
		a.Equal(api0.Method, "DELETE").
			Equal(api1.Method, "POST").
			Equal(1, len(api0.Request.Headers))
	}
}

func TestRecursivePath(t *testing.T) {
	a := assert.New(t)

	opt := &Options{Dir: "./testdir", Recursive: false, Exts: []string{".c", ".h"}}
	paths, err := recursivePath(opt)
	a.NotError(err)
	a.Contains(paths, []string{
		"testdir/testfile.c",
		"testdir/testfile.h",
	})

	opt.Dir = "./testdir"
	opt.Recursive = true
	opt.Exts = []string{".1", ".2"}
	paths, err = recursivePath(opt)
	a.NotError(err)
	a.Contains(paths, []string{
		"testdir/testdir1/testfile.1",
		"testdir/testdir1/testfile.2",
		"testdir/testdir2/testfile.1",
		"testdir/testfile.1",
	})

	opt.Dir = "./testdir/testdir1"
	opt.Recursive = true
	opt.Exts = []string{".1", ".2"}
	paths, err = recursivePath(opt)
	a.NotError(err)
	a.Contains(paths, []string{
		"testdir/testdir1/testfile.1",
		"testdir/testdir1/testfile.2",
	})

	opt.Dir = "./testdir"
	opt.Recursive = true
	opt.Exts = []string{".1"}
	paths, err = recursivePath(opt)
	a.NotError(err)
	a.Contains(paths, []string{
		"testdir/testdir1/testfile.1",
		"testdir/testdir2/testfile.1",
		"testdir/testfile.1",
	})
}

func TestOptions_Init(t *testing.T) {
	a := assert.New(t)

	o := &Options{Dir: "not exists"}
	a.Error(o.Init())

	o.Dir = "./"
	o.Lang = "not exists"
	a.Error(o.Init())

	// 未指定扩展名，则使用系统默认的
	o.Lang = "c"
	a.NotError(o.Init())
	a.Equal(o.Exts, langExts["c"])

	// 指定了 Exts，自动调整扩展名样式。
	o.Lang = "c"
	o.Exts = []string{"c1", ".c2"}
	a.NotError(o.Init())
	a.Equal(o.Exts, []string{".c1", ".c2"})
}
