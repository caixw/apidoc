// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"
	"net/http"
	"testing"

	"github.com/issue9/assert"
)

var (
	status                     = Status(http.StatusOK)
	_      xml.Unmarshaler     = &status
	_      xml.UnmarshalerAttr = &status
	_      xml.Marshaler       = status
	_      xml.MarshalerAttr   = status
)

func TestIsValidStatus(t *testing.T) {
	a := assert.New(t)

	a.True(isValidStatus(100))
	a.True(isValidStatus(500))
	a.False(isValidStatus(1000))
}

func TestStatusXML(t *testing.T) {
	a := assert.New(t)

	type Object struct {
		XMLName struct{} `xml:"xml"`
		Attr    Status   `xml:"attr,attr"`
		Value   Status   `xml:"value"`
	}

	obj := &Object{
		Attr:  http.StatusOK,
		Value: http.StatusCreated,
	}
	str := `<xml attr="200"><value>201</value></xml>`

	data, err := xml.Marshal(obj)
	a.NotError(err).Equal(string(data), str)

	obj1 := &Object{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1, obj)
}
