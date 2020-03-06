// SPDX-License-Identifier: MIT

package mock

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/assert/rest"

	"github.com/caixw/apidoc/v6/message/messagetest"
	"github.com/caixw/apidoc/v6/spec"
	"github.com/caixw/apidoc/v6/spec/spectest"
)

var _ http.Handler = &Mock{}

type tester struct {
	Title string
	Type  *spec.Request
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
		Type:  &spec.Request{Type: spec.None},
		JSON:  "",
		XML:   "",
	},
	{
		Title: "doc.Request{}",
		Type:  &spec.Request{},
		JSON:  "",
		XML:   "",
	},
	{
		Title: "number",
		Type:  &spec.Request{Type: spec.Number, Name: "root"},
		JSON:  "1024",
		XML:   "<root>1024</root>",
	},
	{
		Title: "enum number",
		Type: &spec.Request{
			Name: "root",
			Type: spec.Number,
			Enums: []*spec.Enum{
				{Value: "1024"},
				{Value: "1025"},
			},
		},
		JSON: "1024",
		XML:  `<root>1024</root>`,
	},
	{
		Title: "xml-extract",
		Type: &spec.Request{
			Name: "root",
			Type: spec.Object,
			Items: []*spec.Param{
				{
					Name: "id",
					Type: spec.Number,
					XML:  spec.XML{XMLAttr: true},
				},
				{
					Name: "desc",
					Type: spec.String,
					XML:  spec.XML{XMLExtract: true},
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
		Type: &spec.Request{
			Name: "root",
			Type: spec.String,
			Enums: []*spec.Enum{
				{Value: "1024"},
				{Value: "1025"},
			},
		},
		JSON: `"1024"`,
		XML:  `<root>1024</root>`,
	},
	{ // array
		Title: "[bool]",
		Type: &spec.Request{
			XML:   spec.XML{XMLWrapped: "root"},
			Name:  "arr",
			Type:  spec.Bool,
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
		Type: &spec.Request{
			XML:   spec.XML{XMLWrapped: "root"},
			Name:  "arr",
			Type:  spec.Number,
			Array: true,
			Enums: []*spec.Enum{
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
		Type:  &spec.Request{Type: spec.Bool, Name: "root"},
		JSON:  "true",
		XML:   "<root>true</root>",
	},
	{ // Object
		Title: "Object with wrapped",
		Type: &spec.Request{
			Name: "root",
			Type: spec.Object,
			Items: []*spec.Param{
				{
					Type: spec.String,
					Name: "name",
				},
				{
					Type:  spec.Number,
					Name:  "num",
					Array: true,
					XML:   spec.XML{XMLWrapped: "nums"},
				},
				{
					Type: spec.Number,
					Name: "id",
					XML:  spec.XML{XMLAttr: true},
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
		Type: &spec.Request{
			XML:   spec.XML{XMLWrapped: "root"},
			Name:  "user",
			Type:  spec.Object,
			Array: true,
			Items: []*spec.Param{
				{
					Name: "id",
					Type: spec.Number,
				},
				{
					Name: "name",
					Type: spec.String,
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
		Type: &spec.Request{
			Name: "root",
			Type: spec.Object,
			Headers: []*spec.Param{
				{
					Type: spec.String,
					Name: "content-type",
				},
				{
					Type: spec.String,
					Name: "encoding",
				},
			},
			Items: []*spec.Param{
				{
					Type: spec.Object,
					Name: "name",
					Items: []*spec.Param{
						{
							Type: spec.String,
							Name: "last",
						},
						{
							Type:     spec.String,
							Name:     "first",
							Optional: true,
						},
					},
				},
				{
					Type: spec.Number,
					Name: "age",
					XML:  spec.XML{XMLAttr: true},
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
		Type: &spec.Request{
			Name: "root",
			Type: spec.Object,
			Items: []*spec.Param{
				{
					Type: spec.String,
					Name: "name",
				},
				{
					Type: spec.Number,
					Name: "id",
					XML:  spec.XML{XMLAttr: true},
				},
				{
					Type: spec.Object,
					Name: "group",
					Items: []*spec.Param{
						{
							Type: spec.String,
							Name: "name",
							XML:  spec.XML{XMLAttr: true},
						},
						{
							Type: spec.Number,
							Name: "id",
							XML:  spec.XML{XMLAttr: true},
						},
						{
							Name:  "tags",
							Array: true,
							Type:  spec.Object,
							Items: []*spec.Param{
								{
									Type: spec.String,
									Name: "name",
								},
								{
									Type: spec.Number,
									Name: "id",
									XML:  spec.XML{XMLAttr: true},
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
	d := spec.New()
	a.NotError(d.ParseBlock(&spec.Block{File: "memory.file", Data: []byte(testAPIDoc)}))

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
	mock, err = New(h, &spec.APIDoc{APIDoc: "1.0.1"}, nil)
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
	mock, err = Load(h, spectest.Path(a), map[string]string{"admin": "/admin"})
	h.Stop()
	a.NotError(err).NotNil(mock)

	// loadFromURL
	static := http.FileServer(http.Dir(spectest.Dir(a)))
	srv := httptest.NewServer(static)
	defer srv.Close()

	_, _, h = messagetest.MessageHandler()
	mock, err = Load(h, srv.URL+"/index.xml", map[string]string{"admin": "/admin"})
	h.Stop()
	a.NotError(err).NotNil(mock)
}
