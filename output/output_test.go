// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

import (
	"os"
	"testing"

	"github.com/caixw/apidoc/core"
	"github.com/caixw/apidoc/scanner"
	"github.com/issue9/assert"
)

func TestLexer_Html(t *testing.T) {
	testdir := "./testdir"
	a := assert.New(t)

	// 创建测试目录
	a.NotError(os.MkdirAll(testdir, os.ModePerm), "无法创建测试目录")

	docs, err := core.ScanFiles([]string{"./test.php"}, scanner.CStyle)
	a.NotError(err).
		NotNil(docs).
		True(docs.HasError()). // 第一个注释块会返回一个语法错误
		True(len(docs.Items()) > 0)

	opt := &Options{
		DocDir:     testdir,
		Version:    "doc v0.1",
		AppVersion: "appver 0.1",
		Title:      "TestDoc",
		Elapsed:    1,
	}
	a.NotError(Html(docs.Items(), opt))
}
