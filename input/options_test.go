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

func TestRecursiveDir(t *testing.T) {
	a := assert.New(t)

	files, err := recursiveDir("./testdir", false)
	a.NotError(err)
	a.Equal(len(files), 4)

	files, err = recursiveDir("./testdir", true)
	a.NotError(err)
	a.Equal(len(files), 7)
}

/*func TestDetectDirLang(t *testing.T) {
	a := assert.New(t)

	lang, err := DetectDirLang("./testdir")
	a.NotError(err).Equal(lang, "c++")

	lang, err = DetectDirLang("./testdir/testdir1")
	a.Error(err).Empty(lang)
}
*/
