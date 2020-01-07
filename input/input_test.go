// SPDX-License-Identifier: MIT

package input

import (
	"testing"

	"github.com/issue9/assert"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"

	"github.com/caixw/apidoc/v6/message/messagetest"
)

func TestParse(t *testing.T) {
	a := assert.New(t)

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

	doc, err := Parse(h, php, c)
	a.NotError(err).NotNil(doc).
		Equal(1, len(doc.Apis)).
		Equal(doc.Version, "1.1.1")
	api := doc.Apis[0]
	a.Equal(api.Method, "GET")
	h.Stop()
	a.Empty(erro.String())
}

func TestReadFile(t *testing.T) {
	a := assert.New(t)

	nop, err := readFile("./testdata/gbk.php", encoding.Nop)
	a.NotError(err).
		NotNil(nop).
		NotContains(string(nop), "这是一个 GBK 编码的文件")

	def, err := readFile("./testdata/gbk.php", nil)
	a.NotError(err).
		NotNil(def).
		NotContains(string(def), "这是一个 GBK 编码的文件")
	a.Equal(def, nop)

	data, err := readFile("./testdata/gbk.php", simplifiedchinese.GB18030)
	a.NotError(err).
		NotNil(data).
		Contains(string(data), "这是一个 GBK 编码的文件")
}
