// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"
	"strings"

	"github.com/caixw/apidoc/v5/errors"
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

// 数组作为一个特殊的类型，表示方式与其它类型不同。
// 由一个专门的标记位标记该属性。
const arrayFlag Type = 0b1000_0000

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

// String fmt.Stringer
func (t Type) String() string {
	flag := (arrayFlag & t) == arrayFlag
	if flag {
		t = t & (^arrayFlag)
	}

	ret, found := typeStringMap[t]
	if !found {
		return "none"
	}

	if flag {
		ret = "array." + ret
	}
	return ret
}

func parseType(val string) (Type, error) {
	val = strings.ToLower(val)
	var flag Type

	dotIndex := strings.IndexByte(val, '.')
	if dotIndex > -1 {
		prefix := val[:dotIndex]
		if prefix != "array" {
			return None, errors.New("TODO", "TODO", 0, locale.ErrInvalidTypePrefix, prefix)
		}

		flag = arrayFlag
		val = val[dotIndex+1:]
	}

	t, found := stringTypeMap[val]
	if !found {
		return None, errors.New("TODO", "TODO", 0, locale.ErrInvalidType, val)
	}
	return flag + t, nil
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
	var str string
	if err := d.DecodeElement(&str, &start); err != nil {
		return err
	}

	v, err := parseType(str)
	if err != nil {
		return err
	}

	*t = v
	return nil
}

// MarshalXML xml.Marshaler
func (t Type) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	v := t.String()

	return e.EncodeElement(v, start)
}

// MarshalXMLAttr xml.MarshalerAttr
func (t Type) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{
		Name:  name,
		Value: t.String(),
	}, nil
}
