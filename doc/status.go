// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"
	"net/http"
	"strconv"

	"github.com/caixw/apidoc/v5/internal/locale"
)

// Status 表示 HTTP 状态码
type Status int

func isValidStatus(status int) bool {
	return (status >= http.StatusContinue) &&
		(status <= http.StatusNetworkAuthenticationRequired)
}

// UnmarshalXMLAttr xml.UnmarshalerAttr
func (s *Status) UnmarshalXMLAttr(attr xml.Attr) error {
	v, err := strconv.Atoi(attr.Value)
	if err != nil {
		return err
	}
	if !isValidStatus(v) {
		return locale.Errorf(locale.ErrInvalidStatus, v)
	}

	*s = Status(v)
	return nil
}

// UnmarshalXML xml.Unmarshaler
func (s *Status) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v int
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}

	if !isValidStatus(v) {
		return locale.Errorf(locale.ErrInvalidStatus, v)
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
