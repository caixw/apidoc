// SPDX-License-Identifier: MIT

package doc

import "encoding/xml"

// Richtext 富文本内容
type Richtext struct {
	Text string `xml:",innerxml"`
}

type shadowRichtext Richtext

// MarshalXML xml.Marshaler
func (text Richtext) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if text.Text == "" {
		return nil
	}

	shadow := shadowRichtext(text)
	return e.EncodeElement(shadow, start)
}

// String 返回 Richtext 的文本内容
func (text Richtext) String() string {
	return text.Text
}
