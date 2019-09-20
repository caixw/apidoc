// SPDX-License-Identifier: MIT

package doc

import (
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
	Description Richtext   `xml:"description,omitempty"`
	Requests    []*Request `xml:"request"`
	Responses   []*Request `xml:"response"`
	Callback    *Callback  `xml:"callback,omitempty"`
	Deprecated  Version    `xml:"deprecated,attr,omitempty"`

	Tags    []string `xml:"tag,omitempty"`
	Servers []string `xml:"server,omitempty"`

	line int
	file string
	doc  *Doc
}

// Request 请求内容
type Request struct {
	Param
	Status      int        `xml:"status,attr"`
	Mimetype    string     `xml:"mimetype,attr"`
	Examples    []*Example `xml:"example,omitempty"`
	Headers     []*Header  `xml:"header,omitempty"`
	Description Richtext   `xml:"description,omitempty"`
}

// Path 路径信息
//  <path path="/users/{id}">
//      <param name="id" type="number" summary="summary" />
//      <query name="page" type="number" summary="page" default="1" />
//  </path>
type Path struct {
	Path      string   `xml:"path,attr"`
	Params    []*Param `xml:"param,omitempty"`
	Queries   []*Param `xml:"query,omitempty"`
	Reference string   `xml:"ref,attr,omitempty"`
}

// Param 表示参数类型
type Param struct {
	Name        string   `xml:"name,attr"`
	Type        Type     `xml:"type,attr"`
	Deprecated  Version  `xml:"deprecated,attr,omitempty"`
	Default     string   `xml:"default,attr,omitempty"`
	Required    bool     `xml:"required,attr,omitempty"`
	Enums       []*Enum  `xml:"enum,omitempty"`
	Array       bool     `xml:"array,attr,omitempty"`
	Items       []*Param `xml:"param,omitempty"`
	Reference   string   `xml:"ref,attr,omitempty"`
	Summary     string   `xml:"summary,attr,omitempty"`
	Description Richtext `xml:"description,omitempty"`
}

// IsEnum 是否为一个枚举类型
func (p *Param) IsEnum() bool {
	return len(p.Enums) > 0
}

// Enum 表示枚举值
type Enum struct {
	Deprecated  Version  `xml:"deprecated,attr,omitempty"`
	Value       string   `xml:"value,attr"`
	Description Richtext `xml:",innerxml"`
}

// Header 报头信息
type Header struct {
	Name        string   `xml:"name,attr"`
	Description Richtext `xml:"description,omitempty"`
	Deprecated  Version  `xml:"deprecated,attr,omitempty"`
}

// Example 示例代码
type Example struct {
	Description Richtext `xml:",innerxml"`
	Mimetype    string   `xml:"mimetype,attr"`
	Content     Richtext `xml:",innerxml"`
}

// Callback 回调函数的定义
type Callback struct {
	Param
	Schema      string     `xml:"schema,attr"` // http 或是 https
	Description Richtext   `xml:"description,omitempty"`
	Mimetype    string     `xml:"mimetype,attr"`
	Examples    []*Example `xml:"example,omitempty"`
	Headers     []*Header  `xml:"header,omitempty"`
	Method      Method     `xml:"method,attr"`
	Queries     []*Param   `xml:"queries,omitempty"` // 查询参数
	Requests    []*Request `xml:"requests,omitempty"`
	Deprecated  Version    `xml:"deprecated,attr,omitempty"`
	Reference   string     `xml:"ref,attr,omitempty"`

	// 对回调的返回要求
	Responses []*Request `xml:"response,omitempty"`
}

// NewAPI 返回新的 API 实例
func (doc *Doc) NewAPI(file string, line int) *API {
	api := &API{
		line: line,
		file: file,
		doc:  doc,
	}
	doc.Apis = append(doc.Apis, api)

	return api
}

// FromXML 从 XML 字符串初始化当前的实例
func (api *API) FromXML(data []byte) error {
	err := xml.Unmarshal(data, api)
	if err == nil {
		return nil
	}

	if serr, ok := err.(*message.SyntaxError); ok {
		serr.File = api.file
		serr.Line = api.line

		return err
	}
	return message.WithError(api.file, "", api.line, err)
}

// 检测和修复 api 对象，无法修复返回错误。
//
// NOTE: 需要保证 doc 已经初始化
func (api *API) sanitize() error {
	for _, tag := range api.Tags {
		if !api.doc.tagExists(tag) {
			return message.NewError(api.file, "tag", api.line, locale.ErrInvalidValue)
		}
	}

	for _, srv := range api.Servers {
		if !api.doc.serverExists(srv) {
			return message.NewError(api.file, "server", api.line, locale.ErrInvalidValue)
		}
	}

	// TODO ref

	return nil
}
