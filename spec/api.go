// SPDX-License-Identifier: MIT

package spec

import (
	"bytes"
	"encoding/xml"

	"github.com/caixw/apidoc/v6/core"
	"github.com/caixw/apidoc/v6/internal/locale"
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
	Version     Semver     `xml:"version,attr,omitempty"`
	Method      Method     `xml:"method,attr"`
	ID          string     `xml:"id,attr,omitempty"`
	Path        *Path      `xml:"path"`
	Summary     string     `xml:"summary,attr,omitempty"`
	Description Richtext   `xml:"description,omitempty"`
	Requests    []*Request `xml:"request,omitempty"` // 不同的 mimetype 可能会定义不同
	Responses   []*Request `xml:"response,omitempty"`
	Callback    *Callback  `xml:"callback,omitempty"`
	Deprecated  Semver     `xml:"deprecated,attr,omitempty"`
	Headers     []*Param   `xml:"header,omitempty"`

	Tags    []string `xml:"tag,omitempty"`
	Servers []string `xml:"server,omitempty"`

	Block *core.Block `xml:"-"`
	doc   *APIDoc
}

type shadowAPI API

// UnmarshalXML 实现 xml.Unmarshaler 接口
func (api *API) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	field := "/" + start.Name.Local
	shadow := (*shadowAPI)(api)
	if err := d.DecodeElement(shadow, &start); err != nil {
		// API 可能是嵌套在 apidoc 里的一个子标签。
		// 如果是子标签，则不应该有 doc 变量，也不需要构建错误信息。
		if api.doc == nil {
			return fixedSyntaxError(core.Location{}, err, field)
		}

		line := bytes.Count(api.Block.Data[:d.InputOffset()], []byte{'\n'})
		loc := core.Location{
			URI: api.Block.Location.URI,
			Range: core.Range{
				Start: core.Position{
					Line: api.Block.Location.Range.Start.Line + line,
				},
			},
		}
		return fixedSyntaxError(loc, err, field)
	}

	// 报头不能为 object
	for _, header := range shadow.Headers {
		if header.Type == Object {
			err := locale.Errorf(locale.ErrInvalidValue)
			field = field + "/header[" + header.Name + "].type"
			return fixedSyntaxError(api.Block.Location, err, field)
		}
	}

	return nil
}

// 检测和修复 api 对象，无法修复返回错误。
//
// NOTE: 需要保证 doc 已经初始化
func (api *API) sanitize(field string) error {
	if api.doc == nil {
		panic("api.doc 未获取正确的值")
	}

	for _, tag := range api.Tags {
		if !api.doc.tagExists(tag) {
			return core.NewLocaleError(api.Block.Location, field+"/tag/@name", locale.ErrInvalidValue)
		}
	}

	if len(api.Servers) == 0 {
		return core.NewLocaleError(api.Block.Location, field+"/server", locale.ErrRequired)
	}

	for _, srv := range api.Servers {
		if !api.doc.serverExists(srv) {
			return core.NewLocaleError(api.Block.Location, field+"/server/@name", locale.ErrInvalidValue)
		}
	}

	return nil
}
