// SPDX-License-Identifier: MIT

package mock

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/assert/rest"

	"github.com/caixw/apidoc/v6/doc"
	"github.com/caixw/apidoc/v6/doc/doctest"
	"github.com/caixw/apidoc/v6/input"
	"github.com/caixw/apidoc/v6/message/messagetest"
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
		Title: "xml-extract",
		Type: &doc.Request{
			Name: "root",
			Type: doc.Object,
			Items: []*doc.Param{
				{
					Name: "id",
					Type: doc.Number,
					XML:  doc.XML{XMLAttr: true},
				},
				{
					Name: "desc",
					Type: doc.String,
					XML:  doc.XML{XMLExtract: true},
				},
			},
		},
		JSON: `{
    "id": 1024,
    "desc": "1024"
}`,
		XML: `<root id="1024">1024</root>`,
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
		Title: "Object with wrapped",
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

	{
		Title: "object array",
		Type: &doc.Request{
			XML:   doc.XML{XMLWrapped: "root"},
			Name:  "user",
			Type:  doc.Object,
			Array: true,
			Items: []*doc.Param{
				{
					Name: "id",
					Type: doc.Number,
				},
				{
					Name: "name",
					Type: doc.String,
				},
			},
		},
		JSON: `[
    {
        "id": 1024,
        "name": "1024"
    },
    {
        "id": 1024,
        "name": "1024"
    },
    {
        "id": 1024,
        "name": "1024"
    },
    {
        "id": 1024,
        "name": "1024"
    },
    {
        "id": 1024,
        "name": "1024"
    }
]`,
		XML: `<root>
    <user>
        <id>1024</id>
        <name>1024</name>
    </user>
    <user>
        <id>1024</id>
        <name>1024</name>
    </user>
    <user>
        <id>1024</id>
        <name>1024</name>
    </user>
    <user>
        <id>1024</id>
        <name>1024</name>
    </user>
    <user>
        <id>1024</id>
        <name>1024</name>
    </user>
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

const testAPIDoc = `<apidoc version="1.0.1">
	<title>title</title>
	<server url="https://example.com" name="test" summary="test summary" />
	<mimetype>application/json</mimetype>
	<mimetype>application/xml</mimetype>

	<api method="GET" summary="get users">
		<path path="/users" />
		<response type="object" array="true" xml-wrapped="root" name="user" status="200">
			<param name="id" type="number" summary="id summary" />
			<param name="name" type="string" summary="name summary" />
		</response>
	</api>
	<api method="post" summary="post user">
		<server>test</server>
		<path path="/users" />
		<request type="object" array="true" xml-wrapped="root" name="user">
			<param name="id" type="number" summary="id summary" />
			<param name="name" type="string" summary="name summary" />
		</request>
		<response status="201">
			<header type="string" name="location" summary="新资源的地址" />
		</response>
	</api>
</apidoc>`

func TestNew(t *testing.T) {
	a := assert.New(t)
	d := doc.New()
	a.NotError(d.ParseBlock(&input.Block{File: "memory.file", Data: []byte(testAPIDoc)}))

	erro, _, h := messagetest.MessageHandler()
	mock, err := New(h, d, map[string]string{"test": "/test"})
	a.NotError(err).NotNil(mock)
	srv := rest.NewServer(t, mock, nil)

	// 测试路由是否正常
	srv.Get("/users").Do().Status(http.StatusBadRequest)
	srv.Post("/users", nil).Do().Status(http.StatusMethodNotAllowed)
	srv.Get("/not-found").Do().Status(http.StatusNotFound)

	srv.Post("/test/users", nil).Do().Status(http.StatusBadRequest)
	srv.Get("/test/users").Do().Status(http.StatusMethodNotAllowed)

	h.Stop()
	a.NotEmpty(erro.String())
	srv.Close()

	erro, _, h = messagetest.MessageHandler()
	mock, err = New(h, d, map[string]string{"test": "/test"})
	a.NotError(err).NotNil(mock)
	srv = rest.NewServer(t, mock, nil)

	//
	srv.Post("/test/users", nil).
		Header("accept", "application/json").
		Header("content-type", "application/xml").
		Body([]byte(`<root>
    <user>
        <id>1</id>
        <name>n</name>
    </user>
</root>`)).
		Do().
		Status(http.StatusCreated).
		Header("content-type", "application/json").
		BodyEmpty()

	h.Stop()
	a.Empty(erro.String())

	// 版本号兼容性
	_, _, h = messagetest.MessageHandler()
	mock, err = New(h, &doc.Doc{APIDoc: "1.0.1"}, nil)
	a.Error(err).Nil(mock)
	h.Stop()
}

func TestLoad(t *testing.T) {
	a := assert.New(t)
	_, _, h := messagetest.MessageHandler()
	mock, err := Load(h, "./not-exists", nil)
	h.Stop()
	a.Error(err).Nil(mock)

	// LoadFromPath
	_, _, h = messagetest.MessageHandler()
	mock, err = Load(h, doctest.Path(a), map[string]string{"admin": "/admin"})
	h.Stop()
	a.NotError(err).NotNil(mock)

	// loadFromURL
	static := http.FileServer(http.Dir(doctest.Dir(a)))
	srv := httptest.NewServer(static)
	defer srv.Close()

	_, _, h = messagetest.MessageHandler()
	mock, err = Load(h, srv.URL+"/index.xml", map[string]string{"admin": "/admin"})
	h.Stop()
	a.NotError(err).NotNil(mock)
}
