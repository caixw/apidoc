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

// 确保 typeStringMap 和 stringTypeMap 相同
func TestTypeMap(t *testing.T) {
	a := assert.New(t)

	a.Equal(len(stringTypeMap), len(typeStringMap))

	for k, v := range stringTypeMap {
		a.Equal(typeStringMap[v], k)
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
}
