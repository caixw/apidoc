// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package input

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/issue9/assert"
)

var (
	api1 = []byte(` @api POST /users/login 登录
 group users
 tags: [t1,t2]

 request:
   description: request body
   content:
     application/json:
       schema:
         type: object
         properties:
           username:
             type: string
             description: 登录账号
           password:
             type: string
  	         description: 密码`)

	api2 = []byte(` @api DELETE /users/login 注销登录
 group users
 tags: [t1,t2]

 request:
   description: request body
   content:
     application/json:
       schema:
         type: object
         properties:
           username:
             type: string
             description: 登录账号
           password:
             type: string
   	         description: 密码`)

	doc = []byte(` @apidoc title of api
 version: 2.9
 license:
   name: MIT
   url: https://opensources.org/licenses/MIT
 description:>
   line1
   line2`)
)

func TestParse(t *testing.T) {
	a := assert.New(t)

	testParse(a, "go")
	/*testParse(a, "groovy")
	testParse(a, "java")
	testParse(a, "javascript")
	testParse(a, "pascal")
	testParse(a, "perl")
	testParse(a, "php")
	testParse(a, "python")
	testParse(a, "ruby")
	testParse(a, "rust")
	testParse(a, "swift")*/
}

func testParse(a *assert.Assertion, lang string) {
	o := &Options{
		Lang:      lang,
		Dir:       "./testdata/" + lang,
		Recursive: true,
	}
	a.NotError(o.Sanitize()) // 初始化扩展名信息

	channel, err := Parse(o)
	a.NotError(err).NotNil(channel)

	for b := range channel {
		eq := bytes.Equal(b.Data, api1) ||
			bytes.Equal(b.Data, api2) ||
			bytes.Equal(b.Data, doc)
		a.True(eq, "返回的内容为：%s", string(b.Data))
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
