// SPDX-License-Identifier: MIT

// Package doc 文档格式
package doc

import (
	"bytes"
	"encoding/xml"
	"sort"
	"time"

	xmessage "golang.org/x/text/message"

	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/internal/vars"
	"github.com/caixw/apidoc/v5/message"
)

const createdFormat = time.RFC3339

// Doc 文档
type Doc struct {
	XMLName struct{} `xml:"apidoc"`

	// 程序的版本号
	//
	// 同时也作为文档格式的版本号。客户端可以依此值确定文档格式。
	// 仅用于输出，文档中不需要指定此值。
	APIDoc string `xml:"apidoc,attr,omitempty"`

	// 文档内容的区域信息
	// 如果存在此值，客户端应该尽量根据此值显示相应的界面语言。
	Lang string `xml:"lang,attr,omitempty"`

	// 文档的图标
	//
	// 如果采用默认的 xsl 转换，会替换掉页面上的图标和 favicon 图标
	Logo string `xml:"logo,attr,omitempty"`

	Created     string    `xml:"created,attr,omitempty"` // 文档的生成时间
	Version     Version   `xml:"version,attr,omitempty"` // 文档的版本
	Title       string    `xml:"title"`
	Description CDATA     `xml:"description"`
	Contact     *Contact  `xml:"contact,omitempty"`
	License     *Link     `xml:"license,omitempty"` // 版本信息
	Tags        []*Tag    `xml:"tag,omitempty"`     // 所有的标签
	Servers     []*Server `xml:"server,omitempty"`
	Apis        []*API    `xml:"api,omitempty"`

	// 表示所有 API 都有可能返回的内容
	Responses []*Request `xml:"response,omitempty"`

	file string
	line int
	data []byte
}

// New 返回 Doc 实例
func New() *Doc {
	return &Doc{
		APIDoc:  vars.Version(),
		Created: time.Now().Format(createdFormat),
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
	// doc.Apis 是多线程导入的，无法保证其顺序，
	// 此处可以保证输出内容是按一定顺序排列的。
	sort.SliceStable(doc.Apis, func(i, j int) bool {
		ii := doc.Apis[i]
		jj := doc.Apis[j]

		if ii.Path.Path == jj.Path.Path {
			return ii.Method < jj.Method
		}
		return ii.Path.Path < jj.Path.Path
	})

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
