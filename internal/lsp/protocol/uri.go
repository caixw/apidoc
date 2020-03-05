// SPDX-License-Identifier: MIT

package protocol

import (
	"net/url"

	"github.com/caixw/apidoc/v6/internal/locale"
)

const fileScheme = "file"

// DocumentURI is are transferred as strings.
// The URI’s format is defined in http://tools.ietf.org/html/rfc3986
//
//  foo://example.com:8042/over/there?name=ferret#nose
//  \_/   \______________/\_________/ \_________/ \__/
//   |           |            |            |        |
// scheme     authority       path        query   fragment
//   |   _____________________|__
//  / \ /                        \
//  urn:example:animal:ferret:nose
//
// Many of the interfaces contain fields that correspond to the URI of a document.
// For clarity, the type of such a field is declared as a DocumentUri. Over the wire, it will still
// be transferred as a string, but this guarantees that the contents of that string can be parsed as a valid URI.
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#uri
type DocumentURI string

// File 返回 file:// 协议关联的文件路径
func (uri DocumentURI) File() (string, error) {
	u, err := url.ParseRequestURI(string(uri))
	if err != nil {
		return "", err
	}

	if u.Scheme != fileScheme {
		return "", locale.Errorf(locale.ErrInvalidURIScheme)
	}

	return u.Path, nil
}

// FileURI 根据本地文件路径构建 DocumentURI 实例
func FileURI(path string) DocumentURI {
	u := &url.URL{Scheme: fileScheme, Path: path}
	return DocumentURI(u.String())
}
