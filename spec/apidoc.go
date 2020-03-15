// SPDX-License-Identifier: MIT

package spec

import (
	"bytes"
	"encoding/xml"
	"sort"
	"time"

	xmessage "golang.org/x/text/message"

	"github.com/caixw/apidoc/v6/internal/locale"
	"github.com/caixw/apidoc/v6/internal/vars"
	"github.com/caixw/apidoc/v6/message"
)

const createdFormat = time.RFC3339

// APIDoc 对应 apidoc 元素
type APIDoc struct {
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
	Version     Semver    `xml:"version,attr,omitempty"` // 文档的版本
	Title       string    `xml:"title"`
	Description Richtext  `xml:"description,omitempty"`
	Contact     *Contact  `xml:"contact,omitempty"`
	License     *Link     `xml:"license,omitempty"` // 版本信息
	Tags        []*Tag    `xml:"tag,omitempty"`     // 所有的标签
	Servers     []*Server `xml:"server,omitempty"`
	Apis        []*API    `xml:"api,omitempty"`

	// 表示所有 API 都有可能返回的内容
	Responses []*Request `xml:"response,omitempty"`

	// 表示所有接口都支持的文档类型
	Mimetypes []string `xml:"mimetype"`

	Block *Block `xml:"-"`
}

// Valid 验证文档内容的正确性
func Valid(content []byte) error {
	return NewAPIDoc().fromXML(&Block{Data: content})
}

// NewAPIDoc 返回 APIDoc 实例
func NewAPIDoc() *APIDoc {
	return &APIDoc{
		APIDoc:  vars.Version(),
		Created: time.Now().Format(createdFormat),
	}
}

type shadowDoc APIDoc

// UnmarshalXML 实现 xml.Unmarshaler 接口
//
// 返回的错误信息都为 message.SyntaxError 实例
func (doc *APIDoc) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	shadow := (*shadowDoc)(doc)
	if err := d.DecodeElement(shadow, &start); err != nil {
		line := bytes.Count(doc.Block.Data[:d.InputOffset()], []byte{'\n'})
		return fixedSyntaxError(err, doc.Block.File, "apidoc", doc.Block.Range.Start.Line+line)
	}

	if shadow.Title == "" {
		return doc.Block.localeError("apidoc/title", locale.ErrRequired)
	}

	// Tag.Name 查重
	if findDupTag(shadow.Tags) != "" {
		return doc.Block.localeError("apidoc/tag/@name", locale.ErrDuplicateValue)
	}

	// Server.Name 查重
	if findDupServer(shadow.Servers) != "" {
		return doc.Block.localeError("apidoc/server/@name", locale.ErrDuplicateValue)
	}

	if len(shadow.Mimetypes) == 0 {
		return doc.Block.localeError("apidoc/mimetype", locale.ErrRequired)
	}

	// 操作 clone 进行比较，不影响原文档的排序
	clone := make([]string, len(shadow.Mimetypes))
	copy(clone, shadow.Mimetypes)
	sort.Strings(clone)
	for index := 1; index < len(clone); index++ {
		if clone[index] == clone[index-1] {
			return doc.Block.localeError("apidoc/mimetype", locale.ErrDuplicateValue)
		}
	}

	apis := doc.Apis
	if len(shadow.Apis) > 0 {
		apis = append(apis, shadow.Apis...)
	}

	return nil
}

// Sanitize 检测内容是否合法
func (doc *APIDoc) Sanitize() error {
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

func (doc *APIDoc) tagExists(tag string) bool {
	for _, s := range doc.Tags {
		if s.Name == tag {
			return true
		}
	}
	return false
}

func (doc *APIDoc) serverExists(srv string) bool {
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
