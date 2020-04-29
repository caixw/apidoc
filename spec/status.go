// SPDX-License-Identifier: MIT

package spec

import (
	"encoding/xml"
	"net/http"
	"strconv"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

// Status 表示 HTTP 状态码
type Status int

func isValidStatus(status int) bool {
	return (status >= http.StatusContinue) &&
		(status <= http.StatusNetworkAuthenticationRequired)
}

// UnmarshalXMLAttr xml.UnmarshalerAttr
func (s *Status) UnmarshalXMLAttr(attr xml.Attr) error {
	field := "/@" + attr.Name.Local

	v, err := strconv.Atoi(attr.Value)
	if err != nil {
		return fixedSyntaxError(core.Location{}, err, field)
	}

	if !isValidStatus(v) {
		return newSyntaxError(core.Location{}, field, locale.ErrInvalidFormat)
	}

	*s = Status(v)
	return nil
}

// UnmarshalXML xml.Unmarshaler
func (s *Status) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	field := "/" + start.Name.Local
	var v int
	if err := d.DecodeElement(&v, &start); err != nil {
		return fixedSyntaxError(core.Location{}, err, field)
	}

	if !isValidStatus(v) {
		return newSyntaxError(core.Location{}, field+"/status", locale.ErrInvalidFormat)
	}

	*s = Status(v)
	return nil
}

// MarshalXML xml.Marshaler
func (s Status) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(s.String(), start)
}

// MarshalXMLAttr xml.MarshalerAttr
func (s Status) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{
		Name:  name,
		Value: s.String(),
	}, nil
}

// String fmt.Stringer
func (s Status) String() string {
	return strconv.Itoa(int(s))
}
