// SPDX-License-Identifier: MIT

package core

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"testing"
	"unicode/utf8"

	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/issue9/assert"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func TestFileURI(t *testing.T) {
	a := assert.New(t)

	path := "/path/file"
	uri := FileURI(path)
	a.Equal(uri, SchemeFile+"://"+path)
	file, err := uri.File()
	a.NotError(err).Equal(path, file)

	uri = FileURI("file:///path")
	a.Equal(uri, "file:///path")
	file, err = uri.File()
	a.NotError(err).Equal("/path", file)
}

func TestURI_json(t *testing.T) {
	a := assert.New(t)

	obj := &struct {
		URI URI `json:"uri"`
	}{}

	data := `{"uri":"file:%2F%2F%2Fpath%2F%E4%B8%AD%E6%96%87.txt"}` // file:///path/中文.txt

	a.NotError(json.Unmarshal([]byte(data), obj))
	a.Equal(obj.URI, "file:///path/中文.txt")

	// URI = ""

	obj.URI = ""
	data = `{"uri":""}`
	a.ErrorString(json.Unmarshal([]byte(data), obj), locale.Sprintf(locale.ErrInvalidURI, ""))
	a.Equal(obj.URI, "")
}

func TestURI_Parse(t *testing.T) {
	a := assert.New(t)

	scheme, p := URI("/path/file").Parse()
	a.Empty(scheme).Equal(p, "/path/file")

	scheme, p = URI("c:/path/file").Parse()
	a.Empty(scheme).Equal(p, "c:/path/file")

	scheme, p = URI("file://c:/path/file").Parse()
	a.Equal(scheme, SchemeFile).Equal(p, "c:/path/file")
}

func TestURI_File(t *testing.T) {
	a := assert.New(t)

	uri := URI("file:///path.php")
	file, err := uri.File()
	a.NotError(err).Equal(file, "/path.php")

	// 不管路径格式是否准确
	uri = URI(" :/php")
	file, err = uri.File()
	a.NotError(err).Equal(file, " :/php")

	uri = URI("https://example.com/path.php")
	file, err = uri.File()
	a.Error(err).Empty(file)
}

func TestURI_Append(t *testing.T) {
	a := assert.New(t)

	uri := URI("file://root")
	a.Equal(uri.Append(""), uri)
	a.Equal(uri.Append("path"), "file://root"+string(os.PathSeparator)+"path")
	a.Equal(uri.Append("/path"), "file://root/path")
	a.Equal(uri.Append("//path"), "file://root//path")

	if runtime.GOOS == "windows" {
		a.Equal(uri.Append("\\path"), "file://root\\path")
	}

	uri = URI("file://root/")
	a.Equal(uri.Append(""), uri)
	a.Equal(uri.Append("/path"), "file://root/path")
	a.Equal(uri.Append("//path"), "file://root//path")

	if runtime.GOOS == "windows" {
		a.Equal(uri.Append("\\path"), "file://root/path")
	}
}

func TestURI_Exists(t *testing.T) {
	a := assert.New(t)

	uri := FileURI("./not-exists")
	a.NotEmpty(uri)
	exists, err := uri.Exists()
	a.NotError(err).False(exists)

	uri = FileURI("./uri.go")
	a.NotEmpty(uri)
	exists, err = uri.Exists()
	a.NotError(err).True(exists)

	uri = URI("./uri.go")
	exists, err = uri.Exists()
	a.NotError(err).True(exists)

	// 未知的协议
	uri = URI("unknown:///uri.go")
	exists, err = uri.Exists()
	a.Error(err).False(exists)

	// 无效的 URI
	uri = URI(" :///uri.go")
	exists, err = uri.Exists()
	a.Error(err).False(exists)

	// 测试远程读取
	static := http.FileServer(http.Dir("./"))
	srv := httptest.NewServer(static)
	defer srv.Close()

	uri = URI(srv.URL + "/uri.go")
	exists, err = uri.Exists()
	a.NotError(err).True(exists)

	uri = URI(srv.URL + "/not-exists")
	a.NotError(err).NotEmpty(uri)
	exists, err = uri.Exists()
	a.NotError(err).False(exists)
}

func TestURI_ReadAll(t *testing.T) {
	a := assert.New(t)

	uri := FileURI("./not-exists")
	a.NotEmpty(uri)
	data, err := uri.ReadAll(nil)
	a.Error(err).Nil(data)

	data, err = uri.ReadAll(simplifiedchinese.GB18030)
	a.Error(err).Nil(data)

	uri = FileURI("./uri.go")
	a.NotEmpty(uri)
	data, err = uri.ReadAll(encoding.Nop)
	a.NotError(err).NotNil(data)

	uri = URI("./uri.go")
	data, err = uri.ReadAll(encoding.Nop)
	a.NotError(err).NotNil(data)

	uri = FileURI("./gbk.php")
	a.NotEmpty(uri)
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
	a.Nil(data).ErrorType(err, &HTTPError{})
}

func TestURI_WriteAll(t *testing.T) {
	a := assert.New(t)

	uri := URI(" :///path.php")
	a.Error(uri.WriteAll([]byte("test")))

	// 协议类型错误
	uri = URI("https:///path.php")
	a.Error(uri.WriteAll([]byte("test")))
}
