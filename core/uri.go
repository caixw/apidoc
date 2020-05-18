// SPDX-License-Identifier: MIT

package core

import (
	"io/ioutil"
	"net/http"
	"os"
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

	separator = "://"
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
//
// NOTE: 并非完整的 URI 实现，仅作为了 file:// 和 http:// 支持，
// 也提供对 windows 路径的支持。
type URI string

// FileURI 根据本地文件路径构建 URI 实例
//
// 如果已经存在协议，则不作任何改变返回。
func FileURI(path string) URI {
	if index := strings.Index(path, separator); index > -1 {
		return URI(path)
	}
	return URI(SchemeFile + separator + path)
}

// File 返回 file:// 协议关联的文件路径
func (uri URI) File() (string, error) {
	if scheme, path := uri.Parse(); scheme == SchemeFile || scheme == "" {
		return path, nil
	}
	return "", locale.NewError(locale.ErrInvalidURIScheme)
}

func (uri URI) String() string {
	return string(uri)
}

// Append 追加 path 至 URI 生成新的 URI
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
	scheme, path := uri.Parse()
	switch scheme {
	case SchemeFile, "":
		return utils.FileExists(path), nil
	case SchemeHTTP, SchemeHTTPS:
		return remoteFileIsExists(string(uri))
	default:
		return false, locale.NewError(locale.ErrInvalidURIScheme)
	}
}

// ReadAll 以 enc 编码读取 uri 的内容
//
// 目前仅支持 file、http 和 https 协议
func (uri URI) ReadAll(enc encoding.Encoding) ([]byte, error) {
	scheme, path := uri.Parse()
	switch scheme {
	case SchemeFile, "":
		return readLocalFile(path, enc)
	case SchemeHTTP, SchemeHTTPS:
		return readRemoteFile(string(uri), enc)
	default:
		return nil, locale.NewError(locale.ErrInvalidURIScheme)
	}
}

// WriteAll 写入内容至 uri
func (uri URI) WriteAll(data []byte) error {
	if scheme, path := uri.Parse(); scheme == SchemeFile || scheme == "" {
		return ioutil.WriteFile(path, data, os.ModePerm)
	}
	return locale.NewError(locale.ErrInvalidURIScheme)
}

// Parse 分析 uri，获取其各个部分的内容
func (uri URI) Parse() (schema, path string) {
	uris := string(uri)
	index := strings.Index(uris, separator)
	if index == -1 {
		return "", uris
	}
	return uris[:index], uris[index+len(separator):]
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
		return nil, NewHTTPError(resp.StatusCode, locale.ErrReadRemoteFile, url, resp.StatusCode)
	}

	if enc == nil || enc == encoding.Nop {
		return ioutil.ReadAll(resp.Body)
	}
	reader := transform.NewReader(resp.Body, enc.NewDecoder())
	return ioutil.ReadAll(reader)
}
