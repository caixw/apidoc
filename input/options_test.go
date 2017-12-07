// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package input

import (
	"testing"

	"github.com/issue9/assert"
)

func TestOptions_Sanitize(t *testing.T) {
	a := assert.New(t)

	o := &Options{Dir: "not exists"}
	a.Error(o.Sanitize())

	o.Dir = "./"
	o.Lang = "not exists"
	a.Error(o.Sanitize())

	// 未指定扩展名，则使用系统默认的
	o.Lang = "c++"
	a.NotError(o.Sanitize())
	a.Equal(o.Exts, langExts["c++"])

	// 指定了 Exts，自动调整扩展名样式。
	o.Lang = "c++"
	o.Exts = []string{"c1", ".c2"}
	a.NotError(o.Sanitize())
	a.Equal(o.Exts, []string{".c1", ".c2"})
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
