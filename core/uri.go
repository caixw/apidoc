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

	"github.com/caixw/apidoc/v7/internal/locale"
)

// 目前 URI 支持的协议
const (
	SchemeFile  = "file"
	SchemeHTTP  = "http"
	SchemeHTTPS = "https"
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
//
// 如果是本地相对路径，也可以直接使用 `./path/file` 的形式表示，
// 不需要指定协议。
type URI string

// FileURI 根据本地文件路径构建 URI 实例
//
// NOTE: 如果 path 不是绝对路径，会被转换成绝对路径，
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

	u := &url.URL{Scheme: SchemeFile, Path: path}
	return URI(u.String()), nil
}

// File 返回 file:// 协议关联的文件路径
func (uri URI) File() (string, error) {
	u, err := uri.Parse()
	if err != nil {
		return "", err
	}

	if u.Scheme != SchemeFile && u.Scheme != "" {
		return "", locale.Errorf(locale.ErrInvalidURIScheme)
	}

	return u.Path, nil
}

func (uri URI) String() string {
	return string(uri)
}

// Append 追加 path 至 URI，生成新的 URI。
func (uri URI) Append(path string) URI {
	if path == "" {
		return uri
	}

	str := string(uri)
	last := str[len(str)-1]

	if isPathSeparator(last) {
		if isPathSeparator(path[0]) {
			path = path[1:]
		}
	} else {
		if !isPathSeparator(path[0]) {
			path = "/" + path
		}
	}

	return uri + URI(path)
}

func isPathSeparator(b byte) bool {
	return b == '/' || b == os.PathSeparator
}

// Exists 判断 uri 指向的内容是否存在
//
// 如果是非本地文件，通过 http 的状态码是否为 400 以内加以判断。
func (uri URI) Exists() (bool, error) {
	u, err := uri.Parse()
	if err != nil {
		return false, err
	}

	switch u.Scheme {
	case SchemeFile:
		return localFileIsExists(u.Path), nil
	case SchemeHTTP, SchemeHTTPS:
		return remoteFileIsExists(string(uri))
	default:
		return false, locale.Errorf(locale.ErrInvalidURIScheme)
	}
}

// ReadAll 以 enc 编码读取 uri 的内容
//
// 目前仅支持 file、http 和 https 协议
func (uri URI) ReadAll(enc encoding.Encoding) ([]byte, error) {
	u, err := uri.Parse()
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case SchemeFile:
		return readLocalFile(u.Path, enc)
	case SchemeHTTP, SchemeHTTPS:
		return readRemoteFile(string(uri), enc)
	default:
		return nil, locale.Errorf(locale.ErrInvalidURIScheme)
	}
}

// WriteAll 写入内容至 uri
func (uri URI) WriteAll(data []byte) error {
	u, err := uri.Parse()
	if err != nil {
		return err
	}

	if u.Scheme != SchemeFile {
		return locale.Errorf(locale.ErrInvalidURIScheme)
	}

	return ioutil.WriteFile(u.Path, data, os.ModePerm)
}

// Parse 分析 uri，获取其各个部分的内容
func (uri URI) Parse() (*url.URL, error) {
	if uri.IsNoScheme() {
		return &url.URL{
			Scheme: SchemeFile,
			Path:   string(uri),
		}, nil
	}

	return url.ParseRequestURI(string(uri))
}

// IsNoScheme 是否不包含协议部分
func (uri URI) IsNoScheme() bool {
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
		msg := locale.Sprintf(locale.ErrReadRemoteFile, url, resp.StatusCode)
		return nil, NewHTTPError(resp.StatusCode, msg)
	}

	if enc == nil || enc == encoding.Nop {
		return ioutil.ReadAll(resp.Body)
	}
	reader := transform.NewReader(resp.Body, enc.NewDecoder())
	return ioutil.ReadAll(reader)
}
