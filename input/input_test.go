// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package input

import (
	"path/filepath"
	"testing"

	"github.com/caixw/apidoc/input/encoding"
	"github.com/caixw/apidoc/types"
	"github.com/issue9/assert"
)

var _ types.Sanitizer = &Options{}

func TestParse(t *testing.T) {
	a := assert.New(t)

	testParse(a, "go")
	testParse(a, "groovy")
	testParse(a, "java")
	testParse(a, "javascript")
	testParse(a, "pascal")
	testParse(a, "perl")
	testParse(a, "php")
	testParse(a, "python")
	testParse(a, "ruby")
	testParse(a, "rust")
	testParse(a, "swift")
}

func testParse(a *assert.Assertion, lang string) {
	o := &Options{
		Lang:      lang,
		Dir:       "./testdata/" + lang,
		Recursive: true,
	}
	a.NotError(o.Sanitize()) // 初始化扩展名信息

	docs := types.NewDoc()
	err := Parse(docs, o)
	a.NotError(err).
		Equal(len(docs.Apis), 2)

		// doc.xx
	a.Equal(docs.Title, "title of api").
		Equal(docs.Version, "2.9").
		Equal(docs.BaseURL, "https://api.caixw.io").
		Equal(docs.LicenseName, "MIT").
		Equal(docs.LicenseURL, "https://opensources.org/licenses/MIT").
		Equal(docs.Content, "\n line1\n line2\n")

	// test1.xx
	api0 := docs.Apis[0]
	api1 := docs.Apis[1]
	a.Equal(api0.URL, "/users/login").
		Equal(api1.URL, "/users/login").
		Equal(api0.Group, "users").
		Equal(api1.Group, "users")
}

func TestParseFile(t *testing.T) {
	a := assert.New(t)

	testParseFile(a, "go", "./testdata/go/test1.go")
	testParseFile(a, "groovy", "./testdata/groovy/test1.groovy")
	testParseFile(a, "java", "./testdata/java/test1.java")
	testParseFile(a, "javascript", "./testdata/javascript/test1.js")
	testParseFile(a, "pascal", "./testdata/pascal/test1.pas")
	testParseFile(a, "perl", "./testdata/perl/test1.pl")
	testParseFile(a, "php", "./testdata/php/test1.php")
	testParseFile(a, "python", "./testdata/python/test1.py")
	testParseFile(a, "ruby", "./testdata/ruby/test1.rb")
	testParseFile(a, "rust", "./testdata/rust/test1.rs")
	testParseFile(a, "swift", "./testdata/swift/test1.swift")
}

func testParseFile(a *assert.Assertion, lang, path string) {
	docs := types.NewDoc()
	a.NotNil(docs)

	b, found := langs[lang]
	if !found {
		a.TB().Error("不支持该语言")
	}

	parseFile(docs, path, b, &Options{StartLineNumber: 1, Encoding: encoding.DefaultEncoding})
	a.Equal(2, len(docs.Apis))

	api0 := docs.Apis[0]
	api1 := docs.Apis[1]

	a.TB().Log(api0.Success.Examples[0])
	a.TB().Log(api1.Success.Examples[0])

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
