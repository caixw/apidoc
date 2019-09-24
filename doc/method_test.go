// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"
	"testing"

	"github.com/issue9/assert"
)

var (
	method                     = Method("get")
	_      xml.Unmarshaler     = &method
	_      xml.UnmarshalerAttr = &method
	_      xml.Marshaler       = method
	_      xml.MarshalerAttr   = method
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

	// 不存在的 method
	str = `<xml attr="not-exists" />`
	a.Error(xml.Unmarshal([]byte(str), obj1))
	str = `<xml attr="get"><value>not-exists</value></xml>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 语法错误
	str = `<xml attr="GET"><value>`
	a.Error(xml.Unmarshal([]byte(str), obj1))
}
