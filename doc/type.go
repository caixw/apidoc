// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"
	"strconv"
	"strings"

	"github.com/caixw/apidoc/v5/internal/locale"
)

// Type 表示数据类型
type Type uint8

// 表示支持的各种数据类型
const (
	None Type = iota
	Bool
	Object
	Number
	String
)

var (
	typeStringMap = map[Type]string{
		None:   "none",
		Bool:   "bool",
		Object: "object",
		Number: "number",
		String: "string",
	}

	stringTypeMap = map[string]Type{
		"none":   None,
		"bool":   Bool,
		"object": Object,
		"number": Number,
		"string": String,
	}
)

func parseType(val string) (Type, error) {
	val = strings.ToLower(val)
	if t, found := stringTypeMap[val]; found {
		return t, nil
	}

	return None, locale.Errorf(locale.ErrInvalidFormat)

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
	name := "/" + start.Name.Local
	var str string
	if err := d.DecodeElement(&str, &start); err != nil {
		return fixedSyntaxError(err, "", name, 0)
	}

	v, err := parseType(str)
	if err != nil {
		return fixedSyntaxError(err, "", name+"/type", 0)
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

	attr.Value, err = t.fmtString()
	return attr, err
}

// fmtString
func (t Type) fmtString() (string, error) {
	if v, found := typeStringMap[t]; found {
		return v, nil
	}

	return "", locale.Errorf(locale.ErrInvalidType, strconv.Itoa(int(t)))
}
