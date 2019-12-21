// SPDX-License-Identifier: MIT

package mock

import "github.com/caixw/apidoc/v5/doc"

var xmlTestData = []*jsonTester{
	{ // 无法获取根元素，出错
		Title: "nil",
		Type:  nil,
		Err:   true,
	},
	{ // None 表示不输出内容
		Title: "doc.None",
		Type:  &doc.Request{Type: doc.None},
		Data:  "",
	},
	{ // 等价于 None
		Title: "doc.Request{}",
		Type:  &doc.Request{},
		Data:  "",
	},
	{ // 未指定 name
		Title: "doc.Request{}",
		Type:  &doc.Request{Type: doc.Object},
		Err:   true,
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
	{ // 顶层元素不能是数组
		Title: "array at root element",
		Type: &doc.Request{
			Name:  "root",
			Type:  doc.Bool,
			Array: true,
		},
		Err: true,
	},
	{ // 数组不能是属性值
		Title: "[bool]",
		Type: &doc.Request{
			Name: "root",
			Type: doc.Object,
			Items: []*doc.Param{
				{
					Array: true,
					Name:  "arr",
					Type:  doc.Bool,
					Attr:  true,
				},
			},
		},
		Err: true,
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
		Data: `<root><arr>true</arr><arr>true</arr><arr>true</arr><arr>true</arr><arr>true</arr></root>`,
	},
	{
		Title: "array with enum",
		Type: &doc.Request{
			Name: "root",
			Type: doc.Number,
			Items: []*doc.Param{
				{
					Name:  "arr",
					Array: true,
					Enums: []*doc.Enum{
						{Value: "1"},
						{Value: "2"},
						{Value: "3"},
					},
				},
			},
		},
		Data: `<root><arr>1</arr><arr>1</arr><arr>1</arr><arr>1</arr><arr>1</arr></root>`,
	},
	{
		Title: "bool",
		Type:  &doc.Request{Type: doc.Bool},
		Data:  "true",
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
					Attr: true,
				},
			},
		},
		Data: `<root id="1024"><name>1024</name></root>`,
	},
}
