// SPDX-License-Identifier: MIT

package core

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/issue9/utils"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"

	"github.com/caixw/apidoc/v6/internal/locale"
)

const (
	fileScheme  = "file"
	httpScheme  = "http"
	httpsScheme = "https"
)

// URI 定义 URI
//
// http://tools.ietf.org/html/rfc3986
//
//    foo://example.com:8042/over/there?name=ferret#nose
//    \_/   \______________/\_________/ \_________/ \__/
//     |           |            |            |        |
//  scheme     authority       path        query   fragment
//     |   _____________________|__
//    / \ /                        \
//    urn:example:animal:ferret:nose
type URI string

// FileURI 根据本地文件路径构建 URI 实例
//
// NOTE: 如果 path 不会绝对路径，会被转换成绝对路径，
// 如果仅需要一个表示相对的路径的 URI 类型，可以采用：
//  URI(path)
// 的方式直接将字符串转换成 URI 类型。
func FileURI(path string) (URI, error) {
	if !filepath.IsAbs(path) {
		p, err := filepath.Abs(path)
		if err != nil {
			return "", err
		}
		path = p
	}

	u := &url.URL{Scheme: fileScheme, Path: path}
	return URI(u.String()), nil
}

// File 返回 file:// 协议关联的文件路径
func (uri URI) File() (string, error) {
	if uri.isNoScheme() {
		return string(uri), nil
	}

	u, err := url.ParseRequestURI(string(uri))
	if err != nil {
		return "", err
	}

	if u.Scheme != fileScheme && u.Scheme != "" {
		return "", locale.Errorf(locale.ErrInvalidURIScheme)
	}

	return u.Path, nil
}

func (uri URI) String() string {
	return string(uri)
}

// Append 追加 path 至 URI，生成新的 URI。
func (uri URI) Append(path string) URI {
	return uri + URI(path)
}

// Exists 判断 uri 指向的内容是否存在
//
// 如果是非本地文件，通过 http 的状态码是否为 400 以内加以判断。
func (uri URI) Exists() (bool, error) {
	if uri.isNoScheme() {
		return localFileIsExists(string(uri)), nil
	}

	u, err := url.ParseRequestURI(string(uri))
	if err != nil {
		return false, err
	}

	switch u.Scheme {
	case fileScheme, "":
		return localFileIsExists(u.Path), nil
	case httpScheme, httpsScheme:
		return remoteFileIsExists(string(uri))
	default:
		return false, locale.Errorf(locale.ErrInvalidURIScheme)
	}
}

// ReadAll 以 enc 编码读取 uri 的内容
//
// 目前仅支持 file、http 和 https 协议
func (uri URI) ReadAll(enc encoding.Encoding) ([]byte, error) {
	if uri.isNoScheme() {
		return readLocalFile(string(uri), enc)
	}

	u, err := url.ParseRequestURI(string(uri))
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case fileScheme, "":
		return readLocalFile(u.Path, enc)
	case httpScheme, httpsScheme:
		return readRemoteFile(string(uri), enc)
	default:
		return nil, locale.Errorf(locale.ErrInvalidURIScheme)
	}
}

func (uri URI) isNoScheme() bool {
	str := string(uri)

	if str == "" || str[0] == '.' || str[0] == '/' || str[0] == os.PathSeparator {
		return true
	}

	if index := strings.IndexByte(str, ':'); index > 0 {
		// 可能是 windows 的 c:\path 格式
		return len(str) > index+1 && str[index+1] == '\\'
	}

	return true
}

func localFileIsExists(path string) bool {
	return utils.FileExists(path)
}

func remoteFileIsExists(url string) (bool, error) {
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	return resp.StatusCode < 400, nil
}

// 以指定的编码方式读取本地文件内容
func readLocalFile(path string, enc encoding.Encoding) ([]byte, error) {
	if enc == nil || enc == encoding.Nop {
		return ioutil.ReadFile(path)
	}

	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	reader := transform.NewReader(r, enc.NewDecoder())
	return ioutil.ReadAll(reader)
}

// 以指定的编码方式读取远程文件内容
func readRemoteFile(url string, enc encoding.Encoding) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 300 {
		return nil, locale.Errorf(locale.ErrReadRemoteFile, url, resp.StatusCode)
	}

	if enc == nil || enc == encoding.Nop {
		return ioutil.ReadAll(resp.Body)
	}
	reader := transform.NewReader(resp.Body, enc.NewDecoder())
	return ioutil.ReadAll(reader)
}
