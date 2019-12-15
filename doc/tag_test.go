// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"
	"testing"

	"github.com/issue9/assert"
)

var (
	_ xml.Unmarshaler = &Tag{}
	_ xml.Unmarshaler = &Server{}
)

func TestTag_UnmarshalXML(t *testing.T) {
	a := assert.New(t)

	obj := &Tag{
		Name:       "tag1",
		Title:      "test",
		Deprecated: "1.1.1",
	}
	str := `<Tag name="tag1" title="test" deprecated="1.1.1"></Tag>`

	data, err := xml.Marshal(obj)
	a.NotError(err).Equal(string(data), str)

	obj1 := &Tag{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1, obj)

	// 少 name
	str = `<Tag>test</Tag>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 少 title
	str = `<Tag name="tag1"></Tag>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 语法错误
	str = `<Tag name="tag1" deprecated="x.0.1">desc</Tag>`
	a.Error(xml.Unmarshal([]byte(str), obj1))
}

func TestServer_UnmarshalXML(t *testing.T) {
	a := assert.New(t)

	obj := &Server{
		Name:        "srv1",
		URL:         "https://api.example.com/srv1",
		Deprecated:  "1.1.1",
		Description: Richtext{Text: "<a>test</a>"},
	}
	str := `<Server name="srv1" url="https://api.example.com/srv1" deprecated="1.1.1"><description type="markdown"><![CDATA[<a>test</a>]]></description></Server>`

	data, err := xml.Marshal(obj)
	a.NotError(err).Equal(string(data), str)

	obj1 := &Server{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1.Description.Text, obj.Description.Text).
		NotEqual(obj1.Description.Type, obj.Description.Type) // type 在 marshal 是会能默认址

	// 正常，带 description
	obj1 = &Server{}
	str = `<Server name="tag1" url="https://example.com"><description><![CDATA[text]]></description></Server>`
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1.Description.Text, "text")

	// 少 name
	str = `<Server />`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 少 url
	str = `<Server name="tag1" />`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 少 description 或是 summary
	str = `<Server name="tag1" />`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 语法错误
	str = `<Server name="tag1" deprecated="x.0.1" summary="desc" />`
	a.Error(xml.Unmarshal([]byte(str), obj1))
}

func TestFindDupString(t *testing.T) {
	a := assert.New(t)

	a.Equal(findDupString([]string{"k1", "k2", "K2"}), "")
	a.Equal(findDupString([]string{"k2", "k1", "k2"}), "k2")
}
