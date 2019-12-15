// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"
	"testing"

	"github.com/issue9/assert"
)

var (
	_ xml.Marshaler   = Richtext{}
	_ xml.Unmarshaler = &Richtext{}
)

func TestRichtext_Marshal(t *testing.T) {
	a := assert.New(t)

	type Object struct {
		XMLName struct{} `xml:"xml"`
		Value   Richtext `xml:"value"`
	}

	obj := &Object{
		Value: Richtext{Text: "<a>test</a>"},
	}
	str := `<xml><value type="markdown"><![CDATA[<a>test</a>]]></value></xml>`

	data, err := xml.Marshal(obj)
	a.NotError(err).Equal(string(data), str)

	obj1 := &Object{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.NotEqual(obj1.Value, obj.Value)
	obj.Value.Type = RichtextTypeMarkdown
	a.Equal(obj1.Value, obj.Value)

	// 空值
	obj = &Object{
		Value: Richtext{},
	}
	str = `<xml></xml>`

	data, err = xml.Marshal(obj)
	a.NotError(err).Equal(string(data), str)

	obj1 = &Object{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1, obj)
}

func TestShadowRichtext_Marshal(t *testing.T) {
	a := assert.New(t)

	type Object struct {
		XMLName struct{}       `xml:"xml"`
		Value   shadowRichtext `xml:"value"`
	}

	obj := &Object{
		Value: shadowRichtext{Text: "<p>cdata</p>"},
	}
	str := `<xml><value><![CDATA[<p>cdata</p>]]></value></xml>`

	data, err := xml.Marshal(obj)
	a.NotError(err).Equal(string(data), str)

	obj1 := &Object{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1.Value, obj.Value)

	// 空值
	obj = &Object{
		Value: shadowRichtext{},
	}
	str = `<xml><value></value></xml>`

	data, err = xml.Marshal(obj)
	a.NotError(err).Equal(string(data), str)

	obj1 = &Object{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1, obj)
}
