// SPDX-License-Identifier: MIT

package input

import (
	"path/filepath"
	"testing"

	"github.com/issue9/assert"
	"golang.org/x/text/encoding/simplifiedchinese"

	"github.com/caixw/apidoc/v5/internal/lang"
	opt "github.com/caixw/apidoc/v5/options"
)

func TestBuildOptions(t *testing.T) {
	a := assert.New(t)

	o := &opt.Input{}
	oo, err := buildOptions(o)
	a.Error(err).Nil(oo)

	o.Dir = "not exists"
	oo, err = buildOptions(o)
	a.Error(err).Nil(oo)

	o.Dir = "./"
	oo, err = buildOptions(o)
	a.Error(err).Nil(oo)

	o.Lang = "not exists"
	oo, err = buildOptions(o)
	a.Error(err).Nil(oo)

	// 未指定扩展名，则使用系统默认的
	language := lang.Get("go")
	o.Lang = "go"
	oo, err = buildOptions(o)
	a.NotError(err).NotNil(oo)
	a.Equal(oo.Exts, language.Exts)

	// 指定了 Exts，自动调整扩展名样式。
	o.Lang = "go"
	o.Exts = []string{"go", ".g2"}
	oo, err = buildOptions(o)
	a.NotError(err).NotNil(oo)
	a.Equal(oo.Exts, []string{".go", ".g2"})

	// 特定的编码
	o.Encoding = "GBK"
	oo, err = buildOptions(o)
	a.NotError(err).NotNil(oo)
	a.Equal(oo.encoding, simplifiedchinese.GBK)

	// 不存在的编码
	o.Encoding = "not-exists---"
	oo, err = buildOptions(o)
	a.Error(err).Nil(oo)
}

func TestRecursivePath(t *testing.T) {
	a := assert.New(t)

	opt := &opt.Input{
		Dir:       "./testdir",
		Recursive: false,
		Exts:      []string{".c", ".h"},
	}
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
