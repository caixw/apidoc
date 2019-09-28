// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"
	"testing"

	"github.com/issue9/assert"
)

var (
	_ xml.Unmarshaler = &Header{}
)

func TestHeader_UnmarshalXML(t *testing.T) {
	a := assert.New(t)

	obj := &Header{
		Name:        "text",
		Description: "<a>desc</a>",
	}
	str := `<Header name="text"><![CDATA[<a>desc</a>]]></Header>`

	data, err := xml.Marshal(obj)
	a.NotError(err).Equal(string(data), str)

	obj1 := &Header{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1, obj)

	// 正常
	obj1 = &Header{}
	str = `<Header name="url" deprecated="1.1.1">text</Header>`
	a.NotError(xml.Unmarshal([]byte(str), obj1))

	// 少 name
	obj1 = &Header{}
	str = `<Header url="url">desc</Header>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 少 description
	obj1 = &Header{}
	str = `<Header name="v1"></Header>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 语法错误
	obj1 = &Header{}
	str = `<Header name="url" deprecated="x.1.1">text</Header>`
	a.Error(xml.Unmarshal([]byte(str), obj1))
}
