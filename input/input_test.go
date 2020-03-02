// SPDX-License-Identifier: MIT

package input

import (
	"testing"
	"unicode/utf8"

	"github.com/issue9/assert"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"

	"github.com/caixw/apidoc/v6/message/messagetest"
)

func TestParse(t *testing.T) {
	a := assert.New(t)

	blocks := make(chan Block, 100)
	erro, _, h := messagetest.MessageHandler()
	php := &Options{
		Lang:      "php",
		Dir:       "./testdata",
		Recursive: true,
		Encoding:  "gbk",
	}

	c := &Options{
		Lang:      "c++",
		Dir:       "./testdata",
		Recursive: true,
	}
	a.NotError(Parse(blocks, h, php, c))
	close(blocks)

	a.Equal(5, len(blocks))
	a.Empty(erro.String())
}

func TestParseFile(t *testing.T) {
	a := assert.New(t)

	blocks := make(chan Block, 100)
	erro, _, h := messagetest.MessageHandler()
	o := &Options{
		Lang:      "c++",
		Dir:       "./testdata",
		Recursive: true,
	}
	a.NotError(o.sanitize())
	ParseFile(blocks, h, "./testdata/testfile.c", o)
	close(blocks)

	a.Equal(2, len(blocks))
	a.Empty(erro.String())

	// 非 utf8 编码
	blocks = make(chan Block, 100)
	erro, _, h = messagetest.MessageHandler()
	o = &Options{
		Lang:      "php",
		Dir:       "./testdata",
		Recursive: true,
		Encoding:  "gbk",
	}
	a.NotError(o.Sanitize())
	ParseFile(blocks, h, "./testdata/gbk.php", o)
	close(blocks)

	a.Equal(1, len(blocks))
	a.Empty(erro.String())
}

func TestReadFile(t *testing.T) {
	a := assert.New(t)

	nop, err := readFile("./testdata/gbk.php", encoding.Nop)
	a.NotError(err).
		NotNil(nop).
		NotContains(string(nop), "这是一个 GBK 编码的文件").
		False(utf8.Valid(nop))

	def, err := readFile("./testdata/gbk.php", nil)
	a.NotError(err).
		NotNil(def).
		NotContains(string(def), "这是一个 GBK 编码的文件").
		False(utf8.Valid(def))
	a.Equal(def, nop)

	data, err := readFile("./testdata/gbk.php", simplifiedchinese.GB18030)
	a.NotError(err).
		NotNil(data).
		Contains(string(data), "这是一个 GBK 编码的文件").
		Contains(string(data), "中文").
		True(utf8.Valid(data))
}
