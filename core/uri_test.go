// SPDX-License-Identifier: MIT

package core

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"unicode/utf8"

	"github.com/issue9/assert"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func TestFileURI(t *testing.T) {
	a := assert.New(t)

	path := "/path/file"
	uri, err := FileURI(path)
	a.NotError(err).Equal(uri, fileScheme+"://"+path)
	file, err := uri.File()
	a.NotError(err).Equal(path, file)

	uri = URI(path)
	file, err = uri.File()
	a.NotError(err).Equal(path, file).Equal(path, uri.String())
}

func TestURI_File(t *testing.T) {
	a := assert.New(t)

	uri := URI("file:///path.php")
	file, err := uri.File()
	a.NotError(err).Equal(file, "/path.php")

	uri = URI(" :/php")
	file, err = uri.File()
	a.Error(err).Empty(file)

	uri = URI("https://example.com/path.php")
	file, err = uri.File()
	a.Error(err).Empty(file)
}

func TestURI_isNoScheme(t *testing.T) {
	a := assert.New(t)

	a.True(URI("").isNoScheme())
	a.True(URI("a").isNoScheme())
	a.True(URI("./file.php").isNoScheme())
	a.True(URI("/file.php").isNoScheme())
	a.False(URI("file:///file.php").isNoScheme())

	a.True(URI("c:\\file.php").isNoScheme())
	a.True(URI("c:\\").isNoScheme())
	a.False(URI("c:/file.php").isNoScheme())
}

func TestURI_ReadAll(t *testing.T) {
	a := assert.New(t)

	uri, err := FileURI("./not-exists")
	a.NotError(err).NotEmpty(uri)
	data, err := uri.ReadAll(nil)
	a.Error(err).Nil(data)

	data, err = uri.ReadAll(simplifiedchinese.GB18030)
	a.Error(err).Nil(data)

	uri, err = FileURI("./uri.go")
	a.NotError(err).NotEmpty(uri)
	data, err = uri.ReadAll(encoding.Nop)
	a.NotError(err).NotNil(data)

	uri, err = FileURI("./gbk.php")
	a.NotError(err).NotEmpty(uri)
	data, err = uri.ReadAll(simplifiedchinese.GB18030)
	a.NotError(err).NotNil(data)
	a.Contains(string(data), "这是一个 GBK 编码的文件").
		Contains(string(data), "中文").
		True(utf8.Valid(data))

	// 未知协议
	uri = URI("unknown:///gbk.php")
	data, err = uri.ReadAll(simplifiedchinese.GB18030)
	a.Error(err).Nil(data)

	// 无效的 URI
	uri = URI(" :///gbk.php")
	data, err = uri.ReadAll(simplifiedchinese.GB18030)
	a.Error(err).Nil(data)

	// 测试远程读取
	static := http.FileServer(http.Dir("./"))
	srv := httptest.NewServer(static)
	defer srv.Close()

	uri = URI(srv.URL + "/uri.go")
	data, err = uri.ReadAll(nil)
	a.NotError(err).NotNil(data)

	uri = URI(srv.URL + "/gbk.php")
	data, err = uri.ReadAll(simplifiedchinese.GB18030)
	a.NotError(err).NotNil(data)
	a.Contains(string(data), "这是一个 GBK 编码的文件").
		Contains(string(data), "中文").
		True(utf8.Valid(data))

	// 不存在的远程文件
	uri = URI(srv.URL + "/not-exists")
	data, err = uri.ReadAll(nil)
	a.Error(err).Nil(data)
}
