// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"
	"testing"

	"github.com/issue9/assert"
)

var _ xml.Unmarshaler = &Example{}

func TestExample_UnmarshalXML(t *testing.T) {
	a := assert.New(t)

	obj := &Example{
		Mimetype: "application/xml",
		Content:  `<user name="name" age="18" />`,
	}
	str := `<Example mimetype="application/xml"><![CDATA[<user name="name" age="18" />]]></Example>`

	data, err := xml.Marshal(obj)
	a.NotError(err).Equal(string(data), str)

	obj1 := &Example{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1, obj)

	// 正常
	obj1 = &Example{}
	str = `<Example mimetype="json" summary="summary"><![CDATA[text]]><description>desc</description></Example>`
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1.Summary, "summary")
	a.Equal(obj1.Content, "text")
	a.Equal(obj1.Description.Text, "desc")

	// 少 mimetype
	obj1 = &Example{}
	str = `<Example url="url">desc</Example>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 少 content
	obj1 = &Example{}
	str = `<Example mimetype="json"></Example>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 语法错误
	obj1 = &Example{}
	str = `<Example mimetype="json">text`
	a.Error(xml.Unmarshal([]byte(str), obj1))
}
