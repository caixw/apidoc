// SPDX-License-Identifier: MIT

package input

import (
	"path/filepath"
	"testing"

	"github.com/issue9/assert"
	"golang.org/x/text/encoding/simplifiedchinese"

	"github.com/caixw/apidoc/v5/internal/lang"
)

func TestOptions_Sanitize(t *testing.T) {
	a := assert.New(t)

	var o *Options
	a.Error(o.Sanitize())

	o = &Options{}
	a.Error(o.Sanitize())

	o.Dir = "not exists"
	a.Error(o.Sanitize())

	o.Dir = "./"
	a.Error(o.Sanitize())

	o.Lang = "not exists"
	a.Error(o.Sanitize())

	// 未指定扩展名，则使用系统默认的
	language := lang.Get("go")
	o.Lang = "go"
	a.NotError(o.Sanitize())
	a.Equal(o.Exts, language.Exts)

	// 指定了 Exts，自动调整扩展名样式。
	o.Lang = "go"
	o.Exts = []string{"go", ".g2"}
	a.NotError(o.Sanitize())
	a.Equal(o.Exts, []string{".go", ".g2"})

	// 特定的编码
	o.Encoding = "GBK"
	a.NotError(o.Sanitize())
	a.Equal(o.encoding, simplifiedchinese.GBK)

	// 不存在的编码
	o.Encoding = "not-exists---"
	a.Error(o.Sanitize())
}

func TestRecursivePath(t *testing.T) {
	a := assert.New(t)

	opt := &Options{
		Dir:       "./testdata",
		Recursive: false,
		Exts:      []string{".c", ".h"},
	}
	paths, err := recursivePath(opt)
	a.NotError(err)
	a.Contains(paths, []string{
		filepath.Join("testdata", "testfile.c"),
		filepath.Join("testdata", "testfile.h"),
	})

	opt.Dir = "./testdata"
	opt.Recursive = true
	opt.Exts = []string{".1", ".2"}
	paths, err = recursivePath(opt)
	a.NotError(err)
	a.Contains(paths, []string{
		filepath.Join("testdata", "testdir1", "testfile.1"),
		filepath.Join("testdata", "testdir1", "testfile.2"),
		filepath.Join("testdata", "testdir2", "testfile.1"),
		filepath.Join("testdata", "testfile.1"),
	})

	opt.Dir = "./testdata/testdir1"
	opt.Recursive = true
	opt.Exts = []string{".1", ".2"}
	paths, err = recursivePath(opt)
	a.NotError(err)
	a.Contains(paths, []string{
		filepath.Join("testdata", "testdir1", "testfile.1"),
		filepath.Join("testdata", "testdir1", "testfile.2"),
	})

	opt.Dir = "./testdata"
	opt.Recursive = true
	opt.Exts = []string{".1"}
	paths, err = recursivePath(opt)
	a.NotError(err)
	a.Contains(paths, []string{
		filepath.Join("testdata", "testdir1", "testfile.1"),
		filepath.Join("testdata", "testdir2", "testfile.1"),
		filepath.Join("testdata", "testfile.1"),
	})
}
