// SPDX-License-Identifier: MIT

package build

import (
	"testing"

	"github.com/issue9/assert"
	"golang.org/x/text/encoding/simplifiedchinese"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
	"github.com/caixw/apidoc/v7/internal/lang"
)

func TestParseInputs(t *testing.T) {
	a := assert.New(t)

	blocks := make(chan core.Block, 100)
	erro, _, h := messagetest.MessageHandler()
	php := &Input{
		Lang:      "php",
		Dir:       "./testdata",
		Recursive: true,
		Encoding:  "gbk",
	}
	a.NotError(php.Sanitize())

	c := &Input{
		Lang:      "c++",
		Dir:       "./testdata",
		Recursive: true,
	}
	a.NotError(c.Sanitize())

	parseInputs(blocks, h, php, c)
	close(blocks)

	a.Equal(6, len(blocks))
	a.Empty(erro.String())
}

func TestInput_ParseFile(t *testing.T) {
	a := assert.New(t)

	blocks := make(chan core.Block, 100)
	erro, _, h := messagetest.MessageHandler()
	o := &Input{
		Lang:      "c++",
		Dir:       "./testdata",
		Recursive: true,
	}
	a.NotError(o.Sanitize())
	o.ParseFile(blocks, h, "./testdata/testfile.c")
	h.Stop()
	close(blocks)

	a.Equal(3, len(blocks))
	a.Empty(erro.String())

	// 非 utf8 编码
	blocks = make(chan core.Block, 100)
	erro, _, h = messagetest.MessageHandler()
	o = &Input{
		Lang:      "php",
		Dir:       "./testdata",
		Recursive: true,
		Encoding:  "gbk",
	}
	a.NotError(o.Sanitize())
	o.ParseFile(blocks, h, "./testdata/gbk.php")
	h.Stop()
	close(blocks)
	a.Equal(1, len(blocks))
	blk := <-blocks
	a.Equal(string(blk.Data), `   <api method="GET">
   <path path="/apis/gbk" />
   <description type="markdown"><![CDATA[1223 中文 45 ]]></description>
   <server>test</server>
   </api>
`).
		Equal(blk.Location.Range, core.Range{
			Start: core.Position{Line: 5, Character: 0},
			End:   core.Position{Line: 10, Character: 0},
		})
	a.Empty(erro.String())

	// 文件不存在
	blocks = make(chan core.Block, 100)
	erro, _, h = messagetest.MessageHandler()
	o = &Input{
		Lang: "c++",
		Dir:  "./testdata",
	}
	a.NotError(o.Sanitize())
	o.ParseFile(blocks, h, "./testdata/not-exists.php")
	close(blocks)
	h.Stop()
	a.NotEmpty(erro.String())

	// 没有正确的结束符号
	blocks = make(chan core.Block, 100)
	erro, _, h = messagetest.MessageHandler()
	o = &Input{
		Lang: "c++",
		Dir:  "./testdata",
	}
	a.NotError(o.Sanitize())
	o.ParseFile(blocks, h, "./testdata/testfile.1")
	h.Stop()
	close(blocks)
	a.NotEmpty(erro.String())
}

func TestOptions_Sanitize(t *testing.T) {
	a := assert.New(t)

	var o *Input
	a.Error(o.Sanitize())

	o = &Input{}
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
	o.Encoding = "GbK"
	a.NotError(o.Sanitize())
	a.Equal(o.encoding, simplifiedchinese.GBK)

	// 不存在的编码
	o.Encoding = "not-exists---"
	a.Error(o.Sanitize())
}

func TestRecursivePath(t *testing.T) {
	a := assert.New(t)

	opt := &Input{
		Dir:       "./testdata",
		Recursive: false,
		Exts:      []string{".c", ".h"},
	}
	paths, err := recursivePath(opt)
	a.NotError(err).Equal(2, len(paths))

	opt.Dir = "./testdata"
	opt.Recursive = true
	opt.Exts = []string{".1", ".2"}
	paths, err = recursivePath(opt)
	a.NotError(err).Equal(4, len(paths))

	opt.Dir = "./testdata/testdir1"
	opt.Recursive = true
	opt.Exts = []string{".1", ".2"}
	paths, err = recursivePath(opt)
	a.NotError(err).Equal(2, len(paths))

	opt.Dir = "./testdata"
	opt.Recursive = true
	opt.Exts = []string{".1"}
	paths, err = recursivePath(opt)
	a.NotError(err).Equal(3, len(paths))
}
