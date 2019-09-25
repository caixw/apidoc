// SPDX-License-Identifier: MIT

// Package doc 文档格式
package doc

import (
	"bytes"
	"encoding/xml"

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
	Content string    `xml:"content"`
	Contact *Contact  `xml:"contact"`
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

// New 返回 Doc 实例
func New() *Doc {
	return &Doc{
		APIDoc: vars.Version(),
	}
}

type shadowDoc Doc

// UnmarshalXML xml.Unmarshaler
func (doc *Doc) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var shadow shadowDoc
	if err := d.DecodeElement(&shadow, &start); err != nil {
		line := bytes.Count(doc.data[:d.InputOffset()], []byte{'\n'})
		return message.WithError(doc.file, "", doc.line+line, err)
	}

	// Tag.Name 查重
	if key := findDupTag(shadow.Tags); key != "" {
		return locale.Errorf(locale.ErrDuplicateTag, key)
	}

	// Server.Name 查重
	if key := findDupServer(shadow.Servers); key != "" {
		return locale.Errorf(locale.ErrDuplicateValue, key)
	}

	apis := doc.Apis
	if len(shadow.Apis) > 0 {
		apis = append(apis, shadow.Apis...)
	}

	*doc = Doc(shadow)
	doc.Apis = apis
	doc.APIDoc = vars.Version() // 读取的时候，忽略客户端指定的 apidoc
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
		if err := api.sanitize(); err != nil {
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
