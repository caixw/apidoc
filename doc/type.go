// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"
	"strings"

	"github.com/caixw/apidoc/v5/internal/locale"
)

// Type 表示参数类型
type Type string

// 表示支持的各种数据类型
const (
	None   Type = ""
	Bool        = "bool"
	Object      = "object"
	Number      = "number"
	String      = "string"
)

func parseType(val string) (Type, error) {
	// Deprecated: 下个版本移除对 none 的支持
	if val == "none" {
		return None, nil
	}

	val = strings.ToLower(val)
	switch Type(val) {
	case None, Bool, Object, Number, String:
		return Type(val), nil
	default:
		return None, locale.Errorf(locale.ErrInvalidFormat)
	}
}

// UnmarshalXMLAttr xml.UnmarshalerAttr
func (t *Type) UnmarshalXMLAttr(attr xml.Attr) error {
	v, err := parseType(attr.Value)
	if err != nil {
		return err
	}

	*t = v
	return nil
}

// UnmarshalXML xml.Unmarshaler
func (t *Type) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	field := "/" + start.Name.Local
	var str string
	if err := d.DecodeElement(&str, &start); err != nil {
		return fixedSyntaxError(err, "", field, 0)
	}

	v, err := parseType(str)
	if err != nil {
		return fixedSyntaxError(err, "", field+"/type", 0)
	}

	*t = v
	return nil
}

// MarshalXML xml.Marshaler
func (t Type) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	v, err := t.fmtString()
	if err != nil {
		return err
	}

	return e.EncodeElement(v, start)
}

// MarshalXMLAttr xml.MarshalerAttr
func (t Type) MarshalXMLAttr(name xml.Name) (attr xml.Attr, err error) {
	attr = xml.Attr{Name: name}

	val, err := t.fmtString()
	if err != nil {
		return attr, err
	}

	attr.Value = val
	return attr, nil
}

func (t Type) fmtString() (string, error) {
	valid := t == None ||
		t == Bool ||
		t == Object ||
		t == Number ||
		t == String

	if valid {
		return string(t), nil
	}
	return "", locale.Errorf(locale.ErrInvalidValue)
}
