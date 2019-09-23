// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"

	"github.com/issue9/version"

	"github.com/caixw/apidoc/v5/internal/locale"
)

// Version 版本号
type Version string

// UnmarshalXMLAttr xml.UnmarshalerAttr
func (v *Version) UnmarshalXMLAttr(attr xml.Attr) error {
	if !version.SemVerValid(attr.Value) {
		return newXMLSyntaxError(0, locale.ErrInvalidVersionFormat)
	}

	*v = Version(attr.Value)
	return nil
}

// UnmarshalXML xml.Unmarshaler
func (v *Version) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var str string
	if err := d.DecodeElement(&str, &start); err != nil {
		return err
	}

	if !version.SemVerValid(str) {
		return newXMLSyntaxError(0, locale.ErrInvalidVersionFormat)
	}

	*v = Version(str)
	return nil
}

// MarshalXML xml.Marshaler
func (v Version) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(string(v), start)
}

// MarshalXMLAttr xml.MarshalerAttr
func (v Version) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{
		Name:  name,
		Value: string(v),
	}, nil
}
