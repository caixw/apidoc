// SPDX-License-Identifier: MIT

package doc

import "encoding/xml"

// Richtext 富文本内容
type Richtext struct {
	Text string `xml:",innerxml"`
}

type shadowRichtext Richtext

// CDATA 定义了一个使用 CDATA 的元素
//
// 该元素在没有内容时，不会输出任何内容，包括元素本身。
// 如果需要在无内容时输出元素本身的标签，可以使用 shadowCDATA。
type CDATA struct {
	Text string `xml:",cdata"`
}

type shadowCDATA CDATA

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

// MarshalXML xml.Marshaler
func (cd CDATA) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if cd.Text == "" {
		return nil
	}

	shadow := shadowCDATA(cd)
	return e.EncodeElement(shadow, start)
}

// String 返回 CDATA 的文本内容
func (cd CDATA) String() string {
	return cd.Text
}
