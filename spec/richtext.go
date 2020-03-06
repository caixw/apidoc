// SPDX-License-Identifier: MIT

package spec

import "encoding/xml"

// 富文本可用的类型
const (
	RichtextTypeHTML     = "html"
	RichtextTypeMarkdown = "markdown"
)

// Richtext 富文本内容
type Richtext struct {
	Type string `xml:"type,attr,omitempty"` // 文档类型，可以是 html 或是 markdown
	Text string `xml:",cdata"`
}

type shadowRichtext Richtext

// MarshalXML xml.Marshaler
func (text Richtext) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if text.Text == "" {
		return nil
	}

	if text.Type == "" {
		text.Type = RichtextTypeMarkdown
	}
	shadow := shadowRichtext(text)
	return e.EncodeElement(shadow, start)
}

// UnmarshalXML 实现 xml.Unmarshaler
func (text *Richtext) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	field := "/" + start.Name.Local

	shadow := (*shadowRichtext)(text)
	if err := d.DecodeElement(shadow, &start); err != nil {
		return fixedSyntaxError(err, "", field, 0)
	}

	if shadow.Type == "" {
		shadow.Type = RichtextTypeMarkdown
	}

	return nil
}
