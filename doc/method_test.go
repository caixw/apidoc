// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"
	"testing"

	"github.com/issue9/assert"
)

var (
	m                     = Method("get")
	_ xml.Unmarshaler     = &v
	_ xml.UnmarshalerAttr = &v
	_ xml.Marshaler       = v
	_ xml.MarshalerAttr   = v
)

func TestIsValidMethod(t *testing.T) {
	a := assert.New(t)

	a.True(isValidMethod("GET"))
	a.True(isValidMethod("get"))
	a.False(isValidMethod("not-exists"))
}

func TestMethodXML(t *testing.T) {
	a := assert.New(t)

	type Object struct {
		XMLName struct{} `xml:"xml"`
		Attr    Method   `xml:"attr,attr"`
		Value   Method   `xml:"value"`
	}

	obj := &Object{
		Attr:  "GET",
		Value: "POST",
	}
	str := `<xml attr="GET"><value>POST</value></xml>`

	data, err := xml.Marshal(obj)
	a.NotError(err).Equal(string(data), str)

	obj1 := &Object{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1, obj)
}
