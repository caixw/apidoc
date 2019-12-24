// SPDX-License-Identifier: MIT

package mock

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v5/doc"
)

var xmlTestData = []*tester{
	{ // 无法获取根元素，出错
		Title: "nil",
		Type:  nil,
	},
	{ // None 表示不输出内容
		Title: "doc.None",
		Type:  &doc.Request{Type: doc.None},
	},
	{ // 等价于 None
		Title: "doc.Request{}",
		Type:  &doc.Request{},
	},
	{
		Title: "number",
		Type:  &doc.Request{Type: doc.Number, Name: "root"},
		Data:  "<root>1024</root>",
	},
	{
		Title: "enum number",
		Type: &doc.Request{
			Name: "root",
			Type: doc.Number,
			Enums: []*doc.Enum{
				{Value: "1024"},
				{Value: "1025"},
			},
		},
		Data: "<root>1024</root>",
	},
	{
		Title: "enum string",
		Type: &doc.Request{
			Name: "root",
			Type: doc.String,
			Enums: []*doc.Enum{
				{Value: "1024"},
				{Value: "1025"},
			},
		},
		Data: `<root>1024</root>`,
	},
	{ // array
		Title: "[bool]",
		Type: &doc.Request{
			Name: "root",
			Type: doc.Object,
			Items: []*doc.Param{
				{
					Array: true,
					Name:  "arr",
					Type:  doc.Bool,
				},
			},
		},
		Data: `<root>
    <arr>true</arr>
    <arr>true</arr>
    <arr>true</arr>
    <arr>true</arr>
    <arr>true</arr>
</root>`,
	},
	{
		Title: "array with enum",
		Type: &doc.Request{
			Name: "root",
			Type: doc.Object,
			Items: []*doc.Param{
				{
					Name:  "arr",
					Array: true,
					Type:  doc.Number,
					Enums: []*doc.Enum{
						{Value: "1"},
						{Value: "2"},
						{Value: "3"},
					},
				},
			},
		},
		Data: `<root>
    <arr>1</arr>
    <arr>1</arr>
    <arr>1</arr>
    <arr>1</arr>
    <arr>1</arr>
</root>`,
	},
	{ // Object
		Title: "Object",
		Type: &doc.Request{
			Name: "root",
			Type: doc.Object,
			Items: []*doc.Param{
				{
					Type: doc.String,
					Name: "name",
				},
				{
					Type: doc.Number,
					Name: "id",
					XML:  doc.XML{XMLAttr: true},
				},
			},
		},
		Data: `<root id="1024">
    <name>1024</name>
</root>`,
	},
	{ // 各类型混合
		Title: "Object with array",
		Type: &doc.Request{
			Name: "root",
			Type: doc.Object,
			Items: []*doc.Param{
				{
					Type: doc.String,
					Name: "name",
				},
				{
					Type: doc.Number,
					Name: "id",
					XML:  doc.XML{XMLAttr: true},
				},
				{
					Type: doc.Object,
					Name: "group",
					Items: []*doc.Param{
						{
							Type: doc.Number,
							Name: "id",
							XML:  doc.XML{XMLAttr: true},
						},
						{
							Type: doc.String,
							Name: "name",
							XML:  doc.XML{XMLAttr: true},
						},
						{
							Name:  "tags",
							Array: true,
							Type:  doc.Object,
							Items: []*doc.Param{
								{
									Type: doc.String,
									Name: "name",
								},
								{
									Type: doc.Number,
									Name: "id",
									XML:  doc.XML{XMLAttr: true},
								},
							},
						}, // end tags
					},
				}, // end group
			},
		},
		Data: `<root id="1024">
    <name>1024</name>
    <group id="1024" name="1024">
        <tags id="1024">
            <name>1024</name>
        </tags>
        <tags id="1024">
            <name>1024</name>
        </tags>
        <tags id="1024">
            <name>1024</name>
        </tags>
        <tags id="1024">
            <name>1024</name>
        </tags>
        <tags id="1024">
            <name>1024</name>
        </tags>
    </group>
</root>`,
	},
}

func TestValidXML(t *testing.T) {
	a := assert.New(t)

	for _, item := range xmlTestData {
		err := validXML(item.Type, []byte(item.Data))
		if item.Err {
			a.Error(err, "测试 %s 并未返回错误", item.Title)
		} else {
			a.NotError(err, "测试 %s 时返回错误 %s", item.Title, err)
		}
	}
}

func TestBuildXML(t *testing.T) {
	a := assert.New(t)

	for _, item := range xmlTestData {
		data, err := buildXML(item.Type)
		if item.Err {
			a.Error(err, "测试 %s 并未返回错误", item.Title).
				Nil(data, "测试 %s 存在返回的数据", item.Title)
		} else {
			a.NotError(err, "测试 %s 返回了错误信息 %s", item.Title, err).
				Equal(string(data), item.Data, "测试 %s 返回的数据不相等 v1:%s,v2:%s", item.Title, string(data), item.Data)
		}
	}
}

func TestXMLValidator_find(t *testing.T) {
	item := xmlTestData[len(xmlTestData)-1]

	a := assert.New(t)
	v := &xmlValidator{
		param:   item.Type.ToParam(),
		decoder: xml.NewDecoder(strings.NewReader(item.Data)),
	}

	v.names = []string{}
	p := v.find("")
	a.Nil(p)

	v.names = nil
	p = v.find("")
	a.Nil(p)

	v.names = []string{""}
	p = v.find("")
	a.Nil(p)

	v.names = []string{"root"}
	p = v.find("")
	a.NotNil(p).Equal(p.Type, doc.Object)

	v.names = []string{}
	p = v.find("root")
	a.NotNil(p).Equal(p.Type, doc.Object)

	v.names = []string{"not-exists"}
	p = v.find("")
	a.Nil(p)

	v.names = []string{"root", "group", "id"}
	p = v.find("")
	a.NotNil(p).Equal(p.Type, doc.Number)

	v.names = []string{"root", "group"}
	p = v.find("id")
	a.NotNil(p).Equal(p.Type, doc.Number)

	v.names = []string{"root", "group", "tags", "id"}
	p = v.find("")
	a.NotNil(p).Equal(p.Type, doc.Number)

	v.names = []string{"root", "group", "tags"}
	p = v.find("id")
	a.NotNil(p).Equal(p.Type, doc.Number)
}
