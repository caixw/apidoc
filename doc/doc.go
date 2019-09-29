// SPDX-License-Identifier: MIT

// Package doc 文档格式
package doc

import (
	"bytes"
	"encoding/xml"

	xmessage "golang.org/x/text/message"

	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/internal/vars"
	"github.com/caixw/apidoc/v5/message"
)

// Doc 文档
type Doc struct {
	XMLName struct{} `xml:"apidoc"`

	// 程序的版本号
	//
	// 同时也作为文档格式的版本号。客户端可以依此值确定文档格式。
	// 仅用于输出，文档中不需要指定此值。
	APIDoc string `xml:"apidoc,attr,omitempty"`

	Version Version   `xml:"version,attr,omitempty"` // 文档的版本
	Title   string    `xml:"title"`
	Content cdata     `xml:"content"`
	Contact *Contact  `xml:"contact,omitempty"`
	License *Link     `xml:"license,omitempty"` // 版本信息
	Tags    []*Tag    `xml:"tag,omitempty"`     // 所有的标签
	Servers []*Server `xml:"server,omitempty"`
	Apis    []*API    `xml:"api,omitempty"`

	// 应用于全局的变量
	Mimetype  string     `xml:"mimetype,omitempty"` // 指定可用的 mimetype 类型
	Responses []*Request `xml:"response,omitempty"`
	Requests  []*Request `xml:"Request,omitempty"`

	file string
	line int
	data []byte
}

type cdata struct {
	Text string `xml:",cdata"`
}

// New 返回 Doc 实例
func New() *Doc {
	return &Doc{
		APIDoc: vars.Version(),
	}
}

type shadowDoc Doc

// UnmarshalXML 实现 xml.Unmarshaler 接口
//
// 返回的错误信息都为 message.SyntaxError 实例
func (doc *Doc) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	shadow := (*shadowDoc)(doc)
	if err := d.DecodeElement(shadow, &start); err != nil {
		line := bytes.Count(doc.data[:d.InputOffset()], []byte{'\n'})
		return fixedSyntaxError(err, doc.file, "doc", doc.line+line)
	}

	// Tag.Name 查重
	if key := findDupTag(shadow.Tags); key != "" {
		return message.NewLocaleError(doc.file, "doc/tag#name", doc.line, locale.ErrDuplicateValue)
	}

	// Server.Name 查重
	if key := findDupServer(shadow.Servers); key != "" {
		return message.NewLocaleError(doc.file, "doc/server#name", doc.line, locale.ErrDuplicateValue)
	}

	apis := doc.Apis
	if len(shadow.Apis) > 0 {
		apis = append(apis, shadow.Apis...)
	}

	return nil
}

// FromXML 从 XML 字符串初始化当前的实例
func (doc *Doc) FromXML(file string, line int, data []byte) error {
	doc.file = file
	doc.line = line
	doc.data = data
	return xml.Unmarshal(data, doc)
}

// Sanitize 检测内容是否合法
func (doc *Doc) Sanitize() error {
	for _, api := range doc.Apis { // 查看 API 中的标签是否都存在
		if err := api.sanitize("api"); err != nil {
			return err
		}
	}

	return nil
}

func (doc *Doc) tagExists(tag string) bool {
	for _, s := range doc.Tags {
		if s.Name == tag {
			return true
		}
	}
	return false
}

func (doc *Doc) serverExists(srv string) bool {
	for _, s := range doc.Servers {
		if s.Name == srv {
			return true
		}
	}
	return false
}

func fixedSyntaxError(err error, file, field string, line int) error {
	if serr, ok := err.(*message.SyntaxError); ok {
		serr.File = file
		serr.Line = line

		if serr.Field == "" {
			serr.Field = field
		} else {
			serr.Field = field + serr.Field
		}
		return err
	}

	return message.WithError(file, field, line, err)
}

func newSyntaxError(field string, key xmessage.Reference, val ...interface{}) error {
	return message.NewLocaleError("", field, 0, key, val...)
}
