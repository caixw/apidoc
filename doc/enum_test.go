// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"
	"testing"

	"github.com/issue9/assert"
)

var (
	_ xml.Unmarshaler = &Enum{}
)

func TestEnum_UnmarshalXML(t *testing.T) {
	a := assert.New(t)

	obj := &Enum{
		Value:       "text",
		Description: Richtext{Text: "<a><p>desc</p></a>"},
	}
	str := `<Enum value="text"><description type="markdown"><![CDATA[<a><p>desc</p></a>]]></description></Enum>`

	data, err := xml.Marshal(obj)
	a.NotError(err).Equal(string(data), str)

	obj1 := &Enum{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1.Value, obj.Value).
		NotEmpty(obj1.Description.Type).
		Empty(obj.Description.Type).
		Equal(obj.Description.Text, obj1.Description.Text)

	// 正常
	obj1 = &Enum{}
	str = `<Enum value="url" deprecated="1.1.1"><description>text</description></Enum>`
	a.NotError(xml.Unmarshal([]byte(str), obj1))

	// 少 value
	obj1 = &Enum{}
	str = `<Enum url="url"><description>desc</description></Enum>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 少 description 和 summary
	obj1 = &Enum{}
	str = `<Enum value="v1"></Enum>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 语法错误
	obj1 = &Enum{}
	str = `<Enum value="url" deprecated="x.1.1"><description>text</description></Enum>`
	a.Error(xml.Unmarshal([]byte(str), obj1))
}
