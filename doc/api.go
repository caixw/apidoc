// SPDX-License-Identifier: MIT

package doc

import (
	"bytes"
	"encoding/xml"

	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/message"
)

// API 表示 <api> 顶层元素
//  <api method="GET" version="1.1.1" id="get-user">
//      <path path="/users/{id}">
//          <param name="id" type="number" summary="summary" />
//      </path>
//      <tag>tag1</tag>
//      <server>admin</server>
//      ...
//  </api>
type API struct {
	XMLName     struct{}   `xml:"api"`
	Version     Version    `xml:"version,attr,omitempty"`
	Method      Method     `xml:"method,attr"`
	ID          string     `xml:"id,attr,omitempty"`
	Path        *Path      `xml:"path"`
	Summary     string     `xml:"summary,attr"`
	Description string     `xml:"description,omitempty"`
	Requests    []*Request `xml:"request"`
	Responses   []*Request `xml:"response"`
	Callback    *Callback  `xml:"callback,omitempty"`
	Deprecated  Version    `xml:"deprecated,attr,omitempty"`

	Tags    []string `xml:"tag,omitempty"`
	Servers []string `xml:"server,omitempty"`

	line int
	file string
	data []byte
	doc  *Doc
}

// NewAPI 从 data 中解析新的 API 对象
func (doc *Doc) NewAPI(file string, line int, data []byte) error {
	api := &API{
		file: file,
		line: line,
		data: data,
	}
	if err := xml.Unmarshal(data, api); err != nil {
		return err
	}

	api.doc = doc
	doc.Apis = append(doc.Apis, api)
	return nil
}

type shadowAPI API

// UnmarshalXML xml.Unmarshaler
func (api *API) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// TODO 如果 api.file 为空，则采用 doc.file
	// 这样可以正常解析 apidoc 和 api 合在一起的文档

	var obj shadowAPI
	if err := d.DecodeElement(&obj, &start); err != nil {
		line := bytes.Count(api.data[:d.InputOffset()], []byte{'\n'})
		return message.WithError(api.file, "", api.line+line, err)
	}

	*api = API(obj)
	return nil
}

// 检测和修复 api 对象，无法修复返回错误。
//
// NOTE: 需要保证 doc 已经初始化
func (api *API) sanitize() error {
	for _, tag := range api.Tags {
		if !api.doc.tagExists(tag) {
			return message.NewLocaleError(api.file, "tag", api.line, locale.ErrInvalidValue)
		}
	}

	for _, srv := range api.Servers {
		if !api.doc.serverExists(srv) {
			return message.NewLocaleError(api.file, "server", api.line, locale.ErrInvalidValue)
		}
	}

	// TODO ref

	return nil
}
