// SPDX-License-Identifier: MIT

package spec

import (
	"encoding/xml"
	"testing"

	"github.com/issue9/assert"
)

var (
	v                     = Version("1.1.1")
	_ xml.Unmarshaler     = &v
	_ xml.UnmarshalerAttr = &v
	_ xml.Marshaler       = v
	_ xml.MarshalerAttr   = v
)

func TestVersionXML(t *testing.T) {
	a := assert.New(t)

	type Object struct {
		XMLName struct{} `xml:"xml"`
		Attr    Version  `xml:"attr,attr"`
		Value   Version  `xml:"value"`
	}

	obj := &Object{
		Attr:  "1.0.1",
		Value: "1.0.2",
	}
	str := `<xml attr="1.0.1"><value>1.0.2</value></xml>`

	data, err := xml.Marshal(obj)
	a.NotError(err).Equal(string(data), str)

	obj1 := &Object{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1, obj)

	// 无效的状态
	str = `<xml attr="x.0.1" />`
	a.Error(xml.Unmarshal([]byte(str), obj1))
	str = `<xml attr="0.1.0"><value>.100</value></xml>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 语法错误
	str = `<xml attr="0.1.0"><value>`
	a.Error(xml.Unmarshal([]byte(str), obj1))
}
