// SPDX-License-Identifier: MIT

package doc

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

// 确保 stringTypeMap 的内容都正确存在于 typeStringMap
func TestTypeMap(t *testing.T) {
	a := assert.New(t)

	a.True(len(stringTypeMap) > len(typeStringMap))

	for k, v := range typeStringMap {
		a.Equal(stringTypeMap[v], k)
	}
}

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

	// int 可以正常转换成 number
	str = `<type attr="int" />`
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1.Attr, Number)

	// 无效的类型
	str = `<type attr="not-exists" />`
	a.Error(xml.Unmarshal([]byte(str), obj1))
	str = `<type attr="string"><value>not-exists</value></type>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 语法错误
	str = `<type attr="string"><value>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// fmt
	obj.Value = Type(100)
	data, err = xml.Marshal(obj)
	a.Error(err).Nil(data)
}
