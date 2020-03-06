// SPDX-License-Identifier: MIT

package spec

import (
	"encoding/xml"
	"testing"

	"github.com/issue9/assert"
)

var (
	none                     = None
	_    xml.Unmarshaler     = &none
	_    xml.UnmarshalerAttr = &none
	_    xml.Marshaler       = none
	_    xml.MarshalerAttr   = none
)

func TestTypeXML(t *testing.T) {
	a := assert.New(t)

	type Object struct {
		XMLName struct{} `xml:"type"`
		Attr    Type     `xml:"attr,attr"`
		Value   Type     `xml:"value"`
	}

	obj := &Object{
		Attr:  String,
		Value: Number,
	}
	str := `<type attr="string"><value>number</value></type>`

	data, err := xml.Marshal(obj)
	a.NotError(err)
	a.Equal(string(data), str)

	obj1 := &Object{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj, obj1)

	// 无效的类型
	str = `<type attr="not-exists" />`
	a.Error(xml.Unmarshal([]byte(str), obj1))
	str = `<type attr="string"><value>not-exists</value></type>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 语法错误
	str = `<type attr="string"><value>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// fmt
	obj.Value = Type("100")
	data, err = xml.Marshal(obj)
	a.Error(err).Nil(data)
}

func TestType_None(t *testing.T) {
	a := assert.New(t)

	type Object struct {
		XMLName struct{} `xml:"type"`
		Attr    Type     `xml:"attr,attr,omitempty"`
		Value   Type     `xml:"value"`
	}

	obj := &Object{
		Attr:  None,
		Value: Number,
	}
	str := `<type><value>number</value></type>`

	data, err := xml.Marshal(obj)
	a.NotError(err)
	a.Equal(string(data), str)

	obj1 := &Object{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj, obj1)
}
