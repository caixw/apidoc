// SPDX-License-Identifier: MIT

package spec

import (
	"encoding/xml"

	"github.com/issue9/version"

	"github.com/caixw/apidoc/v6/core"
	"github.com/caixw/apidoc/v6/internal/locale"
)

// Semver 描述文档中与版本相关的信息
//
// 遵守 https://semver.org/lang/zh-CN/ 规则。
type Semver string

// UnmarshalXMLAttr xml.UnmarshalerAttr
func (v *Semver) UnmarshalXMLAttr(attr xml.Attr) error {
	if !version.SemVerValid(attr.Value) {
		field := "/@" + attr.Name.Local
		return newSyntaxError(core.Location{}, field, locale.ErrInvalidFormat)
	}

	*v = Semver(attr.Value)
	return nil
}

// UnmarshalXML xml.Unmarshaler
func (v *Semver) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	field := "/" + start.Name.Local
	var str string
	if err := d.DecodeElement(&str, &start); err != nil {
		return fixedSyntaxError(core.Location{}, err, field)
	}

	if !version.SemVerValid(str) {
		return newSyntaxError(core.Location{}, field, locale.ErrInvalidFormat)
	}

	*v = Semver(str)
	return nil
}

// MarshalXML xml.Marshaler
func (v Semver) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(string(v), start)
}

// MarshalXMLAttr xml.MarshalerAttr
func (v Semver) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{
		Name:  name,
		Value: string(v),
	}, nil
}
