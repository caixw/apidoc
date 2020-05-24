// SPDX-License-Identifier: MIT

package mock

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/assert/rest"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/ast/asttest"
	"github.com/caixw/apidoc/v7/internal/token"
)

var _ http.Handler = &mock{}

type tester struct {
	Title string
	Type  *ast.Request
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
		Type:  &ast.Request{Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeNone}}},
		JSON:  "",
		XML:   "",
	},
	{
		Title: "doc.Request{}",
		Type:  &ast.Request{},
		JSON:  "",
		XML:   "",
	},
	{
		Title: "number",
		Type: &ast.Request{
			Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeNumber}},
			Name: &ast.Attribute{Value: token.String{Value: "root"}},
		},
		JSON: "1024",
		XML:  "<root>1024</root>",
	},
	{
		Title: "enum number",
		Type: &ast.Request{
			Name: &ast.Attribute{Value: token.String{Value: "root"}},
			Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeNumber}},
			Enums: []*ast.Enum{
				{Value: &ast.Attribute{Value: token.String{Value: "1024"}}},
				{Value: &ast.Attribute{Value: token.String{Value: "1025"}}},
			},
		},
		JSON: "1024",
		XML:  `<root>1024</root>`,
	},
	{
		Title: "xml-extract",
		Type: &ast.Request{
			Name: &ast.Attribute{Value: token.String{Value: "root"}},
			Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeObject}},
			Items: []*ast.Param{
				{
					Name: &ast.Attribute{Value: token.String{Value: "id"}},
					Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeNumber}},
					XML:  ast.XML{XMLAttr: &ast.BoolAttribute{Value: ast.Bool{Value: true}}},
				},
				{
					Name: &ast.Attribute{Value: token.String{Value: "desc"}},
					Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeString}},
					XML:  ast.XML{XMLExtract: &ast.BoolAttribute{Value: ast.Bool{Value: true}}},
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
		Type: &ast.Request{
			Name: &ast.Attribute{Value: token.String{Value: "root"}},
			Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeString}},
			Enums: []*ast.Enum{
				{Value: &ast.Attribute{Value: token.String{Value: "1024"}}},
				{Value: &ast.Attribute{Value: token.String{Value: "1025"}}},
			},
		},
		JSON: `"1024"`,
		XML:  `<root>1024</root>`,
	},
	{ // array
		Title: "[bool]",
		Type: &ast.Request{
			XML:   ast.XML{XMLWrapped: &ast.Attribute{Value: token.String{Value: "root"}}},
			Name:  &ast.Attribute{Value: token.String{Value: "arr"}},
			Type:  &ast.TypeAttribute{Value: token.String{Value: ast.TypeBool}},
			Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
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
		Type: &ast.Request{
			XML:   ast.XML{XMLWrapped: &ast.Attribute{Value: token.String{Value: "root"}}},
			Name:  &ast.Attribute{Value: token.String{Value: "arr"}},
			Type:  &ast.TypeAttribute{Value: token.String{Value: ast.TypeNumber}},
			Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
			Enums: []*ast.Enum{
				{Value: &ast.Attribute{Value: token.String{Value: "1"}}},
				{Value: &ast.Attribute{Value: token.String{Value: "2"}}},
				{Value: &ast.Attribute{Value: token.String{Value: "3"}}},
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
		Type: &ast.Request{
			Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeBool}},
			Name: &ast.Attribute{Value: token.String{Value: "root"}},
		},
		JSON: "true",
		XML:  "<root>true</root>",
	},
	{ // Object
		Title: "Object with wrapped",
		Type: &ast.Request{
			Name: &ast.Attribute{Value: token.String{Value: "root"}},
			Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeObject}},
			Items: []*ast.Param{
				{
					Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeString}},
					Name: &ast.Attribute{Value: token.String{Value: "name"}},
				},
				{
					Type:  &ast.TypeAttribute{Value: token.String{Value: ast.TypeNumber}},
					Name:  &ast.Attribute{Value: token.String{Value: "num"}},
					Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
					XML:   ast.XML{XMLWrapped: &ast.Attribute{Value: token.String{Value: "nums"}}},
				},
				{
					Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeNumber}},
					Name: &ast.Attribute{Value: token.String{Value: "id"}},
					XML:  ast.XML{XMLAttr: &ast.BoolAttribute{Value: ast.Bool{Value: true}}},
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
		Type: &ast.Request{
			XML:   ast.XML{XMLWrapped: &ast.Attribute{Value: token.String{Value: "root"}}},
			Name:  &ast.Attribute{Value: token.String{Value: "user"}},
			Type:  &ast.TypeAttribute{Value: token.String{Value: ast.TypeObject}},
			Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
			Items: []*ast.Param{
				{
					Name: &ast.Attribute{Value: token.String{Value: "id"}},
					Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeNumber}},
				},
				{
					Name: &ast.Attribute{Value: token.String{Value: "name"}},
					Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeString}},
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
		Type: &ast.Request{
			Name: &ast.Attribute{Value: token.String{Value: "root"}},
			Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeObject}},
			Headers: []*ast.Param{
				{
					Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeString}},
					Name: &ast.Attribute{Value: token.String{Value: "content-type"}},
				},
				{
					Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeString}},
					Name: &ast.Attribute{Value: token.String{Value: "encoding"}},
				},
			},
			Items: []*ast.Param{
				{
					Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeObject}},
					Name: &ast.Attribute{Value: token.String{Value: "name"}},
					Items: []*ast.Param{
						{
							Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeString}},
							Name: &ast.Attribute{Value: token.String{Value: "last"}},
						},
						{
							Type:     &ast.TypeAttribute{Value: token.String{Value: ast.TypeString}},
							Name:     &ast.Attribute{Value: token.String{Value: "first"}},
							Optional: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
						},
					},
				},
				{
					Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeNumber}},
					Name: &ast.Attribute{Value: token.String{Value: "age"}},
					XML:  ast.XML{XMLAttr: &ast.BoolAttribute{Value: ast.Bool{Value: true}}},
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
		Type: &ast.Request{
			Name: &ast.Attribute{Value: token.String{Value: "root"}},
			Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeObject}},
			Items: []*ast.Param{
				{
					Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeString}},
					Name: &ast.Attribute{Value: token.String{Value: "name"}},
				},
				{
					Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeNumber}},
					Name: &ast.Attribute{Value: token.String{Value: "id"}},
					XML:  ast.XML{XMLAttr: &ast.BoolAttribute{Value: ast.Bool{Value: true}}},
				},
				{
					Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeObject}},
					Name: &ast.Attribute{Value: token.String{Value: "group"}},
					Items: []*ast.Param{
						{
							Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeString}},
							Name: &ast.Attribute{Value: token.String{Value: "name"}},
							XML:  ast.XML{XMLAttr: &ast.BoolAttribute{Value: ast.Bool{Value: true}}},
						},
						{
							Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeNumber}},
							Name: &ast.Attribute{Value: token.String{Value: "id"}},
							XML:  ast.XML{XMLAttr: &ast.BoolAttribute{Value: ast.Bool{Value: true}}},
						},
						{
							Name:  &ast.Attribute{Value: token.String{Value: "tags"}},
							Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
							Type:  &ast.TypeAttribute{Value: token.String{Value: ast.TypeObject}},
							Items: []*ast.Param{
								{
									Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeString}},
									Name: &ast.Attribute{Value: token.String{Value: "name"}},
								},
								{
									Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeNumber}},
									Name: &ast.Attribute{Value: token.String{Value: "id"}},
									XML:  ast.XML{XMLAttr: &ast.BoolAttribute{Value: ast.Bool{Value: true}}},
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
	rslt := messagetest.NewMessageHandler()
	d := &ast.APIDoc{APIDoc: &ast.APIDocVersionAttribute{Value: token.String{Value: ast.Version}}}
	d.Parse(rslt.Handler, core.Block{Data: []byte(testAPIDoc)})
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)

	rslt = messagetest.NewMessageHandler()
	mock, err := New(rslt.Handler, d, map[string]string{"test": "/test"})
	a.NotError(err).NotNil(mock)
	srv := rest.NewServer(t, mock, nil)

	// 测试路由是否正常
	srv.Get("/users").Do().Status(http.StatusBadRequest)
	srv.Post("/users", nil).Do().Status(http.StatusMethodNotAllowed)
	srv.Get("/not-found").Do().Status(http.StatusNotFound)

	srv.Post("/test/users", nil).Do().Status(http.StatusBadRequest)
	srv.Get("/test/users").Do().Status(http.StatusMethodNotAllowed)

	rslt.Handler.Stop()
	a.NotEmpty(rslt.Errors)
	srv.Close()

	rslt = messagetest.NewMessageHandler()
	mock, err = New(rslt.Handler, d, map[string]string{"test": "/test"})
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

	rslt.Handler.Stop()
	a.Empty(rslt.Errors)

	// 版本号兼容性
	rslt = messagetest.NewMessageHandler()
	mock, err = New(rslt.Handler, &ast.APIDoc{APIDoc: &ast.APIDocVersionAttribute{Value: token.String{Value: "1.0.1"}}}, nil)
	a.Error(err).Nil(mock)
	rslt.Handler.Stop()
}

func TestLoad(t *testing.T) {
	a := assert.New(t)
	rslt := messagetest.NewMessageHandler()
	mock, err := Load(rslt.Handler, "./not-exists", nil)
	rslt.Handler.Stop()
	a.Error(err).Nil(mock)

	// LoadFromPath
	rslt = messagetest.NewMessageHandler()
	mock, err = Load(rslt.Handler, asttest.URI(a), map[string]string{"admin": "/admin"})
	rslt.Handler.Stop()
	a.NotError(err).NotNil(mock)

	// loadFromURL
	static := http.FileServer(http.Dir(asttest.Dir(a)))
	srv := httptest.NewServer(static)
	defer srv.Close()

	rslt = messagetest.NewMessageHandler()
	mock, err = Load(rslt.Handler, core.URI(srv.URL+"/index.xml"), map[string]string{"admin": "/admin"})
	rslt.Handler.Stop()
	a.NotError(err).NotNil(mock)
}
