// SPDX-License-Identifier: MIT

package mock

import (
	"net/http"

	"github.com/caixw/apidoc/v5/doc"
)

var _ http.Handler = &Mock{}

type tester struct {
	Title string
	Type  *doc.Request
	JSON  string
	XML   string
}

// 提供了测试 validJSON/buildXML 和 buildJSON/buildXML 的数据
var data = []*tester{
	{
		Title: "nil",
		Type:  nil,
		JSON:  "null",
		XML:   "",
	},
	{
		Title: "doc.None",
		Type:  &doc.Request{Type: doc.None},
		JSON:  "",
		XML:   "",
	},
	{
		Title: "doc.Request{}",
		Type:  &doc.Request{},
		JSON:  "",
		XML:   "",
	},
	{
		Title: "number",
		Type:  &doc.Request{Type: doc.Number, Name: "root"},
		JSON:  "1024",
		XML:   "<root>1024</root>",
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
		JSON: "1024",
		XML:  `<root>1024</root>`,
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
		JSON: `"1024"`,
		XML:  `<root>1024</root>`,
	},
	{ // array
		Title: "[bool]",
		Type: &doc.Request{
			XML:   doc.XML{XMLWrapped: "root"},
			Name:  "arr",
			Type:  doc.Bool,
			Array: true,
		},
		JSON: `[
    true,
    true,
    true,
    true,
    true
]`,
		XML: `<root>
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
			XML:   doc.XML{XMLWrapped: "root"},
			Name:  "arr",
			Type:  doc.Number,
			Array: true,
			Enums: []*doc.Enum{
				{Value: "1"},
				{Value: "2"},
				{Value: "3"},
			},
		},
		JSON: `[
    1,
    1,
    1,
    1,
    1
]`,
		XML: `<root>
    <arr>1</arr>
    <arr>1</arr>
    <arr>1</arr>
    <arr>1</arr>
    <arr>1</arr>
</root>`,
	},
	{
		Title: "bool",
		Type:  &doc.Request{Type: doc.Bool, Name: "root"},
		JSON:  "true",
		XML:   "<root>true</root>",
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
					Type:  doc.Number,
					Name:  "num",
					Array: true,
					XML:   doc.XML{XMLWrapped: "nums"},
				},
				{
					Type: doc.Number,
					Name: "id",
					XML:  doc.XML{XMLAttr: true},
				},
			},
		},
		JSON: `{
    "name": "1024",
    "num": [
        1024,
        1024,
        1024,
        1024,
        1024
    ],
    "id": 1024
}`,
		XML: `<root id="1024">
    <name>1024</name>
    <nums>
        <num>1024</num>
        <num>1024</num>
        <num>1024</num>
        <num>1024</num>
        <num>1024</num>
    </nums>
</root>`,
	},

	// NOTE: 部分测试用例单独引用了该项内容。 必须保持在倒数第二的位置。
	{
		Title: "object with item",
		Type: &doc.Request{
			Name: "root",
			Type: doc.Object,
			Headers: []*doc.Param{
				{
					Type: doc.String,
					Name: "content-type",
				},
				{
					Type: doc.String,
					Name: "encoding",
				},
			},
			Items: []*doc.Param{
				{
					Type: doc.Object,
					Name: "name",
					Items: []*doc.Param{
						{
							Type: doc.String,
							Name: "last",
						},
						{
							Type:     doc.String,
							Name:     "first",
							Optional: true,
						},
					},
				},
				{
					Type: doc.Number,
					Name: "age",
					XML:  doc.XML{XMLAttr: true},
				},
			},
		},
		JSON: `{
    "name": {
        "last": "1024",
        "first": "1024"
    },
    "age": 1024
}`,
		XML: `<root age="1024">
    <name>
        <last>1024</last>
        <first>1024</first>
    </name>
</root>`,
	},

	// NOTE: 部分测试用例单独引用了该项内容。 必须保持在倒数第一的位置。
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
							Type: doc.String,
							Name: "name",
							XML:  doc.XML{XMLAttr: true},
						},
						{
							Type: doc.Number,
							Name: "id",
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
		JSON: `{
    "name": "1024",
    "id": 1024,
    "group": {
        "name": "1024",
        "id": 1024,
        "tags": [
            {
                "name": "1024",
                "id": 1024
            },
            {
                "name": "1024",
                "id": 1024
            },
            {
                "name": "1024",
                "id": 1024
            },
            {
                "name": "1024",
                "id": 1024
            },
            {
                "name": "1024",
                "id": 1024
            }
        ]
    }
}`,
		XML: `<root id="1024">
    <name>1024</name>
    <group name="1024" id="1024">
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
