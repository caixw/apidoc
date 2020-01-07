// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"
	"net/http"
	"strings"

	"github.com/caixw/apidoc/v6/internal/locale"
)

// Method 表示请求方法
type Method string

var validMethods = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodHead,
	http.MethodOptions,
}

func isValidMethod(method string) bool {
	method = strings.ToUpper(method)
	for _, m := range validMethods {
		if m == method {
			return true
		}
	}

	return false
}

// UnmarshalXMLAttr xml.UnmarshalerAttr
func (m *Method) UnmarshalXMLAttr(attr xml.Attr) error {

	if !isValidMethod(attr.Value) {
		field := "/@" + attr.Name.Local
		return newSyntaxError(field, locale.ErrInvalidFormat)
	}

	*m = Method(strings.ToUpper(attr.Value))
	return nil
}

// UnmarshalXML xml.Unmarshaler
func (m *Method) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	field := "/" + start.Name.Local
	var str string
	if err := d.DecodeElement(&str, &start); err != nil {
		return fixedSyntaxError(err, "", field, 0)
	}

	if !isValidMethod(str) {
		return newSyntaxError(field, locale.ErrInvalidValue)
	}

	*m = Method(strings.ToUpper(str))
	return nil
}

// MarshalXML xml.Marshaler
func (m Method) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(strings.ToUpper(string(m)), start)
}

// MarshalXMLAttr xml.MarshalerAttr
func (m Method) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{
		Name:  name,
		Value: strings.ToUpper(string(m)),
	}, nil
}
