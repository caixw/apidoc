// SPDX-License-Identifier: MIT

package path

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"unicode/utf8"

	"github.com/issue9/assert"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func TestRel(t *testing.T) {
	a := assert.New(t)

	data := []*struct {
		path, wd, result string
	}{
		{
			path:   "./test/data",
			wd:     "./test",
			result: filepath.Clean("data"),
		},
		{
			path:   "./test/data",
			wd:     "./tex/data",
			result: filepath.Clean("../../test/data"),
		},

		{ // 无法计算，返回原值
			path:   "/test/data",
			wd:     "./tex/data",
			result: filepath.Clean("/test/data"),
		},
	}

	for index, item := range data {
		result := filepath.Clean(Rel(item.path, item.wd))
		a.Equal(result, item.result, "not equal @%d,v1=%s,v2=%s", index, result, item.result)
	}
}

func TestCurrPath(t *testing.T) {
	a := assert.New(t)

	dir, err := filepath.Abs("./")
	a.NotError(err).NotEmpty(dir)

	d, err := filepath.Abs(CurrPath("./"))
	a.NotError(err).NotEmpty(d)

	a.Equal(d, dir)
}

func TestReadFile(t *testing.T) {
	a := assert.New(t)

	data, err := ReadFile("./not-exists", nil)
	a.Error(err).Nil(data)

	data, err = ReadFile("./not-exists", simplifiedchinese.GB18030)
	a.Error(err).Nil(data)

	data, err = ReadFile("./path.go", encoding.Nop)
	a.NotError(err).NotNil(data)

	data, err = ReadFile("./gbk.php", simplifiedchinese.GB18030)
	a.NotError(err).NotNil(data)
	a.Contains(string(data), "这是一个 GBK 编码的文件").
		Contains(string(data), "中文").
		True(utf8.Valid(data))

	// 测试远程读取
	static := http.FileServer(http.Dir("./"))
	srv := httptest.NewServer(static)
	defer srv.Close()

	data, err = ReadFile(srv.URL+"/path.go", nil)
	a.NotError(err).NotNil(data)

	data, err = ReadFile(srv.URL+"/gbk.php", simplifiedchinese.GB18030)
	a.NotError(err).NotNil(data)
	a.Contains(string(data), "这是一个 GBK 编码的文件").
		Contains(string(data), "中文").
		True(utf8.Valid(data))

	// 不存在的远程文件
	data, err = ReadFile(srv.URL+"/not-exists", nil)
	a.Error(err).Nil(data)
}
